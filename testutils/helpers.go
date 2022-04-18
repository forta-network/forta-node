package testutils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/forta-network/forta-core-go/domain"
)

// mocking code borrowed from go-ethereum/core/types tests

// TestTxs get a list of mock transactions with the following nonces
func TestTxs(nonces ...int) []domain.Transaction {
	var result []domain.Transaction
	for _, nonce := range nonces {
		result = append(result, domain.Transaction{
			Hash:  fmt.Sprintf("%x", nonce),
			Nonce: fmt.Sprintf("%x", nonce),
		})
	}
	return result
}

// TestBlock gets a block with a list of transactions with the following nonces
func TestBlock(nonces ...int) *domain.Block {
	return &domain.Block{
		Number:       "0x0",
		Hash:         "0x1",
		Transactions: TestTxs(nonces...),
	}
}

func TestLogs(indexes ...int) []types.Log {
	var result []types.Log
	for _, index := range indexes {
		result = append(result, types.Log{
			TxHash:  common.HexToHash(fmt.Sprintf("%x", index)),
			TxIndex: uint(index),
		})
	}
	return result
}
