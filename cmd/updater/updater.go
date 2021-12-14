package updater

import (
	"context"
	"fmt"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/updater"
	"github.com/forta-protocol/forta-node/store"
	log "github.com/sirupsen/logrus"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	imgStore, err := store.NewFortaImageStore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create the image store: %v", err)
	}
	dockerClient, err := clients.NewDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create the docker client: %v", err)
	}

	return []services.Service{
		updater.NewUpdater(ctx, cfg, imgStore, dockerClient),
	}, nil
}

// Run runs the updater.
func Run(cfg config.Config) {
	ctx, cancel := services.InitMainContext()
	defer cancel()

	log.Info("starting updater")

	serviceList, err := initServices(ctx, cfg)
	if err != nil {
		log.WithError(err).Error("could not initialize updater services")
		return
	}

	if err := services.StartServices(ctx, serviceList); err != nil {
		log.WithError(err).Error("error running updater services")
	}

	log.Info("stopping updater")
}
