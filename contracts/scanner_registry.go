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

// ScannerRegistryABI is the input ABI used to generate the binding from.
const ScannerRegistryABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"latest\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdminAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdminRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"AgentRegistryChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"latest\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"ScannerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"ScannerAdminAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"ScannerAdminRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"ScannerOwnershipTransfered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"latest\",\"type\":\"bool\"}],\"name\":\"addAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"addAgentAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"addScanner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addScannerAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"agentAdminAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"agentAdminLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"agentAt\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"agentExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"}],\"name\":\"agentLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"agentRef\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"agentUsingLatest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"agentVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"}],\"name\":\"getAgentListHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractMinimalForwarderUpgradeable\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"removeAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"removeAgentAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removeScannerAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"scannerAdminAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"}],\"name\":\"scannerAdminLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"scannerExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"scannerOwners\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setAgentRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferScannerOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_scanner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"latest\",\"type\":\"bool\"}],\"name\":\"updateAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ScannerRegistry is an auto generated Go binding around an Ethereum contract.
type ScannerRegistry struct {
	ScannerRegistryCaller     // Read-only binding to the contract
	ScannerRegistryTransactor // Write-only binding to the contract
	ScannerRegistryFilterer   // Log filterer for contract events
}

// ScannerRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ScannerRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ScannerRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ScannerRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ScannerRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ScannerRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ScannerRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ScannerRegistrySession struct {
	Contract     *ScannerRegistry  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ScannerRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ScannerRegistryCallerSession struct {
	Contract *ScannerRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ScannerRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ScannerRegistryTransactorSession struct {
	Contract     *ScannerRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ScannerRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ScannerRegistryRaw struct {
	Contract *ScannerRegistry // Generic contract binding to access the raw methods on
}

// ScannerRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ScannerRegistryCallerRaw struct {
	Contract *ScannerRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ScannerRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ScannerRegistryTransactorRaw struct {
	Contract *ScannerRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewScannerRegistry creates a new instance of ScannerRegistry, bound to a specific deployed contract.
func NewScannerRegistry(address common.Address, backend bind.ContractBackend) (*ScannerRegistry, error) {
	contract, err := bindScannerRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistry{ScannerRegistryCaller: ScannerRegistryCaller{contract: contract}, ScannerRegistryTransactor: ScannerRegistryTransactor{contract: contract}, ScannerRegistryFilterer: ScannerRegistryFilterer{contract: contract}}, nil
}

// NewScannerRegistryCaller creates a new read-only instance of ScannerRegistry, bound to a specific deployed contract.
func NewScannerRegistryCaller(address common.Address, caller bind.ContractCaller) (*ScannerRegistryCaller, error) {
	contract, err := bindScannerRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryCaller{contract: contract}, nil
}

// NewScannerRegistryTransactor creates a new write-only instance of ScannerRegistry, bound to a specific deployed contract.
func NewScannerRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ScannerRegistryTransactor, error) {
	contract, err := bindScannerRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryTransactor{contract: contract}, nil
}

// NewScannerRegistryFilterer creates a new log filterer instance of ScannerRegistry, bound to a specific deployed contract.
func NewScannerRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ScannerRegistryFilterer, error) {
	contract, err := bindScannerRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryFilterer{contract: contract}, nil
}

// bindScannerRegistry binds a generic wrapper to an already deployed contract.
func bindScannerRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ScannerRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ScannerRegistry *ScannerRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ScannerRegistry.Contract.ScannerRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ScannerRegistry *ScannerRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.ScannerRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ScannerRegistry *ScannerRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.ScannerRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ScannerRegistry *ScannerRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ScannerRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ScannerRegistry *ScannerRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ScannerRegistry *ScannerRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ScannerRegistry *ScannerRegistryCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ScannerRegistry *ScannerRegistrySession) ADMINROLE() ([32]byte, error) {
	return _ScannerRegistry.Contract.ADMINROLE(&_ScannerRegistry.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ScannerRegistry *ScannerRegistryCallerSession) ADMINROLE() ([32]byte, error) {
	return _ScannerRegistry.Contract.ADMINROLE(&_ScannerRegistry.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ScannerRegistry *ScannerRegistryCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ScannerRegistry *ScannerRegistrySession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ScannerRegistry.Contract.DEFAULTADMINROLE(&_ScannerRegistry.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ScannerRegistry *ScannerRegistryCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ScannerRegistry.Contract.DEFAULTADMINROLE(&_ScannerRegistry.CallOpts)
}

// AgentAdminAt is a free data retrieval call binding the contract method 0x2bc0fb12.
//
// Solidity: function agentAdminAt(address _scanner, bytes32 _agentId, uint256 index) view returns(address)
func (_ScannerRegistry *ScannerRegistryCaller) AgentAdminAt(opts *bind.CallOpts, _scanner common.Address, _agentId [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "agentAdminAt", _scanner, _agentId, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentAdminAt is a free data retrieval call binding the contract method 0x2bc0fb12.
//
// Solidity: function agentAdminAt(address _scanner, bytes32 _agentId, uint256 index) view returns(address)
func (_ScannerRegistry *ScannerRegistrySession) AgentAdminAt(_scanner common.Address, _agentId [32]byte, index *big.Int) (common.Address, error) {
	return _ScannerRegistry.Contract.AgentAdminAt(&_ScannerRegistry.CallOpts, _scanner, _agentId, index)
}

// AgentAdminAt is a free data retrieval call binding the contract method 0x2bc0fb12.
//
// Solidity: function agentAdminAt(address _scanner, bytes32 _agentId, uint256 index) view returns(address)
func (_ScannerRegistry *ScannerRegistryCallerSession) AgentAdminAt(_scanner common.Address, _agentId [32]byte, index *big.Int) (common.Address, error) {
	return _ScannerRegistry.Contract.AgentAdminAt(&_ScannerRegistry.CallOpts, _scanner, _agentId, index)
}

// AgentAdminLength is a free data retrieval call binding the contract method 0x5dc8b1ba.
//
// Solidity: function agentAdminLength(address _scanner, bytes32 _agentId) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCaller) AgentAdminLength(opts *bind.CallOpts, _scanner common.Address, _agentId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "agentAdminLength", _scanner, _agentId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentAdminLength is a free data retrieval call binding the contract method 0x5dc8b1ba.
//
// Solidity: function agentAdminLength(address _scanner, bytes32 _agentId) view returns(uint256)
func (_ScannerRegistry *ScannerRegistrySession) AgentAdminLength(_scanner common.Address, _agentId [32]byte) (*big.Int, error) {
	return _ScannerRegistry.Contract.AgentAdminLength(&_ScannerRegistry.CallOpts, _scanner, _agentId)
}

// AgentAdminLength is a free data retrieval call binding the contract method 0x5dc8b1ba.
//
// Solidity: function agentAdminLength(address _scanner, bytes32 _agentId) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCallerSession) AgentAdminLength(_scanner common.Address, _agentId [32]byte) (*big.Int, error) {
	return _ScannerRegistry.Contract.AgentAdminLength(&_ScannerRegistry.CallOpts, _scanner, _agentId)
}

// AgentAt is a free data retrieval call binding the contract method 0xc7cadabc.
//
// Solidity: function agentAt(address _scanner, uint256 index) view returns(bytes32, uint256, bool, string, bool)
func (_ScannerRegistry *ScannerRegistryCaller) AgentAt(opts *bind.CallOpts, _scanner common.Address, index *big.Int) ([32]byte, *big.Int, bool, string, bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "agentAt", _scanner, index)

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

// AgentAt is a free data retrieval call binding the contract method 0xc7cadabc.
//
// Solidity: function agentAt(address _scanner, uint256 index) view returns(bytes32, uint256, bool, string, bool)
func (_ScannerRegistry *ScannerRegistrySession) AgentAt(_scanner common.Address, index *big.Int) ([32]byte, *big.Int, bool, string, bool, error) {
	return _ScannerRegistry.Contract.AgentAt(&_ScannerRegistry.CallOpts, _scanner, index)
}

// AgentAt is a free data retrieval call binding the contract method 0xc7cadabc.
//
// Solidity: function agentAt(address _scanner, uint256 index) view returns(bytes32, uint256, bool, string, bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) AgentAt(_scanner common.Address, index *big.Int) ([32]byte, *big.Int, bool, string, bool, error) {
	return _ScannerRegistry.Contract.AgentAt(&_ScannerRegistry.CallOpts, _scanner, index)
}

// AgentExists is a free data retrieval call binding the contract method 0xa30375f7.
//
// Solidity: function agentExists(address _scanner, bytes32 _agentId) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCaller) AgentExists(opts *bind.CallOpts, _scanner common.Address, _agentId [32]byte) (bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "agentExists", _scanner, _agentId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AgentExists is a free data retrieval call binding the contract method 0xa30375f7.
//
// Solidity: function agentExists(address _scanner, bytes32 _agentId) view returns(bool)
func (_ScannerRegistry *ScannerRegistrySession) AgentExists(_scanner common.Address, _agentId [32]byte) (bool, error) {
	return _ScannerRegistry.Contract.AgentExists(&_ScannerRegistry.CallOpts, _scanner, _agentId)
}

// AgentExists is a free data retrieval call binding the contract method 0xa30375f7.
//
// Solidity: function agentExists(address _scanner, bytes32 _agentId) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) AgentExists(_scanner common.Address, _agentId [32]byte) (bool, error) {
	return _ScannerRegistry.Contract.AgentExists(&_ScannerRegistry.CallOpts, _scanner, _agentId)
}

// AgentLength is a free data retrieval call binding the contract method 0x4aa2f8fa.
//
// Solidity: function agentLength(address _scanner) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCaller) AgentLength(opts *bind.CallOpts, _scanner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "agentLength", _scanner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentLength is a free data retrieval call binding the contract method 0x4aa2f8fa.
//
// Solidity: function agentLength(address _scanner) view returns(uint256)
func (_ScannerRegistry *ScannerRegistrySession) AgentLength(_scanner common.Address) (*big.Int, error) {
	return _ScannerRegistry.Contract.AgentLength(&_ScannerRegistry.CallOpts, _scanner)
}

// AgentLength is a free data retrieval call binding the contract method 0x4aa2f8fa.
//
// Solidity: function agentLength(address _scanner) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCallerSession) AgentLength(_scanner common.Address) (*big.Int, error) {
	return _ScannerRegistry.Contract.AgentLength(&_ScannerRegistry.CallOpts, _scanner)
}

// AgentRef is a free data retrieval call binding the contract method 0x69282121.
//
// Solidity: function agentRef(address _scanner, bytes32 _agentId) view returns(string)
func (_ScannerRegistry *ScannerRegistryCaller) AgentRef(opts *bind.CallOpts, _scanner common.Address, _agentId [32]byte) (string, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "agentRef", _scanner, _agentId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AgentRef is a free data retrieval call binding the contract method 0x69282121.
//
// Solidity: function agentRef(address _scanner, bytes32 _agentId) view returns(string)
func (_ScannerRegistry *ScannerRegistrySession) AgentRef(_scanner common.Address, _agentId [32]byte) (string, error) {
	return _ScannerRegistry.Contract.AgentRef(&_ScannerRegistry.CallOpts, _scanner, _agentId)
}

// AgentRef is a free data retrieval call binding the contract method 0x69282121.
//
// Solidity: function agentRef(address _scanner, bytes32 _agentId) view returns(string)
func (_ScannerRegistry *ScannerRegistryCallerSession) AgentRef(_scanner common.Address, _agentId [32]byte) (string, error) {
	return _ScannerRegistry.Contract.AgentRef(&_ScannerRegistry.CallOpts, _scanner, _agentId)
}

// AgentUsingLatest is a free data retrieval call binding the contract method 0x781ee936.
//
// Solidity: function agentUsingLatest(address , bytes32 ) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCaller) AgentUsingLatest(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte) (bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "agentUsingLatest", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AgentUsingLatest is a free data retrieval call binding the contract method 0x781ee936.
//
// Solidity: function agentUsingLatest(address , bytes32 ) view returns(bool)
func (_ScannerRegistry *ScannerRegistrySession) AgentUsingLatest(arg0 common.Address, arg1 [32]byte) (bool, error) {
	return _ScannerRegistry.Contract.AgentUsingLatest(&_ScannerRegistry.CallOpts, arg0, arg1)
}

// AgentUsingLatest is a free data retrieval call binding the contract method 0x781ee936.
//
// Solidity: function agentUsingLatest(address , bytes32 ) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) AgentUsingLatest(arg0 common.Address, arg1 [32]byte) (bool, error) {
	return _ScannerRegistry.Contract.AgentUsingLatest(&_ScannerRegistry.CallOpts, arg0, arg1)
}

// AgentVersion is a free data retrieval call binding the contract method 0x4d3eb30b.
//
// Solidity: function agentVersion(address , bytes32 ) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCaller) AgentVersion(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "agentVersion", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentVersion is a free data retrieval call binding the contract method 0x4d3eb30b.
//
// Solidity: function agentVersion(address , bytes32 ) view returns(uint256)
func (_ScannerRegistry *ScannerRegistrySession) AgentVersion(arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	return _ScannerRegistry.Contract.AgentVersion(&_ScannerRegistry.CallOpts, arg0, arg1)
}

// AgentVersion is a free data retrieval call binding the contract method 0x4d3eb30b.
//
// Solidity: function agentVersion(address , bytes32 ) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCallerSession) AgentVersion(arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	return _ScannerRegistry.Contract.AgentVersion(&_ScannerRegistry.CallOpts, arg0, arg1)
}

