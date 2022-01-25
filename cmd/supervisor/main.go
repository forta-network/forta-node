package supervisor

import (
	"context"

	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/security"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/supervisor"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	passphrase, err := security.ReadPassphrase()
	if err != nil {
		return nil, err
	}
	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}
	svc, err := supervisor.NewSupervisorService(ctx, supervisor.SupervisorServiceConfig{
		Config:     cfg,
		Passphrase: passphrase,
		Key:        key,
	})
	if err != nil {
		return nil, err
	}
	return []services.Service{
		svc,
		health.NewService(ctx, health.CheckerFrom(svc)),
	}, nil
}

func Run() {
	services.ContainerMain("supervisor", initServices)
}
