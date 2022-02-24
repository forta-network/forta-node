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

	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/store"
	"github.com/forta-protocol/forta-node/utils"
	log "github.com/sirupsen/logrus"
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

	latestReference string
	latestRelease   *config.ReleaseManifest

	lastChecked health.TimeTracker
	lastErr     health.ErrorTracker
}

// NewUpdaterService creates a new updater service.
func NewUpdaterService(ctx context.Context, us store.UpdaterStore, ipfs store.IPFSClient,
	port string, developmentMode bool,
) *UpdaterService {
	return &UpdaterService{
		ctx:             ctx,
		port:            port,
		us:              us,
		ipfs:            ipfs,
		developmentMode: developmentMode,
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

	if err := updater.updateLatestRelease(); err != nil {
		log.WithError(err).Error("error initializing release")
		return err
	}

	utils.GoListenAndServe(updater.server)

	go func() {
		t := time.NewTicker(updateInterval)
		for {
			select {
			case <-updater.ctx.Done():
				log.WithError(updater.ctx.Err()).Info("updater context is done")
				updater.stopServer()
				return
			case <-t.C:
				err := updater.updateLatestRelease()
				updater.lastErr.Set(err)
				updater.lastChecked.Set()
				if err != nil {
					log.WithError(err).Error("error getting release")
					// continue, wait ticker
				}
			}
		}
	}()

	log.Info("updater initialization complete")
	return nil
}

func (updater *UpdaterService) updateLatestRelease() error {
	if updater.developmentMode {
		return updater.readLocalReleaseManifest()
	}

	log.Info("updating latest release")

	ref, err := updater.us.GetLatestReference()
	if err != nil {
		return fmt.Errorf("failed to get the latest release manifest ref: %v", err)
	}
	if ref != updater.latestReference {
		rm, err := updater.ipfs.GetReleaseManifest(ref)
		if err != nil {
			log.WithError(err).Error("error getting release manifest")
			return fmt.Errorf("failed while downloading the release manifest: %v", err)
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
	if updater.server == nil {
		return nil
	}
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

// Health implements the health.Reporter interface.
func (updater *UpdaterService) Health() health.Reports {
	return health.Reports{
		updater.lastChecked.GetReport("event.checked.time"),
		updater.lastErr.GetReport("event.checked.error"),
	}
}
