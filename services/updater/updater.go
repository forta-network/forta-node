package updater

import (
	"context"
	"encoding/json"
	"errors"
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

var (
	errNotAvailable = errors.New("new release not available")

	defaultUpdateCheckIntervalSeconds = 60
)

// UpdaterService receives the release updates.
type UpdaterService struct {
	ctx  context.Context
	port string

	mu             sync.RWMutex
	releaseClient  release.Client
	registryClient registry.Client
	server         *http.Server

	developmentMode bool

	latestReference string
	latestRelease   *release.ReleaseManifest

	updateDelay         time.Duration
	updateCheckInterval time.Duration

	lastChecked health.TimeTracker
	lastErr     health.ErrorTracker
}

// NewUpdaterService creates a new updater service.
func NewUpdaterService(ctx context.Context, registryClient registry.Client, releaseClient release.Client,
	port string, developmentMode bool, updateDelaySeconds, updateCheckIntervalSeconds int,
) *UpdaterService {
	if updateCheckIntervalSeconds == 0 {
		updateCheckIntervalSeconds = defaultUpdateCheckIntervalSeconds
	}

	return &UpdaterService{
		ctx:                 ctx,
		port:                port,
		releaseClient:       releaseClient,
		registryClient:      registryClient,
		developmentMode:     developmentMode,
		updateDelay:         time.Duration(updateDelaySeconds) * time.Second,
		updateCheckInterval: time.Duration(updateCheckIntervalSeconds) * time.Second,
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
		t := time.NewTicker(updater.updateCheckInterval)
		for {
			select {
			case <-updater.ctx.Done():
				log.WithError(updater.ctx.Err()).Info("updater context is done")
				updater.stopServer()
				return
			case <-t.C:
				err := updater.updateLatestReleaseWithDelay(updater.updateCheckInterval)
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

	ref, rm, err := updater.getNewerRelease(updater.latestReference)
	switch err {
	case nil:
		// we downloaded new release info successfully

	case errNotAvailable:
		log.WithFields(log.Fields{
			"release": ref,
		}).Info("no change to release")
		return nil

	default:
		return err
	}

	// so that all scanners don't update simultaneously, this waits a period of time
	if delay > 0 {
		log.WithFields(log.Fields{
			"release": ref, "delay": delay,
		}).Info("delaying update")

		if foundNew := updater.checkNewerReleaseAndWait(ref, delay); foundNew {
			log.Info("detected newer release while delaying current update - aborting")
			return nil
		}

		log.Info("successfully waited before version update")
	}

	updater.mu.Lock()
	defer updater.mu.Unlock()
	updater.latestRelease = rm
	updater.latestReference = ref
	log.WithFields(log.Fields{
		"release": ref,
	}).Info("updating to release")

	return nil
}

func (updater *UpdaterService) getNewerRelease(previousRef string) (string, *release.ReleaseManifest, error) {
	ref, err := updater.compareScannerNodeVersion(previousRef)
	if err != nil {
		return ref, nil, err
	}
	rm, err := updater.releaseClient.GetReleaseManifest(context.Background(), ref)
	if err != nil {
		log.WithError(err).Error("error getting release manifest")
		return ref, nil, fmt.Errorf("failed while downloading the release manifest: %v", err)
	}
	return ref, rm, nil
}

func (updater *UpdaterService) compareScannerNodeVersion(previousRef string) (newRef string, err error) {
	ref, err := updater.registryClient.GetScannerNodeVersion()
	if err != nil {
		log.WithError(err).Error("error getting the latest release manifest ref")
		return "", fmt.Errorf("failed to get the latest release manifest ref: %v", err)
	}
	if ref == previousRef {
		return ref, errNotAvailable
	}
	return ref, nil
}

func (updater *UpdaterService) checkNewerReleaseAndWait(previousRef string, delay time.Duration) (foundNew bool) {
	detectedCh := make(chan struct{})
	defer close(detectedCh)

	ctx, cancel := context.WithCancel(updater.ctx)
	defer cancel()

	go updater.detectNewerRelease(ctx, previousRef, detectedCh)

	select {
	case <-time.After(delay):
		return false
	case <-detectedCh:
		return true
	}
}

func (updater *UpdaterService) detectNewerRelease(ctx context.Context, previousRef string, detectedCh chan struct{}) {
	ticker := time.NewTicker(updater.updateCheckInterval)
	for {
		select {
		case <-ticker.C:
			newRef, _ := updater.compareScannerNodeVersion(previousRef)
			if newRef != previousRef {
				detectedCh <- struct{}{}
				return
			}
		case <-ctx.Done():
			return
		}
	}
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
