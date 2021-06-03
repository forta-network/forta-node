package main

import (
	"context"

	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/services"
	jrp "OpenZeppelin/fortify-node/services/json-rpc"
)

func initJsonRpcProxy(ctx context.Context, cfg config.Config) (*jrp.JsonRpcProxy, error) {
	return jrp.NewJsonRpcProxy(ctx, cfg)
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
