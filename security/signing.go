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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"

	"forta-network/forta-node/protocol"
)

// LoadKey loads the node private key.
func LoadKey(keysDirPath string) (*keystore.Key, error) {
	f, err := os.OpenFile("/passphrase", os.O_RDONLY, 400)
	if err != nil {
		return nil, err
	}

	pw, err := io.ReadAll(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}
	passphrase := string(pw)

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

// SignProtoMessage marshals a message and signs.
func SignProtoMessage(key *keystore.Key, m protoiface.MessageV1) (*protocol.Signature, error) {
	b, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
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

// NewTransactOpts creates new opts with the private key.
func NewTransactOpts(key *keystore.Key) *bind.TransactOpts {
	return bind.NewKeyedTransactor(key.PrivateKey)
}
