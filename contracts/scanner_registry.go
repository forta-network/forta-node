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

// ScannerRegistryMetaData contains all meta data concerning the ScannerRegistry contract.
var ScannerRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddressManager\",\"type\":\"address\"}],\"name\":\"AccessManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"ManagerEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"enumScannerRegistryEnable.Permission\",\"name\":\"permission\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"ScannerEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"name\":\"ScannerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newstakeController\",\"type\":\"address\"}],\"name\":\"StakeControllerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"name\":\"_register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"name\":\"adminRegister\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"name\":\"adminUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"enumScannerRegistryEnable.Permission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"disableScanner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"enumScannerRegistryEnable.Permission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"enableScanner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getManagerAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"getManagerCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"getScanner\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStakeController\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"__manager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"__router\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"__name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"__symbol\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"}],\"name\":\"isManager\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"isRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"}],\"name\":\"scannerAddressToId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newManager\",\"type\":\"address\"}],\"name\":\"setAccessManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"setManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ensRegistry\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"ensName\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stakeController\",\"type\":\"address\"}],\"name\":\"setStakeController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ScannerRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ScannerRegistryMetaData.ABI instead.
var ScannerRegistryABI = ScannerRegistryMetaData.ABI

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

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ScannerRegistry *ScannerRegistrySession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ScannerRegistry.Contract.BalanceOf(&_ScannerRegistry.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ScannerRegistry.Contract.BalanceOf(&_ScannerRegistry.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ScannerRegistry *ScannerRegistryCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ScannerRegistry *ScannerRegistrySession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ScannerRegistry.Contract.GetApproved(&_ScannerRegistry.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ScannerRegistry *ScannerRegistryCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ScannerRegistry.Contract.GetApproved(&_ScannerRegistry.CallOpts, tokenId)
}

// GetManagerAt is a free data retrieval call binding the contract method 0x8e79a369.
//
// Solidity: function getManagerAt(uint256 scannerId, uint256 index) view returns(address)
func (_ScannerRegistry *ScannerRegistryCaller) GetManagerAt(opts *bind.CallOpts, scannerId *big.Int, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "getManagerAt", scannerId, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetManagerAt is a free data retrieval call binding the contract method 0x8e79a369.
//
// Solidity: function getManagerAt(uint256 scannerId, uint256 index) view returns(address)
func (_ScannerRegistry *ScannerRegistrySession) GetManagerAt(scannerId *big.Int, index *big.Int) (common.Address, error) {
	return _ScannerRegistry.Contract.GetManagerAt(&_ScannerRegistry.CallOpts, scannerId, index)
}

// GetManagerAt is a free data retrieval call binding the contract method 0x8e79a369.
//
// Solidity: function getManagerAt(uint256 scannerId, uint256 index) view returns(address)
func (_ScannerRegistry *ScannerRegistryCallerSession) GetManagerAt(scannerId *big.Int, index *big.Int) (common.Address, error) {
	return _ScannerRegistry.Contract.GetManagerAt(&_ScannerRegistry.CallOpts, scannerId, index)
}

// GetManagerCount is a free data retrieval call binding the contract method 0xe11cf71e.
//
// Solidity: function getManagerCount(uint256 scannerId) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCaller) GetManagerCount(opts *bind.CallOpts, scannerId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "getManagerCount", scannerId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetManagerCount is a free data retrieval call binding the contract method 0xe11cf71e.
//
// Solidity: function getManagerCount(uint256 scannerId) view returns(uint256)
func (_ScannerRegistry *ScannerRegistrySession) GetManagerCount(scannerId *big.Int) (*big.Int, error) {
	return _ScannerRegistry.Contract.GetManagerCount(&_ScannerRegistry.CallOpts, scannerId)
}

// GetManagerCount is a free data retrieval call binding the contract method 0xe11cf71e.
//
// Solidity: function getManagerCount(uint256 scannerId) view returns(uint256)
func (_ScannerRegistry *ScannerRegistryCallerSession) GetManagerCount(scannerId *big.Int) (*big.Int, error) {
	return _ScannerRegistry.Contract.GetManagerCount(&_ScannerRegistry.CallOpts, scannerId)
}

// GetScanner is a free data retrieval call binding the contract method 0xa97fe43e.
//
// Solidity: function getScanner(uint256 scannerId) view returns(uint256 chainId, string metadata)
func (_ScannerRegistry *ScannerRegistryCaller) GetScanner(opts *bind.CallOpts, scannerId *big.Int) (struct {
	ChainId  *big.Int
	Metadata string
}, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "getScanner", scannerId)

	outstruct := new(struct {
		ChainId  *big.Int
		Metadata string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Metadata = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

// GetScanner is a free data retrieval call binding the contract method 0xa97fe43e.
//
// Solidity: function getScanner(uint256 scannerId) view returns(uint256 chainId, string metadata)
func (_ScannerRegistry *ScannerRegistrySession) GetScanner(scannerId *big.Int) (struct {
	ChainId  *big.Int
	Metadata string
}, error) {
	return _ScannerRegistry.Contract.GetScanner(&_ScannerRegistry.CallOpts, scannerId)
}

// GetScanner is a free data retrieval call binding the contract method 0xa97fe43e.
//
// Solidity: function getScanner(uint256 scannerId) view returns(uint256 chainId, string metadata)
func (_ScannerRegistry *ScannerRegistryCallerSession) GetScanner(scannerId *big.Int) (struct {
	ChainId  *big.Int
	Metadata string
}, error) {
	return _ScannerRegistry.Contract.GetScanner(&_ScannerRegistry.CallOpts, scannerId)
}

// GetStakeController is a free data retrieval call binding the contract method 0xaebb5150.
//
// Solidity: function getStakeController() view returns(address)
func (_ScannerRegistry *ScannerRegistryCaller) GetStakeController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "getStakeController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetStakeController is a free data retrieval call binding the contract method 0xaebb5150.
//
// Solidity: function getStakeController() view returns(address)
func (_ScannerRegistry *ScannerRegistrySession) GetStakeController() (common.Address, error) {
	return _ScannerRegistry.Contract.GetStakeController(&_ScannerRegistry.CallOpts)
}

// GetStakeController is a free data retrieval call binding the contract method 0xaebb5150.
//
// Solidity: function getStakeController() view returns(address)
func (_ScannerRegistry *ScannerRegistryCallerSession) GetStakeController() (common.Address, error) {
	return _ScannerRegistry.Contract.GetStakeController(&_ScannerRegistry.CallOpts)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ScannerRegistry *ScannerRegistrySession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ScannerRegistry.Contract.IsApprovedForAll(&_ScannerRegistry.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ScannerRegistry.Contract.IsApprovedForAll(&_ScannerRegistry.CallOpts, owner, operator)
}

// IsEnabled is a free data retrieval call binding the contract method 0xc783034c.
//
// Solidity: function isEnabled(uint256 scannerId) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCaller) IsEnabled(opts *bind.CallOpts, scannerId *big.Int) (bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "isEnabled", scannerId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEnabled is a free data retrieval call binding the contract method 0xc783034c.
//
// Solidity: function isEnabled(uint256 scannerId) view returns(bool)
func (_ScannerRegistry *ScannerRegistrySession) IsEnabled(scannerId *big.Int) (bool, error) {
	return _ScannerRegistry.Contract.IsEnabled(&_ScannerRegistry.CallOpts, scannerId)
}

// IsEnabled is a free data retrieval call binding the contract method 0xc783034c.
//
// Solidity: function isEnabled(uint256 scannerId) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) IsEnabled(scannerId *big.Int) (bool, error) {
	return _ScannerRegistry.Contract.IsEnabled(&_ScannerRegistry.CallOpts, scannerId)
}

// IsManager is a free data retrieval call binding the contract method 0x773ed13c.
//
// Solidity: function isManager(uint256 scannerId, address manager) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCaller) IsManager(opts *bind.CallOpts, scannerId *big.Int, manager common.Address) (bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "isManager", scannerId, manager)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsManager is a free data retrieval call binding the contract method 0x773ed13c.
//
// Solidity: function isManager(uint256 scannerId, address manager) view returns(bool)
func (_ScannerRegistry *ScannerRegistrySession) IsManager(scannerId *big.Int, manager common.Address) (bool, error) {
	return _ScannerRegistry.Contract.IsManager(&_ScannerRegistry.CallOpts, scannerId, manager)
}

// IsManager is a free data retrieval call binding the contract method 0x773ed13c.
//
// Solidity: function isManager(uint256 scannerId, address manager) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) IsManager(scannerId *big.Int, manager common.Address) (bool, error) {
	return _ScannerRegistry.Contract.IsManager(&_ScannerRegistry.CallOpts, scannerId, manager)
}

// IsRegistered is a free data retrieval call binding the contract method 0x579a6988.
//
// Solidity: function isRegistered(uint256 scannerId) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCaller) IsRegistered(opts *bind.CallOpts, scannerId *big.Int) (bool, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "isRegistered", scannerId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRegistered is a free data retrieval call binding the contract method 0x579a6988.
//
// Solidity: function isRegistered(uint256 scannerId) view returns(bool)
func (_ScannerRegistry *ScannerRegistrySession) IsRegistered(scannerId *big.Int) (bool, error) {
	return _ScannerRegistry.Contract.IsRegistered(&_ScannerRegistry.CallOpts, scannerId)
}

// IsRegistered is a free data retrieval call binding the contract method 0x579a6988.
//
// Solidity: function isRegistered(uint256 scannerId) view returns(bool)
func (_ScannerRegistry *ScannerRegistryCallerSession) IsRegistered(scannerId *big.Int) (bool, error) {
	return _ScannerRegistry.Contract.IsRegistered(&_ScannerRegistry.CallOpts, scannerId)
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

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ScannerRegistry *ScannerRegistryCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ScannerRegistry *ScannerRegistrySession) Name() (string, error) {
	return _ScannerRegistry.Contract.Name(&_ScannerRegistry.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ScannerRegistry *ScannerRegistryCallerSession) Name() (string, error) {
	return _ScannerRegistry.Contract.Name(&_ScannerRegistry.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ScannerRegistry *ScannerRegistryCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ScannerRegistry *ScannerRegistrySession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ScannerRegistry.Contract.OwnerOf(&_ScannerRegistry.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ScannerRegistry *ScannerRegistryCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ScannerRegistry.Contract.OwnerOf(&_ScannerRegistry.CallOpts, tokenId)
}

// ScannerAddressToId is a free data retrieval call binding the contract method 0x82fe1bcc.
//
// Solidity: function scannerAddressToId(address scanner) pure returns(uint256)
func (_ScannerRegistry *ScannerRegistryCaller) ScannerAddressToId(opts *bind.CallOpts, scanner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "scannerAddressToId", scanner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ScannerAddressToId is a free data retrieval call binding the contract method 0x82fe1bcc.
//
// Solidity: function scannerAddressToId(address scanner) pure returns(uint256)
func (_ScannerRegistry *ScannerRegistrySession) ScannerAddressToId(scanner common.Address) (*big.Int, error) {
	return _ScannerRegistry.Contract.ScannerAddressToId(&_ScannerRegistry.CallOpts, scanner)
}

// ScannerAddressToId is a free data retrieval call binding the contract method 0x82fe1bcc.
//
// Solidity: function scannerAddressToId(address scanner) pure returns(uint256)
func (_ScannerRegistry *ScannerRegistryCallerSession) ScannerAddressToId(scanner common.Address) (*big.Int, error) {
	return _ScannerRegistry.Contract.ScannerAddressToId(&_ScannerRegistry.CallOpts, scanner)
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

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ScannerRegistry *ScannerRegistryCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ScannerRegistry *ScannerRegistrySession) Symbol() (string, error) {
	return _ScannerRegistry.Contract.Symbol(&_ScannerRegistry.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ScannerRegistry *ScannerRegistryCallerSession) Symbol() (string, error) {
	return _ScannerRegistry.Contract.Symbol(&_ScannerRegistry.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ScannerRegistry *ScannerRegistryCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ScannerRegistry *ScannerRegistrySession) TokenURI(tokenId *big.Int) (string, error) {
	return _ScannerRegistry.Contract.TokenURI(&_ScannerRegistry.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ScannerRegistry *ScannerRegistryCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ScannerRegistry.Contract.TokenURI(&_ScannerRegistry.CallOpts, tokenId)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_ScannerRegistry *ScannerRegistryCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ScannerRegistry.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_ScannerRegistry *ScannerRegistrySession) Version() (string, error) {
	return _ScannerRegistry.Contract.Version(&_ScannerRegistry.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_ScannerRegistry *ScannerRegistryCallerSession) Version() (string, error) {
	return _ScannerRegistry.Contract.Version(&_ScannerRegistry.CallOpts)
}

// UnderscoreRegister is a paid mutator transaction binding the contract method 0x713740f2.
//
// Solidity: function _register(address scanner, address owner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) UnderscoreRegister(opts *bind.TransactOpts, scanner common.Address, owner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "_register", scanner, owner, chainId, metadata)
}

// UnderscoreRegister is a paid mutator transaction binding the contract method 0x713740f2.
//
// Solidity: function _register(address scanner, address owner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistrySession) UnderscoreRegister(scanner common.Address, owner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.UnderscoreRegister(&_ScannerRegistry.TransactOpts, scanner, owner, chainId, metadata)
}

// UnderscoreRegister is a paid mutator transaction binding the contract method 0x713740f2.
//
// Solidity: function _register(address scanner, address owner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) UnderscoreRegister(scanner common.Address, owner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.UnderscoreRegister(&_ScannerRegistry.TransactOpts, scanner, owner, chainId, metadata)
}

// AdminRegister is a paid mutator transaction binding the contract method 0x2a91fb34.
//
// Solidity: function adminRegister(address scanner, address owner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) AdminRegister(opts *bind.TransactOpts, scanner common.Address, owner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "adminRegister", scanner, owner, chainId, metadata)
}

// AdminRegister is a paid mutator transaction binding the contract method 0x2a91fb34.
//
// Solidity: function adminRegister(address scanner, address owner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistrySession) AdminRegister(scanner common.Address, owner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AdminRegister(&_ScannerRegistry.TransactOpts, scanner, owner, chainId, metadata)
}

// AdminRegister is a paid mutator transaction binding the contract method 0x2a91fb34.
//
// Solidity: function adminRegister(address scanner, address owner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) AdminRegister(scanner common.Address, owner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AdminRegister(&_ScannerRegistry.TransactOpts, scanner, owner, chainId, metadata)
}

// AdminUpdate is a paid mutator transaction binding the contract method 0xc2dae01d.
//
// Solidity: function adminUpdate(address scanner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) AdminUpdate(opts *bind.TransactOpts, scanner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "adminUpdate", scanner, chainId, metadata)
}

// AdminUpdate is a paid mutator transaction binding the contract method 0xc2dae01d.
//
// Solidity: function adminUpdate(address scanner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistrySession) AdminUpdate(scanner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AdminUpdate(&_ScannerRegistry.TransactOpts, scanner, chainId, metadata)
}

// AdminUpdate is a paid mutator transaction binding the contract method 0xc2dae01d.
//
// Solidity: function adminUpdate(address scanner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) AdminUpdate(scanner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.AdminUpdate(&_ScannerRegistry.TransactOpts, scanner, chainId, metadata)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ScannerRegistry *ScannerRegistrySession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.Approve(&_ScannerRegistry.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.Approve(&_ScannerRegistry.TransactOpts, to, tokenId)
}

// DisableScanner is a paid mutator transaction binding the contract method 0x59238297.
//
// Solidity: function disableScanner(uint256 scannerId, uint8 permission) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) DisableScanner(opts *bind.TransactOpts, scannerId *big.Int, permission uint8) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "disableScanner", scannerId, permission)
}

// DisableScanner is a paid mutator transaction binding the contract method 0x59238297.
//
// Solidity: function disableScanner(uint256 scannerId, uint8 permission) returns()
func (_ScannerRegistry *ScannerRegistrySession) DisableScanner(scannerId *big.Int, permission uint8) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.DisableScanner(&_ScannerRegistry.TransactOpts, scannerId, permission)
}

// DisableScanner is a paid mutator transaction binding the contract method 0x59238297.
//
// Solidity: function disableScanner(uint256 scannerId, uint8 permission) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) DisableScanner(scannerId *big.Int, permission uint8) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.DisableScanner(&_ScannerRegistry.TransactOpts, scannerId, permission)
}

