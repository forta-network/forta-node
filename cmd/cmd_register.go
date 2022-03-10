package cmd

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/forta-protocol/forta-core-go/contracts/contract_scanner_registry"
	"github.com/forta-protocol/forta-core-go/ens"
	"github.com/forta-protocol/forta-core-go/security"
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
	ownerAddress := common.HexToAddress(ownerAddressStr)

	scannerKey, err := security.LoadKeyWithPassphrase(cfg.KeyDirPath, cfg.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to load scanner key: %v", err)
	}
	scannerAddressStr := scannerKey.Address.Hex()

	if strings.EqualFold(scannerAddressStr, ownerAddressStr) {
		redBold("Scanner and owner cannot be the same identity! Please provide a different wallet address of your own.\n")
	}

	ensStore, err := ens.DialENSStore(cfg.Registry.JsonRpc.Url)
	if err != nil {
		return fmt.Errorf("failed to dial ens store: %v", err)
	}
	scannerRegistryAddress, err := ensStore.Resolve(ens.ScannerRegistryContract)
	if err != nil {
		return fmt.Errorf("failed to resolve scanner registry address: %v", err)
	}
	ethClient, err := ethclient.Dial(cfg.Registry.JsonRpc.Url)
	if err != nil {
		return fmt.Errorf("failed to dial api: %v", err)
	}
	registry, err := contract_scanner_registry.NewScannerRegistryTransactor(scannerRegistryAddress, ethClient)
	if err != nil {
		return fmt.Errorf("failed to create contract transactor: %v", err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(scannerKey.PrivateKey, big.NewInt(137))
	if err != nil {
		return fmt.Errorf("failed to create transaction opts: %v", err)
	}

	color.Yellow(fmt.Sprintf("Sending a transaction to register your scanner to chain %d...\n", cfg.ChainID))

	tx, err := registry.Register(opts, ownerAddress, big.NewInt(int64(cfg.ChainID)), "")
	if err != nil {
		return fmt.Errorf("failed to send the transaction: %v", err)
	}
	txHash := tx.Hash().Hex()

	greenBold("Successfully sent the transaction!\n\n")
	whiteBold("Please ensure that https://polygonscan.com/tx/%s succeeds before you do 'forta run'. This can take a while depending on the network load.\n", txHash)

	return nil
}
