package store

import (
	"fmt"
	"testing"
	"time"

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

func TestPopulateShardConfig(t *testing.T) {
	type args struct {
		assignedScanners int
		scannerIndex     int
		chainSettings    map[string]manifest.AgentChainSettings
		chainID          int
	}
	tests := []struct {
		name                string
		args                args
		expectedShardConfig *config.ShardConfig
	}{
		{
			name: "can calculate shard based on chain setting",
			args: args{
				assignedScanners: 6,
				scannerIndex:     4,
				chainSettings: map[string]manifest.AgentChainSettings{
					"1": {
						Target: 2,
						Shards: 3,
					},
				},
				chainID: 1,
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 2,
				Target:  2,
				Shards:  3,
			},
		}, {
			name: "can calculate shard based on default setting",
			args: args{
				assignedScanners: 6,
				scannerIndex:     4,
				chainSettings: map[string]manifest.AgentChainSettings{
					"default": {
						Target: 2,
						Shards: 3,
					},
				},
				chainID: 1,
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 2,
				Target:  2,
				Shards:  3,
			},
		}, {
			name: "chain setting should override default setting",
			args: args{
				assignedScanners: 6,
				scannerIndex:     4,
				chainSettings: map[string]manifest.AgentChainSettings{
					"default": {
						Target: 2,
						Shards: 3,
					},
					"1": {
						Target: 6,
						Shards: 1,
					},
				},
				chainID: 1,
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 0,
				Target:  6,
				Shards:  1,
			},
		}, {
			name: "should return valid config if not sharded",
			args: args{
				assignedScanners: 6,
				scannerIndex:     4,
				chainSettings:    nil,
				chainID:          1,
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 0,
				Target:  6,
				Shards:  1,
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assignment := &registry.Assignment{
					AssignedScanners: tt.args.assignedScanners,
					ScannerIndex:     tt.args.scannerIndex,
				}

				agentManifest := &manifest.SignedAgentManifest{
					Manifest: &manifest.AgentManifest{
						ChainSettings: tt.args.chainSettings,
					},
				}

				shardConfig := populateShardConfig(assignment, agentManifest, tt.args.chainID)
				assert.Equal(t, tt.expectedShardConfig, shardConfig)
			},
		)
	}
}

func TestGetAgentsIfChanged(t *testing.T) {
	scanner := "your-scanner-id"

	// Set up the initial state of the registryStore
	// Modify rs.lastCompletedVersion, rs.lastUpdate, rs.loadedBots, and rs.invalidAssignments as needed
	testBot1 := "test-bot-1"
	testImage1 := "bafybeicc6ce3dnvjjfbrljtxuzncg2np76qkw5xq3w4af5x2c3m2nivwb4@sha256:5cf63050b113ce2df2a106b20d420c6687d30c28ed98cd42498f46475f642458"
	testManifest1 := "Qmex2rYHDsYqHcpSLhjow57MHBLpZMM1unPUSbPDYb5yTa"

	tests := []struct {
		name              string
		assignmentList    []*registry.Assignment
		registryClientErr error
		manifestClientErr error
		expectedAgents    []config.AgentConfig
		expectedUpdate    bool
		expectedErr       error
		manifest          *manifest.SignedAgentManifest
	}{
		{
			name:              "No update needed",
			registryClientErr: fmt.Errorf("failed to get assignments"),
			manifestClientErr: nil,
			expectedAgents:    nil,
			expectedUpdate:    false,
			expectedErr:       fmt.Errorf("failed to get assignments"),
		},
		{
			name:              "Update needed",
			registryClientErr: nil,
			manifestClientErr: nil,
			assignmentList: []*registry.Assignment{
				{
					AgentID:       "test-bot-1",
					AgentManifest: testManifest1,
				},
			},
			expectedAgents: []config.AgentConfig{
				{
					ID:          "test-bot-1",
					Image:       "/" + testImage1,
					Manifest:    testManifest1,
					ShardConfig: &config.ShardConfig{},
				},
			},
			manifest: &manifest.SignedAgentManifest{
				Manifest: &manifest.AgentManifest{
					AgentID:        &testBot1,
					ImageReference: &testImage1,
				},
			},
			expectedUpdate: true,
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				// Create mock objects
				mockRegistryClient := mock_registry.NewMockClient(ctrl)
				mockConfig := config.Config{}
				mockManifestClient := mock_manifest.NewMockClient(ctrl)

				rs := &registryStore{
					rc:                   mockRegistryClient,
					cfg:                  mockConfig,
					mc:                   mockManifestClient,
					lastCompletedVersion: "",
					lastUpdate:           time.Now().Add(-2 * time.Hour),
				}


				// Set up the expectations for the mock objects
				mockRegistryClient.EXPECT().GetAssignmentHash(scanner).Return(&registry.AssignmentHash{}, tt.registryClientErr).MaxTimes(1)
				mockRegistryClient.EXPECT().GetAssignmentList(gomock.Any(), gomock.Any(), scanner).Return(tt.assignmentList, tt.registryClientErr).MaxTimes(1)
				mockRegistryClient.EXPECT().PegLatestBlock().Return(tt.registryClientErr).Do(
					func() {
						return
					},
				).MaxTimes(1)
				mockRegistryClient.EXPECT().ResetOpts().Do(
					func() {
						return
					},
				).MaxTimes(1)

				mockManifestClient.EXPECT().GetAgentManifest(gomock.Any(), gomock.Any()).Return(tt.manifest, tt.manifestClientErr).MaxTimes(1)

				agents, update, err := rs.GetAgentsIfChanged(scanner)

				assert.Equal(t, len(tt.expectedAgents), len(agents))
				assert.Equal(t, tt.expectedUpdate, update)
				assert.Equal(t, tt.expectedErr, err)
			},
		)
	}
}