// EnableScanner is a paid mutator transaction binding the contract method 0x4629f427.
//
// Solidity: function enableScanner(uint256 scannerId, uint8 permission) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) EnableScanner(opts *bind.TransactOpts, scannerId *big.Int, permission uint8) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "enableScanner", scannerId, permission)
}

// EnableScanner is a paid mutator transaction binding the contract method 0x4629f427.
//
// Solidity: function enableScanner(uint256 scannerId, uint8 permission) returns()
func (_ScannerRegistry *ScannerRegistrySession) EnableScanner(scannerId *big.Int, permission uint8) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.EnableScanner(&_ScannerRegistry.TransactOpts, scannerId, permission)
}

// EnableScanner is a paid mutator transaction binding the contract method 0x4629f427.
//
// Solidity: function enableScanner(uint256 scannerId, uint8 permission) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) EnableScanner(scannerId *big.Int, permission uint8) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.EnableScanner(&_ScannerRegistry.TransactOpts, scannerId, permission)
}

// Initialize is a paid mutator transaction binding the contract method 0x2016a0d2.
//
// Solidity: function initialize(address __manager, address __router, string __name, string __symbol) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) Initialize(opts *bind.TransactOpts, __manager common.Address, __router common.Address, __name string, __symbol string) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "initialize", __manager, __router, __name, __symbol)
}

