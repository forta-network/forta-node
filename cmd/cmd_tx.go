package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fatih/color"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-core-go/security/eip712"
	"github.com/forta-network/forta-node/store"
	"github.com/spf13/cobra"
)

func handleFortaRegister(cmd *cobra.Command, args []string) error {
	poolIDStr, err := cmd.Flags().GetString("pool-id")
	if err != nil {
		return err
	}
	poolID, err := hexutil.DecodeBig(poolIDStr)
	if err != nil {
		return fmt.Errorf("failed to decode pool ID: %v", err)
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	scannerKey, err := security.LoadKeyWithPassphrase(cfg.KeyDirPath, cfg.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to load scanner key: %v", err)
	}
	scannerPrivateKey := scannerKey.PrivateKey

	registry, err := store.GetRegistryClient(context.Background(), cfg, registry.ClientConfig{
		JsonRpcUrl: cfg.Registry.JsonRpc.Url,
		ENSAddress: cfg.ENSConfig.ContractAddress,
		Name:       "registry-client",
		PrivateKey: scannerPrivateKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create registry client: %v", err)
	}

	encodedPayload, sig, err := registry.GenerateScannerRegistrationSignature(&eip712.ScannerNodeRegistration{
		Scanner:       scannerKey.Address,
		ScannerPoolId: poolID,
		ChainId:       big.NewInt(int64(cfg.ChainID)),
		Metadata:      "",
		Timestamp:     big.NewInt(time.Now().Unix()),
	})
	if err != nil {
		return fmt.Errorf("failed to generate registration signature: %v", err)
	}

	encodedSig, err := security.EncodeEthereumSignature(sig)
	if err != nil {
		return fmt.Errorf("failed to encode registration signature: %v", err)
	}

	whiteBold("Please use the registration signature below on https://app.forta.network as soon as possible!\n\n")

	if verbose {
		color.New(color.FgYellow).Println("encoded:", hex.EncodeToString(encodedPayload))
		color.New(color.FgYellow).Println("signature:", encodedSig)
	} else {
		color.New(color.FgYellow).Println(encodedSig)
	}

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
