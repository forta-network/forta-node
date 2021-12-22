package updater

import (
	"context"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/updater"
	"github.com/forta-protocol/forta-node/store"
	"github.com/forta-protocol/forta-node/utils"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {

	ipfs := store.NewIPFSClient(cfg.Registry.IPFS.GatewayURL)
	up, err := store.NewContractUpdaterStore(cfg)
	if err != nil {
		return nil, err
	}
	return []services.Service{
		updater.NewUpdaterService(
			ctx, up, ipfs, config.DefaultUpdaterPort,
			utils.ParseBoolEnvVar(config.EnvDevelopment),
			utils.ParseBoolEnvVar(config.EnvNoUpdate),
		),
	}, nil
}

func Run() {
	services.ContainerMain("updater", initServices)
}
