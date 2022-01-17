package runner

import (
	"context"
	"fmt"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/runner"
	"github.com/forta-protocol/forta-node/store"
	log "github.com/sirupsen/logrus"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	imgStore, err := store.NewFortaImageStore(ctx, config.DefaultUpdaterPort)
	if err != nil {
		return nil, fmt.Errorf("failed to create the image store: %v", err)
	}
	dockerClient, err := clients.NewDockerClient("runner")
	if err != nil {
		return nil, fmt.Errorf("failed to create the docker client: %v", err)
	}
	globalDockerClient, err := clients.NewDockerClient("")
	if err != nil {
		return nil, fmt.Errorf("failed to create the docker client: %v", err)
	}

	if cfg.Development {
		log.Warn("running in development mode")
	}

	return []services.Service{
		runner.NewRunner(ctx, cfg, imgStore, dockerClient, globalDockerClient, config.DefaultUpdaterPort),
	}, nil
}

// Run runs the runner.
func Run(cfg config.Config) {
	ctx, cancel := services.InitMainContext()
	defer cancel()

	log.Info("starting runner")

	serviceList, err := initServices(ctx, cfg)
	if err != nil {
		log.WithError(err).Error("could not initialize runner services")
		return
	}

	if err := services.StartServices(ctx, cancel, serviceList); err != nil {
		log.WithError(err).Error("error running runner services")
	}

	log.Info("stopping runner")
}
