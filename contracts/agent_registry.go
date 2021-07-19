// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AgentRegistryABI is the input ABI used to generate the binding from.
const AgentRegistryABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdminAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdminRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"PoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"PoolAdminAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"PoolAdminRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PoolOwnershipTransfered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_ref\",\"type\":\"string\"}],\"name\":\"addAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"addAgentAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"}],\"name\":\"addPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addPoolAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"agentAdminAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"agentAdminLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"agentAt\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"agentExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"}],\"name\":\"agentLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"agentRef\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractMinimalForwarderUpgradeable\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"poolAdminAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"}],\"name\":\"poolAdminLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"poolExistsMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"poolOwners\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"removeAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"removeAgentAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removePoolAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferPoolOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_ref\",\"type\":\"string\"}],\"name\":\"updateAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// AgentRegistry is an auto generated Go binding around an Ethereum contract.
type AgentRegistry struct {
	AgentRegistryCaller     // Read-only binding to the contract
	AgentRegistryTransactor // Write-only binding to the contract
	AgentRegistryFilterer   // Log filterer for contract events
}

// AgentRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AgentRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AgentRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgentRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgentRegistrySession struct {
	Contract     *AgentRegistry    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AgentRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgentRegistryCallerSession struct {
	Contract *AgentRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// AgentRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgentRegistryTransactorSession struct {
	Contract     *AgentRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// AgentRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AgentRegistryRaw struct {
	Contract *AgentRegistry // Generic contract binding to access the raw methods on
}

// AgentRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgentRegistryCallerRaw struct {
	Contract *AgentRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// AgentRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgentRegistryTransactorRaw struct {
	Contract *AgentRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAgentRegistry creates a new instance of AgentRegistry, bound to a specific deployed contract.
func NewAgentRegistry(address common.Address, backend bind.ContractBackend) (*AgentRegistry, error) {
	contract, err := bindAgentRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AgentRegistry{AgentRegistryCaller: AgentRegistryCaller{contract: contract}, AgentRegistryTransactor: AgentRegistryTransactor{contract: contract}, AgentRegistryFilterer: AgentRegistryFilterer{contract: contract}}, nil
}

// NewAgentRegistryCaller creates a new read-only instance of AgentRegistry, bound to a specific deployed contract.
func NewAgentRegistryCaller(address common.Address, caller bind.ContractCaller) (*AgentRegistryCaller, error) {
	contract, err := bindAgentRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryCaller{contract: contract}, nil
}

// NewAgentRegistryTransactor creates a new write-only instance of AgentRegistry, bound to a specific deployed contract.
func NewAgentRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AgentRegistryTransactor, error) {
	contract, err := bindAgentRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryTransactor{contract: contract}, nil
}

// NewAgentRegistryFilterer creates a new log filterer instance of AgentRegistry, bound to a specific deployed contract.
func NewAgentRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AgentRegistryFilterer, error) {
	contract, err := bindAgentRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryFilterer{contract: contract}, nil
}

// bindAgentRegistry binds a generic wrapper to an already deployed contract.
func bindAgentRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AgentRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentRegistry *AgentRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentRegistry.Contract.AgentRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentRegistry *AgentRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AgentRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentRegistry *AgentRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AgentRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentRegistry *AgentRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentRegistry *AgentRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentRegistry *AgentRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentRegistry.Contract.contract.Transact(opts, method, params...)
}

// AgentAdminAt is a free data retrieval call binding the contract method 0xe491fa83.
//
// Solidity: function agentAdminAt(bytes32 _poolId, bytes32 _agentId, uint256 index) view returns(address)
func (_AgentRegistry *AgentRegistryCaller) AgentAdminAt(opts *bind.CallOpts, _poolId [32]byte, _agentId [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentAdminAt", _poolId, _agentId, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentAdminAt is a free data retrieval call binding the contract method 0xe491fa83.
//
// Solidity: function agentAdminAt(bytes32 _poolId, bytes32 _agentId, uint256 index) view returns(address)
func (_AgentRegistry *AgentRegistrySession) AgentAdminAt(_poolId [32]byte, _agentId [32]byte, index *big.Int) (common.Address, error) {
	return _AgentRegistry.Contract.AgentAdminAt(&_AgentRegistry.CallOpts, _poolId, _agentId, index)
}

// AgentAdminAt is a free data retrieval call binding the contract method 0xe491fa83.
//
// Solidity: function agentAdminAt(bytes32 _poolId, bytes32 _agentId, uint256 index) view returns(address)
func (_AgentRegistry *AgentRegistryCallerSession) AgentAdminAt(_poolId [32]byte, _agentId [32]byte, index *big.Int) (common.Address, error) {
	return _AgentRegistry.Contract.AgentAdminAt(&_AgentRegistry.CallOpts, _poolId, _agentId, index)
}

// AgentAdminLength is a free data retrieval call binding the contract method 0x0a260645.
//
// Solidity: function agentAdminLength(bytes32 _poolId, bytes32 _agentId) view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) AgentAdminLength(opts *bind.CallOpts, _poolId [32]byte, _agentId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentAdminLength", _poolId, _agentId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentAdminLength is a free data retrieval call binding the contract method 0x0a260645.
//
// Solidity: function agentAdminLength(bytes32 _poolId, bytes32 _agentId) view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) AgentAdminLength(_poolId [32]byte, _agentId [32]byte) (*big.Int, error) {
	return _AgentRegistry.Contract.AgentAdminLength(&_AgentRegistry.CallOpts, _poolId, _agentId)
}

// AgentAdminLength is a free data retrieval call binding the contract method 0x0a260645.
//
// Solidity: function agentAdminLength(bytes32 _poolId, bytes32 _agentId) view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) AgentAdminLength(_poolId [32]byte, _agentId [32]byte) (*big.Int, error) {
	return _AgentRegistry.Contract.AgentAdminLength(&_AgentRegistry.CallOpts, _poolId, _agentId)
}

