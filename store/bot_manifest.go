package store

import (
	"context"
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/manifest"
	"github.com/patrickmn/go-cache"
)

const (
	botManifestExpiry = time.Hour * 6
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
		manifestCache:  cache.New(botManifestExpiry, time.Hour),
		manifestClient: manifestClient,
		maxRetries:     10,
	}
}

func (bms *botManifestStore) GetBotManifest(ctx context.Context, ref string) (*manifest.SignedAgentManifest, error) {
	cachedManifest, ok := bms.manifestCache.Get(ref)
	if ok {
		bms.manifestCache.Set(ref, cachedManifest, 0)
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
