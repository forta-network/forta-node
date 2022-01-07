package security

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"

	"github.com/forta-protocol/forta-node/protocol"
)

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

// SignAlert signs the alert.
func SignAlert(key *keystore.Key, alert *protocol.Alert) (*protocol.SignedAlert, error) {
	signature, err := SignProtoMessage(key, alert)
	if err != nil {
		return nil, err
	}
	return &protocol.SignedAlert{
		Alert:     alert,
		Signature: signature,
	}, nil
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
		return fmt.Errorf("invalid signature expected=%s, got=%s", signerAddress, addr.Hex())
	}

	return nil
}

func VerifyProtoSignature(m protoiface.MessageV1, signerAddress, sigHex string) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return VerifySignature(b, signerAddress, sigHex)
}

// SignProtoMessage marshals a message and signs.
func SignProtoMessage(key *keystore.Key, m protoiface.MessageV1) (*protocol.Signature, error) {
	b, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	return SignBytes(key, b)
}

// NewTransactOpts creates new opts with the private key.
func NewTransactOpts(key *keystore.Key) *bind.TransactOpts {
	return bind.NewKeyedTransactor(key.PrivateKey)
}
