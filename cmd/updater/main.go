package main

import (
	"context"
	"fmt"
	"os"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/updater"
	"github.com/forta-protocol/forta-node/store"
)

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	natsHost := os.Getenv(config.EnvNatsHost)
	if natsHost == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvNatsHost)
	}
	msgClient := messaging.NewClient("updater", fmt.Sprintf("%s:%s", natsHost, config.DefaultNatsPort))

	imgStore, err := store.NewFortaImagesStore(ctx)
	if err != nil {
		return nil, err
	}

	return []services.Service{
		updater.NewUpdater(ctx, imgStore, msgClient),
	}, nil
}

func main() {
	gethlog.Root().SetHandler(gethlog.StdoutHandler)

	services.ContainerMain("updater", initServices)
}
