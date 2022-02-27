package updater

import (
	"context"
	"github.com/forta-protocol/forta-core-go/release"
	"testing"

	im "github.com/forta-protocol/forta-core-go/release/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	ms "github.com/forta-protocol/forta-node/store/mocks"
)

func TestUpdaterService_UpdateLatestRelease(t *testing.T) {
	c := gomock.NewController(t)

	us := ms.NewMockUpdaterStore(c)
	is := im.NewMockClient(c)
	updater := NewUpdaterService(context.Background(), us, is, "8080", false)

	us.EXPECT().GetLatestReference().Return("reference", nil).Times(1)
	is.EXPECT().GetReleaseManifest(gomock.Any(), "reference").Return(&release.ReleaseManifest{}, nil).Times(1)
	err := updater.updateLatestRelease()
	assert.NoError(t, err)
}

func TestUpdaterService_UpdateLatestReleaseCached(t *testing.T) {
	c := gomock.NewController(t)
	us := ms.NewMockUpdaterStore(c)
	is := im.NewMockClient(c)
	updater := NewUpdaterService(context.Background(), us, is, "8080", false)

	// update twice
	us.EXPECT().GetLatestReference().Return("reference", nil).Times(2)

	// only call ipfs once (because value is the same)
	is.EXPECT().GetReleaseManifest(gomock.Any(), "reference").Return(&release.ReleaseManifest{}, nil).Times(1)
	assert.NoError(t, updater.updateLatestRelease())
	assert.NoError(t, updater.updateLatestRelease())
}

func TestUpdaterService_UpdateLatestReleaseNotCached(t *testing.T) {
	c := gomock.NewController(t)
	us := ms.NewMockUpdaterStore(c)
	is := im.NewMockClient(c)
	updater := NewUpdaterService(context.Background(), us, is, "8080", false)

	// update twice
	us.EXPECT().GetLatestReference().Return("reference1", nil).Times(1)
	us.EXPECT().GetLatestReference().Return("reference2", nil).Times(1)

	// only call ipfs once (because value is the same)
	is.EXPECT().GetReleaseManifest(gomock.Any(), "reference1").Return(&release.ReleaseManifest{}, nil).Times(1)
	is.EXPECT().GetReleaseManifest(gomock.Any(), "reference2").Return(&release.ReleaseManifest{}, nil).Times(1)

	assert.NoError(t, updater.updateLatestRelease())
	assert.NoError(t, updater.updateLatestRelease())
}
