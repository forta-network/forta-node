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
