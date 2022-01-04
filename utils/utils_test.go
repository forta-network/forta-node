package utils

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestAddPercentage(t *testing.T) {
	input := big.NewInt(100)
	AddPercentage(input, 10)
	assert.Equal(t, int64(110), input.Int64())
}
