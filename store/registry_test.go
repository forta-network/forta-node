package store

import (
	"errors"
	"testing"

	"github.com/forta-network/forta-core-go/manifest"
	mock_manifest "github.com/forta-network/forta-core-go/manifest/mocks"
	"github.com/forta-network/forta-core-go/registry"
	mock_registry "github.com/forta-network/forta-core-go/registry/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_calculateShardID(t *testing.T) {
	type args struct {
		target uint
		shards uint
		idx    uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		// [0,0,0,1,1,1]
		{
			name: "2 shards, 3 target",
			args: args{target: 3, shards: 2, idx: 3},
			want: 1,
		},
		// [0,0,1,1,2,2]
		{
			name: "3 shards, 2 target",
			args: args{target: 2, shards: 3, idx: 3},
			want: 1,
		},
		// [0,1,2,3,4,5]
		{
			name: "6 shards, 1 target",
			args: args{target: 1, shards: 6, idx: 5},
			want: 5,
		},
		// 	[0,0,0,0,0,0]
		{
			name: "no redundancy",
			args: args{target: 6, shards: 1, idx: 4},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := calculateShardID(tt.args.target, tt.args.idx); got != tt.want {
					t.Errorf("calculateShardID() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
func TestSetShardInformation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocked registry and manifest clients
	mockRegistryClient := mock_registry.NewMockClient(ctrl)
	mockManifestClient := mock_manifest.NewMockClient(ctrl)

	rs := &registryStore{
		rc: mockRegistryClient,
		mc: mockManifestClient,
	}

	scanner := "scanner-1"

	bots := []config.AgentConfig{
		{
			ID:       "bot-1",
			Manifest: "manifest-1",
		},
		{
			ID:       "bot-2",
			Manifest: "manifest-2",
		},
	}

	assignments := []*registry.Assignment{
		{
			AgentID:          "bot-1",
			AgentManifest:    "manifest-1",
			AssignedScanners: 3,
			ScannerIndex:     1,
		},
		{
			AgentID:          "bot-2",
			AgentManifest:    "manifest-2",
			AssignedScanners: 5,
			ScannerIndex:     2,
		},
	}

	manifests := map[string]*manifest.SignedAgentManifest{
		"manifest-1": {
			Manifest: &manifest.AgentManifest{
				ChainSettings: map[string]manifest.AgentChainSettings{
					keyDefaultChainSetting: {Target: 10, Shards: 2},
				},
			},
		},
		"manifest-2": {
			Manifest: &manifest.AgentManifest{
				ChainSettings: map[string]manifest.AgentChainSettings{
					keyDefaultChainSetting: {Target: 8, Shards: 1},
				},
			},
		},
	}

	t.Run(
		"Successful populateShardedBotConfigs", func(t *testing.T) {
			mockRegistryClient.EXPECT().GetAssignmentList(nil, gomock.Any(), scanner).Return(assignments, nil)

			mockManifestClient.EXPECT().GetAgentManifest(gomock.Any(), "manifest-1").Return(manifests["manifest-1"], nil)
			mockManifestClient.EXPECT().GetAgentManifest(gomock.Any(), "manifest-2").Return(manifests["manifest-2"], nil)

			expectedBots := []config.AgentConfig{
				{
					ID:       "bot-1",
					Manifest: "manifest-1",
					ShardConfig: &config.ShardConfig{
						ShardID: 0,
						Shards:  2,
						Target:  10,
					},
				},
				{
					ID:       "bot-2",
					Manifest: "manifest-2",
					ShardConfig: &config.ShardConfig{
						ShardID: 0,
						Shards:  1,
						Target:  8,
					},
				},
			}

			result, err := rs.populateShardedBotConfigs(scanner, bots)
			assert.NoError(t, err, "populateShardedBotConfigs should not return an error")
			assert.Len(t, result, len(expectedBots), "populateShardedBotConfigs returned a different number of bots")

			for i, bot := range result {
				expectedBot := expectedBots[i]
				assert.Equal(t, expectedBot, bot, "populateShardedBotConfigs returned incorrect shard configuration for bot %s", bot.ID)
			}
		},
	)

	t.Run(
		"Failed GetAssignmentList", func(t *testing.T) {
			expectedError := errors.New("error getting assignments")

			mockRegistryClient.EXPECT().GetAssignmentList(nil, gomock.Any(), scanner).Return(nil, expectedError)

			result, err := rs.populateShardedBotConfigs(scanner, bots)
			assert.Error(t, err, "populateShardedBotConfigs should return an error")
			assert.Nil(t, result, "populateShardedBotConfigs should return nil result")
			assert.EqualError(t, err, expectedError.Error(), "populateShardedBotConfigs returned incorrect error")
		},
	)

	t.Run(
		"Failed GetAgentManifest", func(t *testing.T) {
			expectedError := errors.New("error getting agent manifest")

			mockRegistryClient.EXPECT().GetAssignmentList(nil, gomock.Any(), scanner).Return(assignments, nil)

			mockManifestClient.EXPECT().GetAgentManifest(gomock.Any(), "manifest-1").Return(nil, expectedError)
			mockManifestClient.EXPECT().GetAgentManifest(gomock.Any(), "manifest-2").Return(nil, expectedError)

			result, err := rs.populateShardedBotConfigs(scanner, bots)
			assert.NoError(t, err, "populateShardedBotConfigs should continue when a single bot fails")
			for i, bot := range result {
				expectedBot := bots[i]
				assert.Equal(t, expectedBot, bot, "populateShardedBotConfigs returned incorrect shard configuration for bot %s", bot.ID)
			}
		},
	)

	t.Run(
		"No matching assignments", func(t *testing.T) {
			mockRegistryClient.EXPECT().GetAssignmentList(nil, gomock.Any(), scanner).Return(nil, nil)

			result, err := rs.populateShardedBotConfigs(scanner, bots)
			assert.NoError(t, err, "populateShardedBotConfigs should not return an error")
			assert.Len(t, result, len(bots), "populateShardedBotConfigs should return the same number of bots")
			for _, bot := range result {
				assert.Nil(t, bot.ShardConfig, "populateShardedBotConfigs should not set ShardConfig for bots with no matching assignment")
			}
		},
	)
	t.Run(
		"Failed GetAgentManifest for some bots", func(t *testing.T) {
			expectedError := errors.New("error getting agent manifest")

			mockRegistryClient.EXPECT().GetAssignmentList(nil, gomock.Any(), scanner).Return(assignments, nil)

			mockManifestClient.EXPECT().GetAgentManifest(gomock.Any(), "manifest-1").Return(nil, expectedError)
			mockManifestClient.EXPECT().GetAgentManifest(gomock.Any(), "manifest-2").Return(manifests["manifest-2"], nil)

			result, err := rs.populateShardedBotConfigs(scanner, bots)
			assert.NoError(t, err, "populateShardedBotConfigs should return not an error")
			for _, agentConfig := range result {
				switch agentConfig.ID {
				case "bot-1":
					assert.Nil(t, agentConfig.ShardConfig, "shard config should be nil for bad bots")
				case "bot-2":
					assert.NotNil(t, agentConfig.ShardConfig, "shard config should not be nil for bad bots")
				}
			}
		},
	)
}
