package store

import (
	"context"
	"errors"
	"testing"

	"github.com/forta-network/forta-core-go/manifest"
	mock_manifest "github.com/forta-network/forta-core-go/manifest/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestBotManifestStore(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	manifestClient := mock_manifest.NewMockClient(ctrl)

	testManifest := &manifest.SignedAgentManifest{
		Manifest: &manifest.AgentManifest{
			AgentID:        &testBot1,
			ImageReference: &testImage1,
			ChainSettings: map[string]manifest.AgentChainSettings{
				"123": {
					Shards: 6,
					Target: 1,
				},
			},
		},
	}
	testManifestRef := "test-manifest-ref"

	manifestStore := NewBotManifestStore(manifestClient)
	manifestStore.maxRetries = 1 // override the default

	// the first call fails: hit the client and then return
	manifestClient.EXPECT().GetAgentManifest(gomock.Any(), testManifestRef).Return(nil, errors.New("failed"))
	manifest, err := manifestStore.GetBotManifest(context.Background(), testManifestRef)
	r.Error(err)
	r.Nil(manifest)

	// the second call succeeds: hit the client, set the cache and return
	manifestClient.EXPECT().GetAgentManifest(gomock.Any(), testManifestRef).Return(testManifest, nil)
	manifest, err = manifestStore.GetBotManifest(context.Background(), testManifestRef)
	r.NoError(err)
	r.Equal(testManifest, manifest)

	// third call succeeds: hit the cache and return
	manifest, err = manifestStore.GetBotManifest(context.Background(), testManifestRef)
	r.NoError(err)
	r.Equal(testManifest, manifest)
}
