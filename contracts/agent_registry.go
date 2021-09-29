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

// AgentRegistryMetadataAgentMetadata is an auto generated low-level Go binding around an user-defined struct.
type AgentRegistryMetadataAgentMetadata struct {
	Version  *big.Int
	Metadata string
	ChainIds []*big.Int
}

// AgentRegistryMetaData contains all meta data concerning the AgentRegistry contract.
var AgentRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddressManager\",\"type\":\"address\"}],\"name\":\"AccessManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"commit\",\"type\":\"bytes32\"}],\"name\":\"AgentCommitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumAgentRegistryEnable.Permission\",\"name\":\"permission\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"AgentEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"}],\"name\":\"AgentUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"}],\"name\":\"createAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"enumAgentRegistryEnable.Permission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"disableAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"enumAgentRegistryEnable.Permission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"enableAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"getAgent\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"}],\"internalType\":\"structAgentRegistryMetadata.AgentMetadata\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getAgentByChainAndIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getAgentByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAgentCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"name\":\"getAgentCountByChain\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"__manager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"__router\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"__name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"__symbol\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commit\",\"type\":\"bytes32\"}],\"name\":\"prepareAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newManager\",\"type\":\"address\"}],\"name\":\"setAccessManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ensRegistry\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"ensName\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"}],\"name\":\"updateAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
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

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistryCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistrySession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgentRegistry.Contract.DEFAULTADMINROLE(&_AgentRegistry.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistryCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgentRegistry.Contract.DEFAULTADMINROLE(&_AgentRegistry.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _AgentRegistry.Contract.BalanceOf(&_AgentRegistry.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _AgentRegistry.Contract.BalanceOf(&_AgentRegistry.CallOpts, owner)
}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((uint256,string,uint256[]))
func (_AgentRegistry *AgentRegistryCaller) GetAgent(opts *bind.CallOpts, agentId *big.Int) (AgentRegistryMetadataAgentMetadata, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "getAgent", agentId)

	if err != nil {
		return *new(AgentRegistryMetadataAgentMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new(AgentRegistryMetadataAgentMetadata)).(*AgentRegistryMetadataAgentMetadata)

	return out0, err

}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((uint256,string,uint256[]))
func (_AgentRegistry *AgentRegistrySession) GetAgent(agentId *big.Int) (AgentRegistryMetadataAgentMetadata, error) {
	return _AgentRegistry.Contract.GetAgent(&_AgentRegistry.CallOpts, agentId)
}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((uint256,string,uint256[]))
func (_AgentRegistry *AgentRegistryCallerSession) GetAgent(agentId *big.Int) (AgentRegistryMetadataAgentMetadata, error) {
	return _AgentRegistry.Contract.GetAgent(&_AgentRegistry.CallOpts, agentId)
}

// GetAgentByChainAndIndex is a free data retrieval call binding the contract method 0xac388ff5.
//
// Solidity: function getAgentByChainAndIndex(uint256 chainId, uint256 index) view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) GetAgentByChainAndIndex(opts *bind.CallOpts, chainId *big.Int, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "getAgentByChainAndIndex", chainId, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAgentByChainAndIndex is a free data retrieval call binding the contract method 0xac388ff5.
//
// Solidity: function getAgentByChainAndIndex(uint256 chainId, uint256 index) view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) GetAgentByChainAndIndex(chainId *big.Int, index *big.Int) (*big.Int, error) {
	return _AgentRegistry.Contract.GetAgentByChainAndIndex(&_AgentRegistry.CallOpts, chainId, index)
}

// GetAgentByChainAndIndex is a free data retrieval call binding the contract method 0xac388ff5.
//
// Solidity: function getAgentByChainAndIndex(uint256 chainId, uint256 index) view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) GetAgentByChainAndIndex(chainId *big.Int, index *big.Int) (*big.Int, error) {
	return _AgentRegistry.Contract.GetAgentByChainAndIndex(&_AgentRegistry.CallOpts, chainId, index)
}

