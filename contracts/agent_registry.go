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

// AgentRegistryABI is the input ABI used to generate the binding from.
const AgentRegistryABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdminAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentAdminRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"agentId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AgentUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"PoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"PoolAdminAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"PoolAdminRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PoolOwnershipTransfered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_ref\",\"type\":\"string\"}],\"name\":\"addAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addAgentAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"}],\"name\":\"addPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addPoolAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"agentAdmins\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"agentAt\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"}],\"name\":\"agentLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"poolAdmins\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"poolExistsMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"poolOwners\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"}],\"name\":\"removeAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removeAgentAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removePoolAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferPoolOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_agentId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_ref\",\"type\":\"string\"}],\"name\":\"updateAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// AgentRegistryFuncSigs maps the 4-byte function signature to its string representation.
var AgentRegistryFuncSigs = map[string]string{
	"0944cd34": "addAgent(bytes32,bytes32,string)",
	"63438a97": "addAgentAdmin(bytes32,bytes32,address)",
	"8a0033fc": "addPool(bytes32)",
	"84930b76": "addPoolAdmin(bytes32,address)",
	"9dfeddae": "agentAdmins(bytes32,bytes32,address)",
	"89d0a013": "agentAt(bytes32,uint256)",
	"1e73feb1": "agentLength(bytes32)",
	"250ca0f0": "poolAdmins(bytes32,address)",
	"0ead48e6": "poolExistsMap(bytes32)",
	"219d2d09": "poolOwners(bytes32)",
	"61e62a0f": "removeAgent(bytes32,bytes32)",
	"a7f2f5ab": "removeAgentAdmin(bytes32,bytes32,address)",
	"acce5630": "removePoolAdmin(bytes32,address)",
	"1b6feb92": "transferPoolOwnership(bytes32,address)",
	"e0434e3f": "updateAgent(bytes32,bytes32,string)",
}

