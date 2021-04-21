package services

import (
	"context"
	"math/big"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/safe-node/clients"
	"OpenZeppelin/safe-node/feeds"
)

// TxStreamService pulls TX info from providers and emits to channel
type TxStreamService struct {
	cfg    TxStreamServiceConfig
	ctx    context.Context
	output chan *feeds.TransactionEvent
	txFeed feeds.TransactionFeed
}

type TxStreamServiceConfig struct {
	Url        string
	StartBlock *big.Int
}

func (t *TxStreamService) ReadOnlyStream() <-chan *feeds.TransactionEvent {
	return t.output
}

func (t *TxStreamService) Start() error {
	log.Infof("Starting %s", t.Name())
	defer close(t.output)
	return t.txFeed.ForEachTransaction(func(evt *feeds.TransactionEvent) error {
		log.Debug("<- TxStreamService putting event in stream")
		t.output <- evt
		return nil
	})
}

func (t *TxStreamService) Stop() error {
	log.Infof("Stopping %s", t.Name())
	return nil
}

func (t *TxStreamService) Name() string {
	return "TxStreamService"
}

func NewTxStreamService(ctx context.Context, cfg TxStreamServiceConfig) (*TxStreamService, error) {
	output := make(chan *feeds.TransactionEvent)
	ethClient, err := clients.NewEthClient(ctx, cfg.Url)
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
		output,
		txFeed,
	}, nil
}
