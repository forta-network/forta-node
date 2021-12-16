package supervisor

import (
	"context"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/supervisor"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	svc, err := supervisor.NewSupervisorService(ctx, supervisor.SupervisorServiceConfig{
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