// AgentRegistryBin is the compiled bytecode used for deploying new contracts.
var AgentRegistryBin = "0x608060405234801561001057600080fd5b506117b2806100206000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806363438a97116100975780639dfeddae116100665780639dfeddae14610257578063a7f2f5ab1461028b578063acce56301461029e578063e0434e3f146102b157600080fd5b806363438a97146101fd57806384930b761461021057806389d0a013146102235780638a0033fc1461024457600080fd5b80631e73feb1116100d35780631e73feb11461015a578063219d2d091461017b578063250ca0f0146101bc57806361e62a0f146101ea57600080fd5b80630944cd34146100fa5780630ead48e61461010f5780631b6feb9214610147575b600080fd5b61010d6101083660046113ea565b6102c4565b005b61013261011d366004611352565b60016020526000908152604090205460ff1681565b60405190151581526020015b60405180910390f35b61010d61015536600461136a565b61041a565b61016d610168366004611352565b610587565b60405190815260200161013e565b6101a4610189366004611352565b6003602052600090815260409020546001600160a01b031681565b6040516001600160a01b03909116815260200161013e565b6101326101ca36600461136a565b600260209081526000928352604080842090915290825290205460ff1681565b61010d6101f8366004611395565b6105a4565b61010d61020b3660046113b6565b6106f5565b61010d61021e36600461136a565b610875565b610236610231366004611395565b6109aa565b60405161013e9291906115c7565b61010d610252366004611352565b610a5c565b6101326102653660046113b6565b600460209081526000938452604080852082529284528284209052825290205460ff1681565b61010d6102993660046113b6565b610b24565b61010d6102ac36600461136a565b610c96565b61010d6102bf3660046113ea565b610da1565b600083815260036020526040902054839083906001600160a01b03163314806103065750600082815260026020908152604080832033845290915290205460ff165b8061033257506000828152600460209081526040808320848452825280832033845290915290205460ff165b6103575760405162461bcd60e51b815260040161034e90611674565b60405180910390fd5b600085815260208190526040902061036f9085610f81565b156103bc5760405162461bcd60e51b815260206004820152601c60248201527f4167656e7420616c726561647920657869737473206f6e20706f6f6c00000000604482015260640161034e565b60008581526020819052604090206103d5908585610f94565b507fdd1417d47ca73bdaf24e5c1df4158e61cbbb73fda6562b14a16215ccfa201aeb8585853360405161040b949392919061158f565b60405180910390a15050505050565b60008281526003602052604090205482906001600160a01b031633146104525760405162461bcd60e51b815260040161034e906116d1565b6001600160a01b0382166104a85760405162461bcd60e51b815260206004820152601960248201527f61646472657373283029206973206e6f7420616c6c6f77656400000000000000604482015260640161034e565b6000838152600360205260409020546001600160a01b03838116911614156105125760405162461bcd60e51b815260206004820152601860248201527f4164647265737320697320616c7265616479206f776e65720000000000000000604482015260640161034e565b60008381526003602090815260409182902080546001600160a01b0319166001600160a01b03861690811790915582518681523392810192909252918101919091527fecc1c418cc7d2c3062b60c1d8f0b7748e2c8c10a109df4f6089b285f26b993d6906060015b60405180910390a1505050565b600081815260208190526040812061059e90610fc7565b92915050565b600082815260036020526040902054829082906001600160a01b03163314806105e65750600082815260026020908152604080832033845290915290205460ff165b8061061257506000828152600460209081526040808320848452825280832033845290915290205460ff165b61062e5760405162461bcd60e51b815260040161034e90611674565b60008481526020819052604090206106469084610f81565b6106925760405162461bcd60e51b815260206004820152601c60248201527f4167656e7420646f6573206e6f74206578697374206f6e20706f6f6c00000000604482015260640161034e565b60008481526020819052604090206106aa9084610fd2565b50604080518581526020810185905233918101919091527fccc61238b7da918476155ec78f1d6f23cd587a3f796774b80c0c79dde8f57e6b906060015b60405180910390a150505050565b60008381526003602052604090205483906001600160a01b03163314806107355750600081815260026020908152604080832033845290915290205460ff165b6107515760405162461bcd60e51b815260040161034e90611617565b6001600160a01b0382166107775760405162461bcd60e51b815260040161034e906115e0565b600084815260046020908152604080832086845282528083206001600160a01b038616845290915290205460ff16156107fc5760405162461bcd60e51b815260206004820152602160248201527f4164647265737320697320616c726561647920616e206167656e74206f776e656044820152603960f91b606482015260840161034e565b600084815260046020908152604080832086845282528083206001600160a01b03861680855290835292819020805460ff1916600117905580518781529182018690528101919091523360608201527f11cbfe7067a564a8ddb9a780914987915bc0d29835c8b0554a9d9db2e51064ed906080016106e7565b60008281526003602052604090205482906001600160a01b031633146108ad5760405162461bcd60e51b815260040161034e906116d1565b6001600160a01b0382166108d35760405162461bcd60e51b815260040161034e906115e0565b60008381526002602090815260408083206001600160a01b038616845290915290205460ff16156109465760405162461bcd60e51b815260206004820152601b60248201527f4164647265737320697320616c726561647920616e2061646d696e0000000000604482015260640161034e565b60008381526002602090815260408083206001600160a01b03861680855290835292819020805460ff191660011790558051868152918201929092527f4c2f0be874d7047c389e9fadf454854a8f5b72b49de199c7d60831b6c462a9a0910161057a565b60008281526020819052604081206060906109c59084610ff6565b8080546109d19061172b565b80601f01602080910402602001604051908101604052809291908181526020018280546109fd9061172b565b8015610a4a5780601f10610a1f57610100808354040283529160200191610a4a565b820191906000526020600020905b815481529060010190602001808311610a2d57829003601f168201915b50505050509050915091509250929050565b60008181526001602052604090205460ff1615610ab15760405162461bcd60e51b8152602060048201526013602482015272506f6f6c20616c72656164792065786973747360681b604482015260640161034e565b600081815260036020908152604080832080546001600160a01b03191633908117909155600180845293829020805460ff19169094179093558051848152918201929092527ffd0fa7919fbe3857a4236750e8d3e42ac691881d2f31c50ebe4cfc2a7705ee80910160405180910390a150565b60008381526003602052604090205483906001600160a01b0316331480610b645750600081815260026020908152604080832033845290915290205460ff165b610b805760405162461bcd60e51b815260040161034e90611617565b6001600160a01b038216610ba65760405162461bcd60e51b815260040161034e906115e0565b600084815260046020908152604080832086845282528083206001600160a01b038616845290915290205460ff16610c205760405162461bcd60e51b815260206004820152601d60248201527f41646472657373206973206e6f7420616e206167656e74206f776e6572000000604482015260640161034e565b600084815260046020908152604080832086845282528083206001600160a01b03861680855290835292819020805460ff1916905580518781529182018690528101919091523360608201527f801fe4d9d51c0445c2938970ad5e042c3f2f2c387d67de4111af1bc60d22dba9906080016106e7565b60008281526003602052604090205482906001600160a01b03163314610cce5760405162461bcd60e51b815260040161034e906116d1565b60008381526002602090815260408083206001600160a01b038616845290915290205460ff16610d405760405162461bcd60e51b815260206004820152601760248201527f41646472657373206973206e6f7420616e2061646d696e000000000000000000604482015260640161034e565b60008381526002602090815260408083206001600160a01b03861680855290835292819020805460ff191690558051868152918201929092527f25fcb4cb225521d9a561fd90af5809d997ae0fa769a5c95d90750162623bf624910161057a565b600083815260036020526040902054839083906001600160a01b0316331480610de35750600082815260026020908152604080832033845290915290205460ff165b80610e0f57506000828152600460209081526040808320848452825280832033845290915290205460ff165b610e2b5760405162461bcd60e51b815260040161034e90611674565b6000858152602081905260409020610e439085610f81565b610e8f5760405162461bcd60e51b815260206004820152601e60248201527f4167656e74206d75737420657869737420746f20626520757064617465640000604482015260640161034e565b82516020808501919091206000878152808352604080822088835260020190935291909120604051610ec191906114f4565b60405180910390201415610f325760405162461bcd60e51b815260206004820152603260248201527f4e6577207265666572656e6365206d75737420626520646966666572656e74206044820152717468616e206f6c64207265666572656e636560701b606482015260840161034e565b6000858152602081905260409020610f4b908585610f94565b507f8a33912f31f80df098dcad04aba17afc8c52d0f700baf2750ea3899f9230a8508585853360405161040b949392919061158f565b6000610f8d8383611020565b9392505050565b6000828152600284016020908152604082208351610fb492850190611260565b50610fbf8484611038565b949350505050565b600061059e82611044565b60008181526002830160205260408120610fec90826112e4565b610f8d838361104e565b60008080611004858561105a565b6000818152600296909601602052604090952094959350505050565b60008181526001830160205260408120541515610f8d565b6000610f8d8383611066565b600061059e825490565b6000610f8d83836110b5565b6000610f8d83836111cc565b60008181526001830160205260408120546110ad5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561059e565b50600061059e565b600081815260018301602052604081205480156111c25760006110d9600183611708565b85549091506000906110ed90600190611708565b9050600086600001828154811061111457634e487b7160e01b600052603260045260246000fd5b906000526020600020015490508087600001848154811061114557634e487b7160e01b600052603260045260246000fd5b60009182526020808320909101929092558281526001890190915260409020849055865487908061118657634e487b7160e01b600052603160045260246000fd5b6001900381819060005260206000200160009055905586600101600087815260200190815260200160002060009055600194505050505061059e565b600091505061059e565b8154600090821061122a5760405162461bcd60e51b815260206004820152602260248201527f456e756d657261626c655365743a20696e646578206f7574206f6620626f756e604482015261647360f01b606482015260840161034e565b82600001828154811061124d57634e487b7160e01b600052603260045260246000fd5b9060005260206000200154905092915050565b82805461126c9061172b565b90600052602060002090601f01602090048101928261128e57600085556112d4565b82601f106112a757805160ff19168380011785556112d4565b828001600101855582156112d4579182015b828111156112d45782518255916020019190600101906112b9565b506112e0929150611321565b5090565b5080546112f09061172b565b6000825580601f10611300575050565b601f01602090049060005260206000209081019061131e9190611321565b50565b5b808211156112e05760008155600101611322565b80356001600160a01b038116811461134d57600080fd5b919050565b600060208284031215611363578081fd5b5035919050565b6000806040838503121561137c578081fd5b8235915061138c60208401611336565b90509250929050565b600080604083850312156113a7578182fd5b50508035926020909101359150565b6000806000606084860312156113ca578081fd5b83359250602084013591506113e160408501611336565b90509250925092565b6000806000606084860312156113fe578283fd5b8335925060208401359150604084013567ffffffffffffffff80821115611423578283fd5b818601915086601f830112611436578283fd5b81358181111561144857611448611766565b604051601f8201601f19908116603f0116810190838211818310171561147057611470611766565b81604052828152896020848701011115611488578586fd5b82602086016020830137856020848301015280955050505050509250925092565b60008151808452815b818110156114ce576020818501810151868301820152016114b2565b818111156114df5782602083870101525b50601f01601f19169290920160200192915050565b600080835482600182811c91508083168061151057607f831692505b602080841082141561153057634e487b7160e01b87526022600452602487fd5b818015611544576001811461155557611581565b60ff19861689528489019650611581565b60008a815260209020885b868110156115795781548b820152908501908301611560565b505084890196505b509498975050505050505050565b8481528360208201526080604082015260006115ae60808301856114a9565b905060018060a01b038316606083015295945050505050565b828152604060208201526000610fbf60408301846114a9565b60208082526019908201527f41646472657373283029206973206e6f7420616c6c6f77656400000000000000604082015260600190565b60208082526038908201527f4f6e6c7920706f6f6c206f776e6572206f7220706f6f6c2061646d696e20636160408201527f6e20706572666f726d2074686973206f7065726174696f6e0000000000000000606082015260800190565b6020808252603c908201527f4f6e6c7920706f6f6c206f776e65722c20706f6f6c2061646d696e2c206f722060408201527f6167656e74206f776e65722063616e20757064617465206167656e7400000000606082015260800190565b6020808252601b908201527f43616c6c6572206973206e6f74206f776e6572206f6620706f6f6c0000000000604082015260600190565b60008282101561172657634e487b7160e01b81526011600452602481fd5b500390565b600181811c9082168061173f57607f821691505b6020821081141561176057634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052604160045260246000fdfea2646970667358221220133e2349e500916364bdf46b79df9ef7004bab638f34354db68faa23605afe5064736f6c63430008040033"

