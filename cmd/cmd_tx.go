package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/store"
	"github.com/spf13/cobra"
)

func handleFortaRegister(cmd *cobra.Command, args []string) error {
	ownerAddressStr, err := cmd.Flags().GetString("owner-address")
	if err != nil {
		return err
	}
	if !common.IsHexAddress(ownerAddressStr) {
		return errors.New("invalid owner address provided")
	}

	scannerKey, err := security.LoadKeyWithPassphrase(cfg.KeyDirPath, cfg.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to load scanner key: %v", err)
	}
	scannerPrivateKey := scannerKey.PrivateKey
	scannerAddressStr := scannerKey.Address.Hex()

	if strings.EqualFold(scannerAddressStr, ownerAddressStr) {
		redBold("Scanner and owner cannot be the same identity! Please provide a different wallet address of your own.\n")
	}

	registry, err := store.GetRegistryClient(context.Background(), cfg, registry.ClientConfig{
		JsonRpcUrl: cfg.Registry.JsonRpc.Url,
		ENSAddress: cfg.ENSConfig.ContractAddress,
		Name:       "registry-client",
		PrivateKey: scannerPrivateKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create registry client: %v", err)
	}

	color.Yellow(fmt.Sprintf("Sending a transaction to register your scan node to chain %d...\n", cfg.ChainID))

	txHash, err := registry.RegisterScannerOld(ownerAddressStr, int64(cfg.ChainID), "")
	if err != nil && strings.Contains(err.Error(), "insufficient funds") {
		yellowBold("This action requires Polygon (Mainnet) MATIC. Have you funded your address %s yet?\n", scannerAddressStr)
	}
	if err != nil {
		return fmt.Errorf("failed to send the transaction: %v", err)
	}

	greenBold("Successfully sent the transaction!\n\n")
	whiteBold("Please ensure that https://polygonscan.com/tx/%s succeeds before you do 'forta run'. This can take a while depending on the network load.\n", txHash)

	return nil
}

func handleFortaEnable(cmd *cobra.Command, args []string) error {
	scannerKey, err := security.LoadKeyWithPassphrase(cfg.KeyDirPath, cfg.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to load scanner key: %v", err)
	}
	scannerPrivateKey := scannerKey.PrivateKey
	scannerAddressStr := scannerKey.Address.Hex()

	reg, err := store.GetRegistryClient(context.Background(), cfg, registry.ClientConfig{
		JsonRpcUrl: cfg.Registry.JsonRpc.Url,
		ENSAddress: cfg.ENSConfig.ContractAddress,
		Name:       "registry-client",
		PrivateKey: scannerPrivateKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create registry client: %v", err)
	}

	color.Yellow("Sending a transaction to enable your scan node...\n")

	txHash, err := reg.EnableScanner(registry.ScannerPermissionSelf, scannerAddressStr)
	if err != nil && strings.Contains(err.Error(), "insufficient funds") {
		yellowBold("This action requires Polygon (Mainnet) MATIC. Have you funded your address %s yet?\n", scannerAddressStr)
	}
	if err != nil {
		return fmt.Errorf("failed to send the transaction: %v", err)
	}

	greenBold("Successfully sent the transaction!\n\n")
	whiteBold("https://polygonscan.com/tx/%s\n", txHash)

	return nil
}

func handleFortaDisable(cmd *cobra.Command, args []string) error {
	scannerKey, err := security.LoadKeyWithPassphrase(cfg.KeyDirPath, cfg.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to load scanner key: %v", err)
	}
	scannerPrivateKey := scannerKey.PrivateKey
	scannerAddressStr := scannerKey.Address.Hex()

	reg, err := store.GetRegistryClient(context.Background(), cfg, registry.ClientConfig{
		JsonRpcUrl: cfg.Registry.JsonRpc.Url,
		ENSAddress: cfg.ENSConfig.ContractAddress,
		Name:       "registry-client",
		PrivateKey: scannerPrivateKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create registry client: %v", err)
	}

	color.Yellow("Sending a transaction to disable your scan node...\n")

	txHash, err := reg.DisableScanner(registry.ScannerPermissionSelf, scannerAddressStr)
	if err != nil && strings.Contains(err.Error(), "insufficient funds") {
		yellowBold("This action requires Polygon (Mainnet) MATIC. Have you funded your address %s yet?\n", scannerAddressStr)
	}
	if err != nil {
		return fmt.Errorf("failed to send the transaction: %v", err)
	}

	greenBold("Successfully sent the transaction!\n\n")
	whiteBold("https://polygonscan.com/tx/%s\n", txHash)

	return nil
}
