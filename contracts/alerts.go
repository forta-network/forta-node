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

// AlertsMetaData contains all meta data concerning the Alerts contract.
var AlertsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddressManager\",\"type\":\"address\"}],\"name\":\"AccessManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"alertsId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockStart\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockEnd\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"alertCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxSeverity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"}],\"name\":\"AlertBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"ScannerRegistryChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_blockStart\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_blockEnd\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_alertCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxSeverity\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_ref\",\"type\":\"string\"}],\"name\":\"addAlertBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"__manager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"__router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"__scannerRegistry\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"scannerRegistry\",\"outputs\":[{\"internalType\":\"contractScannerRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newManager\",\"type\":\"address\"}],\"name\":\"setAccessManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ensRegistry\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"ensName\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newScannerRegistry\",\"type\":\"address\"}],\"name\":\"setScannerRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// AlertsABI is the input ABI used to generate the binding from.
// Deprecated: Use AlertsMetaData.ABI instead.
var AlertsABI = AlertsMetaData.ABI

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

// ScannerRegistry is a free data retrieval call binding the contract method 0x5e9f88b1.
//
// Solidity: function scannerRegistry() view returns(address)
func (_Alerts *AlertsCaller) ScannerRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Alerts.contract.Call(opts, &out, "scannerRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ScannerRegistry is a free data retrieval call binding the contract method 0x5e9f88b1.
//
// Solidity: function scannerRegistry() view returns(address)
func (_Alerts *AlertsSession) ScannerRegistry() (common.Address, error) {
	return _Alerts.Contract.ScannerRegistry(&_Alerts.CallOpts)
}

// ScannerRegistry is a free data retrieval call binding the contract method 0x5e9f88b1.
//
// Solidity: function scannerRegistry() view returns(address)
func (_Alerts *AlertsCallerSession) ScannerRegistry() (common.Address, error) {
	return _Alerts.Contract.ScannerRegistry(&_Alerts.CallOpts)
}

// AddAlertBatch is a paid mutator transaction binding the contract method 0x8defd825.
//
// Solidity: function addAlertBatch(uint256 _chainId, uint256 _blockStart, uint256 _blockEnd, uint256 _alertCount, uint256 _maxSeverity, string _ref) returns()
func (_Alerts *AlertsTransactor) AddAlertBatch(opts *bind.TransactOpts, _chainId *big.Int, _blockStart *big.Int, _blockEnd *big.Int, _alertCount *big.Int, _maxSeverity *big.Int, _ref string) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "addAlertBatch", _chainId, _blockStart, _blockEnd, _alertCount, _maxSeverity, _ref)
}

// AddAlertBatch is a paid mutator transaction binding the contract method 0x8defd825.
//
// Solidity: function addAlertBatch(uint256 _chainId, uint256 _blockStart, uint256 _blockEnd, uint256 _alertCount, uint256 _maxSeverity, string _ref) returns()
func (_Alerts *AlertsSession) AddAlertBatch(_chainId *big.Int, _blockStart *big.Int, _blockEnd *big.Int, _alertCount *big.Int, _maxSeverity *big.Int, _ref string) (*types.Transaction, error) {
	return _Alerts.Contract.AddAlertBatch(&_Alerts.TransactOpts, _chainId, _blockStart, _blockEnd, _alertCount, _maxSeverity, _ref)
}

// AddAlertBatch is a paid mutator transaction binding the contract method 0x8defd825.
//
// Solidity: function addAlertBatch(uint256 _chainId, uint256 _blockStart, uint256 _blockEnd, uint256 _alertCount, uint256 _maxSeverity, string _ref) returns()
func (_Alerts *AlertsTransactorSession) AddAlertBatch(_chainId *big.Int, _blockStart *big.Int, _blockEnd *big.Int, _alertCount *big.Int, _maxSeverity *big.Int, _ref string) (*types.Transaction, error) {
	return _Alerts.Contract.AddAlertBatch(&_Alerts.TransactOpts, _chainId, _blockStart, _blockEnd, _alertCount, _maxSeverity, _ref)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address __manager, address __router, address __scannerRegistry) returns()
