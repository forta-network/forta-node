package security

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func CreateBotJWT(key *keystore.Key, agentID string, claims map[string]interface{}, creator func(key *keystore.Key, claims map[string]interface{}) (string, error)) (string, error) {
	if claims == nil {
		claims = make(map[string]interface{})
	}

	claims["bot-id"] = agentID

	return creator(key, claims)
}
