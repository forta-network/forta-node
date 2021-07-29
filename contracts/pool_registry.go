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

// PoolRegistryABI is the input ABI used to generate the binding from.
const PoolRegistryABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdminAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdminRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"AgentRegistryChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"PoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"PoolAdminAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"PoolAdminRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PoolOwnershipTransfered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"latest\",\"type\":\"bool\"}],\"name\":\"addAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"addAgentAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"}],\"name\":\"addPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addPoolAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"agentAdminAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"agentAdminLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"agentAt\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"agentExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"}],\"name\":\"agentLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"agentRef\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"agentUsingLatest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"agentVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"}],\"name\":\"getPoolHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractMinimalForwarderUpgradeable\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"poolAdminAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"}],\"name\":\"poolAdminLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"poolExistsMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"poolOwners\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"removeAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"removeAgentAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removePoolAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setAgentRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferPoolOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"latest\",\"type\":\"bool\"}],\"name\":\"updateAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// PoolRegistry is an auto generated Go binding around an Ethereum contract.
type PoolRegistry struct {
	PoolRegistryCaller     // Read-only binding to the contract
	PoolRegistryTransactor // Write-only binding to the contract
	PoolRegistryFilterer   // Log filterer for contract events
}

// PoolRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoolRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoolRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoolRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoolRegistrySession struct {
	Contract     *PoolRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoolRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoolRegistryCallerSession struct {
	Contract *PoolRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// PoolRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoolRegistryTransactorSession struct {
	Contract     *PoolRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PoolRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoolRegistryRaw struct {
	Contract *PoolRegistry // Generic contract binding to access the raw methods on
}

// PoolRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoolRegistryCallerRaw struct {
	Contract *PoolRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// PoolRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoolRegistryTransactorRaw struct {
	Contract *PoolRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoolRegistry creates a new instance of PoolRegistry, bound to a specific deployed contract.
func NewPoolRegistry(address common.Address, backend bind.ContractBackend) (*PoolRegistry, error) {
	contract, err := bindPoolRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PoolRegistry{PoolRegistryCaller: PoolRegistryCaller{contract: contract}, PoolRegistryTransactor: PoolRegistryTransactor{contract: contract}, PoolRegistryFilterer: PoolRegistryFilterer{contract: contract}}, nil
}

// NewPoolRegistryCaller creates a new read-only instance of PoolRegistry, bound to a specific deployed contract.
func NewPoolRegistryCaller(address common.Address, caller bind.ContractCaller) (*PoolRegistryCaller, error) {
	contract, err := bindPoolRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoolRegistryCaller{contract: contract}, nil
}

// NewPoolRegistryTransactor creates a new write-only instance of PoolRegistry, bound to a specific deployed contract.
func NewPoolRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*PoolRegistryTransactor, error) {
	contract, err := bindPoolRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoolRegistryTransactor{contract: contract}, nil
}

// NewPoolRegistryFilterer creates a new log filterer instance of PoolRegistry, bound to a specific deployed contract.
func NewPoolRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*PoolRegistryFilterer, error) {
	contract, err := bindPoolRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoolRegistryFilterer{contract: contract}, nil
}

// bindPoolRegistry binds a generic wrapper to an already deployed contract.
func bindPoolRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PoolRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolRegistry *PoolRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolRegistry.Contract.PoolRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolRegistry *PoolRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolRegistry.Contract.PoolRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolRegistry *PoolRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolRegistry.Contract.PoolRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolRegistry *PoolRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolRegistry *PoolRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolRegistry *PoolRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolRegistry.Contract.contract.Transact(opts, method, params...)
}