func (_Alerts *AlertsTransactor) Initialize(opts *bind.TransactOpts, __manager common.Address, __router common.Address, __scannerRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "initialize", __manager, __router, __scannerRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address __manager, address __router, address __scannerRegistry) returns()
func (_Alerts *AlertsSession) Initialize(__manager common.Address, __router common.Address, __scannerRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.Initialize(&_Alerts.TransactOpts, __manager, __router, __scannerRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address __manager, address __router, address __scannerRegistry) returns()
func (_Alerts *AlertsTransactorSession) Initialize(__manager common.Address, __router common.Address, __scannerRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.Initialize(&_Alerts.TransactOpts, __manager, __router, __scannerRegistry)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Alerts *AlertsTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Alerts *AlertsSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Alerts.Contract.Multicall(&_Alerts.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Alerts *AlertsTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Alerts.Contract.Multicall(&_Alerts.TransactOpts, data)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_Alerts *AlertsTransactor) SetAccessManager(opts *bind.TransactOpts, newManager common.Address) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "setAccessManager", newManager)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_Alerts *AlertsSession) SetAccessManager(newManager common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.SetAccessManager(&_Alerts.TransactOpts, newManager)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_Alerts *AlertsTransactorSession) SetAccessManager(newManager common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.SetAccessManager(&_Alerts.TransactOpts, newManager)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_Alerts *AlertsTransactor) SetName(opts *bind.TransactOpts, ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "setName", ensRegistry, ensName)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_Alerts *AlertsSession) SetName(ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _Alerts.Contract.SetName(&_Alerts.TransactOpts, ensRegistry, ensName)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_Alerts *AlertsTransactorSession) SetName(ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _Alerts.Contract.SetName(&_Alerts.TransactOpts, ensRegistry, ensName)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_Alerts *AlertsTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "setRouter", newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_Alerts *AlertsSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.SetRouter(&_Alerts.TransactOpts, newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_Alerts *AlertsTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.SetRouter(&_Alerts.TransactOpts, newRouter)
}

// SetScannerRegistry is a paid mutator transaction binding the contract method 0x6b254492.
//
// Solidity: function setScannerRegistry(address newScannerRegistry) returns()
func (_Alerts *AlertsTransactor) SetScannerRegistry(opts *bind.TransactOpts, newScannerRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "setScannerRegistry", newScannerRegistry)
}

// SetScannerRegistry is a paid mutator transaction binding the contract method 0x6b254492.
//
// Solidity: function setScannerRegistry(address newScannerRegistry) returns()
func (_Alerts *AlertsSession) SetScannerRegistry(newScannerRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.SetScannerRegistry(&_Alerts.TransactOpts, newScannerRegistry)
}

