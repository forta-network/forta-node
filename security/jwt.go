package security

import (
	"crypto/ecdsa"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"
)

const (
	alg = "ETH"
)

func init() {
	jwt.RegisterSigningMethod(alg, func() jwt.SigningMethod {
		return ethSigningMethod{}
	})
}

type ethSigningMethod struct{}

func (e ethSigningMethod) Verify(signingString, signature string, key interface{}) error {
	address := key.(string)
	hash := crypto.Keccak256([]byte(signingString))
	sig, err := base64.RawURLEncoding.DecodeString(signature)

	if err != nil {
		return err
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return err
	}

	if sigPublicKeyECDSA == nil {
		return errors.New("could not recover address (pub is nil)")
	}
	addr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	if addr.Hex() != address {
		return fmt.Errorf("signature invalid: expected=%s, got=%s", address, addr.Hex())
	}

	return nil
}

func (e ethSigningMethod) Sign(signingString string, key interface{}) (string, error) {
	pk := key.(*ecdsa.PrivateKey)
	hash := crypto.Keccak256([]byte(signingString))
	sig, err := crypto.Sign(hash, pk)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sig), nil
}

func (e ethSigningMethod) Alg() string {
	return alg
}

func VerifyJWT(tokenString string, address string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(ethSigningMethod); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return address, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func CreateJWT(key *keystore.Key, claims map[string]interface{}) (string, error) {
	mapClaims := map[string]interface{}{
		"sub": key.Address.Hex(),
		"nbf": time.Now().UTC(),
	}
	for k, v := range claims {
		mapClaims[k] = v
	}
	token := jwt.NewWithClaims(&ethSigningMethod{}, jwt.MapClaims(mapClaims))
	str, err := token.SignedString(key.PrivateKey)
	if err != nil {
		return "", err
	}
	return str, nil
}
