package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/forta-protocol/forta-core-go/contracts/contract_scanner_registry"
	"github.com/forta-protocol/forta-core-go/security"
	"github.com/forta-protocol/forta-core-go/utils"
	"github.com/forta-protocol/forta-node/cmd/runner"
	"github.com/spf13/cobra"
)

func handleFortaRun(cmd *cobra.Command, args []string) error {
	if err := checkScannerState(); err != nil {
		return err
	}
	runner.Run(cfg)
	return nil
}

func checkScannerState() error {
	if parsedArgs.NoCheck {
		return nil
	}

	ethClient, err := ethclient.Dial(cfg.Registry.JsonRpc.Url)
	if err != nil {
		return fmt.Errorf("failed to dial api: %v", err)
	}
	registry, err := contract_scanner_registry.NewScannerRegistryCaller(
		common.HexToAddress(cfg.ScannerRegistryContractAddress),
		ethClient,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract caller: %v", err)
	}
	scannerKey, err := security.LoadKeyWithPassphrase(cfg.KeyDirPath, cfg.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to load scanner key: %v", err)
	}
	scannerState, err := registry.GetScannerState(nil, utils.ScannerIDHexToBigInt(scannerKey.Address.Hex()))
	if err != nil && !strings.Contains(err.Error(), "reverted") {
		return fmt.Errorf("failed to check scanner enablement state: %v", err)
	}
	// treat reverts the same as non-registered
	if !scannerState.Registered {
		yellowBold("Scanner not registered - please make sure you register with 'forta register' first.\n")
		return errors.New("cannot run scanner")
	}
	if !scannerState.Enabled {
		yellowBold("Scanner not enabled - please ensure that you have registered with 'forta register' first and staked minimum required amount of FORT.\n")
		return errors.New("cannot run scanner")
	}
	return nil
}
