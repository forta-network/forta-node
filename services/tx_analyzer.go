package services

import (
	"context"

	"github.com/golang/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"

	"OpenZeppelin/safe-node/feeds"
)

// TxAnalyzerService reads TX info, calls agents, and emits results
type TxAnalyzerService struct {
	cfg TxAnalyzerServiceConfig
	ctx context.Context
}

type TxAnalyzerServiceConfig struct {
	TxChannel <-chan *feeds.TransactionEvent
}

func (t *TxAnalyzerService) Start() error {
	log.Infof("Starting %s", t.Name())
	for tx := range t.cfg.TxChannel {
		//log.Infof("%s, %d, %s", tx.BlockEvent.EventType, tx.BlockEvent.Block.NumberU64(), tx.Transaction.Hash().Hex())
		msg, err := tx.ToMessage()
		if err != nil {
			return err
		}
		jm := jsonpb.Marshaler{}
		str, err := jm.MarshalToString(msg)
		if err != nil {
			return err
		}
		log.Info(str)
	}
	return nil
}

func (t *TxAnalyzerService) Stop() error {
	log.Infof("Stopping %s", t.Name())
	return nil
}

func (t *TxAnalyzerService) Name() string {
	return "TxAnalyzerService"
}

func NewTxAnalyzerService(ctx context.Context, cfg TxAnalyzerServiceConfig) *TxAnalyzerService {
	return &TxAnalyzerService{
		cfg,
		ctx,
	}
}
