package updater

import (
	"context"
	"github.com/forta-network/forta-node/store"
	mock_store "github.com/forta-network/forta-node/store/mocks"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

const (
	testUpdateCheckIntervalSeconds = 1
	testUpdateDelaySeconds         = 15
)

func TestUpdaterService_UpdateLatestRelease(t *testing.T) {
	r := require.New(t)

	svs := mock_store.NewMockScannerReleaseStore(gomock.NewController(t))
	updater := NewUpdaterService(
		context.Background(), svs, "8080", testUpdateDelaySeconds, testUpdateCheckIntervalSeconds,
	)

	svs.EXPECT().GetRelease(gomock.Any()).Return(&store.ScannerRelease{
		Reference: "reference",
	}, nil).Times(1)

	err := updater.updateLatestReleaseWithDelay(0)
	r.NoError(err)
}

func TestUpdaterService_UpdateLatestReleaseNotCached(t *testing.T) {
	r := require.New(t)

	svs := mock_store.NewMockScannerReleaseStore(gomock.NewController(t))
	updater := NewUpdaterService(
		context.Background(), svs, "8080", testUpdateDelaySeconds, testUpdateCheckIntervalSeconds,
	)

	svs.EXPECT().GetRelease(gomock.Any()).Return(&store.ScannerRelease{
		Reference: "reference1",
	}, nil).Times(1)

	svs.EXPECT().GetRelease(gomock.Any()).Return(&store.ScannerRelease{
		Reference: "reference2",
	}, nil).Times(1)

	r.NoError(updater.updateLatestReleaseWithDelay(0))
	r.Equal("reference1", updater.latestReference)

	r.NoError(updater.updateLatestReleaseWithDelay(0))
	r.Equal("reference2", updater.latestReference)
}

func TestUpdaterService_UpdateLatestReleaseAbort(t *testing.T) {
	r := require.New(t)

	svs := mock_store.NewMockScannerReleaseStore(gomock.NewController(t))
	updater := NewUpdaterService(
		context.Background(), svs, "8080", testUpdateDelaySeconds, testUpdateCheckIntervalSeconds,
	)

	initalLatestRef := updater.latestReference

	svs.EXPECT().GetRelease(gomock.Any()).Return(&store.ScannerRelease{
		Reference: "reference1",
	}, nil).Times(1)

	svs.EXPECT().GetRelease(gomock.Any()).Return(&store.ScannerRelease{
		Reference: "reference2",
	}, nil).Times(1)

	r.NoError(updater.updateLatestReleaseWithDelay(updater.updateDelay))

	// update should be ineffective and be aborted
	r.Equal(initalLatestRef, updater.latestReference)
}
