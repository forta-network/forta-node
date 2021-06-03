package security

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"

	"fortify-node/protocol"
)

func SignAlert(key *keystore.Key, alert *protocol.Alert) (*protocol.SignedAlert, error) {
	b, err := proto.Marshal(alert)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256(b)
	sig, err := crypto.Sign(hash, key.PrivateKey)
	if err != nil {
		return nil, err
	}
	signature := fmt.Sprintf("0x%s", hex.EncodeToString(sig))
	return &protocol.SignedAlert{
		Alert: alert,
		Signature: &protocol.Signature{
			Signature: signature,
			Algorithm: "ECDSA",
		},
	}, nil
}
