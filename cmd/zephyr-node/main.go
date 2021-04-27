package main

import (
	"context"
	"fmt"
	"math/big"

	"OpenZeppelin/zephyr-node/config"
	"OpenZeppelin/zephyr-node/services"
)

func initTxStream(ctx context.Context, cfg config.Config) (*services.TxStreamService, error) {
	url := cfg.Ethereum.JsonRpcUrl
	startBlock := cfg.Ethereum.StartBlock
	var sb *big.Int
	if startBlock != 0 {
		sb = big.NewInt(int64(startBlock))
	}
	if url == "" {
		return nil, fmt.Errorf("ethereum.jsonRpcUrl is required")
	}
	return services.NewTxStreamService(ctx, services.TxStreamServiceConfig{
		Url:        url,
		StartBlock: sb,
	})
}

func initTxAnalyzer(ctx context.Context, cfg config.Config, stream *services.TxStreamService) (*services.TxAnalyzerService, error) {
	return services.NewTxAnalyzerService(ctx, services.TxAnalyzerServiceConfig{
		TxChannel:      stream.ReadOnlyStream(),
		AgentAddresses: cfg.AgentContainerNames(),
	}), nil
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	txStream, err := initTxStream(ctx, cfg)
	if err != nil {
		return nil, err
	}
	txAnalyzer, err := initTxAnalyzer(ctx, cfg, txStream)
	if err != nil {
		return nil, err
	}
	return []services.Service{
		txStream,
		txAnalyzer,
		services.NewTxLogger(ctx),
	}, nil
}

func main() {
	services.ContainerMain("zephyr-node", initServices)
}
