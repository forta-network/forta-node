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
	ctx            context.Context
	port           string
	scannerAddress string

	mu             sync.RWMutex
	releaseClient  release.Client
	registryClient registry.Client
	srs            store.ScannerReleaseStore
	server         *http.Server

	latestReference string
	latestRelease   *release.ReleaseManifest

	overrideUpdateDelaySeconds *int
	updateCheckInterval        time.Duration

	errCounter *nodeutils.ErrorCounter

	lastChecked        health.TimeTracker
	lastErr            health.ErrorTracker
	finalCheck         health.TimeTracker
	finalCheckErr      health.ErrorTracker
	latestVersion      health.MessageTracker
	latestIsPrerelease health.MessageTracker
}

// NewUpdaterService creates a new updater service.
func NewUpdaterService(ctx context.Context, svs store.ScannerReleaseStore,
	port, scannerAddress string, overrideUpdateDelaySeconds *int, updateCheckIntervalSeconds int,
) *UpdaterService {
	if updateCheckIntervalSeconds == 0 {
		updateCheckIntervalSeconds = defaultUpdateCheckIntervalSeconds
	}

	return &UpdaterService{
		ctx:                        ctx,
		port:                       port,
		scannerAddress:             scannerAddress,
		srs:                        svs,
		overrideUpdateDelaySeconds: overrideUpdateDelaySeconds,
		updateCheckInterval:        time.Duration(updateCheckIntervalSeconds) * time.Second,
		errCounter: nodeutils.NewErrorCounter(uint(maxConsecutiveUpdateErrors), func(err error) bool {
			return err != nil // all non-nil errors are critical errors
		}),
	}
}

func (updater *UpdaterService) calculateDelay() time.Duration {
	if updater.overrideUpdateDelaySeconds != nil {
		return time.Duration(*updater.overrideUpdateDelaySeconds) * time.Second
	}
	// if anything goes wrong, just stick to the default schedule
	maxUpdateDelay := time.Hour * release.DefaultAutoUpdateHours
	// take the max auto-update delay from the release manifest if it's non-zero
	if updater.latestRelease != nil &&
		updater.latestRelease.Release.Config.AutoUpdateInHours > 0 {
		maxUpdateDelay = time.Duration(updater.latestRelease.Release.Config.AutoUpdateInHours) * time.Hour
	}
	return CalculateReleaseDelay(
		updater.scannerAddress,
		maxUpdateDelay,
	)
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
				err := updater.updateLatestReleaseWithDelay(updater.calculateDelay())
				if updater.errCounter.TooManyErrs(err) {
					log.WithError(err).Panic("too many update errors - exiting")
				}
				updater.lastErr.Set(err)
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
	updater.lastChecked.Set()
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
	log.WithFields(log.Fields{"release": latest.Reference, "delay": delay}).Info("delaying update")
	if delay > 0 {
		<-time.After(delay)
	}
	log.Info("successfully waited before version update")

	for {
		updater.finalCheck.Set()
		latest, err = updater.srs.GetRelease(updater.ctx)
		if err == nil {
			break
		}
		updater.finalCheckErr.Set(err)
		log.WithError(err).Error("failed to get the latest release just after the delay is over - retrying")
		time.Sleep(time.Second * 5)
	}
	updater.finalCheckErr.Set(nil)
	log.WithFields(log.Fields{"release": latest.Reference, "delay": delay}).
		Info("successfully got the latest release once more after the delay")

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
		updater.finalCheck.GetReport("event.checked.final.time"),
		updater.finalCheckErr.GetReport("event.checked.final.error"),
		updater.latestVersion.GetReport("latest.version"),
		updater.latestIsPrerelease.GetReport("latest.is-prerelease"),
	}
}
