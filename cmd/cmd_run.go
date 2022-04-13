package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/forta-protocol/forta-core-go/registry"
	"github.com/forta-protocol/forta-core-go/security"
	"github.com/forta-protocol/forta-node/cmd/runner"
	"github.com/forta-protocol/forta-node/store"
	"github.com/spf13/cobra"
)

// errors
var (
	ErrCannotRunScanner = errors.New("cannot run scanner")
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

	scannerKey, err := security.LoadKeyWithPassphrase(cfg.KeyDirPath, cfg.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to load scanner key: %v", err)
	}
	scannerAddressStr := scannerKey.Address.Hex()

	registry, err := store.GetRegistryClient(context.Background(), cfg, registry.ClientConfig{
		JsonRpcUrl: cfg.Registry.JsonRpc.Url,
		ENSAddress: cfg.ENSConfig.ContractAddress,
		Name:       "registry-client",
	})
	if err != nil {
		return fmt.Errorf("failed to create registry client: %v", err)
	}
	scanner, err := registry.GetScanner(scannerAddressStr)
	if err != nil {
		return fmt.Errorf("failed to check scanner state: %v", err)
	}

	// treat reverts the same as non-registered
	if scanner == nil {
		yellowBold("Scanner not registered - please make sure you register with 'forta register' first.\n")
		toStderr("You can disable this behaviour with --no-check flag.\n")
		return ErrCannotRunScanner
	}
	if !scanner.Enabled {
		yellowBold("Scanner not enabled - please ensure that you have registered with 'forta register' first and staked minimum required amount of FORT.\n")
		toStderr("You can disable this behaviour with --no-check flag.\n")
		return ErrCannotRunScanner
	}
	return nil
}
