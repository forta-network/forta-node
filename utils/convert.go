package utils

import (
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
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

func BytesToHex(b []byte) string {
	return common.BytesToHash(b).Hex()
}

func HexToInt64(hex string) int64 {
	bigInt, err := HexToBigInt(hex)
	if err != nil {
		log.WithField("hex", hex).Error("could not convert hex to number")
		return 0
	}
	return bigInt.Int64()
}
