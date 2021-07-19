package contracts

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	agentAddedTopicHash   = crypto.Keccak256Hash([]byte("AgentAdded(bytes32,bytes32,string,address)"))
	agentUpdatedTopicHash = crypto.Keccak256Hash([]byte("AgentUpdated(bytes32,bytes32,string,address)"))
	agentRemovedTopicHash = crypto.Keccak256Hash([]byte("AgentRemoved(bytes32,bytes32,address)"))
)

// Errors from event parsing
var (
	ErrNoEvents = errors.New("no event found")
)

// BindAgentRegistry binds a generic wrapper to an already deployed contract.
func BindAgentRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AgentRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// AgentLogUnpacker unpacks logs using the bound contract.
type AgentLogUnpacker struct {
	contract *bind.BoundContract
}

func NewAgentLogUnpacker(address common.Address) *AgentLogUnpacker {
	contract, _ := BindAgentRegistry(address, nil, nil, nil)
	return &AgentLogUnpacker{
		contract: contract,
	}
}

// UnpackAgentRegistryAgentAdded unpacks the event.
func (alu *AgentLogUnpacker) UnpackAgentRegistryAgentAdded(log *types.Log) (*AgentRegistryAgentAdded, error) {
	if !hasTopic(log.Topics, agentAddedTopicHash) {
		return nil, ErrNoEvents
	}
	var agentEvent AgentRegistryAgentAdded
	return &agentEvent, alu.contract.UnpackLog(&agentEvent, "AgentAdded", *log)
}

// UnpackAgentRegistryAgentUpdated unpacks the event.
func (alu *AgentLogUnpacker) UnpackAgentRegistryAgentUpdated(log *types.Log) (*AgentRegistryAgentUpdated, error) {
	if !hasTopic(log.Topics, agentUpdatedTopicHash) {
		return nil, ErrNoEvents
	}
	var agentEvent AgentRegistryAgentUpdated
	return &agentEvent, alu.contract.UnpackLog(&agentEvent, "AgentUpdated", *log)
}

// UnpackAgentRegistryAgentRemoved unpacks the event.
func (alu *AgentLogUnpacker) UnpackAgentRegistryAgentRemoved(log *types.Log) (*AgentRegistryAgentRemoved, error) {
	if !hasTopic(log.Topics, agentRemovedTopicHash) {
		return nil, ErrNoEvents
	}
	var agentEvent AgentRegistryAgentRemoved
	return &agentEvent, alu.contract.UnpackLog(&agentEvent, "AgentRemoved", *log)
}

func hasTopic(topics []common.Hash, expected common.Hash) bool {
	for _, topic := range topics {
		if topic.String() == expected.String() {
			return true
		}
	}
	return false
}
