package security

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateJWT(t *testing.T) {
	MockPassphrase = "Forta123"
	key, err := LoadKey("testkey")
	assert.NoError(t, err)

	address := "0xeE0D82ac806efe2b9a0003a27a785458bC67bbf0"
	claims := map[string]interface{}{
		"claim1": "claim",
		"nbf":    12345,
	}

	token, err := CreateJWT(key, claims)
	assert.NoError(t, err)
	t.Log(token)

	validToken, err := VerifyJWT(token, address)
	assert.NoError(t, err)

	c, ok := validToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, address, c["sub"])
}