// GetAgentListHash is a free data retrieval call binding the contract method 0x51cc51b2.
//
// Solidity: function getAgentListHash(address _scanner) view returns(bytes32)
func (_ScannerRegistry *ScannerRegistryCaller) GetAgentListHash(opts *bind.CallOpts, _scanner common.Address) ([32]byte, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "getAgentListHash", _scanner)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetAgentListHash is a free data retrieval call binding the contract method 0x51cc51b2.
//
// Solidity: function getAgentListHash(address _scanner) view returns(bytes32)
func (_ScannerRegistry *ScannerRegistrySession) GetAgentListHash(_scanner common.Address) ([32]byte, error) {
	return _ScannerRegistry.Contract.GetAgentListHash(&_ScannerRegistry.CallOpts, _scanner)
}

// GetAgentListHash is a free data retrieval call binding the contract method 0x51cc51b2.
//
// Solidity: function getAgentListHash(address _scanner) view returns(bytes32)
func (_ScannerRegistry *ScannerRegistryCallerSession) GetAgentListHash(_scanner common.Address) ([32]byte, error) {
	return _ScannerRegistry.Contract.GetAgentListHash(&_ScannerRegistry.CallOpts, _scanner)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ScannerRegistry *ScannerRegistryCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ScannerRegistry *ScannerRegistrySession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ScannerRegistry.Contract.GetRoleAdmin(&_ScannerRegistry.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ScannerRegistry *ScannerRegistryCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ScannerRegistry.Contract.GetRoleAdmin(&_ScannerRegistry.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ScannerRegistry *ScannerRegistrySession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ScannerRegistry.Contract.HasRole(&_ScannerRegistry.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ScannerRegistry.Contract.HasRole(&_ScannerRegistry.CallOpts, role, account)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCaller) IsTrustedForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "isTrustedForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_ScannerRegistry *ScannerRegistrySession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _ScannerRegistry.Contract.IsTrustedForwarder(&_ScannerRegistry.CallOpts, forwarder)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _ScannerRegistry.Contract.IsTrustedForwarder(&_ScannerRegistry.CallOpts, forwarder)
}

// ScannerAdminAt is a free data retrieval call binding the contract method 0x6271b6ed.
//
// Solidity: function scannerAdminAt(address _scanner, uint256 index) view returns(address)
func (_ScannerRegistry *ScannerRegistryCaller) ScannerAdminAt(opts *bind.CallOpts, _scanner common.Address, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "scannerAdminAt", _scanner, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ScannerAdminAt is a free data retrieval call binding the contract method 0x6271b6ed.
//
// Solidity: function scannerAdminAt(address _scanner, uint256 index) view returns(address)
func (_ScannerRegistry *ScannerRegistrySession) ScannerAdminAt(_scanner common.Address, index *big.Int) (common.Address, error) {
	return _ScannerRegistry.Contract.ScannerAdminAt(&_ScannerRegistry.CallOpts, _scanner, index)
}

// ScannerAdminAt is a free data retrieval call binding the contract method 0x6271b6ed.
//
// Solidity: function scannerAdminAt(address _scanner, uint256 index) view returns(address)
func (_ScannerRegistry *ScannerRegistryCallerSession) ScannerAdminAt(_scanner common.Address, index *big.Int) (common.Address, error) {
	return _ScannerRegistry.Contract.ScannerAdminAt(&_ScannerRegistry.CallOpts, _scanner, index)
}

// ScannerAdminLength is a free data retrieval call binding the contract method 0x74c45dbf.
//
// Solidity: function scannerAdminLength(address _scanner) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCaller) ScannerAdminLength(opts *bind.CallOpts, _scanner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "scannerAdminLength", _scanner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ScannerAdminLength is a free data retrieval call binding the contract method 0x74c45dbf.
//
// Solidity: function scannerAdminLength(address _scanner) view returns(uint256)
func (_ScannerRegistry *ScannerRegistrySession) ScannerAdminLength(_scanner common.Address) (*big.Int, error) {
	return _ScannerRegistry.Contract.ScannerAdminLength(&_ScannerRegistry.CallOpts, _scanner)
}

// ScannerAdminLength is a free data retrieval call binding the contract method 0x74c45dbf.
//
// Solidity: function scannerAdminLength(address _scanner) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCallerSession) ScannerAdminLength(_scanner common.Address) (*big.Int, error) {
	return _ScannerRegistry.Contract.ScannerAdminLength(&_ScannerRegistry.CallOpts, _scanner)
}

// ScannerExists is a free data retrieval call binding the contract method 0xb4891277.
//
// Solidity: function scannerExists(address ) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCaller) ScannerExists(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "scannerExists", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ScannerExists is a free data retrieval call binding the contract method 0xb4891277.
//
// Solidity: function scannerExists(address ) view returns(bool)
func (_ScannerRegistry *ScannerRegistrySession) ScannerExists(arg0 common.Address) (bool, error) {
	return _ScannerRegistry.Contract.ScannerExists(&_ScannerRegistry.CallOpts, arg0)
}

// ScannerExists is a free data retrieval call binding the contract method 0xb4891277.
//
// Solidity: function scannerExists(address ) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) ScannerExists(arg0 common.Address) (bool, error) {
	return _ScannerRegistry.Contract.ScannerExists(&_ScannerRegistry.CallOpts, arg0)
}

// ScannerOwners is a free data retrieval call binding the contract method 0x93a53a2f.
//
// Solidity: function scannerOwners(address ) view returns(address)
func (_ScannerRegistry *ScannerRegistryCaller) ScannerOwners(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "scannerOwners", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ScannerOwners is a free data retrieval call binding the contract method 0x93a53a2f.
//
// Solidity: function scannerOwners(address ) view returns(address)
func (_ScannerRegistry *ScannerRegistrySession) ScannerOwners(arg0 common.Address) (common.Address, error) {
	return _ScannerRegistry.Contract.ScannerOwners(&_ScannerRegistry.CallOpts, arg0)
}

// ScannerOwners is a free data retrieval call binding the contract method 0x93a53a2f.
//
// Solidity: function scannerOwners(address ) view returns(address)
func (_ScannerRegistry *ScannerRegistryCallerSession) ScannerOwners(arg0 common.Address) (common.Address, error) {
	return _ScannerRegistry.Contract.ScannerOwners(&_ScannerRegistry.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ScannerRegistry *ScannerRegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ScannerRegistry.Contract.SupportsInterface(&_ScannerRegistry.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ScannerRegistry.Contract.SupportsInterface(&_ScannerRegistry.CallOpts, interfaceId)
}

// AddAgent is a paid mutator transaction binding the contract method 0x822a323c.
//
// Solidity: function addAgent(address _scanner, bytes32 _agentId, uint256 version, bool latest) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) AddAgent(opts *bind.TransactOpts, _scanner common.Address, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "addAgent", _scanner, _agentId, version, latest)
}

// AddAgent is a paid mutator transaction binding the contract method 0x822a323c.
//
// Solidity: function addAgent(address _scanner, bytes32 _agentId, uint256 version, bool latest) returns()
func (_ScannerRegistry *ScannerRegistrySession) AddAgent(_scanner common.Address, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AddAgent(&_ScannerRegistry.TransactOpts, _scanner, _agentId, version, latest)
}

// AddAgent is a paid mutator transaction binding the contract method 0x822a323c.
//
// Solidity: function addAgent(address _scanner, bytes32 _agentId, uint256 version, bool latest) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) AddAgent(_scanner common.Address, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AddAgent(&_ScannerRegistry.TransactOpts, _scanner, _agentId, version, latest)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x37a5e835.
//
// Solidity: function addAgentAdmin(address _scanner, bytes32 _agentId, address _admin) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) AddAgentAdmin(opts *bind.TransactOpts, _scanner common.Address, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "addAgentAdmin", _scanner, _agentId, _admin)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x37a5e835.
//
// Solidity: function addAgentAdmin(address _scanner, bytes32 _agentId, address _admin) returns()
func (_ScannerRegistry *ScannerRegistrySession) AddAgentAdmin(_scanner common.Address, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AddAgentAdmin(&_ScannerRegistry.TransactOpts, _scanner, _agentId, _admin)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x37a5e835.
//
// Solidity: function addAgentAdmin(address _scanner, bytes32 _agentId, address _admin) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) AddAgentAdmin(_scanner common.Address, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AddAgentAdmin(&_ScannerRegistry.TransactOpts, _scanner, _agentId, _admin)
}

// AddScanner is a paid mutator transaction binding the contract method 0x3122512c.
//
// Solidity: function addScanner(address _owner) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) AddScanner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "addScanner", _owner)
}

// AddScanner is a paid mutator transaction binding the contract method 0x3122512c.
//
// Solidity: function addScanner(address _owner) returns()
func (_ScannerRegistry *ScannerRegistrySession) AddScanner(_owner common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AddScanner(&_ScannerRegistry.TransactOpts, _owner)
}

// AddScanner is a paid mutator transaction binding the contract method 0x3122512c.
//
// Solidity: function addScanner(address _owner) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) AddScanner(_owner common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AddScanner(&_ScannerRegistry.TransactOpts, _owner)
}

// AddScannerAdmin is a paid mutator transaction binding the contract method 0x7090ce6c.
//
// Solidity: function addScannerAdmin(address _scanner, address admin) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) AddScannerAdmin(opts *bind.TransactOpts, _scanner common.Address, admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "addScannerAdmin", _scanner, admin)
}

