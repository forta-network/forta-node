// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AgentRegistryMetaData contains all meta data concerning the AgentRegistry contract.
var AgentRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dev\",\"type\":\"address\"}],\"name\":\"AgentDevAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dev\",\"type\":\"address\"}],\"name\":\"AgentDevRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"AgentOwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentVersionPublished\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_dev\",\"type\":\"address\"}],\"name\":\"addAgentDev\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"agentDisabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"agentExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"agentLatestVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"agentReference\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_ref\",\"type\":\"string\"}],\"name\":\"createAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"disableAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"enableAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"getLatestRef\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractMinimalForwarderUpgradeable\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_dev\",\"type\":\"address\"}],\"name\":\"removeAgentDev\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferAgentOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_ref\",\"type\":\"string\"}],\"name\":\"updateAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AgentRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use AgentRegistryMetaData.ABI instead.
var AgentRegistryABI = AgentRegistryMetaData.ABI

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

// AgentDisabled is a free data retrieval call binding the contract method 0x0d8cb278.
//
// Solidity: function agentDisabled(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) AgentDisabled(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentDisabled", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AgentDisabled is a free data retrieval call binding the contract method 0x0d8cb278.
//
// Solidity: function agentDisabled(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) AgentDisabled(arg0 [32]byte) (bool, error) {
	return _AgentRegistry.Contract.AgentDisabled(&_AgentRegistry.CallOpts, arg0)
}

// AgentDisabled is a free data retrieval call binding the contract method 0x0d8cb278.
//
// Solidity: function agentDisabled(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) AgentDisabled(arg0 [32]byte) (bool, error) {
	return _AgentRegistry.Contract.AgentDisabled(&_AgentRegistry.CallOpts, arg0)
}

// AgentExists is a free data retrieval call binding the contract method 0xa5b925d2.
//
// Solidity: function agentExists(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) AgentExists(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentExists", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AgentExists is a free data retrieval call binding the contract method 0xa5b925d2.
//
// Solidity: function agentExists(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) AgentExists(arg0 [32]byte) (bool, error) {
	return _AgentRegistry.Contract.AgentExists(&_AgentRegistry.CallOpts, arg0)
}

// AgentExists is a free data retrieval call binding the contract method 0xa5b925d2.
//
// Solidity: function agentExists(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) AgentExists(arg0 [32]byte) (bool, error) {
	return _AgentRegistry.Contract.AgentExists(&_AgentRegistry.CallOpts, arg0)
}

// AgentLatestVersion is a free data retrieval call binding the contract method 0x380619bb.
//
// Solidity: function agentLatestVersion(bytes32 ) view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) AgentLatestVersion(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentLatestVersion", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentLatestVersion is a free data retrieval call binding the contract method 0x380619bb.
//
// Solidity: function agentLatestVersion(bytes32 ) view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) AgentLatestVersion(arg0 [32]byte) (*big.Int, error) {
	return _AgentRegistry.Contract.AgentLatestVersion(&_AgentRegistry.CallOpts, arg0)
}

// AgentLatestVersion is a free data retrieval call binding the contract method 0x380619bb.
//
// Solidity: function agentLatestVersion(bytes32 ) view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) AgentLatestVersion(arg0 [32]byte) (*big.Int, error) {
	return _AgentRegistry.Contract.AgentLatestVersion(&_AgentRegistry.CallOpts, arg0)
}

