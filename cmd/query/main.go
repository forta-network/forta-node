package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"forta-network/forta-node/config"
	"forta-network/forta-node/security"
	"forta-network/forta-node/services"
	"forta-network/forta-node/services/query"
	"forta-network/forta-node/store"
)

func initApi(ctx context.Context, as store.AlertStore, cfg config.Config) (*query.AlertApi, error) {
	return query.NewAlertApi(ctx, as, query.AlertApiConfig{Port: 80})
}

func initListener(ctx context.Context, as store.AlertStore, cfg config.Config) (*query.AlertListener, error) {
	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}
	return query.NewAlertListener(ctx, as, query.AlertListenerConfig{
		Port:            8770,
		ChainID:         cfg.Scanner.ChainID,
		Key:             key,
		PublisherConfig: cfg.Query.PublishTo,
	})
}

func initPruner(ctx context.Context, as store.AlertStore, cfg config.Config) (*query.DBPruner, error) {
	return query.NewDBPruner(ctx, as)
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	as, err := store.NewBadgerAlertStore()
	if err != nil {
		log.Errorf("Error while initializing BadgerDB: %s", err.Error())
		return nil, err
	}

	api, err := initApi(ctx, as, cfg)
	if err != nil {
		log.Errorf("Error while initializing API: %s", err.Error())
		return nil, err
	}

	listener, err := initListener(ctx, as, cfg)
	if err != nil {
		log.Errorf("Error while initializing Listener: %s", err.Error())
		return nil, err
	}

	pruner, err := initPruner(ctx, as, cfg)
	if err != nil {
		log.Errorf("Error while initializing Pruner: %s", err.Error())
		return nil, err
	}

	return []services.Service{
		api,
		listener,
		pruner,
	}, nil
}

func main() {
	services.ContainerMain("query", initServices)
}
