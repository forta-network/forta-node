package main

import (
	"context"
	"flag"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/safe-node/config"
	"OpenZeppelin/safe-node/services"
)

func initServices(cfg config.Config, ctx context.Context) ([]services.Service, error) {
	svc, err := services.NewTxNodeService(ctx, services.TxNodeServiceConfig{
		JsonRpcUrl: cfg.Ethereum.JsonRpcUrl,
		LogLevel:   cfg.Log.Level,
	})
	if err != nil {
		return nil, err
	}
	return []services.Service{
		svc,
	}, nil
}

func main() {

	ctx, cancel := services.InitMainContext()
	defer cancel()

	cfgFile := flag.String("config", "config.yml", "filename for configuration yaml")

	flag.Parse()

	cfg, err := config.GetConfig(*cfgFile)
	if err != nil {
		log.Error("could not read config file", err)
		return
	}
	if err := config.InitLogLevel(cfg); err != nil {
		log.Error("error initializing log level", err)
		return
	}

	log.Info("Starting Node")

	serviceList, err := initServices(cfg, ctx)
	if err != nil {
		log.Error("could not initialize services", err)
		return
	}

	if err := services.StartServices(ctx, serviceList); err != nil {
		log.Error("error running services", err)
	}

	log.Info("Stopping Node")
}