// SetScannerRegistry is a paid mutator transaction binding the contract method 0x6b254492.
//
// Solidity: function setScannerRegistry(address newScannerRegistry) returns()
func (_Alerts *AlertsTransactorSession) SetScannerRegistry(newScannerRegistry common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.SetScannerRegistry(&_Alerts.TransactOpts, newScannerRegistry)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Alerts *AlertsTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Alerts *AlertsSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.UpgradeTo(&_Alerts.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Alerts *AlertsTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Alerts.Contract.UpgradeTo(&_Alerts.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Alerts *AlertsTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Alerts.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Alerts *AlertsSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Alerts.Contract.UpgradeToAndCall(&_Alerts.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Alerts *AlertsTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Alerts.Contract.UpgradeToAndCall(&_Alerts.TransactOpts, newImplementation, data)
}

// AlertsAccessManagerUpdatedIterator is returned from FilterAccessManagerUpdated and is used to iterate over the raw logs and unpacked data for AccessManagerUpdated events raised by the Alerts contract.
type AlertsAccessManagerUpdatedIterator struct {
	Event *AlertsAccessManagerUpdated // Event containing the contract specifics and raw log

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
func (it *AlertsAccessManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsAccessManagerUpdated)
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
		it.Event = new(AlertsAccessManagerUpdated)
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
func (it *AlertsAccessManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsAccessManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsAccessManagerUpdated represents a AccessManagerUpdated event raised by the Alerts contract.
type AlertsAccessManagerUpdated struct {
	NewAddressManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterAccessManagerUpdated is a free log retrieval operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_Alerts *AlertsFilterer) FilterAccessManagerUpdated(opts *bind.FilterOpts, newAddressManager []common.Address) (*AlertsAccessManagerUpdatedIterator, error) {

	var newAddressManagerRule []interface{}
	for _, newAddressManagerItem := range newAddressManager {
		newAddressManagerRule = append(newAddressManagerRule, newAddressManagerItem)
	}

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "AccessManagerUpdated", newAddressManagerRule)
	if err != nil {
		return nil, err
	}
	return &AlertsAccessManagerUpdatedIterator{contract: _Alerts.contract, event: "AccessManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchAccessManagerUpdated is a free log subscription operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_Alerts *AlertsFilterer) WatchAccessManagerUpdated(opts *bind.WatchOpts, sink chan<- *AlertsAccessManagerUpdated, newAddressManager []common.Address) (event.Subscription, error) {

	var newAddressManagerRule []interface{}
	for _, newAddressManagerItem := range newAddressManager {
		newAddressManagerRule = append(newAddressManagerRule, newAddressManagerItem)
	}

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "AccessManagerUpdated", newAddressManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsAccessManagerUpdated)
				if err := _Alerts.contract.UnpackLog(event, "AccessManagerUpdated", log); err != nil {
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

// ParseAccessManagerUpdated is a log parse operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_Alerts *AlertsFilterer) ParseAccessManagerUpdated(log types.Log) (*AlertsAccessManagerUpdated, error) {
	event := new(AlertsAccessManagerUpdated)
	if err := _Alerts.contract.UnpackLog(event, "AccessManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the Alerts contract.
type AlertsAdminChangedIterator struct {
	Event *AlertsAdminChanged // Event containing the contract specifics and raw log

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
func (it *AlertsAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsAdminChanged)
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
		it.Event = new(AlertsAdminChanged)
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
func (it *AlertsAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsAdminChanged represents a AdminChanged event raised by the Alerts contract.
type AlertsAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Alerts *AlertsFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*AlertsAdminChangedIterator, error) {

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &AlertsAdminChangedIterator{contract: _Alerts.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Alerts *AlertsFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *AlertsAdminChanged) (event.Subscription, error) {

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsAdminChanged)
				if err := _Alerts.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Alerts *AlertsFilterer) ParseAdminChanged(log types.Log) (*AlertsAdminChanged, error) {
	event := new(AlertsAdminChanged)
	if err := _Alerts.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsAlertBatchIterator is returned from FilterAlertBatch and is used to iterate over the raw logs and unpacked data for AlertBatch events raised by the Alerts contract.
type AlertsAlertBatchIterator struct {
	Event *AlertsAlertBatch // Event containing the contract specifics and raw log

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
func (it *AlertsAlertBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsAlertBatch)
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
		it.Event = new(AlertsAlertBatch)
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
func (it *AlertsAlertBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsAlertBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsAlertBatch represents a AlertBatch event raised by the Alerts contract.
type AlertsAlertBatch struct {
	AlertsId    [32]byte
	Scanner     common.Address
	ChainId     *big.Int
	BlockStart  *big.Int
	BlockEnd    *big.Int
	AlertCount  *big.Int
	MaxSeverity *big.Int
	Ref         string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAlertBatch is a free log retrieval operation binding the contract event 0x28be462f7201a7fb00d382965861c289c80b933d88532b002b59f45beaee83f2.
//
// Solidity: event AlertBatch(bytes32 indexed alertsId, address indexed scanner, uint256 indexed chainId, uint256 blockStart, uint256 blockEnd, uint256 alertCount, uint256 maxSeverity, string ref)
func (_Alerts *AlertsFilterer) FilterAlertBatch(opts *bind.FilterOpts, alertsId [][32]byte, scanner []common.Address, chainId []*big.Int) (*AlertsAlertBatchIterator, error) {

	var alertsIdRule []interface{}
	for _, alertsIdItem := range alertsId {
		alertsIdRule = append(alertsIdRule, alertsIdItem)
	}
	var scannerRule []interface{}
	for _, scannerItem := range scanner {
		scannerRule = append(scannerRule, scannerItem)
	}
	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "AlertBatch", alertsIdRule, scannerRule, chainIdRule)
	if err != nil {
		return nil, err
	}
	return &AlertsAlertBatchIterator{contract: _Alerts.contract, event: "AlertBatch", logs: logs, sub: sub}, nil
}

// WatchAlertBatch is a free log subscription operation binding the contract event 0x28be462f7201a7fb00d382965861c289c80b933d88532b002b59f45beaee83f2.
//
// Solidity: event AlertBatch(bytes32 indexed alertsId, address indexed scanner, uint256 indexed chainId, uint256 blockStart, uint256 blockEnd, uint256 alertCount, uint256 maxSeverity, string ref)
func (_Alerts *AlertsFilterer) WatchAlertBatch(opts *bind.WatchOpts, sink chan<- *AlertsAlertBatch, alertsId [][32]byte, scanner []common.Address, chainId []*big.Int) (event.Subscription, error) {

	var alertsIdRule []interface{}
	for _, alertsIdItem := range alertsId {
		alertsIdRule = append(alertsIdRule, alertsIdItem)
	}
	var scannerRule []interface{}
	for _, scannerItem := range scanner {
		scannerRule = append(scannerRule, scannerItem)
	}
	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "AlertBatch", alertsIdRule, scannerRule, chainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsAlertBatch)
				if err := _Alerts.contract.UnpackLog(event, "AlertBatch", log); err != nil {
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

// ParseAlertBatch is a log parse operation binding the contract event 0x28be462f7201a7fb00d382965861c289c80b933d88532b002b59f45beaee83f2.
//
// Solidity: event AlertBatch(bytes32 indexed alertsId, address indexed scanner, uint256 indexed chainId, uint256 blockStart, uint256 blockEnd, uint256 alertCount, uint256 maxSeverity, string ref)
func (_Alerts *AlertsFilterer) ParseAlertBatch(log types.Log) (*AlertsAlertBatch, error) {
	event := new(AlertsAlertBatch)
	if err := _Alerts.contract.UnpackLog(event, "AlertBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the Alerts contract.
type AlertsBeaconUpgradedIterator struct {
	Event *AlertsBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *AlertsBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsBeaconUpgraded)
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
		it.Event = new(AlertsBeaconUpgraded)
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
func (it *AlertsBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsBeaconUpgraded represents a BeaconUpgraded event raised by the Alerts contract.
type AlertsBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Alerts *AlertsFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*AlertsBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &AlertsBeaconUpgradedIterator{contract: _Alerts.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Alerts *AlertsFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *AlertsBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsBeaconUpgraded)
				if err := _Alerts.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Alerts *AlertsFilterer) ParseBeaconUpgraded(log types.Log) (*AlertsBeaconUpgraded, error) {
	event := new(AlertsBeaconUpgraded)
	if err := _Alerts.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsRouterUpdatedIterator is returned from FilterRouterUpdated and is used to iterate over the raw logs and unpacked data for RouterUpdated events raised by the Alerts contract.
type AlertsRouterUpdatedIterator struct {
	Event *AlertsRouterUpdated // Event containing the contract specifics and raw log

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
func (it *AlertsRouterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsRouterUpdated)
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
		it.Event = new(AlertsRouterUpdated)
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
func (it *AlertsRouterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsRouterUpdated represents a RouterUpdated event raised by the Alerts contract.
type AlertsRouterUpdated struct {
	Router common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRouterUpdated is a free log retrieval operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_Alerts *AlertsFilterer) FilterRouterUpdated(opts *bind.FilterOpts, router []common.Address) (*AlertsRouterUpdatedIterator, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "RouterUpdated", routerRule)
	if err != nil {
		return nil, err
	}
	return &AlertsRouterUpdatedIterator{contract: _Alerts.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

// WatchRouterUpdated is a free log subscription operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_Alerts *AlertsFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *AlertsRouterUpdated, router []common.Address) (event.Subscription, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "RouterUpdated", routerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsRouterUpdated)
				if err := _Alerts.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

// ParseRouterUpdated is a log parse operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_Alerts *AlertsFilterer) ParseRouterUpdated(log types.Log) (*AlertsRouterUpdated, error) {
	event := new(AlertsRouterUpdated)
	if err := _Alerts.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsScannerRegistryChangedIterator is returned from FilterScannerRegistryChanged and is used to iterate over the raw logs and unpacked data for ScannerRegistryChanged events raised by the Alerts contract.
type AlertsScannerRegistryChangedIterator struct {
	Event *AlertsScannerRegistryChanged // Event containing the contract specifics and raw log

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
func (it *AlertsScannerRegistryChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsScannerRegistryChanged)
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
		it.Event = new(AlertsScannerRegistryChanged)
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
func (it *AlertsScannerRegistryChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsScannerRegistryChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsScannerRegistryChanged represents a ScannerRegistryChanged event raised by the Alerts contract.
type AlertsScannerRegistryChanged struct {
	From common.Address
	To   common.Address
	By   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterScannerRegistryChanged is a free log retrieval operation binding the contract event 0x86d76b9fe0c8674850798bf28a78d85728ea9754f9272989d85d53d4622d7952.
//
// Solidity: event ScannerRegistryChanged(address from, address to, address by)
func (_Alerts *AlertsFilterer) FilterScannerRegistryChanged(opts *bind.FilterOpts) (*AlertsScannerRegistryChangedIterator, error) {

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "ScannerRegistryChanged")
	if err != nil {
		return nil, err
	}
	return &AlertsScannerRegistryChangedIterator{contract: _Alerts.contract, event: "ScannerRegistryChanged", logs: logs, sub: sub}, nil
}

// WatchScannerRegistryChanged is a free log subscription operation binding the contract event 0x86d76b9fe0c8674850798bf28a78d85728ea9754f9272989d85d53d4622d7952.
//
// Solidity: event ScannerRegistryChanged(address from, address to, address by)
func (_Alerts *AlertsFilterer) WatchScannerRegistryChanged(opts *bind.WatchOpts, sink chan<- *AlertsScannerRegistryChanged) (event.Subscription, error) {

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "ScannerRegistryChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsScannerRegistryChanged)
				if err := _Alerts.contract.UnpackLog(event, "ScannerRegistryChanged", log); err != nil {
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

// ParseScannerRegistryChanged is a log parse operation binding the contract event 0x86d76b9fe0c8674850798bf28a78d85728ea9754f9272989d85d53d4622d7952.
//
// Solidity: event ScannerRegistryChanged(address from, address to, address by)
func (_Alerts *AlertsFilterer) ParseScannerRegistryChanged(log types.Log) (*AlertsScannerRegistryChanged, error) {
	event := new(AlertsScannerRegistryChanged)
	if err := _Alerts.contract.UnpackLog(event, "ScannerRegistryChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlertsUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Alerts contract.
type AlertsUpgradedIterator struct {
	Event *AlertsUpgraded // Event containing the contract specifics and raw log

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
func (it *AlertsUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlertsUpgraded)
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
		it.Event = new(AlertsUpgraded)
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
func (it *AlertsUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlertsUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlertsUpgraded represents a Upgraded event raised by the Alerts contract.
type AlertsUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Alerts *AlertsFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*AlertsUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Alerts.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &AlertsUpgradedIterator{contract: _Alerts.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Alerts *AlertsFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *AlertsUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Alerts.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlertsUpgraded)
				if err := _Alerts.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Alerts *AlertsFilterer) ParseUpgraded(log types.Log) (*AlertsUpgraded, error) {
	event := new(AlertsUpgraded)
	if err := _Alerts.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
