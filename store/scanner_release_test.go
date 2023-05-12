package store

import (
	"context"
	"github.com/forta-network/forta-core-go/release"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockReleaseStore struct {
	mockRm  *release.ReleaseManifest
	mockErr error
}

func (mrs *mockReleaseStore) GetReleaseManifest(ctx context.Context, reference string) (*release.ReleaseManifest, error) {
	return mrs.mockRm, mrs.mockErr
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
			assert.Nil(t, res, tst.name)
			assert.Error(t, err, tst.expectedErr, tst.name)
			continue
		}
		assert.NoError(t, err, tst.name)
		assert.Equal(t, tst.expectedRef, res.Reference, tst.name)
	}
}
