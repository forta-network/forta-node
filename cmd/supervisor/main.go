package supervisor

import (
	"context"
	"strconv"

	"github.com/forta-protocol/forta-core-go/clients/health"
	"github.com/forta-protocol/forta-core-go/security"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/supervisor"
)

const (
	expectedContainersCount = 4
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	passphrase, err := security.ReadPassphrase()
	if err != nil {
		return nil, err
	}
	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}
	svc, err := supervisor.NewSupervisorService(ctx, supervisor.SupervisorServiceConfig{
		Config:     cfg,
		Passphrase: passphrase,
		Key:        key,
	})
	if err != nil {
		return nil, err
	}
	return []services.Service{
		health.NewService(ctx, health.CheckerFrom(summarizeReports, svc)),
		svc,
	}, nil
}

func summarizeReports(reports health.Reports) *health.Report {
	summary := health.NewSummary()

	containersManager, ok := reports.NameContains("containers.managed")
	if ok {
		count, _ := strconv.Atoi(containersManager.Details)
		if count < expectedContainersCount {
			summary.Addf("missing %d containers.", expectedContainersCount-count)
			summary.Status(health.StatusFailing)
		} else {
			summary.Addf("all %d service containers are running.", expectedContainersCount)
		}
	}

	telemetryErr, ok := reports.NameContains("telemetry-sync.error")
	if ok && len(telemetryErr.Details) > 0 {
		summary.Addf("telemetry sync is failing with error '%s' (non-critical).", telemetryErr.Details)
		// do not change status - non critical
	}

	return summary.Finish()
}

func Run() {
	services.ContainerMain("supervisor", initServices)
}
