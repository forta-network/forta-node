package main

import (
	"context"
	"fmt"
	"github.com/forta-protocol/forta-node/clients/messaging"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/security"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/query"
	"github.com/forta-protocol/forta-node/store"
)

func initApi(ctx context.Context, as store.AlertStore, cfg config.Config) (*query.AlertApi, error) {
	return query.NewAlertApi(ctx, as, query.AlertApiConfig{Port: 80})
}

func initListener(ctx context.Context, as store.AlertStore, cfg config.Config) (*query.AlertListener, error) {
	natsHost := os.Getenv(config.EnvNatsHost)
	if natsHost == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvNatsHost)
	}
	mc := messaging.NewClient("metrics", fmt.Sprintf("%s:%s", natsHost, config.DefaultNatsPort))

	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}
	return query.NewAlertListener(ctx, as, mc, query.AlertListenerConfig{
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