// DeployAgentRegistry deploys a new Ethereum contract, binding an instance of AgentRegistry to it.
func DeployAgentRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AgentRegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(AgentRegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AgentRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AgentRegistry{AgentRegistryCaller: AgentRegistryCaller{contract: contract}, AgentRegistryTransactor: AgentRegistryTransactor{contract: contract}, AgentRegistryFilterer: AgentRegistryFilterer{contract: contract}}, nil
}

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

// AgentAdmins is a free data retrieval call binding the contract method 0x9dfeddae.
//
// Solidity: function agentAdmins(bytes32 , bytes32 , address ) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) AgentAdmins(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte, arg2 common.Address) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentAdmins", arg0, arg1, arg2)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AgentAdmins is a free data retrieval call binding the contract method 0x9dfeddae.
//
// Solidity: function agentAdmins(bytes32 , bytes32 , address ) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) AgentAdmins(arg0 [32]byte, arg1 [32]byte, arg2 common.Address) (bool, error) {
	return _AgentRegistry.Contract.AgentAdmins(&_AgentRegistry.CallOpts, arg0, arg1, arg2)
}

// AgentAdmins is a free data retrieval call binding the contract method 0x9dfeddae.
//
// Solidity: function agentAdmins(bytes32 , bytes32 , address ) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) AgentAdmins(arg0 [32]byte, arg1 [32]byte, arg2 common.Address) (bool, error) {
	return _AgentRegistry.Contract.AgentAdmins(&_AgentRegistry.CallOpts, arg0, arg1, arg2)
}

// AgentAt is a free data retrieval call binding the contract method 0x89d0a013.
//
// Solidity: function agentAt(bytes32 _poolId, uint256 index) view returns(bytes32, string)
func (_AgentRegistry *AgentRegistryCaller) AgentAt(opts *bind.CallOpts, _poolId [32]byte, index *big.Int) ([32]byte, string, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentAt", _poolId, index)

	if err != nil {
		return *new([32]byte), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// AgentAt is a free data retrieval call binding the contract method 0x89d0a013.
//
// Solidity: function agentAt(bytes32 _poolId, uint256 index) view returns(bytes32, string)
func (_AgentRegistry *AgentRegistrySession) AgentAt(_poolId [32]byte, index *big.Int) ([32]byte, string, error) {
	return _AgentRegistry.Contract.AgentAt(&_AgentRegistry.CallOpts, _poolId, index)
}

// AgentAt is a free data retrieval call binding the contract method 0x89d0a013.
//
// Solidity: function agentAt(bytes32 _poolId, uint256 index) view returns(bytes32, string)
func (_AgentRegistry *AgentRegistryCallerSession) AgentAt(_poolId [32]byte, index *big.Int) ([32]byte, string, error) {
	return _AgentRegistry.Contract.AgentAt(&_AgentRegistry.CallOpts, _poolId, index)
}

// AgentLength is a free data retrieval call binding the contract method 0x1e73feb1.
//
// Solidity: function agentLength(bytes32 _poolId) view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) AgentLength(opts *bind.CallOpts, _poolId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentLength", _poolId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentLength is a free data retrieval call binding the contract method 0x1e73feb1.
//
// Solidity: function agentLength(bytes32 _poolId) view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) AgentLength(_poolId [32]byte) (*big.Int, error) {
	return _AgentRegistry.Contract.AgentLength(&_AgentRegistry.CallOpts, _poolId)
}

// AgentLength is a free data retrieval call binding the contract method 0x1e73feb1.
//
// Solidity: function agentLength(bytes32 _poolId) view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) AgentLength(_poolId [32]byte) (*big.Int, error) {
	return _AgentRegistry.Contract.AgentLength(&_AgentRegistry.CallOpts, _poolId)
}

// PoolAdmins is a free data retrieval call binding the contract method 0x250ca0f0.
//
// Solidity: function poolAdmins(bytes32 , address ) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) PoolAdmins(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "poolAdmins", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PoolAdmins is a free data retrieval call binding the contract method 0x250ca0f0.
//
// Solidity: function poolAdmins(bytes32 , address ) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) PoolAdmins(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _AgentRegistry.Contract.PoolAdmins(&_AgentRegistry.CallOpts, arg0, arg1)
}

// PoolAdmins is a free data retrieval call binding the contract method 0x250ca0f0.
//
// Solidity: function poolAdmins(bytes32 , address ) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) PoolAdmins(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _AgentRegistry.Contract.PoolAdmins(&_AgentRegistry.CallOpts, arg0, arg1)
}

// PoolExistsMap is a free data retrieval call binding the contract method 0x0ead48e6.
//
// Solidity: function poolExistsMap(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) PoolExistsMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "poolExistsMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PoolExistsMap is a free data retrieval call binding the contract method 0x0ead48e6.
//
// Solidity: function poolExistsMap(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) PoolExistsMap(arg0 [32]byte) (bool, error) {
	return _AgentRegistry.Contract.PoolExistsMap(&_AgentRegistry.CallOpts, arg0)
}

// PoolExistsMap is a free data retrieval call binding the contract method 0x0ead48e6.
//
// Solidity: function poolExistsMap(bytes32 ) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) PoolExistsMap(arg0 [32]byte) (bool, error) {
	return _AgentRegistry.Contract.PoolExistsMap(&_AgentRegistry.CallOpts, arg0)
}

// PoolOwners is a free data retrieval call binding the contract method 0x219d2d09.
//
// Solidity: function poolOwners(bytes32 ) view returns(address)
func (_AgentRegistry *AgentRegistryCaller) PoolOwners(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "poolOwners", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolOwners is a free data retrieval call binding the contract method 0x219d2d09.
//
// Solidity: function poolOwners(bytes32 ) view returns(address)
func (_AgentRegistry *AgentRegistrySession) PoolOwners(arg0 [32]byte) (common.Address, error) {
	return _AgentRegistry.Contract.PoolOwners(&_AgentRegistry.CallOpts, arg0)
}

// PoolOwners is a free data retrieval call binding the contract method 0x219d2d09.
//
// Solidity: function poolOwners(bytes32 ) view returns(address)
func (_AgentRegistry *AgentRegistryCallerSession) PoolOwners(arg0 [32]byte) (common.Address, error) {
	return _AgentRegistry.Contract.PoolOwners(&_AgentRegistry.CallOpts, arg0)
}

// AddAgent is a paid mutator transaction binding the contract method 0x0944cd34.
//
// Solidity: function addAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactor) AddAgent(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "addAgent", _poolId, _agentId, _ref)
}

