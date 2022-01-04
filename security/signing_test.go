package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	signature = "0x815136705413e8608fb33c7eab05057d1c697db2b8f8fc22e4e29c0d980002626a292cf12a863192f162c19576288c937e162301bc79dcbd006b1e76aea264b101"
	signer    = "0xeE0D82ac806efe2b9a0003a27a785458bC67bbf0"
	ref       = "Qmc2Dmb3wAycyeg3E7Nf6AANqDeBhiX4rdSy3ZJqg2PpMP"
)

func TestSignString(t *testing.T) {
	MockPassphrase = "Forta123"
	key, err := LoadKey("testkey")
	assert.NoError(t, err)

	res, err := SignString(key, ref)
	assert.NoError(t, err)
	assert.Equal(t, signature, res.Signature)
	assert.Equal(t, signer, res.Signer)
}

func TestVerifySignature(t *testing.T) {
	err := VerifySignature([]byte(ref), signer, signature)
	assert.NoError(t, err)
}
