package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func BigIntToHex(i *big.Int) string {
	return hexutil.EncodeBig(i)
}

func HexToBigInt(hex string) (*big.Int, error) {
	return hexutil.DecodeBig(hex)
}
