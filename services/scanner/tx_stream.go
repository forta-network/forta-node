package scanner

import (
	"context"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/domain"
	"github.com/forta-protocol/forta-node/ethereum"
	"github.com/forta-protocol/forta-node/feeds"

	log "github.com/sirupsen/logrus"
)

// TxStreamService pulls TX info from providers and emits to channel
type TxStreamService struct {
	cfg         TxStreamServiceConfig
	ctx         context.Context
	blockOutput chan *domain.BlockEvent
	txOutput    chan *domain.TransactionEvent
	txFeed      feeds.TransactionFeed
}

type TxStreamServiceConfig struct {
	JsonRpcConfig      config.JsonRpcConfig
	TraceJsonRpcConfig config.JsonRpcConfig
}

func (t *TxStreamService) ReadOnlyBlockStream() <-chan *domain.BlockEvent {
	return t.blockOutput
}

func (t *TxStreamService) ReadOnlyTxStream() <-chan *domain.TransactionEvent {
	return t.txOutput
}

func (t *TxStreamService) handleBlock(evt *domain.BlockEvent) error {
	t.blockOutput <- evt
	return nil
}

func (t *TxStreamService) handleTx(evt *domain.TransactionEvent) error {
	t.txOutput <- evt
	return nil
}

func (t *TxStreamService) Start() error {
	log.Infof("Starting %s", t.Name())
	go func() {
		if err := t.txFeed.ForEachTransaction(t.handleBlock, t.handleTx); err != nil {
			log.WithError(err).Error("tx feed error")
		}
	}()
	return nil
}

func (t *TxStreamService) Stop() error {
	log.Infof("Stopping %s", t.Name())
	close(t.txOutput)
	close(t.blockOutput)
	return nil
}

func (t *TxStreamService) Name() string {
	return "TxStream"
}

func NewTxStreamService(ctx context.Context, ethClient ethereum.Client, blockFeed feeds.BlockFeed, cfg TxStreamServiceConfig) (*TxStreamService, error) {
	txOutput := make(chan *domain.TransactionEvent)
	blockOutput := make(chan *domain.BlockEvent)

	txFeed, err := feeds.NewTransactionFeed(ctx, ethClient, blockFeed, 10)
	if err != nil {
		return nil, err
	}

	return &TxStreamService{
		cfg,
		ctx,
		blockOutput,
		txOutput,
		txFeed,
	}, nil
}
