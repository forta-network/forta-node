package updater

import (
	"context"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/store"
)

// Updater receives and publishes the latest Forta node image updates.
type Updater struct {
	imgStore  store.FortaImagesStore
	msgClient clients.MessageClient
}

// Start starts the service.
func (up *Updater) Start() error {
	go up.receive()
	return nil
}

// Name returns the name of the service.
func (up *Updater) Name() string {
	return "updater"
}

// Stop stops the service
func (up *Updater) Stop() error {
	return nil
}

// NewUpdater creates a new updater.
func NewUpdater(ctx context.Context, imgStore store.FortaImagesStore, msgClient clients.MessageClient) *Updater {
	return &Updater{
		imgStore:  imgStore,
		msgClient: msgClient,
	}
}

func (up *Updater) receive() {
	for latestRefs := range up.imgStore.Latest() {
		up.msgClient.Publish(messaging.SubjectImagesLatest, latestRefs)
	}
}
