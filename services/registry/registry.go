package registry

import (
	"fmt"
	"math/big"
	"sync"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/contracts"
	"OpenZeppelin/fortify-node/feeds"
	"OpenZeppelin/fortify-node/services"
	"OpenZeppelin/fortify-node/services/registry/regtypes"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type agentUpdate struct {
	Config     config.AgentConfig
	IsCreation bool
	IsUpdate   bool
	IsRemoval  bool
}

// RegistryService listens to the agent pool changes so the node can stay in sync.
type RegistryService struct {
	cfg       config.Config
	poolID    common.Hash
	msgClient clients.MessageClient
	logFeed   feeds.LogFeed

	contract    ContractRegistryCaller
	logUnpacker LogUnpacker
	ipfsClient  IPFSClient

	agentsConfigs  []*config.AgentConfig
	agentUpdates   chan *agentUpdate
	agentUpdatesWg sync.WaitGroup
	done           chan struct{}
}

// LogUnpacker unpacks agent events from logs.
type LogUnpacker interface {
	UnpackAgentRegistryAgentAdded(log *types.Log) (*contracts.AgentRegistryAgentAdded, error)
	UnpackAgentRegistryAgentUpdated(log *types.Log) (*contracts.AgentRegistryAgentUpdated, error)
	UnpackAgentRegistryAgentRemoved(log *types.Log) (*contracts.AgentRegistryAgentRemoved, error)
}

// ContractRegistryCaller calls the contract registry.
type ContractRegistryCaller interface {
	AgentLength(opts *bind.CallOpts, _poolId [32]byte) (*big.Int, error)
	AgentAt(opts *bind.CallOpts, _poolId [32]byte, index *big.Int) ([32]byte, string, error)
}

// IPFSClient interacts with an IPFS Gateway.
type IPFSClient interface {
	GetAgentFile(cid string) (*regtypes.AgentFile, error)
}

// New creates a new service.
func New(cfg config.Config, msgClient clients.MessageClient, logFeed feeds.LogFeed) services.Service {
	var ipfsURL string
	if cfg.Registry.IPFSGateway != nil {
		ipfsURL = *cfg.Registry.IPFSGateway
	} else {
		ipfsURL = config.DefaultIPFSGateway
	}

	return &RegistryService{
		cfg:          cfg,
		poolID:       common.HexToHash(cfg.Registry.PoolID),
		msgClient:    msgClient,
		logFeed:      logFeed,
		logUnpacker:  contracts.NewAgentLogUnpacker(common.HexToAddress(cfg.Registry.ContractAddress)),
		ipfsClient:   &ipfsClient{ipfsURL},
		agentUpdates: make(chan *agentUpdate, 100),
		done:         make(chan struct{}),
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
	rs.contract, err = contracts.NewAgentRegistryCaller(common.HexToAddress(rs.cfg.Registry.ContractAddress), ethclient.NewClient(rpcClient))
	if err != nil {
		return fmt.Errorf("failed to create the agent registry caller: %v", err)
	}
	return rs.start()
}

func (rs *RegistryService) start() error {
	// Start detecting and buffering events.
	go func() {
		log.Info("registry: ForEachTransaction")
		err := rs.logFeed.ForEachLog(rs.detectAgentEvents)
		if err != nil {
			panic(err)
		}
	}()

	// Start to handle agent updates but wait until initialization is complete.
	rs.agentUpdatesWg.Add(1)

	go func() {
		log.Info("registry: listenToAgentUpdates")
		rs.listenToAgentUpdates()
		log.Warn("registry: listenToAgentUpdates is DONE!")
	}()

	if err := rs.publishLatestAgents(); err != nil {
		return fmt.Errorf("failed to publish the latest agents: %v", err)
	}
	log.Info("registry: publishLatestAgents complete")

	// Continue by processing buffered updates.
	rs.agentUpdatesWg.Done()

	return nil
}

func (rs *RegistryService) publishLatestAgents() (err error) {
	rs.agentsConfigs, err = rs.getLatestAgents()
	if err != nil {
		return
	}
	rs.msgClient.Publish(messaging.SubjectAgentsVersionsLatest, rs.agentsConfigs)
	return
}

func (rs *RegistryService) getLatestAgents() ([]*config.AgentConfig, error) {

	lengthBig, err := rs.contract.AgentLength(nil, rs.poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the pool agents length: %v", err)
	}
	log.Infof("registry: getLatestAgents(%s) = %s", rs.poolID.Hex(), lengthBig.Text(10))
	// TODO: If we are going to get 100s of agents, we probably need to batch the calls here.
	var agentConfigs []*config.AgentConfig
	length := int(lengthBig.Int64())
	for i := 0; i < length; i++ {
		agentID, agentRef, err := rs.contract.AgentAt(nil, rs.poolID, big.NewInt(int64(i)))
		if err != nil {
			return nil, fmt.Errorf("failed to get agent at index '%d' in pool '%s': %v", i, rs.poolID.String(), err)
		}
		agentCfg, err := rs.makeAgentConfig(agentID, agentRef)
		if err != nil {
			return nil, fmt.Errorf("failed to make agent config: %v", err)
		}
		agentConfigs = append(agentConfigs, &agentCfg)
	}
	return agentConfigs, nil
}

// Stop stops the registry service.
func (rs *RegistryService) Stop() error {
	close(rs.agentUpdates)
	return nil
}

// Name returns the name of the service.
func (rs *RegistryService) Name() string {
	return "RegistryService"
}