// AgentAt is a free data retrieval call binding the contract method 0x89d0a013.
//
// Solidity: function agentAt(bytes32 _poolId, uint256 index) view returns(bytes32, string)
func (_AgentRegistry *AgentRegistryCaller) AgentAt(opts *bind.CallOpts, _poolId [32]byte, index *big.Int) ([32]byte, string, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentAt", _poolId, index)

	if err != nil {
		return *new([32]byte), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// AgentAt is a free data retrieval call binding the contract method 0x89d0a013.
//
// Solidity: function agentAt(bytes32 _poolId, uint256 index) view returns(bytes32, string)
func (_AgentRegistry *AgentRegistrySession) AgentAt(_poolId [32]byte, index *big.Int) ([32]byte, string, error) {
	return _AgentRegistry.Contract.AgentAt(&_AgentRegistry.CallOpts, _poolId, index)
}

// AgentAt is a free data retrieval call binding the contract method 0x89d0a013.
//
// Solidity: function agentAt(bytes32 _poolId, uint256 index) view returns(bytes32, string)
func (_AgentRegistry *AgentRegistryCallerSession) AgentAt(_poolId [32]byte, index *big.Int) ([32]byte, string, error) {
	return _AgentRegistry.Contract.AgentAt(&_AgentRegistry.CallOpts, _poolId, index)
}

// AgentExists is a free data retrieval call binding the contract method 0x839ec6fd.
//
// Solidity: function agentExists(bytes32 _poolId, bytes32 _agentId) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) AgentExists(opts *bind.CallOpts, _poolId [32]byte, _agentId [32]byte) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentExists", _poolId, _agentId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AgentExists is a free data retrieval call binding the contract method 0x839ec6fd.
//
// Solidity: function agentExists(bytes32 _poolId, bytes32 _agentId) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) AgentExists(_poolId [32]byte, _agentId [32]byte) (bool, error) {
	return _AgentRegistry.Contract.AgentExists(&_AgentRegistry.CallOpts, _poolId, _agentId)
}

// AgentExists is a free data retrieval call binding the contract method 0x839ec6fd.
//
// Solidity: function agentExists(bytes32 _poolId, bytes32 _agentId) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) AgentExists(_poolId [32]byte, _agentId [32]byte) (bool, error) {
	return _AgentRegistry.Contract.AgentExists(&_AgentRegistry.CallOpts, _poolId, _agentId)
}

// AgentLength is a free data retrieval call binding the contract method 0x1e73feb1.
//
// Solidity: function agentLength(bytes32 _poolId) view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) AgentLength(opts *bind.CallOpts, _poolId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentLength", _poolId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentLength is a free data retrieval call binding the contract method 0x1e73feb1.
//
// Solidity: function agentLength(bytes32 _poolId) view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) AgentLength(_poolId [32]byte) (*big.Int, error) {
	return _AgentRegistry.Contract.AgentLength(&_AgentRegistry.CallOpts, _poolId)
}

// AgentLength is a free data retrieval call binding the contract method 0x1e73feb1.
//
// Solidity: function agentLength(bytes32 _poolId) view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) AgentLength(_poolId [32]byte) (*big.Int, error) {
	return _AgentRegistry.Contract.AgentLength(&_AgentRegistry.CallOpts, _poolId)
}

