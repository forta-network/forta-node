package services

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

type TxLogger struct {
	ctx context.Context
}

func (t *TxLogger) Start() error {
	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		if t.ctx.Err() != nil {
			return t.ctx.Err()
		}
		log.Info("tx-logger tick")
	}
	return nil
}

func (t *TxLogger) Stop() error {
	log.Infof("Stopping %s", t.Name())
	return nil
}

func (t *TxLogger) Name() string {
	return "TxLogger"
}

func NewTxLogger(ctx context.Context) *TxLogger {
	return &TxLogger{
		ctx: ctx,
	}
}
