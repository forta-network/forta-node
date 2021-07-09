package registry

import (
	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/contracts"
	"OpenZeppelin/fortify-node/services"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// RegistryService listens to the agent pool changes so the node
// can stay in sync.
// TODO: Instead of publishing messages for the static config,
// check or listen to an actively maintained resource (e.g. a smart contract).
// TODO: The registry service or the config should construct unique names for the agents.
type RegistryService struct {
	cfg       config.Config
	client    *ethclient.Client
	msgClient clients.MessageClient
}

// New creates a new service.
func New(cfg config.Config, msgClient clients.MessageClient) services.Service {
	return &RegistryService{
		cfg:       cfg,
		msgClient: msgClient,
	}
}

// Start starts the registry service.
func (rs *RegistryService) Start() error {
	rpcClient, err := rpc.Dial(rs.cfg.Registry.JSONRPCURL)
	if err != nil {
		return err
	}
	rs.client = ethclient.NewClient(rpcClient)

	// TODO: Meanwhile we read the contract, agent updates can happen
	// so we need to start by listening to new events first and buffer them.
	return rs.publishLatestAgents()
}

func (rs *RegistryService) publishLatestAgents() error {
	agents, err := rs.getLatestAgents()
	if err != nil {
		return err
	}
	rs.msgClient.Publish(messaging.SubjectAgentsVersionsLatest, agents)
	return nil
}

func (rs *RegistryService) getLatestAgents() ([]config.AgentConfig, error) {
	contract, err := contracts.NewAgentRegistryCaller(common.HexToAddress(rs.cfg.Registry.ContractAddress), rs.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create the agent registry caller: %v", err)
	}
	poolID := common.BytesToHash([]byte(rs.cfg.Registry.PoolID))
	lengthBig, err := contract.AgentLength(nil, poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the pool agents length: %v", err)
	}
	// TODO: If we are going to get 100s of agents, we need to batch the calls here.
	var agentConfigs []config.AgentConfig
	length := int(lengthBig.Int64())
	for i := 0; i < length; i++ {
		_, agentRef, err := contract.AgentAt(nil, poolID, big.NewInt(int64(i)))
		if err != nil {
			return nil, fmt.Errorf("failed to get agent at index '%d' in pool '%s': %v", i, poolID.String(), err)
		}
		// TODO: Maybe we can just use single reference?
		agentConfigs = append(agentConfigs, config.AgentConfig{
			Name:  agentRef,
			Image: agentRef,
		})
	}
	return agentConfigs, nil
}

// Stop stops the registry service.
func (rs *RegistryService) Stop() error {
	return nil
}

// Name returns the name of the service.
func (rs *RegistryService) Name() string {
	return "RegistryService"
}
