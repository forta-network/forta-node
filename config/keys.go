package config

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/forta-network/forta-core-go/security"
)

// LoadKeyInContainer loads the key in the service container depending on the config.
func LoadKeyInContainer(cfg Config) (*keystore.Key, error) {
	if len(cfg.LocalModeConfig.PrivateKeyHex) > 0 {
		privKey, err := crypto.HexToECDSA(cfg.LocalModeConfig.PrivateKeyHex)
		if err != nil {
			return nil, fmt.Errorf("invalid local mode key: %v", err)
		}
		publicAddr := crypto.PubkeyToAddress(privKey.PublicKey)
		return &keystore.Key{
			PrivateKey: privKey,
			Address:    publicAddr,
		}, nil
	}
	return security.LoadKey(DefaultContainerKeyDirPath)
}
