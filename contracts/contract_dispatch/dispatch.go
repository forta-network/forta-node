// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract_dispatch

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

// DispatchMetaData contains all meta data concerning the Dispatch contract.
var DispatchMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MissingRole\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddressManager\",\"type\":\"address\"}],\"name\":\"AccessManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"Link\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"agentHash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"manifest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"agentRefAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"agentVersion\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"agentRegistry\",\"outputs\":[{\"internalType\":\"contractAgentRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"agentsAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"agentsFor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"__manager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"__router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"__agents\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"__scanners\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"link\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"scannerHash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"manifest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"scannerRefAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"scannerRegistry\",\"outputs\":[{\"internalType\":\"contractScannerRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"scannersAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"scannersFor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newManager\",\"type\":\"address\"}],\"name\":\"setAccessManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAgentRegistry\",\"type\":\"address\"}],\"name\":\"setAgentRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ensRegistry\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"ensName\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newScannerRegistry\",\"type\":\"address\"}],\"name\":\"setScannerRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"unlink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DispatchABI is the input ABI used to generate the binding from.
// Deprecated: Use DispatchMetaData.ABI instead.
var DispatchABI = DispatchMetaData.ABI

// Dispatch is an auto generated Go binding around an Ethereum contract.
type Dispatch struct {
	DispatchCaller     // Read-only binding to the contract
	DispatchTransactor // Write-only binding to the contract
	DispatchFilterer   // Log filterer for contract events
}