// GetAgentByIndex is a free data retrieval call binding the contract method 0x89432d40.
//
// Solidity: function getAgentByIndex(uint256 index) view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) GetAgentByIndex(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "getAgentByIndex", index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAgentByIndex is a free data retrieval call binding the contract method 0x89432d40.
//
// Solidity: function getAgentByIndex(uint256 index) view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) GetAgentByIndex(index *big.Int) (*big.Int, error) {
	return _AgentRegistry.Contract.GetAgentByIndex(&_AgentRegistry.CallOpts, index)
}

// GetAgentByIndex is a free data retrieval call binding the contract method 0x89432d40.
//
// Solidity: function getAgentByIndex(uint256 index) view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) GetAgentByIndex(index *big.Int) (*big.Int, error) {
	return _AgentRegistry.Contract.GetAgentByIndex(&_AgentRegistry.CallOpts, index)
}

// GetAgentCount is a free data retrieval call binding the contract method 0x91cab63e.
//
// Solidity: function getAgentCount() view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) GetAgentCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "getAgentCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAgentCount is a free data retrieval call binding the contract method 0x91cab63e.
//
// Solidity: function getAgentCount() view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) GetAgentCount() (*big.Int, error) {
	return _AgentRegistry.Contract.GetAgentCount(&_AgentRegistry.CallOpts)
}

// GetAgentCount is a free data retrieval call binding the contract method 0x91cab63e.
//
// Solidity: function getAgentCount() view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) GetAgentCount() (*big.Int, error) {
	return _AgentRegistry.Contract.GetAgentCount(&_AgentRegistry.CallOpts)
}

// GetAgentCountByChain is a free data retrieval call binding the contract method 0x8eea276f.
//
// Solidity: function getAgentCountByChain(uint256 chainId) view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) GetAgentCountByChain(opts *bind.CallOpts, chainId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "getAgentCountByChain", chainId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAgentCountByChain is a free data retrieval call binding the contract method 0x8eea276f.
//
// Solidity: function getAgentCountByChain(uint256 chainId) view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) GetAgentCountByChain(chainId *big.Int) (*big.Int, error) {
	return _AgentRegistry.Contract.GetAgentCountByChain(&_AgentRegistry.CallOpts, chainId)
}

