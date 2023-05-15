package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

func handleFortaAccountAddress(cmd *cobra.Command, args []string) error {
	ks := keystore.NewKeyStore(cfg.KeyDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	accounts := ks.Accounts()
	if len(accounts) > 1 {
		redBold("You have multiple accounts. Please import your scanner account again with 'forta account import'.")
		fmt.Println("Your current account addresses:")
		for _, account := range accounts {
			fmt.Println(account.Address.Hex())
		}
		return errors.New("multiple accounts")
	}

	if len(accounts) == 0 {
		redBold("You have no accounts. Please import your scanner account with 'forta account import'.")
		return errors.New("no accounts")
	}

	fmt.Println(accounts[0].Address.Hex())
	return nil
}

func handleFortaAccountImport(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read the private key: %v", err)
	}
	hexKey := strings.TrimSpace(string(b))

	if len(cfg.Passphrase) == 0 {
		redBold("Your passphrase is not set. Please set it with FORTA_PASSPHRASE environment variable or provide it with the --passphrase flag.\n")
		return errors.New("empty passhphrase")
	}

	os.RemoveAll(cfg.KeyDirPath)
	ks := keystore.NewKeyStore(cfg.KeyDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return fmt.Errorf("could not parse the private key hex: %v", err)
	}
	account, err := ks.ImportECDSA(privateKey, cfg.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to import: %v", err)
	}
	fmt.Println(account.Address.Hex())
	return nil
}
