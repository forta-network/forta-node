package ethaccounts

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	DefaultPassword = "0"

	GethNodeKey, _  = crypto.HexToECDSA("0da4d32840b0ef3e30c82d9d4772c8e1bbcd1ac6417b46d958fb5c7db0be99c1")
	GethNodeAddress = common.HexToAddress("0x1111e291778AE830cfE4e34185e4e560E94047c7")

	ScannerKey, _  = crypto.HexToECDSA("412ae1bd0021a1489a8824e11edc4a017b4b0e12c39be936089b350ea55e997d")
	ScannerAddress = common.HexToAddress("0x222244861C15A8F2A05fbD15E747Ea8F20c2C0c9")

	DeployerKey, _  = crypto.HexToECDSA("02b432bb5b53daf8edf652c028b2eef9a383d688806b6a0bc2b253b3392195b2")
	DeployerAddress = common.HexToAddress("0x3333C25Cb71F00F113425c60E0CbF551c00cEf49")

	ProxyAdminKey, _  = crypto.HexToECDSA("065d74b69b496014c8d23a7eb60edf1da2928f17a868962d1c3f70b6983bde2d")
	ProxyAdminAddress = common.HexToAddress("0x44443b6c4899e3c11Ff666fD98B1cc9bF283174F")

	AccessAdminKey, _  = crypto.HexToECDSA("25582762d32c064e5c7d86c62e3cebe612a6f7bfac45b6403d395d99f90037ed")
	AccessAdminAddress = common.HexToAddress("0x55557b2a04394aBf4bb216f85628686E496C5aaF")

	ExploiterKey, _  = crypto.HexToECDSA("9983c18517758908acd2fa32909dd1490949eee0fa63501b0adeb02802481773")
	ExploiterAddress = common.HexToAddress("0x66664f69BCFE12A7bA2857fAbC42a02729e5c160")

	MiscKey, _  = crypto.HexToECDSA("4c2c30fe62230a6e4550b7f56e4e877c9fcb5aa9f468f3c0e2f94b24785c019a")
	MiscAddress = common.HexToAddress("0x1337B4cBAe461949A00854EECd27Bc331CcaD2f1")

	ForwarderKey, _  = crypto.HexToECDSA("ec80d29573324a3ba1b4e9e9f8376282816bfb9876100af955bc098bf6c986c6")
	ForwarderAddress = common.HexToAddress("0x2337608875c0D3eFDf5232aFf3343a43C73a900F")
)