// Initialize is a paid mutator transaction binding the contract method 0x2016a0d2.
//
// Solidity: function initialize(address __manager, address __router, string __name, string __symbol) returns()
func (_ScannerRegistry *ScannerRegistrySession) Initialize(__manager common.Address, __router common.Address, __name string, __symbol string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.Initialize(&_ScannerRegistry.TransactOpts, __manager, __router, __name, __symbol)
}

// Initialize is a paid mutator transaction binding the contract method 0x2016a0d2.
//
// Solidity: function initialize(address __manager, address __router, string __name, string __symbol) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) Initialize(__manager common.Address, __router common.Address, __name string, __symbol string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.Initialize(&_ScannerRegistry.TransactOpts, __manager, __router, __name, __symbol)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ScannerRegistry *ScannerRegistryTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ScannerRegistry *ScannerRegistrySession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.Multicall(&_ScannerRegistry.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ScannerRegistry *ScannerRegistryTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.Multicall(&_ScannerRegistry.TransactOpts, data)
}

// Register is a paid mutator transaction binding the contract method 0xf11b1b88.
//
// Solidity: function register(address owner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) Register(opts *bind.TransactOpts, owner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "register", owner, chainId, metadata)
}

// Register is a paid mutator transaction binding the contract method 0xf11b1b88.
//
// Solidity: function register(address owner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistrySession) Register(owner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.Register(&_ScannerRegistry.TransactOpts, owner, chainId, metadata)
}

