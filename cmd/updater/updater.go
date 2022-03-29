package updater

import (
	"context"
	"time"

	"github.com/forta-protocol/forta-core-go/registry"
	"github.com/forta-protocol/forta-core-go/release"

	"github.com/forta-protocol/forta-core-go/clients/health"
	"github.com/forta-protocol/forta-core-go/utils"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/healthutils"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/updater"
	"github.com/forta-protocol/forta-node/store"
	log "github.com/sirupsen/logrus"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	cfg.Registry.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Registry.JsonRpc.Url)
	cfg.Registry.IPFS.APIURL = utils.ConvertToDockerHostURL(cfg.Registry.IPFS.APIURL)
	cfg.Registry.IPFS.GatewayURL = utils.ConvertToDockerHostURL(cfg.Registry.IPFS.GatewayURL)

	rc, err := release.NewClient(cfg.Registry.IPFS.GatewayURL)
	if err != nil {
		return nil, err
	}
	rg, err := store.GetRegistryClient(ctx, cfg, registry.ClientConfig{
		JsonRpcUrl: cfg.Registry.JsonRpc.Url,
		ENSAddress: cfg.ENSConfig.ContractAddress,
		Name:       "updater",
	})
	if err != nil {
		return nil, err
	}

	developmentMode := utils.ParseBoolEnvVar(config.EnvDevelopment)

	log.WithFields(log.Fields{
		"developmentMode": developmentMode,
	}).Info("updater modes")

	updaterService := updater.NewUpdaterService(
		ctx, rg, rc, config.DefaultContainerPort,
		developmentMode,
	)

	return []services.Service{
		health.NewService(
			ctx, "", healthutils.DefaultHealthServerErrHandler,
			health.CheckerFrom(summarizeReports, updaterService),
		),
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