// AgentAdminAt is a free data retrieval call binding the contract method 0xe491fa83.
//
// Solidity: function agentAdminAt(bytes32 _poolId, bytes32 _agentId, uint256 index) view returns(address)
func (_PoolRegistry *PoolRegistryCaller) AgentAdminAt(opts *bind.CallOpts, _poolId [32]byte, _agentId [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "agentAdminAt", _poolId, _agentId, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentAdminAt is a free data retrieval call binding the contract method 0xe491fa83.
//
// Solidity: function agentAdminAt(bytes32 _poolId, bytes32 _agentId, uint256 index) view returns(address)
func (_PoolRegistry *PoolRegistrySession) AgentAdminAt(_poolId [32]byte, _agentId [32]byte, index *big.Int) (common.Address, error) {
	return _PoolRegistry.Contract.AgentAdminAt(&_PoolRegistry.CallOpts, _poolId, _agentId, index)
}

// AgentAdminAt is a free data retrieval call binding the contract method 0xe491fa83.
//
// Solidity: function agentAdminAt(bytes32 _poolId, bytes32 _agentId, uint256 index) view returns(address)
func (_PoolRegistry *PoolRegistryCallerSession) AgentAdminAt(_poolId [32]byte, _agentId [32]byte, index *big.Int) (common.Address, error) {
	return _PoolRegistry.Contract.AgentAdminAt(&_PoolRegistry.CallOpts, _poolId, _agentId, index)
}

// AgentAdminLength is a free data retrieval call binding the contract method 0x0a260645.
//
// Solidity: function agentAdminLength(bytes32 _poolId, bytes32 _agentId) view returns(uint256)
func (_PoolRegistry *PoolRegistryCaller) AgentAdminLength(opts *bind.CallOpts, _poolId [32]byte, _agentId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "agentAdminLength", _poolId, _agentId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentAdminLength is a free data retrieval call binding the contract method 0x0a260645.
//
// Solidity: function agentAdminLength(bytes32 _poolId, bytes32 _agentId) view returns(uint256)
func (_PoolRegistry *PoolRegistrySession) AgentAdminLength(_poolId [32]byte, _agentId [32]byte) (*big.Int, error) {
	return _PoolRegistry.Contract.AgentAdminLength(&_PoolRegistry.CallOpts, _poolId, _agentId)
}

// AgentAdminLength is a free data retrieval call binding the contract method 0x0a260645.
//
// Solidity: function agentAdminLength(bytes32 _poolId, bytes32 _agentId) view returns(uint256)
func (_PoolRegistry *PoolRegistryCallerSession) AgentAdminLength(_poolId [32]byte, _agentId [32]byte) (*big.Int, error) {
	return _PoolRegistry.Contract.AgentAdminLength(&_PoolRegistry.CallOpts, _poolId, _agentId)
}

// AgentAt is a free data retrieval call binding the contract method 0x89d0a013.
//
// Solidity: function agentAt(bytes32 _poolId, uint256 index) view returns(bytes32, uint256, bool, string, bool)
func (_PoolRegistry *PoolRegistryCaller) AgentAt(opts *bind.CallOpts, _poolId [32]byte, index *big.Int) ([32]byte, *big.Int, bool, string, bool, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "agentAt", _poolId, index)

	if err != nil {
		return *new([32]byte), *new(*big.Int), *new(bool), *new(string), *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(bool)).(*bool)
	out3 := *abi.ConvertType(out[3], new(string)).(*string)
	out4 := *abi.ConvertType(out[4], new(bool)).(*bool)

	return out0, out1, out2, out3, out4, err

}

// AgentAt is a free data retrieval call binding the contract method 0x89d0a013.
//
// Solidity: function agentAt(bytes32 _poolId, uint256 index) view returns(bytes32, uint256, bool, string, bool)
func (_PoolRegistry *PoolRegistrySession) AgentAt(_poolId [32]byte, index *big.Int) ([32]byte, *big.Int, bool, string, bool, error) {
	return _PoolRegistry.Contract.AgentAt(&_PoolRegistry.CallOpts, _poolId, index)
}

// AgentAt is a free data retrieval call binding the contract method 0x89d0a013.
//
// Solidity: function agentAt(bytes32 _poolId, uint256 index) view returns(bytes32, uint256, bool, string, bool)
func (_PoolRegistry *PoolRegistryCallerSession) AgentAt(_poolId [32]byte, index *big.Int) ([32]byte, *big.Int, bool, string, bool, error) {
	return _PoolRegistry.Contract.AgentAt(&_PoolRegistry.CallOpts, _poolId, index)
}

// AgentExists is a free data retrieval call binding the contract method 0x839ec6fd.
//
// Solidity: function agentExists(bytes32 _poolId, bytes32 _agentId) view returns(bool)
func (_PoolRegistry *PoolRegistryCaller) AgentExists(opts *bind.CallOpts, _poolId [32]byte, _agentId [32]byte) (bool, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "agentExists", _poolId, _agentId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AgentExists is a free data retrieval call binding the contract method 0x839ec6fd.
//
// Solidity: function agentExists(bytes32 _poolId, bytes32 _agentId) view returns(bool)
func (_PoolRegistry *PoolRegistrySession) AgentExists(_poolId [32]byte, _agentId [32]byte) (bool, error) {
	return _PoolRegistry.Contract.AgentExists(&_PoolRegistry.CallOpts, _poolId, _agentId)
}

// AgentExists is a free data retrieval call binding the contract method 0x839ec6fd.
//
// Solidity: function agentExists(bytes32 _poolId, bytes32 _agentId) view returns(bool)
func (_PoolRegistry *PoolRegistryCallerSession) AgentExists(_poolId [32]byte, _agentId [32]byte) (bool, error) {
	return _PoolRegistry.Contract.AgentExists(&_PoolRegistry.CallOpts, _poolId, _agentId)
}

// AgentLength is a free data retrieval call binding the contract method 0x1e73feb1.
//
// Solidity: function agentLength(bytes32 _poolId) view returns(uint256)
func (_PoolRegistry *PoolRegistryCaller) AgentLength(opts *bind.CallOpts, _poolId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "agentLength", _poolId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentLength is a free data retrieval call binding the contract method 0x1e73feb1.
//
// Solidity: function agentLength(bytes32 _poolId) view returns(uint256)
func (_PoolRegistry *PoolRegistrySession) AgentLength(_poolId [32]byte) (*big.Int, error) {
	return _PoolRegistry.Contract.AgentLength(&_PoolRegistry.CallOpts, _poolId)
}

// AgentLength is a free data retrieval call binding the contract method 0x1e73feb1.
//
// Solidity: function agentLength(bytes32 _poolId) view returns(uint256)
func (_PoolRegistry *PoolRegistryCallerSession) AgentLength(_poolId [32]byte) (*big.Int, error) {
	return _PoolRegistry.Contract.AgentLength(&_PoolRegistry.CallOpts, _poolId)
}

// AgentRef is a free data retrieval call binding the contract method 0x6a2189d2.
//
// Solidity: function agentRef(bytes32 _poolId, bytes32 _agentId) view returns(string)
func (_PoolRegistry *PoolRegistryCaller) AgentRef(opts *bind.CallOpts, _poolId [32]byte, _agentId [32]byte) (string, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "agentRef", _poolId, _agentId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AgentRef is a free data retrieval call binding the contract method 0x6a2189d2.
//
// Solidity: function agentRef(bytes32 _poolId, bytes32 _agentId) view returns(string)
func (_PoolRegistry *PoolRegistrySession) AgentRef(_poolId [32]byte, _agentId [32]byte) (string, error) {
	return _PoolRegistry.Contract.AgentRef(&_PoolRegistry.CallOpts, _poolId, _agentId)
}

// AgentRef is a free data retrieval call binding the contract method 0x6a2189d2.
//
// Solidity: function agentRef(bytes32 _poolId, bytes32 _agentId) view returns(string)
func (_PoolRegistry *PoolRegistryCallerSession) AgentRef(_poolId [32]byte, _agentId [32]byte) (string, error) {
	return _PoolRegistry.Contract.AgentRef(&_PoolRegistry.CallOpts, _poolId, _agentId)
}

// AgentUsingLatest is a free data retrieval call binding the contract method 0x4edd5119.
//
// Solidity: function agentUsingLatest(bytes32 , bytes32 ) view returns(bool)
func (_PoolRegistry *PoolRegistryCaller) AgentUsingLatest(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) (bool, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "agentUsingLatest", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AgentUsingLatest is a free data retrieval call binding the contract method 0x4edd5119.
//
// Solidity: function agentUsingLatest(bytes32 , bytes32 ) view returns(bool)
func (_PoolRegistry *PoolRegistrySession) AgentUsingLatest(arg0 [32]byte, arg1 [32]byte) (bool, error) {
	return _PoolRegistry.Contract.AgentUsingLatest(&_PoolRegistry.CallOpts, arg0, arg1)
}

// AgentUsingLatest is a free data retrieval call binding the contract method 0x4edd5119.
//
// Solidity: function agentUsingLatest(bytes32 , bytes32 ) view returns(bool)
func (_PoolRegistry *PoolRegistryCallerSession) AgentUsingLatest(arg0 [32]byte, arg1 [32]byte) (bool, error) {
	return _PoolRegistry.Contract.AgentUsingLatest(&_PoolRegistry.CallOpts, arg0, arg1)
}

// AgentVersion is a free data retrieval call binding the contract method 0x021ef050.
//
// Solidity: function agentVersion(bytes32 , bytes32 ) view returns(uint256)
func (_PoolRegistry *PoolRegistryCaller) AgentVersion(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "agentVersion", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentVersion is a free data retrieval call binding the contract method 0x021ef050.
//
// Solidity: function agentVersion(bytes32 , bytes32 ) view returns(uint256)
func (_PoolRegistry *PoolRegistrySession) AgentVersion(arg0 [32]byte, arg1 [32]byte) (*big.Int, error) {
	return _PoolRegistry.Contract.AgentVersion(&_PoolRegistry.CallOpts, arg0, arg1)
}

// AgentVersion is a free data retrieval call binding the contract method 0x021ef050.
//
// Solidity: function agentVersion(bytes32 , bytes32 ) view returns(uint256)
func (_PoolRegistry *PoolRegistryCallerSession) AgentVersion(arg0 [32]byte, arg1 [32]byte) (*big.Int, error) {
	return _PoolRegistry.Contract.AgentVersion(&_PoolRegistry.CallOpts, arg0, arg1)
}

// GetPoolHash is a free data retrieval call binding the contract method 0xaf28390f.
//
// Solidity: function getPoolHash(bytes32 _poolId) view returns(bytes32)
func (_PoolRegistry *PoolRegistryCaller) GetPoolHash(opts *bind.CallOpts, _poolId [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "getPoolHash", _poolId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPoolHash is a free data retrieval call binding the contract method 0xaf28390f.
//
// Solidity: function getPoolHash(bytes32 _poolId) view returns(bytes32)
func (_PoolRegistry *PoolRegistrySession) GetPoolHash(_poolId [32]byte) ([32]byte, error) {
	return _PoolRegistry.Contract.GetPoolHash(&_PoolRegistry.CallOpts, _poolId)
}

// GetPoolHash is a free data retrieval call binding the contract method 0xaf28390f.
//
// Solidity: function getPoolHash(bytes32 _poolId) view returns(bytes32)
func (_PoolRegistry *PoolRegistryCallerSession) GetPoolHash(_poolId [32]byte) ([32]byte, error) {
	return _PoolRegistry.Contract.GetPoolHash(&_PoolRegistry.CallOpts, _poolId)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_PoolRegistry *PoolRegistryCaller) IsTrustedForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "isTrustedForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_PoolRegistry *PoolRegistrySession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _PoolRegistry.Contract.IsTrustedForwarder(&_PoolRegistry.CallOpts, forwarder)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_PoolRegistry *PoolRegistryCallerSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _PoolRegistry.Contract.IsTrustedForwarder(&_PoolRegistry.CallOpts, forwarder)
}

// PoolAdminAt is a free data retrieval call binding the contract method 0x43e32207.
//
// Solidity: function poolAdminAt(bytes32 _poolId, uint256 index) view returns(address)
func (_PoolRegistry *PoolRegistryCaller) PoolAdminAt(opts *bind.CallOpts, _poolId [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "poolAdminAt", _poolId, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolAdminAt is a free data retrieval call binding the contract method 0x43e32207.
//
// Solidity: function poolAdminAt(bytes32 _poolId, uint256 index) view returns(address)
func (_PoolRegistry *PoolRegistrySession) PoolAdminAt(_poolId [32]byte, index *big.Int) (common.Address, error) {
	return _PoolRegistry.Contract.PoolAdminAt(&_PoolRegistry.CallOpts, _poolId, index)
}

// PoolAdminAt is a free data retrieval call binding the contract method 0x43e32207.
//
// Solidity: function poolAdminAt(bytes32 _poolId, uint256 index) view returns(address)
func (_PoolRegistry *PoolRegistryCallerSession) PoolAdminAt(_poolId [32]byte, index *big.Int) (common.Address, error) {
	return _PoolRegistry.Contract.PoolAdminAt(&_PoolRegistry.CallOpts, _poolId, index)
}

// PoolAdminLength is a free data retrieval call binding the contract method 0xa0135277.
//
// Solidity: function poolAdminLength(bytes32 _poolId) view returns(uint256)
func (_PoolRegistry *PoolRegistryCaller) PoolAdminLength(opts *bind.CallOpts, _poolId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "poolAdminLength", _poolId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PoolAdminLength is a free data retrieval call binding the contract method 0xa0135277.
//
// Solidity: function poolAdminLength(bytes32 _poolId) view returns(uint256)
func (_PoolRegistry *PoolRegistrySession) PoolAdminLength(_poolId [32]byte) (*big.Int, error) {
	return _PoolRegistry.Contract.PoolAdminLength(&_PoolRegistry.CallOpts, _poolId)
}

// PoolAdminLength is a free data retrieval call binding the contract method 0xa0135277.
//
// Solidity: function poolAdminLength(bytes32 _poolId) view returns(uint256)
func (_PoolRegistry *PoolRegistryCallerSession) PoolAdminLength(_poolId [32]byte) (*big.Int, error) {
	return _PoolRegistry.Contract.PoolAdminLength(&_PoolRegistry.CallOpts, _poolId)
}

// PoolExistsMap is a free data retrieval call binding the contract method 0x0ead48e6.
//
// Solidity: function poolExistsMap(bytes32 ) view returns(bool)
func (_PoolRegistry *PoolRegistryCaller) PoolExistsMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "poolExistsMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PoolExistsMap is a free data retrieval call binding the contract method 0x0ead48e6.
//
// Solidity: function poolExistsMap(bytes32 ) view returns(bool)
func (_PoolRegistry *PoolRegistrySession) PoolExistsMap(arg0 [32]byte) (bool, error) {
	return _PoolRegistry.Contract.PoolExistsMap(&_PoolRegistry.CallOpts, arg0)
}

// PoolExistsMap is a free data retrieval call binding the contract method 0x0ead48e6.
//
// Solidity: function poolExistsMap(bytes32 ) view returns(bool)
func (_PoolRegistry *PoolRegistryCallerSession) PoolExistsMap(arg0 [32]byte) (bool, error) {
	return _PoolRegistry.Contract.PoolExistsMap(&_PoolRegistry.CallOpts, arg0)
}

// PoolOwners is a free data retrieval call binding the contract method 0x219d2d09.
//
// Solidity: function poolOwners(bytes32 ) view returns(address)
func (_PoolRegistry *PoolRegistryCaller) PoolOwners(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _PoolRegistry.contract.Call(opts, &out, "poolOwners", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolOwners is a free data retrieval call binding the contract method 0x219d2d09.
//
// Solidity: function poolOwners(bytes32 ) view returns(address)
func (_PoolRegistry *PoolRegistrySession) PoolOwners(arg0 [32]byte) (common.Address, error) {
	return _PoolRegistry.Contract.PoolOwners(&_PoolRegistry.CallOpts, arg0)
}

// PoolOwners is a free data retrieval call binding the contract method 0x219d2d09.
//
// Solidity: function poolOwners(bytes32 ) view returns(address)
func (_PoolRegistry *PoolRegistryCallerSession) PoolOwners(arg0 [32]byte) (common.Address, error) {
	return _PoolRegistry.Contract.PoolOwners(&_PoolRegistry.CallOpts, arg0)
}

// AddAgent is a paid mutator transaction binding the contract method 0x9814a81e.
//
// Solidity: function addAgent(bytes32 _poolId, bytes32 _agentId, uint256 version, bool latest) returns()
func (_PoolRegistry *PoolRegistryTransactor) AddAgent(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "addAgent", _poolId, _agentId, version, latest)
}

// AddAgent is a paid mutator transaction binding the contract method 0x9814a81e.
//
// Solidity: function addAgent(bytes32 _poolId, bytes32 _agentId, uint256 version, bool latest) returns()
func (_PoolRegistry *PoolRegistrySession) AddAgent(_poolId [32]byte, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _PoolRegistry.Contract.AddAgent(&_PoolRegistry.TransactOpts, _poolId, _agentId, version, latest)
}

// AddAgent is a paid mutator transaction binding the contract method 0x9814a81e.
//
// Solidity: function addAgent(bytes32 _poolId, bytes32 _agentId, uint256 version, bool latest) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) AddAgent(_poolId [32]byte, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _PoolRegistry.Contract.AddAgent(&_PoolRegistry.TransactOpts, _poolId, _agentId, version, latest)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x63438a97.
//
// Solidity: function addAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_PoolRegistry *PoolRegistryTransactor) AddAgentAdmin(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "addAgentAdmin", _poolId, _agentId, _admin)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x63438a97.
//
// Solidity: function addAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_PoolRegistry *PoolRegistrySession) AddAgentAdmin(_poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.AddAgentAdmin(&_PoolRegistry.TransactOpts, _poolId, _agentId, _admin)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x63438a97.
//
// Solidity: function addAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) AddAgentAdmin(_poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.AddAgentAdmin(&_PoolRegistry.TransactOpts, _poolId, _agentId, _admin)
}

// AddPool is a paid mutator transaction binding the contract method 0x8a0033fc.
//
// Solidity: function addPool(bytes32 _poolId) returns()
func (_PoolRegistry *PoolRegistryTransactor) AddPool(opts *bind.TransactOpts, _poolId [32]byte) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "addPool", _poolId)
}

// AddPool is a paid mutator transaction binding the contract method 0x8a0033fc.
//
// Solidity: function addPool(bytes32 _poolId) returns()
func (_PoolRegistry *PoolRegistrySession) AddPool(_poolId [32]byte) (*types.Transaction, error) {
	return _PoolRegistry.Contract.AddPool(&_PoolRegistry.TransactOpts, _poolId)
}

// AddPool is a paid mutator transaction binding the contract method 0x8a0033fc.
//
// Solidity: function addPool(bytes32 _poolId) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) AddPool(_poolId [32]byte) (*types.Transaction, error) {
	return _PoolRegistry.Contract.AddPool(&_PoolRegistry.TransactOpts, _poolId)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x84930b76.
//
// Solidity: function addPoolAdmin(bytes32 _poolId, address admin) returns()
func (_PoolRegistry *PoolRegistryTransactor) AddPoolAdmin(opts *bind.TransactOpts, _poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "addPoolAdmin", _poolId, admin)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x84930b76.
//
// Solidity: function addPoolAdmin(bytes32 _poolId, address admin) returns()
func (_PoolRegistry *PoolRegistrySession) AddPoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.AddPoolAdmin(&_PoolRegistry.TransactOpts, _poolId, admin)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x84930b76.
//
// Solidity: function addPoolAdmin(bytes32 _poolId, address admin) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) AddPoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.AddPoolAdmin(&_PoolRegistry.TransactOpts, _poolId, admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address forwarder) returns()
func (_PoolRegistry *PoolRegistryTransactor) Initialize(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "initialize", forwarder)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address forwarder) returns()
func (_PoolRegistry *PoolRegistrySession) Initialize(forwarder common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.Initialize(&_PoolRegistry.TransactOpts, forwarder)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address forwarder) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) Initialize(forwarder common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.Initialize(&_PoolRegistry.TransactOpts, forwarder)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x61e62a0f.
//
// Solidity: function removeAgent(bytes32 _poolId, bytes32 _agentId) returns()
func (_PoolRegistry *PoolRegistryTransactor) RemoveAgent(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "removeAgent", _poolId, _agentId)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x61e62a0f.
//
// Solidity: function removeAgent(bytes32 _poolId, bytes32 _agentId) returns()
func (_PoolRegistry *PoolRegistrySession) RemoveAgent(_poolId [32]byte, _agentId [32]byte) (*types.Transaction, error) {
	return _PoolRegistry.Contract.RemoveAgent(&_PoolRegistry.TransactOpts, _poolId, _agentId)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x61e62a0f.
//
// Solidity: function removeAgent(bytes32 _poolId, bytes32 _agentId) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) RemoveAgent(_poolId [32]byte, _agentId [32]byte) (*types.Transaction, error) {
	return _PoolRegistry.Contract.RemoveAgent(&_PoolRegistry.TransactOpts, _poolId, _agentId)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0xa7f2f5ab.
//
// Solidity: function removeAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_PoolRegistry *PoolRegistryTransactor) RemoveAgentAdmin(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "removeAgentAdmin", _poolId, _agentId, _admin)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0xa7f2f5ab.
//
// Solidity: function removeAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_PoolRegistry *PoolRegistrySession) RemoveAgentAdmin(_poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.RemoveAgentAdmin(&_PoolRegistry.TransactOpts, _poolId, _agentId, _admin)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0xa7f2f5ab.
//
// Solidity: function removeAgentAdmin(bytes32 _poolId, bytes32 _agentId, address _admin) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) RemoveAgentAdmin(_poolId [32]byte, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.RemoveAgentAdmin(&_PoolRegistry.TransactOpts, _poolId, _agentId, _admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xacce5630.
//
// Solidity: function removePoolAdmin(bytes32 _poolId, address admin) returns()
func (_PoolRegistry *PoolRegistryTransactor) RemovePoolAdmin(opts *bind.TransactOpts, _poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "removePoolAdmin", _poolId, admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xacce5630.
//
// Solidity: function removePoolAdmin(bytes32 _poolId, address admin) returns()
func (_PoolRegistry *PoolRegistrySession) RemovePoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.RemovePoolAdmin(&_PoolRegistry.TransactOpts, _poolId, admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xacce5630.
//
// Solidity: function removePoolAdmin(bytes32 _poolId, address admin) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) RemovePoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.RemovePoolAdmin(&_PoolRegistry.TransactOpts, _poolId, admin)
}

// SetAgentRegistry is a paid mutator transaction binding the contract method 0x28342ecf.
//
// Solidity: function setAgentRegistry(address _addr) returns()
func (_PoolRegistry *PoolRegistryTransactor) SetAgentRegistry(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "setAgentRegistry", _addr)
}

// SetAgentRegistry is a paid mutator transaction binding the contract method 0x28342ecf.
//
// Solidity: function setAgentRegistry(address _addr) returns()
func (_PoolRegistry *PoolRegistrySession) SetAgentRegistry(_addr common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.SetAgentRegistry(&_PoolRegistry.TransactOpts, _addr)
}

// SetAgentRegistry is a paid mutator transaction binding the contract method 0x28342ecf.
//
// Solidity: function setAgentRegistry(address _addr) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) SetAgentRegistry(_addr common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.SetAgentRegistry(&_PoolRegistry.TransactOpts, _addr)
}

// TransferPoolOwnership is a paid mutator transaction binding the contract method 0x1b6feb92.
//
// Solidity: function transferPoolOwnership(bytes32 _poolId, address _to) returns()
func (_PoolRegistry *PoolRegistryTransactor) TransferPoolOwnership(opts *bind.TransactOpts, _poolId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "transferPoolOwnership", _poolId, _to)
}

// TransferPoolOwnership is a paid mutator transaction binding the contract method 0x1b6feb92.
//
// Solidity: function transferPoolOwnership(bytes32 _poolId, address _to) returns()
func (_PoolRegistry *PoolRegistrySession) TransferPoolOwnership(_poolId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.TransferPoolOwnership(&_PoolRegistry.TransactOpts, _poolId, _to)
}

// TransferPoolOwnership is a paid mutator transaction binding the contract method 0x1b6feb92.
//
// Solidity: function transferPoolOwnership(bytes32 _poolId, address _to) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) TransferPoolOwnership(_poolId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _PoolRegistry.Contract.TransferPoolOwnership(&_PoolRegistry.TransactOpts, _poolId, _to)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xf587fecb.
//
// Solidity: function updateAgent(bytes32 _poolId, bytes32 _agentId, uint256 version, bool latest) returns()
func (_PoolRegistry *PoolRegistryTransactor) UpdateAgent(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _PoolRegistry.contract.Transact(opts, "updateAgent", _poolId, _agentId, version, latest)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xf587fecb.
//
// Solidity: function updateAgent(bytes32 _poolId, bytes32 _agentId, uint256 version, bool latest) returns()
func (_PoolRegistry *PoolRegistrySession) UpdateAgent(_poolId [32]byte, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _PoolRegistry.Contract.UpdateAgent(&_PoolRegistry.TransactOpts, _poolId, _agentId, version, latest)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xf587fecb.
//
// Solidity: function updateAgent(bytes32 _poolId, bytes32 _agentId, uint256 version, bool latest) returns()
func (_PoolRegistry *PoolRegistryTransactorSession) UpdateAgent(_poolId [32]byte, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _PoolRegistry.Contract.UpdateAgent(&_PoolRegistry.TransactOpts, _poolId, _agentId, version, latest)
}

// PoolRegistryAgentAddedIterator is returned from FilterAgentAdded and is used to iterate over the raw logs and unpacked data for AgentAdded events raised by the PoolRegistry contract.
type PoolRegistryAgentAddedIterator struct {
	Event *PoolRegistryAgentAdded // Event containing the contract specifics and raw log

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
func (it *PoolRegistryAgentAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRegistryAgentAdded)
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
		it.Event = new(PoolRegistryAgentAdded)
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
func (it *PoolRegistryAgentAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRegistryAgentAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRegistryAgentAdded represents a AgentAdded event raised by the PoolRegistry contract.
type PoolRegistryAgentAdded struct {
	PoolId  [32]byte
	AgentId [32]byte
	Ref     string
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdded is a free log retrieval operation binding the contract event 0xdd1417d47ca73bdaf24e5c1df4158e61cbbb73fda6562b14a16215ccfa201aeb.
//
// Solidity: event AgentAdded(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_PoolRegistry *PoolRegistryFilterer) FilterAgentAdded(opts *bind.FilterOpts) (*PoolRegistryAgentAddedIterator, error) {

	logs, sub, err := _PoolRegistry.contract.FilterLogs(opts, "AgentAdded")
	if err != nil {
		return nil, err
	}
	return &PoolRegistryAgentAddedIterator{contract: _PoolRegistry.contract, event: "AgentAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdded is a free log subscription operation binding the contract event 0xdd1417d47ca73bdaf24e5c1df4158e61cbbb73fda6562b14a16215ccfa201aeb.
//
// Solidity: event AgentAdded(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_PoolRegistry *PoolRegistryFilterer) WatchAgentAdded(opts *bind.WatchOpts, sink chan<- *PoolRegistryAgentAdded) (event.Subscription, error) {

	logs, sub, err := _PoolRegistry.contract.WatchLogs(opts, "AgentAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRegistryAgentAdded)
				if err := _PoolRegistry.contract.UnpackLog(event, "AgentAdded", log); err != nil {
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
func (_PoolRegistry *PoolRegistryFilterer) ParseAgentAdded(log types.Log) (*PoolRegistryAgentAdded, error) {
	event := new(PoolRegistryAgentAdded)
	if err := _PoolRegistry.contract.UnpackLog(event, "AgentAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRegistryAgentAdminAddedIterator is returned from FilterAgentAdminAdded and is used to iterate over the raw logs and unpacked data for AgentAdminAdded events raised by the PoolRegistry contract.
type PoolRegistryAgentAdminAddedIterator struct {
	Event *PoolRegistryAgentAdminAdded // Event containing the contract specifics and raw log

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
func (it *PoolRegistryAgentAdminAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRegistryAgentAdminAdded)
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
		it.Event = new(PoolRegistryAgentAdminAdded)
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
func (it *PoolRegistryAgentAdminAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRegistryAgentAdminAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRegistryAgentAdminAdded represents a AgentAdminAdded event raised by the PoolRegistry contract.
type PoolRegistryAgentAdminAdded struct {
	PoolId  [32]byte
	AgentId [32]byte
	Admin   common.Address
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdminAdded is a free log retrieval operation binding the contract event 0x11cbfe7067a564a8ddb9a780914987915bc0d29835c8b0554a9d9db2e51064ed.
//
// Solidity: event AgentAdminAdded(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_PoolRegistry *PoolRegistryFilterer) FilterAgentAdminAdded(opts *bind.FilterOpts) (*PoolRegistryAgentAdminAddedIterator, error) {

	logs, sub, err := _PoolRegistry.contract.FilterLogs(opts, "AgentAdminAdded")
	if err != nil {
		return nil, err
	}
	return &PoolRegistryAgentAdminAddedIterator{contract: _PoolRegistry.contract, event: "AgentAdminAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdminAdded is a free log subscription operation binding the contract event 0x11cbfe7067a564a8ddb9a780914987915bc0d29835c8b0554a9d9db2e51064ed.
//
// Solidity: event AgentAdminAdded(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_PoolRegistry *PoolRegistryFilterer) WatchAgentAdminAdded(opts *bind.WatchOpts, sink chan<- *PoolRegistryAgentAdminAdded) (event.Subscription, error) {

	logs, sub, err := _PoolRegistry.contract.WatchLogs(opts, "AgentAdminAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRegistryAgentAdminAdded)
				if err := _PoolRegistry.contract.UnpackLog(event, "AgentAdminAdded", log); err != nil {
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
func (_PoolRegistry *PoolRegistryFilterer) ParseAgentAdminAdded(log types.Log) (*PoolRegistryAgentAdminAdded, error) {
	event := new(PoolRegistryAgentAdminAdded)
	if err := _PoolRegistry.contract.UnpackLog(event, "AgentAdminAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRegistryAgentAdminRemovedIterator is returned from FilterAgentAdminRemoved and is used to iterate over the raw logs and unpacked data for AgentAdminRemoved events raised by the PoolRegistry contract.
type PoolRegistryAgentAdminRemovedIterator struct {
	Event *PoolRegistryAgentAdminRemoved // Event containing the contract specifics and raw log

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
func (it *PoolRegistryAgentAdminRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRegistryAgentAdminRemoved)
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
		it.Event = new(PoolRegistryAgentAdminRemoved)
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
func (it *PoolRegistryAgentAdminRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRegistryAgentAdminRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRegistryAgentAdminRemoved represents a AgentAdminRemoved event raised by the PoolRegistry contract.
type PoolRegistryAgentAdminRemoved struct {
	PoolId  [32]byte
	AgentId [32]byte
	Admin   common.Address
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdminRemoved is a free log retrieval operation binding the contract event 0x801fe4d9d51c0445c2938970ad5e042c3f2f2c387d67de4111af1bc60d22dba9.
//
// Solidity: event AgentAdminRemoved(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_PoolRegistry *PoolRegistryFilterer) FilterAgentAdminRemoved(opts *bind.FilterOpts) (*PoolRegistryAgentAdminRemovedIterator, error) {

	logs, sub, err := _PoolRegistry.contract.FilterLogs(opts, "AgentAdminRemoved")
	if err != nil {
		return nil, err
	}
	return &PoolRegistryAgentAdminRemovedIterator{contract: _PoolRegistry.contract, event: "AgentAdminRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentAdminRemoved is a free log subscription operation binding the contract event 0x801fe4d9d51c0445c2938970ad5e042c3f2f2c387d67de4111af1bc60d22dba9.
//
// Solidity: event AgentAdminRemoved(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_PoolRegistry *PoolRegistryFilterer) WatchAgentAdminRemoved(opts *bind.WatchOpts, sink chan<- *PoolRegistryAgentAdminRemoved) (event.Subscription, error) {

	logs, sub, err := _PoolRegistry.contract.WatchLogs(opts, "AgentAdminRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRegistryAgentAdminRemoved)
				if err := _PoolRegistry.contract.UnpackLog(event, "AgentAdminRemoved", log); err != nil {
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
func (_PoolRegistry *PoolRegistryFilterer) ParseAgentAdminRemoved(log types.Log) (*PoolRegistryAgentAdminRemoved, error) {
	event := new(PoolRegistryAgentAdminRemoved)
	if err := _PoolRegistry.contract.UnpackLog(event, "AgentAdminRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRegistryAgentRegistryChangedIterator is returned from FilterAgentRegistryChanged and is used to iterate over the raw logs and unpacked data for AgentRegistryChanged events raised by the PoolRegistry contract.
type PoolRegistryAgentRegistryChangedIterator struct {
	Event *PoolRegistryAgentRegistryChanged // Event containing the contract specifics and raw log

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
func (it *PoolRegistryAgentRegistryChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRegistryAgentRegistryChanged)
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
		it.Event = new(PoolRegistryAgentRegistryChanged)
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
func (it *PoolRegistryAgentRegistryChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRegistryAgentRegistryChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRegistryAgentRegistryChanged represents a AgentRegistryChanged event raised by the PoolRegistry contract.
type PoolRegistryAgentRegistryChanged struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAgentRegistryChanged is a free log retrieval operation binding the contract event 0x2d35b9d2e073404ea7e01756776790bfe7c7e525789490689493de40955f5632.
//
// Solidity: event AgentRegistryChanged(address from, address to)
func (_PoolRegistry *PoolRegistryFilterer) FilterAgentRegistryChanged(opts *bind.FilterOpts) (*PoolRegistryAgentRegistryChangedIterator, error) {

	logs, sub, err := _PoolRegistry.contract.FilterLogs(opts, "AgentRegistryChanged")
	if err != nil {
		return nil, err
	}
	return &PoolRegistryAgentRegistryChangedIterator{contract: _PoolRegistry.contract, event: "AgentRegistryChanged", logs: logs, sub: sub}, nil
}

// WatchAgentRegistryChanged is a free log subscription operation binding the contract event 0x2d35b9d2e073404ea7e01756776790bfe7c7e525789490689493de40955f5632.
//
// Solidity: event AgentRegistryChanged(address from, address to)
func (_PoolRegistry *PoolRegistryFilterer) WatchAgentRegistryChanged(opts *bind.WatchOpts, sink chan<- *PoolRegistryAgentRegistryChanged) (event.Subscription, error) {

	logs, sub, err := _PoolRegistry.contract.WatchLogs(opts, "AgentRegistryChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRegistryAgentRegistryChanged)
				if err := _PoolRegistry.contract.UnpackLog(event, "AgentRegistryChanged", log); err != nil {
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

// ParseAgentRegistryChanged is a log parse operation binding the contract event 0x2d35b9d2e073404ea7e01756776790bfe7c7e525789490689493de40955f5632.
//
// Solidity: event AgentRegistryChanged(address from, address to)
func (_PoolRegistry *PoolRegistryFilterer) ParseAgentRegistryChanged(log types.Log) (*PoolRegistryAgentRegistryChanged, error) {
	event := new(PoolRegistryAgentRegistryChanged)
	if err := _PoolRegistry.contract.UnpackLog(event, "AgentRegistryChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRegistryAgentRemovedIterator is returned from FilterAgentRemoved and is used to iterate over the raw logs and unpacked data for AgentRemoved events raised by the PoolRegistry contract.
type PoolRegistryAgentRemovedIterator struct {
	Event *PoolRegistryAgentRemoved // Event containing the contract specifics and raw log

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
func (it *PoolRegistryAgentRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRegistryAgentRemoved)
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
		it.Event = new(PoolRegistryAgentRemoved)
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
func (it *PoolRegistryAgentRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRegistryAgentRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRegistryAgentRemoved represents a AgentRemoved event raised by the PoolRegistry contract.
type PoolRegistryAgentRemoved struct {
	PoolId  [32]byte
	AgentId [32]byte
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentRemoved is a free log retrieval operation binding the contract event 0xccc61238b7da918476155ec78f1d6f23cd587a3f796774b80c0c79dde8f57e6b.
//
// Solidity: event AgentRemoved(bytes32 poolId, bytes32 agentId, address by)
func (_PoolRegistry *PoolRegistryFilterer) FilterAgentRemoved(opts *bind.FilterOpts) (*PoolRegistryAgentRemovedIterator, error) {

	logs, sub, err := _PoolRegistry.contract.FilterLogs(opts, "AgentRemoved")
	if err != nil {
		return nil, err
	}
	return &PoolRegistryAgentRemovedIterator{contract: _PoolRegistry.contract, event: "AgentRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentRemoved is a free log subscription operation binding the contract event 0xccc61238b7da918476155ec78f1d6f23cd587a3f796774b80c0c79dde8f57e6b.
//
// Solidity: event AgentRemoved(bytes32 poolId, bytes32 agentId, address by)
func (_PoolRegistry *PoolRegistryFilterer) WatchAgentRemoved(opts *bind.WatchOpts, sink chan<- *PoolRegistryAgentRemoved) (event.Subscription, error) {

	logs, sub, err := _PoolRegistry.contract.WatchLogs(opts, "AgentRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRegistryAgentRemoved)
				if err := _PoolRegistry.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
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
func (_PoolRegistry *PoolRegistryFilterer) ParseAgentRemoved(log types.Log) (*PoolRegistryAgentRemoved, error) {
	event := new(PoolRegistryAgentRemoved)
	if err := _PoolRegistry.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRegistryAgentUpdatedIterator is returned from FilterAgentUpdated and is used to iterate over the raw logs and unpacked data for AgentUpdated events raised by the PoolRegistry contract.
type PoolRegistryAgentUpdatedIterator struct {
	Event *PoolRegistryAgentUpdated // Event containing the contract specifics and raw log

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
func (it *PoolRegistryAgentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRegistryAgentUpdated)
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
		it.Event = new(PoolRegistryAgentUpdated)
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
func (it *PoolRegistryAgentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRegistryAgentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRegistryAgentUpdated represents a AgentUpdated event raised by the PoolRegistry contract.
type PoolRegistryAgentUpdated struct {
	PoolId  [32]byte
	AgentId [32]byte
	Ref     string
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentUpdated is a free log retrieval operation binding the contract event 0x8a33912f31f80df098dcad04aba17afc8c52d0f700baf2750ea3899f9230a850.
//
// Solidity: event AgentUpdated(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_PoolRegistry *PoolRegistryFilterer) FilterAgentUpdated(opts *bind.FilterOpts) (*PoolRegistryAgentUpdatedIterator, error) {

	logs, sub, err := _PoolRegistry.contract.FilterLogs(opts, "AgentUpdated")
	if err != nil {
		return nil, err
	}
	return &PoolRegistryAgentUpdatedIterator{contract: _PoolRegistry.contract, event: "AgentUpdated", logs: logs, sub: sub}, nil
}

// WatchAgentUpdated is a free log subscription operation binding the contract event 0x8a33912f31f80df098dcad04aba17afc8c52d0f700baf2750ea3899f9230a850.
//
// Solidity: event AgentUpdated(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_PoolRegistry *PoolRegistryFilterer) WatchAgentUpdated(opts *bind.WatchOpts, sink chan<- *PoolRegistryAgentUpdated) (event.Subscription, error) {

	logs, sub, err := _PoolRegistry.contract.WatchLogs(opts, "AgentUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRegistryAgentUpdated)
				if err := _PoolRegistry.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
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
func (_PoolRegistry *PoolRegistryFilterer) ParseAgentUpdated(log types.Log) (*PoolRegistryAgentUpdated, error) {
	event := new(PoolRegistryAgentUpdated)
	if err := _PoolRegistry.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRegistryPoolAddedIterator is returned from FilterPoolAdded and is used to iterate over the raw logs and unpacked data for PoolAdded events raised by the PoolRegistry contract.
type PoolRegistryPoolAddedIterator struct {
	Event *PoolRegistryPoolAdded // Event containing the contract specifics and raw log

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
func (it *PoolRegistryPoolAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRegistryPoolAdded)
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
		it.Event = new(PoolRegistryPoolAdded)
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
func (it *PoolRegistryPoolAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRegistryPoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRegistryPoolAdded represents a PoolAdded event raised by the PoolRegistry contract.
type PoolRegistryPoolAdded struct {
	PoolId [32]byte
	By     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolAdded is a free log retrieval operation binding the contract event 0xfd0fa7919fbe3857a4236750e8d3e42ac691881d2f31c50ebe4cfc2a7705ee80.
//
// Solidity: event PoolAdded(bytes32 poolId, address by)
func (_PoolRegistry *PoolRegistryFilterer) FilterPoolAdded(opts *bind.FilterOpts) (*PoolRegistryPoolAddedIterator, error) {

	logs, sub, err := _PoolRegistry.contract.FilterLogs(opts, "PoolAdded")
	if err != nil {
		return nil, err
	}
	return &PoolRegistryPoolAddedIterator{contract: _PoolRegistry.contract, event: "PoolAdded", logs: logs, sub: sub}, nil
}

// WatchPoolAdded is a free log subscription operation binding the contract event 0xfd0fa7919fbe3857a4236750e8d3e42ac691881d2f31c50ebe4cfc2a7705ee80.
//
// Solidity: event PoolAdded(bytes32 poolId, address by)
func (_PoolRegistry *PoolRegistryFilterer) WatchPoolAdded(opts *bind.WatchOpts, sink chan<- *PoolRegistryPoolAdded) (event.Subscription, error) {

	logs, sub, err := _PoolRegistry.contract.WatchLogs(opts, "PoolAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRegistryPoolAdded)
				if err := _PoolRegistry.contract.UnpackLog(event, "PoolAdded", log); err != nil {
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
func (_PoolRegistry *PoolRegistryFilterer) ParsePoolAdded(log types.Log) (*PoolRegistryPoolAdded, error) {
	event := new(PoolRegistryPoolAdded)
	if err := _PoolRegistry.contract.UnpackLog(event, "PoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRegistryPoolAdminAddedIterator is returned from FilterPoolAdminAdded and is used to iterate over the raw logs and unpacked data for PoolAdminAdded events raised by the PoolRegistry contract.
type PoolRegistryPoolAdminAddedIterator struct {
	Event *PoolRegistryPoolAdminAdded // Event containing the contract specifics and raw log

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
func (it *PoolRegistryPoolAdminAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRegistryPoolAdminAdded)
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
		it.Event = new(PoolRegistryPoolAdminAdded)
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
func (it *PoolRegistryPoolAdminAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRegistryPoolAdminAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRegistryPoolAdminAdded represents a PoolAdminAdded event raised by the PoolRegistry contract.
type PoolRegistryPoolAdminAdded struct {
	PoolId [32]byte
	Addr   common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolAdminAdded is a free log retrieval operation binding the contract event 0x4c2f0be874d7047c389e9fadf454854a8f5b72b49de199c7d60831b6c462a9a0.
//
// Solidity: event PoolAdminAdded(bytes32 poolId, address addr)
func (_PoolRegistry *PoolRegistryFilterer) FilterPoolAdminAdded(opts *bind.FilterOpts) (*PoolRegistryPoolAdminAddedIterator, error) {

	logs, sub, err := _PoolRegistry.contract.FilterLogs(opts, "PoolAdminAdded")
	if err != nil {
		return nil, err
	}
	return &PoolRegistryPoolAdminAddedIterator{contract: _PoolRegistry.contract, event: "PoolAdminAdded", logs: logs, sub: sub}, nil
}

// WatchPoolAdminAdded is a free log subscription operation binding the contract event 0x4c2f0be874d7047c389e9fadf454854a8f5b72b49de199c7d60831b6c462a9a0.
//
// Solidity: event PoolAdminAdded(bytes32 poolId, address addr)
func (_PoolRegistry *PoolRegistryFilterer) WatchPoolAdminAdded(opts *bind.WatchOpts, sink chan<- *PoolRegistryPoolAdminAdded) (event.Subscription, error) {

	logs, sub, err := _PoolRegistry.contract.WatchLogs(opts, "PoolAdminAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRegistryPoolAdminAdded)
				if err := _PoolRegistry.contract.UnpackLog(event, "PoolAdminAdded", log); err != nil {
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
func (_PoolRegistry *PoolRegistryFilterer) ParsePoolAdminAdded(log types.Log) (*PoolRegistryPoolAdminAdded, error) {
	event := new(PoolRegistryPoolAdminAdded)
	if err := _PoolRegistry.contract.UnpackLog(event, "PoolAdminAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRegistryPoolAdminRemovedIterator is returned from FilterPoolAdminRemoved and is used to iterate over the raw logs and unpacked data for PoolAdminRemoved events raised by the PoolRegistry contract.
type PoolRegistryPoolAdminRemovedIterator struct {
	Event *PoolRegistryPoolAdminRemoved // Event containing the contract specifics and raw log

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
func (it *PoolRegistryPoolAdminRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRegistryPoolAdminRemoved)
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
		it.Event = new(PoolRegistryPoolAdminRemoved)
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
func (it *PoolRegistryPoolAdminRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRegistryPoolAdminRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRegistryPoolAdminRemoved represents a PoolAdminRemoved event raised by the PoolRegistry contract.
type PoolRegistryPoolAdminRemoved struct {
	PoolId [32]byte
	Addr   common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolAdminRemoved is a free log retrieval operation binding the contract event 0x25fcb4cb225521d9a561fd90af5809d997ae0fa769a5c95d90750162623bf624.
//
// Solidity: event PoolAdminRemoved(bytes32 poolId, address addr)
func (_PoolRegistry *PoolRegistryFilterer) FilterPoolAdminRemoved(opts *bind.FilterOpts) (*PoolRegistryPoolAdminRemovedIterator, error) {

	logs, sub, err := _PoolRegistry.contract.FilterLogs(opts, "PoolAdminRemoved")
	if err != nil {
		return nil, err
	}
	return &PoolRegistryPoolAdminRemovedIterator{contract: _PoolRegistry.contract, event: "PoolAdminRemoved", logs: logs, sub: sub}, nil
}

// WatchPoolAdminRemoved is a free log subscription operation binding the contract event 0x25fcb4cb225521d9a561fd90af5809d997ae0fa769a5c95d90750162623bf624.
//
// Solidity: event PoolAdminRemoved(bytes32 poolId, address addr)
func (_PoolRegistry *PoolRegistryFilterer) WatchPoolAdminRemoved(opts *bind.WatchOpts, sink chan<- *PoolRegistryPoolAdminRemoved) (event.Subscription, error) {

	logs, sub, err := _PoolRegistry.contract.WatchLogs(opts, "PoolAdminRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRegistryPoolAdminRemoved)
				if err := _PoolRegistry.contract.UnpackLog(event, "PoolAdminRemoved", log); err != nil {
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
func (_PoolRegistry *PoolRegistryFilterer) ParsePoolAdminRemoved(log types.Log) (*PoolRegistryPoolAdminRemoved, error) {
	event := new(PoolRegistryPoolAdminRemoved)
	if err := _PoolRegistry.contract.UnpackLog(event, "PoolAdminRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRegistryPoolOwnershipTransferedIterator is returned from FilterPoolOwnershipTransfered and is used to iterate over the raw logs and unpacked data for PoolOwnershipTransfered events raised by the PoolRegistry contract.
type PoolRegistryPoolOwnershipTransferedIterator struct {
	Event *PoolRegistryPoolOwnershipTransfered // Event containing the contract specifics and raw log

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
func (it *PoolRegistryPoolOwnershipTransferedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRegistryPoolOwnershipTransfered)
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
		it.Event = new(PoolRegistryPoolOwnershipTransfered)
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
func (it *PoolRegistryPoolOwnershipTransferedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRegistryPoolOwnershipTransferedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRegistryPoolOwnershipTransfered represents a PoolOwnershipTransfered event raised by the PoolRegistry contract.
type PoolRegistryPoolOwnershipTransfered struct {
	PoolId [32]byte
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolOwnershipTransfered is a free log retrieval operation binding the contract event 0xecc1c418cc7d2c3062b60c1d8f0b7748e2c8c10a109df4f6089b285f26b993d6.
//
// Solidity: event PoolOwnershipTransfered(bytes32 poolId, address from, address to)
func (_PoolRegistry *PoolRegistryFilterer) FilterPoolOwnershipTransfered(opts *bind.FilterOpts) (*PoolRegistryPoolOwnershipTransferedIterator, error) {

	logs, sub, err := _PoolRegistry.contract.FilterLogs(opts, "PoolOwnershipTransfered")
	if err != nil {
		return nil, err
	}
	return &PoolRegistryPoolOwnershipTransferedIterator{contract: _PoolRegistry.contract, event: "PoolOwnershipTransfered", logs: logs, sub: sub}, nil
}

// WatchPoolOwnershipTransfered is a free log subscription operation binding the contract event 0xecc1c418cc7d2c3062b60c1d8f0b7748e2c8c10a109df4f6089b285f26b993d6.
//
// Solidity: event PoolOwnershipTransfered(bytes32 poolId, address from, address to)
func (_PoolRegistry *PoolRegistryFilterer) WatchPoolOwnershipTransfered(opts *bind.WatchOpts, sink chan<- *PoolRegistryPoolOwnershipTransfered) (event.Subscription, error) {

	logs, sub, err := _PoolRegistry.contract.WatchLogs(opts, "PoolOwnershipTransfered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRegistryPoolOwnershipTransfered)
				if err := _PoolRegistry.contract.UnpackLog(event, "PoolOwnershipTransfered", log); err != nil {
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
func (_PoolRegistry *PoolRegistryFilterer) ParsePoolOwnershipTransfered(log types.Log) (*PoolRegistryPoolOwnershipTransfered, error) {
	event := new(PoolRegistryPoolOwnershipTransfered)
	if err := _PoolRegistry.contract.UnpackLog(event, "PoolOwnershipTransfered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
