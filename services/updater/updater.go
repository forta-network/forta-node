package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/release"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/store"

	"github.com/forta-network/forta-node/nodeutils"
	log "github.com/sirupsen/logrus"
)

var (
	defaultUpdateCheckIntervalSeconds = 60
	maxConsecutiveUpdateErrors        = 60
)

// UpdaterService receives the release updates.
type UpdaterService struct {
	ctx  context.Context
	port string

	mu             sync.RWMutex
	releaseClient  release.Client
	registryClient registry.Client
	srs            store.ScannerReleaseStore
	server         *http.Server

	latestReference string
	latestRelease   *release.ReleaseManifest

	updateDelay         time.Duration
	updateCheckInterval time.Duration

	errCounter *nodeutils.ErrorCounter

	lastChecked        health.TimeTracker
	lastErr            health.ErrorTracker
	latestVersion      health.MessageTracker
	latestIsPrerelease health.MessageTracker
}

// NewUpdaterService creates a new updater service.
func NewUpdaterService(ctx context.Context, svs store.ScannerReleaseStore,
	port string, updateDelaySeconds, updateCheckIntervalSeconds int,
) *UpdaterService {
	if updateCheckIntervalSeconds == 0 {
		updateCheckIntervalSeconds = defaultUpdateCheckIntervalSeconds
	}

	return &UpdaterService{
		ctx:                 ctx,
		port:                port,
		srs:                 svs,
		updateDelay:         time.Duration(updateDelaySeconds) * time.Second,
		updateCheckInterval: time.Duration(updateCheckIntervalSeconds) * time.Second,
		errCounter: nodeutils.NewErrorCounter(uint(maxConsecutiveUpdateErrors), func(err error) bool {
			return err != nil // all non-nil errors are critical errors
		}),
	}
}

func (updater *UpdaterService) handleGetVersion(w http.ResponseWriter, r *http.Request) {
	updater.mu.RLock()
	defer updater.mu.RUnlock()

	if updater.latestRelease == nil || updater.latestReference == "" {
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

	if err := updater.updateLatestReleaseWithDelay(0); err != nil {
		log.WithError(err).Error("error initializing release")
		return err
	}

	utils.GoListenAndServe(updater.server)

	go func() {
		t := time.NewTicker(updater.updateCheckInterval)
		for {
			select {
			case <-updater.ctx.Done():
				log.WithError(updater.ctx.Err()).Info("updater context is done")
				updater.stopServer()
				return
			case <-t.C:
				err := updater.updateLatestReleaseWithDelay(updater.updateDelay)
				if updater.errCounter.TooManyErrs(err) {
					log.WithError(err).Panic("too many update errors - exiting")
				}
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

func (updater *UpdaterService) updateLatestReleaseWithDelay(delay time.Duration) error {
	log.Info("updating latest release")

	// note: if reference is blank, this returns an error
	latest, err := updater.srs.GetRelease(updater.ctx)
	if err != nil {
		return err
	}

	updater.mu.RLock()
	latestRef := updater.latestReference
	updater.mu.RUnlock()

	// if the same as before, return the value (blank isn't possible here)
	if latest.Reference == latestRef {
		log.WithFields(log.Fields{
			"release": latest.Reference,
		}).Info("no change to release")
		return nil
	}

	// so that all scanners don't update simultaneously, this waits a period of time
	if delay > 0 {
		log.WithFields(log.Fields{
			"release": latest.Reference, "delay": delay,
		}).Info("delaying update")

		// if a newer release is found while waiting, this returns and tries again
		// (this resets the delay clock)
		if foundNew := updater.waitForDelayOrNewerRelease(latest.Reference, delay); foundNew {
			log.Info("detected newer release while delaying current update - aborting")
			return nil
		}

		log.Info("successfully waited before version update")
	}

	updater.latestVersion.Set(latest.ReleaseManifest.Release.Version)
	updater.latestIsPrerelease.Set(strconv.FormatBool(latest.IsPrerelease))

	updater.mu.Lock()
	defer updater.mu.Unlock()
	updater.latestRelease = &latest.ReleaseManifest
	updater.latestReference = latest.Reference
	log.WithFields(log.Fields{
		"release": latest.Reference,
	}).Info("updating to release")

	return nil
}

// returns true if a newer release is detected, otherwise waits for delay and returns false
func (updater *UpdaterService) waitForDelayOrNewerRelease(currentRef string, delay time.Duration) bool {
	detectedCh := make(chan struct{})

	ctx, cancel := context.WithCancel(updater.ctx)
	defer cancel()

	go updater.waitForNewerRelease(ctx, currentRef, detectedCh)

	select {
	case <-time.After(delay):
		return false
	case <-detectedCh:
		return true
	}
}

// notifies channel is a newer version is detected
func (updater *UpdaterService) waitForNewerRelease(ctx context.Context, currentRef string, detectedCh chan struct{}) {
	ticker := time.NewTicker(updater.updateCheckInterval)
	for {
		select {
		case <-ticker.C:
			if rel, err := updater.srs.GetRelease(updater.ctx); err != nil {
				log.WithError(err).Error("error getting release during delay (ignoring intermittent)")
				continue
			} else if rel.Reference != currentRef {
				detectedCh <- struct{}{}
				return
			}
		case <-ctx.Done():
			return
		}
	}
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
		updater.latestVersion.GetReport("latest.version"),
		updater.latestIsPrerelease.GetReport("latest.is-prerelease"),
	}
}
