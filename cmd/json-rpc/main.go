package json_rpc

import (
	"context"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	jrp "github.com/forta-protocol/forta-node/services/json-rpc"
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

func Run() {
	services.ContainerMain("json-rpc", initServices)
}
