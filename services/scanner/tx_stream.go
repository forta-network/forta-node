package scanner

import (
	"context"
	"math/big"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/feeds"
)

// TxStreamService pulls TX info from providers and emits to channel
type TxStreamService struct {
	cfg         TxStreamServiceConfig
	ctx         context.Context
	blockOutput chan *feeds.BlockEvent
	txOutput    chan *feeds.TransactionEvent
	txFeed      feeds.TransactionFeed
}

type TxStreamServiceConfig struct {
	Url        string
	StartBlock *big.Int
}

func (t *TxStreamService) ReadOnlyBlockStream() <-chan *feeds.BlockEvent {
	return t.blockOutput
}

func (t *TxStreamService) ReadOnlyTxStream() <-chan *feeds.TransactionEvent {
	return t.txOutput
}

func (t *TxStreamService) handleBlock(evt *feeds.BlockEvent) error {
	log.Debug("<- TxStream putting block in stream")
	t.blockOutput <- evt
	return nil
}

func (t *TxStreamService) handleTx(evt *feeds.TransactionEvent) error {
	log.Debug("<- TxStream putting tx in stream")
	t.txOutput <- evt
	return nil
}

func (t *TxStreamService) Start() error {
	log.Infof("Starting %s", t.Name())
	defer close(t.txOutput)
	defer close(t.blockOutput)
	return t.txFeed.ForEachTransaction(t.handleBlock, t.handleTx)
}

func (t *TxStreamService) Stop() error {
	log.Infof("Stopping %s", t.Name())
	return nil
}

func (t *TxStreamService) Name() string {
	return "TxStream"
}

func NewTxStreamService(ctx context.Context, cfg TxStreamServiceConfig) (*TxStreamService, error) {
	txOutput := make(chan *feeds.TransactionEvent)
	blockOutput := make(chan *feeds.BlockEvent)

	ethClient, err := clients.NewStreamEthClient(ctx, cfg.Url)
	if err != nil {
		return nil, err
	}
	txFeed, err := feeds.NewTransactionFeed(ctx, ethClient, cfg.StartBlock, 10)
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