// GetAgentCountByChain is a free data retrieval call binding the contract method 0x8eea276f.
//
// Solidity: function getAgentCountByChain(uint256 chainId) view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) GetAgentCountByChain(chainId *big.Int) (*big.Int, error) {
	return _AgentRegistry.Contract.GetAgentCountByChain(&_AgentRegistry.CallOpts, chainId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_AgentRegistry *AgentRegistryCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_AgentRegistry *AgentRegistrySession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _AgentRegistry.Contract.GetApproved(&_AgentRegistry.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_AgentRegistry *AgentRegistryCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _AgentRegistry.Contract.GetApproved(&_AgentRegistry.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _AgentRegistry.Contract.IsApprovedForAll(&_AgentRegistry.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _AgentRegistry.Contract.IsApprovedForAll(&_AgentRegistry.CallOpts, owner, operator)
}

// IsEnabled is a free data retrieval call binding the contract method 0xc783034c.
//
// Solidity: function isEnabled(uint256 agentId) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) IsEnabled(opts *bind.CallOpts, agentId *big.Int) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "isEnabled", agentId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEnabled is a free data retrieval call binding the contract method 0xc783034c.
//
// Solidity: function isEnabled(uint256 agentId) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) IsEnabled(agentId *big.Int) (bool, error) {
	return _AgentRegistry.Contract.IsEnabled(&_AgentRegistry.CallOpts, agentId)
}

// IsEnabled is a free data retrieval call binding the contract method 0xc783034c.
//
// Solidity: function isEnabled(uint256 agentId) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) IsEnabled(agentId *big.Int) (bool, error) {
	return _AgentRegistry.Contract.IsEnabled(&_AgentRegistry.CallOpts, agentId)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AgentRegistry *AgentRegistryCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AgentRegistry *AgentRegistrySession) Name() (string, error) {
	return _AgentRegistry.Contract.Name(&_AgentRegistry.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AgentRegistry *AgentRegistryCallerSession) Name() (string, error) {
	return _AgentRegistry.Contract.Name(&_AgentRegistry.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_AgentRegistry *AgentRegistryCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_AgentRegistry *AgentRegistrySession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _AgentRegistry.Contract.OwnerOf(&_AgentRegistry.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_AgentRegistry *AgentRegistryCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _AgentRegistry.Contract.OwnerOf(&_AgentRegistry.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgentRegistry.Contract.SupportsInterface(&_AgentRegistry.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgentRegistry.Contract.SupportsInterface(&_AgentRegistry.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AgentRegistry *AgentRegistryCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AgentRegistry *AgentRegistrySession) Symbol() (string, error) {
	return _AgentRegistry.Contract.Symbol(&_AgentRegistry.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AgentRegistry *AgentRegistryCallerSession) Symbol() (string, error) {
	return _AgentRegistry.Contract.Symbol(&_AgentRegistry.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_AgentRegistry *AgentRegistryCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_AgentRegistry *AgentRegistrySession) TokenURI(tokenId *big.Int) (string, error) {
	return _AgentRegistry.Contract.TokenURI(&_AgentRegistry.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_AgentRegistry *AgentRegistryCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _AgentRegistry.Contract.TokenURI(&_AgentRegistry.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_AgentRegistry *AgentRegistryTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_AgentRegistry *AgentRegistrySession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.Approve(&_AgentRegistry.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.Approve(&_AgentRegistry.TransactOpts, to, tokenId)
}

// CreateAgent is a paid mutator transaction binding the contract method 0x7935d5b4.
//
// Solidity: function createAgent(uint256 agentId, address owner, string metadata, uint256[] chainIds) returns()
func (_AgentRegistry *AgentRegistryTransactor) CreateAgent(opts *bind.TransactOpts, agentId *big.Int, owner common.Address, metadata string, chainIds []*big.Int) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "createAgent", agentId, owner, metadata, chainIds)
}

// CreateAgent is a paid mutator transaction binding the contract method 0x7935d5b4.
//
// Solidity: function createAgent(uint256 agentId, address owner, string metadata, uint256[] chainIds) returns()
func (_AgentRegistry *AgentRegistrySession) CreateAgent(agentId *big.Int, owner common.Address, metadata string, chainIds []*big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.CreateAgent(&_AgentRegistry.TransactOpts, agentId, owner, metadata, chainIds)
}

// CreateAgent is a paid mutator transaction binding the contract method 0x7935d5b4.
//
// Solidity: function createAgent(uint256 agentId, address owner, string metadata, uint256[] chainIds) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) CreateAgent(agentId *big.Int, owner common.Address, metadata string, chainIds []*big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.CreateAgent(&_AgentRegistry.TransactOpts, agentId, owner, metadata, chainIds)
}

// DisableAgent is a paid mutator transaction binding the contract method 0x21095d65.
//
// Solidity: function disableAgent(uint256 agentId, uint8 permission) returns()
func (_AgentRegistry *AgentRegistryTransactor) DisableAgent(opts *bind.TransactOpts, agentId *big.Int, permission uint8) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "disableAgent", agentId, permission)
}

// DisableAgent is a paid mutator transaction binding the contract method 0x21095d65.
//
// Solidity: function disableAgent(uint256 agentId, uint8 permission) returns()
func (_AgentRegistry *AgentRegistrySession) DisableAgent(agentId *big.Int, permission uint8) (*types.Transaction, error) {
	return _AgentRegistry.Contract.DisableAgent(&_AgentRegistry.TransactOpts, agentId, permission)
}

// DisableAgent is a paid mutator transaction binding the contract method 0x21095d65.
//
// Solidity: function disableAgent(uint256 agentId, uint8 permission) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) DisableAgent(agentId *big.Int, permission uint8) (*types.Transaction, error) {
	return _AgentRegistry.Contract.DisableAgent(&_AgentRegistry.TransactOpts, agentId, permission)
}

// EnableAgent is a paid mutator transaction binding the contract method 0xd0c6464b.
//
// Solidity: function enableAgent(uint256 agentId, uint8 permission) returns()
func (_AgentRegistry *AgentRegistryTransactor) EnableAgent(opts *bind.TransactOpts, agentId *big.Int, permission uint8) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "enableAgent", agentId, permission)
}

// EnableAgent is a paid mutator transaction binding the contract method 0xd0c6464b.
//
// Solidity: function enableAgent(uint256 agentId, uint8 permission) returns()
func (_AgentRegistry *AgentRegistrySession) EnableAgent(agentId *big.Int, permission uint8) (*types.Transaction, error) {
	return _AgentRegistry.Contract.EnableAgent(&_AgentRegistry.TransactOpts, agentId, permission)
}

// EnableAgent is a paid mutator transaction binding the contract method 0xd0c6464b.
//
// Solidity: function enableAgent(uint256 agentId, uint8 permission) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) EnableAgent(agentId *big.Int, permission uint8) (*types.Transaction, error) {
	return _AgentRegistry.Contract.EnableAgent(&_AgentRegistry.TransactOpts, agentId, permission)
}

// Initialize is a paid mutator transaction binding the contract method 0x2016a0d2.
//
// Solidity: function initialize(address __manager, address __router, string __name, string __symbol) returns()
func (_AgentRegistry *AgentRegistryTransactor) Initialize(opts *bind.TransactOpts, __manager common.Address, __router common.Address, __name string, __symbol string) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "initialize", __manager, __router, __name, __symbol)
}

// Initialize is a paid mutator transaction binding the contract method 0x2016a0d2.
//
// Solidity: function initialize(address __manager, address __router, string __name, string __symbol) returns()
func (_AgentRegistry *AgentRegistrySession) Initialize(__manager common.Address, __router common.Address, __name string, __symbol string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.Initialize(&_AgentRegistry.TransactOpts, __manager, __router, __name, __symbol)
}

// Initialize is a paid mutator transaction binding the contract method 0x2016a0d2.
//
// Solidity: function initialize(address __manager, address __router, string __name, string __symbol) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) Initialize(__manager common.Address, __router common.Address, __name string, __symbol string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.Initialize(&_AgentRegistry.TransactOpts, __manager, __router, __name, __symbol)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_AgentRegistry *AgentRegistryTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_AgentRegistry *AgentRegistrySession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.Multicall(&_AgentRegistry.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_AgentRegistry *AgentRegistryTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.Multicall(&_AgentRegistry.TransactOpts, data)
}

