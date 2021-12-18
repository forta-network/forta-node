package publisher

import (
	"context"
	"fmt"
	"os"

	"github.com/forta-protocol/forta-node/clients/messaging"
	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/security"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/publisher"
)

func initListener(ctx context.Context, cfg config.Config) (*publisher.Publisher, error) {
	natsHost := os.Getenv(config.EnvNatsHost)
	if natsHost == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvNatsHost)
	}
	mc := messaging.NewClient("metrics", fmt.Sprintf("%s:%s", natsHost, config.DefaultNatsPort))

	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}

	log.Info("listener init complete")
	return publisher.NewPublisher(ctx, mc, publisher.PublisherConfig{
		ChainID:         cfg.ChainID,
		Key:             key,
		PublisherConfig: cfg.Publish,
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
