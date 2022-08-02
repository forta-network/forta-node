package jwt_provider

import (
	"context"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	jwt_provider "github.com/forta-network/forta-node/services/jwt-provider"
)

func initJWTProvider(cfg config.Config) (*jwt_provider.JWTProvider, error) {
	return jwt_provider.NewJWTProvider(cfg)
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	jwtProvider, err := initJWTProvider(cfg)
	if err != nil {
		return nil, err
	}

	return []services.Service{
		health.NewService(
			ctx, "", healthutils.DefaultHealthServerErrHandler,
			health.CheckerFrom(nil, jwtProvider),
		),
		jwtProvider,
	}, nil
}

func Run() {
	services.ContainerMain("jwt-provider", initServices)
}
