package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"OpenZeppelin/zephyr-node/config"
	"OpenZeppelin/zephyr-node/services"
)

func initTxStream(ctx context.Context, cfg config.Config) (*services.TxStreamService, error) {
	url := cfg.Scanner.Ethereum.JsonRpcUrl
	startBlock := cfg.Scanner.Ethereum.StartBlock
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
	qn := os.Getenv(config.EnvQueryNode)
	if qn == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvQueryNode)
	}
	return services.NewTxAnalyzerService(ctx, services.TxAnalyzerServiceConfig{
		TxChannel:     stream.ReadOnlyStream(),
		AgentConfigs:  cfg.Agents,
		QueryNodeAddr: qn,
	})
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
	services.ContainerMain("scanner", initServices)
}