// AgentRef is a free data retrieval call binding the contract method 0x6a2189d2.
//
// Solidity: function agentRef(bytes32 _poolId, bytes32 _agentId) view returns(string)
func (_AgentRegistry *AgentRegistryCaller) AgentRef(opts *bind.CallOpts, _poolId [32]byte, _agentId [32]byte) (string, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentRef", _poolId, _agentId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AgentRef is a free data retrieval call binding the contract method 0x6a2189d2.
//
// Solidity: function agentRef(bytes32 _poolId, bytes32 _agentId) view returns(string)
func (_AgentRegistry *AgentRegistrySession) AgentRef(_poolId [32]byte, _agentId [32]byte) (string, error) {
	return _AgentRegistry.Contract.AgentRef(&_AgentRegistry.CallOpts, _poolId, _agentId)
}

// AgentRef is a free data retrieval call binding the contract method 0x6a2189d2.
//
// Solidity: function agentRef(bytes32 _poolId, bytes32 _agentId) view returns(string)
func (_AgentRegistry *AgentRegistryCallerSession) AgentRef(_poolId [32]byte, _agentId [32]byte) (string, error) {
	return _AgentRegistry.Contract.AgentRef(&_AgentRegistry.CallOpts, _poolId, _agentId)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) IsTrustedForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "isTrustedForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _AgentRegistry.Contract.IsTrustedForwarder(&_AgentRegistry.CallOpts, forwarder)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _AgentRegistry.Contract.IsTrustedForwarder(&_AgentRegistry.CallOpts, forwarder)
}

// PoolAdminAt is a free data retrieval call binding the contract method 0x43e32207.
//
// Solidity: function poolAdminAt(bytes32 _poolId, uint256 index) view returns(address)
func (_AgentRegistry *AgentRegistryCaller) PoolAdminAt(opts *bind.CallOpts, _poolId [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "poolAdminAt", _poolId, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolAdminAt is a free data retrieval call binding the contract method 0x43e32207.
//
// Solidity: function poolAdminAt(bytes32 _poolId, uint256 index) view returns(address)
func (_AgentRegistry *AgentRegistrySession) PoolAdminAt(_poolId [32]byte, index *big.Int) (common.Address, error) {
	return _AgentRegistry.Contract.PoolAdminAt(&_AgentRegistry.CallOpts, _poolId, index)
}

// PoolAdminAt is a free data retrieval call binding the contract method 0x43e32207.
//
// Solidity: function poolAdminAt(bytes32 _poolId, uint256 index) view returns(address)
func (_AgentRegistry *AgentRegistryCallerSession) PoolAdminAt(_poolId [32]byte, index *big.Int) (common.Address, error) {
	return _AgentRegistry.Contract.PoolAdminAt(&_AgentRegistry.CallOpts, _poolId, index)
}

// PoolAdminLength is a free data retrieval call binding the contract method 0xa0135277.
//
// Solidity: function poolAdminLength(bytes32 _poolId) view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) PoolAdminLength(opts *bind.CallOpts, _poolId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "poolAdminLength", _poolId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PoolAdminLength is a free data retrieval call binding the contract method 0xa0135277.
//
// Solidity: function poolAdminLength(bytes32 _poolId) view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) PoolAdminLength(_poolId [32]byte) (*big.Int, error) {
	return _AgentRegistry.Contract.PoolAdminLength(&_AgentRegistry.CallOpts, _poolId)
}

// PoolAdminLength is a free data retrieval call binding the contract method 0xa0135277.
//
// Solidity: function poolAdminLength(bytes32 _poolId) view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) PoolAdminLength(_poolId [32]byte) (*big.Int, error) {
	return _AgentRegistry.Contract.PoolAdminLength(&_AgentRegistry.CallOpts, _poolId)
}

// PoolExistsMap is a free data retrieval call binding the contract method 0x0ead48e6.
//
// Solidity: function poolExistsMap(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) PoolExistsMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "poolExistsMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PoolExistsMap is a free data retrieval call binding the contract method 0x0ead48e6.
//
// Solidity: function poolExistsMap(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) PoolExistsMap(arg0 [32]byte) (bool, error) {
	return _AgentRegistry.Contract.PoolExistsMap(&_AgentRegistry.CallOpts, arg0)
}

// PoolExistsMap is a free data retrieval call binding the contract method 0x0ead48e6.
//
// Solidity: function poolExistsMap(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) PoolExistsMap(arg0 [32]byte) (bool, error) {
	return _AgentRegistry.Contract.PoolExistsMap(&_AgentRegistry.CallOpts, arg0)
}