// AddScannerAdmin is a paid mutator transaction binding the contract method 0x7090ce6c.
//
// Solidity: function addScannerAdmin(address _scanner, address admin) returns()
func (_ScannerRegistry *ScannerRegistrySession) AddScannerAdmin(_scanner common.Address, admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AddScannerAdmin(&_ScannerRegistry.TransactOpts, _scanner, admin)
}

// AddScannerAdmin is a paid mutator transaction binding the contract method 0x7090ce6c.
//
// Solidity: function addScannerAdmin(address _scanner, address admin) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) AddScannerAdmin(_scanner common.Address, admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AddScannerAdmin(&_ScannerRegistry.TransactOpts, _scanner, admin)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ScannerRegistry *ScannerRegistrySession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.GrantRole(&_ScannerRegistry.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.GrantRole(&_ScannerRegistry.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address forwarder) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) Initialize(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "initialize", forwarder)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address forwarder) returns()
func (_ScannerRegistry *ScannerRegistrySession) Initialize(forwarder common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.Initialize(&_ScannerRegistry.TransactOpts, forwarder)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address forwarder) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) Initialize(forwarder common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.Initialize(&_ScannerRegistry.TransactOpts, forwarder)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x32784004.
//
// Solidity: function removeAgent(address _scanner, bytes32 _agentId) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) RemoveAgent(opts *bind.TransactOpts, _scanner common.Address, _agentId [32]byte) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "removeAgent", _scanner, _agentId)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x32784004.
//
// Solidity: function removeAgent(address _scanner, bytes32 _agentId) returns()
func (_ScannerRegistry *ScannerRegistrySession) RemoveAgent(_scanner common.Address, _agentId [32]byte) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.RemoveAgent(&_ScannerRegistry.TransactOpts, _scanner, _agentId)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x32784004.
//
// Solidity: function removeAgent(address _scanner, bytes32 _agentId) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) RemoveAgent(_scanner common.Address, _agentId [32]byte) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.RemoveAgent(&_ScannerRegistry.TransactOpts, _scanner, _agentId)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0x6c74d8ec.
//
// Solidity: function removeAgentAdmin(address _scanner, bytes32 _agentId, address _admin) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) RemoveAgentAdmin(opts *bind.TransactOpts, _scanner common.Address, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "removeAgentAdmin", _scanner, _agentId, _admin)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0x6c74d8ec.
//
// Solidity: function removeAgentAdmin(address _scanner, bytes32 _agentId, address _admin) returns()
func (_ScannerRegistry *ScannerRegistrySession) RemoveAgentAdmin(_scanner common.Address, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.RemoveAgentAdmin(&_ScannerRegistry.TransactOpts, _scanner, _agentId, _admin)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0x6c74d8ec.
//
// Solidity: function removeAgentAdmin(address _scanner, bytes32 _agentId, address _admin) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) RemoveAgentAdmin(_scanner common.Address, _agentId [32]byte, _admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.RemoveAgentAdmin(&_ScannerRegistry.TransactOpts, _scanner, _agentId, _admin)
}

// RemoveScannerAdmin is a paid mutator transaction binding the contract method 0xe6173357.
//
// Solidity: function removeScannerAdmin(address _scanner, address admin) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) RemoveScannerAdmin(opts *bind.TransactOpts, _scanner common.Address, admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "removeScannerAdmin", _scanner, admin)
}

