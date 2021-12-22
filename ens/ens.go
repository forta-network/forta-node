package ens

import (
	"fmt"
	"github.com/forta-protocol/forta-node/store"
)

type FortaContracts struct {
	Dispatch       string
	Agent          string
	Alerts         string
	ScannerVersion string
}

// ENS contains the default names.
type ENS struct {
	Dispatch       string
	Alerts         string
	Agents         string
	ScannerVersion string
}

// GetENSNames returns the default ENS names.
func GetENSNames() *ENS {
	return &ENS{
		Dispatch:       "dispatch.forta.eth",
		Alerts:         "alerts.forta.eth",
		Agents:         "agents.registries.forta.eth",
		ScannerVersion: "scanner-node-version.forta.eth",
	}
}

func ResolveFortaContracts(jsonRpcUrl, resolverAddr string) (*FortaContracts, error) {
	ens, err := store.DialENSStoreAt(jsonRpcUrl, resolverAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve contract addresses from ENS: %v", err)
	}

	names := GetENSNames()
	dispatch, err := ens.Resolve(names.Dispatch)
	if err != nil {
		return nil, err
	}

	alerts, err := ens.Resolve(names.Alerts)
	if err != nil {
		return nil, err
	}

	agents, err := ens.Resolve(names.Agents)
	if err != nil {
		return nil, err
	}

	snv, err := ens.Resolve(names.ScannerVersion)
	if err != nil {
		return nil, err
	}

	return &FortaContracts{
		Dispatch:       dispatch.Hex(),
		Agent:          agents.Hex(),
		Alerts:         alerts.Hex(),
		ScannerVersion: snv.Hex(),
	}, nil
}
