package store

import (
	"context"
	"time"

	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
)

const defaultImageCheckInterval = time.Minute

// FortaImagesStore keeps track of the latest versions of Forta node container images.
type FortaImagesStore interface {
	Latest() <-chan config.FortaImages
}

type fortaImagesStore struct {
	latestCh chan config.FortaImages
}

// NewFortaImagesStore creates a new store.
func NewFortaImagesStore(ctx context.Context) (*fortaImagesStore, error) {
	store := &fortaImagesStore{
		latestCh: make(chan config.FortaImages),
	}
	go store.loop(ctx)
	store.attachTestSource()
	return store, nil
}

func (store *fortaImagesStore) loop(ctx context.Context) {
	ticker := time.NewTicker(defaultImageCheckInterval)
	for range ticker.C {
		store.check(ctx)
	}
}

func (store *fortaImagesStore) check(ctx context.Context) {
	// TODO: Listen from the contract, compare with known images, replace and push to the channel.
}

func (store *fortaImagesStore) attachTestSource() {
	msgClient := messaging.NewClient("containerstore", "forta-nats:4222")
	msgClient.Subscribe("images.test", messaging.ImagesHandler(func(payload messaging.ImagesPayload) error {
		store.latestCh <- (config.FortaImages)(payload)
		return nil
	}))
}

// Latest returns a channel that provides the latest image references.
func (store *fortaImagesStore) Latest() <-chan config.FortaImages {
	return store.latestCh
}
