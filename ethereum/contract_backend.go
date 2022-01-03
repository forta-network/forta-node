package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

const (
	maxNonceDrift = 50
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
	logger := log.WithField("address", account.Hex())
	serverNonce, err := cb.ContractBackend.PendingNonceAt(ctx, account)
	if err != nil {
		logger.WithError(err).Error("failed to get pending nonce from server")
		return 0, err
	}
	logger = logger.WithFields(log.Fields{
		"serverNonce": serverNonce,
		"localNonce":  cb.nonce,
	})
	switch {
	case cb.nonce > serverNonce && cb.nonce-serverNonce >= maxNonceDrift:
		logger.Warn("resetted local nonce")
		cb.nonce = serverNonce
		return serverNonce, nil

	case serverNonce > cb.nonce:
		logger.Info("using server nonce")
		return serverNonce, nil

	default:
		logger.Info("using local nonce")
		return cb.nonce, nil
	}
}

// SendTransaction sends the transaction with the most up-to-date nonce.
func (cb *contractBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	logger := getTxLogger(tx)
	logger.Info("sending")
	if err := cb.ContractBackend.SendTransaction(ctx, tx); err != nil {
		logger.WithError(err).Error("failed to send")
		return err
	}
	logger.Info("sent")
	// count it locally: if sending the tx is successful than that's the previous nonce for sure
	newNonce := tx.Nonce() + 1
	if newNonce > cb.nonce {
		cb.nonce = newNonce
	}
	return nil
}

func getTxLogger(tx *types.Transaction) *log.Entry {
	return log.WithFields(log.Fields{
		"to":       tx.To().Hex(),
		"nonce":    tx.Nonce(),
		"gasLimit": tx.Gas(),
		"gasPrice": tx.GasPrice().Uint64(),
		"hash":     tx.Hash().Hex(),
	})
}