// PoolOwners is a free data retrieval call binding the contract method 0x219d2d09.
//
// Solidity: function poolOwners(bytes32 ) view returns(address)
func (_AgentRegistry *AgentRegistryCaller) PoolOwners(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "poolOwners", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolOwners is a free data retrieval call binding the contract method 0x219d2d09.
//
// Solidity: function poolOwners(bytes32 ) view returns(address)
func (_AgentRegistry *AgentRegistrySession) PoolOwners(arg0 [32]byte) (common.Address, error) {
	return _AgentRegistry.Contract.PoolOwners(&_AgentRegistry.CallOpts, arg0)
}

// PoolOwners is a free data retrieval call binding the contract method 0x219d2d09.
//
// Solidity: function poolOwners(bytes32 ) view returns(address)
func (_AgentRegistry *AgentRegistryCallerSession) PoolOwners(arg0 [32]byte) (common.Address, error) {
	return _AgentRegistry.Contract.PoolOwners(&_AgentRegistry.CallOpts, arg0)
}

// AddAgent is a paid mutator transaction binding the contract method 0x0944cd34.
//
// Solidity: function addAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactor) AddAgent(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "addAgent", _poolId, _agentId, _ref)
}

// AddAgent is a paid mutator transaction binding the contract method 0x0944cd34.
//
// Solidity: function addAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistrySession) AddAgent(_poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId, _ref)
}

// AddAgent is a paid mutator transaction binding the contract method 0x0944cd34.
//
// Solidity: function addAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) AddAgent(_poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId, _ref)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x63438a97.
//
// Solidity: function addAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_AgentRegistry *AgentRegistryTransactor) AddAgentAdmin(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "addAgentAdmin", _poolId, _agentId, _admin)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x63438a97.
//
// Solidity: function addAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_AgentRegistry *AgentRegistrySession) AddAgentAdmin(_poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddAgentAdmin(&_AgentRegistry.TransactOpts, _poolId, _agentId, _admin)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x63438a97.
//
// Solidity: function addAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) AddAgentAdmin(_poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddAgentAdmin(&_AgentRegistry.TransactOpts, _poolId, _agentId, _admin)
}

// AddPool is a paid mutator transaction binding the contract method 0x8a0033fc.
//
// Solidity: function addPool(bytes32 _poolId) returns()
func (_AgentRegistry *AgentRegistryTransactor) AddPool(opts *bind.TransactOpts, _poolId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "addPool", _poolId)
}

// AddPool is a paid mutator transaction binding the contract method 0x8a0033fc.
//
// Solidity: function addPool(bytes32 _poolId) returns()
func (_AgentRegistry *AgentRegistrySession) AddPool(_poolId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddPool(&_AgentRegistry.TransactOpts, _poolId)
}

// AddPool is a paid mutator transaction binding the contract method 0x8a0033fc.
//
// Solidity: function addPool(bytes32 _poolId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) AddPool(_poolId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddPool(&_AgentRegistry.TransactOpts, _poolId)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x84930b76.
//
// Solidity: function addPoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactor) AddPoolAdmin(opts *bind.TransactOpts, _poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "addPoolAdmin", _poolId, admin)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x84930b76.
//
// Solidity: function addPoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistrySession) AddPoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddPoolAdmin(&_AgentRegistry.TransactOpts, _poolId, admin)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x84930b76.
//
// Solidity: function addPoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) AddPoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddPoolAdmin(&_AgentRegistry.TransactOpts, _poolId, admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address forwarder) returns()
func (_AgentRegistry *AgentRegistryTransactor) Initialize(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "initialize", forwarder)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address forwarder) returns()
func (_AgentRegistry *AgentRegistrySession) Initialize(forwarder common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.Initialize(&_AgentRegistry.TransactOpts, forwarder)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address forwarder) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) Initialize(forwarder common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.Initialize(&_AgentRegistry.TransactOpts, forwarder)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x61e62a0f.
//
// Solidity: function removeAgent(bytes32 _poolId, bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistryTransactor) RemoveAgent(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "removeAgent", _poolId, _agentId)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x61e62a0f.
//
// Solidity: function removeAgent(bytes32 _poolId, bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistrySession) RemoveAgent(_poolId [32]byte, _agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemoveAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x61e62a0f.
//
// Solidity: function removeAgent(bytes32 _poolId, bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) RemoveAgent(_poolId [32]byte, _agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemoveAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0xa7f2f5ab.
//
// Solidity: function removeAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_AgentRegistry *AgentRegistryTransactor) RemoveAgentAdmin(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "removeAgentAdmin", _poolId, _agentId, _admin)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0xa7f2f5ab.
//
// Solidity: function removeAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_AgentRegistry *AgentRegistrySession) RemoveAgentAdmin(_poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemoveAgentAdmin(&_AgentRegistry.TransactOpts, _poolId, _agentId, _admin)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0xa7f2f5ab.
//
// Solidity: function removeAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) RemoveAgentAdmin(_poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemoveAgentAdmin(&_AgentRegistry.TransactOpts, _poolId, _agentId, _admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xacce5630.
//
// Solidity: function removePoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactor) RemovePoolAdmin(opts *bind.TransactOpts, _poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "removePoolAdmin", _poolId, admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xacce5630.
//
// Solidity: function removePoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistrySession) RemovePoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemovePoolAdmin(&_AgentRegistry.TransactOpts, _poolId, admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xacce5630.
//
// Solidity: function removePoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) RemovePoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemovePoolAdmin(&_AgentRegistry.TransactOpts, _poolId, admin)
}

