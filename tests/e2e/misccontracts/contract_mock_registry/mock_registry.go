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
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"__scannerNodeVersion\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"__agentManifest\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"agentRefAt\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"agentVersion\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"disabledFlags\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"getAgent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"agentVersion\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"chainIds\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scanner\",\"type\":\"address\"}],\"name\":\"getScanner\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"disabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"scannerPoolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"internalType\":\"structMockRegistry.ScannerNode\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"getScanner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkTestAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"numAgentsFor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"numScannersFor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"}],\"name\":\"scannerHash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"manifest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"scannerNodeVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"scannerRefAt\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"scannerId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"disabledFlags\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlinkTestAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200135e3803806200135e833981810160405281019062000037919062000210565b8160009081620000489190620004e0565b506000600181905550600060028190555080600390816200006a9190620004e0565b505050620005c7565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620000dc8262000091565b810181811067ffffffffffffffff82111715620000fe57620000fd620000a2565b5b80604052505050565b60006200011362000073565b9050620001218282620000d1565b919050565b600067ffffffffffffffff821115620001445762000143620000a2565b5b6200014f8262000091565b9050602081019050919050565b60005b838110156200017c5780820151818401526020810190506200015f565b838111156200018c576000848401525b50505050565b6000620001a9620001a38462000126565b62000107565b905082815260208101848484011115620001c857620001c76200008c565b5b620001d58482856200015c565b509392505050565b600082601f830112620001f557620001f462000087565b5b81516200020784826020860162000192565b91505092915050565b600080604083850312156200022a57620002296200007d565b5b600083015167ffffffffffffffff8111156200024b576200024a62000082565b5b6200025985828601620001dd565b925050602083015167ffffffffffffffff8111156200027d576200027c62000082565b5b6200028b85828601620001dd565b9150509250929050565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620002e857607f821691505b602082108103620002fe57620002fd620002a0565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620003687fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000329565b62000374868362000329565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b6000620003c1620003bb620003b5846200038c565b62000396565b6200038c565b9050919050565b6000819050919050565b620003dd83620003a0565b620003f5620003ec82620003c8565b84845462000336565b825550505050565b600090565b6200040c620003fd565b62000419818484620003d2565b505050565b5b8181101562000441576200043560008262000402565b6001810190506200041f565b5050565b601f82111562000490576200045a8162000304565b620004658462000319565b8101602085101562000475578190505b6200048d620004848562000319565b8301826200041e565b50505b505050565b600082821c905092915050565b6000620004b56000198460080262000495565b1980831691505092915050565b6000620004d08383620004a2565b9150826002028217905092915050565b620004eb8262000295565b67ffffffffffffffff811115620005075762000506620000a2565b5b620005138254620002cf565b6200052082828562000445565b600060209050601f83116001811462000558576000841562000543578287015190505b6200054f8582620004c2565b865550620005bf565b601f198416620005688662000304565b60005b8281101562000592578489015182556001820191506020850194506020810190506200056b565b86831015620005b25784890151620005ae601f891682620004a2565b8355505b6001600288020188555050505b505050505050565b610d8780620005d76000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c80637c9acefc1161008c578063b1774f9d11610066578063b1774f9d1461024e578063bd3c3a1a1461027f578063c783034c146102af578063d20f1f64146102df576100cf565b80637c9acefc146101db5780638b2e98d6146101e5578063a97fe43e1461021b576100cf565b80632de5aaf7146100d457806332dee2f614610108578063345db3e11461013f5780633820d2431461015d57806354fd4d501461018d5780636877063a146101ab575b600080fd5b6100ee60048036038101906100e991906106d1565b6102e9565b6040516100ff9594939291906108c0565b60405180910390f35b610122600480360381019061011d9190610921565b61040a565b604051610136989796959493929190610961565b60405180910390f35b61014761045b565b60405161015491906109ed565b60405180910390f35b610177600480360381019061017291906106d1565b6104e9565b6040516101849190610a0f565b60405180910390f35b6101956104f4565b6040516101a291906109ed565b60405180910390f35b6101c560048036038101906101c09190610a56565b61052d565b6040516101d29190610b52565b60405180910390f35b6101e361057e565b005b6101ff60048036038101906101fa9190610921565b610590565b6040516102129796959493929190610b74565b60405180910390f35b610235600480360381019061023091906106d1565b6105fa565b6040516102459493929190610bea565b60405180910390f35b610268600480360381019061026391906106d1565b610626565b604051610276929190610c4f565b60405180910390f35b610299600480360381019061029491906106d1565b61063b565b6040516102a69190610a0f565b60405180910390f35b6102c960048036038101906102c491906106d1565b610647565b6040516102d69190610c78565b60405180910390f35b6102e7610652565b005b60008060006060806000600167ffffffffffffffff81111561030e5761030d610c93565b5b60405190808252806020026020018201604052801561033c5781602001602082028036833780820191505090505b50905060898160008151811061035557610354610cc2565b5b60200260200101818152505060016000600160038481805461037690610d20565b80601f01602080910402602001604051908101604052809291908181526020018280546103a290610d20565b80156103ef5780601f106103c4576101008083540402835291602001916103ef565b820191906000526020600020905b8154815290600101906020018083116103d257829003601f168201915b50505050509150955095509550955095505091939590929450565b60008060008060608060008061042060016102e9565b809750819850829950839b50849c50505050505087878787878760016000975097509750975097509750975097509295985092959890939650565b6000805461046890610d20565b80601f016020809104026020016040519081016040528092919081815260200182805461049490610d20565b80156104e15780601f106104b6576101008083540402835291602001916104e1565b820191906000526020600020905b8154815290600101906020018083116104c457829003601f168201915b505050505081565b600060019050919050565b6040518060400160405280600581526020017f302e302e3100000000000000000000000000000000000000000000000000000081525081565b610535610663565b61053d610663565b600181600001901515908115158152505060018160200190151590811515815250506001816040018181525050608981606001818152505080915050919050565b60026001819055506000600281905550565b6000806000806060600080600173222244861c15a8f2a05fbd15e747ea8f20c2c0c973ffffffffffffffffffffffffffffffffffffffff16600060896040518060200160405280600081525060016000965096509650965096509650965092959891949750929550565b600080600060606001600060896040518060200160405280600081525093509350935093509193509193565b60008060025460015460001b91509150915091565b60006002549050919050565b600060019050919050565b600180819055506001600281905550565b6040518060a001604052806000151581526020016000151581526020016000815260200160008152602001606081525090565b600080fd5b6000819050919050565b6106ae8161069b565b81146106b957600080fd5b50565b6000813590506106cb816106a5565b92915050565b6000602082840312156106e7576106e6610696565b5b60006106f5848285016106bc565b91505092915050565b60008115159050919050565b610713816106fe565b82525050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061074482610719565b9050919050565b61075481610739565b82525050565b6107638161069b565b82525050565b600081519050919050565b600082825260208201905092915050565b60005b838110156107a3578082015181840152602081019050610788565b838111156107b2576000848401525b50505050565b6000601f19601f8301169050919050565b60006107d482610769565b6107de8185610774565b93506107ee818560208601610785565b6107f7816107b8565b840191505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6108378161069b565b82525050565b6000610849838361082e565b60208301905092915050565b6000602082019050919050565b600061086d82610802565b610877818561080d565b93506108828361081e565b8060005b838110156108b357815161089a888261083d565b97506108a583610855565b925050600181019050610886565b5085935050505092915050565b600060a0820190506108d5600083018861070a565b6108e2602083018761074b565b6108ef604083018661075a565b818103606083015261090181856107c9565b905081810360808301526109158184610862565b90509695505050505050565b6000806040838503121561093857610937610696565b5b6000610946858286016106bc565b9250506020610957858286016106bc565b9150509250929050565b600061010082019050610977600083018b61070a565b610984602083018a61074b565b610991604083018961075a565b61099e606083018861075a565b81810360808301526109b081876107c9565b905081810360a08301526109c48186610862565b90506109d360c083018561070a565b6109e060e083018461075a565b9998505050505050505050565b60006020820190508181036000830152610a0781846107c9565b905092915050565b6000602082019050610a24600083018461075a565b92915050565b610a3381610739565b8114610a3e57600080fd5b50565b600081359050610a5081610a2a565b92915050565b600060208284031215610a6c57610a6b610696565b5b6000610a7a84828501610a41565b91505092915050565b610a8c816106fe565b82525050565b600082825260208201905092915050565b6000610aae82610769565b610ab88185610a92565b9350610ac8818560208601610785565b610ad1816107b8565b840191505092915050565b600060a083016000830151610af46000860182610a83565b506020830151610b076020860182610a83565b506040830151610b1a604086018261082e565b506060830151610b2d606086018261082e565b5060808301518482036080860152610b458282610aa3565b9150508091505092915050565b60006020820190508181036000830152610b6c8184610adc565b905092915050565b600060e082019050610b89600083018a61070a565b610b96602083018961075a565b610ba3604083018861074b565b610bb0606083018761075a565b8181036080830152610bc281866107c9565b9050610bd160a083018561070a565b610bde60c083018461075a565b98975050505050505050565b6000608082019050610bff600083018761070a565b610c0c602083018661074b565b610c19604083018561075a565b8181036060830152610c2b81846107c9565b905095945050505050565b6000819050919050565b610c4981610c36565b82525050565b6000604082019050610c64600083018561075a565b610c716020830184610c40565b9392505050565b6000602082019050610c8d600083018461070a565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610d3857607f821691505b602082108103610d4b57610d4a610cf1565b5b5091905056fea2646970667358221220127ec93d6fd02127e30de9a1d3f0e1e5f1bcc71974456c2a0952d5d2bfe316ce64736f6c634300080f0033",
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

