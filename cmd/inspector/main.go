package inspector

import (
	"context"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/inspector"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	inspector, err := inspector.NewInspector(ctx, inspector.InspectorConfig{
		Config:    cfg,
		ProxyHost: config.DockerJSONRPCProxyContainerName,
		ProxyPort: config.DefaultJSONRPCProxyPort,
	})
	if err != nil {
		return nil, err
	}

	return []services.Service{
		health.NewService(
			ctx, "", healthutils.DefaultHealthServerErrHandler,
			health.CheckerFrom(summarizeReports, inspector),
		),
		inspector,
	}, nil
}

func summarizeReports(reports health.Reports) *health.Report {
	summary := health.NewSummary()
	return summary.Finish()
}

func Run() {
	services.ContainerMain("inspector", initServices)
}