// RemoveScannerAdmin is a paid mutator transaction binding the contract method 0xe6173357.
//
// Solidity: function removeScannerAdmin(address _scanner, address admin) returns()
func (_ScannerRegistry *ScannerRegistrySession) RemoveScannerAdmin(_scanner common.Address, admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.RemoveScannerAdmin(&_ScannerRegistry.TransactOpts, _scanner, admin)
}

// RemoveScannerAdmin is a paid mutator transaction binding the contract method 0xe6173357.
//
// Solidity: function removeScannerAdmin(address _scanner, address admin) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) RemoveScannerAdmin(_scanner common.Address, admin common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.RemoveScannerAdmin(&_ScannerRegistry.TransactOpts, _scanner, admin)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ScannerRegistry *ScannerRegistrySession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.RenounceRole(&_ScannerRegistry.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.RenounceRole(&_ScannerRegistry.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ScannerRegistry *ScannerRegistrySession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.RevokeRole(&_ScannerRegistry.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.RevokeRole(&_ScannerRegistry.TransactOpts, role, account)
}

// SetAgentRegistry is a paid mutator transaction binding the contract method 0x28342ecf.
//
// Solidity: function setAgentRegistry(address _addr) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) SetAgentRegistry(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "setAgentRegistry", _addr)
}

