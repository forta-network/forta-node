package cmd

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/require"
)

func TestShowAccounts(t *testing.T) {
	r := require.New(t)

	dir := t.TempDir()
	cfg.KeyDirPath = dir
	scannerKeyStore := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := scannerKeyStore.NewAccount("Forta123")
	r.NoError(err)

	r.NoError(handleFortaAccountAddress(nil, nil))
}

func TestShowAccounts_NoAccounts(t *testing.T) {
	r := require.New(t)

	dir := t.TempDir()
	cfg.KeyDirPath = dir
	keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	r.Error(handleFortaAccountAddress(nil, nil))
}

func TestShowAccounts_MultipleAccounts(t *testing.T) {
	r := require.New(t)

	dir := t.TempDir()
	cfg.KeyDirPath = dir
	scannerKeyStore := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := scannerKeyStore.NewAccount("Forta123")
	r.NoError(err)

	_, err = scannerKeyStore.NewAccount("Forta456")
	r.NoError(err)

	r.Error(handleFortaAccountAddress(nil, nil))
}
