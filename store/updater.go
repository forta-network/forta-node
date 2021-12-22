package store

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/contracts"
	"github.com/forta-protocol/forta-node/ethereum"
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

	snvc, err := contracts.NewScannerNodeVersionCaller(common.HexToAddress(cfg.ScannerVersionContractAddress), client)
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
