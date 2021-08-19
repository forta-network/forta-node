package runner

import (
	"context"
	"forta-network/forta-node/config"
	"forta-network/forta-node/services"
	"forta-network/forta-node/services/containers"

	log "github.com/sirupsen/logrus"
)

func initServices(cfg config.Config, passphrase string, ctx context.Context) ([]services.Service, error) {
	svc, err := containers.NewTxNodeService(ctx, containers.TxNodeServiceConfig{
		Config:     cfg,
		Passphrase: passphrase,
	})
	if err != nil {
		return nil, err
	}
	return []services.Service{
		svc,
	}, nil
}

// Run runs the node.
func Run(cfg config.Config) {
	ctx, cancel := services.InitMainContext()
	defer cancel()

	log.Info("Starting Node")

	serviceList, err := initServices(cfg, cfg.Passphrase, ctx)
	if err != nil {
		log.Error("could not initialize services", err)
		return
	}

	if err := services.StartServices(ctx, serviceList); err != nil {
		log.Error("error running services", err)
	}

	log.Info("Stopping Node")
}
