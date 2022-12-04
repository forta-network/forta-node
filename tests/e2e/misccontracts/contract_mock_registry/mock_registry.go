// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract_mock_registry

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

// MockRegistryScannerNode is an auto generated low-level Go binding around an user-defined struct.
type MockRegistryScannerNode struct {
	Registered    bool
	Disabled      bool
	ScannerPoolId *big.Int
	ChainId       *big.Int
	Metadata      string
}

// MockRegistryMetaData contains all meta data concerning the MockRegistry contract.
var MockRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"__scannerNodeVersion\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"__agentManifest\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"agentRefAt\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"agentVersion\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"disabledFlags\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"getAgent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"agentVersion\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"}],\"name\":\"getScanner\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"disabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"scannerPoolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"internalType\":\"structMockRegistry.ScannerNode\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"getScanner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkTestAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"numAgentsFor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"scannerHash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"manifest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"scannerNodeVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlinkTestAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200118538038062001185833981810160405281019062000037919062000210565b8160009081620000489190620004e0565b506000600181905550600060028190555080600390816200006a9190620004e0565b505050620005c7565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620000dc8262000091565b810181811067ffffffffffffffff82111715620000fe57620000fd620000a2565b5b80604052505050565b60006200011362000073565b9050620001218282620000d1565b919050565b600067ffffffffffffffff821115620001445762000143620000a2565b5b6200014f8262000091565b9050602081019050919050565b60005b838110156200017c5780820151818401526020810190506200015f565b838111156200018c576000848401525b50505050565b6000620001a9620001a38462000126565b62000107565b905082815260208101848484011115620001c857620001c76200008c565b5b620001d58482856200015c565b509392505050565b600082601f830112620001f557620001f462000087565b5b81516200020784826020860162000192565b91505092915050565b600080604083850312156200022a57620002296200007d565b5b600083015167ffffffffffffffff8111156200024b576200024a62000082565b5b6200025985828601620001dd565b925050602083015167ffffffffffffffff8111156200027d576200027c62000082565b5b6200028b85828601620001dd565b9150509250929050565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620002e857607f821691505b602082108103620002fe57620002fd620002a0565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620003687fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000329565b62000374868362000329565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b6000620003c1620003bb620003b5846200038c565b62000396565b6200038c565b9050919050565b6000819050919050565b620003dd83620003a0565b620003f5620003ec82620003c8565b84845462000336565b825550505050565b600090565b6200040c620003fd565b62000419818484620003d2565b505050565b5b8181101562000441576200043560008262000402565b6001810190506200041f565b5050565b601f82111562000490576200045a8162000304565b620004658462000319565b8101602085101562000475578190505b6200048d620004848562000319565b8301826200041e565b50505b505050565b600082821c905092915050565b6000620004b56000198460080262000495565b1980831691505092915050565b6000620004d08383620004a2565b9150826002028217905092915050565b620004eb8262000295565b67ffffffffffffffff811115620005075762000506620000a2565b5b620005138254620002cf565b6200052082828562000445565b600060209050601f83116001811462000558576000841562000543578287015190505b6200054f8582620004c2565b865550620005bf565b601f198416620005688662000304565b60005b8281101562000592578489015182556001820191506020850194506020810190506200056b565b86831015620005b25784890151620005ae601f891682620004a2565b8355505b6001600288020188555050505b505050505050565b610bae80620005d76000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c8063a97fe43e11610066578063a97fe43e14610166578063b1774f9d14610199578063bd3c3a1a146101ca578063c783034c146101fa578063d20f1f641461022a5761009e565b80632de5aaf7146100a357806332dee2f6146100d7578063345db3e11461010e5780636877063a1461012c5780637c9acefc1461015c575b600080fd5b6100bd60048036038101906100b8919061056e565b610234565b6040516100ce95949392919061075d565b60405180910390f35b6100f160048036038101906100ec91906107be565b610355565b6040516101059897969594939291906107fe565b60405180910390f35b6101166103a6565b604051610123919061088a565b60405180910390f35b610146600480360381019061014191906108d8565b610434565b60405161015391906109d4565b60405180910390f35b610164610485565b005b610180600480360381019061017b919061056e565b610497565b60405161019094939291906109f6565b60405180910390f35b6101b360048036038101906101ae919061056e565b6104c3565b6040516101c1929190610a5b565b60405180910390f35b6101e460048036038101906101df919061056e565b6104d8565b6040516101f19190610a84565b60405180910390f35b610214600480360381019061020f919061056e565b6104e4565b6040516102219190610a9f565b60405180910390f35b6102326104ef565b005b60008060006060806000600167ffffffffffffffff81111561025957610258610aba565b5b6040519080825280602002602001820160405280156102875781602001602082028036833780820191505090505b5090506089816000815181106102a05761029f610ae9565b5b6020026020010181815250506001600060016003848180546102c190610b47565b80601f01602080910402602001604051908101604052809291908181526020018280546102ed90610b47565b801561033a5780601f1061030f5761010080835404028352916020019161033a565b820191906000526020600020905b81548152906001019060200180831161031d57829003601f168201915b50505050509150955095509550955095505091939590929450565b60008060008060608060008061036b6001610234565b809750819850829950839b50849c50505050505087878787878760016000975097509750975097509750975097509295985092959890939650565b600080546103b390610b47565b80601f01602080910402602001604051908101604052809291908181526020018280546103df90610b47565b801561042c5780601f106104015761010080835404028352916020019161042c565b820191906000526020600020905b81548152906001019060200180831161040f57829003601f168201915b505050505081565b61043c610500565b610444610500565b600181600001901515908115158152505060018160200190151590811515815250506001816040018181525050608981606001818152505080915050919050565b60026001819055506000600281905550565b600080600060606001600060896040518060200160405280600081525093509350935093509193509193565b60008060025460015460001b91509150915091565b60006002549050919050565b600060019050919050565b600180819055506001600281905550565b6040518060a001604052806000151581526020016000151581526020016000815260200160008152602001606081525090565b600080fd5b6000819050919050565b61054b81610538565b811461055657600080fd5b50565b60008135905061056881610542565b92915050565b60006020828403121561058457610583610533565b5b600061059284828501610559565b91505092915050565b60008115159050919050565b6105b08161059b565b82525050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006105e1826105b6565b9050919050565b6105f1816105d6565b82525050565b61060081610538565b82525050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610640578082015181840152602081019050610625565b8381111561064f576000848401525b50505050565b6000601f19601f8301169050919050565b600061067182610606565b61067b8185610611565b935061068b818560208601610622565b61069481610655565b840191505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6106d481610538565b82525050565b60006106e683836106cb565b60208301905092915050565b6000602082019050919050565b600061070a8261069f565b61071481856106aa565b935061071f836106bb565b8060005b8381101561075057815161073788826106da565b9750610742836106f2565b925050600181019050610723565b5085935050505092915050565b600060a08201905061077260008301886105a7565b61077f60208301876105e8565b61078c60408301866105f7565b818103606083015261079e8185610666565b905081810360808301526107b281846106ff565b90509695505050505050565b600080604083850312156107d5576107d4610533565b5b60006107e385828601610559565b92505060206107f485828601610559565b9150509250929050565b600061010082019050610814600083018b6105a7565b610821602083018a6105e8565b61082e60408301896105f7565b61083b60608301886105f7565b818103608083015261084d8187610666565b905081810360a083015261086181866106ff565b905061087060c08301856105a7565b61087d60e08301846105f7565b9998505050505050505050565b600060208201905081810360008301526108a48184610666565b905092915050565b6108b5816105d6565b81146108c057600080fd5b50565b6000813590506108d2816108ac565b92915050565b6000602082840312156108ee576108ed610533565b5b60006108fc848285016108c3565b91505092915050565b61090e8161059b565b82525050565b600082825260208201905092915050565b600061093082610606565b61093a8185610914565b935061094a818560208601610622565b61095381610655565b840191505092915050565b600060a0830160008301516109766000860182610905565b5060208301516109896020860182610905565b50604083015161099c60408601826106cb565b5060608301516109af60608601826106cb565b50608083015184820360808601526109c78282610925565b9150508091505092915050565b600060208201905081810360008301526109ee818461095e565b905092915050565b6000608082019050610a0b60008301876105a7565b610a1860208301866105e8565b610a2560408301856105f7565b8181036060830152610a378184610666565b905095945050505050565b6000819050919050565b610a5581610a42565b82525050565b6000604082019050610a7060008301856105f7565b610a7d6020830184610a4c565b9392505050565b6000602082019050610a9960008301846105f7565b92915050565b6000602082019050610ab460008301846105a7565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610b5f57607f821691505b602082108103610b7257610b71610b18565b5b5091905056fea2646970667358221220945d9d50a69c5664de93823a132844c4c9137f4f910690971de703f631ff112864736f6c634300080f0033",
}

// MockRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use MockRegistryMetaData.ABI instead.
var MockRegistryABI = MockRegistryMetaData.ABI

// MockRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockRegistryMetaData.Bin instead.
var MockRegistryBin = MockRegistryMetaData.Bin

// DeployMockRegistry deploys a new Ethereum contract, binding an instance of MockRegistry to it.
func DeployMockRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, __scannerNodeVersion string, __agentManifest string) (common.Address, *types.Transaction, *MockRegistry, error) {
	parsed, err := MockRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockRegistryBin), backend, __scannerNodeVersion, __agentManifest)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockRegistry{MockRegistryCaller: MockRegistryCaller{contract: contract}, MockRegistryTransactor: MockRegistryTransactor{contract: contract}, MockRegistryFilterer: MockRegistryFilterer{contract: contract}}, nil
}

// MockRegistry is an auto generated Go binding around an Ethereum contract.
type MockRegistry struct {
	MockRegistryCaller     // Read-only binding to the contract
	MockRegistryTransactor // Write-only binding to the contract
	MockRegistryFilterer   // Log filterer for contract events
}

// MockRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockRegistrySession struct {
	Contract     *MockRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MockRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockRegistryCallerSession struct {
	Contract *MockRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// MockRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockRegistryTransactorSession struct {
	Contract     *MockRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// MockRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockRegistryRaw struct {
	Contract *MockRegistry // Generic contract binding to access the raw methods on
}

// MockRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockRegistryCallerRaw struct {
	Contract *MockRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// MockRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockRegistryTransactorRaw struct {
	Contract *MockRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockRegistry creates a new instance of MockRegistry, bound to a specific deployed contract.
func NewMockRegistry(address common.Address, backend bind.ContractBackend) (*MockRegistry, error) {
	contract, err := bindMockRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockRegistry{MockRegistryCaller: MockRegistryCaller{contract: contract}, MockRegistryTransactor: MockRegistryTransactor{contract: contract}, MockRegistryFilterer: MockRegistryFilterer{contract: contract}}, nil
}

// NewMockRegistryCaller creates a new read-only instance of MockRegistry, bound to a specific deployed contract.
func NewMockRegistryCaller(address common.Address, caller bind.ContractCaller) (*MockRegistryCaller, error) {
	contract, err := bindMockRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockRegistryCaller{contract: contract}, nil
}

// NewMockRegistryTransactor creates a new write-only instance of MockRegistry, bound to a specific deployed contract.
func NewMockRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*MockRegistryTransactor, error) {
	contract, err := bindMockRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockRegistryTransactor{contract: contract}, nil
}

// NewMockRegistryFilterer creates a new log filterer instance of MockRegistry, bound to a specific deployed contract.
func NewMockRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*MockRegistryFilterer, error) {
	contract, err := bindMockRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockRegistryFilterer{contract: contract}, nil
}

// bindMockRegistry binds a generic wrapper to an already deployed contract.
func bindMockRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MockRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockRegistry *MockRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockRegistry.Contract.MockRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockRegistry *MockRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockRegistry.Contract.MockRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockRegistry *MockRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockRegistry.Contract.MockRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockRegistry *MockRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockRegistry *MockRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockRegistry *MockRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockRegistry.Contract.contract.Transact(opts, method, params...)
}

