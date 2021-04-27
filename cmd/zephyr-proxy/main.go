package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"OpenZeppelin/zephyr-node/config"
	"OpenZeppelin/zephyr-node/services"
)

func initProxies(ctx context.Context, cfg config.Config) (*services.TxProxyService, error) {
	keys := os.Getenv(config.EnvAgentKeys)
	if keys == "" {
		return nil, fmt.Errorf("%s is required", config.EnvAgentKeys)
	}
	return services.NewTxProxyService(ctx, services.TxProxyConfig{
		AgentAddresses: cfg.AgentContainerNames(),
		Keys:           strings.Split(keys, ","),
	})
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	proxies, err := initProxies(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return []services.Service{
		proxies,
	}, nil
}

func main() {
	services.ContainerMain("zephyr-proxy", initServices)
}
