package store

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/wealdtech/go-ens/v3"
)

// ENS resolves inputs.
type ENS interface {
	Resolve(input string) (common.Address, error)
}

// ENSStore wraps the ENS client which interacts with namespace contract(s).
type ENSStore struct {
	backend bind.ContractBackend
}

// NewENSStore creates a new store.
func NewENSStore(backend bind.ContractBackend) ENS {
	return &ENSStore{backend: backend}
}

// DialENSStore dials an Ethereum API and creates a new store.
func DialENSStore(rpcUrl string) (ENS, error) {
	client, err := rpc.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	return &ENSStore{backend: ethclient.NewClient(client)}, nil
}

// Resolve resolves an input to an address.
func (ensstore *ENSStore) Resolve(input string) (common.Address, error) {
	return ens.Resolve(ensstore.backend, input)
}
