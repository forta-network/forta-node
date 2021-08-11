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

// AlertsABI is the input ABI used to generate the binding from.
const AlertsABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentRegistryChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"alertId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"alertRef\",\"type\":\"string\"}],\"name\":\"Alert\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"PoolRegistryChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"type\":\"function\",\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_alertId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_alertRef\",\"type\":\"string\"}],\"name\":\"addAlert\",\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractMinimalForwarderUpgradeable\",\"name\":\"_forwarder\",\"type\":\"address\"},{\"internalType\":\"contractPoolRegistry\",\"name\":\"_poolRegistry\",\"type\":\"address\"},{\"internalType\":\"contractAgentRegistry\",\"name\":\"_agentRegistry\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractPoolRegistry\",\"name\":\"_poolRegistry\",\"type\":\"address\"},{\"internalType\":\"contractAgentRegistry\",\"name\":\"_agentRegistry\",\"type\":\"address\"}],\"name\":\"setRegistries\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Alerts is an auto generated Go binding around an Ethereum contract.
type Alerts struct {
	AlertsCaller     // Read-only binding to the contract
	AlertsTransactor // Write-only binding to the contract
	AlertsFilterer   // Log filterer for contract events
}

// AlertsCaller is an auto generated read-only Go binding around an Ethereum contract.
type AlertsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AlertsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AlertsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AlertsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AlertsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AlertsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AlertsSession struct {
	Contract     *Alerts           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AlertsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AlertsCallerSession struct {
	Contract *AlertsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AlertsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AlertsTransactorSession struct {
	Contract     *AlertsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AlertsRaw is an auto generated low-level Go binding around an Ethereum contract.
type AlertsRaw struct {
	Contract *Alerts // Generic contract binding to access the raw methods on
}

// AlertsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AlertsCallerRaw struct {
	Contract *AlertsCaller // Generic read-only contract binding to access the raw methods on
}

// AlertsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AlertsTransactorRaw struct {
	Contract *AlertsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAlerts creates a new instance of Alerts, bound to a specific deployed contract.
func NewAlerts(address common.Address, backend bind.ContractBackend) (*Alerts, error) {
	contract, err := bindAlerts(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Alerts{AlertsCaller: AlertsCaller{contract: contract}, AlertsTransactor: AlertsTransactor{contract: contract}, AlertsFilterer: AlertsFilterer{contract: contract}}, nil
}

// NewAlertsCaller creates a new read-only instance of Alerts, bound to a specific deployed contract.
func NewAlertsCaller(address common.Address, caller bind.ContractCaller) (*AlertsCaller, error) {
	contract, err := bindAlerts(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AlertsCaller{contract: contract}, nil
}

// NewAlertsTransactor creates a new write-only instance of Alerts, bound to a specific deployed contract.
func NewAlertsTransactor(address common.Address, transactor bind.ContractTransactor) (*AlertsTransactor, error) {
	contract, err := bindAlerts(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AlertsTransactor{contract: contract}, nil
}

// NewAlertsFilterer creates a new log filterer instance of Alerts, bound to a specific deployed contract.
func NewAlertsFilterer(address common.Address, filterer bind.ContractFilterer) (*AlertsFilterer, error) {
	contract, err := bindAlerts(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AlertsFilterer{contract: contract}, nil
}

// bindAlerts binds a generic wrapper to an already deployed contract.
func bindAlerts(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AlertsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Alerts *AlertsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Alerts.Contract.AlertsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Alerts *AlertsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Alerts.Contract.AlertsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Alerts *AlertsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Alerts.Contract.AlertsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Alerts *AlertsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Alerts.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Alerts *AlertsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Alerts.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Alerts *AlertsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Alerts.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Alerts *AlertsCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Alerts.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Alerts *AlertsSession) ADMINROLE() ([32]byte, error) {
	return _Alerts.Contract.ADMINROLE(&_Alerts.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Alerts *AlertsCallerSession) ADMINROLE() ([32]byte, error) {
	return _Alerts.Contract.ADMINROLE(&_Alerts.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Alerts *AlertsCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Alerts.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Alerts *AlertsSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Alerts.Contract.DEFAULTADMINROLE(&_Alerts.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Alerts *AlertsCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Alerts.Contract.DEFAULTADMINROLE(&_Alerts.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Alerts *AlertsCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Alerts.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Alerts *AlertsSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Alerts.Contract.GetRoleAdmin(&_Alerts.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Alerts *AlertsCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Alerts.Contract.GetRoleAdmin(&_Alerts.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Alerts *AlertsCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Alerts.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Alerts *AlertsSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Alerts.Contract.HasRole(&_Alerts.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Alerts *AlertsCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Alerts.Contract.HasRole(&_Alerts.CallOpts, role, account)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_Alerts *AlertsCaller) IsTrustedForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _Alerts.contract.Call(opts, &out, "isTrustedForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_Alerts *AlertsSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _Alerts.Contract.IsTrustedForwarder(&_Alerts.CallOpts, forwarder)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_Alerts *AlertsCallerSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _Alerts.Contract.IsTrustedForwarder(&_Alerts.CallOpts, forwarder)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Alerts *AlertsCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Alerts.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Alerts *AlertsSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Alerts.Contract.SupportsInterface(&_Alerts.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Alerts *AlertsCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Alerts.Contract.SupportsInterface(&_Alerts.CallOpts, interfaceId)
}

// AddAlert is a paid mutator transaction binding the contract method 0x1a8e12b8.
//
// Solidity: function addAlert(bytes32 _poolId, bytes32 _agentId, bytes32 _alertId, string _alertRef) returns()
func (_Alerts *AlertsTransactor) AddAlert(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, _alertId [32]byte, _alertRef string) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "addAlert", _poolId, _agentId, _alertId, _alertRef)
}

// AddAlert is a paid mutator transaction binding the contract method 0x1a8e12b8.
//
// Solidity: function addAlert(bytes32 _poolId, bytes32 _agentId, bytes32 _alertId, string _alertRef) returns()
func (_Alerts *AlertsSession) AddAlert(_poolId [32]byte, _agentId [32]byte, _alertId [32]byte, _alertRef string) (*types.Transaction, error) {
	return _Alerts.Contract.AddAlert(&_Alerts.TransactOpts, _poolId, _agentId, _alertId, _alertRef)
}

// AddAlert is a paid mutator transaction binding the contract method 0x1a8e12b8.
//
// Solidity: function addAlert(bytes32 _poolId, bytes32 _agentId, bytes32 _alertId, string _alertRef) returns()
func (_Alerts *AlertsTransactorSession) AddAlert(_poolId [32]byte, _agentId [32]byte, _alertId [32]byte, _alertRef string) (*types.Transaction, error) {
	return _Alerts.Contract.AddAlert(&_Alerts.TransactOpts, _poolId, _agentId, _alertId, _alertRef)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Alerts *AlertsTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Alerts *AlertsSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.GrantRole(&_Alerts.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Alerts *AlertsTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.GrantRole(&_Alerts.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _forwarder, address _poolRegistry, address _agentRegistry) returns()
func (_Alerts *AlertsTransactor) Initialize(opts *bind.TransactOpts, _forwarder common.Address, _poolRegistry common.Address, _agentRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "initialize", _forwarder, _poolRegistry, _agentRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _forwarder, address _poolRegistry, address _agentRegistry) returns()
func (_Alerts *AlertsSession) Initialize(_forwarder common.Address, _poolRegistry common.Address, _agentRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.Initialize(&_Alerts.TransactOpts, _forwarder, _poolRegistry, _agentRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _forwarder, address _poolRegistry, address _agentRegistry) returns()
func (_Alerts *AlertsTransactorSession) Initialize(_forwarder common.Address, _poolRegistry common.Address, _agentRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.Initialize(&_Alerts.TransactOpts, _forwarder, _poolRegistry, _agentRegistry)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Alerts *AlertsTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Alerts *AlertsSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.RenounceRole(&_Alerts.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Alerts *AlertsTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.RenounceRole(&_Alerts.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Alerts *AlertsTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Alerts *AlertsSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.RevokeRole(&_Alerts.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Alerts *AlertsTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.RevokeRole(&_Alerts.TransactOpts, role, account)
}

// SetRegistries is a paid mutator transaction binding the contract method 0x399bdf2a.
//
// Solidity: function setRegistries(address _poolRegistry, address _agentRegistry) returns()
func (_Alerts *AlertsTransactor) SetRegistries(opts *bind.TransactOpts, _poolRegistry common.Address, _agentRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "setRegistries", _poolRegistry, _agentRegistry)
}

// SetRegistries is a paid mutator transaction binding the contract method 0x399bdf2a.
//
// Solidity: function setRegistries(address _poolRegistry, address _agentRegistry) returns()
func (_Alerts *AlertsSession) SetRegistries(_poolRegistry common.Address, _agentRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.SetRegistries(&_Alerts.TransactOpts, _poolRegistry, _agentRegistry)
}

// SetRegistries is a paid mutator transaction binding the contract method 0x399bdf2a.
//
// Solidity: function setRegistries(address _poolRegistry, address _agentRegistry) returns()
func (_Alerts *AlertsTransactorSession) SetRegistries(_poolRegistry common.Address, _agentRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.SetRegistries(&_Alerts.TransactOpts, _poolRegistry, _agentRegistry)
}

// AlertsAgentRegistryChangedIterator is returned from FilterAgentRegistryChanged and is used to iterate over the raw logs and unpacked data for AgentRegistryChanged events raised by the Alerts contract.
type AlertsAgentRegistryChangedIterator struct {
	Event *AlertsAgentRegistryChanged // Event containing the contract specifics and raw log

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
func (it *AlertsAgentRegistryChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsAgentRegistryChanged)
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
		it.Event = new(AlertsAgentRegistryChanged)
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
func (it *AlertsAgentRegistryChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsAgentRegistryChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsAgentRegistryChanged represents a AgentRegistryChanged event raised by the Alerts contract.
type AlertsAgentRegistryChanged struct {
	From common.Address
	To   common.Address
	By   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAgentRegistryChanged is a free log retrieval operation binding the contract event 0x9258ae620f73184170c242afd2d6549360b80d3615dcdb0caeb30f369bc45d28.
//
// Solidity: event AgentRegistryChanged(address from, address to, address by)
func (_Alerts *AlertsFilterer) FilterAgentRegistryChanged(opts *bind.FilterOpts) (*AlertsAgentRegistryChangedIterator, error) {

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "AgentRegistryChanged")
	if err != nil {
		return nil, err
	}
	return &AlertsAgentRegistryChangedIterator{contract: _Alerts.contract, event: "AgentRegistryChanged", logs: logs, sub: sub}, nil
}

// WatchAgentRegistryChanged is a free log subscription operation binding the contract event 0x9258ae620f73184170c242afd2d6549360b80d3615dcdb0caeb30f369bc45d28.
//
// Solidity: event AgentRegistryChanged(address from, address to, address by)
func (_Alerts *AlertsFilterer) WatchAgentRegistryChanged(opts *bind.WatchOpts, sink chan<- *AlertsAgentRegistryChanged) (event.Subscription, error) {

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "AgentRegistryChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsAgentRegistryChanged)
				if err := _Alerts.contract.UnpackLog(event, "AgentRegistryChanged", log); err != nil {
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

// ParseAgentRegistryChanged is a log parse operation binding the contract event 0x9258ae620f73184170c242afd2d6549360b80d3615dcdb0caeb30f369bc45d28.
//
// Solidity: event AgentRegistryChanged(address from, address to, address by)
func (_Alerts *AlertsFilterer) ParseAgentRegistryChanged(log types.Log) (*AlertsAgentRegistryChanged, error) {
	event := new(AlertsAgentRegistryChanged)
	if err := _Alerts.contract.UnpackLog(event, "AgentRegistryChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsAlertIterator is returned from FilterAlert and is used to iterate over the raw logs and unpacked data for Alert events raised by the Alerts contract.
type AlertsAlertIterator struct {
	Event *AlertsAlert // Event containing the contract specifics and raw log

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
func (it *AlertsAlertIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsAlert)
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
		it.Event = new(AlertsAlert)
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
func (it *AlertsAlertIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsAlertIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsAlert represents a Alert event raised by the Alerts contract.
type AlertsAlert struct {
	PoolId   [32]byte
	Scanner  common.Address
	AgentId  [32]byte
	AlertId  [32]byte
	AlertRef string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAlert is a free log retrieval operation binding the contract event 0xd097506e69e86a46632e8e4d68ffc302eaeb63dfe08ffc44ed924f03fcd1a039.
//
// Solidity: event Alert(bytes32 indexed poolId, address indexed scanner, bytes32 indexed agentId, bytes32 alertId, string alertRef)
func (_Alerts *AlertsFilterer) FilterAlert(opts *bind.FilterOpts, poolId [][32]byte, scanner []common.Address, agentId [][32]byte) (*AlertsAlertIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var scannerRule []interface{}
	for _, scannerItem := range scanner {
		scannerRule = append(scannerRule, scannerItem)
	}
	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "Alert", poolIdRule, scannerRule, agentIdRule)
	if err != nil {
		return nil, err
	}
	return &AlertsAlertIterator{contract: _Alerts.contract, event: "Alert", logs: logs, sub: sub}, nil
}

// WatchAlert is a free log subscription operation binding the contract event 0xd097506e69e86a46632e8e4d68ffc302eaeb63dfe08ffc44ed924f03fcd1a039.
//
// Solidity: event Alert(bytes32 indexed poolId, address indexed scanner, bytes32 indexed agentId, bytes32 alertId, string alertRef)
func (_Alerts *AlertsFilterer) WatchAlert(opts *bind.WatchOpts, sink chan<- *AlertsAlert, poolId [][32]byte, scanner []common.Address, agentId [][32]byte) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var scannerRule []interface{}
	for _, scannerItem := range scanner {
		scannerRule = append(scannerRule, scannerItem)
	}
	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "Alert", poolIdRule, scannerRule, agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsAlert)
				if err := _Alerts.contract.UnpackLog(event, "Alert", log); err != nil {
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

// ParseAlert is a log parse operation binding the contract event 0xd097506e69e86a46632e8e4d68ffc302eaeb63dfe08ffc44ed924f03fcd1a039.
//
// Solidity: event Alert(bytes32 indexed poolId, address indexed scanner, bytes32 indexed agentId, bytes32 alertId, string alertRef)
func (_Alerts *AlertsFilterer) ParseAlert(log types.Log) (*AlertsAlert, error) {
	event := new(AlertsAlert)
	if err := _Alerts.contract.UnpackLog(event, "Alert", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsPoolRegistryChangedIterator is returned from FilterPoolRegistryChanged and is used to iterate over the raw logs and unpacked data for PoolRegistryChanged events raised by the Alerts contract.
type AlertsPoolRegistryChangedIterator struct {
	Event *AlertsPoolRegistryChanged // Event containing the contract specifics and raw log

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
func (it *AlertsPoolRegistryChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsPoolRegistryChanged)
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
		it.Event = new(AlertsPoolRegistryChanged)
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
func (it *AlertsPoolRegistryChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsPoolRegistryChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsPoolRegistryChanged represents a PoolRegistryChanged event raised by the Alerts contract.
type AlertsPoolRegistryChanged struct {
	From common.Address
	To   common.Address
	By   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterPoolRegistryChanged is a free log retrieval operation binding the contract event 0xf4f63f049796a39f99b6d5c64f26eae479cf9a51dfce0b3db21f69bf5d6c5829.
//
// Solidity: event PoolRegistryChanged(address from, address to, address by)
func (_Alerts *AlertsFilterer) FilterPoolRegistryChanged(opts *bind.FilterOpts) (*AlertsPoolRegistryChangedIterator, error) {

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "PoolRegistryChanged")
	if err != nil {
		return nil, err
	}
	return &AlertsPoolRegistryChangedIterator{contract: _Alerts.contract, event: "PoolRegistryChanged", logs: logs, sub: sub}, nil
}

// WatchPoolRegistryChanged is a free log subscription operation binding the contract event 0xf4f63f049796a39f99b6d5c64f26eae479cf9a51dfce0b3db21f69bf5d6c5829.
//
// Solidity: event PoolRegistryChanged(address from, address to, address by)
func (_Alerts *AlertsFilterer) WatchPoolRegistryChanged(opts *bind.WatchOpts, sink chan<- *AlertsPoolRegistryChanged) (event.Subscription, error) {

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "PoolRegistryChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsPoolRegistryChanged)
				if err := _Alerts.contract.UnpackLog(event, "PoolRegistryChanged", log); err != nil {
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

// ParsePoolRegistryChanged is a log parse operation binding the contract event 0xf4f63f049796a39f99b6d5c64f26eae479cf9a51dfce0b3db21f69bf5d6c5829.
//
// Solidity: event PoolRegistryChanged(address from, address to, address by)
func (_Alerts *AlertsFilterer) ParsePoolRegistryChanged(log types.Log) (*AlertsPoolRegistryChanged, error) {
	event := new(AlertsPoolRegistryChanged)
	if err := _Alerts.contract.UnpackLog(event, "PoolRegistryChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Alerts contract.
type AlertsRoleAdminChangedIterator struct {
	Event *AlertsRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *AlertsRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsRoleAdminChanged)
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
		it.Event = new(AlertsRoleAdminChanged)
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
func (it *AlertsRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsRoleAdminChanged represents a RoleAdminChanged event raised by the Alerts contract.
type AlertsRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Alerts *AlertsFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*AlertsRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &AlertsRoleAdminChangedIterator{contract: _Alerts.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Alerts *AlertsFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *AlertsRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsRoleAdminChanged)
				if err := _Alerts.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Alerts *AlertsFilterer) ParseRoleAdminChanged(log types.Log) (*AlertsRoleAdminChanged, error) {
	event := new(AlertsRoleAdminChanged)
	if err := _Alerts.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Alerts contract.
type AlertsRoleGrantedIterator struct {
	Event *AlertsRoleGranted // Event containing the contract specifics and raw log

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
func (it *AlertsRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsRoleGranted)
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
		it.Event = new(AlertsRoleGranted)
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
func (it *AlertsRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsRoleGranted represents a RoleGranted event raised by the Alerts contract.
type AlertsRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Alerts *AlertsFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AlertsRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AlertsRoleGrantedIterator{contract: _Alerts.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Alerts *AlertsFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *AlertsRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsRoleGranted)
				if err := _Alerts.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Alerts *AlertsFilterer) ParseRoleGranted(log types.Log) (*AlertsRoleGranted, error) {
	event := new(AlertsRoleGranted)
	if err := _Alerts.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Alerts contract.
type AlertsRoleRevokedIterator struct {
	Event *AlertsRoleRevoked // Event containing the contract specifics and raw log

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
func (it *AlertsRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsRoleRevoked)
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
		it.Event = new(AlertsRoleRevoked)
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
func (it *AlertsRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsRoleRevoked represents a RoleRevoked event raised by the Alerts contract.
type AlertsRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Alerts *AlertsFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AlertsRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AlertsRoleRevokedIterator{contract: _Alerts.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Alerts *AlertsFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *AlertsRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsRoleRevoked)
				if err := _Alerts.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Alerts *AlertsFilterer) ParseRoleRevoked(log types.Log) (*AlertsRoleRevoked, error) {
	event := new(AlertsRoleRevoked)
	if err := _Alerts.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
