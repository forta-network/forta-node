package store

import (
	"context"
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/manifest"
	"github.com/patrickmn/go-cache"
)

// BotManifestStore loads bot manifests.
type BotManifestStore interface {
	GetBotManifest(ctx context.Context, ref string) (*manifest.SignedAgentManifest, error)
}

type botManifestStore struct {
	manifestCache  *cache.Cache
	manifestClient manifest.Client
	maxRetries     int
}

var _ BotManifestStore = &botManifestStore{}

// NewBotManifestStore creates a new bot manifest store.
func NewBotManifestStore(manifestClient manifest.Client) *botManifestStore {
	return &botManifestStore{
		manifestCache:  cache.New(time.Hour*6, time.Hour),
		manifestClient: manifestClient,
		maxRetries:     10,
	}
}

func (bms *botManifestStore) GetBotManifest(ctx context.Context, ref string) (*manifest.SignedAgentManifest, error) {
	cachedManifest, ok := bms.manifestCache.Get(ref)
	if ok {
		return cachedManifest.(*manifest.SignedAgentManifest), nil
	}

	var (
		loadedManifest *manifest.SignedAgentManifest
		err            error
	)
	for i := 0; i < bms.maxRetries; i++ {
		loadedManifest, err = bms.manifestClient.GetAgentManifest(ctx, ref)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load the bot manifest: %v", err)
	}

	bms.manifestCache.Set(ref, loadedManifest, 0)
	return loadedManifest, err
}