// DispatchCaller is an auto generated read-only Go binding around an Ethereum contract.
type DispatchCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DispatchTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DispatchTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DispatchFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DispatchFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DispatchSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DispatchSession struct {
	Contract     *Dispatch         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DispatchCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DispatchCallerSession struct {
	Contract *DispatchCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// DispatchTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DispatchTransactorSession struct {
	Contract     *DispatchTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DispatchRaw is an auto generated low-level Go binding around an Ethereum contract.
type DispatchRaw struct {
	Contract *Dispatch // Generic contract binding to access the raw methods on
}

// DispatchCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DispatchCallerRaw struct {
	Contract *DispatchCaller // Generic read-only contract binding to access the raw methods on
}

// DispatchTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DispatchTransactorRaw struct {
	Contract *DispatchTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDispatch creates a new instance of Dispatch, bound to a specific deployed contract.
func NewDispatch(address common.Address, backend bind.ContractBackend) (*Dispatch, error) {
	contract, err := bindDispatch(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Dispatch{DispatchCaller: DispatchCaller{contract: contract}, DispatchTransactor: DispatchTransactor{contract: contract}, DispatchFilterer: DispatchFilterer{contract: contract}}, nil
}

// NewDispatchCaller creates a new read-only instance of Dispatch, bound to a specific deployed contract.
func NewDispatchCaller(address common.Address, caller bind.ContractCaller) (*DispatchCaller, error) {
	contract, err := bindDispatch(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DispatchCaller{contract: contract}, nil
}

// NewDispatchTransactor creates a new write-only instance of Dispatch, bound to a specific deployed contract.
func NewDispatchTransactor(address common.Address, transactor bind.ContractTransactor) (*DispatchTransactor, error) {
	contract, err := bindDispatch(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DispatchTransactor{contract: contract}, nil
}

// NewDispatchFilterer creates a new log filterer instance of Dispatch, bound to a specific deployed contract.
func NewDispatchFilterer(address common.Address, filterer bind.ContractFilterer) (*DispatchFilterer, error) {
	contract, err := bindDispatch(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DispatchFilterer{contract: contract}, nil
}

// bindDispatch binds a generic wrapper to an already deployed contract.
func bindDispatch(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DispatchABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dispatch *DispatchRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dispatch.Contract.DispatchCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dispatch *DispatchRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dispatch.Contract.DispatchTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dispatch *DispatchRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dispatch.Contract.DispatchTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dispatch *DispatchCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dispatch.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dispatch *DispatchTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dispatch.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dispatch *DispatchTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dispatch.Contract.contract.Transact(opts, method, params...)
}

// AgentHash is a free data retrieval call binding the contract method 0xc2c2e46a.
//
// Solidity: function agentHash(uint256 agentId) view returns(uint256 length, bytes32 manifest)
func (_Dispatch *DispatchCaller) AgentHash(opts *bind.CallOpts, agentId *big.Int) (struct {
	Length   *big.Int
	Manifest [32]byte
}, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "agentHash", agentId)

	outstruct := new(struct {
		Length   *big.Int
		Manifest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Length = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Manifest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// AgentHash is a free data retrieval call binding the contract method 0xc2c2e46a.
//
// Solidity: function agentHash(uint256 agentId) view returns(uint256 length, bytes32 manifest)
func (_Dispatch *DispatchSession) AgentHash(agentId *big.Int) (struct {
	Length   *big.Int
	Manifest [32]byte
}, error) {
	return _Dispatch.Contract.AgentHash(&_Dispatch.CallOpts, agentId)
}

// AgentHash is a free data retrieval call binding the contract method 0xc2c2e46a.
//
// Solidity: function agentHash(uint256 agentId) view returns(uint256 length, bytes32 manifest)
func (_Dispatch *DispatchCallerSession) AgentHash(agentId *big.Int) (struct {
	Length   *big.Int
	Manifest [32]byte
}, error) {
	return _Dispatch.Contract.AgentHash(&_Dispatch.CallOpts, agentId)
}

// AgentRefAt is a free data retrieval call binding the contract method 0x32dee2f6.
//
// Solidity: function agentRefAt(uint256 scannerId, uint256 pos) view returns(uint256 agentId, bool enabled, uint256 agentVersion, string metadata, uint256[] chainIds)
func (_Dispatch *DispatchCaller) AgentRefAt(opts *bind.CallOpts, scannerId *big.Int, pos *big.Int) (struct {
	AgentId      *big.Int
	Enabled      bool
	AgentVersion *big.Int
	Metadata     string
	ChainIds     []*big.Int
}, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "agentRefAt", scannerId, pos)

	outstruct := new(struct {
		AgentId      *big.Int
		Enabled      bool
		AgentVersion *big.Int
		Metadata     string
		ChainIds     []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AgentId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Enabled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.AgentVersion = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Metadata = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.ChainIds = *abi.ConvertType(out[4], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// AgentRefAt is a free data retrieval call binding the contract method 0x32dee2f6.
//
// Solidity: function agentRefAt(uint256 scannerId, uint256 pos) view returns(uint256 agentId, bool enabled, uint256 agentVersion, string metadata, uint256[] chainIds)
func (_Dispatch *DispatchSession) AgentRefAt(scannerId *big.Int, pos *big.Int) (struct {
	AgentId      *big.Int
	Enabled      bool
	AgentVersion *big.Int
	Metadata     string
	ChainIds     []*big.Int
}, error) {
	return _Dispatch.Contract.AgentRefAt(&_Dispatch.CallOpts, scannerId, pos)
}

// AgentRefAt is a free data retrieval call binding the contract method 0x32dee2f6.
//
// Solidity: function agentRefAt(uint256 scannerId, uint256 pos) view returns(uint256 agentId, bool enabled, uint256 agentVersion, string metadata, uint256[] chainIds)
func (_Dispatch *DispatchCallerSession) AgentRefAt(scannerId *big.Int, pos *big.Int) (struct {
	AgentId      *big.Int
	Enabled      bool
	AgentVersion *big.Int
	Metadata     string
	ChainIds     []*big.Int
}, error) {
	return _Dispatch.Contract.AgentRefAt(&_Dispatch.CallOpts, scannerId, pos)
}

// AgentRegistry is a free data retrieval call binding the contract method 0x0d1cfcae.
//
// Solidity: function agentRegistry() view returns(address)
func (_Dispatch *DispatchCaller) AgentRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "agentRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentRegistry is a free data retrieval call binding the contract method 0x0d1cfcae.
//
// Solidity: function agentRegistry() view returns(address)
func (_Dispatch *DispatchSession) AgentRegistry() (common.Address, error) {
	return _Dispatch.Contract.AgentRegistry(&_Dispatch.CallOpts)
}

// AgentRegistry is a free data retrieval call binding the contract method 0x0d1cfcae.
//
// Solidity: function agentRegistry() view returns(address)
func (_Dispatch *DispatchCallerSession) AgentRegistry() (common.Address, error) {
	return _Dispatch.Contract.AgentRegistry(&_Dispatch.CallOpts)
}

// AgentsAt is a free data retrieval call binding the contract method 0xe323cba5.
//
// Solidity: function agentsAt(uint256 scannerId, uint256 pos) view returns(uint256)
func (_Dispatch *DispatchCaller) AgentsAt(opts *bind.CallOpts, scannerId *big.Int, pos *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "agentsAt", scannerId, pos)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentsAt is a free data retrieval call binding the contract method 0xe323cba5.
//
// Solidity: function agentsAt(uint256 scannerId, uint256 pos) view returns(uint256)
func (_Dispatch *DispatchSession) AgentsAt(scannerId *big.Int, pos *big.Int) (*big.Int, error) {
	return _Dispatch.Contract.AgentsAt(&_Dispatch.CallOpts, scannerId, pos)
}

// AgentsAt is a free data retrieval call binding the contract method 0xe323cba5.
//
// Solidity: function agentsAt(uint256 scannerId, uint256 pos) view returns(uint256)
func (_Dispatch *DispatchCallerSession) AgentsAt(scannerId *big.Int, pos *big.Int) (*big.Int, error) {
	return _Dispatch.Contract.AgentsAt(&_Dispatch.CallOpts, scannerId, pos)
}

// AgentsFor is a free data retrieval call binding the contract method 0xded838c7.
//
// Solidity: function agentsFor(uint256 scannerId) view returns(uint256)
func (_Dispatch *DispatchCaller) AgentsFor(opts *bind.CallOpts, scannerId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "agentsFor", scannerId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentsFor is a free data retrieval call binding the contract method 0xded838c7.
//
// Solidity: function agentsFor(uint256 scannerId) view returns(uint256)
func (_Dispatch *DispatchSession) AgentsFor(scannerId *big.Int) (*big.Int, error) {
	return _Dispatch.Contract.AgentsFor(&_Dispatch.CallOpts, scannerId)
}

// AgentsFor is a free data retrieval call binding the contract method 0xded838c7.
//
// Solidity: function agentsFor(uint256 scannerId) view returns(uint256)
func (_Dispatch *DispatchCallerSession) AgentsFor(scannerId *big.Int) (*big.Int, error) {
	return _Dispatch.Contract.AgentsFor(&_Dispatch.CallOpts, scannerId)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_Dispatch *DispatchCaller) IsTrustedForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "isTrustedForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_Dispatch *DispatchSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _Dispatch.Contract.IsTrustedForwarder(&_Dispatch.CallOpts, forwarder)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_Dispatch *DispatchCallerSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _Dispatch.Contract.IsTrustedForwarder(&_Dispatch.CallOpts, forwarder)
}

// ScannerHash is a free data retrieval call binding the contract method 0xb1774f9d.
//
// Solidity: function scannerHash(uint256 scannerId) view returns(uint256 length, bytes32 manifest)
func (_Dispatch *DispatchCaller) ScannerHash(opts *bind.CallOpts, scannerId *big.Int) (struct {
	Length   *big.Int
	Manifest [32]byte
}, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "scannerHash", scannerId)

	outstruct := new(struct {
		Length   *big.Int
		Manifest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Length = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Manifest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// ScannerHash is a free data retrieval call binding the contract method 0xb1774f9d.
//
// Solidity: function scannerHash(uint256 scannerId) view returns(uint256 length, bytes32 manifest)
func (_Dispatch *DispatchSession) ScannerHash(scannerId *big.Int) (struct {
	Length   *big.Int
	Manifest [32]byte
}, error) {
	return _Dispatch.Contract.ScannerHash(&_Dispatch.CallOpts, scannerId)
}

// ScannerHash is a free data retrieval call binding the contract method 0xb1774f9d.
//
// Solidity: function scannerHash(uint256 scannerId) view returns(uint256 length, bytes32 manifest)
func (_Dispatch *DispatchCallerSession) ScannerHash(scannerId *big.Int) (struct {
	Length   *big.Int
	Manifest [32]byte
}, error) {
	return _Dispatch.Contract.ScannerHash(&_Dispatch.CallOpts, scannerId)
}

// ScannerRefAt is a free data retrieval call binding the contract method 0x8b2e98d6.
//
// Solidity: function scannerRefAt(uint256 agentId, uint256 pos) view returns(uint256 scannerId, bool enabled)
func (_Dispatch *DispatchCaller) ScannerRefAt(opts *bind.CallOpts, agentId *big.Int, pos *big.Int) (struct {
	ScannerId *big.Int
	Enabled   bool
}, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "scannerRefAt", agentId, pos)

	outstruct := new(struct {
		ScannerId *big.Int
		Enabled   bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScannerId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Enabled = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// ScannerRefAt is a free data retrieval call binding the contract method 0x8b2e98d6.
//
// Solidity: function scannerRefAt(uint256 agentId, uint256 pos) view returns(uint256 scannerId, bool enabled)
func (_Dispatch *DispatchSession) ScannerRefAt(agentId *big.Int, pos *big.Int) (struct {
	ScannerId *big.Int
	Enabled   bool
}, error) {
	return _Dispatch.Contract.ScannerRefAt(&_Dispatch.CallOpts, agentId, pos)
}

// ScannerRefAt is a free data retrieval call binding the contract method 0x8b2e98d6.
//
// Solidity: function scannerRefAt(uint256 agentId, uint256 pos) view returns(uint256 scannerId, bool enabled)
func (_Dispatch *DispatchCallerSession) ScannerRefAt(agentId *big.Int, pos *big.Int) (struct {
	ScannerId *big.Int
	Enabled   bool
}, error) {
	return _Dispatch.Contract.ScannerRefAt(&_Dispatch.CallOpts, agentId, pos)
}

// ScannerRegistry is a free data retrieval call binding the contract method 0x5e9f88b1.
//
// Solidity: function scannerRegistry() view returns(address)
func (_Dispatch *DispatchCaller) ScannerRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "scannerRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ScannerRegistry is a free data retrieval call binding the contract method 0x5e9f88b1.
//
// Solidity: function scannerRegistry() view returns(address)
func (_Dispatch *DispatchSession) ScannerRegistry() (common.Address, error) {
	return _Dispatch.Contract.ScannerRegistry(&_Dispatch.CallOpts)
}

// ScannerRegistry is a free data retrieval call binding the contract method 0x5e9f88b1.
//
// Solidity: function scannerRegistry() view returns(address)
func (_Dispatch *DispatchCallerSession) ScannerRegistry() (common.Address, error) {
	return _Dispatch.Contract.ScannerRegistry(&_Dispatch.CallOpts)
}

// ScannersAt is a free data retrieval call binding the contract method 0x0a2bc370.
//
// Solidity: function scannersAt(uint256 agentId, uint256 pos) view returns(uint256)
func (_Dispatch *DispatchCaller) ScannersAt(opts *bind.CallOpts, agentId *big.Int, pos *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "scannersAt", agentId, pos)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ScannersAt is a free data retrieval call binding the contract method 0x0a2bc370.
//
// Solidity: function scannersAt(uint256 agentId, uint256 pos) view returns(uint256)
func (_Dispatch *DispatchSession) ScannersAt(agentId *big.Int, pos *big.Int) (*big.Int, error) {
	return _Dispatch.Contract.ScannersAt(&_Dispatch.CallOpts, agentId, pos)
}

// ScannersAt is a free data retrieval call binding the contract method 0x0a2bc370.
//
// Solidity: function scannersAt(uint256 agentId, uint256 pos) view returns(uint256)
func (_Dispatch *DispatchCallerSession) ScannersAt(agentId *big.Int, pos *big.Int) (*big.Int, error) {
	return _Dispatch.Contract.ScannersAt(&_Dispatch.CallOpts, agentId, pos)
}

// ScannersFor is a free data retrieval call binding the contract method 0x08eb7d4f.
//
// Solidity: function scannersFor(uint256 agentId) view returns(uint256)
func (_Dispatch *DispatchCaller) ScannersFor(opts *bind.CallOpts, agentId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "scannersFor", agentId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ScannersFor is a free data retrieval call binding the contract method 0x08eb7d4f.
//
// Solidity: function scannersFor(uint256 agentId) view returns(uint256)
func (_Dispatch *DispatchSession) ScannersFor(agentId *big.Int) (*big.Int, error) {
	return _Dispatch.Contract.ScannersFor(&_Dispatch.CallOpts, agentId)
}

// ScannersFor is a free data retrieval call binding the contract method 0x08eb7d4f.
//
// Solidity: function scannersFor(uint256 agentId) view returns(uint256)
func (_Dispatch *DispatchCallerSession) ScannersFor(agentId *big.Int) (*big.Int, error) {
	return _Dispatch.Contract.ScannersFor(&_Dispatch.CallOpts, agentId)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Dispatch *DispatchCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Dispatch.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Dispatch *DispatchSession) Version() (string, error) {
	return _Dispatch.Contract.Version(&_Dispatch.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Dispatch *DispatchCallerSession) Version() (string, error) {
	return _Dispatch.Contract.Version(&_Dispatch.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address __manager, address __router, address __agents, address __scanners) returns()
func (_Dispatch *DispatchTransactor) Initialize(opts *bind.TransactOpts, __manager common.Address, __router common.Address, __agents common.Address, __scanners common.Address) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "initialize", __manager, __router, __agents, __scanners)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address __manager, address __router, address __agents, address __scanners) returns()
func (_Dispatch *DispatchSession) Initialize(__manager common.Address, __router common.Address, __agents common.Address, __scanners common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.Initialize(&_Dispatch.TransactOpts, __manager, __router, __agents, __scanners)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address __manager, address __router, address __agents, address __scanners) returns()
func (_Dispatch *DispatchTransactorSession) Initialize(__manager common.Address, __router common.Address, __agents common.Address, __scanners common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.Initialize(&_Dispatch.TransactOpts, __manager, __router, __agents, __scanners)
}

// Link is a paid mutator transaction binding the contract method 0x86cf48e7.
//
// Solidity: function link(uint256 agentId, uint256 scannerId) returns()
func (_Dispatch *DispatchTransactor) Link(opts *bind.TransactOpts, agentId *big.Int, scannerId *big.Int) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "link", agentId, scannerId)
}

// Link is a paid mutator transaction binding the contract method 0x86cf48e7.
//
// Solidity: function link(uint256 agentId, uint256 scannerId) returns()
func (_Dispatch *DispatchSession) Link(agentId *big.Int, scannerId *big.Int) (*types.Transaction, error) {
	return _Dispatch.Contract.Link(&_Dispatch.TransactOpts, agentId, scannerId)
}

// Link is a paid mutator transaction binding the contract method 0x86cf48e7.
//
// Solidity: function link(uint256 agentId, uint256 scannerId) returns()
func (_Dispatch *DispatchTransactorSession) Link(agentId *big.Int, scannerId *big.Int) (*types.Transaction, error) {
	return _Dispatch.Contract.Link(&_Dispatch.TransactOpts, agentId, scannerId)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Dispatch *DispatchTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Dispatch *DispatchSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Dispatch.Contract.Multicall(&_Dispatch.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Dispatch *DispatchTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Dispatch.Contract.Multicall(&_Dispatch.TransactOpts, data)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_Dispatch *DispatchTransactor) SetAccessManager(opts *bind.TransactOpts, newManager common.Address) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "setAccessManager", newManager)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_Dispatch *DispatchSession) SetAccessManager(newManager common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.SetAccessManager(&_Dispatch.TransactOpts, newManager)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_Dispatch *DispatchTransactorSession) SetAccessManager(newManager common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.SetAccessManager(&_Dispatch.TransactOpts, newManager)
}

// SetAgentRegistry is a paid mutator transaction binding the contract method 0x28342ecf.
//
// Solidity: function setAgentRegistry(address newAgentRegistry) returns()
func (_Dispatch *DispatchTransactor) SetAgentRegistry(opts *bind.TransactOpts, newAgentRegistry common.Address) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "setAgentRegistry", newAgentRegistry)
}

// SetAgentRegistry is a paid mutator transaction binding the contract method 0x28342ecf.
//
// Solidity: function setAgentRegistry(address newAgentRegistry) returns()
func (_Dispatch *DispatchSession) SetAgentRegistry(newAgentRegistry common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.SetAgentRegistry(&_Dispatch.TransactOpts, newAgentRegistry)
}

// SetAgentRegistry is a paid mutator transaction binding the contract method 0x28342ecf.
//
// Solidity: function setAgentRegistry(address newAgentRegistry) returns()
func (_Dispatch *DispatchTransactorSession) SetAgentRegistry(newAgentRegistry common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.SetAgentRegistry(&_Dispatch.TransactOpts, newAgentRegistry)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_Dispatch *DispatchTransactor) SetName(opts *bind.TransactOpts, ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "setName", ensRegistry, ensName)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_Dispatch *DispatchSession) SetName(ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _Dispatch.Contract.SetName(&_Dispatch.TransactOpts, ensRegistry, ensName)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_Dispatch *DispatchTransactorSession) SetName(ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _Dispatch.Contract.SetName(&_Dispatch.TransactOpts, ensRegistry, ensName)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_Dispatch *DispatchTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "setRouter", newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_Dispatch *DispatchSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.SetRouter(&_Dispatch.TransactOpts, newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_Dispatch *DispatchTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.SetRouter(&_Dispatch.TransactOpts, newRouter)
}

// SetScannerRegistry is a paid mutator transaction binding the contract method 0x6b254492.
//
// Solidity: function setScannerRegistry(address newScannerRegistry) returns()
func (_Dispatch *DispatchTransactor) SetScannerRegistry(opts *bind.TransactOpts, newScannerRegistry common.Address) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "setScannerRegistry", newScannerRegistry)
}

// SetScannerRegistry is a paid mutator transaction binding the contract method 0x6b254492.
//
// Solidity: function setScannerRegistry(address newScannerRegistry) returns()
func (_Dispatch *DispatchSession) SetScannerRegistry(newScannerRegistry common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.SetScannerRegistry(&_Dispatch.TransactOpts, newScannerRegistry)
}

// SetScannerRegistry is a paid mutator transaction binding the contract method 0x6b254492.
//
// Solidity: function setScannerRegistry(address newScannerRegistry) returns()
func (_Dispatch *DispatchTransactorSession) SetScannerRegistry(newScannerRegistry common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.SetScannerRegistry(&_Dispatch.TransactOpts, newScannerRegistry)
}

// Unlink is a paid mutator transaction binding the contract method 0x0c65b39d.
//
// Solidity: function unlink(uint256 agentId, uint256 scannerId) returns()
func (_Dispatch *DispatchTransactor) Unlink(opts *bind.TransactOpts, agentId *big.Int, scannerId *big.Int) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "unlink", agentId, scannerId)
}

// Unlink is a paid mutator transaction binding the contract method 0x0c65b39d.
//
// Solidity: function unlink(uint256 agentId, uint256 scannerId) returns()
func (_Dispatch *DispatchSession) Unlink(agentId *big.Int, scannerId *big.Int) (*types.Transaction, error) {
	return _Dispatch.Contract.Unlink(&_Dispatch.TransactOpts, agentId, scannerId)
}

// Unlink is a paid mutator transaction binding the contract method 0x0c65b39d.
//
// Solidity: function unlink(uint256 agentId, uint256 scannerId) returns()
func (_Dispatch *DispatchTransactorSession) Unlink(agentId *big.Int, scannerId *big.Int) (*types.Transaction, error) {
	return _Dispatch.Contract.Unlink(&_Dispatch.TransactOpts, agentId, scannerId)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Dispatch *DispatchTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Dispatch *DispatchSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.UpgradeTo(&_Dispatch.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Dispatch *DispatchTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Dispatch.Contract.UpgradeTo(&_Dispatch.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Dispatch *DispatchTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Dispatch.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Dispatch *DispatchSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Dispatch.Contract.UpgradeToAndCall(&_Dispatch.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Dispatch *DispatchTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Dispatch.Contract.UpgradeToAndCall(&_Dispatch.TransactOpts, newImplementation, data)
}

// DispatchAccessManagerUpdatedIterator is returned from FilterAccessManagerUpdated and is used to iterate over the raw logs and unpacked data for AccessManagerUpdated events raised by the Dispatch contract.
type DispatchAccessManagerUpdatedIterator struct {
	Event *DispatchAccessManagerUpdated // Event containing the contract specifics and raw log

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
func (it *DispatchAccessManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DispatchAccessManagerUpdated)
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
		it.Event = new(DispatchAccessManagerUpdated)
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
func (it *DispatchAccessManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DispatchAccessManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DispatchAccessManagerUpdated represents a AccessManagerUpdated event raised by the Dispatch contract.
type DispatchAccessManagerUpdated struct {
	NewAddressManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterAccessManagerUpdated is a free log retrieval operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_Dispatch *DispatchFilterer) FilterAccessManagerUpdated(opts *bind.FilterOpts, newAddressManager []common.Address) (*DispatchAccessManagerUpdatedIterator, error) {

	var newAddressManagerRule []interface{}
	for _, newAddressManagerItem := range newAddressManager {
		newAddressManagerRule = append(newAddressManagerRule, newAddressManagerItem)
	}

	logs, sub, err := _Dispatch.contract.FilterLogs(opts, "AccessManagerUpdated", newAddressManagerRule)
	if err != nil {
		return nil, err
	}
	return &DispatchAccessManagerUpdatedIterator{contract: _Dispatch.contract, event: "AccessManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchAccessManagerUpdated is a free log subscription operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_Dispatch *DispatchFilterer) WatchAccessManagerUpdated(opts *bind.WatchOpts, sink chan<- *DispatchAccessManagerUpdated, newAddressManager []common.Address) (event.Subscription, error) {

	var newAddressManagerRule []interface{}
	for _, newAddressManagerItem := range newAddressManager {
		newAddressManagerRule = append(newAddressManagerRule, newAddressManagerItem)
	}

	logs, sub, err := _Dispatch.contract.WatchLogs(opts, "AccessManagerUpdated", newAddressManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DispatchAccessManagerUpdated)
				if err := _Dispatch.contract.UnpackLog(event, "AccessManagerUpdated", log); err != nil {
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
func (_Dispatch *DispatchFilterer) ParseAccessManagerUpdated(log types.Log) (*DispatchAccessManagerUpdated, error) {
	event := new(DispatchAccessManagerUpdated)
	if err := _Dispatch.contract.UnpackLog(event, "AccessManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DispatchAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the Dispatch contract.
type DispatchAdminChangedIterator struct {
	Event *DispatchAdminChanged // Event containing the contract specifics and raw log

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
func (it *DispatchAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DispatchAdminChanged)
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
		it.Event = new(DispatchAdminChanged)
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
func (it *DispatchAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DispatchAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DispatchAdminChanged represents a AdminChanged event raised by the Dispatch contract.
type DispatchAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Dispatch *DispatchFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*DispatchAdminChangedIterator, error) {

	logs, sub, err := _Dispatch.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &DispatchAdminChangedIterator{contract: _Dispatch.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Dispatch *DispatchFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *DispatchAdminChanged) (event.Subscription, error) {

	logs, sub, err := _Dispatch.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DispatchAdminChanged)
				if err := _Dispatch.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_Dispatch *DispatchFilterer) ParseAdminChanged(log types.Log) (*DispatchAdminChanged, error) {
	event := new(DispatchAdminChanged)
	if err := _Dispatch.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DispatchBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the Dispatch contract.
type DispatchBeaconUpgradedIterator struct {
	Event *DispatchBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *DispatchBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DispatchBeaconUpgraded)
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
		it.Event = new(DispatchBeaconUpgraded)
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
func (it *DispatchBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DispatchBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DispatchBeaconUpgraded represents a BeaconUpgraded event raised by the Dispatch contract.
type DispatchBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Dispatch *DispatchFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*DispatchBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Dispatch.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &DispatchBeaconUpgradedIterator{contract: _Dispatch.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Dispatch *DispatchFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *DispatchBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Dispatch.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DispatchBeaconUpgraded)
				if err := _Dispatch.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_Dispatch *DispatchFilterer) ParseBeaconUpgraded(log types.Log) (*DispatchBeaconUpgraded, error) {
	event := new(DispatchBeaconUpgraded)
	if err := _Dispatch.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DispatchLinkIterator is returned from FilterLink and is used to iterate over the raw logs and unpacked data for Link events raised by the Dispatch contract.
type DispatchLinkIterator struct {
	Event *DispatchLink // Event containing the contract specifics and raw log

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
func (it *DispatchLinkIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DispatchLink)
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
		it.Event = new(DispatchLink)
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
func (it *DispatchLinkIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DispatchLinkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DispatchLink represents a Link event raised by the Dispatch contract.
type DispatchLink struct {
	AgentId   *big.Int
	ScannerId *big.Int
	Enable    bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLink is a free log retrieval operation binding the contract event 0xf1b8cb2c3105270e747f9df25ec991871d6732bb7c7b86a088fe7d59c9272bbf.
//
// Solidity: event Link(uint256 agentId, uint256 scannerId, bool enable)
func (_Dispatch *DispatchFilterer) FilterLink(opts *bind.FilterOpts) (*DispatchLinkIterator, error) {

	logs, sub, err := _Dispatch.contract.FilterLogs(opts, "Link")
	if err != nil {
		return nil, err
	}
	return &DispatchLinkIterator{contract: _Dispatch.contract, event: "Link", logs: logs, sub: sub}, nil
}

// WatchLink is a free log subscription operation binding the contract event 0xf1b8cb2c3105270e747f9df25ec991871d6732bb7c7b86a088fe7d59c9272bbf.
//
// Solidity: event Link(uint256 agentId, uint256 scannerId, bool enable)
func (_Dispatch *DispatchFilterer) WatchLink(opts *bind.WatchOpts, sink chan<- *DispatchLink) (event.Subscription, error) {

	logs, sub, err := _Dispatch.contract.WatchLogs(opts, "Link")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DispatchLink)
				if err := _Dispatch.contract.UnpackLog(event, "Link", log); err != nil {
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

// ParseLink is a log parse operation binding the contract event 0xf1b8cb2c3105270e747f9df25ec991871d6732bb7c7b86a088fe7d59c9272bbf.
//
// Solidity: event Link(uint256 agentId, uint256 scannerId, bool enable)
func (_Dispatch *DispatchFilterer) ParseLink(log types.Log) (*DispatchLink, error) {
	event := new(DispatchLink)
	if err := _Dispatch.contract.UnpackLog(event, "Link", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DispatchRouterUpdatedIterator is returned from FilterRouterUpdated and is used to iterate over the raw logs and unpacked data for RouterUpdated events raised by the Dispatch contract.
type DispatchRouterUpdatedIterator struct {
	Event *DispatchRouterUpdated // Event containing the contract specifics and raw log

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
func (it *DispatchRouterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DispatchRouterUpdated)
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
		it.Event = new(DispatchRouterUpdated)
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
func (it *DispatchRouterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DispatchRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DispatchRouterUpdated represents a RouterUpdated event raised by the Dispatch contract.
type DispatchRouterUpdated struct {
	Router common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRouterUpdated is a free log retrieval operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_Dispatch *DispatchFilterer) FilterRouterUpdated(opts *bind.FilterOpts, router []common.Address) (*DispatchRouterUpdatedIterator, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _Dispatch.contract.FilterLogs(opts, "RouterUpdated", routerRule)
	if err != nil {
		return nil, err
	}
	return &DispatchRouterUpdatedIterator{contract: _Dispatch.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

// WatchRouterUpdated is a free log subscription operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_Dispatch *DispatchFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *DispatchRouterUpdated, router []common.Address) (event.Subscription, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _Dispatch.contract.WatchLogs(opts, "RouterUpdated", routerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DispatchRouterUpdated)
				if err := _Dispatch.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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
func (_Dispatch *DispatchFilterer) ParseRouterUpdated(log types.Log) (*DispatchRouterUpdated, error) {
	event := new(DispatchRouterUpdated)
	if err := _Dispatch.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DispatchUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Dispatch contract.
type DispatchUpgradedIterator struct {
	Event *DispatchUpgraded // Event containing the contract specifics and raw log

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
func (it *DispatchUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DispatchUpgraded)
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
		it.Event = new(DispatchUpgraded)
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
func (it *DispatchUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DispatchUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DispatchUpgraded represents a Upgraded event raised by the Dispatch contract.
type DispatchUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Dispatch *DispatchFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*DispatchUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Dispatch.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &DispatchUpgradedIterator{contract: _Dispatch.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Dispatch *DispatchFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *DispatchUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Dispatch.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DispatchUpgraded)
				if err := _Dispatch.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Dispatch *DispatchFilterer) ParseUpgraded(log types.Log) (*DispatchUpgraded, error) {
	event := new(DispatchUpgraded)
	if err := _Dispatch.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
