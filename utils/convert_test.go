package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexToBigInt(t *testing.T) {
	hex := "0xabc"
	res := HexToBigInt(hex)
	assert.Equal(t, int64(2748), res.Int64())
}
