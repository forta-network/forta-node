package sla_tracker

import (
	"context"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/sla-tracker"
)

func initSLATracker(ctx context.Context, cfg sla_tracker.SLATrackerConfig) (*sla_tracker.SLATracker, error) {
	return sla_tracker.NewSLATracker(ctx, cfg)
}

func initServices(ctx context.Context, _ config.Config) ([]services.Service, error) {
	// can't dial localhost - need to dial host gateway from container
	cfg := sla_tracker.SLATrackerConfig{}
	cfg.JSONRpcHost = config.DockerJSONRPCProxyContainerName
	cfg.JSONRpcPort = config.DefaultJSONRPCProxyPort

	proxy, err := initSLATracker(ctx, cfg)
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

	apiErr, ok := reports.NameContains("service.sla-tracker.api")
	if ok && len(apiErr.Details) > 0 {
		summary.Addf("last time the api failed with error '%s'.", apiErr.Details)
	}

	return summary.Finish()
}

func Run() {
	services.ContainerMain("sla-tracker", initServices)
}
