package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/zephyr-node/config"
	"OpenZeppelin/zephyr-node/services"
)

func initTxStream(ctx context.Context) (*services.TxStreamService, error) {
	url := os.Getenv(services.EnvJsonRpcUrl)
	startBlock := os.Getenv(services.EnvStartBlock)
	var sb *big.Int
	if startBlock != "" {
		sbVal, err := strconv.ParseInt(startBlock, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%s must be numeric", services.EnvStartBlock)
		}
		sb = big.NewInt(sbVal)
	}
	if url == "" {
		return nil, fmt.Errorf("%s is a required env var", services.EnvJsonRpcUrl)
	}
	return services.NewTxStreamService(ctx, services.TxStreamServiceConfig{
		Url:        url,
		StartBlock: sb,
	})
}

func initServices(ctx context.Context) ([]services.Service, error) {
	txStream, err := initTxStream(ctx)
	if err != nil {
		return nil, err
	}
	txAnalyzer := services.NewTxAnalyzerService(ctx, services.TxAnalyzerServiceConfig{
		TxChannel: txStream.ReadOnlyStream(),
	})
	return []services.Service{
		txStream,
		txAnalyzer,
		services.NewTxLogger(ctx),
	}, nil
}

func main() {
	logLevel := os.Getenv(config.EnvLogLevel)
	if logLevel == "" {
		logLevel = "info"
	}
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Error("could not initialize log level", err)
		return
	}
	log.SetLevel(lvl)
	log.Info("Starting Node")

	ctx, cancel := services.InitMainContext()
	defer cancel()

	serviceList, err := initServices(ctx)
	if err != nil {
		log.Error("could not initialize services", err)
		return
	}

	if err := services.StartServices(ctx, serviceList); err != nil {
		log.Error("error running services", err)
	}

	log.Info("Stopping Node")
}
