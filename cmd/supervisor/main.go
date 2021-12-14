package supervisor

import (
	"context"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/containers"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	svc, err := containers.NewTxNodeService(ctx, containers.TxNodeServiceConfig{
		Config:     cfg,
		Passphrase: cfg.Passphrase,
	})
	if err != nil {
		return nil, err
	}
	return []services.Service{
		svc,
	}, nil
}

func Run() {
	services.ContainerMain("supervisor", initServices)
}
