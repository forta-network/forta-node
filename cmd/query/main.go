package main

import (
	"context"
	"fmt"
	"os"

	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/services"
	"OpenZeppelin/fortify-node/services/query"
	"OpenZeppelin/fortify-node/store"
)

func initApi(ctx context.Context, as store.AlertStore, cfg config.Config) (*query.AlertApi, error) {
	return query.NewAlertApi(ctx, as, query.AlertApiConfig{Port: 80})
}

func initListener(ctx context.Context, as store.AlertStore, cfg config.Config) (*query.AlertListener, error) {
	natsHost := os.Getenv(config.EnvNatsHost)
	if natsHost == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvNatsHost)
	}
	msgClient := messaging.NewClient("query", fmt.Sprintf("%s:%s", natsHost, config.DefaultNatsPort))
	return query.NewAlertListener(ctx, as, query.AlertListenerConfig{Port: 8770}, msgClient)
}

func initPruner(ctx context.Context, as store.AlertStore, cfg config.Config) (*query.DBPruner, error) {
	return query.NewDBPruner(ctx, as)
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	as, err := store.NewBadgerAlertStore()
	if err != nil {
		return nil, err
	}

	api, err := initApi(ctx, as, cfg)
	if err != nil {
		return nil, err
	}

	listener, err := initListener(ctx, as, cfg)
	if err != nil {
		return nil, err
	}

	pruner, err := initPruner(ctx, as, cfg)
	if err != nil {
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
