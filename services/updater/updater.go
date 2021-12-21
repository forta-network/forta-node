package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/store"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var updateInterval = 1 * time.Hour

// UpdaterService receives the release updates.
type UpdaterService struct {
	ctx  context.Context
	port string

	mu     sync.RWMutex
	ipfs   store.IPFSClient
	us     store.UpdaterStore
	cancel context.CancelFunc
	server *http.Server

	latestReference string
	latestRelease   *config.ReleaseManifest
}

// NewUpdaterService creates a new updater service.
func NewUpdaterService(ctx context.Context, us store.UpdaterStore, ipfs store.IPFSClient, port string) *UpdaterService {
	ctx, cancel := context.WithCancel(ctx)
	return &UpdaterService{
		ctx:    ctx,
		port:   port,
		us:     us,
		ipfs:   ipfs,
		cancel: cancel,
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

	b, _ := json.Marshal(updater.latestRelease)
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

	t := time.NewTicker(updateInterval)
	grp.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			if err := updater.updateLatestRelease(); err != nil {
				log.WithError(err).Error("error getting release")
				return err
			}
		}
		return nil
	})

	//initialize at start
	if err := updater.updateLatestRelease(); err != nil {
		log.WithError(err).Error("error initializing release")
		return err
	}

	if err := grp.Wait(); err != nil {
		log.WithError(err).Error("error returned while running updater")
	}
	return nil
}

func (updater *UpdaterService) updateLatestRelease() error {
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

// Name returns the name of the service.
func (updater *UpdaterService) Name() string {
	return "updater"
}

// Stop stops the service
func (updater *UpdaterService) Stop() error {
	log.Infof("stopping %s", updater.Name())
	if updater.server != nil {
		if err := updater.server.Shutdown(updater.ctx); err != nil {
			log.WithError(err).Error("error stopping server")
			return err
		}
	}
	updater.cancel()
	return nil
}