// PrepareAgent is a paid mutator transaction binding the contract method 0xa8439d98.
//
// Solidity: function prepareAgent(bytes32 commit) returns()
func (_AgentRegistry *AgentRegistryTransactor) PrepareAgent(opts *bind.TransactOpts, commit [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "prepareAgent", commit)
}

// PrepareAgent is a paid mutator transaction binding the contract method 0xa8439d98.
//
// Solidity: function prepareAgent(bytes32 commit) returns()
func (_AgentRegistry *AgentRegistrySession) PrepareAgent(commit [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.PrepareAgent(&_AgentRegistry.TransactOpts, commit)
}

// PrepareAgent is a paid mutator transaction binding the contract method 0xa8439d98.
//
// Solidity: function prepareAgent(bytes32 commit) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) PrepareAgent(commit [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.PrepareAgent(&_AgentRegistry.TransactOpts, commit)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_AgentRegistry *AgentRegistryTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_AgentRegistry *AgentRegistrySession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SafeTransferFrom(&_AgentRegistry.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SafeTransferFrom(&_AgentRegistry.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_AgentRegistry *AgentRegistryTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, _data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_AgentRegistry *AgentRegistrySession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SafeTransferFrom0(&_AgentRegistry.TransactOpts, from, to, tokenId, _data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SafeTransferFrom0(&_AgentRegistry.TransactOpts, from, to, tokenId, _data)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_AgentRegistry *AgentRegistryTransactor) SetAccessManager(opts *bind.TransactOpts, newManager common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "setAccessManager", newManager)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_AgentRegistry *AgentRegistrySession) SetAccessManager(newManager common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetAccessManager(&_AgentRegistry.TransactOpts, newManager)
}