// SetAgentRegistry is a paid mutator transaction binding the contract method 0x28342ecf.
//
// Solidity: function setAgentRegistry(address _addr) returns()
func (_ScannerRegistry *ScannerRegistrySession) SetAgentRegistry(_addr common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetAgentRegistry(&_ScannerRegistry.TransactOpts, _addr)
}

// SetAgentRegistry is a paid mutator transaction binding the contract method 0x28342ecf.
//
// Solidity: function setAgentRegistry(address _addr) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) SetAgentRegistry(_addr common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetAgentRegistry(&_ScannerRegistry.TransactOpts, _addr)
}

// TransferScannerOwnership is a paid mutator transaction binding the contract method 0xe3acf437.
//
// Solidity: function transferScannerOwnership(address _scanner, address _to) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) TransferScannerOwnership(opts *bind.TransactOpts, _scanner common.Address, _to common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "transferScannerOwnership", _scanner, _to)
}

// TransferScannerOwnership is a paid mutator transaction binding the contract method 0xe3acf437.
//
// Solidity: function transferScannerOwnership(address _scanner, address _to) returns()
func (_ScannerRegistry *ScannerRegistrySession) TransferScannerOwnership(_scanner common.Address, _to common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.TransferScannerOwnership(&_ScannerRegistry.TransactOpts, _scanner, _to)
}

// TransferScannerOwnership is a paid mutator transaction binding the contract method 0xe3acf437.
//
// Solidity: function transferScannerOwnership(address _scanner, address _to) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) TransferScannerOwnership(_scanner common.Address, _to common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.TransferScannerOwnership(&_ScannerRegistry.TransactOpts, _scanner, _to)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0x1f50fa6d.
//
// Solidity: function updateAgent(address _scanner, bytes32 _agentId, uint256 version, bool latest) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) UpdateAgent(opts *bind.TransactOpts, _scanner common.Address, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "updateAgent", _scanner, _agentId, version, latest)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0x1f50fa6d.
//
// Solidity: function updateAgent(address _scanner, bytes32 _agentId, uint256 version, bool latest) returns()
func (_ScannerRegistry *ScannerRegistrySession) UpdateAgent(_scanner common.Address, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.UpdateAgent(&_ScannerRegistry.TransactOpts, _scanner, _agentId, version, latest)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0x1f50fa6d.
//
// Solidity: function updateAgent(address _scanner, bytes32 _agentId, uint256 version, bool latest) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) UpdateAgent(_scanner common.Address, _agentId [32]byte, version *big.Int, latest bool) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.UpdateAgent(&_ScannerRegistry.TransactOpts, _scanner, _agentId, version, latest)
}

