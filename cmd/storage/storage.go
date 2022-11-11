package storage

import (
	"context"
	"fmt"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/storage"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	service, err := storage.NewStorage(
		ctx, fmt.Sprintf("http://%s:5001", config.DockerIpfsContainerName),
		cfg.StorageConfig.Provide,
	)
	if err != nil {
		return nil, err
	}

	return []services.Service{
		health.NewService(
			ctx, "", healthutils.DefaultHealthServerErrHandler,
			health.CheckerFrom(nil, service),
		),
		service,
	}, nil
}

func Run() {
	services.ContainerMain("storage", initServices)
}