// SetAccessManager is a paid mutator transaction binding the contract method 0xc9580804.
//
// Solidity: function setAccessManager(address newManager) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) SetAccessManager(newManager common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetAccessManager(&_AgentRegistry.TransactOpts, newManager)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_AgentRegistry *AgentRegistryTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_AgentRegistry *AgentRegistrySession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetApprovalForAll(&_AgentRegistry.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetApprovalForAll(&_AgentRegistry.TransactOpts, operator, approved)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_AgentRegistry *AgentRegistryTransactor) SetName(opts *bind.TransactOpts, ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "setName", ensRegistry, ensName)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_AgentRegistry *AgentRegistrySession) SetName(ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetName(&_AgentRegistry.TransactOpts, ensRegistry, ensName)
}

// SetName is a paid mutator transaction binding the contract method 0x3121db1c.
//
// Solidity: function setName(address ensRegistry, string ensName) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) SetName(ensRegistry common.Address, ensName string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetName(&_AgentRegistry.TransactOpts, ensRegistry, ensName)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_AgentRegistry *AgentRegistryTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "setRouter", newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_AgentRegistry *AgentRegistrySession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetRouter(&_AgentRegistry.TransactOpts, newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetRouter(&_AgentRegistry.TransactOpts, newRouter)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_AgentRegistry *AgentRegistryTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_AgentRegistry *AgentRegistrySession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.TransferFrom(&_AgentRegistry.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.TransferFrom(&_AgentRegistry.TransactOpts, from, to, tokenId)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xaa9ac6c6.
//
// Solidity: function updateAgent(uint256 agentId, string metadata, uint256[] chainIds) returns()
func (_AgentRegistry *AgentRegistryTransactor) UpdateAgent(opts *bind.TransactOpts, agentId *big.Int, metadata string, chainIds []*big.Int) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "updateAgent", agentId, metadata, chainIds)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xaa9ac6c6.
//
// Solidity: function updateAgent(uint256 agentId, string metadata, uint256[] chainIds) returns()
func (_AgentRegistry *AgentRegistrySession) UpdateAgent(agentId *big.Int, metadata string, chainIds []*big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpdateAgent(&_AgentRegistry.TransactOpts, agentId, metadata, chainIds)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xaa9ac6c6.
//
// Solidity: function updateAgent(uint256 agentId, string metadata, uint256[] chainIds) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) UpdateAgent(agentId *big.Int, metadata string, chainIds []*big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpdateAgent(&_AgentRegistry.TransactOpts, agentId, metadata, chainIds)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_AgentRegistry *AgentRegistryTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_AgentRegistry *AgentRegistrySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpgradeTo(&_AgentRegistry.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpgradeTo(&_AgentRegistry.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AgentRegistry *AgentRegistryTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AgentRegistry *AgentRegistrySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpgradeToAndCall(&_AgentRegistry.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AgentRegistry *AgentRegistryTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpgradeToAndCall(&_AgentRegistry.TransactOpts, newImplementation, data)
}

