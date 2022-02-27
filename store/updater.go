package store

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/forta-protocol/forta-core-go/contracts/contract_scanner_node_version"
	"github.com/forta-protocol/forta-core-go/ethereum"
	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/config"
)

type UpdaterStore interface {
	GetLatestReference() (string, error)
}

type contractUpdaterStore struct {
	versionContract scannerNodeVersion
}

func NewContractUpdaterStore(cfg config.Config) (*contractUpdaterStore, error) {
	rpc, err := ethereum.NewRpcClient(cfg.Registry.JsonRpc.Url)
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpc)

	log.WithField("address", cfg.ScannerVersionContractAddress).Info("attaching to scanner version contract")
	snvc, err := contract_scanner_node_version.NewScannerNodeVersionCaller(common.HexToAddress(cfg.ScannerVersionContractAddress), client)
	if err != nil {
		return nil, err
	}
	return &contractUpdaterStore{
		versionContract: snvc,
	}, nil
}

func (cus *contractUpdaterStore) GetLatestReference() (string, error) {
	return cus.versionContract.ScannerNodeVersion(nil)
}
