package publisher

import (
	"context"
	"fmt"
	"github.com/forta-protocol/forta-node/clients/alertapi"
	"os"

	"github.com/forta-protocol/forta-node/clients/messaging"
	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/security"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/publisher"
)

func initListener(ctx context.Context, cfg config.Config) (*publisher.Publisher, error) {
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
	})
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {

	listener, err := initListener(ctx, cfg)
	if err != nil {
		log.Errorf("Error while initializing Listener: %s", err.Error())
		return nil, err
	}

	return []services.Service{
		listener,
	}, nil
}

func Run() {
	services.ContainerMain("publisher", initServices)
}