// AgentReference is a free data retrieval call binding the contract method 0x6a47f4f6.
//
// Solidity: function agentReference(bytes32 , uint256 ) view returns(string)
func (_AgentRegistry *AgentRegistryCaller) AgentReference(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (string, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentReference", arg0, arg1)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AgentReference is a free data retrieval call binding the contract method 0x6a47f4f6.
//
// Solidity: function agentReference(bytes32 , uint256 ) view returns(string)
func (_AgentRegistry *AgentRegistrySession) AgentReference(arg0 [32]byte, arg1 *big.Int) (string, error) {
	return _AgentRegistry.Contract.AgentReference(&_AgentRegistry.CallOpts, arg0, arg1)
}

// AgentReference is a free data retrieval call binding the contract method 0x6a47f4f6.
//
// Solidity: function agentReference(bytes32 , uint256 ) view returns(string)
func (_AgentRegistry *AgentRegistryCallerSession) AgentReference(arg0 [32]byte, arg1 *big.Int) (string, error) {
	return _AgentRegistry.Contract.AgentReference(&_AgentRegistry.CallOpts, arg0, arg1)
}

// GetLatestRef is a free data retrieval call binding the contract method 0xe15091a2.
//
// Solidity: function getLatestRef(bytes32 _agentId) view returns(string)
func (_AgentRegistry *AgentRegistryCaller) GetLatestRef(opts *bind.CallOpts, _agentId [32]byte) (string, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "getLatestRef", _agentId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetLatestRef is a free data retrieval call binding the contract method 0xe15091a2.
//
// Solidity: function getLatestRef(bytes32 _agentId) view returns(string)
func (_AgentRegistry *AgentRegistrySession) GetLatestRef(_agentId [32]byte) (string, error) {
	return _AgentRegistry.Contract.GetLatestRef(&_AgentRegistry.CallOpts, _agentId)
}

// GetLatestRef is a free data retrieval call binding the contract method 0xe15091a2.
//
// Solidity: function getLatestRef(bytes32 _agentId) view returns(string)
func (_AgentRegistry *AgentRegistryCallerSession) GetLatestRef(_agentId [32]byte) (string, error) {
	return _AgentRegistry.Contract.GetLatestRef(&_AgentRegistry.CallOpts, _agentId)
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

// AddAgentDev is a paid mutator transaction binding the contract method 0x6dfa2fef.
//
// Solidity: function addAgentDev(bytes32 _agentId, address _dev) returns()
func (_AgentRegistry *AgentRegistryTransactor) AddAgentDev(opts *bind.TransactOpts, _agentId [32]byte, _dev common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "addAgentDev", _agentId, _dev)
}

// AddAgentDev is a paid mutator transaction binding the contract method 0x6dfa2fef.
//
// Solidity: function addAgentDev(bytes32 _agentId, address _dev) returns()
func (_AgentRegistry *AgentRegistrySession) AddAgentDev(_agentId [32]byte, _dev common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddAgentDev(&_AgentRegistry.TransactOpts, _agentId, _dev)
}

// AddAgentDev is a paid mutator transaction binding the contract method 0x6dfa2fef.
//
// Solidity: function addAgentDev(bytes32 _agentId, address _dev) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) AddAgentDev(_agentId [32]byte, _dev common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddAgentDev(&_AgentRegistry.TransactOpts, _agentId, _dev)
}

// CreateAgent is a paid mutator transaction binding the contract method 0x2a083465.
//
// Solidity: function createAgent(bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactor) CreateAgent(opts *bind.TransactOpts, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "createAgent", _agentId, _ref)
}

// CreateAgent is a paid mutator transaction binding the contract method 0x2a083465.
//
// Solidity: function createAgent(bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistrySession) CreateAgent(_agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.CreateAgent(&_AgentRegistry.TransactOpts, _agentId, _ref)
}

// CreateAgent is a paid mutator transaction binding the contract method 0x2a083465.
//
// Solidity: function createAgent(bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) CreateAgent(_agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.CreateAgent(&_AgentRegistry.TransactOpts, _agentId, _ref)
}

// DisableAgent is a paid mutator transaction binding the contract method 0x09c53a27.
//
// Solidity: function disableAgent(bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistryTransactor) DisableAgent(opts *bind.TransactOpts, _agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "disableAgent", _agentId)
}

// DisableAgent is a paid mutator transaction binding the contract method 0x09c53a27.
//
// Solidity: function disableAgent(bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistrySession) DisableAgent(_agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.DisableAgent(&_AgentRegistry.TransactOpts, _agentId)
}

// DisableAgent is a paid mutator transaction binding the contract method 0x09c53a27.
//
// Solidity: function disableAgent(bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) DisableAgent(_agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.DisableAgent(&_AgentRegistry.TransactOpts, _agentId)
}

