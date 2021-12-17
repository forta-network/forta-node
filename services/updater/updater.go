package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/forta-protocol/forta-node/config"
)

// UpdaterService receives the release updates.
type UpdaterService struct {
	ctx  context.Context
	port string

	latestRelease *config.ReleaseManifest
	mu            sync.RWMutex
}

// NewUpdaterService creates a new updater service.
func NewUpdaterService(ctx context.Context, port string) *UpdaterService {
	return &UpdaterService{
		ctx:  ctx,
		port: port,
	}
}

// Start starts the service.
func (updater *UpdaterService) Start() error {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		updater.mu.RLock()
		defer updater.mu.RUnlock()

		if updater.latestRelease == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		b, _ := json.Marshal(updater.latestRelease)
		w.Write(b)
	}))

	go http.ListenAndServe(fmt.Sprintf(":%s", updater.port), nil)
	go updater.findLatestRelease()

	return nil
}

func (updater *UpdaterService) findLatestRelease() {
	// TODO: Find the latest release, unmarshal, lock, update in-memory, unlock.
}

// Name returns the name of the service.
func (updater *UpdaterService) Name() string {
	return "updater"
}

// Stop stops the service
func (updater *UpdaterService) Stop() error {
	return nil
}
