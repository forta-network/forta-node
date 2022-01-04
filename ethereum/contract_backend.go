package ethereum

import (
	"context"
	"github.com/forta-protocol/forta-node/utils"
	"math/big"
	"time"

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
	nonce           uint64
	gasPrice        *big.Int
	gasPriceUpdated time.Time
	maxPrice        *big.Int
	ContractBackend
}

// NewContractBackend creates a new contract backend by wrapping `ethclient.Client`.
func NewContractBackend(client *rpc.Client, maxPrice *big.Int) bind.ContractBackend {
	return &contractBackend{
		ContractBackend: ethclient.NewClient(client),
		maxPrice:        maxPrice,
	}
}

// SuggestGasPrice retrieves the currently suggested gas price and adds 10%
func (cb *contractBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if cb.gasPrice != nil && time.Since(cb.gasPriceUpdated) < 1*time.Minute {
		return cb.gasPrice, nil
	}
	gp, err := cb.ContractBackend.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	utils.AddPercentage(gp, 10)
	if cb.maxPrice != nil {
		if gp.Cmp(cb.maxPrice) == 1 {
			log.WithFields(log.Fields{
				"suggested": gp.Int64(),
				"maximum":   cb.maxPrice.Int64(),
			}).Warn("returning maximum price")
			return cb.maxPrice, nil
		}
	}
	//TODO: drop to debug
	log.WithFields(log.Fields{
		"gasPrice": gp.Int64(),
	}).Info("returning gas price")
	cb.gasPriceUpdated = time.Now()
	cb.gasPrice = gp
	return gp, nil
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

func incrementableError(err error) bool {
	return err.Error() == "replacement transaction underpriced"
}

// SendTransaction sends the transaction with the most up-to-date nonce.
func (cb *contractBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	logger := getTxLogger(tx)
	logger.Info("sending")
	if err := cb.ContractBackend.SendTransaction(ctx, tx); err != nil {
		logger.WithError(err).Error("failed to send")
		if incrementableError(err) {
			cb.incrementNonce(tx)
		}
		return err
	}
	logger.Info("sent")
	// count it locally: if sending the tx is successful than that's the previous nonce for sure
	cb.incrementNonce(tx)
	return nil
}

func (cb *contractBackend) incrementNonce(tx *types.Transaction) {
	newNonce := tx.Nonce() + 1
	if newNonce > cb.nonce {
		cb.nonce = newNonce
	}
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
