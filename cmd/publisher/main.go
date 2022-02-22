package publisher

import (
	"context"
	"fmt"
	"os"

	"github.com/forta-protocol/forta-node/clients/alertapi"
	"github.com/forta-protocol/forta-node/clients/health"

	"github.com/forta-protocol/forta-node/clients/messaging"
	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/security"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/publisher"
)

func initPublisher(ctx context.Context, cfg config.Config) (*publisher.Publisher, error) {
	mc := messaging.NewClient("metrics", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}

	releaseInfoStr := os.Getenv(config.EnvReleaseInfo)
	var releaseSummary *config.ReleaseSummary
	if len(releaseInfoStr) > 0 {
		releaseInfo := config.ReleaseInfoFromString(releaseInfoStr)
		releaseSummary = config.MakeSummaryFromReleaseInfo(releaseInfo)
	}

	apiClient := alertapi.NewClient(cfg.Publish.APIURL)

	return publisher.NewPublisher(ctx, mc, apiClient, publisher.PublisherConfig{
		ChainID:         cfg.ChainID,
		Key:             key,
		PublisherConfig: cfg.Publish,
		ReleaseSummary:  releaseSummary,
		Config:          cfg,
	})
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	publisher, err := initPublisher(ctx, cfg)
	if err != nil {
		log.Errorf("Error while initializing Listener: %s", err.Error())
		return nil, err
	}

	return []services.Service{
		health.NewService(ctx, health.CheckerFrom(summarizeReports, publisher)),
		publisher,
	}, nil
}

func summarizeReports(reports health.Reports) *health.Report {
	summary := health.NewSummary()

	batchPublishErr, ok := reports.NameContains("publisher.event.batch-publish.error")
	if ok && len(batchPublishErr.Details) > 0 {
		summary.Addf("failed to publish the last batch with error '%s'", batchPublishErr.Details)
		summary.Status(health.StatusFailing)
	}

	return summary.Finish()
}

func Run() {
	services.ContainerMain("publisher", initServices)
}