// EnableAgent is a paid mutator transaction binding the contract method 0x7879ae7b.
//
// Solidity: function enableAgent(bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistryTransactor) EnableAgent(opts *bind.TransactOpts, _agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "enableAgent", _agentId)
}

// EnableAgent is a paid mutator transaction binding the contract method 0x7879ae7b.
//
// Solidity: function enableAgent(bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistrySession) EnableAgent(_agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.EnableAgent(&_AgentRegistry.TransactOpts, _agentId)
}

// EnableAgent is a paid mutator transaction binding the contract method 0x7879ae7b.
//
// Solidity: function enableAgent(bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) EnableAgent(_agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.EnableAgent(&_AgentRegistry.TransactOpts, _agentId)
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

// RemoveAgentDev is a paid mutator transaction binding the contract method 0x0024cdf2.
//
// Solidity: function removeAgentDev(bytes32 _agentId, address _dev) returns()
func (_AgentRegistry *AgentRegistryTransactor) RemoveAgentDev(opts *bind.TransactOpts, _agentId [32]byte, _dev common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "removeAgentDev", _agentId, _dev)
}

// RemoveAgentDev is a paid mutator transaction binding the contract method 0x0024cdf2.
//
// Solidity: function removeAgentDev(bytes32 _agentId, address _dev) returns()
func (_AgentRegistry *AgentRegistrySession) RemoveAgentDev(_agentId [32]byte, _dev common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemoveAgentDev(&_AgentRegistry.TransactOpts, _agentId, _dev)
}

// RemoveAgentDev is a paid mutator transaction binding the contract method 0x0024cdf2.
//
// Solidity: function removeAgentDev(bytes32 _agentId, address _dev) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) RemoveAgentDev(_agentId [32]byte, _dev common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemoveAgentDev(&_AgentRegistry.TransactOpts, _agentId, _dev)
}

// TransferAgentOwnership is a paid mutator transaction binding the contract method 0x6a7095e3.
//
// Solidity: function transferAgentOwnership(bytes32 _agentId, address _to) returns()
func (_AgentRegistry *AgentRegistryTransactor) TransferAgentOwnership(opts *bind.TransactOpts, _agentId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "transferAgentOwnership", _agentId, _to)
}

// TransferAgentOwnership is a paid mutator transaction binding the contract method 0x6a7095e3.
//
// Solidity: function transferAgentOwnership(bytes32 _agentId, address _to) returns()
func (_AgentRegistry *AgentRegistrySession) TransferAgentOwnership(_agentId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.TransferAgentOwnership(&_AgentRegistry.TransactOpts, _agentId, _to)
}

// TransferAgentOwnership is a paid mutator transaction binding the contract method 0x6a7095e3.
//
// Solidity: function transferAgentOwnership(bytes32 _agentId, address _to) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) TransferAgentOwnership(_agentId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.TransferAgentOwnership(&_AgentRegistry.TransactOpts, _agentId, _to)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xf8e66c6e.
//
// Solidity: function updateAgent(bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactor) UpdateAgent(opts *bind.TransactOpts, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "updateAgent", _agentId, _ref)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xf8e66c6e.
//
// Solidity: function updateAgent(bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistrySession) UpdateAgent(_agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpdateAgent(&_AgentRegistry.TransactOpts, _agentId, _ref)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xf8e66c6e.
//
// Solidity: function updateAgent(bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) UpdateAgent(_agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpdateAgent(&_AgentRegistry.TransactOpts, _agentId, _ref)
}