// AddAgent is a paid mutator transaction binding the contract method 0x0944cd34.
//
// Solidity: function addAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistrySession) AddAgent(_poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId, _ref)
}

// AddAgent is a paid mutator transaction binding the contract method 0x0944cd34.
//
// Solidity: function addAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) AddAgent(_poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId, _ref)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x63438a97.
//
// Solidity: function addAgentAdmin(bytes32 _poolId, bytes32 _agentId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactor) AddAgentAdmin(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "addAgentAdmin", _poolId, _agentId, admin)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x63438a97.
//
// Solidity: function addAgentAdmin(bytes32 _poolId, bytes32 _agentId, address admin) returns()
func (_AgentRegistry *AgentRegistrySession) AddAgentAdmin(_poolId [32]byte, _agentId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddAgentAdmin(&_AgentRegistry.TransactOpts, _poolId, _agentId, admin)
}

// AddAgentAdmin is a paid mutator transaction binding the contract method 0x63438a97.
//
// Solidity: function addAgentAdmin(bytes32 _poolId, bytes32 _agentId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) AddAgentAdmin(_poolId [32]byte, _agentId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddAgentAdmin(&_AgentRegistry.TransactOpts, _poolId, _agentId, admin)
}

// AddPool is a paid mutator transaction binding the contract method 0x8a0033fc.
//
// Solidity: function addPool(bytes32 _poolId) returns()
func (_AgentRegistry *AgentRegistryTransactor) AddPool(opts *bind.TransactOpts, _poolId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "addPool", _poolId)
}

// AddPool is a paid mutator transaction binding the contract method 0x8a0033fc.
//
// Solidity: function addPool(bytes32 _poolId) returns()
func (_AgentRegistry *AgentRegistrySession) AddPool(_poolId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddPool(&_AgentRegistry.TransactOpts, _poolId)
}

// AddPool is a paid mutator transaction binding the contract method 0x8a0033fc.
//
// Solidity: function addPool(bytes32 _poolId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) AddPool(_poolId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddPool(&_AgentRegistry.TransactOpts, _poolId)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x84930b76.
//
// Solidity: function addPoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactor) AddPoolAdmin(opts *bind.TransactOpts, _poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "addPoolAdmin", _poolId, admin)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x84930b76.
//
// Solidity: function addPoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistrySession) AddPoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddPoolAdmin(&_AgentRegistry.TransactOpts, _poolId, admin)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x84930b76.
//
// Solidity: function addPoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) AddPoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AddPoolAdmin(&_AgentRegistry.TransactOpts, _poolId, admin)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x61e62a0f.
//
// Solidity: function removeAgent(bytes32 _poolId, bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistryTransactor) RemoveAgent(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "removeAgent", _poolId, _agentId)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x61e62a0f.
//
// Solidity: function removeAgent(bytes32 _poolId, bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistrySession) RemoveAgent(_poolId [32]byte, _agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemoveAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x61e62a0f.
//
// Solidity: function removeAgent(bytes32 _poolId, bytes32 _agentId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) RemoveAgent(_poolId [32]byte, _agentId [32]byte) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemoveAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0xa7f2f5ab.
//
// Solidity: function removeAgentAdmin(bytes32 _poolId, bytes32 _agentId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactor) RemoveAgentAdmin(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "removeAgentAdmin", _poolId, _agentId, admin)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0xa7f2f5ab.
//
// Solidity: function removeAgentAdmin(bytes32 _poolId, bytes32 _agentId, address admin) returns()
func (_AgentRegistry *AgentRegistrySession) RemoveAgentAdmin(_poolId [32]byte, _agentId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemoveAgentAdmin(&_AgentRegistry.TransactOpts, _poolId, _agentId, admin)
}

// RemoveAgentAdmin is a paid mutator transaction binding the contract method 0xa7f2f5ab.
//
// Solidity: function removeAgentAdmin(bytes32 _poolId, bytes32 _agentId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) RemoveAgentAdmin(_poolId [32]byte, _agentId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemoveAgentAdmin(&_AgentRegistry.TransactOpts, _poolId, _agentId, admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xacce5630.
//
// Solidity: function removePoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactor) RemovePoolAdmin(opts *bind.TransactOpts, _poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "removePoolAdmin", _poolId, admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xacce5630.
//
// Solidity: function removePoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistrySession) RemovePoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemovePoolAdmin(&_AgentRegistry.TransactOpts, _poolId, admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xacce5630.
//
// Solidity: function removePoolAdmin(bytes32 _poolId, address admin) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) RemovePoolAdmin(_poolId [32]byte, admin common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RemovePoolAdmin(&_AgentRegistry.TransactOpts, _poolId, admin)
}

// TransferPoolOwnership is a paid mutator transaction binding the contract method 0x1b6feb92.
//
// Solidity: function transferPoolOwnership(bytes32 _poolId, address _to) returns()
func (_AgentRegistry *AgentRegistryTransactor) TransferPoolOwnership(opts *bind.TransactOpts, _poolId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "transferPoolOwnership", _poolId, _to)
}

// TransferPoolOwnership is a paid mutator transaction binding the contract method 0x1b6feb92.
//
// Solidity: function transferPoolOwnership(bytes32 _poolId, address _to) returns()
func (_AgentRegistry *AgentRegistrySession) TransferPoolOwnership(_poolId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.TransferPoolOwnership(&_AgentRegistry.TransactOpts, _poolId, _to)
}

// TransferPoolOwnership is a paid mutator transaction binding the contract method 0x1b6feb92.
//
// Solidity: function transferPoolOwnership(bytes32 _poolId, address _to) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) TransferPoolOwnership(_poolId [32]byte, _to common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.TransferPoolOwnership(&_AgentRegistry.TransactOpts, _poolId, _to)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xe0434e3f.
//
// Solidity: function updateAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactor) UpdateAgent(opts *bind.TransactOpts, _poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "updateAgent", _poolId, _agentId, _ref)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xe0434e3f.
//
// Solidity: function updateAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistrySession) UpdateAgent(_poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpdateAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId, _ref)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xe0434e3f.
//
// Solidity: function updateAgent(bytes32 _poolId, bytes32 _agentId, string _ref) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) UpdateAgent(_poolId [32]byte, _agentId [32]byte, _ref string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.UpdateAgent(&_AgentRegistry.TransactOpts, _poolId, _agentId, _ref)
}

