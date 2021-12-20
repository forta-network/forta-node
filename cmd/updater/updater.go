package updater

import (
	"context"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/updater"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	return []services.Service{
		updater.NewUpdaterService(ctx, config.DefaultUpdaterPort),
	}, nil
}

func Run() {
	services.ContainerMain("updater", initServices)
}
