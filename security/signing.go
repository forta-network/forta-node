package security

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-protocol/forta-node/encoding"
	"github.com/forta-protocol/forta-node/utils"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/forta-protocol/forta-node/protocol"
)

var ErrMissingSignature = errors.New("missing signature")
var ErrInvalidSignature = errors.New("invalid signature")

func ReadPassphrase() (string, error) {
	f, err := os.OpenFile("/passphrase", os.O_RDONLY, 400)
	if err != nil {
		return "", err
	}
	pw, err := io.ReadAll(bufio.NewReader(f))
	if err != nil {
		return "", err
	}
	return string(pw), nil
}

// LoadKey loads the passphrase and the node private key.
func LoadKey(keysDirPath string) (*keystore.Key, error) {
	passphrase, err := ReadPassphrase()
	if err != nil {
		return nil, err
	}
	return LoadKeyWithPassphrase(keysDirPath, passphrase)
}

// LoadKeyWithPassphrase decrypts and loads the node private key using provided passphrase.
func LoadKeyWithPassphrase(keysDirPath, passphrase string) (*keystore.Key, error) {
	files, err := ioutil.ReadDir(keysDirPath)
	if err != nil {
		return nil, err
	}

	if len(files) != 1 {
		return nil, errors.New("there must be only one key in key directory")
	}

	keyBytes, err := ioutil.ReadFile(path.Join(keysDirPath, files[0].Name()))
	if err != nil {
		return nil, err
	}

	return keystore.DecryptKey(keyBytes, passphrase)
}

func alertHash(alert *protocol.Alert) common.Hash {
	metadata := utils.MapToList(alert.Metadata)
	alertStr := fmt.Sprintf("%s%s%s", alert.Id, strings.Join(metadata, ""), alert.Timestamp)
	return crypto.Keccak256Hash([]byte(alertStr))
}

// SignAlert signs the alert using the alertID and deterministicly formatted Metadata
func SignAlert(key *keystore.Key, alert *protocol.Alert) (*protocol.SignedAlert, error) {
	hash := alertHash(alert)
	signature, err := SignBytes(key, hash.Bytes())
	if err != nil {
		return nil, err
	}
	return &protocol.SignedAlert{
		Alert:     alert,
		Signature: signature,
	}, nil
}

// VerifyAlertSignature returns an error if the signature for the signed alert is invalid
func VerifyAlertSignature(sa *protocol.SignedAlert) error {
	if sa.Signature == nil {
		return ErrMissingSignature
	}
	hash := alertHash(sa.Alert)
	return VerifySignature(hash.Bytes(), sa.Signature.Signer, sa.Signature.Signature)
}

func SignBytes(key *keystore.Key, b []byte) (*protocol.Signature, error) {
	hash := crypto.Keccak256(b)
	sig, err := crypto.Sign(hash, key.PrivateKey)

	if err != nil {
		return nil, err
	}

	return &protocol.Signature{
		Signature: fmt.Sprintf("0x%s", hex.EncodeToString(sig)),
		Algorithm: "ECDSA",
		Signer:    key.Address.Hex(),
	}, nil
}

func SignString(key *keystore.Key, input string) (*protocol.Signature, error) {
	return SignBytes(key, []byte(input))
}

func VerifySignature(message []byte, signerAddress string, sigHex string) error {
	hash := crypto.Keccak256Hash(message)
	sigHex = strings.ReplaceAll(sigHex, "0x", "")
	signature, err := hex.DecodeString(sigHex)

	if err != nil {
		return err
	}

	pubKey, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return err
	}

	if pubKey == nil {
		return errors.New("could not recover address (pub is nil)")
	}

	addr := crypto.PubkeyToAddress(*pubKey)

	if addr.Hex() != signerAddress {
		return ErrInvalidSignature
	}

	return nil
}

// SignBatch will sign an alert batch and return a SignedAlertBatch
func SignBatch(key *keystore.Key, batch *protocol.AlertBatch) (*protocol.SignedAlertBatch, error) {
	encoded, err := encoding.EncodeBatch(batch)
	if err != nil {
		return nil, err
	}
	signature, err := SignString(key, encoded)
	if err != nil {
		return nil, err
	}

	return &protocol.SignedAlertBatch{
		//TODO: remove Data in subsequent deploy
		Data:      batch,
		Encoded:   encoded,
		Signature: signature,
	}, nil
}

// VerifyBatchSignature will return an error if the signature fails to validate
func VerifyBatchSignature(signedBatch *protocol.SignedAlertBatch) error {
	if signedBatch.Signature == nil {
		return ErrMissingSignature
	}
	return VerifySignature([]byte(signedBatch.Encoded), signedBatch.Signature.Signer, signedBatch.Signature.Signature)
}

// NewTransactOpts creates new opts with the private key.
func NewTransactOpts(key *keystore.Key) *bind.TransactOpts {
	return bind.NewKeyedTransactor(key.PrivateKey)
}