// TransferPoolOwnership is a paid mutator transaction binding the contract method 0x1b6feb92.
//
// Solidity: function transferPoolOwnership(bytes32 _poolId, address _to) returns()
func (_AgentRegistry *AgentRegistryTransactor) TransferPoolOwnership(opts *bind.TransactOpts, _poolId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "transferPoolOwnership", _poolId, _to)
}

// TransferPoolOwnership is a paid mutator transaction binding the contract method 0x1b6feb92.
//
// Solidity: function transferPoolOwnership(bytes32 _poolId, address _to) returns()
func (_AgentRegistry *AgentRegistrySession) TransferPoolOwnership(_poolId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.TransferPoolOwnership(&_AgentRegistry.TransactOpts, _poolId, _to)
}

// TransferPoolOwnership is a paid mutator transaction binding the contract method 0x1b6feb92.
//
// Solidity: function transferPoolOwnership(bytes32 _poolId, address _to) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) TransferPoolOwnership(_poolId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.TransferPoolOwnership(&_AgentRegistry.TransactOpts, _poolId, _to)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xe0434e3f.
//
// Solidity: function updateAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactor) UpdateAgent(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "updateAgent", _poolId, _agentId, _ref)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xe0434e3f.
//
// Solidity: function updateAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistrySession) UpdateAgent(_poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpdateAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId, _ref)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xe0434e3f.
//
// Solidity: function updateAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) UpdateAgent(_poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpdateAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId, _ref)
}