// NumScannersFor is a free data retrieval call binding the contract method 0x3820d243.
//
// Solidity: function numScannersFor(uint256 agentId) view returns(uint256 count)
func (_MockRegistry *MockRegistryCaller) NumScannersFor(opts *bind.CallOpts, agentId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "numScannersFor", agentId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumScannersFor is a free data retrieval call binding the contract method 0x3820d243.
//
// Solidity: function numScannersFor(uint256 agentId) view returns(uint256 count)
func (_MockRegistry *MockRegistrySession) NumScannersFor(agentId *big.Int) (*big.Int, error) {
	return _MockRegistry.Contract.NumScannersFor(&_MockRegistry.CallOpts, agentId)
}

// NumScannersFor is a free data retrieval call binding the contract method 0x3820d243.
//
// Solidity: function numScannersFor(uint256 agentId) view returns(uint256 count)
func (_MockRegistry *MockRegistryCallerSession) NumScannersFor(agentId *big.Int) (*big.Int, error) {
	return _MockRegistry.Contract.NumScannersFor(&_MockRegistry.CallOpts, agentId)
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

// ScannerRefAt is a free data retrieval call binding the contract method 0x8b2e98d6.
//
// Solidity: function scannerRefAt(uint256 agentId, uint256 pos) view returns(bool registered, uint256 scannerId, address owner, uint256 chainId, string metadata, bool enabled, uint256 disabledFlags)
func (_MockRegistry *MockRegistryCaller) ScannerRefAt(opts *bind.CallOpts, agentId *big.Int, pos *big.Int) (struct {
	Registered    bool
	ScannerId     *big.Int
	Owner         common.Address
	ChainId       *big.Int
	Metadata      string
	Enabled       bool
	DisabledFlags *big.Int
}, error) {
	var out []interface{}
	err := _MockRegistry.contract.Call(opts, &out, "scannerRefAt", agentId, pos)

	outstruct := new(struct {
		Registered    bool
		ScannerId     *big.Int
		Owner         common.Address
		ChainId       *big.Int
		Metadata      string
		Enabled       bool
		DisabledFlags *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Registered = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ScannerId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Owner = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Metadata = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.Enabled = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.DisabledFlags = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ScannerRefAt is a free data retrieval call binding the contract method 0x8b2e98d6.
//
// Solidity: function scannerRefAt(uint256 agentId, uint256 pos) view returns(bool registered, uint256 scannerId, address owner, uint256 chainId, string metadata, bool enabled, uint256 disabledFlags)
func (_MockRegistry *MockRegistrySession) ScannerRefAt(agentId *big.Int, pos *big.Int) (struct {
	Registered    bool
	ScannerId     *big.Int
	Owner         common.Address
	ChainId       *big.Int
	Metadata      string
	Enabled       bool
	DisabledFlags *big.Int
}, error) {
	return _MockRegistry.Contract.ScannerRefAt(&_MockRegistry.CallOpts, agentId, pos)
}

// ScannerRefAt is a free data retrieval call binding the contract method 0x8b2e98d6.
//
// Solidity: function scannerRefAt(uint256 agentId, uint256 pos) view returns(bool registered, uint256 scannerId, address owner, uint256 chainId, string metadata, bool enabled, uint256 disabledFlags)
func (_MockRegistry *MockRegistryCallerSession) ScannerRefAt(agentId *big.Int, pos *big.Int) (struct {
	Registered    bool
	ScannerId     *big.Int
	Owner         common.Address
	ChainId       *big.Int
	Metadata      string
	Enabled       bool
	DisabledFlags *big.Int
}, error) {
	return _MockRegistry.Contract.ScannerRefAt(&_MockRegistry.CallOpts, agentId, pos)
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
