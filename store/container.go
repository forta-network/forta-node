package store

import (
	"context"
	"time"

	"github.com/forta-protocol/forta-node/config"
)

const defaultImageCheckInterval = time.Minute

// FortaImageStore keeps track of the latest Forta node image.
type FortaImageStore interface {
	Latest() <-chan string
}

type fortaImageStore struct {
	latestCh  chan string
	latestImg string
}

// NewFortaImageStore creates a new store.
func NewFortaImageStore(ctx context.Context) (*fortaImageStore, error) {
	store := &fortaImageStore{
		latestCh: make(chan string),
	}
	go store.loop(ctx)
	return store, nil
}

func (store *fortaImageStore) loop(ctx context.Context) {
	store.check(ctx)
	ticker := time.NewTicker(defaultImageCheckInterval)
	for range ticker.C {
		store.check(ctx)
	}
}

func (store *fortaImageStore) check(ctx context.Context) {
	// TODO: Improve this later to check a contract or something.
	if len(store.latestImg) > 0 {
		return
	}
	store.latestImg = config.DockerScannerNodeImage
	store.latestCh <- config.DockerScannerNodeImage
}

// Latest returns a channel that provides the latest image reference.
func (store *fortaImageStore) Latest() <-chan string {
	return store.latestCh
}
