package store

import (
	"bytes"
	"errors"

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
	backend      bind.ContractBackend
	resolverAddr string
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

// DialENSStoreAt dials an Ethereum API and creates a new store that works with a resolver at given address.
func DialENSStoreAt(rpcUrl, resolverAddr string) (ENS, error) {
	client, err := rpc.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	return &ENSStore{backend: ethclient.NewClient(client), resolverAddr: resolverAddr}, nil
}

// Resolve resolves an input to an address.
func (ensstore *ENSStore) Resolve(input string) (common.Address, error) {
	if len(ensstore.resolverAddr) == 0 {
		return ens.Resolve(ensstore.backend, input)
	}
	resolver, err := ens.NewResolverAt(ensstore.backend, input, common.HexToAddress(ensstore.resolverAddr))
	if err != nil {
		return common.Address{}, err
	}
	// Resolve the domain
	address, err := resolver.Address()
	if err != nil {
		return ens.UnknownAddress, err
	}
	if bytes.Equal(address.Bytes(), ens.UnknownAddress.Bytes()) {
		return ens.UnknownAddress, errors.New("no address")
	}
	return address, nil
}
