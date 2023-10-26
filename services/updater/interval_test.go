package updater

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalculateReleaseDelay(t *testing.T) {
	r := require.New(t)

	scannerAddr := "0x3f88c2b3e267e6b8e9dE017cdB47a59aC9Ecb284"
	delay := CalculateReleaseDelay(scannerAddr, time.Hour*2)
	expected := time.Duration(3661000000000) // 1h1m1s
	r.Equal(expected, delay)

	scannerAddr = "0x593Ad20C41660a73C1D708FcF4a2B3653063a183"
	delay = CalculateReleaseDelay(scannerAddr, time.Hour*2)
	expected = time.Duration(5666000000000) // 1h34m26s
	r.Equal(expected, delay)

	scannerAddr = "0x0Cfa28a1293E8bb83bC3a9cb08Cc4E68b21B8cb9"
	delay = CalculateReleaseDelay(scannerAddr, time.Hour*2)
	expected = time.Duration(1739000000000) // 28m59s
	r.Equal(expected, delay)
}
