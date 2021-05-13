package main

import (
	"context"

	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/services"
)

func initJsonRpcProxy(ctx context.Context, cfg config.Config) (*services.JsonRpcProxy, error) {
	return services.NewJsonRpcProxy(ctx, cfg)
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	proxy, err := initJsonRpcProxy(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return []services.Service{
		proxy,
	}, nil
}

func main() {
	services.ContainerMain("json-rpc", initServices)
}
