package jwt_provider

import (
	"context"
	
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	botjwt "github.com/forta-network/forta-node/services/bot-jwt"
)

func initJWTProvider(cfg config.Config) (*botjwt.JWTProvider, error) {
	return botjwt.NewBotJWTProvider(cfg)
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	jwtProvider, err := initJWTProvider(cfg)
	if err != nil {
		return nil, err
	}
	
	return []services.Service{
		health.NewService(
			ctx, "", healthutils.DefaultHealthServerErrHandler,
			health.CheckerFrom(summarizeReports, jwtProvider),
		),
		jwtProvider,
	}, nil
}

func summarizeReports(reports health.Reports) *health.Report {
	summary := health.NewSummary()
	
	apiErr, ok := reports.NameContains("service.jwt-provider.api")
	if ok && len(apiErr.Details) > 0 {
		summary.Addf("last time the api failed with error '%s'.", apiErr.Details)
	}
	
	return summary.Finish()
}

func Run() {
	services.ContainerMain("bot-jwt-provider", initServices)
}
