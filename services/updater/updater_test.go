package updater

import (
	"context"
	"github.com/forta-protocol/forta-core-go/release"
	"testing"

	rm "github.com/forta-protocol/forta-core-go/registry/mocks"
	im "github.com/forta-protocol/forta-core-go/release/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdaterService_UpdateLatestRelease(t *testing.T) {
	c := gomock.NewController(t)

	rg := rm.NewMockClient(c)
	is := im.NewMockClient(c)
	updater := NewUpdaterService(context.Background(), rg, is, "8080", false)

	rg.EXPECT().GetScannerNodeVersion().Return("reference", nil).Times(1)
	is.EXPECT().GetReleaseManifest(gomock.Any(), "reference").Return(&release.ReleaseManifest{}, nil).Times(1)
	err := updater.updateLatestRelease()
	assert.NoError(t, err)
}

func TestUpdaterService_UpdateLatestReleaseCached(t *testing.T) {
	c := gomock.NewController(t)
	rg := rm.NewMockClient(c)
	is := im.NewMockClient(c)
	updater := NewUpdaterService(context.Background(), rg, is, "8080", false)

	// update twice
	rg.EXPECT().GetScannerNodeVersion().Return("reference", nil).Times(2)

	// only call ipfs once (because value is the same)
	is.EXPECT().GetReleaseManifest(gomock.Any(), "reference").Return(&release.ReleaseManifest{}, nil).Times(1)
	assert.NoError(t, updater.updateLatestRelease())
	assert.NoError(t, updater.updateLatestRelease())
}

func TestUpdaterService_UpdateLatestReleaseNotCached(t *testing.T) {
	c := gomock.NewController(t)
	rg := rm.NewMockClient(c)
	is := im.NewMockClient(c)
	updater := NewUpdaterService(context.Background(), rg, is, "8080", false)

	// update twice
	rg.EXPECT().GetScannerNodeVersion().Return("reference1", nil).Times(1)
	rg.EXPECT().GetScannerNodeVersion().Return("reference2", nil).Times(1)

	// only call ipfs once (because value is the same)
	is.EXPECT().GetReleaseManifest(gomock.Any(), "reference1").Return(&release.ReleaseManifest{}, nil).Times(1)
	is.EXPECT().GetReleaseManifest(gomock.Any(), "reference2").Return(&release.ReleaseManifest{}, nil).Times(1)

	assert.NoError(t, updater.updateLatestRelease())
	assert.NoError(t, updater.updateLatestRelease())
}
