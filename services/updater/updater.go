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

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/release"
	"github.com/forta-network/forta-core-go/utils"

	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
)

// UpdaterService receives the release updates.
type UpdaterService struct {
	ctx  context.Context
	port string

	mu     sync.RWMutex
	rl     release.Client
	rg     registry.Client
	server *http.Server

	developmentMode bool

	latestReference string
	latestRelease   *release.ReleaseManifest

	delaySeconds int

	lastChecked health.TimeTracker
	lastErr     health.ErrorTracker
}

// NewUpdaterService creates a new updater service.
func NewUpdaterService(ctx context.Context, rg registry.Client, rc release.Client,
	port string, developmentMode bool, delaySeconds int,
) *UpdaterService {
	return &UpdaterService{
		ctx:             ctx,
		port:            port,
		rg:              rg,
		rl:              rc,
		developmentMode: developmentMode,
		delaySeconds:    delaySeconds,
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

	b, _ := json.Marshal(&release.ReleaseInfo{
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
		t := time.NewTicker(time.Minute)
		for {
			select {
			case <-updater.ctx.Done():
				log.WithError(updater.ctx.Err()).Info("updater context is done")
				updater.stopServer()
				return
			case <-t.C:
				err := updater.updateLatestReleaseWithDelay(time.Duration(updater.delaySeconds) * time.Second)
				updater.lastErr.Set(err)
				updater.lastChecked.Set()
				if err != nil {
					log.WithError(err).Error("error getting release")
				}
			}
		}
	}()

	log.Info("updater initialization complete")
	return nil
}

func (updater *UpdaterService) updateLatestRelease() error {
	return updater.updateLatestReleaseWithDelay(0)
}

func (updater *UpdaterService) updateLatestReleaseWithDelay(delay time.Duration) error {
	if updater.developmentMode {
		return updater.readLocalReleaseManifest()
	}

	log.Info("updating latest release")

	ref, err := updater.rg.GetScannerNodeVersion()
	if err != nil {
		return fmt.Errorf("failed to get the latest release manifest ref: %v", err)
	}
	if ref != updater.latestReference {
		rm, err := updater.rl.GetReleaseManifest(context.Background(), ref)
		if err != nil {
			log.WithError(err).Error("error getting release manifest")
			return fmt.Errorf("failed while downloading the release manifest: %v", err)
		}

		// so that all scanners don't update simultaneously, this waits a period of time
		if delay > 0 {
			log.WithFields(log.Fields{
				"release": ref, "delay": delay,
			}).Info("delaying update")
			time.Sleep(delay)
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
	var release release.ReleaseManifest
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
