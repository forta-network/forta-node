package security

import (
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	signature = "0x815136705413e8608fb33c7eab05057d1c697db2b8f8fc22e4e29c0d980002626a292cf12a863192f162c19576288c937e162301bc79dcbd006b1e76aea264b101"
	signer    = "0xeE0D82ac806efe2b9a0003a27a785458bC67bbf0"
	ref       = "Qmc2Dmb3wAycyeg3E7Nf6AANqDeBhiX4rdSy3ZJqg2PpMP"
)

func getTestAlert() *protocol.Alert {
	entropicMap := make(map[string]string)
	for i := 0; i < 10; i++ {
		entropicMap[uuid.Must(uuid.NewUUID()).String()] = "true"
	}
	return &protocol.Alert{
		Id:   "0xabcdefg",
		Type: 2,
		Finding: &protocol.Finding{
			Protocol:    "metadata",
			Severity:    5,
			Metadata:    entropicMap,
			Type:        2,
			AlertId:     "finding",
			Name:        "name",
			Description: "description",
		},
		Timestamp: time.Now().String(),
		Metadata:  entropicMap,
		Agent: &protocol.AgentInfo{
			Image:     "image",
			ImageHash: "hash",
			Id:        "id",
			Manifest:  "manifest",
		},
		Tags:    entropicMap,
		Scanner: &protocol.ScannerInfo{Address: "0xaddress"},
	}
}

func TestSignString(t *testing.T) {
	key, err := LoadKeyWithPassphrase("testkey", "Forta123")
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

func TestVerifyAlertSignature(t *testing.T) {
	key, err := LoadKeyWithPassphrase("testkey", "Forta123")
	assert.NoError(t, err)

	alert := getTestAlert()

	signed, err := SignAlert(key, alert)
	assert.NoError(t, err)
	assert.NotNil(t, signed.Signature)

	t.Log(signed.Signature.Signature)

	assert.NoError(t, VerifyAlertSignature(signed))
}

func TestVerifyAlertSignature_Bad(t *testing.T) {
	key, err := LoadKeyWithPassphrase("testkey", "Forta123")
	assert.NoError(t, err)

	alert1 := getTestAlert()
	alert2 := getTestAlert()

	signed1, err := SignAlert(key, alert1)
	assert.NoError(t, err)

	signed2, err := SignAlert(key, alert2)
	assert.NoError(t, err)

	//copy over the signature from #2
	signed1.Signature = signed2.Signature

	assert.ErrorIs(t, VerifyAlertSignature(signed1), ErrInvalidSignature)
}

func TestVerifyAlertSignature_Missing(t *testing.T) {
	key, err := LoadKeyWithPassphrase("testkey", "Forta123")
	assert.NoError(t, err)

	alert := getTestAlert()

	signed, err := SignAlert(key, alert)
	assert.NoError(t, err)
	assert.NotNil(t, signed.Signature)

	signed.Signature = nil

	assert.ErrorIs(t, VerifyAlertSignature(signed), ErrMissingSignature)
}
