package updater

import (
	"context"
	"testing"

	"github.com/forta-protocol/forta-node/config"
	ms "github.com/forta-protocol/forta-node/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdaterService_UpdateLatestRelease(t *testing.T) {
	c := gomock.NewController(t)
	us := ms.NewMockUpdaterStore(c)
	is := ms.NewMockIPFSClient(c)
	updater := NewUpdaterService(context.Background(), us, is, "8080")

	us.EXPECT().GetLatestReference().Return("reference", nil).Times(1)
	is.EXPECT().GetReleaseManifest("reference").Return(&config.ReleaseManifest{}, nil).Times(1)
	err := updater.updateLatestRelease()
	assert.NoError(t, err)
}

func TestUpdaterService_UpdateLatestReleaseCached(t *testing.T) {
	c := gomock.NewController(t)
	us := ms.NewMockUpdaterStore(c)
	is := ms.NewMockIPFSClient(c)
	updater := NewUpdaterService(context.Background(), us, is, "8080")

	// update twice
	us.EXPECT().GetLatestReference().Return("reference", nil).Times(2)

	// only call ipfs once (because value is the same)
	is.EXPECT().GetReleaseManifest("reference").Return(&config.ReleaseManifest{}, nil).Times(1)
	assert.NoError(t, updater.updateLatestRelease())
	assert.NoError(t, updater.updateLatestRelease())
}

func TestUpdaterService_UpdateLatestReleaseNotCached(t *testing.T) {
	c := gomock.NewController(t)
	us := ms.NewMockUpdaterStore(c)
	is := ms.NewMockIPFSClient(c)
	updater := NewUpdaterService(context.Background(), us, is, "8080")

	// update twice
	us.EXPECT().GetLatestReference().Return("reference1", nil).Times(1)
	us.EXPECT().GetLatestReference().Return("reference2", nil).Times(1)

	// only call ipfs once (because value is the same)
	is.EXPECT().GetReleaseManifest("reference1").Return(&config.ReleaseManifest{}, nil).Times(1)
	is.EXPECT().GetReleaseManifest("reference2").Return(&config.ReleaseManifest{}, nil).Times(1)

	assert.NoError(t, updater.updateLatestRelease())
	assert.NoError(t, updater.updateLatestRelease())
}
