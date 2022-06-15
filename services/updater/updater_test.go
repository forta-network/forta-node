package updater

import (
	"context"
	"testing"

	"github.com/forta-network/forta-core-go/release"
	"github.com/stretchr/testify/require"

	rm "github.com/forta-network/forta-core-go/registry/mocks"
	im "github.com/forta-network/forta-core-go/release/mocks"
	"github.com/golang/mock/gomock"
)

const (
	testUpdateCheckIntervalSeconds = 1
	testUpdateDelaySeconds         = 15
)

func TestUpdaterService_UpdateLatestRelease(t *testing.T) {
	r := require.New(t)

	registryClient := rm.NewMockClient(gomock.NewController(t))
	releaseClient := im.NewMockClient(gomock.NewController(t))
	updater := NewUpdaterService(
		context.Background(), registryClient, releaseClient, "8080", false, false,
		testUpdateDelaySeconds, testUpdateCheckIntervalSeconds,
	)

	registryClient.EXPECT().GetScannerNodeVersion().Return("reference", nil).Times(1)
	releaseClient.EXPECT().GetReleaseManifest(gomock.Any(), "reference").Return(&release.ReleaseManifest{}, nil).
		Times(1)
	err := updater.updateLatestRelease()
	r.NoError(err)
}

func TestUpdaterService_UpdateLatestReleaseCached(t *testing.T) {
	r := require.New(t)

	registryClient := rm.NewMockClient(gomock.NewController(t))
	releaseClient := im.NewMockClient(gomock.NewController(t))
	updater := NewUpdaterService(
		context.Background(), registryClient, releaseClient, "8080", false, false,
		testUpdateDelaySeconds, testUpdateCheckIntervalSeconds,
	)

	// request latest version ref twice
	registryClient.EXPECT().GetScannerNodeVersion().Return("reference", nil).Times(2)
	// only get release manifest once (because release ref value is the same)
	releaseClient.EXPECT().GetReleaseManifest(gomock.Any(), "reference").Return(&release.ReleaseManifest{}, nil).
		Times(1)

	r.NoError(updater.updateLatestRelease())
	r.NoError(updater.updateLatestRelease())
}

func TestUpdaterService_UpdateLatestReleaseNotCached(t *testing.T) {
	r := require.New(t)

	registryClient := rm.NewMockClient(gomock.NewController(t))
	releaseClient := im.NewMockClient(gomock.NewController(t))
	updater := NewUpdaterService(
		context.Background(), registryClient, releaseClient, "8080", false, false,
		testUpdateDelaySeconds, testUpdateCheckIntervalSeconds,
	)

	// update twice

	registryClient.EXPECT().GetScannerNodeVersion().Return("reference1", nil).Times(1)
	releaseClient.EXPECT().GetReleaseManifest(gomock.Any(), "reference1").Return(&release.ReleaseManifest{}, nil).Times(1)
	registryClient.EXPECT().GetScannerNodeVersion().Return("reference2", nil).Times(1)
	releaseClient.EXPECT().GetReleaseManifest(gomock.Any(), "reference2").Return(&release.ReleaseManifest{}, nil).Times(1)

	r.NoError(updater.updateLatestRelease())
	r.NoError(updater.updateLatestRelease())
}

func TestUpdaterService_UpdateLatestReleaseAbort(t *testing.T) {
	r := require.New(t)

	registryClient := rm.NewMockClient(gomock.NewController(t))
	releaseClient := im.NewMockClient(gomock.NewController(t))
	updater := NewUpdaterService(
		context.Background(), registryClient, releaseClient, "8080", false, false,
		testUpdateDelaySeconds, testUpdateCheckIntervalSeconds,
	)

	initalLatestRef := updater.latestReference

	// get new version
	registryClient.EXPECT().GetScannerNodeVersion().Return("reference1", nil).Times(1)
	releaseClient.EXPECT().GetReleaseManifest(gomock.Any(), "reference1").Return(&release.ReleaseManifest{}, nil).
		Times(1)
	// receive new ref during the async wait
	registryClient.EXPECT().GetScannerNodeVersion().Return("reference2", nil).Times(1)
	r.NoError(updater.updateLatestReleaseWithDelay(updater.updateDelay))

	// update should be ineffective and be aborted
	r.Equal(initalLatestRef, updater.latestReference)
}
