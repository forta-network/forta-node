package store

import (
	"io/ioutil"
	"path"

	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-network/forta-node/config"
	"github.com/goccy/go-json"
)

type ensOverrideStore struct {
	contractsMap map[string]string
}

func NewENSOverrideResolver(cfg config.Config) (*ensOverrideStore, error) {
	var store ensOverrideStore
	b, err := ioutil.ReadFile(path.Join(cfg.FortaDir, "ens-override.json"))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &store.contractsMap); err != nil {
		return nil, err
	}
	return &store, nil
}

func (store *ensOverrideStore) Resolve(input string) (common.Address, error) {
	return common.HexToAddress(store.contractsMap[input]), nil
}
