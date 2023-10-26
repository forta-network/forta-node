package updater

import (
	"context"
	"testing"

	"github.com/forta-network/forta-node/store"
	mock_store "github.com/forta-network/forta-node/store/mocks"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

var (
	testUpdateCheckIntervalSeconds = 1
	testUpdateDelaySeconds         = 15
)

func TestUpdaterService_UpdateLatestRelease(t *testing.T) {
	r := require.New(t)

	svs := mock_store.NewMockScannerReleaseStore(gomock.NewController(t))
	updater := NewUpdaterService(
		context.Background(), svs, "8080", "", &testUpdateDelaySeconds, testUpdateCheckIntervalSeconds,
	)

	svs.EXPECT().GetRelease(gomock.Any()).Return(&store.ScannerRelease{
		Reference: "reference",
	}, nil).Times(2)

	err := updater.updateLatestRelease(true)
	r.NoError(err)
}

func TestUpdaterService_UpdateLatestRelease_SingleEachTime(t *testing.T) {
	r := require.New(t)

	svs := mock_store.NewMockScannerReleaseStore(gomock.NewController(t))
	updater := NewUpdaterService(
		context.Background(), svs, "8080", "", &testUpdateDelaySeconds, testUpdateCheckIntervalSeconds,
	)

	svs.EXPECT().GetRelease(gomock.Any()).Return(&store.ScannerRelease{
		Reference: "reference1",
	}, nil).Times(2)

	svs.EXPECT().GetRelease(gomock.Any()).Return(&store.ScannerRelease{
		Reference: "reference2",
	}, nil).Times(2)

	r.NoError(updater.updateLatestRelease(true))
	r.Equal("reference1", updater.latestReference)

	r.NoError(updater.updateLatestRelease(true))
	r.Equal("reference2", updater.latestReference)
}

func TestUpdaterService_UpdateLatestRelease_TwoInARow(t *testing.T) {
	r := require.New(t)

	svs := mock_store.NewMockScannerReleaseStore(gomock.NewController(t))
	updater := NewUpdaterService(
		context.Background(), svs, "8080", "", &testUpdateDelaySeconds, testUpdateCheckIntervalSeconds,
	)

	finalRef := "reference2"

	svs.EXPECT().GetRelease(gomock.Any()).Return(&store.ScannerRelease{
		Reference: "reference1",
	}, nil).Times(1)

	svs.EXPECT().GetRelease(gomock.Any()).Return(&store.ScannerRelease{
		Reference: "reference2",
	}, nil).Times(1)

	r.NoError(updater.updateLatestRelease(false))

	// should update to the latest one
	r.Equal(finalRef, updater.latestReference)
}
