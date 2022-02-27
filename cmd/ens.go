package cmd

import (
	"github.com/forta-protocol/forta-core-go/ens"
	"github.com/goccy/go-json"
	"io/ioutil"
	"path"
	"time"
)

const (
	contractAddressCacheExpiry = time.Hour
)

// useEnsDefaults gets and uses ENS defaults if needed.
func useEnsDefaults() error {
	if cfg.Registry.ContractAddress != "" {
		return nil
	}

	return ensureLatestContractAddresses()
}

// useEnsAgentReg finds the agent registry from a contract.
func useEnsAgentReg() error {
	return ensureLatestContractAddresses()
}

func ensureLatestContractAddresses() error {
	now := time.Now().UTC()

	cache, ok := getContractAddressCache()
	if ok && now.Before(cache.ExpiresAt) {
		setContractAddressesFromCache(cache)
		return nil
	}

	whiteBold("Refreshing contract address cache...\n")

	if cfg.ENSConfig.DefaultContract {
		cfg.ENSConfig.ContractAddress = ""
	}
	es, err := ens.DialENSStoreAt(cfg.ENSConfig.JsonRpc.Url, cfg.ENSConfig.ContractAddress)
	if err != nil {
		return err
	}

	contracts, err := es.ResolveRegistryContracts()
	if err != nil {
		return err
	}

	cache.Dispatch = contracts.Dispatch.Hex()
	cache.Agents = contracts.AgentRegistry.Hex()
	cache.ScannerVersion = contracts.ScannerRegistry.Hex()
	cache.ExpiresAt = time.Now().UTC().Add(contractAddressCacheExpiry)

	b, err := json.MarshalIndent(&cache, "", "  ") // indent by two spaces
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(cfg.FortaDir, "contracts.json"), b, 0644); err != nil {
		return err
	}

	setContractAddressesFromCache(cache)
	return nil
}

// sets only if not overridden
func setContractAddressesFromCache(cache contractAddressCache) {
	if cfg.Registry.ContractAddress == "" {
		cfg.Registry.ContractAddress = cache.Dispatch
	}
	cfg.AgentRegistryContractAddress = cache.Agents
}

type contractAddressCache struct {
	Dispatch       string    `json:"dispatch"`
	Agents         string    `json:"agents"`
	ScannerVersion string    `json:"scannerVersion"`
	ExpiresAt      time.Time `json:"expiresAt"`
}

func getContractAddressCache() (cache contractAddressCache, ok bool) {
	b, err := ioutil.ReadFile(path.Join(cfg.FortaDir, "contracts.json"))
	if err != nil {
		return
	}

	if err := json.Unmarshal(b, &cache); err != nil {
		return
	}

	ok = true
	return
}
