package testutils

import (
	"hash"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"
)

// mocking code borrowed from go-ethereum/core/types tests

// TestTxs get a list of mock transactions with the following nonces
func TestTxs(nonces ...int) []*types.Transaction {
	testAddr := common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	var result []*types.Transaction
	for _, nonce := range nonces {
		emptyEip2718Tx := types.NewTx(&types.AccessListTx{
			ChainID:  big.NewInt(1),
			Nonce:    uint64(nonce),
			To:       &testAddr,
			Value:    big.NewInt(10),
			Gas:      25000,
			GasPrice: big.NewInt(1),
			Data:     common.FromHex("5544"),
		})

		tx, _ := emptyEip2718Tx.WithSignature(
			types.NewEIP2930Signer(big.NewInt(1)),
			common.Hex2Bytes("c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b266032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d3752101"),
		)

		result = append(result, tx)
	}
	return result
}

// TestBlock gets a block with a list of transactions with the following nonces
func TestBlock(nonces ...int) *types.Block {
	return types.NewBlock(&types.Header{
		Difficulty: math.BigPow(11, 11),
		Number:     math.BigPow(2, 9),
		GasLimit:   12345678,
		GasUsed:    1476322,
		Time:       9876543,
		Extra:      []byte("coolest block on chain"),
	}, TestTxs(nonces...), nil, nil, newHasher())
}

type testHasher struct {
	hasher hash.Hash
}

func newHasher() *testHasher {
	return &testHasher{hasher: sha3.NewLegacyKeccak256()}
}

func (h *testHasher) Reset() {
	h.hasher.Reset()
}

func (h *testHasher) Update(key, val []byte) {
	h.hasher.Write(key)
	h.hasher.Write(val)
}

func (h *testHasher) Hash() common.Hash {
	return common.BytesToHash(h.hasher.Sum(nil))
}