// AgentRegistryAgentCreatedIterator is returned from FilterAgentCreated and is used to iterate over the raw logs and unpacked data for AgentCreated events raised by the AgentRegistry contract.
type AgentRegistryAgentCreatedIterator struct {
	Event *AgentRegistryAgentCreated // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentCreated)
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
		it.Event = new(AgentRegistryAgentCreated)
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
func (it *AgentRegistryAgentCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentCreated represents a AgentCreated event raised by the AgentRegistry contract.
type AgentRegistryAgentCreated struct {
	AgentId [32]byte
	Ref     string
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentCreated is a free log retrieval operation binding the contract event 0x3fb68868592f9cbe88b7ab6f760982ce46981f0a76b4609e2ef8ba1e36273af7.
//
// Solidity: event AgentCreated(bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentCreated(opts *bind.FilterOpts) (*AgentRegistryAgentCreatedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentCreated")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentCreatedIterator{contract: _AgentRegistry.contract, event: "AgentCreated", logs: logs, sub: sub}, nil
}

// WatchAgentCreated is a free log subscription operation binding the contract event 0x3fb68868592f9cbe88b7ab6f760982ce46981f0a76b4609e2ef8ba1e36273af7.
//
// Solidity: event AgentCreated(bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentCreated(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentCreated) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentCreated)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentCreated", log); err != nil {
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

// ParseAgentCreated is a log parse operation binding the contract event 0x3fb68868592f9cbe88b7ab6f760982ce46981f0a76b4609e2ef8ba1e36273af7.
//
// Solidity: event AgentCreated(bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentCreated(log types.Log) (*AgentRegistryAgentCreated, error) {
	event := new(AgentRegistryAgentCreated)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentDevAddedIterator is returned from FilterAgentDevAdded and is used to iterate over the raw logs and unpacked data for AgentDevAdded events raised by the AgentRegistry contract.
type AgentRegistryAgentDevAddedIterator struct {
	Event *AgentRegistryAgentDevAdded // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentDevAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentDevAdded)
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
		it.Event = new(AgentRegistryAgentDevAdded)
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
func (it *AgentRegistryAgentDevAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentDevAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentDevAdded represents a AgentDevAdded event raised by the AgentRegistry contract.
type AgentRegistryAgentDevAdded struct {
	AgentId [32]byte
	Dev     common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentDevAdded is a free log retrieval operation binding the contract event 0x562a4954d968dadacb6b407d782f2cad3d9dbc2406388a79f7c822d4f446741b.
//
// Solidity: event AgentDevAdded(bytes32 agentId, address dev)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentDevAdded(opts *bind.FilterOpts) (*AgentRegistryAgentDevAddedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentDevAdded")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentDevAddedIterator{contract: _AgentRegistry.contract, event: "AgentDevAdded", logs: logs, sub: sub}, nil
}

// WatchAgentDevAdded is a free log subscription operation binding the contract event 0x562a4954d968dadacb6b407d782f2cad3d9dbc2406388a79f7c822d4f446741b.
//
// Solidity: event AgentDevAdded(bytes32 agentId, address dev)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentDevAdded(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentDevAdded) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentDevAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentDevAdded)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentDevAdded", log); err != nil {
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

// ParseAgentDevAdded is a log parse operation binding the contract event 0x562a4954d968dadacb6b407d782f2cad3d9dbc2406388a79f7c822d4f446741b.
//
// Solidity: event AgentDevAdded(bytes32 agentId, address dev)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentDevAdded(log types.Log) (*AgentRegistryAgentDevAdded, error) {
	event := new(AgentRegistryAgentDevAdded)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentDevAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentDevRemovedIterator is returned from FilterAgentDevRemoved and is used to iterate over the raw logs and unpacked data for AgentDevRemoved events raised by the AgentRegistry contract.
type AgentRegistryAgentDevRemovedIterator struct {
	Event *AgentRegistryAgentDevRemoved // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentDevRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentDevRemoved)
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
		it.Event = new(AgentRegistryAgentDevRemoved)
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
func (it *AgentRegistryAgentDevRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentDevRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentDevRemoved represents a AgentDevRemoved event raised by the AgentRegistry contract.
type AgentRegistryAgentDevRemoved struct {
	AgentId [32]byte
	Dev     common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentDevRemoved is a free log retrieval operation binding the contract event 0x3abd19868328c96de107509ea68ef7fbe8b0d916643077d14c5f33dbf572a536.
//
// Solidity: event AgentDevRemoved(bytes32 agentId, address dev)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentDevRemoved(opts *bind.FilterOpts) (*AgentRegistryAgentDevRemovedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentDevRemoved")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentDevRemovedIterator{contract: _AgentRegistry.contract, event: "AgentDevRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentDevRemoved is a free log subscription operation binding the contract event 0x3abd19868328c96de107509ea68ef7fbe8b0d916643077d14c5f33dbf572a536.
//
// Solidity: event AgentDevRemoved(bytes32 agentId, address dev)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentDevRemoved(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentDevRemoved) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentDevRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentDevRemoved)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentDevRemoved", log); err != nil {
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

// ParseAgentDevRemoved is a log parse operation binding the contract event 0x3abd19868328c96de107509ea68ef7fbe8b0d916643077d14c5f33dbf572a536.
//
// Solidity: event AgentDevRemoved(bytes32 agentId, address dev)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentDevRemoved(log types.Log) (*AgentRegistryAgentDevRemoved, error) {
	event := new(AgentRegistryAgentDevRemoved)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentDevRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentDisabledIterator is returned from FilterAgentDisabled and is used to iterate over the raw logs and unpacked data for AgentDisabled events raised by the AgentRegistry contract.
type AgentRegistryAgentDisabledIterator struct {
	Event *AgentRegistryAgentDisabled // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentDisabled)
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
		it.Event = new(AgentRegistryAgentDisabled)
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
func (it *AgentRegistryAgentDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentDisabled represents a AgentDisabled event raised by the AgentRegistry contract.
type AgentRegistryAgentDisabled struct {
	AgentId [32]byte
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentDisabled is a free log retrieval operation binding the contract event 0x0a570f21f16393b7b462abfe9536bee30f7911547097671f80fb7d2f2908344e.
//
// Solidity: event AgentDisabled(bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentDisabled(opts *bind.FilterOpts) (*AgentRegistryAgentDisabledIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentDisabled")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentDisabledIterator{contract: _AgentRegistry.contract, event: "AgentDisabled", logs: logs, sub: sub}, nil
}

// WatchAgentDisabled is a free log subscription operation binding the contract event 0x0a570f21f16393b7b462abfe9536bee30f7911547097671f80fb7d2f2908344e.
//
// Solidity: event AgentDisabled(bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentDisabled(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentDisabled) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentDisabled)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentDisabled", log); err != nil {
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

// ParseAgentDisabled is a log parse operation binding the contract event 0x0a570f21f16393b7b462abfe9536bee30f7911547097671f80fb7d2f2908344e.
//
// Solidity: event AgentDisabled(bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentDisabled(log types.Log) (*AgentRegistryAgentDisabled, error) {
	event := new(AgentRegistryAgentDisabled)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentEnabledIterator is returned from FilterAgentEnabled and is used to iterate over the raw logs and unpacked data for AgentEnabled events raised by the AgentRegistry contract.
type AgentRegistryAgentEnabledIterator struct {
	Event *AgentRegistryAgentEnabled // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentEnabled)
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
		it.Event = new(AgentRegistryAgentEnabled)
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
func (it *AgentRegistryAgentEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentEnabled represents a AgentEnabled event raised by the AgentRegistry contract.
type AgentRegistryAgentEnabled struct {
	AgentId [32]byte
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentEnabled is a free log retrieval operation binding the contract event 0xb02042f4079bf702a976aa95de4238d66f8dcfe9127c71987600a8be91245a72.
//
// Solidity: event AgentEnabled(bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentEnabled(opts *bind.FilterOpts) (*AgentRegistryAgentEnabledIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentEnabled")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentEnabledIterator{contract: _AgentRegistry.contract, event: "AgentEnabled", logs: logs, sub: sub}, nil
}

// WatchAgentEnabled is a free log subscription operation binding the contract event 0xb02042f4079bf702a976aa95de4238d66f8dcfe9127c71987600a8be91245a72.
//
// Solidity: event AgentEnabled(bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentEnabled(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentEnabled) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentEnabled)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentEnabled", log); err != nil {
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

// ParseAgentEnabled is a log parse operation binding the contract event 0xb02042f4079bf702a976aa95de4238d66f8dcfe9127c71987600a8be91245a72.
//
// Solidity: event AgentEnabled(bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentEnabled(log types.Log) (*AgentRegistryAgentEnabled, error) {
	event := new(AgentRegistryAgentEnabled)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentOwnershipTransferredIterator is returned from FilterAgentOwnershipTransferred and is used to iterate over the raw logs and unpacked data for AgentOwnershipTransferred events raised by the AgentRegistry contract.
type AgentRegistryAgentOwnershipTransferredIterator struct {
	Event *AgentRegistryAgentOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentOwnershipTransferred)
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
		it.Event = new(AgentRegistryAgentOwnershipTransferred)
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
func (it *AgentRegistryAgentOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentOwnershipTransferred represents a AgentOwnershipTransferred event raised by the AgentRegistry contract.
type AgentRegistryAgentOwnershipTransferred struct {
	AgentId [32]byte
	From    common.Address
	To      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentOwnershipTransferred is a free log retrieval operation binding the contract event 0xc5f5c5cbe4eacc5301c81a3faa85abee1fc13485967f1dd5f8b8b170d511db9b.
//
// Solidity: event AgentOwnershipTransferred(bytes32 agentId, address from, address to)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentOwnershipTransferred(opts *bind.FilterOpts) (*AgentRegistryAgentOwnershipTransferredIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentOwnershipTransferred")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentOwnershipTransferredIterator{contract: _AgentRegistry.contract, event: "AgentOwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchAgentOwnershipTransferred is a free log subscription operation binding the contract event 0xc5f5c5cbe4eacc5301c81a3faa85abee1fc13485967f1dd5f8b8b170d511db9b.
//
// Solidity: event AgentOwnershipTransferred(bytes32 agentId, address from, address to)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentOwnershipTransferred) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentOwnershipTransferred")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentOwnershipTransferred)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentOwnershipTransferred", log); err != nil {
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

// ParseAgentOwnershipTransferred is a log parse operation binding the contract event 0xc5f5c5cbe4eacc5301c81a3faa85abee1fc13485967f1dd5f8b8b170d511db9b.
//
// Solidity: event AgentOwnershipTransferred(bytes32 agentId, address from, address to)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentOwnershipTransferred(log types.Log) (*AgentRegistryAgentOwnershipTransferred, error) {
	event := new(AgentRegistryAgentOwnershipTransferred)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentOwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentVersionPublishedIterator is returned from FilterAgentVersionPublished and is used to iterate over the raw logs and unpacked data for AgentVersionPublished events raised by the AgentRegistry contract.
type AgentRegistryAgentVersionPublishedIterator struct {
	Event *AgentRegistryAgentVersionPublished // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentVersionPublishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentVersionPublished)
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
		it.Event = new(AgentRegistryAgentVersionPublished)
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
func (it *AgentRegistryAgentVersionPublishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentVersionPublishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentVersionPublished represents a AgentVersionPublished event raised by the AgentRegistry contract.
type AgentRegistryAgentVersionPublished struct {
	AgentId [32]byte
	Version *big.Int
	Ref     string
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentVersionPublished is a free log retrieval operation binding the contract event 0x5703130761babe8d891d09dc2e0ba0e0f545c2798f55daa85d215947ef3a60b5.
//
// Solidity: event AgentVersionPublished(bytes32 agentId, uint256 version, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentVersionPublished(opts *bind.FilterOpts) (*AgentRegistryAgentVersionPublishedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentVersionPublished")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentVersionPublishedIterator{contract: _AgentRegistry.contract, event: "AgentVersionPublished", logs: logs, sub: sub}, nil
}

// WatchAgentVersionPublished is a free log subscription operation binding the contract event 0x5703130761babe8d891d09dc2e0ba0e0f545c2798f55daa85d215947ef3a60b5.
//
// Solidity: event AgentVersionPublished(bytes32 agentId, uint256 version, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentVersionPublished(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentVersionPublished) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentVersionPublished")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentVersionPublished)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentVersionPublished", log); err != nil {
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

// ParseAgentVersionPublished is a log parse operation binding the contract event 0x5703130761babe8d891d09dc2e0ba0e0f545c2798f55daa85d215947ef3a60b5.
//
// Solidity: event AgentVersionPublished(bytes32 agentId, uint256 version, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentVersionPublished(log types.Log) (*AgentRegistryAgentVersionPublished, error) {
	event := new(AgentRegistryAgentVersionPublished)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentVersionPublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
