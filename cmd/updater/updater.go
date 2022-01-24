package updater

import (
	"context"

	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/updater"
	"github.com/forta-protocol/forta-node/store"
	"github.com/forta-protocol/forta-node/utils"
	log "github.com/sirupsen/logrus"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {

	ipfs := store.NewIPFSClient(cfg.Registry.IPFS.GatewayURL)
	up, err := store.NewContractUpdaterStore(cfg)
	if err != nil {
		return nil, err
	}

	developmentMode := utils.ParseBoolEnvVar(config.EnvDevelopment)
	noUpdate := utils.ParseBoolEnvVar(config.EnvNoUpdate)

	log.WithFields(log.Fields{
		"developmentMode": developmentMode,
		"noUpdate":        noUpdate,
	}).Info("updater modes")

	updaterService := updater.NewUpdaterService(
		ctx, up, ipfs, config.DefaultContainerPort,
		developmentMode, noUpdate,
	)

	return []services.Service{
		updaterService,
		health.NewService(ctx, health.CheckerFrom(updaterService)),
	}, nil
}

func Run() {
	services.ContainerMain("updater", initServices)
}