// AgentRegistryAgentAddedIterator is returned from FilterAgentAdded and is used to iterate over the raw logs and unpacked data for AgentAdded events raised by the AgentRegistry contract.
type AgentRegistryAgentAddedIterator struct {
	Event *AgentRegistryAgentAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentRegistryAgentAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentRegistryAgentAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentRegistryAgentAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentAdded represents a AgentAdded event raised by the AgentRegistry contract.
type AgentRegistryAgentAdded struct {
	PoolId  [32]byte
	AgentId [32]byte
	Ref     string
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdded is a free log retrieval operation binding the contract event 0xdd1417d47ca73bdaf24e5c1df4158e61cbbb73fda6562b14a16215ccfa201aeb.
//
// Solidity: event AgentAdded(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentAdded(opts *bind.FilterOpts) (*AgentRegistryAgentAddedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentAdded")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentAddedIterator{contract: _AgentRegistry.contract, event: "AgentAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdded is a free log subscription operation binding the contract event 0xdd1417d47ca73bdaf24e5c1df4158e61cbbb73fda6562b14a16215ccfa201aeb.
//
// Solidity: event AgentAdded(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentAdded(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentAdded) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentAdded)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentAdded is a log parse operation binding the contract event 0xdd1417d47ca73bdaf24e5c1df4158e61cbbb73fda6562b14a16215ccfa201aeb.
//
// Solidity: event AgentAdded(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentAdded(log types.Log) (*AgentRegistryAgentAdded, error) {
	event := new(AgentRegistryAgentAdded)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentAdminAddedIterator is returned from FilterAgentAdminAdded and is used to iterate over the raw logs and unpacked data for AgentAdminAdded events raised by the AgentRegistry contract.
type AgentRegistryAgentAdminAddedIterator struct {
	Event *AgentRegistryAgentAdminAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentRegistryAgentAdminAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentAdminAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentRegistryAgentAdminAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentRegistryAgentAdminAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentAdminAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentAdminAdded represents a AgentAdminAdded event raised by the AgentRegistry contract.
type AgentRegistryAgentAdminAdded struct {
	PoolId  [32]byte
	AgentId [32]byte
	Admin   common.Address
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdminAdded is a free log retrieval operation binding the contract event 0x11cbfe7067a564a8ddb9a780914987915bc0d29835c8b0554a9d9db2e51064ed.
//
// Solidity: event AgentAdminAdded(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentAdminAdded(opts *bind.FilterOpts) (*AgentRegistryAgentAdminAddedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentAdminAdded")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentAdminAddedIterator{contract: _AgentRegistry.contract, event: "AgentAdminAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdminAdded is a free log subscription operation binding the contract event 0x11cbfe7067a564a8ddb9a780914987915bc0d29835c8b0554a9d9db2e51064ed.
//
// Solidity: event AgentAdminAdded(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentAdminAdded(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentAdminAdded) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentAdminAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentAdminAdded)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdminAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentAdminAdded is a log parse operation binding the contract event 0x11cbfe7067a564a8ddb9a780914987915bc0d29835c8b0554a9d9db2e51064ed.
//
// Solidity: event AgentAdminAdded(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentAdminAdded(log types.Log) (*AgentRegistryAgentAdminAdded, error) {
	event := new(AgentRegistryAgentAdminAdded)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdminAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentAdminRemovedIterator is returned from FilterAgentAdminRemoved and is used to iterate over the raw logs and unpacked data for AgentAdminRemoved events raised by the AgentRegistry contract.
type AgentRegistryAgentAdminRemovedIterator struct {
	Event *AgentRegistryAgentAdminRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentRegistryAgentAdminRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentAdminRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentRegistryAgentAdminRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentRegistryAgentAdminRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentAdminRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentAdminRemoved represents a AgentAdminRemoved event raised by the AgentRegistry contract.
type AgentRegistryAgentAdminRemoved struct {
	PoolId  [32]byte
	AgentId [32]byte
	Admin   common.Address
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdminRemoved is a free log retrieval operation binding the contract event 0x801fe4d9d51c0445c2938970ad5e042c3f2f2c387d67de4111af1bc60d22dba9.
//
// Solidity: event AgentAdminRemoved(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentAdminRemoved(opts *bind.FilterOpts) (*AgentRegistryAgentAdminRemovedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentAdminRemoved")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentAdminRemovedIterator{contract: _AgentRegistry.contract, event: "AgentAdminRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentAdminRemoved is a free log subscription operation binding the contract event 0x801fe4d9d51c0445c2938970ad5e042c3f2f2c387d67de4111af1bc60d22dba9.
//
// Solidity: event AgentAdminRemoved(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentAdminRemoved(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentAdminRemoved) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentAdminRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentAdminRemoved)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdminRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentAdminRemoved is a log parse operation binding the contract event 0x801fe4d9d51c0445c2938970ad5e042c3f2f2c387d67de4111af1bc60d22dba9.
//
// Solidity: event AgentAdminRemoved(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentAdminRemoved(log types.Log) (*AgentRegistryAgentAdminRemoved, error) {
	event := new(AgentRegistryAgentAdminRemoved)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdminRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentRemovedIterator is returned from FilterAgentRemoved and is used to iterate over the raw logs and unpacked data for AgentRemoved events raised by the AgentRegistry contract.
type AgentRegistryAgentRemovedIterator struct {
	Event *AgentRegistryAgentRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentRegistryAgentRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentRegistryAgentRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentRegistryAgentRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentRemoved represents a AgentRemoved event raised by the AgentRegistry contract.
type AgentRegistryAgentRemoved struct {
	PoolId  [32]byte
	AgentId [32]byte
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentRemoved is a free log retrieval operation binding the contract event 0xccc61238b7da918476155ec78f1d6f23cd587a3f796774b80c0c79dde8f57e6b.
//
// Solidity: event AgentRemoved(bytes32 poolId, bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentRemoved(opts *bind.FilterOpts) (*AgentRegistryAgentRemovedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentRemoved")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentRemovedIterator{contract: _AgentRegistry.contract, event: "AgentRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentRemoved is a free log subscription operation binding the contract event 0xccc61238b7da918476155ec78f1d6f23cd587a3f796774b80c0c79dde8f57e6b.
//
// Solidity: event AgentRemoved(bytes32 poolId, bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentRemoved(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentRemoved) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentRemoved)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentRemoved is a log parse operation binding the contract event 0xccc61238b7da918476155ec78f1d6f23cd587a3f796774b80c0c79dde8f57e6b.
//
// Solidity: event AgentRemoved(bytes32 poolId, bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentRemoved(log types.Log) (*AgentRegistryAgentRemoved, error) {
	event := new(AgentRegistryAgentRemoved)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentUpdatedIterator is returned from FilterAgentUpdated and is used to iterate over the raw logs and unpacked data for AgentUpdated events raised by the AgentRegistry contract.
type AgentRegistryAgentUpdatedIterator struct {
	Event *AgentRegistryAgentUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentRegistryAgentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentRegistryAgentUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentRegistryAgentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentUpdated represents a AgentUpdated event raised by the AgentRegistry contract.
type AgentRegistryAgentUpdated struct {
	PoolId  [32]byte
	AgentId [32]byte
	Ref     string
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentUpdated is a free log retrieval operation binding the contract event 0x8a33912f31f80df098dcad04aba17afc8c52d0f700baf2750ea3899f9230a850.
//
// Solidity: event AgentUpdated(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentUpdated(opts *bind.FilterOpts) (*AgentRegistryAgentUpdatedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentUpdated")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentUpdatedIterator{contract: _AgentRegistry.contract, event: "AgentUpdated", logs: logs, sub: sub}, nil
}

// WatchAgentUpdated is a free log subscription operation binding the contract event 0x8a33912f31f80df098dcad04aba17afc8c52d0f700baf2750ea3899f9230a850.
//
// Solidity: event AgentUpdated(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentUpdated(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentUpdated) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentUpdated)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentUpdated is a log parse operation binding the contract event 0x8a33912f31f80df098dcad04aba17afc8c52d0f700baf2750ea3899f9230a850.
//
// Solidity: event AgentUpdated(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentUpdated(log types.Log) (*AgentRegistryAgentUpdated, error) {
	event := new(AgentRegistryAgentUpdated)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryPoolAddedIterator is returned from FilterPoolAdded and is used to iterate over the raw logs and unpacked data for PoolAdded events raised by the AgentRegistry contract.
type AgentRegistryPoolAddedIterator struct {
	Event *AgentRegistryPoolAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentRegistryPoolAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryPoolAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentRegistryPoolAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentRegistryPoolAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryPoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryPoolAdded represents a PoolAdded event raised by the AgentRegistry contract.
type AgentRegistryPoolAdded struct {
	PoolId [32]byte
	By     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolAdded is a free log retrieval operation binding the contract event 0xfd0fa7919fbe3857a4236750e8d3e42ac691881d2f31c50ebe4cfc2a7705ee80.
//
// Solidity: event PoolAdded(bytes32 poolId, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterPoolAdded(opts *bind.FilterOpts) (*AgentRegistryPoolAddedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "PoolAdded")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryPoolAddedIterator{contract: _AgentRegistry.contract, event: "PoolAdded", logs: logs, sub: sub}, nil
}

// WatchPoolAdded is a free log subscription operation binding the contract event 0xfd0fa7919fbe3857a4236750e8d3e42ac691881d2f31c50ebe4cfc2a7705ee80.
//
// Solidity: event PoolAdded(bytes32 poolId, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchPoolAdded(opts *bind.WatchOpts, sink chan<- *AgentRegistryPoolAdded) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "PoolAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryPoolAdded)
				if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePoolAdded is a log parse operation binding the contract event 0xfd0fa7919fbe3857a4236750e8d3e42ac691881d2f31c50ebe4cfc2a7705ee80.
//
// Solidity: event PoolAdded(bytes32 poolId, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParsePoolAdded(log types.Log) (*AgentRegistryPoolAdded, error) {
	event := new(AgentRegistryPoolAdded)
	if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryPoolAdminAddedIterator is returned from FilterPoolAdminAdded and is used to iterate over the raw logs and unpacked data for PoolAdminAdded events raised by the AgentRegistry contract.
type AgentRegistryPoolAdminAddedIterator struct {
	Event *AgentRegistryPoolAdminAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentRegistryPoolAdminAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryPoolAdminAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentRegistryPoolAdminAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentRegistryPoolAdminAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryPoolAdminAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryPoolAdminAdded represents a PoolAdminAdded event raised by the AgentRegistry contract.
type AgentRegistryPoolAdminAdded struct {
	PoolId [32]byte
	Addr   common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolAdminAdded is a free log retrieval operation binding the contract event 0x4c2f0be874d7047c389e9fadf454854a8f5b72b49de199c7d60831b6c462a9a0.
//
// Solidity: event PoolAdminAdded(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) FilterPoolAdminAdded(opts *bind.FilterOpts) (*AgentRegistryPoolAdminAddedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "PoolAdminAdded")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryPoolAdminAddedIterator{contract: _AgentRegistry.contract, event: "PoolAdminAdded", logs: logs, sub: sub}, nil
}

// WatchPoolAdminAdded is a free log subscription operation binding the contract event 0x4c2f0be874d7047c389e9fadf454854a8f5b72b49de199c7d60831b6c462a9a0.
//
// Solidity: event PoolAdminAdded(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) WatchPoolAdminAdded(opts *bind.WatchOpts, sink chan<- *AgentRegistryPoolAdminAdded) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "PoolAdminAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryPoolAdminAdded)
				if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdminAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePoolAdminAdded is a log parse operation binding the contract event 0x4c2f0be874d7047c389e9fadf454854a8f5b72b49de199c7d60831b6c462a9a0.
//
// Solidity: event PoolAdminAdded(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) ParsePoolAdminAdded(log types.Log) (*AgentRegistryPoolAdminAdded, error) {
	event := new(AgentRegistryPoolAdminAdded)
	if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdminAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryPoolAdminRemovedIterator is returned from FilterPoolAdminRemoved and is used to iterate over the raw logs and unpacked data for PoolAdminRemoved events raised by the AgentRegistry contract.
type AgentRegistryPoolAdminRemovedIterator struct {
	Event *AgentRegistryPoolAdminRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentRegistryPoolAdminRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryPoolAdminRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentRegistryPoolAdminRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentRegistryPoolAdminRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryPoolAdminRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryPoolAdminRemoved represents a PoolAdminRemoved event raised by the AgentRegistry contract.
type AgentRegistryPoolAdminRemoved struct {
	PoolId [32]byte
	Addr   common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolAdminRemoved is a free log retrieval operation binding the contract event 0x25fcb4cb225521d9a561fd90af5809d997ae0fa769a5c95d90750162623bf624.
//
// Solidity: event PoolAdminRemoved(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) FilterPoolAdminRemoved(opts *bind.FilterOpts) (*AgentRegistryPoolAdminRemovedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "PoolAdminRemoved")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryPoolAdminRemovedIterator{contract: _AgentRegistry.contract, event: "PoolAdminRemoved", logs: logs, sub: sub}, nil
}

// WatchPoolAdminRemoved is a free log subscription operation binding the contract event 0x25fcb4cb225521d9a561fd90af5809d997ae0fa769a5c95d90750162623bf624.
//
// Solidity: event PoolAdminRemoved(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) WatchPoolAdminRemoved(opts *bind.WatchOpts, sink chan<- *AgentRegistryPoolAdminRemoved) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "PoolAdminRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryPoolAdminRemoved)
				if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdminRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePoolAdminRemoved is a log parse operation binding the contract event 0x25fcb4cb225521d9a561fd90af5809d997ae0fa769a5c95d90750162623bf624.
//
// Solidity: event PoolAdminRemoved(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) ParsePoolAdminRemoved(log types.Log) (*AgentRegistryPoolAdminRemoved, error) {
	event := new(AgentRegistryPoolAdminRemoved)
	if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdminRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryPoolOwnershipTransferedIterator is returned from FilterPoolOwnershipTransfered and is used to iterate over the raw logs and unpacked data for PoolOwnershipTransfered events raised by the AgentRegistry contract.
type AgentRegistryPoolOwnershipTransferedIterator struct {
	Event *AgentRegistryPoolOwnershipTransfered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentRegistryPoolOwnershipTransferedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryPoolOwnershipTransfered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentRegistryPoolOwnershipTransfered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentRegistryPoolOwnershipTransferedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryPoolOwnershipTransferedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryPoolOwnershipTransfered represents a PoolOwnershipTransfered event raised by the AgentRegistry contract.
type AgentRegistryPoolOwnershipTransfered struct {
	PoolId [32]byte
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolOwnershipTransfered is a free log retrieval operation binding the contract event 0xecc1c418cc7d2c3062b60c1d8f0b7748e2c8c10a109df4f6089b285f26b993d6.
//
// Solidity: event PoolOwnershipTransfered(bytes32 poolId, address from, address to)
func (_AgentRegistry *AgentRegistryFilterer) FilterPoolOwnershipTransfered(opts *bind.FilterOpts) (*AgentRegistryPoolOwnershipTransferedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "PoolOwnershipTransfered")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryPoolOwnershipTransferedIterator{contract: _AgentRegistry.contract, event: "PoolOwnershipTransfered", logs: logs, sub: sub}, nil
}

// WatchPoolOwnershipTransfered is a free log subscription operation binding the contract event 0xecc1c418cc7d2c3062b60c1d8f0b7748e2c8c10a109df4f6089b285f26b993d6.
//
// Solidity: event PoolOwnershipTransfered(bytes32 poolId, address from, address to)
func (_AgentRegistry *AgentRegistryFilterer) WatchPoolOwnershipTransfered(opts *bind.WatchOpts, sink chan<- *AgentRegistryPoolOwnershipTransfered) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "PoolOwnershipTransfered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryPoolOwnershipTransfered)
				if err := _AgentRegistry.contract.UnpackLog(event, "PoolOwnershipTransfered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePoolOwnershipTransfered is a log parse operation binding the contract event 0xecc1c418cc7d2c3062b60c1d8f0b7748e2c8c10a109df4f6089b285f26b993d6.
//
// Solidity: event PoolOwnershipTransfered(bytes32 poolId, address from, address to)
func (_AgentRegistry *AgentRegistryFilterer) ParsePoolOwnershipTransfered(log types.Log) (*AgentRegistryPoolOwnershipTransfered, error) {
	event := new(AgentRegistryPoolOwnershipTransfered)
	if err := _AgentRegistry.contract.UnpackLog(event, "PoolOwnershipTransfered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
