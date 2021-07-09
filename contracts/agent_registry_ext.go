package contracts

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BindAgentRegistry binds a generic wrapper to an already deployed contract.
func BindAgentRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AgentRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// UnpackAgentRegistryAgentAdded unpacks the event.
func UnpackAgentRegistryAgentAdded(contract *bind.BoundContract, log *types.Log) (*AgentRegistryAgentAdded, error) {
	var transferEvent AgentRegistryAgentAdded
	return &transferEvent, contract.UnpackLog(&transferEvent, "AgentAdded", *log)
}

// UnpackAgentRegistryAgentRemoved unpacks the event.
func UnpackAgentRegistryAgentRemoved(contract *bind.BoundContract, log *types.Log) (*AgentRegistryAgentRemoved, error) {
	var transferEvent AgentRegistryAgentRemoved
	return &transferEvent, contract.UnpackLog(&transferEvent, "AgentRemoved", *log)
}
