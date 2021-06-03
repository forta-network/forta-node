package testutils

import (
	"fmt"

	"github.com/OpenZeppelin/fortify-node/domain"
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
		Hash:         "0x1",
		Transactions: TestTxs(nonces...),
	}
}
