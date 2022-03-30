package store

import (
	"io/ioutil"
	"path"

	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-protocol/forta-core-go/domain/registry"
	"github.com/forta-protocol/forta-core-go/ens"
	"github.com/forta-protocol/forta-node/config"
	"github.com/goccy/go-json"
)

type ensOverrideStore struct {
	contracts    registry.RegistryContracts
	contractsMap map[string]string
}

func NewENSOverrideStore(cfg config.Config) (*ensOverrideStore, error) {
	var store ensOverrideStore
	b, err := ioutil.ReadFile(path.Join(cfg.FortaDir, "ens-override.json"))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &store.contractsMap); err != nil {
		return nil, err
	}
	store.contracts.Dispatch = common.HexToAddress(store.contractsMap[ens.DispatchContract])
	store.contracts.AgentRegistry = common.HexToAddress(store.contractsMap[ens.AgentRegistryContract])
	store.contracts.ScannerRegistry = common.HexToAddress(store.contractsMap[ens.ScannerRegistryContract])
	store.contracts.ScannerNodeVersion = common.HexToAddress(store.contractsMap[ens.ScannerNodeVersionContract])
	store.contracts.FortaStaking = common.HexToAddress(store.contractsMap[ens.StakingContract])
	return &store, nil
}

func (store *ensOverrideStore) Resolve(input string) (common.Address, error) {
	return common.HexToAddress(store.contractsMap[input]), nil
}

func (store *ensOverrideStore) ResolveRegistryContracts() (*registry.RegistryContracts, error) {
	return &store.contracts, nil
}
