package store

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/release"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/config"
)

var ErrBlankReference = errors.New("reference is blank")

type ScannerReleaseStore interface {
	GetRelease(ctx context.Context) (*ScannerRelease, error)
}

type ScannerRelease struct {
	Reference       string
	ReleaseManifest release.ReleaseManifest
	IsPrerelease    bool
	IsDevMode       bool
}

type lookupVersionStore struct {
	rc            release.Client
	lookup        func() (string, error)
	isPrerelease  bool
	cachedRelease *ScannerRelease
	mux           sync.Mutex
}

// GetRelease is a thread-safe lookup of the current version
// it caches the latest version and returns it if the reference matches
func (l *lookupVersionStore) GetRelease(ctx context.Context) (*ScannerRelease, error) {
	l.mux.Lock()
	defer l.mux.Unlock()

	ref, err := l.lookup()
	if err != nil {
		log.WithError(err).Error("error calling contract for version")
		return nil, err
	}
	if ref == "" {
		log.WithError(ErrBlankReference).Error("version ref is blank")
		return nil, ErrBlankReference
	}

	if l.cachedRelease != nil && l.cachedRelease.Reference == ref {
		return l.cachedRelease, nil
	}

	rm, err := loadRef(ctx, l.rc, ref)
	if err != nil {
		return nil, err
	}
	if rm == nil {
		return nil, errors.New("release manifest is nil")
	}
	res := &ScannerRelease{
		Reference:       ref,
		ReleaseManifest: *rm,
		IsPrerelease:    l.isPrerelease,
	}

	l.cachedRelease = res
	return res, nil
}

func loadRef(ctx context.Context, rc release.Client, ref string) (*release.ReleaseManifest, error) {
	res, err := rc.GetReleaseManifest(ctx, ref)
	if err != nil {
		return nil, err
	}
	if res.Release.Version == "" {
		return nil, errors.New("release was blank")
	}
	return res, nil
}

type devScannerVersionStore struct{}

func (d *devScannerVersionStore) GetRelease(ctx context.Context) (*ScannerRelease, error) {
	b, err := os.ReadFile(path.Join(config.DefaultContainerFortaDirPath, "local-release.json"))
	if err != nil {
		log.WithError(err).Info("could not read the test release manifest file - ignoring error")
		return nil, err
	}
	var rm release.ReleaseManifest
	if err := json.Unmarshal(b, &rm); err != nil {
		log.WithError(err).Info("could not unmarshal the test release manifest - ignoring error")
		return nil, err
	}
	return &ScannerRelease{
		Reference:       rm.Release.Commit,
		ReleaseManifest: rm,
		IsDevMode:       true,
	}, nil
}
func NewScannerReleaseStore(ctx context.Context, cfg config.Config) (ScannerReleaseStore, error) {
	developmentMode := utils.ParseBoolEnvVar(config.EnvDevelopment)
	if developmentMode {
		return &devScannerVersionStore{}, nil
	}

	releaseClient, err := release.NewClient(cfg.Registry.IPFS.GatewayURL)
	if err != nil {
		return nil, err
	}

	registryClient, err := GetRegistryClient(ctx, cfg, registry.ClientConfig{
		JsonRpcUrl: cfg.Registry.JsonRpc.Url,
		ENSAddress: cfg.ENSConfig.ContractAddress,
		Name:       "updater",
	})
	if err != nil {
		return nil, err
	}

	lookup := registryClient.GetScannerNodeVersion
	if cfg.AutoUpdate.TrackPrereleases {
		lookup = registryClient.GetScannerNodePrereleaseVersion
	}
	return &lookupVersionStore{
		rc:           releaseClient,
		lookup:       lookup,
		isPrerelease: cfg.AutoUpdate.TrackPrereleases,
		mux:          sync.Mutex{},
	}, nil
}