// AgentRefAt is a free data retrieval call binding the contract method 0x32dee2f6.
//
// Solidity: function agentRefAt(uint256 scannerId, uint256 pos) view returns(bool registered, address owner, uint256 agentId, uint256 agentVersion, string metadata, uint256[] chainIds, bool enabled, uint256 disabledFlags)
func (_MockRegistry *MockRegistryCaller) AgentRefAt(opts *bind.CallOpts, scannerId *big.Int, pos *big.Int) (struct {
	Registered    bool
	Owner         common.Address
	AgentId       *big.Int
	AgentVersion  *big.Int
	Metadata      string
	ChainIds      []*big.Int
	Enabled       bool
	DisabledFlags *big.Int
}, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "agentRefAt", scannerId, pos)

	outstruct := new(struct {
		Registered    bool
		Owner         common.Address
		AgentId       *big.Int
		AgentVersion  *big.Int
		Metadata      string
		ChainIds      []*big.Int
		Enabled       bool
		DisabledFlags *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Registered = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Owner = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.AgentId = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.AgentVersion = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Metadata = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.ChainIds = *abi.ConvertType(out[5], new([]*big.Int)).(*[]*big.Int)
	outstruct.Enabled = *abi.ConvertType(out[6], new(bool)).(*bool)
	outstruct.DisabledFlags = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// AgentRefAt is a free data retrieval call binding the contract method 0x32dee2f6.
//
// Solidity: function agentRefAt(uint256 scannerId, uint256 pos) view returns(bool registered, address owner, uint256 agentId, uint256 agentVersion, string metadata, uint256[] chainIds, bool enabled, uint256 disabledFlags)
func (_MockRegistry *MockRegistrySession) AgentRefAt(scannerId *big.Int, pos *big.Int) (struct {
	Registered    bool
	Owner         common.Address
	AgentId       *big.Int
	AgentVersion  *big.Int
	Metadata      string
	ChainIds      []*big.Int
	Enabled       bool
	DisabledFlags *big.Int
}, error) {
	return _MockRegistry.Contract.AgentRefAt(&_MockRegistry.CallOpts, scannerId, pos)
}

// AgentRefAt is a free data retrieval call binding the contract method 0x32dee2f6.
//
// Solidity: function agentRefAt(uint256 scannerId, uint256 pos) view returns(bool registered, address owner, uint256 agentId, uint256 agentVersion, string metadata, uint256[] chainIds, bool enabled, uint256 disabledFlags)
func (_MockRegistry *MockRegistryCallerSession) AgentRefAt(scannerId *big.Int, pos *big.Int) (struct {
	Registered    bool
	Owner         common.Address
	AgentId       *big.Int
	AgentVersion  *big.Int
	Metadata      string
	ChainIds      []*big.Int
	Enabled       bool
	DisabledFlags *big.Int
}, error) {
	return _MockRegistry.Contract.AgentRefAt(&_MockRegistry.CallOpts, scannerId, pos)
}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns(bool registered, address owner, uint256 agentVersion, string metadata, uint256[] chainIds)
func (_MockRegistry *MockRegistryCaller) GetAgent(opts *bind.CallOpts, agentId *big.Int) (struct {
	Registered   bool
	Owner        common.Address
	AgentVersion *big.Int
	Metadata     string
	ChainIds     []*big.Int
}, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "getAgent", agentId)

	outstruct := new(struct {
		Registered   bool
		Owner        common.Address
		AgentVersion *big.Int
		Metadata     string
		ChainIds     []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Registered = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Owner = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.AgentVersion = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Metadata = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.ChainIds = *abi.ConvertType(out[4], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns(bool registered, address owner, uint256 agentVersion, string metadata, uint256[] chainIds)
func (_MockRegistry *MockRegistrySession) GetAgent(agentId *big.Int) (struct {
	Registered   bool
	Owner        common.Address
	AgentVersion *big.Int
	Metadata     string
	ChainIds     []*big.Int
}, error) {
	return _MockRegistry.Contract.GetAgent(&_MockRegistry.CallOpts, agentId)
}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns(bool registered, address owner, uint256 agentVersion, string metadata, uint256[] chainIds)
func (_MockRegistry *MockRegistryCallerSession) GetAgent(agentId *big.Int) (struct {
	Registered   bool
	Owner        common.Address
	AgentVersion *big.Int
	Metadata     string
	ChainIds     []*big.Int
}, error) {
	return _MockRegistry.Contract.GetAgent(&_MockRegistry.CallOpts, agentId)
}

// GetScanner is a free data retrieval call binding the contract method 0x6877063a.
//
// Solidity: function getScanner(address scanner) view returns((bool,bool,uint256,uint256,string))
func (_MockRegistry *MockRegistryCaller) GetScanner(opts *bind.CallOpts, scanner common.Address) (MockRegistryScannerNode, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "getScanner", scanner)

	if err != nil {
		return *new(MockRegistryScannerNode), err
	}

	out0 := *abi.ConvertType(out[0], new(MockRegistryScannerNode)).(*MockRegistryScannerNode)

	return out0, err

}

// GetScanner is a free data retrieval call binding the contract method 0x6877063a.
//
// Solidity: function getScanner(address scanner) view returns((bool,bool,uint256,uint256,string))
func (_MockRegistry *MockRegistrySession) GetScanner(scanner common.Address) (MockRegistryScannerNode, error) {
	return _MockRegistry.Contract.GetScanner(&_MockRegistry.CallOpts, scanner)
}

// GetScanner is a free data retrieval call binding the contract method 0x6877063a.
//
// Solidity: function getScanner(address scanner) view returns((bool,bool,uint256,uint256,string))
func (_MockRegistry *MockRegistryCallerSession) GetScanner(scanner common.Address) (MockRegistryScannerNode, error) {
	return _MockRegistry.Contract.GetScanner(&_MockRegistry.CallOpts, scanner)
}

// GetScanner0 is a free data retrieval call binding the contract method 0xa97fe43e.
//
// Solidity: function getScanner(uint256 scannerId) view returns(bool registered, address owner, uint256 chainId, string metadata)
func (_MockRegistry *MockRegistryCaller) GetScanner0(opts *bind.CallOpts, scannerId *big.Int) (struct {
	Registered bool
	Owner      common.Address
	ChainId    *big.Int
	Metadata   string
}, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "getScanner0", scannerId)

	outstruct := new(struct {
		Registered bool
		Owner      common.Address
		ChainId    *big.Int
		Metadata   string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Registered = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Owner = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.ChainId = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Metadata = *abi.ConvertType(out[3], new(string)).(*string)

	return *outstruct, err

}

// GetScanner0 is a free data retrieval call binding the contract method 0xa97fe43e.
//
// Solidity: function getScanner(uint256 scannerId) view returns(bool registered, address owner, uint256 chainId, string metadata)
func (_MockRegistry *MockRegistrySession) GetScanner0(scannerId *big.Int) (struct {
	Registered bool
	Owner      common.Address
	ChainId    *big.Int
	Metadata   string
}, error) {
	return _MockRegistry.Contract.GetScanner0(&_MockRegistry.CallOpts, scannerId)
}

// GetScanner0 is a free data retrieval call binding the contract method 0xa97fe43e.
//
// Solidity: function getScanner(uint256 scannerId) view returns(bool registered, address owner, uint256 chainId, string metadata)
func (_MockRegistry *MockRegistryCallerSession) GetScanner0(scannerId *big.Int) (struct {
	Registered bool
	Owner      common.Address
	ChainId    *big.Int
	Metadata   string
}, error) {
	return _MockRegistry.Contract.GetScanner0(&_MockRegistry.CallOpts, scannerId)
}

// IsEnabled is a free data retrieval call binding the contract method 0xc783034c.
//
// Solidity: function isEnabled(uint256 id) view returns(bool)
func (_MockRegistry *MockRegistryCaller) IsEnabled(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "isEnabled", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEnabled is a free data retrieval call binding the contract method 0xc783034c.
//
// Solidity: function isEnabled(uint256 id) view returns(bool)
func (_MockRegistry *MockRegistrySession) IsEnabled(id *big.Int) (bool, error) {
	return _MockRegistry.Contract.IsEnabled(&_MockRegistry.CallOpts, id)
}

// IsEnabled is a free data retrieval call binding the contract method 0xc783034c.
//
// Solidity: function isEnabled(uint256 id) view returns(bool)
func (_MockRegistry *MockRegistryCallerSession) IsEnabled(id *big.Int) (bool, error) {
	return _MockRegistry.Contract.IsEnabled(&_MockRegistry.CallOpts, id)
}

// NumAgentsFor is a free data retrieval call binding the contract method 0xbd3c3a1a.
//
// Solidity: function numAgentsFor(uint256 scannerId) view returns(uint256)
func (_MockRegistry *MockRegistryCaller) NumAgentsFor(opts *bind.CallOpts, scannerId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "numAgentsFor", scannerId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumAgentsFor is a free data retrieval call binding the contract method 0xbd3c3a1a.
//
// Solidity: function numAgentsFor(uint256 scannerId) view returns(uint256)
func (_MockRegistry *MockRegistrySession) NumAgentsFor(scannerId *big.Int) (*big.Int, error) {
	return _MockRegistry.Contract.NumAgentsFor(&_MockRegistry.CallOpts, scannerId)
}

// NumAgentsFor is a free data retrieval call binding the contract method 0xbd3c3a1a.
//
// Solidity: function numAgentsFor(uint256 scannerId) view returns(uint256)
func (_MockRegistry *MockRegistryCallerSession) NumAgentsFor(scannerId *big.Int) (*big.Int, error) {
	return _MockRegistry.Contract.NumAgentsFor(&_MockRegistry.CallOpts, scannerId)
}

// ScannerHash is a free data retrieval call binding the contract method 0xb1774f9d.
//
// Solidity: function scannerHash(uint256 scannerId) view returns(uint256 length, bytes32 manifest)
func (_MockRegistry *MockRegistryCaller) ScannerHash(opts *bind.CallOpts, scannerId *big.Int) (struct {
	Length   *big.Int
	Manifest [32]byte
}, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "scannerHash", scannerId)

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
func (_MockRegistry *MockRegistrySession) ScannerHash(scannerId *big.Int) (struct {
	Length   *big.Int
	Manifest [32]byte
}, error) {
	return _MockRegistry.Contract.ScannerHash(&_MockRegistry.CallOpts, scannerId)
}

// ScannerHash is a free data retrieval call binding the contract method 0xb1774f9d.
//
// Solidity: function scannerHash(uint256 scannerId) view returns(uint256 length, bytes32 manifest)
func (_MockRegistry *MockRegistryCallerSession) ScannerHash(scannerId *big.Int) (struct {
	Length   *big.Int
	Manifest [32]byte
}, error) {
	return _MockRegistry.Contract.ScannerHash(&_MockRegistry.CallOpts, scannerId)
}

// ScannerNodeVersion is a free data retrieval call binding the contract method 0x345db3e1.
//
// Solidity: function scannerNodeVersion() view returns(string)
func (_MockRegistry *MockRegistryCaller) ScannerNodeVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "scannerNodeVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ScannerNodeVersion is a free data retrieval call binding the contract method 0x345db3e1.
//
// Solidity: function scannerNodeVersion() view returns(string)
func (_MockRegistry *MockRegistrySession) ScannerNodeVersion() (string, error) {
	return _MockRegistry.Contract.ScannerNodeVersion(&_MockRegistry.CallOpts)
}

// ScannerNodeVersion is a free data retrieval call binding the contract method 0x345db3e1.
//
// Solidity: function scannerNodeVersion() view returns(string)
func (_MockRegistry *MockRegistryCallerSession) ScannerNodeVersion() (string, error) {
	return _MockRegistry.Contract.ScannerNodeVersion(&_MockRegistry.CallOpts)
}

// LinkTestAgent is a paid mutator transaction binding the contract method 0xd20f1f64.
//
// Solidity: function linkTestAgent() returns()
func (_MockRegistry *MockRegistryTransactor) LinkTestAgent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockRegistry.contract.Transact(opts, "linkTestAgent")
}

// LinkTestAgent is a paid mutator transaction binding the contract method 0xd20f1f64.
//
// Solidity: function linkTestAgent() returns()
func (_MockRegistry *MockRegistrySession) LinkTestAgent() (*types.Transaction, error) {
	return _MockRegistry.Contract.LinkTestAgent(&_MockRegistry.TransactOpts)
}

// LinkTestAgent is a paid mutator transaction binding the contract method 0xd20f1f64.
//
// Solidity: function linkTestAgent() returns()
func (_MockRegistry *MockRegistryTransactorSession) LinkTestAgent() (*types.Transaction, error) {
	return _MockRegistry.Contract.LinkTestAgent(&_MockRegistry.TransactOpts)
}

// UnlinkTestAgent is a paid mutator transaction binding the contract method 0x7c9acefc.
//
// Solidity: function unlinkTestAgent() returns()
func (_MockRegistry *MockRegistryTransactor) UnlinkTestAgent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockRegistry.contract.Transact(opts, "unlinkTestAgent")
}

// UnlinkTestAgent is a paid mutator transaction binding the contract method 0x7c9acefc.
//
// Solidity: function unlinkTestAgent() returns()
func (_MockRegistry *MockRegistrySession) UnlinkTestAgent() (*types.Transaction, error) {
	return _MockRegistry.Contract.UnlinkTestAgent(&_MockRegistry.TransactOpts)
}

// UnlinkTestAgent is a paid mutator transaction binding the contract method 0x7c9acefc.
//
// Solidity: function unlinkTestAgent() returns()
func (_MockRegistry *MockRegistryTransactorSession) UnlinkTestAgent() (*types.Transaction, error) {
	return _MockRegistry.Contract.UnlinkTestAgent(&_MockRegistry.TransactOpts)
}
