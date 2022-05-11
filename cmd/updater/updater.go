package updater

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-network/forta-core-go/security"
	"math/big"
	"time"

	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/release"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/updater"
	"github.com/forta-network/forta-node/store"
	log "github.com/sirupsen/logrus"
)

const minUpdateInterval = 1 * time.Minute
const maxUpdateInterval = 1 * time.Hour

func generateIntervalMs(addr common.Address) int64 {
	interval := big.NewInt(0)
	interval.Mod(utils.ScannerIDHexToBigInt(addr.Hex()), big.NewInt((maxUpdateInterval).Milliseconds()))
	return interval.Int64() + minUpdateInterval.Milliseconds()
}

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

	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}

	intervalMs := generateIntervalMs(key.Address)
	updateDelay := int(intervalMs / 1000)
	if cfg.AutoUpdate.UpdateDelay != nil {
		updateDelay = *cfg.AutoUpdate.UpdateDelay
	}

	updaterService := updater.NewUpdaterService(
		ctx, rg, rc, config.DefaultContainerPort,
		developmentMode, updateDelay,
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
