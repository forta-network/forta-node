package ens

import (
	"fmt"
	"github.com/forta-protocol/forta-node/store"
	log "github.com/sirupsen/logrus"
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
		Agents:         "agents.registries.forta.eth",
		ScannerVersion: "scanner-node-version.forta.eth",
	}
}

func ResolveFortaContracts(jsonRpcUrl, resolverAddr string) (*FortaContracts, error) {
	ens, err := store.DialENSStoreAt(jsonRpcUrl, resolverAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve ens contract addresses: %v", err)
	}

	names := GetENSNames()
	dispatch, err := ens.Resolve(names.Dispatch)
	if err != nil {
		log.WithFields(log.Fields{
			"address": names.Dispatch,
		}).WithError(err).Error("ens cannot resolve dispatch contract")
		return nil, err
	}

	alerts, err := ens.Resolve(names.Alerts)
	if err != nil {
		log.WithFields(log.Fields{
			"address": names.Alerts,
		}).WithError(err).Error("ens cannot resolve alerts contract")
		return nil, err
	}

	agents, err := ens.Resolve(names.Agents)
	if err != nil {
		log.WithFields(log.Fields{
			"address": names.Agents,
		}).WithError(err).Error("ens cannot resolve agents contract")
		return nil, err
	}

	snv, err := ens.Resolve(names.ScannerVersion)
	if err != nil {
		log.WithFields(log.Fields{
			"address": names.ScannerVersion,
		}).WithError(err).Error("ens cannot resolve scanner version contract")
		return nil, err
	}

	return &FortaContracts{
		Dispatch:       dispatch.Hex(),
		Agent:          agents.Hex(),
		Alerts:         alerts.Hex(),
		ScannerVersion: snv.Hex(),
	}, nil
}
