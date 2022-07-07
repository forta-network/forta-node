package jwt_provider

import (
	"context"
	
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	botjwt "github.com/forta-network/forta-node/services/bot-jwt"
)

func initJWTProvider(cfg *botjwt.JWTProviderConfig) (*botjwt.JWTProvider, error) {
	return botjwt.NewBotJWTProvider(cfg)
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	// can't dial localhost - need to dial host gateway from container
	cfg.Scan.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
	cfg.JsonRpcProxy.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.JsonRpcProxy.JsonRpc.Url)
	
	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}
	
	proxy, err := initJWTProvider(&botjwt.JWTProviderConfig{Key: key})
	if err != nil {
		return nil, err
	}
	
	return []services.Service{
		health.NewService(
			ctx, "", healthutils.DefaultHealthServerErrHandler,
			health.CheckerFrom(summarizeReports, proxy),
		),
		proxy,
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
