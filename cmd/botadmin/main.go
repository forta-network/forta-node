package botadmin

import (
	"context"
	"os"

	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/network"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	server, err := network.NewBotAdminServer(os.Getenv(config.EnvContainerName))
	if err != nil {
		return nil, err
	}

	return []services.Service{
		server,
	}, nil
}

func Run() {
	services.ContainerMain("bot-admin", initServices)
}
