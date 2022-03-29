package store

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/forta-protocol/forta-core-go/release"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/config"
)

const defaultImageCheckInterval = time.Second * 5

// FortaImageStore keeps track of the latest Forta node image.
type FortaImageStore interface {
	Latest() <-chan ImageRefs
	EmbeddedImageRefs() ImageRefs
}

// ImageRefs contains the latest image references.
type ImageRefs struct {
	Supervisor  string
	Updater     string
	ReleaseInfo *release.ReleaseInfo
}

type fortaImageStore struct {
	updaterPort string
	latestCh    chan ImageRefs
	latestImgs  ImageRefs
}

// NewFortaImageStore creates a new store.
func NewFortaImageStore(ctx context.Context, updaterPort string, autoUpdate bool) (*fortaImageStore, error) {
	store := &fortaImageStore{
		updaterPort: updaterPort,
		latestCh:    make(chan ImageRefs),
	}
	if autoUpdate {
		go store.loop(ctx)
	}
	return store, nil
}

func (store *fortaImageStore) loop(ctx context.Context) {
	store.check(ctx)
	ticker := time.NewTicker(defaultImageCheckInterval)
	for range ticker.C {
		store.check(ctx)
	}
}

func (store *fortaImageStore) EmbeddedImageRefs() ImageRefs {
	return ImageRefs{
		Supervisor:  config.DockerSupervisorImage,
		Updater:     config.DockerUpdaterImage,
		ReleaseInfo: config.GetBuildReleaseInfo(),
	}
}

func (store *fortaImageStore) check(ctx context.Context) {
	latestReleaseInfo, err := store.getFromUpdater(ctx)
	if err != nil {
		log.WithError(err).Warn("failed to get the latest release from the updater")
	}

	if latestReleaseInfo == nil {
		return
	}

	serviceImgs := latestReleaseInfo.Manifest.Release.Services
	if serviceImgs.Supervisor != store.latestImgs.Supervisor || serviceImgs.Updater != store.latestImgs.Updater {
		log.WithField("commit", latestReleaseInfo.Manifest.Release.Commit).Info("got newer release from updater")

		store.latestImgs = ImageRefs{
			Supervisor:  serviceImgs.Supervisor,
			Updater:     serviceImgs.Updater,
			ReleaseInfo: latestReleaseInfo,
		}
		store.latestCh <- store.latestImgs
	}
}

func (store *fortaImageStore) getFromUpdater(ctx context.Context) (*release.ReleaseInfo, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s", store.updaterPort))
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound { // 404 == not ready yet
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected updater response with code %d: %s", resp.StatusCode, string(respBody))
	}
	var releaseInfo release.ReleaseInfo
	return &releaseInfo, json.Unmarshal(respBody, &releaseInfo)
}

// Latest returns a channel that provides the latest image reference.
func (store *fortaImageStore) Latest() <-chan ImageRefs {
	return store.latestCh
}
