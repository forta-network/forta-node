package scanner

import (
	"context"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-core-go/feeds"
	"github.com/forta-network/forta-node/config"

	log "github.com/sirupsen/logrus"
)

// TxStreamService pulls TX info from providers and emits to channel
type TxStreamService struct {
	cfg         TxStreamServiceConfig
	ctx         context.Context
	blockOutput chan *domain.BlockEvent
	txOutput    chan *domain.TransactionEvent
	txFeed      feeds.TransactionFeed

	lastBlockActivity health.TimeTracker
	lastTxActivity    health.TimeTracker
}

type TxStreamServiceConfig struct {
	JsonRpcConfig       config.JsonRpcConfig
	TraceJsonRpcConfig  config.JsonRpcConfig
	SkipBlocksOlderThan *time.Duration
}

func (t *TxStreamService) ReadOnlyBlockStream() <-chan *domain.BlockEvent {
	return t.blockOutput
}

func (t *TxStreamService) ReadOnlyTxStream() <-chan *domain.TransactionEvent {
	return t.txOutput
}

func (t *TxStreamService) handleBlock(evt *domain.BlockEvent) error {
	t.blockOutput <- evt
	t.lastBlockActivity.Set()
	return nil
}

func (t *TxStreamService) handleTx(evt *domain.TransactionEvent) error {
	t.txOutput <- evt
	t.lastTxActivity.Set()
	return nil
}

func (t *TxStreamService) Start() error {
	go func() {
		if err := t.txFeed.ForEachTransaction(t.handleBlock, t.handleTx); err != nil {
			logger := log.WithError(err)
			if err != context.Canceled {
				logger.Panic("tx feed error")
			}
			logger.Info("tx feed stopped")
		}
	}()
	return nil
}

func (t *TxStreamService) Stop() error {
	if t.txOutput != nil {
		// drain and close tx channel
		func(c chan *domain.TransactionEvent) {
			for {
				select {
				case tx := <-c:
					log.WithFields(log.Fields{"tx": tx.Transaction.Hash}).Info("gracefully draining transaction")
				default:
					close(c)
					return
				}
			}
		}(t.txOutput)
	}
	if t.blockOutput != nil {
		// drain and close block channel
		func(c chan *domain.BlockEvent) {
			for {
				select {
				case block := <-c:
					log.WithFields(log.Fields{"tx": block.Block.Hash}).Info("gracefully draining block")
				default:
					close(c)
					return
				}
			}
		}(t.blockOutput)
	}
	return nil
}

func (t *TxStreamService) Name() string {
	return "tx-stream"
}

// Health implements health.Reporter interface.
func (t *TxStreamService) Health() health.Reports {
	return health.Reports{
		t.lastBlockActivity.GetReport("event.block.time"),
		t.lastTxActivity.GetReport("event.transaction.time"),
	}
}

func NewTxStreamService(ctx context.Context, ethClient ethereum.Client, blockFeed feeds.BlockFeed, cfg TxStreamServiceConfig) (*TxStreamService, error) {
	txOutput := make(chan *domain.TransactionEvent)
	blockOutput := make(chan *domain.BlockEvent)

	txFeed, err := feeds.NewTransactionFeed(ctx, ethClient, blockFeed, cfg.SkipBlocksOlderThan, 10)
	if err != nil {
		return nil, err
	}

	return &TxStreamService{
		cfg:         cfg,
		ctx:         ctx,
		blockOutput: blockOutput,
		txOutput:    txOutput,
		txFeed:      txFeed,
	}, nil
}