// Register is a paid mutator transaction binding the contract method 0xf11b1b88.
//
// Solidity: function register(address owner, uint256 chainId, string metadata) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) Register(owner common.Address, chainId *big.Int, metadata string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.Register(&_ScannerRegistry.TransactOpts, owner, chainId, metadata)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ScannerRegistry *ScannerRegistrySession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SafeTransferFrom(&_ScannerRegistry.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SafeTransferFrom(&_ScannerRegistry.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, _data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_ScannerRegistry *ScannerRegistrySession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SafeTransferFrom0(&_ScannerRegistry.TransactOpts, from, to, tokenId, _data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SafeTransferFrom0(&_ScannerRegistry.TransactOpts, from, to, tokenId, _data)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) SetAccessManager(opts *bind.TransactOpts, newManager common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "setAccessManager", newManager)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_ScannerRegistry *ScannerRegistrySession) SetAccessManager(newManager common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetAccessManager(&_ScannerRegistry.TransactOpts, newManager)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) SetAccessManager(newManager common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetAccessManager(&_ScannerRegistry.TransactOpts, newManager)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ScannerRegistry *ScannerRegistrySession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetApprovalForAll(&_ScannerRegistry.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetApprovalForAll(&_ScannerRegistry.TransactOpts, operator, approved)
}

// SetManager is a paid mutator transaction binding the contract method 0x5a74fc29.
//
// Solidity: function setManager(uint256 scannerId, address manager, bool enable) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) SetManager(opts *bind.TransactOpts, scannerId *big.Int, manager common.Address, enable bool) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "setManager", scannerId, manager, enable)
}

// SetManager is a paid mutator transaction binding the contract method 0x5a74fc29.
//
// Solidity: function setManager(uint256 scannerId, address manager, bool enable) returns()
func (_ScannerRegistry *ScannerRegistrySession) SetManager(scannerId *big.Int, manager common.Address, enable bool) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetManager(&_ScannerRegistry.TransactOpts, scannerId, manager, enable)
}

// SetManager is a paid mutator transaction binding the contract method 0x5a74fc29.
//
// Solidity: function setManager(uint256 scannerId, address manager, bool enable) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) SetManager(scannerId *big.Int, manager common.Address, enable bool) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetManager(&_ScannerRegistry.TransactOpts, scannerId, manager, enable)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) SetName(opts *bind.TransactOpts, ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "setName", ensRegistry, ensName)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_ScannerRegistry *ScannerRegistrySession) SetName(ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetName(&_ScannerRegistry.TransactOpts, ensRegistry, ensName)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) SetName(ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetName(&_ScannerRegistry.TransactOpts, ensRegistry, ensName)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "setRouter", newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_ScannerRegistry *ScannerRegistrySession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetRouter(&_ScannerRegistry.TransactOpts, newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetRouter(&_ScannerRegistry.TransactOpts, newRouter)
}

