package publisher

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
)

func initListener(ctx context.Context, cfg config.Config) (*query.AlertListener, error) {
	natsHost := os.Getenv(config.EnvNatsHost)
	if natsHost == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvNatsHost)
	}
	mc := messaging.NewClient("metrics", fmt.Sprintf("%s:%s", natsHost, config.DefaultNatsPort))

	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}
	return query.NewAlertListener(ctx, mc, query.AlertListenerConfig{
		Port:            8770,
		ChainID:         cfg.Scanner.ChainID,
		Key:             key,
		PublisherConfig: cfg.Query.PublishTo,
	})
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {

	listener, err := initListener(ctx, cfg)
	if err != nil {
		log.Errorf("Error while initializing Listener: %s", err.Error())
		return nil, err
	}

	return []services.Service{
		listener,
	}, nil
}

func Run() {
	services.ContainerMain("publisher", initServices)
}
