package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/store"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var updateInterval = 1 * time.Minute

// UpdaterService receives the release updates.
type UpdaterService struct {
	ctx  context.Context
	port string

	mu     sync.RWMutex
	ipfs   store.IPFSClient
	us     store.UpdaterStore
	server *http.Server

	developmentMode bool
	noUpdate        bool

	latestReference string
	latestRelease   *config.ReleaseManifest
}

// NewUpdaterService creates a new updater service.
func NewUpdaterService(ctx context.Context, us store.UpdaterStore, ipfs store.IPFSClient,
	port string, developmentMode, noUpdate bool,
) *UpdaterService {
	return &UpdaterService{
		ctx:             ctx,
		port:            port,
		us:              us,
		ipfs:            ipfs,
		developmentMode: developmentMode,
		noUpdate:        noUpdate,
	}
}

func (updater *UpdaterService) handleGetVersion(w http.ResponseWriter, r *http.Request) {
	updater.mu.RLock()
	defer updater.mu.RUnlock()

	if updater.latestRelease == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{
		"release": updater.latestReference,
	}).Info("release response")

	b, _ := json.Marshal(&config.ReleaseInfo{
		IPFS:     updater.latestReference,
		Manifest: *updater.latestRelease,
	})
	w.Write(b)
}

// Start starts the service.
func (updater *UpdaterService) Start() error {
	updater.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", updater.port),
		Handler: http.HandlerFunc(updater.handleGetVersion),
	}

	grp, ctx := errgroup.WithContext(updater.ctx)
	grp.Go(func() error {
		return updater.server.ListenAndServe()
	})

	grp.Go(func() error {
		if updater.noUpdate {
			return nil
		}

		t := time.NewTicker(updateInterval)
		for {
			select {
			case <-ctx.Done():
				log.WithError(ctx.Err()).Info("updater context is done")
				updater.stopServer()
				return ctx.Err()
			case <-t.C:
				if err := updater.updateLatestRelease(); err != nil {
					log.WithError(err).Error("error getting release")
					// continue, wait ticker
				}
			}
		}
	})

	if !updater.noUpdate {
		//initialize at start
		if err := updater.updateLatestRelease(); err != nil {
			log.WithError(err).Error("error initializing release")
			return err
		}
	}

	log.Info("updater initialization complete")
	if err := grp.Wait(); err != nil {
		log.WithError(err).Error("error returned while running updater")
		return err
	}
	return nil
}

func (updater *UpdaterService) updateLatestRelease() error {
	if updater.developmentMode {
		return updater.readLocalReleaseManifest()
	}

	log.Info("updating latest release")

	ref, err := updater.us.GetLatestReference()
	if err != nil {
		return err
	}
	if ref != updater.latestReference {
		rm, err := updater.ipfs.GetReleaseManifest(ref)
		if err != nil {
			log.WithError(err).Error("error getting release manifest")
			return err
		}
		updater.mu.Lock()
		defer updater.mu.Unlock()
		updater.latestRelease = rm
		updater.latestReference = ref
		log.WithFields(log.Fields{
			"release": ref,
		}).Info("updating to release")
	} else {
		log.WithFields(log.Fields{
			"release": ref,
		}).Info("no change to release")
	}
	return nil
}

func (updater *UpdaterService) readLocalReleaseManifest() error {
	b, err := ioutil.ReadFile(path.Join(config.DefaultContainerFortaDirPath, "test-release.json"))
	if err != nil {
		log.WithError(err).Info("could not read the test release manifest file - ignoring error")
		return nil
	}
	var release config.ReleaseManifest
	if err := json.Unmarshal(b, &release); err != nil {
		log.WithError(err).Info("could not unmarshal the test release manifest - ignoring error")
		return nil
	}
	updater.mu.Lock()
	defer updater.mu.Unlock()
	updater.latestReference = "test-release.json"
	updater.latestRelease = &release
	return nil
}

// Name returns the name of the service.
func (updater *UpdaterService) Name() string {
	return "updater"
}

func (updater *UpdaterService) stopServer() error {
	log.Info("stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := updater.server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("error stopping server (ignored)")
	}
	return nil
}

// Stop stops the service
func (updater *UpdaterService) Stop() error {
	log.Infof("stopping %s", updater.Name())
	return updater.stopServer()
}