// SetStakeController is a paid mutator transaction binding the contract method 0x9a9d67bd.
//
// Solidity: function setStakeController(address stakeController) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) SetStakeController(opts *bind.TransactOpts, stakeController common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "setStakeController", stakeController)
}

// SetStakeController is a paid mutator transaction binding the contract method 0x9a9d67bd.
//
// Solidity: function setStakeController(address stakeController) returns()
func (_ScannerRegistry *ScannerRegistrySession) SetStakeController(stakeController common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetStakeController(&_ScannerRegistry.TransactOpts, stakeController)
}

// SetStakeController is a paid mutator transaction binding the contract method 0x9a9d67bd.
//
// Solidity: function setStakeController(address stakeController) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) SetStakeController(stakeController common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.SetStakeController(&_ScannerRegistry.TransactOpts, stakeController)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ScannerRegistry *ScannerRegistrySession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.TransferFrom(&_ScannerRegistry.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.TransferFrom(&_ScannerRegistry.TransactOpts, from, to, tokenId)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ScannerRegistry *ScannerRegistryTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ScannerRegistry *ScannerRegistrySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.UpgradeTo(&_ScannerRegistry.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.UpgradeTo(&_ScannerRegistry.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ScannerRegistry *ScannerRegistryTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ScannerRegistry.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ScannerRegistry *ScannerRegistrySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.UpgradeToAndCall(&_ScannerRegistry.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ScannerRegistry *ScannerRegistryTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ScannerRegistry.Contract.UpgradeToAndCall(&_ScannerRegistry.TransactOpts, newImplementation, data)
}

