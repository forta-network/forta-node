// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract_scanner_node_version

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

// ScannerNodeVersionMetaData contains all meta data concerning the ScannerNodeVersion contract.
var ScannerNodeVersionMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MissingRole\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddressManager\",\"type\":\"address\"}],\"name\":\"AccessManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newVersion\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"oldVersion\",\"type\":\"string\"}],\"name\":\"ScannerNodeVersionUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"__manager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"__router\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"scannerNodeVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newManager\",\"type\":\"address\"}],\"name\":\"setAccessManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ensRegistry\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"ensName\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newVersion\",\"type\":\"string\"}],\"name\":\"setScannerNodeVersion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ScannerNodeVersionABI is the input ABI used to generate the binding from.
// Deprecated: Use ScannerNodeVersionMetaData.ABI instead.
var ScannerNodeVersionABI = ScannerNodeVersionMetaData.ABI

// ScannerNodeVersion is an auto generated Go binding around an Ethereum contract.
type ScannerNodeVersion struct {
	ScannerNodeVersionCaller     // Read-only binding to the contract
	ScannerNodeVersionTransactor // Write-only binding to the contract
	ScannerNodeVersionFilterer   // Log filterer for contract events
}

// ScannerNodeVersionCaller is an auto generated read-only Go binding around an Ethereum contract.
type ScannerNodeVersionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ScannerNodeVersionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ScannerNodeVersionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ScannerNodeVersionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ScannerNodeVersionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ScannerNodeVersionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ScannerNodeVersionSession struct {
	Contract     *ScannerNodeVersion // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ScannerNodeVersionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ScannerNodeVersionCallerSession struct {
	Contract *ScannerNodeVersionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ScannerNodeVersionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ScannerNodeVersionTransactorSession struct {
	Contract     *ScannerNodeVersionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ScannerNodeVersionRaw is an auto generated low-level Go binding around an Ethereum contract.
type ScannerNodeVersionRaw struct {
	Contract *ScannerNodeVersion // Generic contract binding to access the raw methods on
}

// ScannerNodeVersionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ScannerNodeVersionCallerRaw struct {
	Contract *ScannerNodeVersionCaller // Generic read-only contract binding to access the raw methods on
}

// ScannerNodeVersionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ScannerNodeVersionTransactorRaw struct {
	Contract *ScannerNodeVersionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewScannerNodeVersion creates a new instance of ScannerNodeVersion, bound to a specific deployed contract.
func NewScannerNodeVersion(address common.Address, backend bind.ContractBackend) (*ScannerNodeVersion, error) {
	contract, err := bindScannerNodeVersion(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ScannerNodeVersion{ScannerNodeVersionCaller: ScannerNodeVersionCaller{contract: contract}, ScannerNodeVersionTransactor: ScannerNodeVersionTransactor{contract: contract}, ScannerNodeVersionFilterer: ScannerNodeVersionFilterer{contract: contract}}, nil
}

// NewScannerNodeVersionCaller creates a new read-only instance of ScannerNodeVersion, bound to a specific deployed contract.
func NewScannerNodeVersionCaller(address common.Address, caller bind.ContractCaller) (*ScannerNodeVersionCaller, error) {
	contract, err := bindScannerNodeVersion(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ScannerNodeVersionCaller{contract: contract}, nil
}

// NewScannerNodeVersionTransactor creates a new write-only instance of ScannerNodeVersion, bound to a specific deployed contract.
func NewScannerNodeVersionTransactor(address common.Address, transactor bind.ContractTransactor) (*ScannerNodeVersionTransactor, error) {
	contract, err := bindScannerNodeVersion(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ScannerNodeVersionTransactor{contract: contract}, nil
}

// NewScannerNodeVersionFilterer creates a new log filterer instance of ScannerNodeVersion, bound to a specific deployed contract.
func NewScannerNodeVersionFilterer(address common.Address, filterer bind.ContractFilterer) (*ScannerNodeVersionFilterer, error) {
	contract, err := bindScannerNodeVersion(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ScannerNodeVersionFilterer{contract: contract}, nil
}

// bindScannerNodeVersion binds a generic wrapper to an already deployed contract.
func bindScannerNodeVersion(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ScannerNodeVersionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ScannerNodeVersion *ScannerNodeVersionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ScannerNodeVersion.Contract.ScannerNodeVersionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ScannerNodeVersion *ScannerNodeVersionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.ScannerNodeVersionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ScannerNodeVersion *ScannerNodeVersionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.ScannerNodeVersionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ScannerNodeVersion *ScannerNodeVersionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ScannerNodeVersion.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ScannerNodeVersion *ScannerNodeVersionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ScannerNodeVersion *ScannerNodeVersionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.contract.Transact(opts, method, params...)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_ScannerNodeVersion *ScannerNodeVersionCaller) IsTrustedForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _ScannerNodeVersion.contract.Call(opts, &out, "isTrustedForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_ScannerNodeVersion *ScannerNodeVersionSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _ScannerNodeVersion.Contract.IsTrustedForwarder(&_ScannerNodeVersion.CallOpts, forwarder)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_ScannerNodeVersion *ScannerNodeVersionCallerSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _ScannerNodeVersion.Contract.IsTrustedForwarder(&_ScannerNodeVersion.CallOpts, forwarder)
}

// ScannerNodeVersion is a free data retrieval call binding the contract method 0x345db3e1.
//
// Solidity: function scannerNodeVersion() view returns(string)
func (_ScannerNodeVersion *ScannerNodeVersionCaller) ScannerNodeVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ScannerNodeVersion.contract.Call(opts, &out, "scannerNodeVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ScannerNodeVersion is a free data retrieval call binding the contract method 0x345db3e1.
//
// Solidity: function scannerNodeVersion() view returns(string)
func (_ScannerNodeVersion *ScannerNodeVersionSession) ScannerNodeVersion() (string, error) {
	return _ScannerNodeVersion.Contract.ScannerNodeVersion(&_ScannerNodeVersion.CallOpts)
}

// ScannerNodeVersion is a free data retrieval call binding the contract method 0x345db3e1.
//
// Solidity: function scannerNodeVersion() view returns(string)
func (_ScannerNodeVersion *ScannerNodeVersionCallerSession) ScannerNodeVersion() (string, error) {
	return _ScannerNodeVersion.Contract.ScannerNodeVersion(&_ScannerNodeVersion.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_ScannerNodeVersion *ScannerNodeVersionCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ScannerNodeVersion.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_ScannerNodeVersion *ScannerNodeVersionSession) Version() (string, error) {
	return _ScannerNodeVersion.Contract.Version(&_ScannerNodeVersion.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_ScannerNodeVersion *ScannerNodeVersionCallerSession) Version() (string, error) {
	return _ScannerNodeVersion.Contract.Version(&_ScannerNodeVersion.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address __manager, address __router) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactor) Initialize(opts *bind.TransactOpts, __manager common.Address, __router common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.contract.Transact(opts, "initialize", __manager, __router)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address __manager, address __router) returns()
func (_ScannerNodeVersion *ScannerNodeVersionSession) Initialize(__manager common.Address, __router common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.Initialize(&_ScannerNodeVersion.TransactOpts, __manager, __router)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address __manager, address __router) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactorSession) Initialize(__manager common.Address, __router common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.Initialize(&_ScannerNodeVersion.TransactOpts, __manager, __router)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ScannerNodeVersion *ScannerNodeVersionTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _ScannerNodeVersion.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ScannerNodeVersion *ScannerNodeVersionSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.Multicall(&_ScannerNodeVersion.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ScannerNodeVersion *ScannerNodeVersionTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.Multicall(&_ScannerNodeVersion.TransactOpts, data)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactor) SetAccessManager(opts *bind.TransactOpts, newManager common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.contract.Transact(opts, "setAccessManager", newManager)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_ScannerNodeVersion *ScannerNodeVersionSession) SetAccessManager(newManager common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.SetAccessManager(&_ScannerNodeVersion.TransactOpts, newManager)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactorSession) SetAccessManager(newManager common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.SetAccessManager(&_ScannerNodeVersion.TransactOpts, newManager)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactor) SetName(opts *bind.TransactOpts, ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _ScannerNodeVersion.contract.Transact(opts, "setName", ensRegistry, ensName)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_ScannerNodeVersion *ScannerNodeVersionSession) SetName(ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.SetName(&_ScannerNodeVersion.TransactOpts, ensRegistry, ensName)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactorSession) SetName(ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.SetName(&_ScannerNodeVersion.TransactOpts, ensRegistry, ensName)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.contract.Transact(opts, "setRouter", newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_ScannerNodeVersion *ScannerNodeVersionSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.SetRouter(&_ScannerNodeVersion.TransactOpts, newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.SetRouter(&_ScannerNodeVersion.TransactOpts, newRouter)
}

// SetScannerNodeVersion is a paid mutator transaction binding the contract method 0x94bb55a2.
//
// Solidity: function setScannerNodeVersion(string newVersion) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactor) SetScannerNodeVersion(opts *bind.TransactOpts, newVersion string) (*types.Transaction, error) {
	return _ScannerNodeVersion.contract.Transact(opts, "setScannerNodeVersion", newVersion)
}

// SetScannerNodeVersion is a paid mutator transaction binding the contract method 0x94bb55a2.
//
// Solidity: function setScannerNodeVersion(string newVersion) returns()
func (_ScannerNodeVersion *ScannerNodeVersionSession) SetScannerNodeVersion(newVersion string) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.SetScannerNodeVersion(&_ScannerNodeVersion.TransactOpts, newVersion)
}

// SetScannerNodeVersion is a paid mutator transaction binding the contract method 0x94bb55a2.
//
// Solidity: function setScannerNodeVersion(string newVersion) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactorSession) SetScannerNodeVersion(newVersion string) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.SetScannerNodeVersion(&_ScannerNodeVersion.TransactOpts, newVersion)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ScannerNodeVersion *ScannerNodeVersionSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.UpgradeTo(&_ScannerNodeVersion.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.UpgradeTo(&_ScannerNodeVersion.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ScannerNodeVersion.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ScannerNodeVersion *ScannerNodeVersionSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.UpgradeToAndCall(&_ScannerNodeVersion.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ScannerNodeVersion *ScannerNodeVersionTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ScannerNodeVersion.Contract.UpgradeToAndCall(&_ScannerNodeVersion.TransactOpts, newImplementation, data)
}

// ScannerNodeVersionAccessManagerUpdatedIterator is returned from FilterAccessManagerUpdated and is used to iterate over the raw logs and unpacked data for AccessManagerUpdated events raised by the ScannerNodeVersion contract.
type ScannerNodeVersionAccessManagerUpdatedIterator struct {
	Event *ScannerNodeVersionAccessManagerUpdated // Event containing the contract specifics and raw log

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
func (it *ScannerNodeVersionAccessManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerNodeVersionAccessManagerUpdated)
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
		it.Event = new(ScannerNodeVersionAccessManagerUpdated)
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
func (it *ScannerNodeVersionAccessManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerNodeVersionAccessManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerNodeVersionAccessManagerUpdated represents a AccessManagerUpdated event raised by the ScannerNodeVersion contract.
type ScannerNodeVersionAccessManagerUpdated struct {
	NewAddressManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterAccessManagerUpdated is a free log retrieval operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) FilterAccessManagerUpdated(opts *bind.FilterOpts, newAddressManager []common.Address) (*ScannerNodeVersionAccessManagerUpdatedIterator, error) {

	var newAddressManagerRule []interface{}
	for _, newAddressManagerItem := range newAddressManager {
		newAddressManagerRule = append(newAddressManagerRule, newAddressManagerItem)
	}

	logs, sub, err := _ScannerNodeVersion.contract.FilterLogs(opts, "AccessManagerUpdated", newAddressManagerRule)
	if err != nil {
		return nil, err
	}
	return &ScannerNodeVersionAccessManagerUpdatedIterator{contract: _ScannerNodeVersion.contract, event: "AccessManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchAccessManagerUpdated is a free log subscription operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) WatchAccessManagerUpdated(opts *bind.WatchOpts, sink chan<- *ScannerNodeVersionAccessManagerUpdated, newAddressManager []common.Address) (event.Subscription, error) {

	var newAddressManagerRule []interface{}
	for _, newAddressManagerItem := range newAddressManager {
		newAddressManagerRule = append(newAddressManagerRule, newAddressManagerItem)
	}

	logs, sub, err := _ScannerNodeVersion.contract.WatchLogs(opts, "AccessManagerUpdated", newAddressManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerNodeVersionAccessManagerUpdated)
				if err := _ScannerNodeVersion.contract.UnpackLog(event, "AccessManagerUpdated", log); err != nil {
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
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) ParseAccessManagerUpdated(log types.Log) (*ScannerNodeVersionAccessManagerUpdated, error) {
	event := new(ScannerNodeVersionAccessManagerUpdated)
	if err := _ScannerNodeVersion.contract.UnpackLog(event, "AccessManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerNodeVersionAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ScannerNodeVersion contract.
type ScannerNodeVersionAdminChangedIterator struct {
	Event *ScannerNodeVersionAdminChanged // Event containing the contract specifics and raw log

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
func (it *ScannerNodeVersionAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerNodeVersionAdminChanged)
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
		it.Event = new(ScannerNodeVersionAdminChanged)
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
func (it *ScannerNodeVersionAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerNodeVersionAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerNodeVersionAdminChanged represents a AdminChanged event raised by the ScannerNodeVersion contract.
type ScannerNodeVersionAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ScannerNodeVersionAdminChangedIterator, error) {

	logs, sub, err := _ScannerNodeVersion.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ScannerNodeVersionAdminChangedIterator{contract: _ScannerNodeVersion.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ScannerNodeVersionAdminChanged) (event.Subscription, error) {

	logs, sub, err := _ScannerNodeVersion.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerNodeVersionAdminChanged)
				if err := _ScannerNodeVersion.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) ParseAdminChanged(log types.Log) (*ScannerNodeVersionAdminChanged, error) {
	event := new(ScannerNodeVersionAdminChanged)
	if err := _ScannerNodeVersion.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerNodeVersionBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the ScannerNodeVersion contract.
type ScannerNodeVersionBeaconUpgradedIterator struct {
	Event *ScannerNodeVersionBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *ScannerNodeVersionBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerNodeVersionBeaconUpgraded)
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
		it.Event = new(ScannerNodeVersionBeaconUpgraded)
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
func (it *ScannerNodeVersionBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerNodeVersionBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerNodeVersionBeaconUpgraded represents a BeaconUpgraded event raised by the ScannerNodeVersion contract.
type ScannerNodeVersionBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ScannerNodeVersionBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ScannerNodeVersion.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ScannerNodeVersionBeaconUpgradedIterator{contract: _ScannerNodeVersion.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ScannerNodeVersionBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ScannerNodeVersion.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerNodeVersionBeaconUpgraded)
				if err := _ScannerNodeVersion.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) ParseBeaconUpgraded(log types.Log) (*ScannerNodeVersionBeaconUpgraded, error) {
	event := new(ScannerNodeVersionBeaconUpgraded)
	if err := _ScannerNodeVersion.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerNodeVersionRouterUpdatedIterator is returned from FilterRouterUpdated and is used to iterate over the raw logs and unpacked data for RouterUpdated events raised by the ScannerNodeVersion contract.
type ScannerNodeVersionRouterUpdatedIterator struct {
	Event *ScannerNodeVersionRouterUpdated // Event containing the contract specifics and raw log

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
func (it *ScannerNodeVersionRouterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerNodeVersionRouterUpdated)
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
		it.Event = new(ScannerNodeVersionRouterUpdated)
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
func (it *ScannerNodeVersionRouterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerNodeVersionRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerNodeVersionRouterUpdated represents a RouterUpdated event raised by the ScannerNodeVersion contract.
type ScannerNodeVersionRouterUpdated struct {
	Router common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRouterUpdated is a free log retrieval operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) FilterRouterUpdated(opts *bind.FilterOpts, router []common.Address) (*ScannerNodeVersionRouterUpdatedIterator, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _ScannerNodeVersion.contract.FilterLogs(opts, "RouterUpdated", routerRule)
	if err != nil {
		return nil, err
	}
	return &ScannerNodeVersionRouterUpdatedIterator{contract: _ScannerNodeVersion.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

// WatchRouterUpdated is a free log subscription operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *ScannerNodeVersionRouterUpdated, router []common.Address) (event.Subscription, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _ScannerNodeVersion.contract.WatchLogs(opts, "RouterUpdated", routerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerNodeVersionRouterUpdated)
				if err := _ScannerNodeVersion.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) ParseRouterUpdated(log types.Log) (*ScannerNodeVersionRouterUpdated, error) {
	event := new(ScannerNodeVersionRouterUpdated)
	if err := _ScannerNodeVersion.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerNodeVersionScannerNodeVersionUpdatedIterator is returned from FilterScannerNodeVersionUpdated and is used to iterate over the raw logs and unpacked data for ScannerNodeVersionUpdated events raised by the ScannerNodeVersion contract.
type ScannerNodeVersionScannerNodeVersionUpdatedIterator struct {
	Event *ScannerNodeVersionScannerNodeVersionUpdated // Event containing the contract specifics and raw log

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
func (it *ScannerNodeVersionScannerNodeVersionUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerNodeVersionScannerNodeVersionUpdated)
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
		it.Event = new(ScannerNodeVersionScannerNodeVersionUpdated)
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
func (it *ScannerNodeVersionScannerNodeVersionUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerNodeVersionScannerNodeVersionUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerNodeVersionScannerNodeVersionUpdated represents a ScannerNodeVersionUpdated event raised by the ScannerNodeVersion contract.
type ScannerNodeVersionScannerNodeVersionUpdated struct {
	NewVersion string
	OldVersion string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterScannerNodeVersionUpdated is a free log retrieval operation binding the contract event 0xaa8f23ad4857a5d22df0189ebc6150d51b3152d8c621c3e8d24c66387897819b.
//
// Solidity: event ScannerNodeVersionUpdated(string newVersion, string oldVersion)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) FilterScannerNodeVersionUpdated(opts *bind.FilterOpts) (*ScannerNodeVersionScannerNodeVersionUpdatedIterator, error) {

	logs, sub, err := _ScannerNodeVersion.contract.FilterLogs(opts, "ScannerNodeVersionUpdated")
	if err != nil {
		return nil, err
	}
	return &ScannerNodeVersionScannerNodeVersionUpdatedIterator{contract: _ScannerNodeVersion.contract, event: "ScannerNodeVersionUpdated", logs: logs, sub: sub}, nil
}

// WatchScannerNodeVersionUpdated is a free log subscription operation binding the contract event 0xaa8f23ad4857a5d22df0189ebc6150d51b3152d8c621c3e8d24c66387897819b.
//
// Solidity: event ScannerNodeVersionUpdated(string newVersion, string oldVersion)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) WatchScannerNodeVersionUpdated(opts *bind.WatchOpts, sink chan<- *ScannerNodeVersionScannerNodeVersionUpdated) (event.Subscription, error) {

	logs, sub, err := _ScannerNodeVersion.contract.WatchLogs(opts, "ScannerNodeVersionUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerNodeVersionScannerNodeVersionUpdated)
				if err := _ScannerNodeVersion.contract.UnpackLog(event, "ScannerNodeVersionUpdated", log); err != nil {
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

// ParseScannerNodeVersionUpdated is a log parse operation binding the contract event 0xaa8f23ad4857a5d22df0189ebc6150d51b3152d8c621c3e8d24c66387897819b.
//
// Solidity: event ScannerNodeVersionUpdated(string newVersion, string oldVersion)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) ParseScannerNodeVersionUpdated(log types.Log) (*ScannerNodeVersionScannerNodeVersionUpdated, error) {
	event := new(ScannerNodeVersionScannerNodeVersionUpdated)
	if err := _ScannerNodeVersion.contract.UnpackLog(event, "ScannerNodeVersionUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerNodeVersionUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ScannerNodeVersion contract.
type ScannerNodeVersionUpgradedIterator struct {
	Event *ScannerNodeVersionUpgraded // Event containing the contract specifics and raw log

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
func (it *ScannerNodeVersionUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerNodeVersionUpgraded)
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
		it.Event = new(ScannerNodeVersionUpgraded)
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
func (it *ScannerNodeVersionUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerNodeVersionUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerNodeVersionUpgraded represents a Upgraded event raised by the ScannerNodeVersion contract.
type ScannerNodeVersionUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ScannerNodeVersionUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ScannerNodeVersion.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ScannerNodeVersionUpgradedIterator{contract: _ScannerNodeVersion.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ScannerNodeVersionUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ScannerNodeVersion.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerNodeVersionUpgraded)
				if err := _ScannerNodeVersion.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ScannerNodeVersion *ScannerNodeVersionFilterer) ParseUpgraded(log types.Log) (*ScannerNodeVersionUpgraded, error) {
	event := new(ScannerNodeVersionUpgraded)
	if err := _ScannerNodeVersion.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