// AgentRegistryAccessManagerUpdatedIterator is returned from FilterAccessManagerUpdated and is used to iterate over the raw logs and unpacked data for AccessManagerUpdated events raised by the AgentRegistry contract.
type AgentRegistryAccessManagerUpdatedIterator struct {
	Event *AgentRegistryAccessManagerUpdated // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAccessManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAccessManagerUpdated)
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
		it.Event = new(AgentRegistryAccessManagerUpdated)
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
func (it *AgentRegistryAccessManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAccessManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAccessManagerUpdated represents a AccessManagerUpdated event raised by the AgentRegistry contract.
type AgentRegistryAccessManagerUpdated struct {
	NewAddressManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterAccessManagerUpdated is a free log retrieval operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_AgentRegistry *AgentRegistryFilterer) FilterAccessManagerUpdated(opts *bind.FilterOpts, newAddressManager []common.Address) (*AgentRegistryAccessManagerUpdatedIterator, error) {

	var newAddressManagerRule []interface{}
	for _, newAddressManagerItem := range newAddressManager {
		newAddressManagerRule = append(newAddressManagerRule, newAddressManagerItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AccessManagerUpdated", newAddressManagerRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAccessManagerUpdatedIterator{contract: _AgentRegistry.contract, event: "AccessManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchAccessManagerUpdated is a free log subscription operation binding the contract event 0xa5bc17e575e3b53b23d0e93e121a5a66d1de4d5edb4dfde6027b14d79b7f2b9c.
//
// Solidity: event AccessManagerUpdated(address indexed newAddressManager)
func (_AgentRegistry *AgentRegistryFilterer) WatchAccessManagerUpdated(opts *bind.WatchOpts, sink chan<- *AgentRegistryAccessManagerUpdated, newAddressManager []common.Address) (event.Subscription, error) {

	var newAddressManagerRule []interface{}
	for _, newAddressManagerItem := range newAddressManager {
		newAddressManagerRule = append(newAddressManagerRule, newAddressManagerItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AccessManagerUpdated", newAddressManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAccessManagerUpdated)
				if err := _AgentRegistry.contract.UnpackLog(event, "AccessManagerUpdated", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseAccessManagerUpdated(log types.Log) (*AgentRegistryAccessManagerUpdated, error) {
	event := new(AgentRegistryAccessManagerUpdated)
	if err := _AgentRegistry.contract.UnpackLog(event, "AccessManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the AgentRegistry contract.
type AgentRegistryAdminChangedIterator struct {
	Event *AgentRegistryAdminChanged // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAdminChanged)
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
		it.Event = new(AgentRegistryAdminChanged)
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
func (it *AgentRegistryAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAdminChanged represents a AdminChanged event raised by the AgentRegistry contract.
type AgentRegistryAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_AgentRegistry *AgentRegistryFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*AgentRegistryAdminChangedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAdminChangedIterator{contract: _AgentRegistry.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_AgentRegistry *AgentRegistryFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *AgentRegistryAdminChanged) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAdminChanged)
				if err := _AgentRegistry.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseAdminChanged(log types.Log) (*AgentRegistryAdminChanged, error) {
	event := new(AgentRegistryAdminChanged)
	if err := _AgentRegistry.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentCommittedIterator is returned from FilterAgentCommitted and is used to iterate over the raw logs and unpacked data for AgentCommitted events raised by the AgentRegistry contract.
type AgentRegistryAgentCommittedIterator struct {
	Event *AgentRegistryAgentCommitted // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentCommittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentCommitted)
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
		it.Event = new(AgentRegistryAgentCommitted)
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
func (it *AgentRegistryAgentCommittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentCommittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentCommitted represents a AgentCommitted event raised by the AgentRegistry contract.
type AgentRegistryAgentCommitted struct {
	Commit [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAgentCommitted is a free log retrieval operation binding the contract event 0x43343771105cc71439fde05af6660ad660519c971ddc249e5d5cd9c29da46d38.
//
// Solidity: event AgentCommitted(bytes32 indexed commit)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentCommitted(opts *bind.FilterOpts, commit [][32]byte) (*AgentRegistryAgentCommittedIterator, error) {

	var commitRule []interface{}
	for _, commitItem := range commit {
		commitRule = append(commitRule, commitItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentCommitted", commitRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentCommittedIterator{contract: _AgentRegistry.contract, event: "AgentCommitted", logs: logs, sub: sub}, nil
}

// WatchAgentCommitted is a free log subscription operation binding the contract event 0x43343771105cc71439fde05af6660ad660519c971ddc249e5d5cd9c29da46d38.
//
// Solidity: event AgentCommitted(bytes32 indexed commit)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentCommitted(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentCommitted, commit [][32]byte) (event.Subscription, error) {

	var commitRule []interface{}
	for _, commitItem := range commit {
		commitRule = append(commitRule, commitItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentCommitted", commitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentCommitted)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentCommitted", log); err != nil {
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

// ParseAgentCommitted is a log parse operation binding the contract event 0x43343771105cc71439fde05af6660ad660519c971ddc249e5d5cd9c29da46d38.
//
// Solidity: event AgentCommitted(bytes32 indexed commit)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentCommitted(log types.Log) (*AgentRegistryAgentCommitted, error) {
	event := new(AgentRegistryAgentCommitted)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentCommitted", log); err != nil {
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
	AgentId    *big.Int
	Permission uint8
	Enabled    bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAgentEnabled is a free log retrieval operation binding the contract event 0x9aae0324e55c8e6801ea6ce0a07d77640a94e8a03658c6b34aa55d48a51ec4ef.
//
// Solidity: event AgentEnabled(uint256 indexed agentId, uint8 permission, bool enabled)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentEnabled(opts *bind.FilterOpts, agentId []*big.Int) (*AgentRegistryAgentEnabledIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentEnabled", agentIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentEnabledIterator{contract: _AgentRegistry.contract, event: "AgentEnabled", logs: logs, sub: sub}, nil
}

// WatchAgentEnabled is a free log subscription operation binding the contract event 0x9aae0324e55c8e6801ea6ce0a07d77640a94e8a03658c6b34aa55d48a51ec4ef.
//
// Solidity: event AgentEnabled(uint256 indexed agentId, uint8 permission, bool enabled)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentEnabled(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentEnabled, agentId []*big.Int) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentEnabled", agentIdRule)
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

// ParseAgentEnabled is a log parse operation binding the contract event 0x9aae0324e55c8e6801ea6ce0a07d77640a94e8a03658c6b34aa55d48a51ec4ef.
//
// Solidity: event AgentEnabled(uint256 indexed agentId, uint8 permission, bool enabled)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentEnabled(log types.Log) (*AgentRegistryAgentEnabled, error) {
	event := new(AgentRegistryAgentEnabled)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentEnabled", log); err != nil {
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
	AgentId  *big.Int
	Metadata string
	ChainIds []*big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAgentUpdated is a free log retrieval operation binding the contract event 0x808318fff0e666157af89dad699f694ddee6735f06dba7769a172a8251d34b3f.
//
// Solidity: event AgentUpdated(uint256 indexed agentId, string metadata, uint256[] chainIds)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentUpdated(opts *bind.FilterOpts, agentId []*big.Int) (*AgentRegistryAgentUpdatedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentUpdated", agentIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentUpdatedIterator{contract: _AgentRegistry.contract, event: "AgentUpdated", logs: logs, sub: sub}, nil
}

// WatchAgentUpdated is a free log subscription operation binding the contract event 0x808318fff0e666157af89dad699f694ddee6735f06dba7769a172a8251d34b3f.
//
// Solidity: event AgentUpdated(uint256 indexed agentId, string metadata, uint256[] chainIds)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentUpdated(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentUpdated, agentId []*big.Int) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentUpdated", agentIdRule)
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

// ParseAgentUpdated is a log parse operation binding the contract event 0x808318fff0e666157af89dad699f694ddee6735f06dba7769a172a8251d34b3f.
//
// Solidity: event AgentUpdated(uint256 indexed agentId, string metadata, uint256[] chainIds)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentUpdated(log types.Log) (*AgentRegistryAgentUpdated, error) {
	event := new(AgentRegistryAgentUpdated)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the AgentRegistry contract.
type AgentRegistryApprovalIterator struct {
	Event *AgentRegistryApproval // Event containing the contract specifics and raw log

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
func (it *AgentRegistryApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryApproval)
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
		it.Event = new(AgentRegistryApproval)
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
func (it *AgentRegistryApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryApproval represents a Approval event raised by the AgentRegistry contract.
type AgentRegistryApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_AgentRegistry *AgentRegistryFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*AgentRegistryApprovalIterator, error) {

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

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryApprovalIterator{contract: _AgentRegistry.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_AgentRegistry *AgentRegistryFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AgentRegistryApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryApproval)
				if err := _AgentRegistry.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseApproval(log types.Log) (*AgentRegistryApproval, error) {
	event := new(AgentRegistryApproval)
	if err := _AgentRegistry.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the AgentRegistry contract.
type AgentRegistryApprovalForAllIterator struct {
	Event *AgentRegistryApprovalForAll // Event containing the contract specifics and raw log

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
func (it *AgentRegistryApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryApprovalForAll)
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
		it.Event = new(AgentRegistryApprovalForAll)
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
func (it *AgentRegistryApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryApprovalForAll represents a ApprovalForAll event raised by the AgentRegistry contract.
type AgentRegistryApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_AgentRegistry *AgentRegistryFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*AgentRegistryApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryApprovalForAllIterator{contract: _AgentRegistry.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_AgentRegistry *AgentRegistryFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *AgentRegistryApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryApprovalForAll)
				if err := _AgentRegistry.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseApprovalForAll(log types.Log) (*AgentRegistryApprovalForAll, error) {
	event := new(AgentRegistryApprovalForAll)
	if err := _AgentRegistry.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the AgentRegistry contract.
type AgentRegistryBeaconUpgradedIterator struct {
	Event *AgentRegistryBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *AgentRegistryBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryBeaconUpgraded)
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
		it.Event = new(AgentRegistryBeaconUpgraded)
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
func (it *AgentRegistryBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryBeaconUpgraded represents a BeaconUpgraded event raised by the AgentRegistry contract.
type AgentRegistryBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_AgentRegistry *AgentRegistryFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*AgentRegistryBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryBeaconUpgradedIterator{contract: _AgentRegistry.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_AgentRegistry *AgentRegistryFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *AgentRegistryBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryBeaconUpgraded)
				if err := _AgentRegistry.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseBeaconUpgraded(log types.Log) (*AgentRegistryBeaconUpgraded, error) {
	event := new(AgentRegistryBeaconUpgraded)
	if err := _AgentRegistry.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryRouterUpdatedIterator is returned from FilterRouterUpdated and is used to iterate over the raw logs and unpacked data for RouterUpdated events raised by the AgentRegistry contract.
type AgentRegistryRouterUpdatedIterator struct {
	Event *AgentRegistryRouterUpdated // Event containing the contract specifics and raw log

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
func (it *AgentRegistryRouterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryRouterUpdated)
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
		it.Event = new(AgentRegistryRouterUpdated)
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
func (it *AgentRegistryRouterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryRouterUpdated represents a RouterUpdated event raised by the AgentRegistry contract.
type AgentRegistryRouterUpdated struct {
	Router common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRouterUpdated is a free log retrieval operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_AgentRegistry *AgentRegistryFilterer) FilterRouterUpdated(opts *bind.FilterOpts, router []common.Address) (*AgentRegistryRouterUpdatedIterator, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "RouterUpdated", routerRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryRouterUpdatedIterator{contract: _AgentRegistry.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

// WatchRouterUpdated is a free log subscription operation binding the contract event 0x7aed1d3e8155a07ccf395e44ea3109a0e2d6c9b29bbbe9f142d9790596f4dc80.
//
// Solidity: event RouterUpdated(address indexed router)
func (_AgentRegistry *AgentRegistryFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *AgentRegistryRouterUpdated, router []common.Address) (event.Subscription, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "RouterUpdated", routerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryRouterUpdated)
				if err := _AgentRegistry.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseRouterUpdated(log types.Log) (*AgentRegistryRouterUpdated, error) {
	event := new(AgentRegistryRouterUpdated)
	if err := _AgentRegistry.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the AgentRegistry contract.
type AgentRegistryTransferIterator struct {
	Event *AgentRegistryTransfer // Event containing the contract specifics and raw log

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
func (it *AgentRegistryTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryTransfer)
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
		it.Event = new(AgentRegistryTransfer)
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
func (it *AgentRegistryTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryTransfer represents a Transfer event raised by the AgentRegistry contract.
type AgentRegistryTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_AgentRegistry *AgentRegistryFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*AgentRegistryTransferIterator, error) {

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

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryTransferIterator{contract: _AgentRegistry.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_AgentRegistry *AgentRegistryFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AgentRegistryTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryTransfer)
				if err := _AgentRegistry.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseTransfer(log types.Log) (*AgentRegistryTransfer, error) {
	event := new(AgentRegistryTransfer)
	if err := _AgentRegistry.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the AgentRegistry contract.
type AgentRegistryUpgradedIterator struct {
	Event *AgentRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *AgentRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryUpgraded)
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
		it.Event = new(AgentRegistryUpgraded)
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
func (it *AgentRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryUpgraded represents a Upgraded event raised by the AgentRegistry contract.
type AgentRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AgentRegistry *AgentRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*AgentRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryUpgradedIterator{contract: _AgentRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AgentRegistry *AgentRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *AgentRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryUpgraded)
				if err := _AgentRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseUpgraded(log types.Log) (*AgentRegistryUpgraded, error) {
	event := new(AgentRegistryUpgraded)
	if err := _AgentRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
