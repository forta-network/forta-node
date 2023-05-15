package cmd

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/registry"
	mock_registry "github.com/forta-network/forta-core-go/registry/mocks"
	"github.com/forta-network/forta-core-go/security"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	testPoolID = int64(123)
)

func TestPoolAuthorization(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	regClient := mock_registry.NewMockClient(ctrl)

	dir := t.TempDir()
	scannerKeyStore := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := scannerKeyStore.NewAccount("Forta123")
	r.NoError(err)

	scannerKey, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	r.NoError(err)

	regClient.EXPECT().SetRegistryChainID(cfg.Registry.ChainID)
	regClient.EXPECT().GetPoolScanner(scannerKey.Address.Hex()).Return(nil, nil)
	regClient.EXPECT().WillNewScannerShutdownPool(big.NewInt(testPoolID)).Return(false, nil)
	regClient.EXPECT().GenerateScannerRegistrationSignature(gomock.Any()).Return(&registry.ScannerRegistrationInfo{}, nil)

	err = authorizePoolWithRegistry(regClient, scannerKey, testPoolID, false, false, false)
	r.NoError(err)
}

func TestPoolAuthorization_ScannerExists(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	regClient := mock_registry.NewMockClient(ctrl)

	dir := t.TempDir()
	scannerKeyStore := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := scannerKeyStore.NewAccount("Forta123")
	r.NoError(err)

	scannerKey, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	r.NoError(err)

	regClient.EXPECT().SetRegistryChainID(cfg.Registry.ChainID)
	regClient.EXPECT().GetPoolScanner(scannerKey.Address.Hex()).Return(&registry.Scanner{}, nil)

	err = authorizePoolWithRegistry(regClient, scannerKey, testPoolID, false, false, false)
	r.NoError(err)
}

func TestPoolAuthorization_Polygonscan(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	regClient := mock_registry.NewMockClient(ctrl)

	dir := t.TempDir()
	scannerKeyStore := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := scannerKeyStore.NewAccount("Forta123")
	r.NoError(err)

	scannerKey, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	r.NoError(err)

	regClient.EXPECT().SetRegistryChainID(cfg.Registry.ChainID)
	regClient.EXPECT().GetPoolScanner(scannerKey.Address.Hex()).Return(nil, nil)
	regClient.EXPECT().WillNewScannerShutdownPool(big.NewInt(testPoolID)).Return(false, nil)
	regClient.EXPECT().GenerateScannerRegistrationSignature(gomock.Any()).Return(&registry.ScannerRegistrationInfo{}, nil)

	err = authorizePoolWithRegistry(regClient, scannerKey, testPoolID, true, false, false)
	r.NoError(err)
}

func TestPoolAuthorization_GenerateForRegisteredScanner(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	regClient := mock_registry.NewMockClient(ctrl)

	dir := t.TempDir()
	scannerKeyStore := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := scannerKeyStore.NewAccount("Forta123")
	r.NoError(err)

	scannerKey, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	r.NoError(err)

	regClient.EXPECT().SetRegistryChainID(cfg.Registry.ChainID)
	regClient.EXPECT().GetPoolScanner(scannerKey.Address.Hex()).Return(&registry.Scanner{}, nil)
	regClient.EXPECT().WillNewScannerShutdownPool(big.NewInt(testPoolID)).Return(false, nil)
	regClient.EXPECT().GenerateScannerRegistrationSignature(gomock.Any()).Return(&registry.ScannerRegistrationInfo{}, nil)

	err = authorizePoolWithRegistry(regClient, scannerKey, testPoolID, false, true, false)
	r.NoError(err)
}

func TestPoolAuthorization_CleanOutput(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	regClient := mock_registry.NewMockClient(ctrl)

	dir := t.TempDir()
	scannerKeyStore := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := scannerKeyStore.NewAccount("Forta123")
	r.NoError(err)

	scannerKey, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	r.NoError(err)

	regClient.EXPECT().SetRegistryChainID(cfg.Registry.ChainID)
	regClient.EXPECT().GetPoolScanner(scannerKey.Address.Hex()).Return(nil, nil)
	regClient.EXPECT().WillNewScannerShutdownPool(big.NewInt(testPoolID)).Return(false, nil)
	regClient.EXPECT().GenerateScannerRegistrationSignature(gomock.Any()).Return(&registry.ScannerRegistrationInfo{}, nil)

	err = authorizePoolWithRegistry(regClient, scannerKey, testPoolID, false, false, true)
	r.NoError(err)
}
