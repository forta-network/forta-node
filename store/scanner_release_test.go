package store

import (
	"context"
	"testing"

	"github.com/forta-network/forta-core-go/release"
	"github.com/forta-network/forta-node/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockReleaseStore struct {
	mockRm  *release.ReleaseManifest
	mockErr error
}

func (mrs *mockReleaseStore) GetReleaseManifest(reference string) (*release.ReleaseManifest, error) {
	return mrs.mockRm, mrs.mockErr
}

// TestNewScannerReleaseStore returns a release
func TestNewScannerReleaseStore(t *testing.T) {
	rs, err := NewScannerReleaseStore(context.Background(), config.Config{
		Registry: config.RegistryConfig{
			JsonRpc:                config.JsonRpcConfig{Url: "https://polygon-rpc.com"},
			ReleaseDistributionUrl: "https://dist.forta.network/manifests/releases",
			IPFS: config.IPFSConfig{
				GatewayURL: "https://ipfs.forta.network",
			},
		},
		ENSConfig: config.ENSConfig{ContractAddress: "0x08f42fcc52a9C2F391bF507C4E8688D0b53e1bd7"},
	})
	assert.NoError(t, err)

	rls, err := rs.GetRelease(context.Background())
	assert.NoError(t, err)

	assert.NotNil(t, rls)
	assert.True(t, len(rls.Reference) > 0)
}

func TestLookupVersionStore_GetRelease(t *testing.T) {
	mockRm := &release.ReleaseManifest{
		Release: release.Release{
			Version: "version",
			Commit:  "commit",
		},
	}
	type test struct {
		name string

		mockRef       string
		mockLookupErr error
		mockRM        *release.ReleaseManifest
		mockCached    *ScannerRelease

		expectedRef string
		expectedErr error
	}
	tests := []test{
		{
			name:        "update-no-cache",
			mockRef:     "test",
			mockRM:      mockRm,
			expectedRef: "test",
		},
		{
			name:    "update-with-cache",
			mockRef: "test",
			mockRM:  mockRm,
			mockCached: &ScannerRelease{
				Reference: "stale",
			},
			expectedRef: "test",
		},
		{
			name:    "not-updated-cached",
			mockRef: "test",
			mockRM:  nil,
			mockCached: &ScannerRelease{
				Reference: "test",
			},
			expectedRef: "test",
		},
		{
			name:    "update-with-blank-reference",
			mockRef: "",
			mockCached: &ScannerRelease{
				Reference: "stale",
			},
			expectedRef: "",
			expectedErr: ErrBlankReference,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			r := require.New(t)

			lookup := func() (string, error) {
				return tst.mockRef, tst.mockLookupErr
			}
			lvs := &lookupVersionStore{
				rc:            &mockReleaseStore{mockRm: tst.mockRM},
				lookup:        lookup,
				cachedRelease: tst.mockCached,
			}
			res, err := lvs.GetRelease(context.Background())
			if tst.expectedErr != nil {
				r.Nil(res)
				r.Error(err, tst.expectedErr)
				return
			}
			r.NoError(err)
			r.Equal(tst.expectedRef, res.Reference)
		})
	}
}
