package json_rpc

import (
	"context"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/components/registry"
	jrp "github.com/forta-network/forta-node/services/json-rpc"
	jrpcache "github.com/forta-network/forta-node/services/json-rpc/cache"
)

func initJsonRpcProxy(ctx context.Context, cfg config.Config) (*jrp.JsonRpcProxy, error) {
	return jrp.NewJsonRpcProxy(ctx, cfg)
}

func initJsonRpcCache(ctx context.Context, cfg config.Config, botRegistry registry.BotRegistry) (*jrpcache.JsonRpcCache, error) {
	return jrpcache.NewJsonRpcCache(ctx, cfg.JsonRpcCache, botRegistry)
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	// can't dial localhost - need to dial host gateway from container
	cfg.Scan.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
	cfg.JsonRpcProxy.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.JsonRpcProxy.JsonRpc.Url)
	cfg.Registry.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Registry.JsonRpc.Url)

	proxy, err := initJsonRpcProxy(ctx, cfg)
	if err != nil {
		return nil, err
	}

	key, err := config.LoadKeyInContainer(cfg)
	if err != nil {
		return nil, err
	}

	botRegistry, err := registry.New(cfg, key.Address)
	if err != nil {
		return nil, err
	}

	cache, err := initJsonRpcCache(ctx, cfg, botRegistry)
	if err != nil {
		return nil, err
	}

	return []services.Service{
		health.NewService(
			ctx, "", healthutils.DefaultHealthServerErrHandler,
			health.CheckerFrom(summarizeReports, proxy),
		),
		proxy,
		cache,
	}, nil
}

func summarizeReports(reports health.Reports) *health.Report {
	summary := health.NewSummary()

	apiErr, ok := reports.NameContains("service.json-rpc-proxy.api")
	if ok && len(apiErr.Details) > 0 {
		summary.Addf("last time the api failed with error '%s'.", apiErr.Details)
	}

	return summary.Finish()
}

func Run() {
	services.ContainerMain("json-rpc", initServices)
}
