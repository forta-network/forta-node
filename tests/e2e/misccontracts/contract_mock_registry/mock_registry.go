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
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"__scannerNodeVersion\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"__agentManifest\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"agentRefAt\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"agentVersion\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"disabledFlags\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"getAgent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"agentVersion\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"}],\"name\":\"getScanner\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"disabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"scannerPoolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"internalType\":\"structMockRegistry.ScannerNode\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"getScanner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkTestAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"numAgentsFor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"scannerHash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"manifest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"scannerNodeVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlinkTestAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620011e7380380620011e7833981810160405281019062000037919062000210565b8160009081620000489190620004e0565b506000600181905550600060028190555080600390816200006a9190620004e0565b505050620005c7565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620000dc8262000091565b810181811067ffffffffffffffff82111715620000fe57620000fd620000a2565b5b80604052505050565b60006200011362000073565b9050620001218282620000d1565b919050565b600067ffffffffffffffff821115620001445762000143620000a2565b5b6200014f8262000091565b9050602081019050919050565b60005b838110156200017c5780820151818401526020810190506200015f565b838111156200018c576000848401525b50505050565b6000620001a9620001a38462000126565b62000107565b905082815260208101848484011115620001c857620001c76200008c565b5b620001d58482856200015c565b509392505050565b600082601f830112620001f557620001f462000087565b5b81516200020784826020860162000192565b91505092915050565b600080604083850312156200022a57620002296200007d565b5b600083015167ffffffffffffffff8111156200024b576200024a62000082565b5b6200025985828601620001dd565b925050602083015167ffffffffffffffff8111156200027d576200027c62000082565b5b6200028b85828601620001dd565b9150509250929050565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620002e857607f821691505b602082108103620002fe57620002fd620002a0565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620003687fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000329565b62000374868362000329565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b6000620003c1620003bb620003b5846200038c565b62000396565b6200038c565b9050919050565b6000819050919050565b620003dd83620003a0565b620003f5620003ec82620003c8565b84845462000336565b825550505050565b600090565b6200040c620003fd565b62000419818484620003d2565b505050565b5b8181101562000441576200043560008262000402565b6001810190506200041f565b5050565b601f82111562000490576200045a8162000304565b620004658462000319565b8101602085101562000475578190505b6200048d620004848562000319565b8301826200041e565b50505b505050565b600082821c905092915050565b6000620004b56000198460080262000495565b1980831691505092915050565b6000620004d08383620004a2565b9150826002028217905092915050565b620004eb8262000295565b67ffffffffffffffff811115620005075762000506620000a2565b5b620005138254620002cf565b6200052082828562000445565b600060209050601f83116001811462000558576000841562000543578287015190505b6200054f8582620004c2565b865550620005bf565b601f198416620005688662000304565b60005b8281101562000592578489015182556001820191506020850194506020810190506200056b565b86831015620005b25784890151620005ae601f891682620004a2565b8355505b6001600288020188555050505b505050505050565b610c1080620005d76000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80637c9acefc116100715780637c9acefc14610185578063a97fe43e1461018f578063b1774f9d146101c2578063bd3c3a1a146101f3578063c783034c14610223578063d20f1f6414610253576100a9565b80632de5aaf7146100ae57806332dee2f6146100e2578063345db3e11461011957806354fd4d50146101375780636877063a14610155575b600080fd5b6100c860048036038101906100c391906105d0565b61025d565b6040516100d99594939291906107bf565b60405180910390f35b6100fc60048036038101906100f79190610820565b61037e565b604051610110989796959493929190610860565b60405180910390f35b6101216103cf565b60405161012e91906108ec565b60405180910390f35b61013f61045d565b60405161014c91906108ec565b60405180910390f35b61016f600480360381019061016a919061093a565b610496565b60405161017c9190610a36565b60405180910390f35b61018d6104e7565b005b6101a960048036038101906101a491906105d0565b6104f9565b6040516101b99493929190610a58565b60405180910390f35b6101dc60048036038101906101d791906105d0565b610525565b6040516101ea929190610abd565b60405180910390f35b61020d600480360381019061020891906105d0565b61053a565b60405161021a9190610ae6565b60405180910390f35b61023d600480360381019061023891906105d0565b610546565b60405161024a9190610b01565b60405180910390f35b61025b610551565b005b60008060006060806000600167ffffffffffffffff81111561028257610281610b1c565b5b6040519080825280602002602001820160405280156102b05781602001602082028036833780820191505090505b5090506089816000815181106102c9576102c8610b4b565b5b6020026020010181815250506001600060016003848180546102ea90610ba9565b80601f016020809104026020016040519081016040528092919081815260200182805461031690610ba9565b80156103635780601f1061033857610100808354040283529160200191610363565b820191906000526020600020905b81548152906001019060200180831161034657829003601f168201915b50505050509150955095509550955095505091939590929450565b600080600080606080600080610394600161025d565b809750819850829950839b50849c50505050505087878787878760016000975097509750975097509750975097509295985092959890939650565b600080546103dc90610ba9565b80601f016020809104026020016040519081016040528092919081815260200182805461040890610ba9565b80156104555780601f1061042a57610100808354040283529160200191610455565b820191906000526020600020905b81548152906001019060200180831161043857829003601f168201915b505050505081565b6040518060400160405280600581526020017f302e302e3100000000000000000000000000000000000000000000000000000081525081565b61049e610562565b6104a6610562565b600181600001901515908115158152505060018160200190151590811515815250506001816040018181525050608981606001818152505080915050919050565b60026001819055506000600281905550565b600080600060606001600060896040518060200160405280600081525093509350935093509193509193565b60008060025460015460001b91509150915091565b60006002549050919050565b600060019050919050565b600180819055506001600281905550565b6040518060a001604052806000151581526020016000151581526020016000815260200160008152602001606081525090565b600080fd5b6000819050919050565b6105ad8161059a565b81146105b857600080fd5b50565b6000813590506105ca816105a4565b92915050565b6000602082840312156105e6576105e5610595565b5b60006105f4848285016105bb565b91505092915050565b60008115159050919050565b610612816105fd565b82525050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061064382610618565b9050919050565b61065381610638565b82525050565b6106628161059a565b82525050565b600081519050919050565b600082825260208201905092915050565b60005b838110156106a2578082015181840152602081019050610687565b838111156106b1576000848401525b50505050565b6000601f19601f8301169050919050565b60006106d382610668565b6106dd8185610673565b93506106ed818560208601610684565b6106f6816106b7565b840191505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6107368161059a565b82525050565b6000610748838361072d565b60208301905092915050565b6000602082019050919050565b600061076c82610701565b610776818561070c565b93506107818361071d565b8060005b838110156107b2578151610799888261073c565b97506107a483610754565b925050600181019050610785565b5085935050505092915050565b600060a0820190506107d46000830188610609565b6107e1602083018761064a565b6107ee6040830186610659565b818103606083015261080081856106c8565b905081810360808301526108148184610761565b90509695505050505050565b6000806040838503121561083757610836610595565b5b6000610845858286016105bb565b9250506020610856858286016105bb565b9150509250929050565b600061010082019050610876600083018b610609565b610883602083018a61064a565b6108906040830189610659565b61089d6060830188610659565b81810360808301526108af81876106c8565b905081810360a08301526108c38186610761565b90506108d260c0830185610609565b6108df60e0830184610659565b9998505050505050505050565b6000602082019050818103600083015261090681846106c8565b905092915050565b61091781610638565b811461092257600080fd5b50565b6000813590506109348161090e565b92915050565b6000602082840312156109505761094f610595565b5b600061095e84828501610925565b91505092915050565b610970816105fd565b82525050565b600082825260208201905092915050565b600061099282610668565b61099c8185610976565b93506109ac818560208601610684565b6109b5816106b7565b840191505092915050565b600060a0830160008301516109d86000860182610967565b5060208301516109eb6020860182610967565b5060408301516109fe604086018261072d565b506060830151610a11606086018261072d565b5060808301518482036080860152610a298282610987565b9150508091505092915050565b60006020820190508181036000830152610a5081846109c0565b905092915050565b6000608082019050610a6d6000830187610609565b610a7a602083018661064a565b610a876040830185610659565b8181036060830152610a9981846106c8565b905095945050505050565b6000819050919050565b610ab781610aa4565b82525050565b6000604082019050610ad26000830185610659565b610adf6020830184610aae565b9392505050565b6000602082019050610afb6000830184610659565b92915050565b6000602082019050610b166000830184610609565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610bc157607f821691505b602082108103610bd457610bd3610b7a565b5b5091905056fea26469706673582212201d7918d7b4aedebb60684f503e2f0b99fc7c9ece2c4a0d19ec3ed550aebe981064736f6c634300080f0033",
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

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_MockRegistry *MockRegistryCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_MockRegistry *MockRegistrySession) Version() (string, error) {
	return _MockRegistry.Contract.Version(&_MockRegistry.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_MockRegistry *MockRegistryCallerSession) Version() (string, error) {
	return _MockRegistry.Contract.Version(&_MockRegistry.CallOpts)
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
