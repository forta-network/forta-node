package registry

import (
	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/contracts"
	"OpenZeppelin/fortify-node/feeds"
	"OpenZeppelin/fortify-node/services"
	"fmt"
	"io"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	ipfsapi "github.com/ipfs/go-ipfs-api"
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
	client    *ethclient.Client
	msgClient clients.MessageClient
	txFeed    feeds.TransactionFeed

	logUnpacker LogUnpacker
	ipfsClient  IPFSClient

	agentsConfigs  []config.AgentConfig
	agentUpdates   chan *agentUpdate
	agentUpdatesWg sync.WaitGroup
}

// LogUnpacker unpacks agent events from logs.
type LogUnpacker interface {
	UnpackAgentRegistryAgentAdded(log *types.Log) (*contracts.AgentRegistryAgentAdded, error)
	UnpackAgentRegistryAgentUpdated(log *types.Log) (*contracts.AgentRegistryAgentUpdated, error)
	UnpackAgentRegistryAgentRemoved(log *types.Log) (*contracts.AgentRegistryAgentRemoved, error)
}

// IPFSClient interacts with an IPFS API/Gateway.
type IPFSClient interface {
	Cat(path string) (io.ReadCloser, error)
}

// New creates a new service.
func New(cfg config.Config, msgClient clients.MessageClient, txFeed feeds.TransactionFeed) services.Service {
	var ipfsURL string
	if cfg.Registry.IPFS != nil {
		ipfsURL = *cfg.Registry.IPFS
	} else {
		ipfsURL = config.DefaultIPFSGateway
	}

	return &RegistryService{
		cfg:          cfg,
		poolID:       common.HexToHash(cfg.Registry.PoolID),
		msgClient:    msgClient,
		txFeed:       txFeed,
		logUnpacker:  contracts.NewAgentLogUnpacker(common.HexToAddress(cfg.Registry.ContractAddress)),
		ipfsClient:   ipfsapi.NewShell(ipfsURL),
		agentUpdates: make(chan *agentUpdate, 100),
	}
}

// Start starts the registry service.
func (rs *RegistryService) Start() error {
	rpcClient, err := rpc.Dial(rs.cfg.Registry.JSONRPCURL)
	if err != nil {
		return err
	}
	rs.client = ethclient.NewClient(rpcClient)

	// Start detecting and buffering events.
	go rs.txFeed.ForEachTransaction(nil, rs.detectAgentEvents)

	// Start to handle agent updates but wait until initialization is complete.
	rs.agentUpdatesWg.Add(1)
	go rs.listenToAgentUpdates()

	if err := rs.publishLatestAgents(); err != nil {
		return fmt.Errorf("failed to publish the latest agents: %v", err)
	}

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

func (rs *RegistryService) getLatestAgents() ([]config.AgentConfig, error) {
	contract, err := contracts.NewAgentRegistryCaller(common.HexToAddress(rs.cfg.Registry.ContractAddress), rs.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create the agent registry caller: %v", err)
	}
	lengthBig, err := contract.AgentLength(nil, rs.poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the pool agents length: %v", err)
	}
	// TODO: If we are going to get 100s of agents, we probably need to batch the calls here.
	var agentConfigs []config.AgentConfig
	length := int(lengthBig.Int64())
	for i := 0; i < length; i++ {
		agentID, agentRef, err := contract.AgentAt(nil, rs.poolID, big.NewInt(int64(i)))
		if err != nil {
			return nil, fmt.Errorf("failed to get agent at index '%d' in pool '%s': %v", i, rs.poolID.String(), err)
		}
		agentCfg, err := rs.makeAgentConfig(agentID, agentRef)
		if err != nil {
			return nil, fmt.Errorf("failed to make agent config: %v", err)
		}
		agentConfigs = append(agentConfigs, agentCfg)
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
