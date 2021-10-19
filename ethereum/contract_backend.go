package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// ContractBackend is the same interface.
type ContractBackend interface {
	bind.ContractBackend
}

// contractBackend is a wrapper of go-ethereum client. This is useful for implementing
// extra features. It's not thread-safe.
type contractBackend struct {
	nonce uint64
	ContractBackend
}

// NewContractBackend creates a new contract backend by wrapping `ethclient.Client`.
func NewContractBackend(client *rpc.Client) bind.ContractBackend {
	return &contractBackend{ContractBackend: ethclient.NewClient(client)}
}

// PendingNonceAt helps us count the nonce more robustly.
func (cb *contractBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	pendingNonce, err := cb.ContractBackend.PendingNonceAt(ctx, account)
	if err != nil {
		return 0, err
	}
	if pendingNonce > cb.nonce {
		return pendingNonce, nil
	}
	return cb.nonce, nil
}

// SendTransaction sends the transaction with the most up-to-date nonce.
func (cb *contractBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if err := cb.ContractBackend.SendTransaction(ctx, tx); err != nil {
		return err
	}
	// count it locally: if sending the tx is successful than that's the previous nonce for sure
	newNonce := tx.Nonce() + 1
	if newNonce > cb.nonce {
		cb.nonce = newNonce
	}
	return nil
}
