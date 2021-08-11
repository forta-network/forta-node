package registry

import (
	"context"
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/contracts"
	"OpenZeppelin/fortify-node/ethereum"
	"OpenZeppelin/fortify-node/services"
	"OpenZeppelin/fortify-node/services/registry/regtypes"
	"OpenZeppelin/fortify-node/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// RegistryService listens to the agent scanner list changes so the node can stay in sync.
type RegistryService struct {
	cfg            config.Config
	scannerAddress common.Address
	msgClient      clients.MessageClient

	contract   ContractRegistryCaller
	ipfsClient IPFSClient
	ethClient  EthClient

	agentsConfigs []*config.AgentConfig
	done          chan struct{}
	version       string
	sem           *semaphore.Weighted
}

// ContractRegistryCaller calls the contract registry.
type ContractRegistryCaller interface {
	AgentLength(opts *bind.CallOpts, scanner common.Address) (*big.Int, error)
	AgentAt(opts *bind.CallOpts, scanner common.Address, index *big.Int) ([32]byte, *big.Int, bool, string, bool, error)
	GetAgentListHash(opts *bind.CallOpts, scanner common.Address) ([32]byte, error)
}

// IPFSClient interacts with an IPFS Gateway.
type IPFSClient interface {
	GetAgentFile(cid string) (*regtypes.AgentFile, error)
}

// EthClient interacts with the Ethereum API.
type EthClient interface {
	ethereum.Client
}

// New creates a new service.
func New(cfg config.Config, scannerAddress common.Address, msgClient clients.MessageClient) services.Service {
	var ipfsURL string
	if cfg.Registry.IPFSGateway != nil {
		ipfsURL = *cfg.Registry.IPFSGateway
	} else {
		ipfsURL = config.DefaultIPFSGateway
	}

	return &RegistryService{
		cfg:            cfg,
		scannerAddress: scannerAddress,
		msgClient:      msgClient,
		ipfsClient:     &ipfsClient{ipfsURL},
		done:           make(chan struct{}),
	}
}

// Start starts the registry service.
func (rs *RegistryService) Start() error {
	log.Infof("Starting %s", rs.Name())
	rpcClient, err := rpc.Dial(rs.cfg.Registry.Ethereum.JsonRpcUrl)
	if err != nil {
		return err
	}
	log.Infof("Creating Caller: %s", rs.Name())

	// used for getting the latest block number so that we can query consistent state
	ethClient, err := ethereum.NewStreamEthClient(context.Background(), rs.cfg.Registry.Ethereum.JsonRpcUrl)
	if err != nil {
		return err
	}
	rs.ethClient = ethClient

	// init registry contract
	rs.contract, err = contracts.NewScannerRegistryCaller(common.HexToAddress(rs.cfg.Registry.ContractAddress), ethclient.NewClient(rpcClient))
	if err != nil {
		return fmt.Errorf("failed to create the agent registry caller: %v", err)
	}
	rs.sem = semaphore.NewWeighted(1)
	return rs.start()
}

func (rs *RegistryService) start() error {
	go func() {
		//TODO: possibly make this configurable, but 15s per block is normal
		ticker := time.NewTicker(15 * time.Second)
		for {
			if err := rs.publishLatestAgents(); err != nil {
				log.Errorf("failed to publish the latest agents: %v", err)
			}
			<-ticker.C
		}
	}()

	return nil
}

func (rs *RegistryService) publishLatestAgents() error {
	// only allow one executor at a time, even if slow
	if rs.sem.TryAcquire(1) {
		defer rs.sem.Release(1)
		// opts is nil so we get the latest scanner list version
		version, err := rs.contract.GetAgentListHash(nil, rs.scannerAddress)
		if err != nil {
			return fmt.Errorf("failed to get the scanner list agents version: %v", err)
		}
		versionStr := string(version[:])
		// if versions change, then get the full list of agents
		if rs.version == "" || rs.version != versionStr {
			log.Infof("registry: agent version changed %s->%s", rs.version, versionStr)
			rs.version = versionStr
			rs.agentsConfigs, err = rs.getLatestAgents()
			if err != nil {
				return fmt.Errorf("failed to get latest agents: %v", err)
			}
			rs.msgClient.Publish(messaging.SubjectAgentsVersionsLatest, rs.agentsConfigs)
		} else {
			log.Info("registry: no agent changes detected")
		}
	}
	return nil
}

func (rs *RegistryService) getLatestAgents() ([]*config.AgentConfig, error) {
	var agentConfigs []*config.AgentConfig
	blk, err := rs.ethClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the block for agents: %v", err)
	}

	num, err := utils.HexToBigInt(blk.Number)
	opts := &bind.CallOpts{
		BlockNumber: num,
	}

	lengthBig, err := rs.contract.AgentLength(opts, rs.scannerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get the scanner list agents length: %v", err)
	}
	log.Infof("registry: getLatestAgents(%s) = %s", rs.scannerAddress.Hex(), lengthBig.Text(10))
	// TODO: If we are going to get 100s of agents, we probably need to batch the calls here.
	length := int(lengthBig.Int64())
	for i := 0; i < length; i++ {
		agentID, _, _, agentRef, disabled, err := rs.contract.AgentAt(opts, rs.scannerAddress, big.NewInt(int64(i)))
		if err != nil {
			return nil, fmt.Errorf("failed to get agent at index '%d' in scanner list '%s': %v", i, rs.scannerAddress.String(), err)
		}
		// if agent dev disables agent, this will prevent it from running
		if !disabled {
			agentCfg, err := rs.makeAgentConfig(agentID, agentRef)
			if err != nil {
				return nil, fmt.Errorf("failed to make agent config: %v", err)
			}
			agentConfigs = append(agentConfigs, &agentCfg)
		}
	}

	return agentConfigs, nil
}

func (rs *RegistryService) makeAgentConfig(agentID [32]byte, ref string) (agentCfg config.AgentConfig, err error) {
	agentCfg.ID = (common.Hash)(agentID).String()
	if len(ref) == 0 {
		return
	}

	var (
		agentData *regtypes.AgentFile
	)
	for i := 0; i < 10; i++ {
		agentData, err = rs.ipfsClient.GetAgentFile(ref)
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("failed to load the agent file using ipfs ref: %v", err)
		return
	}

	var ok bool
	agentCfg.Image, ok = utils.ValidateImageRef(rs.cfg.Registry.ContainerRegistry, agentData.Manifest.ImageReference)
	if !ok {
		log.Warnf("invalid agent reference - skipping: %s", agentCfg.Image)
	}

	return
}

// Stop stops the registry service.
func (rs *RegistryService) Stop() error {
	return nil
}

// Name returns the name of the service.
func (rs *RegistryService) Name() string {
	return "RegistryService"
}
