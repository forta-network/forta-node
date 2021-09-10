package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func BigIntToHex(i *big.Int) string {
	return hexutil.EncodeBig(i)
}

func HexToBigInt(hex string) (*big.Int, error) {
	return hexutil.DecodeBig(hex)
}

func Bytes32ToHex(b [32]byte) string {
	return common.BytesToHash(b[:]).Hex()
}
