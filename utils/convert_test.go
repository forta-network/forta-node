package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexToBigInt(t *testing.T) {
	hex := "0xabc"
	res, err := HexToBigInt(hex)
	assert.NoError(t, err)
	assert.Equal(t, int64(2748), res.Int64())
}

func TestBytesToHex(t *testing.T) {
	value := "ABCDEFG12382837581235uASDFASDFASDFzxcvzxcvzxcv"
	assert.Equal(t, "0x37353831323335754153444641534446415344467a7863767a7863767a786376", BytesToHex([]byte(value)))
}
