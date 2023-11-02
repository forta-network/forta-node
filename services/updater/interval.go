package updater

import (
	"math/big"
	"time"

	"github.com/forta-network/forta-core-go/utils"
)

const minUpdateDelay = 1 * time.Minute

// CalculateReleaseDelay calculates a release delay for a given scanner address
// that is within the max update interval.
func CalculateReleaseDelay(scannerAddr string, maxUpdateDelay time.Duration) time.Duration {
	interval := big.NewInt(0)
	interval.Mod(
		utils.ScannerIDHexToBigInt(scannerAddr), // scanner address is naturally random number
		// taking module using max gives a remainder that is smaller than max
		// if we use 24h, the result is smaller than that
		big.NewInt((maxUpdateDelay).Milliseconds()),
	)
	intervalMs := interval.Int64() + minUpdateDelay.Milliseconds()
	return time.Duration(intervalMs/int64(1000)) * time.Second
}
