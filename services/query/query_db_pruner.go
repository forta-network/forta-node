package query

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/store"
)

// DBPruner periodically prunes db of old records
type DBPruner struct {
	ctx    context.Context
	store  store.AlertStore
	cancel context.CancelFunc
}

func (t *DBPruner) Start() error {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		if t.ctx.Err() != nil {
			return t.ctx.Err()
		}
		log.Info("start db prune operation")
		if err := t.store.Prune(); err != nil {
			if err == store.ErrNoPruneNeeded {
				log.Info(err.Error())
			} else {
				log.Errorf("unexpected error while pruning db: %s", err.Error())
			}
		}
		log.Info("end db prune operation")
	}
	return nil
}

func (t *DBPruner) Stop() error {
	defer t.cancel()
	log.Infof("Stopping %s", t.Name())
	return nil
}

func (t *DBPruner) Name() string {
	return "DBPruner"
}

func NewDBPruner(ctx context.Context, s store.AlertStore) (*DBPruner, error) {
	ctx, cancel := context.WithCancel(ctx)
	return &DBPruner{
		ctx:    ctx,
		store:  s,
		cancel: cancel,
	}, nil
}
