package security

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
)

func TestCreateBotJWT(t *testing.T) {
	testErr := errors.New("test")
	testCases := []struct {
		name        string
		claims      map[string]interface{}
		jwtFunc     func(key *keystore.Key, c map[string]interface{}) (string, error)
		expectedJWT string
		expectedErr error
	}{
		{
			name: "valid claims",
			claims: map[string]interface{}{
				"test": "value",
			},
			jwtFunc: func(key *keystore.Key, c map[string]interface{}) (string, error) {
				assert.Equal(t, "value", c["test"])
				assert.Equal(t, "botID", c["bot-id"])
				return "jwt", nil
			},
			expectedJWT: "jwt",
			expectedErr: nil,
		},
		{
			name:   "default claims",
			claims: nil,
			jwtFunc: func(key *keystore.Key, c map[string]interface{}) (string, error) {
				assert.Equal(t, "botID", c["bot-id"])
				return "jwt", nil
			},
			expectedJWT: "jwt",
			expectedErr: nil,
		},
		{
			name:   "jwt creation error",
			claims: nil,
			jwtFunc: func(key *keystore.Key, c map[string]interface{}) (string, error) {
				assert.Equal(t, "botID", c["bot-id"])
				return "", testErr
			},
			expectedJWT: "",
			expectedErr: testErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jwt, err := CreateBotJWT(nil, "botID", tc.claims, tc.jwtFunc)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedJWT, jwt)
		})
	}
}
