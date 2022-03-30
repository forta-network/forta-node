package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/forta-protocol/forta-core-go/registry"
	"github.com/forta-protocol/forta-core-go/security"
	"github.com/forta-protocol/forta-node/store"
	"github.com/spf13/cobra"
)

const registryPermissionSelf uint8 = 1

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
	})
	if err != nil {
		return fmt.Errorf("failed to create registry client: %v", err)
	}

	color.Yellow(fmt.Sprintf("Sending a transaction to register your scanner to chain %d...\n", cfg.ChainID))

	txHash, err := registry.RegisterScanner(scannerPrivateKey, ownerAddressStr, int64(cfg.ChainID), "")
	if err != nil {
		return fmt.Errorf("failed to send the transaction: %v", err)
	}

	greenBold("Successfully sent the transaction!\n\n")
	whiteBold("Please ensure that https://polygonscan.com/tx/%s succeeds before you do 'forta run'. This can take a while depending on the network load.\n", txHash)

	return nil
}
