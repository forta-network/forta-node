package security

import (
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestCreateJWT(t *testing.T) {
	key, err := LoadKeyWithPassphrase("testkey", "Forta123")
	assert.NoError(t, err)

	address := "0xeE0D82ac806efe2b9a0003a27a785458bC67bbf0"

	token, err := CreateScannerJWT(key, map[string]interface{}{
		"batch": "QmNvoaBmvjaVukfSyZtnHYYzN3iBHV4V3WyKHNwTnoubNf",
	})
	assert.NoError(t, err)
	t.Log(token)

	validToken, err := VerifyScannerJWT(token)
	assert.NoError(t, err)

	claims := validToken.Token.Claims.(jwt.MapClaims)
	for k, v := range claims {
		t.Logf("%s = %v", k, v)
	}

	assert.Equal(t, address, validToken.Scanner)
}
