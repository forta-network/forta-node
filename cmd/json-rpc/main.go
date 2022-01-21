package json_rpc

import (
	"context"

	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	jrp "github.com/forta-protocol/forta-node/services/json-rpc"
	"github.com/forta-protocol/forta-node/utils"
)

func initJsonRpcProxy(ctx context.Context, cfg config.Config) (*jrp.JsonRpcProxy, error) {
	return jrp.NewJsonRpcProxy(ctx, cfg)
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	// can't dial localhost - need to dial host gateway from container
	cfg.Scan.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
	cfg.JsonRpcProxy.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.JsonRpcProxy.JsonRpc.Url)

	proxy, err := initJsonRpcProxy(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return []services.Service{
		proxy,
		health.NewService(ctx, health.CheckerFrom(proxy)),
	}, nil
}

func Run() {
	services.ContainerMain("json-rpc", initServices)
}