// AgentRegistryAgentAddedIterator is returned from FilterAgentAdded and is used to iterate over the raw logs and unpacked data for AgentAdded events raised by the AgentRegistry contract.
type AgentRegistryAgentAddedIterator struct {
	Event *AgentRegistryAgentAdded // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentAdded)
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
		it.Event = new(AgentRegistryAgentAdded)
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
func (it *AgentRegistryAgentAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentAdded represents a AgentAdded event raised by the AgentRegistry contract.
type AgentRegistryAgentAdded struct {
	PoolId  [32]byte
	AgentId [32]byte
	Ref     string
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdded is a free log retrieval operation binding the contract event 0xdd1417d47ca73bdaf24e5c1df4158e61cbbb73fda6562b14a16215ccfa201aeb.
//
// Solidity: event AgentAdded(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentAdded(opts *bind.FilterOpts) (*AgentRegistryAgentAddedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentAdded")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentAddedIterator{contract: _AgentRegistry.contract, event: "AgentAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdded is a free log subscription operation binding the contract event 0xdd1417d47ca73bdaf24e5c1df4158e61cbbb73fda6562b14a16215ccfa201aeb.
//
// Solidity: event AgentAdded(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentAdded(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentAdded) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentAdded)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdded", log); err != nil {
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

// ParseAgentAdded is a log parse operation binding the contract event 0xdd1417d47ca73bdaf24e5c1df4158e61cbbb73fda6562b14a16215ccfa201aeb.
//
// Solidity: event AgentAdded(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentAdded(log types.Log) (*AgentRegistryAgentAdded, error) {
	event := new(AgentRegistryAgentAdded)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentAdminAddedIterator is returned from FilterAgentAdminAdded and is used to iterate over the raw logs and unpacked data for AgentAdminAdded events raised by the AgentRegistry contract.
type AgentRegistryAgentAdminAddedIterator struct {
	Event *AgentRegistryAgentAdminAdded // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentAdminAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentAdminAdded)
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
		it.Event = new(AgentRegistryAgentAdminAdded)
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
func (it *AgentRegistryAgentAdminAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentAdminAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentAdminAdded represents a AgentAdminAdded event raised by the AgentRegistry contract.
type AgentRegistryAgentAdminAdded struct {
	PoolId  [32]byte
	AgentId [32]byte
	Admin   common.Address
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdminAdded is a free log retrieval operation binding the contract event 0x11cbfe7067a564a8ddb9a780914987915bc0d29835c8b0554a9d9db2e51064ed.
//
// Solidity: event AgentAdminAdded(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentAdminAdded(opts *bind.FilterOpts) (*AgentRegistryAgentAdminAddedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentAdminAdded")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentAdminAddedIterator{contract: _AgentRegistry.contract, event: "AgentAdminAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdminAdded is a free log subscription operation binding the contract event 0x11cbfe7067a564a8ddb9a780914987915bc0d29835c8b0554a9d9db2e51064ed.
//
// Solidity: event AgentAdminAdded(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentAdminAdded(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentAdminAdded) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentAdminAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentAdminAdded)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdminAdded", log); err != nil {
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

// ParseAgentAdminAdded is a log parse operation binding the contract event 0x11cbfe7067a564a8ddb9a780914987915bc0d29835c8b0554a9d9db2e51064ed.
//
// Solidity: event AgentAdminAdded(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentAdminAdded(log types.Log) (*AgentRegistryAgentAdminAdded, error) {
	event := new(AgentRegistryAgentAdminAdded)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdminAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentAdminRemovedIterator is returned from FilterAgentAdminRemoved and is used to iterate over the raw logs and unpacked data for AgentAdminRemoved events raised by the AgentRegistry contract.
type AgentRegistryAgentAdminRemovedIterator struct {
	Event *AgentRegistryAgentAdminRemoved // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentAdminRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentAdminRemoved)
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
		it.Event = new(AgentRegistryAgentAdminRemoved)
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
func (it *AgentRegistryAgentAdminRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentAdminRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentAdminRemoved represents a AgentAdminRemoved event raised by the AgentRegistry contract.
type AgentRegistryAgentAdminRemoved struct {
	PoolId  [32]byte
	AgentId [32]byte
	Admin   common.Address
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdminRemoved is a free log retrieval operation binding the contract event 0x801fe4d9d51c0445c2938970ad5e042c3f2f2c387d67de4111af1bc60d22dba9.
//
// Solidity: event AgentAdminRemoved(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentAdminRemoved(opts *bind.FilterOpts) (*AgentRegistryAgentAdminRemovedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentAdminRemoved")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentAdminRemovedIterator{contract: _AgentRegistry.contract, event: "AgentAdminRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentAdminRemoved is a free log subscription operation binding the contract event 0x801fe4d9d51c0445c2938970ad5e042c3f2f2c387d67de4111af1bc60d22dba9.
//
// Solidity: event AgentAdminRemoved(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentAdminRemoved(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentAdminRemoved) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentAdminRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentAdminRemoved)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdminRemoved", log); err != nil {
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

// ParseAgentAdminRemoved is a log parse operation binding the contract event 0x801fe4d9d51c0445c2938970ad5e042c3f2f2c387d67de4111af1bc60d22dba9.
//
// Solidity: event AgentAdminRemoved(bytes32 poolId, bytes32 agentId, address admin, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentAdminRemoved(log types.Log) (*AgentRegistryAgentAdminRemoved, error) {
	event := new(AgentRegistryAgentAdminRemoved)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentAdminRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentRemovedIterator is returned from FilterAgentRemoved and is used to iterate over the raw logs and unpacked data for AgentRemoved events raised by the AgentRegistry contract.
type AgentRegistryAgentRemovedIterator struct {
	Event *AgentRegistryAgentRemoved // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentRemoved)
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
		it.Event = new(AgentRegistryAgentRemoved)
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
func (it *AgentRegistryAgentRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentRemoved represents a AgentRemoved event raised by the AgentRegistry contract.
type AgentRegistryAgentRemoved struct {
	PoolId  [32]byte
	AgentId [32]byte
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentRemoved is a free log retrieval operation binding the contract event 0xccc61238b7da918476155ec78f1d6f23cd587a3f796774b80c0c79dde8f57e6b.
//
// Solidity: event AgentRemoved(bytes32 poolId, bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentRemoved(opts *bind.FilterOpts) (*AgentRegistryAgentRemovedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentRemoved")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentRemovedIterator{contract: _AgentRegistry.contract, event: "AgentRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentRemoved is a free log subscription operation binding the contract event 0xccc61238b7da918476155ec78f1d6f23cd587a3f796774b80c0c79dde8f57e6b.
//
// Solidity: event AgentRemoved(bytes32 poolId, bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentRemoved(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentRemoved) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentRemoved)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
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

// ParseAgentRemoved is a log parse operation binding the contract event 0xccc61238b7da918476155ec78f1d6f23cd587a3f796774b80c0c79dde8f57e6b.
//
// Solidity: event AgentRemoved(bytes32 poolId, bytes32 agentId, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentRemoved(log types.Log) (*AgentRegistryAgentRemoved, error) {
	event := new(AgentRegistryAgentRemoved)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
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
	PoolId  [32]byte
	AgentId [32]byte
	Ref     string
	By      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentUpdated is a free log retrieval operation binding the contract event 0x8a33912f31f80df098dcad04aba17afc8c52d0f700baf2750ea3899f9230a850.
//
// Solidity: event AgentUpdated(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentUpdated(opts *bind.FilterOpts) (*AgentRegistryAgentUpdatedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentUpdated")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentUpdatedIterator{contract: _AgentRegistry.contract, event: "AgentUpdated", logs: logs, sub: sub}, nil
}

// WatchAgentUpdated is a free log subscription operation binding the contract event 0x8a33912f31f80df098dcad04aba17afc8c52d0f700baf2750ea3899f9230a850.
//
// Solidity: event AgentUpdated(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentUpdated(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentUpdated) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentUpdated")
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

// ParseAgentUpdated is a log parse operation binding the contract event 0x8a33912f31f80df098dcad04aba17afc8c52d0f700baf2750ea3899f9230a850.
//
// Solidity: event AgentUpdated(bytes32 poolId, bytes32 agentId, string ref, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentUpdated(log types.Log) (*AgentRegistryAgentUpdated, error) {
	event := new(AgentRegistryAgentUpdated)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryPoolAddedIterator is returned from FilterPoolAdded and is used to iterate over the raw logs and unpacked data for PoolAdded events raised by the AgentRegistry contract.
type AgentRegistryPoolAddedIterator struct {
	Event *AgentRegistryPoolAdded // Event containing the contract specifics and raw log

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
func (it *AgentRegistryPoolAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryPoolAdded)
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
		it.Event = new(AgentRegistryPoolAdded)
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
func (it *AgentRegistryPoolAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryPoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryPoolAdded represents a PoolAdded event raised by the AgentRegistry contract.
type AgentRegistryPoolAdded struct {
	PoolId [32]byte
	By     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolAdded is a free log retrieval operation binding the contract event 0xfd0fa7919fbe3857a4236750e8d3e42ac691881d2f31c50ebe4cfc2a7705ee80.
//
// Solidity: event PoolAdded(bytes32 poolId, address by)
func (_AgentRegistry *AgentRegistryFilterer) FilterPoolAdded(opts *bind.FilterOpts) (*AgentRegistryPoolAddedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "PoolAdded")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryPoolAddedIterator{contract: _AgentRegistry.contract, event: "PoolAdded", logs: logs, sub: sub}, nil
}

// WatchPoolAdded is a free log subscription operation binding the contract event 0xfd0fa7919fbe3857a4236750e8d3e42ac691881d2f31c50ebe4cfc2a7705ee80.
//
// Solidity: event PoolAdded(bytes32 poolId, address by)
func (_AgentRegistry *AgentRegistryFilterer) WatchPoolAdded(opts *bind.WatchOpts, sink chan<- *AgentRegistryPoolAdded) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "PoolAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryPoolAdded)
				if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdded", log); err != nil {
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

// ParsePoolAdded is a log parse operation binding the contract event 0xfd0fa7919fbe3857a4236750e8d3e42ac691881d2f31c50ebe4cfc2a7705ee80.
//
// Solidity: event PoolAdded(bytes32 poolId, address by)
func (_AgentRegistry *AgentRegistryFilterer) ParsePoolAdded(log types.Log) (*AgentRegistryPoolAdded, error) {
	event := new(AgentRegistryPoolAdded)
	if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryPoolAdminAddedIterator is returned from FilterPoolAdminAdded and is used to iterate over the raw logs and unpacked data for PoolAdminAdded events raised by the AgentRegistry contract.
type AgentRegistryPoolAdminAddedIterator struct {
	Event *AgentRegistryPoolAdminAdded // Event containing the contract specifics and raw log

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
func (it *AgentRegistryPoolAdminAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryPoolAdminAdded)
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
		it.Event = new(AgentRegistryPoolAdminAdded)
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
func (it *AgentRegistryPoolAdminAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryPoolAdminAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryPoolAdminAdded represents a PoolAdminAdded event raised by the AgentRegistry contract.
type AgentRegistryPoolAdminAdded struct {
	PoolId [32]byte
	Addr   common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolAdminAdded is a free log retrieval operation binding the contract event 0x4c2f0be874d7047c389e9fadf454854a8f5b72b49de199c7d60831b6c462a9a0.
//
// Solidity: event PoolAdminAdded(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) FilterPoolAdminAdded(opts *bind.FilterOpts) (*AgentRegistryPoolAdminAddedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "PoolAdminAdded")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryPoolAdminAddedIterator{contract: _AgentRegistry.contract, event: "PoolAdminAdded", logs: logs, sub: sub}, nil
}

// WatchPoolAdminAdded is a free log subscription operation binding the contract event 0x4c2f0be874d7047c389e9fadf454854a8f5b72b49de199c7d60831b6c462a9a0.
//
// Solidity: event PoolAdminAdded(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) WatchPoolAdminAdded(opts *bind.WatchOpts, sink chan<- *AgentRegistryPoolAdminAdded) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "PoolAdminAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryPoolAdminAdded)
				if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdminAdded", log); err != nil {
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

// ParsePoolAdminAdded is a log parse operation binding the contract event 0x4c2f0be874d7047c389e9fadf454854a8f5b72b49de199c7d60831b6c462a9a0.
//
// Solidity: event PoolAdminAdded(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) ParsePoolAdminAdded(log types.Log) (*AgentRegistryPoolAdminAdded, error) {
	event := new(AgentRegistryPoolAdminAdded)
	if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdminAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryPoolAdminRemovedIterator is returned from FilterPoolAdminRemoved and is used to iterate over the raw logs and unpacked data for PoolAdminRemoved events raised by the AgentRegistry contract.
type AgentRegistryPoolAdminRemovedIterator struct {
	Event *AgentRegistryPoolAdminRemoved // Event containing the contract specifics and raw log

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
func (it *AgentRegistryPoolAdminRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryPoolAdminRemoved)
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
		it.Event = new(AgentRegistryPoolAdminRemoved)
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
func (it *AgentRegistryPoolAdminRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryPoolAdminRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryPoolAdminRemoved represents a PoolAdminRemoved event raised by the AgentRegistry contract.
type AgentRegistryPoolAdminRemoved struct {
	PoolId [32]byte
	Addr   common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolAdminRemoved is a free log retrieval operation binding the contract event 0x25fcb4cb225521d9a561fd90af5809d997ae0fa769a5c95d90750162623bf624.
//
// Solidity: event PoolAdminRemoved(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) FilterPoolAdminRemoved(opts *bind.FilterOpts) (*AgentRegistryPoolAdminRemovedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "PoolAdminRemoved")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryPoolAdminRemovedIterator{contract: _AgentRegistry.contract, event: "PoolAdminRemoved", logs: logs, sub: sub}, nil
}

// WatchPoolAdminRemoved is a free log subscription operation binding the contract event 0x25fcb4cb225521d9a561fd90af5809d997ae0fa769a5c95d90750162623bf624.
//
// Solidity: event PoolAdminRemoved(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) WatchPoolAdminRemoved(opts *bind.WatchOpts, sink chan<- *AgentRegistryPoolAdminRemoved) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "PoolAdminRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryPoolAdminRemoved)
				if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdminRemoved", log); err != nil {
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

// ParsePoolAdminRemoved is a log parse operation binding the contract event 0x25fcb4cb225521d9a561fd90af5809d997ae0fa769a5c95d90750162623bf624.
//
// Solidity: event PoolAdminRemoved(bytes32 poolId, address addr)
func (_AgentRegistry *AgentRegistryFilterer) ParsePoolAdminRemoved(log types.Log) (*AgentRegistryPoolAdminRemoved, error) {
	event := new(AgentRegistryPoolAdminRemoved)
	if err := _AgentRegistry.contract.UnpackLog(event, "PoolAdminRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryPoolOwnershipTransferedIterator is returned from FilterPoolOwnershipTransfered and is used to iterate over the raw logs and unpacked data for PoolOwnershipTransfered events raised by the AgentRegistry contract.
type AgentRegistryPoolOwnershipTransferedIterator struct {
	Event *AgentRegistryPoolOwnershipTransfered // Event containing the contract specifics and raw log

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
func (it *AgentRegistryPoolOwnershipTransferedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryPoolOwnershipTransfered)
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
		it.Event = new(AgentRegistryPoolOwnershipTransfered)
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
func (it *AgentRegistryPoolOwnershipTransferedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryPoolOwnershipTransferedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryPoolOwnershipTransfered represents a PoolOwnershipTransfered event raised by the AgentRegistry contract.
type AgentRegistryPoolOwnershipTransfered struct {
	PoolId [32]byte
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPoolOwnershipTransfered is a free log retrieval operation binding the contract event 0xecc1c418cc7d2c3062b60c1d8f0b7748e2c8c10a109df4f6089b285f26b993d6.
//
// Solidity: event PoolOwnershipTransfered(bytes32 poolId, address from, address to)
func (_AgentRegistry *AgentRegistryFilterer) FilterPoolOwnershipTransfered(opts *bind.FilterOpts) (*AgentRegistryPoolOwnershipTransferedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "PoolOwnershipTransfered")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryPoolOwnershipTransferedIterator{contract: _AgentRegistry.contract, event: "PoolOwnershipTransfered", logs: logs, sub: sub}, nil
}

// WatchPoolOwnershipTransfered is a free log subscription operation binding the contract event 0xecc1c418cc7d2c3062b60c1d8f0b7748e2c8c10a109df4f6089b285f26b993d6.
//
// Solidity: event PoolOwnershipTransfered(bytes32 poolId, address from, address to)
func (_AgentRegistry *AgentRegistryFilterer) WatchPoolOwnershipTransfered(opts *bind.WatchOpts, sink chan<- *AgentRegistryPoolOwnershipTransfered) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "PoolOwnershipTransfered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryPoolOwnershipTransfered)
				if err := _AgentRegistry.contract.UnpackLog(event, "PoolOwnershipTransfered", log); err != nil {
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

// ParsePoolOwnershipTransfered is a log parse operation binding the contract event 0xecc1c418cc7d2c3062b60c1d8f0b7748e2c8c10a109df4f6089b285f26b993d6.
//
// Solidity: event PoolOwnershipTransfered(bytes32 poolId, address from, address to)
func (_AgentRegistry *AgentRegistryFilterer) ParsePoolOwnershipTransfered(log types.Log) (*AgentRegistryPoolOwnershipTransfered, error) {
	event := new(AgentRegistryPoolOwnershipTransfered)
	if err := _AgentRegistry.contract.UnpackLog(event, "PoolOwnershipTransfered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EnumerableSetUpgradeableABI is the input ABI used to generate the binding from.
const EnumerableSetUpgradeableABI = "[]"

// EnumerableSetUpgradeableBin is the compiled bytecode used for deploying new contracts.
var EnumerableSetUpgradeableBin = "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212205aad5bc4264d189a21bec035629c731fcea42df7f3f53a3822d4795748b2c02764736f6c63430008040033"

// DeployEnumerableSetUpgradeable deploys a new Ethereum contract, binding an instance of EnumerableSetUpgradeable to it.
func DeployEnumerableSetUpgradeable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EnumerableSetUpgradeable, error) {
	parsed, err := abi.JSON(strings.NewReader(EnumerableSetUpgradeableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EnumerableSetUpgradeableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EnumerableSetUpgradeable{EnumerableSetUpgradeableCaller: EnumerableSetUpgradeableCaller{contract: contract}, EnumerableSetUpgradeableTransactor: EnumerableSetUpgradeableTransactor{contract: contract}, EnumerableSetUpgradeableFilterer: EnumerableSetUpgradeableFilterer{contract: contract}}, nil
}

// EnumerableSetUpgradeable is an auto generated Go binding around an Ethereum contract.
type EnumerableSetUpgradeable struct {
	EnumerableSetUpgradeableCaller     // Read-only binding to the contract
	EnumerableSetUpgradeableTransactor // Write-only binding to the contract
	EnumerableSetUpgradeableFilterer   // Log filterer for contract events
}

// EnumerableSetUpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type EnumerableSetUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableSetUpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EnumerableSetUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableSetUpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EnumerableSetUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableSetUpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EnumerableSetUpgradeableSession struct {
	Contract     *EnumerableSetUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// EnumerableSetUpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EnumerableSetUpgradeableCallerSession struct {
	Contract *EnumerableSetUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// EnumerableSetUpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EnumerableSetUpgradeableTransactorSession struct {
	Contract     *EnumerableSetUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// EnumerableSetUpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type EnumerableSetUpgradeableRaw struct {
	Contract *EnumerableSetUpgradeable // Generic contract binding to access the raw methods on
}

// EnumerableSetUpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EnumerableSetUpgradeableCallerRaw struct {
	Contract *EnumerableSetUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// EnumerableSetUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EnumerableSetUpgradeableTransactorRaw struct {
	Contract *EnumerableSetUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEnumerableSetUpgradeable creates a new instance of EnumerableSetUpgradeable, bound to a specific deployed contract.
func NewEnumerableSetUpgradeable(address common.Address, backend bind.ContractBackend) (*EnumerableSetUpgradeable, error) {
	contract, err := bindEnumerableSetUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EnumerableSetUpgradeable{EnumerableSetUpgradeableCaller: EnumerableSetUpgradeableCaller{contract: contract}, EnumerableSetUpgradeableTransactor: EnumerableSetUpgradeableTransactor{contract: contract}, EnumerableSetUpgradeableFilterer: EnumerableSetUpgradeableFilterer{contract: contract}}, nil
}

// NewEnumerableSetUpgradeableCaller creates a new read-only instance of EnumerableSetUpgradeable, bound to a specific deployed contract.
func NewEnumerableSetUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*EnumerableSetUpgradeableCaller, error) {
	contract, err := bindEnumerableSetUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EnumerableSetUpgradeableCaller{contract: contract}, nil
}

// NewEnumerableSetUpgradeableTransactor creates a new write-only instance of EnumerableSetUpgradeable, bound to a specific deployed contract.
func NewEnumerableSetUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*EnumerableSetUpgradeableTransactor, error) {
	contract, err := bindEnumerableSetUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EnumerableSetUpgradeableTransactor{contract: contract}, nil
}

// NewEnumerableSetUpgradeableFilterer creates a new log filterer instance of EnumerableSetUpgradeable, bound to a specific deployed contract.
func NewEnumerableSetUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*EnumerableSetUpgradeableFilterer, error) {
	contract, err := bindEnumerableSetUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EnumerableSetUpgradeableFilterer{contract: contract}, nil
}

// bindEnumerableSetUpgradeable binds a generic wrapper to an already deployed contract.
func bindEnumerableSetUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EnumerableSetUpgradeableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnumerableSetUpgradeable *EnumerableSetUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnumerableSetUpgradeable.Contract.EnumerableSetUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnumerableSetUpgradeable *EnumerableSetUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnumerableSetUpgradeable.Contract.EnumerableSetUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnumerableSetUpgradeable *EnumerableSetUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnumerableSetUpgradeable.Contract.EnumerableSetUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnumerableSetUpgradeable *EnumerableSetUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnumerableSetUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnumerableSetUpgradeable *EnumerableSetUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnumerableSetUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnumerableSetUpgradeable *EnumerableSetUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnumerableSetUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// EnumerableStringMapABI is the input ABI used to generate the binding from.
const EnumerableStringMapABI = "[]"

// EnumerableStringMapBin is the compiled bytecode used for deploying new contracts.
var EnumerableStringMapBin = "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220c5db38e644ce0ebae64ba2a93a7fc4968b975a27a70e9b7a40434d75507bff3c64736f6c63430008040033"

// DeployEnumerableStringMap deploys a new Ethereum contract, binding an instance of EnumerableStringMap to it.
func DeployEnumerableStringMap(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EnumerableStringMap, error) {
	parsed, err := abi.JSON(strings.NewReader(EnumerableStringMapABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EnumerableStringMapBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EnumerableStringMap{EnumerableStringMapCaller: EnumerableStringMapCaller{contract: contract}, EnumerableStringMapTransactor: EnumerableStringMapTransactor{contract: contract}, EnumerableStringMapFilterer: EnumerableStringMapFilterer{contract: contract}}, nil
}

// EnumerableStringMap is an auto generated Go binding around an Ethereum contract.
type EnumerableStringMap struct {
	EnumerableStringMapCaller     // Read-only binding to the contract
	EnumerableStringMapTransactor // Write-only binding to the contract
	EnumerableStringMapFilterer   // Log filterer for contract events
}

// EnumerableStringMapCaller is an auto generated read-only Go binding around an Ethereum contract.
type EnumerableStringMapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableStringMapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EnumerableStringMapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableStringMapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EnumerableStringMapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableStringMapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EnumerableStringMapSession struct {
	Contract     *EnumerableStringMap // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// EnumerableStringMapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EnumerableStringMapCallerSession struct {
	Contract *EnumerableStringMapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// EnumerableStringMapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EnumerableStringMapTransactorSession struct {
	Contract     *EnumerableStringMapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// EnumerableStringMapRaw is an auto generated low-level Go binding around an Ethereum contract.
type EnumerableStringMapRaw struct {
	Contract *EnumerableStringMap // Generic contract binding to access the raw methods on
}

// EnumerableStringMapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EnumerableStringMapCallerRaw struct {
	Contract *EnumerableStringMapCaller // Generic read-only contract binding to access the raw methods on
}

// EnumerableStringMapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EnumerableStringMapTransactorRaw struct {
	Contract *EnumerableStringMapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEnumerableStringMap creates a new instance of EnumerableStringMap, bound to a specific deployed contract.
func NewEnumerableStringMap(address common.Address, backend bind.ContractBackend) (*EnumerableStringMap, error) {
	contract, err := bindEnumerableStringMap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EnumerableStringMap{EnumerableStringMapCaller: EnumerableStringMapCaller{contract: contract}, EnumerableStringMapTransactor: EnumerableStringMapTransactor{contract: contract}, EnumerableStringMapFilterer: EnumerableStringMapFilterer{contract: contract}}, nil
}

// NewEnumerableStringMapCaller creates a new read-only instance of EnumerableStringMap, bound to a specific deployed contract.
func NewEnumerableStringMapCaller(address common.Address, caller bind.ContractCaller) (*EnumerableStringMapCaller, error) {
	contract, err := bindEnumerableStringMap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EnumerableStringMapCaller{contract: contract}, nil
}

// NewEnumerableStringMapTransactor creates a new write-only instance of EnumerableStringMap, bound to a specific deployed contract.
func NewEnumerableStringMapTransactor(address common.Address, transactor bind.ContractTransactor) (*EnumerableStringMapTransactor, error) {
	contract, err := bindEnumerableStringMap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EnumerableStringMapTransactor{contract: contract}, nil
}

// NewEnumerableStringMapFilterer creates a new log filterer instance of EnumerableStringMap, bound to a specific deployed contract.
func NewEnumerableStringMapFilterer(address common.Address, filterer bind.ContractFilterer) (*EnumerableStringMapFilterer, error) {
	contract, err := bindEnumerableStringMap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EnumerableStringMapFilterer{contract: contract}, nil
}

// bindEnumerableStringMap binds a generic wrapper to an already deployed contract.
func bindEnumerableStringMap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EnumerableStringMapABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnumerableStringMap *EnumerableStringMapRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnumerableStringMap.Contract.EnumerableStringMapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnumerableStringMap *EnumerableStringMapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnumerableStringMap.Contract.EnumerableStringMapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnumerableStringMap *EnumerableStringMapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnumerableStringMap.Contract.EnumerableStringMapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnumerableStringMap *EnumerableStringMapCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnumerableStringMap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnumerableStringMap *EnumerableStringMapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnumerableStringMap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnumerableStringMap *EnumerableStringMapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnumerableStringMap.Contract.contract.Transact(opts, method, params...)
}