// ScannerRegistryAccessManagerUpdatedIterator is returned from FilterAccessManagerUpdated and is used to iterate over the raw logs and unpacked data for AccessManagerUpdated events raised by the ScannerRegistry contract.
type ScannerRegistryAccessManagerUpdatedIterator struct {
	Event *ScannerRegistryAccessManagerUpdated // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryAccessManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryAccessManagerUpdated)
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
		it.Event = new(ScannerRegistryAccessManagerUpdated)
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
func (it *ScannerRegistryAccessManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryAccessManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryAccessManagerUpdated represents a AccessManagerUpdated event raised by the ScannerRegistry contract.
type ScannerRegistryAccessManagerUpdated struct {
	NewAddressManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterAccessManagerUpdated is a free log retrieval operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterAccessManagerUpdated(opts *bind.FilterOpts, newAddressManager []common.Address) (*ScannerRegistryAccessManagerUpdatedIterator, error) {

	var newAddressManagerRule []interface{}
	for _, newAddressManagerItem := range newAddressManager {
		newAddressManagerRule = append(newAddressManagerRule, newAddressManagerItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "AccessManagerUpdated", newAddressManagerRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryAccessManagerUpdatedIterator{contract: _ScannerRegistry.contract, event: "AccessManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchAccessManagerUpdated is a free log subscription operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchAccessManagerUpdated(opts *bind.WatchOpts, sink chan<- *ScannerRegistryAccessManagerUpdated, newAddressManager []common.Address) (event.Subscription, error) {

	var newAddressManagerRule []interface{}
	for _, newAddressManagerItem := range newAddressManager {
		newAddressManagerRule = append(newAddressManagerRule, newAddressManagerItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "AccessManagerUpdated", newAddressManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryAccessManagerUpdated)
				if err := _ScannerRegistry.contract.UnpackLog(event, "AccessManagerUpdated", log); err != nil {
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
func (_ScannerRegistry *ScannerRegistryFilterer) ParseAccessManagerUpdated(log types.Log) (*ScannerRegistryAccessManagerUpdated, error) {
	event := new(ScannerRegistryAccessManagerUpdated)
	if err := _ScannerRegistry.contract.UnpackLog(event, "AccessManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ScannerRegistry contract.
type ScannerRegistryAdminChangedIterator struct {
	Event *ScannerRegistryAdminChanged // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryAdminChanged)
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
		it.Event = new(ScannerRegistryAdminChanged)
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
func (it *ScannerRegistryAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryAdminChanged represents a AdminChanged event raised by the ScannerRegistry contract.
type ScannerRegistryAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ScannerRegistryAdminChangedIterator, error) {

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryAdminChangedIterator{contract: _ScannerRegistry.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ScannerRegistryAdminChanged) (event.Subscription, error) {

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryAdminChanged)
				if err := _ScannerRegistry.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_ScannerRegistry *ScannerRegistryFilterer) ParseAdminChanged(log types.Log) (*ScannerRegistryAdminChanged, error) {
	event := new(ScannerRegistryAdminChanged)
	if err := _ScannerRegistry.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ScannerRegistry contract.
type ScannerRegistryApprovalIterator struct {
	Event *ScannerRegistryApproval // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryApproval)
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
		it.Event = new(ScannerRegistryApproval)
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
func (it *ScannerRegistryApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryApproval represents a Approval event raised by the ScannerRegistry contract.
type ScannerRegistryApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*ScannerRegistryApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryApprovalIterator{contract: _ScannerRegistry.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ScannerRegistryApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryApproval)
				if err := _ScannerRegistry.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseApproval(log types.Log) (*ScannerRegistryApproval, error) {
	event := new(ScannerRegistryApproval)
	if err := _ScannerRegistry.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ScannerRegistry contract.
type ScannerRegistryApprovalForAllIterator struct {
	Event *ScannerRegistryApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryApprovalForAll)
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
		it.Event = new(ScannerRegistryApprovalForAll)
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
func (it *ScannerRegistryApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryApprovalForAll represents a ApprovalForAll event raised by the ScannerRegistry contract.
type ScannerRegistryApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ScannerRegistryApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryApprovalForAllIterator{contract: _ScannerRegistry.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ScannerRegistryApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryApprovalForAll)
				if err := _ScannerRegistry.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseApprovalForAll(log types.Log) (*ScannerRegistryApprovalForAll, error) {
	event := new(ScannerRegistryApprovalForAll)
	if err := _ScannerRegistry.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the ScannerRegistry contract.
type ScannerRegistryBeaconUpgradedIterator struct {
	Event *ScannerRegistryBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryBeaconUpgraded)
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
		it.Event = new(ScannerRegistryBeaconUpgraded)
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
func (it *ScannerRegistryBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryBeaconUpgraded represents a BeaconUpgraded event raised by the ScannerRegistry contract.
type ScannerRegistryBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ScannerRegistryBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryBeaconUpgradedIterator{contract: _ScannerRegistry.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ScannerRegistryBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryBeaconUpgraded)
				if err := _ScannerRegistry.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_ScannerRegistry *ScannerRegistryFilterer) ParseBeaconUpgraded(log types.Log) (*ScannerRegistryBeaconUpgraded, error) {
	event := new(ScannerRegistryBeaconUpgraded)
	if err := _ScannerRegistry.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryManagerEnabledIterator is returned from FilterManagerEnabled and is used to iterate over the raw logs and unpacked data for ManagerEnabled events raised by the ScannerRegistry contract.
type ScannerRegistryManagerEnabledIterator struct {
	Event *ScannerRegistryManagerEnabled // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryManagerEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryManagerEnabled)
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
		it.Event = new(ScannerRegistryManagerEnabled)
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
func (it *ScannerRegistryManagerEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryManagerEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryManagerEnabled represents a ManagerEnabled event raised by the ScannerRegistry contract.
type ScannerRegistryManagerEnabled struct {
	ScannerId *big.Int
	Manager   common.Address
	Enabled   bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterManagerEnabled is a free log retrieval operation binding the contract event 0x538b6537a6fe8f0deae9f3b86ad1924d5e5b3d5a683055276b2824f918be043e.
//
// Solidity: event ManagerEnabled(uint256 indexed scannerId, address indexed manager, bool enabled)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterManagerEnabled(opts *bind.FilterOpts, scannerId []*big.Int, manager []common.Address) (*ScannerRegistryManagerEnabledIterator, error) {

	var scannerIdRule []interface{}
	for _, scannerIdItem := range scannerId {
		scannerIdRule = append(scannerIdRule, scannerIdItem)
	}
	var managerRule []interface{}
	for _, managerItem := range manager {
		managerRule = append(managerRule, managerItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "ManagerEnabled", scannerIdRule, managerRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryManagerEnabledIterator{contract: _ScannerRegistry.contract, event: "ManagerEnabled", logs: logs, sub: sub}, nil
}

// WatchManagerEnabled is a free log subscription operation binding the contract event 0x538b6537a6fe8f0deae9f3b86ad1924d5e5b3d5a683055276b2824f918be043e.
//
// Solidity: event ManagerEnabled(uint256 indexed scannerId, address indexed manager, bool enabled)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchManagerEnabled(opts *bind.WatchOpts, sink chan<- *ScannerRegistryManagerEnabled, scannerId []*big.Int, manager []common.Address) (event.Subscription, error) {

	var scannerIdRule []interface{}
	for _, scannerIdItem := range scannerId {
		scannerIdRule = append(scannerIdRule, scannerIdItem)
	}
	var managerRule []interface{}
	for _, managerItem := range manager {
		managerRule = append(managerRule, managerItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "ManagerEnabled", scannerIdRule, managerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryManagerEnabled)
				if err := _ScannerRegistry.contract.UnpackLog(event, "ManagerEnabled", log); err != nil {
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

// ParseManagerEnabled is a log parse operation binding the contract event 0x538b6537a6fe8f0deae9f3b86ad1924d5e5b3d5a683055276b2824f918be043e.
//
// Solidity: event ManagerEnabled(uint256 indexed scannerId, address indexed manager, bool enabled)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseManagerEnabled(log types.Log) (*ScannerRegistryManagerEnabled, error) {
	event := new(ScannerRegistryManagerEnabled)
	if err := _ScannerRegistry.contract.UnpackLog(event, "ManagerEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryRouterUpdatedIterator is returned from FilterRouterUpdated and is used to iterate over the raw logs and unpacked data for RouterUpdated events raised by the ScannerRegistry contract.
type ScannerRegistryRouterUpdatedIterator struct {
	Event *ScannerRegistryRouterUpdated // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryRouterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryRouterUpdated)
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
		it.Event = new(ScannerRegistryRouterUpdated)
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
func (it *ScannerRegistryRouterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryRouterUpdated represents a RouterUpdated event raised by the ScannerRegistry contract.
type ScannerRegistryRouterUpdated struct {
	Router common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRouterUpdated is a free log retrieval operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterRouterUpdated(opts *bind.FilterOpts, router []common.Address) (*ScannerRegistryRouterUpdatedIterator, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "RouterUpdated", routerRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryRouterUpdatedIterator{contract: _ScannerRegistry.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

// WatchRouterUpdated is a free log subscription operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *ScannerRegistryRouterUpdated, router []common.Address) (event.Subscription, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "RouterUpdated", routerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryRouterUpdated)
				if err := _ScannerRegistry.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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
func (_ScannerRegistry *ScannerRegistryFilterer) ParseRouterUpdated(log types.Log) (*ScannerRegistryRouterUpdated, error) {
	event := new(ScannerRegistryRouterUpdated)
	if err := _ScannerRegistry.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryScannerEnabledIterator is returned from FilterScannerEnabled and is used to iterate over the raw logs and unpacked data for ScannerEnabled events raised by the ScannerRegistry contract.
type ScannerRegistryScannerEnabledIterator struct {
	Event *ScannerRegistryScannerEnabled // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryScannerEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryScannerEnabled)
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
		it.Event = new(ScannerRegistryScannerEnabled)
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
func (it *ScannerRegistryScannerEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryScannerEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryScannerEnabled represents a ScannerEnabled event raised by the ScannerRegistry contract.
type ScannerRegistryScannerEnabled struct {
	ScannerId  *big.Int
	Enabled    bool
	Permission uint8
	Value      bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterScannerEnabled is a free log retrieval operation binding the contract event 0xcde09e0ec4a155a87ef3eb8c971163d71fa1a87d4916cebef6ae4c9d296b25d4.
//
// Solidity: event ScannerEnabled(uint256 indexed scannerId, bool indexed enabled, uint8 permission, bool value)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterScannerEnabled(opts *bind.FilterOpts, scannerId []*big.Int, enabled []bool) (*ScannerRegistryScannerEnabledIterator, error) {

	var scannerIdRule []interface{}
	for _, scannerIdItem := range scannerId {
		scannerIdRule = append(scannerIdRule, scannerIdItem)
	}
	var enabledRule []interface{}
	for _, enabledItem := range enabled {
		enabledRule = append(enabledRule, enabledItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "ScannerEnabled", scannerIdRule, enabledRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryScannerEnabledIterator{contract: _ScannerRegistry.contract, event: "ScannerEnabled", logs: logs, sub: sub}, nil
}

// WatchScannerEnabled is a free log subscription operation binding the contract event 0xcde09e0ec4a155a87ef3eb8c971163d71fa1a87d4916cebef6ae4c9d296b25d4.
//
// Solidity: event ScannerEnabled(uint256 indexed scannerId, bool indexed enabled, uint8 permission, bool value)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchScannerEnabled(opts *bind.WatchOpts, sink chan<- *ScannerRegistryScannerEnabled, scannerId []*big.Int, enabled []bool) (event.Subscription, error) {

	var scannerIdRule []interface{}
	for _, scannerIdItem := range scannerId {
		scannerIdRule = append(scannerIdRule, scannerIdItem)
	}
	var enabledRule []interface{}
	for _, enabledItem := range enabled {
		enabledRule = append(enabledRule, enabledItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "ScannerEnabled", scannerIdRule, enabledRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryScannerEnabled)
				if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerEnabled", log); err != nil {
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

// ParseScannerEnabled is a log parse operation binding the contract event 0xcde09e0ec4a155a87ef3eb8c971163d71fa1a87d4916cebef6ae4c9d296b25d4.
//
// Solidity: event ScannerEnabled(uint256 indexed scannerId, bool indexed enabled, uint8 permission, bool value)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseScannerEnabled(log types.Log) (*ScannerRegistryScannerEnabled, error) {
	event := new(ScannerRegistryScannerEnabled)
	if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryScannerUpdatedIterator is returned from FilterScannerUpdated and is used to iterate over the raw logs and unpacked data for ScannerUpdated events raised by the ScannerRegistry contract.
type ScannerRegistryScannerUpdatedIterator struct {
	Event *ScannerRegistryScannerUpdated // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryScannerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryScannerUpdated)
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
		it.Event = new(ScannerRegistryScannerUpdated)
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
func (it *ScannerRegistryScannerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryScannerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryScannerUpdated represents a ScannerUpdated event raised by the ScannerRegistry contract.
type ScannerRegistryScannerUpdated struct {
	ScannerId *big.Int
	ChainId   *big.Int
	Metadata  string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterScannerUpdated is a free log retrieval operation binding the contract event 0x72d55569a8bd2d7bfb37627de4db16b8967136cfd50c423227036b24055e440d.
//
// Solidity: event ScannerUpdated(uint256 indexed scannerId, uint256 indexed chainId, string metadata)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterScannerUpdated(opts *bind.FilterOpts, scannerId []*big.Int, chainId []*big.Int) (*ScannerRegistryScannerUpdatedIterator, error) {

	var scannerIdRule []interface{}
	for _, scannerIdItem := range scannerId {
		scannerIdRule = append(scannerIdRule, scannerIdItem)
	}
	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "ScannerUpdated", scannerIdRule, chainIdRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryScannerUpdatedIterator{contract: _ScannerRegistry.contract, event: "ScannerUpdated", logs: logs, sub: sub}, nil
}

// WatchScannerUpdated is a free log subscription operation binding the contract event 0x72d55569a8bd2d7bfb37627de4db16b8967136cfd50c423227036b24055e440d.
//
// Solidity: event ScannerUpdated(uint256 indexed scannerId, uint256 indexed chainId, string metadata)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchScannerUpdated(opts *bind.WatchOpts, sink chan<- *ScannerRegistryScannerUpdated, scannerId []*big.Int, chainId []*big.Int) (event.Subscription, error) {

	var scannerIdRule []interface{}
	for _, scannerIdItem := range scannerId {
		scannerIdRule = append(scannerIdRule, scannerIdItem)
	}
	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "ScannerUpdated", scannerIdRule, chainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryScannerUpdated)
				if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerUpdated", log); err != nil {
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

// ParseScannerUpdated is a log parse operation binding the contract event 0x72d55569a8bd2d7bfb37627de4db16b8967136cfd50c423227036b24055e440d.
//
// Solidity: event ScannerUpdated(uint256 indexed scannerId, uint256 indexed chainId, string metadata)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseScannerUpdated(log types.Log) (*ScannerRegistryScannerUpdated, error) {
	event := new(ScannerRegistryScannerUpdated)
	if err := _ScannerRegistry.contract.UnpackLog(event, "ScannerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryStakeControllerUpdatedIterator is returned from FilterStakeControllerUpdated and is used to iterate over the raw logs and unpacked data for StakeControllerUpdated events raised by the ScannerRegistry contract.
type ScannerRegistryStakeControllerUpdatedIterator struct {
	Event *ScannerRegistryStakeControllerUpdated // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryStakeControllerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryStakeControllerUpdated)
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
		it.Event = new(ScannerRegistryStakeControllerUpdated)
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
func (it *ScannerRegistryStakeControllerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryStakeControllerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryStakeControllerUpdated represents a StakeControllerUpdated event raised by the ScannerRegistry contract.
type ScannerRegistryStakeControllerUpdated struct {
	NewstakeController common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterStakeControllerUpdated is a free log retrieval operation binding the contract event 0xcaa3d2f2b6f12475c0d16d986f57f334e0e8b9fff5335d3b6aafbca251da3f00.
//
// Solidity: event StakeControllerUpdated(address indexed newstakeController)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterStakeControllerUpdated(opts *bind.FilterOpts, newstakeController []common.Address) (*ScannerRegistryStakeControllerUpdatedIterator, error) {

	var newstakeControllerRule []interface{}
	for _, newstakeControllerItem := range newstakeController {
		newstakeControllerRule = append(newstakeControllerRule, newstakeControllerItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "StakeControllerUpdated", newstakeControllerRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryStakeControllerUpdatedIterator{contract: _ScannerRegistry.contract, event: "StakeControllerUpdated", logs: logs, sub: sub}, nil
}

// WatchStakeControllerUpdated is a free log subscription operation binding the contract event 0xcaa3d2f2b6f12475c0d16d986f57f334e0e8b9fff5335d3b6aafbca251da3f00.
//
// Solidity: event StakeControllerUpdated(address indexed newstakeController)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchStakeControllerUpdated(opts *bind.WatchOpts, sink chan<- *ScannerRegistryStakeControllerUpdated, newstakeController []common.Address) (event.Subscription, error) {

	var newstakeControllerRule []interface{}
	for _, newstakeControllerItem := range newstakeController {
		newstakeControllerRule = append(newstakeControllerRule, newstakeControllerItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "StakeControllerUpdated", newstakeControllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryStakeControllerUpdated)
				if err := _ScannerRegistry.contract.UnpackLog(event, "StakeControllerUpdated", log); err != nil {
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

// ParseStakeControllerUpdated is a log parse operation binding the contract event 0xcaa3d2f2b6f12475c0d16d986f57f334e0e8b9fff5335d3b6aafbca251da3f00.
//
// Solidity: event StakeControllerUpdated(address indexed newstakeController)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseStakeControllerUpdated(log types.Log) (*ScannerRegistryStakeControllerUpdated, error) {
	event := new(ScannerRegistryStakeControllerUpdated)
	if err := _ScannerRegistry.contract.UnpackLog(event, "StakeControllerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ScannerRegistry contract.
type ScannerRegistryTransferIterator struct {
	Event *ScannerRegistryTransfer // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryTransfer)
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
		it.Event = new(ScannerRegistryTransfer)
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
func (it *ScannerRegistryTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryTransfer represents a Transfer event raised by the ScannerRegistry contract.
type ScannerRegistryTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*ScannerRegistryTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryTransferIterator{contract: _ScannerRegistry.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ScannerRegistryTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryTransfer)
				if err := _ScannerRegistry.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ScannerRegistry *ScannerRegistryFilterer) ParseTransfer(log types.Log) (*ScannerRegistryTransfer, error) {
	event := new(ScannerRegistryTransfer)
	if err := _ScannerRegistry.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ScannerRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ScannerRegistry contract.
type ScannerRegistryUpgradedIterator struct {
	Event *ScannerRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *ScannerRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScannerRegistryUpgraded)
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
		it.Event = new(ScannerRegistryUpgraded)
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
func (it *ScannerRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ScannerRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ScannerRegistryUpgraded represents a Upgraded event raised by the ScannerRegistry contract.
type ScannerRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ScannerRegistry *ScannerRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ScannerRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ScannerRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ScannerRegistryUpgradedIterator{contract: _ScannerRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ScannerRegistry *ScannerRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ScannerRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ScannerRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ScannerRegistryUpgraded)
				if err := _ScannerRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ScannerRegistry *ScannerRegistryFilterer) ParseUpgraded(log types.Log) (*ScannerRegistryUpgraded, error) {
	event := new(ScannerRegistryUpgraded)
	if err := _ScannerRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
