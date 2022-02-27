package updater

import (
	"context"
	"github.com/forta-protocol/forta-core-go/release"
	"time"

	"github.com/forta-protocol/forta-core-go/clients/health"
	"github.com/forta-protocol/forta-core-go/utils"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/updater"
	"github.com/forta-protocol/forta-node/store"
	log "github.com/sirupsen/logrus"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {

	rc, err := release.NewClient(cfg.Registry.IPFS.GatewayURL)
	if err != nil {
		return nil, err
	}

	up, err := store.NewContractUpdaterStore(cfg)
	if err != nil {
		return nil, err
	}

	developmentMode := utils.ParseBoolEnvVar(config.EnvDevelopment)

	log.WithFields(log.Fields{
		"developmentMode": developmentMode,
	}).Info("updater modes")

	updaterService := updater.NewUpdaterService(
		ctx, up, rc, config.DefaultContainerPort,
		developmentMode,
	)

	return []services.Service{
		health.NewService(ctx, health.CheckerFrom(summarizeReports, updaterService)),
		updaterService,
	}, nil
}

func summarizeReports(reports health.Reports) *health.Report {
	summary := health.NewSummary()

	checkedErr, ok := reports.NameContains("event.checked.error")
	if !ok {
		summary.Fail()
		return summary.Finish()
	}
	if len(checkedErr.Details) > 0 {
		summary.Addf("auto-updater is failing to check new versions with error '%s'", checkedErr.Details)
		summary.Status(health.StatusFailing)
	}

	checkedTime, ok := reports.NameContains("event.checked.time")
	if ok {
		t, ok := checkedTime.Time()
		if ok {
			checkDelay := time.Since(*t)
			if checkDelay > time.Minute*10 {
				summary.Addf("and late for %d minutes", int64(checkDelay.Minutes()))
				summary.Status(health.StatusFailing)
			}
		}
	}
	summary.Punc(".")
	return summary.Finish()
}

func Run() {
	services.ContainerMain("updater", initServices)
}