// ScannerRegistryAgentAddedIterator is returned from FilterAgentAdded and is used to iterate over the raw logs and unpacked data for AgentAdded events raised by the ScannerRegistry contract.
type ScannerRegistryAgentAddedIterator struct {
	Event *ScannerRegistryAgentAdded // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryAgentAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryAgentAdded)
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
		it.Event = new(ScannerRegistryAgentAdded)
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
func (it *ScannerRegistryAgentAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryAgentAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryAgentAdded represents a AgentAdded event raised by the ScannerRegistry contract.
type ScannerRegistryAgentAdded struct {
	Scanner common.Address
	AgentId [32]byte
	Version *big.Int
	Latest  bool
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdded is a free log retrieval operation binding the contract event 0xcae65e9685779ff56dcf7b49328a09581f1644d32aebecd51e56ed7761e52107.
//
// Solidity: event AgentAdded(address scanner, bytes32 agentId, uint256 version, bool latest, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterAgentAdded(opts *bind.FilterOpts) (*ScannerRegistryAgentAddedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "AgentAdded")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryAgentAddedIterator{contract: _ScannerRegistry.contract, event: "AgentAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdded is a free log subscription operation binding the contract event 0xcae65e9685779ff56dcf7b49328a09581f1644d32aebecd51e56ed7761e52107.
//
// Solidity: event AgentAdded(address scanner, bytes32 agentId, uint256 version, bool latest, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchAgentAdded(opts *bind.WatchOpts, sink chan<- *ScannerRegistryAgentAdded) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "AgentAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryAgentAdded)
				if err := _ScannerRegistry.contract.UnpackLog(event, "AgentAdded", log); err != nil {
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

// ParseAgentAdded is a log parse operation binding the contract event 0xcae65e9685779ff56dcf7b49328a09581f1644d32aebecd51e56ed7761e52107.
//
// Solidity: event AgentAdded(address scanner, bytes32 agentId, uint256 version, bool latest, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseAgentAdded(log types.Log) (*ScannerRegistryAgentAdded, error) {
	event := new(ScannerRegistryAgentAdded)
	if err := _ScannerRegistry.contract.UnpackLog(event, "AgentAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryAgentAdminAddedIterator is returned from FilterAgentAdminAdded and is used to iterate over the raw logs and unpacked data for AgentAdminAdded events raised by the ScannerRegistry contract.
type ScannerRegistryAgentAdminAddedIterator struct {
	Event *ScannerRegistryAgentAdminAdded // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryAgentAdminAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryAgentAdminAdded)
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
		it.Event = new(ScannerRegistryAgentAdminAdded)
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
func (it *ScannerRegistryAgentAdminAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryAgentAdminAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryAgentAdminAdded represents a AgentAdminAdded event raised by the ScannerRegistry contract.
type ScannerRegistryAgentAdminAdded struct {
	Scanner common.Address
	AgentId [32]byte
	Admin   common.Address
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdminAdded is a free log retrieval operation binding the contract event 0x2b0df8efc673280b2db12c377c053917c699f2fe33596e8fb5741f0b4fa3ea78.
//
// Solidity: event AgentAdminAdded(address scanner, bytes32 agentId, address admin, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterAgentAdminAdded(opts *bind.FilterOpts) (*ScannerRegistryAgentAdminAddedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "AgentAdminAdded")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryAgentAdminAddedIterator{contract: _ScannerRegistry.contract, event: "AgentAdminAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdminAdded is a free log subscription operation binding the contract event 0x2b0df8efc673280b2db12c377c053917c699f2fe33596e8fb5741f0b4fa3ea78.
//
// Solidity: event AgentAdminAdded(address scanner, bytes32 agentId, address admin, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchAgentAdminAdded(opts *bind.WatchOpts, sink chan<- *ScannerRegistryAgentAdminAdded) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "AgentAdminAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryAgentAdminAdded)
				if err := _ScannerRegistry.contract.UnpackLog(event, "AgentAdminAdded", log); err != nil {
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

// ParseAgentAdminAdded is a log parse operation binding the contract event 0x2b0df8efc673280b2db12c377c053917c699f2fe33596e8fb5741f0b4fa3ea78.
//
// Solidity: event AgentAdminAdded(address scanner, bytes32 agentId, address admin, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseAgentAdminAdded(log types.Log) (*ScannerRegistryAgentAdminAdded, error) {
	event := new(ScannerRegistryAgentAdminAdded)
	if err := _ScannerRegistry.contract.UnpackLog(event, "AgentAdminAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryAgentAdminRemovedIterator is returned from FilterAgentAdminRemoved and is used to iterate over the raw logs and unpacked data for AgentAdminRemoved events raised by the ScannerRegistry contract.
type ScannerRegistryAgentAdminRemovedIterator struct {
	Event *ScannerRegistryAgentAdminRemoved // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryAgentAdminRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryAgentAdminRemoved)
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
		it.Event = new(ScannerRegistryAgentAdminRemoved)
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
func (it *ScannerRegistryAgentAdminRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryAgentAdminRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryAgentAdminRemoved represents a AgentAdminRemoved event raised by the ScannerRegistry contract.
type ScannerRegistryAgentAdminRemoved struct {
	Scanner common.Address
	AgentId [32]byte
	Admin   common.Address
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdminRemoved is a free log retrieval operation binding the contract event 0xd5c42544cfb37e2b94af896761d5ec8609d3511ef873e757b70e3f635e82285b.
//
// Solidity: event AgentAdminRemoved(address scanner, bytes32 agentId, address admin, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterAgentAdminRemoved(opts *bind.FilterOpts) (*ScannerRegistryAgentAdminRemovedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "AgentAdminRemoved")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryAgentAdminRemovedIterator{contract: _ScannerRegistry.contract, event: "AgentAdminRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentAdminRemoved is a free log subscription operation binding the contract event 0xd5c42544cfb37e2b94af896761d5ec8609d3511ef873e757b70e3f635e82285b.
//
// Solidity: event AgentAdminRemoved(address scanner, bytes32 agentId, address admin, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchAgentAdminRemoved(opts *bind.WatchOpts, sink chan<- *ScannerRegistryAgentAdminRemoved) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "AgentAdminRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryAgentAdminRemoved)
				if err := _ScannerRegistry.contract.UnpackLog(event, "AgentAdminRemoved", log); err != nil {
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

// ParseAgentAdminRemoved is a log parse operation binding the contract event 0xd5c42544cfb37e2b94af896761d5ec8609d3511ef873e757b70e3f635e82285b.
//
// Solidity: event AgentAdminRemoved(address scanner, bytes32 agentId, address admin, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseAgentAdminRemoved(log types.Log) (*ScannerRegistryAgentAdminRemoved, error) {
	event := new(ScannerRegistryAgentAdminRemoved)
	if err := _ScannerRegistry.contract.UnpackLog(event, "AgentAdminRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryAgentRegistryChangedIterator is returned from FilterAgentRegistryChanged and is used to iterate over the raw logs and unpacked data for AgentRegistryChanged events raised by the ScannerRegistry contract.
type ScannerRegistryAgentRegistryChangedIterator struct {
	Event *ScannerRegistryAgentRegistryChanged // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryAgentRegistryChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryAgentRegistryChanged)
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
		it.Event = new(ScannerRegistryAgentRegistryChanged)
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
func (it *ScannerRegistryAgentRegistryChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryAgentRegistryChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryAgentRegistryChanged represents a AgentRegistryChanged event raised by the ScannerRegistry contract.
type ScannerRegistryAgentRegistryChanged struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAgentRegistryChanged is a free log retrieval operation binding the contract event 0x2d35b9d2e073404ea7e01756776790bfe7c7e525789490689493de40955f5632.
//
// Solidity: event AgentRegistryChanged(address from, address to)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterAgentRegistryChanged(opts *bind.FilterOpts) (*ScannerRegistryAgentRegistryChangedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "AgentRegistryChanged")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryAgentRegistryChangedIterator{contract: _ScannerRegistry.contract, event: "AgentRegistryChanged", logs: logs, sub: sub}, nil
}

// WatchAgentRegistryChanged is a free log subscription operation binding the contract event 0x2d35b9d2e073404ea7e01756776790bfe7c7e525789490689493de40955f5632.
//
// Solidity: event AgentRegistryChanged(address from, address to)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchAgentRegistryChanged(opts *bind.WatchOpts, sink chan<- *ScannerRegistryAgentRegistryChanged) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "AgentRegistryChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryAgentRegistryChanged)
				if err := _ScannerRegistry.contract.UnpackLog(event, "AgentRegistryChanged", log); err != nil {
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
func (_ScannerRegistry *ScannerRegistryFilterer) ParseAgentRegistryChanged(log types.Log) (*ScannerRegistryAgentRegistryChanged, error) {
	event := new(ScannerRegistryAgentRegistryChanged)
	if err := _ScannerRegistry.contract.UnpackLog(event, "AgentRegistryChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryAgentRemovedIterator is returned from FilterAgentRemoved and is used to iterate over the raw logs and unpacked data for AgentRemoved events raised by the ScannerRegistry contract.
type ScannerRegistryAgentRemovedIterator struct {
	Event *ScannerRegistryAgentRemoved // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryAgentRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryAgentRemoved)
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
		it.Event = new(ScannerRegistryAgentRemoved)
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
func (it *ScannerRegistryAgentRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryAgentRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryAgentRemoved represents a AgentRemoved event raised by the ScannerRegistry contract.
type ScannerRegistryAgentRemoved struct {
	Scanner common.Address
	AgentId [32]byte
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentRemoved is a free log retrieval operation binding the contract event 0xdd06fa26adb1d18eabdccff199e230a4e909f17bebf02835c9d118b5cff9a09e.
//
// Solidity: event AgentRemoved(address scanner, bytes32 agentId, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterAgentRemoved(opts *bind.FilterOpts) (*ScannerRegistryAgentRemovedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "AgentRemoved")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryAgentRemovedIterator{contract: _ScannerRegistry.contract, event: "AgentRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentRemoved is a free log subscription operation binding the contract event 0xdd06fa26adb1d18eabdccff199e230a4e909f17bebf02835c9d118b5cff9a09e.
//
// Solidity: event AgentRemoved(address scanner, bytes32 agentId, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchAgentRemoved(opts *bind.WatchOpts, sink chan<- *ScannerRegistryAgentRemoved) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "AgentRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryAgentRemoved)
				if err := _ScannerRegistry.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
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

// ParseAgentRemoved is a log parse operation binding the contract event 0xdd06fa26adb1d18eabdccff199e230a4e909f17bebf02835c9d118b5cff9a09e.
//
// Solidity: event AgentRemoved(address scanner, bytes32 agentId, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseAgentRemoved(log types.Log) (*ScannerRegistryAgentRemoved, error) {
	event := new(ScannerRegistryAgentRemoved)
	if err := _ScannerRegistry.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryAgentUpdatedIterator is returned from FilterAgentUpdated and is used to iterate over the raw logs and unpacked data for AgentUpdated events raised by the ScannerRegistry contract.
type ScannerRegistryAgentUpdatedIterator struct {
	Event *ScannerRegistryAgentUpdated // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryAgentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryAgentUpdated)
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
		it.Event = new(ScannerRegistryAgentUpdated)
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
func (it *ScannerRegistryAgentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryAgentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryAgentUpdated represents a AgentUpdated event raised by the ScannerRegistry contract.
type ScannerRegistryAgentUpdated struct {
	Scanner common.Address
	AgentId [32]byte
	Version *big.Int
	Latest  bool
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentUpdated is a free log retrieval operation binding the contract event 0x6089b03be22d9d86d1c7e19f861bde296fe19236ac93c105d140e79d489e394a.
//
// Solidity: event AgentUpdated(address scanner, bytes32 agentId, uint256 version, bool latest, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterAgentUpdated(opts *bind.FilterOpts) (*ScannerRegistryAgentUpdatedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "AgentUpdated")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryAgentUpdatedIterator{contract: _ScannerRegistry.contract, event: "AgentUpdated", logs: logs, sub: sub}, nil
}

// WatchAgentUpdated is a free log subscription operation binding the contract event 0x6089b03be22d9d86d1c7e19f861bde296fe19236ac93c105d140e79d489e394a.
//
// Solidity: event AgentUpdated(address scanner, bytes32 agentId, uint256 version, bool latest, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchAgentUpdated(opts *bind.WatchOpts, sink chan<- *ScannerRegistryAgentUpdated) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "AgentUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryAgentUpdated)
				if err := _ScannerRegistry.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
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

// ParseAgentUpdated is a log parse operation binding the contract event 0x6089b03be22d9d86d1c7e19f861bde296fe19236ac93c105d140e79d489e394a.
//
// Solidity: event AgentUpdated(address scanner, bytes32 agentId, uint256 version, bool latest, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseAgentUpdated(log types.Log) (*ScannerRegistryAgentUpdated, error) {
	event := new(ScannerRegistryAgentUpdated)
	if err := _ScannerRegistry.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ScannerRegistry contract.
type ScannerRegistryRoleAdminChangedIterator struct {
	Event *ScannerRegistryRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryRoleAdminChanged)
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
		it.Event = new(ScannerRegistryRoleAdminChanged)
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
func (it *ScannerRegistryRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryRoleAdminChanged represents a RoleAdminChanged event raised by the ScannerRegistry contract.
type ScannerRegistryRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ScannerRegistryRoleAdminChangedIterator, error) {

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

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryRoleAdminChangedIterator{contract: _ScannerRegistry.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ScannerRegistryRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryRoleAdminChanged)
				if err := _ScannerRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_ScannerRegistry *ScannerRegistryFilterer) ParseRoleAdminChanged(log types.Log) (*ScannerRegistryRoleAdminChanged, error) {
	event := new(ScannerRegistryRoleAdminChanged)
	if err := _ScannerRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ScannerRegistry contract.
type ScannerRegistryRoleGrantedIterator struct {
	Event *ScannerRegistryRoleGranted // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryRoleGranted)
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
		it.Event = new(ScannerRegistryRoleGranted)
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
func (it *ScannerRegistryRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryRoleGranted represents a RoleGranted event raised by the ScannerRegistry contract.
type ScannerRegistryRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ScannerRegistryRoleGrantedIterator, error) {

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

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryRoleGrantedIterator{contract: _ScannerRegistry.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ScannerRegistryRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryRoleGranted)
				if err := _ScannerRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_ScannerRegistry *ScannerRegistryFilterer) ParseRoleGranted(log types.Log) (*ScannerRegistryRoleGranted, error) {
	event := new(ScannerRegistryRoleGranted)
	if err := _ScannerRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ScannerRegistry contract.
type ScannerRegistryRoleRevokedIterator struct {
	Event *ScannerRegistryRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryRoleRevoked)
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
		it.Event = new(ScannerRegistryRoleRevoked)
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
func (it *ScannerRegistryRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryRoleRevoked represents a RoleRevoked event raised by the ScannerRegistry contract.
type ScannerRegistryRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ScannerRegistryRoleRevokedIterator, error) {

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

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryRoleRevokedIterator{contract: _ScannerRegistry.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ScannerRegistryRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryRoleRevoked)
				if err := _ScannerRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_ScannerRegistry *ScannerRegistryFilterer) ParseRoleRevoked(log types.Log) (*ScannerRegistryRoleRevoked, error) {
	event := new(ScannerRegistryRoleRevoked)
	if err := _ScannerRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryScannerAddedIterator is returned from FilterScannerAdded and is used to iterate over the raw logs and unpacked data for ScannerAdded events raised by the ScannerRegistry contract.
type ScannerRegistryScannerAddedIterator struct {
	Event *ScannerRegistryScannerAdded // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryScannerAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryScannerAdded)
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
		it.Event = new(ScannerRegistryScannerAdded)
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
func (it *ScannerRegistryScannerAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryScannerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryScannerAdded represents a ScannerAdded event raised by the ScannerRegistry contract.
type ScannerRegistryScannerAdded struct {
	Scanner common.Address
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterScannerAdded is a free log retrieval operation binding the contract event 0x1f2879e55165661fe1587e9878a094323b1b2604e2ccc622906a3caa5af5ecae.
//
// Solidity: event ScannerAdded(address scanner, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterScannerAdded(opts *bind.FilterOpts) (*ScannerRegistryScannerAddedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "ScannerAdded")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryScannerAddedIterator{contract: _ScannerRegistry.contract, event: "ScannerAdded", logs: logs, sub: sub}, nil
}

// WatchScannerAdded is a free log subscription operation binding the contract event 0x1f2879e55165661fe1587e9878a094323b1b2604e2ccc622906a3caa5af5ecae.
//
// Solidity: event ScannerAdded(address scanner, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchScannerAdded(opts *bind.WatchOpts, sink chan<- *ScannerRegistryScannerAdded) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "ScannerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryScannerAdded)
				if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerAdded", log); err != nil {
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

// ParseScannerAdded is a log parse operation binding the contract event 0x1f2879e55165661fe1587e9878a094323b1b2604e2ccc622906a3caa5af5ecae.
//
// Solidity: event ScannerAdded(address scanner, address by)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseScannerAdded(log types.Log) (*ScannerRegistryScannerAdded, error) {
	event := new(ScannerRegistryScannerAdded)
	if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryScannerAdminAddedIterator is returned from FilterScannerAdminAdded and is used to iterate over the raw logs and unpacked data for ScannerAdminAdded events raised by the ScannerRegistry contract.
type ScannerRegistryScannerAdminAddedIterator struct {
	Event *ScannerRegistryScannerAdminAdded // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryScannerAdminAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryScannerAdminAdded)
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
		it.Event = new(ScannerRegistryScannerAdminAdded)
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
func (it *ScannerRegistryScannerAdminAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryScannerAdminAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryScannerAdminAdded represents a ScannerAdminAdded event raised by the ScannerRegistry contract.
type ScannerRegistryScannerAdminAdded struct {
	Scanner common.Address
	Addr    common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterScannerAdminAdded is a free log retrieval operation binding the contract event 0xcbb62b060e77e89ac31e9fd7ec0a03690304dead091d562b0154bd29b530e0a9.
//
// Solidity: event ScannerAdminAdded(address scanner, address addr)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterScannerAdminAdded(opts *bind.FilterOpts) (*ScannerRegistryScannerAdminAddedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "ScannerAdminAdded")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryScannerAdminAddedIterator{contract: _ScannerRegistry.contract, event: "ScannerAdminAdded", logs: logs, sub: sub}, nil
}

// WatchScannerAdminAdded is a free log subscription operation binding the contract event 0xcbb62b060e77e89ac31e9fd7ec0a03690304dead091d562b0154bd29b530e0a9.
//
// Solidity: event ScannerAdminAdded(address scanner, address addr)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchScannerAdminAdded(opts *bind.WatchOpts, sink chan<- *ScannerRegistryScannerAdminAdded) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "ScannerAdminAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryScannerAdminAdded)
				if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerAdminAdded", log); err != nil {
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

// ParseScannerAdminAdded is a log parse operation binding the contract event 0xcbb62b060e77e89ac31e9fd7ec0a03690304dead091d562b0154bd29b530e0a9.
//
// Solidity: event ScannerAdminAdded(address scanner, address addr)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseScannerAdminAdded(log types.Log) (*ScannerRegistryScannerAdminAdded, error) {
	event := new(ScannerRegistryScannerAdminAdded)
	if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerAdminAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryScannerAdminRemovedIterator is returned from FilterScannerAdminRemoved and is used to iterate over the raw logs and unpacked data for ScannerAdminRemoved events raised by the ScannerRegistry contract.
type ScannerRegistryScannerAdminRemovedIterator struct {
	Event *ScannerRegistryScannerAdminRemoved // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryScannerAdminRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryScannerAdminRemoved)
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
		it.Event = new(ScannerRegistryScannerAdminRemoved)
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
func (it *ScannerRegistryScannerAdminRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryScannerAdminRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryScannerAdminRemoved represents a ScannerAdminRemoved event raised by the ScannerRegistry contract.
type ScannerRegistryScannerAdminRemoved struct {
	Scanner common.Address
	Addr    common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterScannerAdminRemoved is a free log retrieval operation binding the contract event 0x5ba22ac69a875368bb837a9d99307fc62477a5ac90bf69d9722c6c1991f82a25.
//
// Solidity: event ScannerAdminRemoved(address scanner, address addr)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterScannerAdminRemoved(opts *bind.FilterOpts) (*ScannerRegistryScannerAdminRemovedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "ScannerAdminRemoved")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryScannerAdminRemovedIterator{contract: _ScannerRegistry.contract, event: "ScannerAdminRemoved", logs: logs, sub: sub}, nil
}

// WatchScannerAdminRemoved is a free log subscription operation binding the contract event 0x5ba22ac69a875368bb837a9d99307fc62477a5ac90bf69d9722c6c1991f82a25.
//
// Solidity: event ScannerAdminRemoved(address scanner, address addr)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchScannerAdminRemoved(opts *bind.WatchOpts, sink chan<- *ScannerRegistryScannerAdminRemoved) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "ScannerAdminRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryScannerAdminRemoved)
				if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerAdminRemoved", log); err != nil {
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

// ParseScannerAdminRemoved is a log parse operation binding the contract event 0x5ba22ac69a875368bb837a9d99307fc62477a5ac90bf69d9722c6c1991f82a25.
//
// Solidity: event ScannerAdminRemoved(address scanner, address addr)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseScannerAdminRemoved(log types.Log) (*ScannerRegistryScannerAdminRemoved, error) {
	event := new(ScannerRegistryScannerAdminRemoved)
	if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerAdminRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryScannerOwnershipTransferedIterator is returned from FilterScannerOwnershipTransfered and is used to iterate over the raw logs and unpacked data for ScannerOwnershipTransfered events raised by the ScannerRegistry contract.
type ScannerRegistryScannerOwnershipTransferedIterator struct {
	Event *ScannerRegistryScannerOwnershipTransfered // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryScannerOwnershipTransferedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryScannerOwnershipTransfered)
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
		it.Event = new(ScannerRegistryScannerOwnershipTransfered)
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
func (it *ScannerRegistryScannerOwnershipTransferedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryScannerOwnershipTransferedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryScannerOwnershipTransfered represents a ScannerOwnershipTransfered event raised by the ScannerRegistry contract.
type ScannerRegistryScannerOwnershipTransfered struct {
	Scanner common.Address
	From    common.Address
	To      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterScannerOwnershipTransfered is a free log retrieval operation binding the contract event 0x1edce13ffb6c6fe8dfdfb55ffc977f12572840735b1bea1c2c7767a60eae23bb.
//
// Solidity: event ScannerOwnershipTransfered(address scanner, address from, address to)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterScannerOwnershipTransfered(opts *bind.FilterOpts) (*ScannerRegistryScannerOwnershipTransferedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "ScannerOwnershipTransfered")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryScannerOwnershipTransferedIterator{contract: _ScannerRegistry.contract, event: "ScannerOwnershipTransfered", logs: logs, sub: sub}, nil
}

// WatchScannerOwnershipTransfered is a free log subscription operation binding the contract event 0x1edce13ffb6c6fe8dfdfb55ffc977f12572840735b1bea1c2c7767a60eae23bb.
//
// Solidity: event ScannerOwnershipTransfered(address scanner, address from, address to)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchScannerOwnershipTransfered(opts *bind.WatchOpts, sink chan<- *ScannerRegistryScannerOwnershipTransfered) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "ScannerOwnershipTransfered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryScannerOwnershipTransfered)
				if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerOwnershipTransfered", log); err != nil {
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

// ParseScannerOwnershipTransfered is a log parse operation binding the contract event 0x1edce13ffb6c6fe8dfdfb55ffc977f12572840735b1bea1c2c7767a60eae23bb.
//
// Solidity: event ScannerOwnershipTransfered(address scanner, address from, address to)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseScannerOwnershipTransfered(log types.Log) (*ScannerRegistryScannerOwnershipTransfered, error) {
	event := new(ScannerRegistryScannerOwnershipTransfered)
	if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerOwnershipTransfered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
