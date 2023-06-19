package store

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/manifest"
	mock_manifest "github.com/forta-network/forta-core-go/manifest/mocks"
	"github.com/forta-network/forta-core-go/registry"
	mock_registry "github.com/forta-network/forta-core-go/registry/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testScannerID = "your-scanner-id"
	testBot1      = "test-bot-1"
	testImage1    = "bafybeicc6ce3dnvjjfbrljtxuzncg2np76qkw5xq3w4af5x2c3m2nivwb4@sha256:5cf63050b113ce2df2a106b20d420c6687d30c28ed98cd42498f46475f642458"
	testManifest1 = "Qmex2rYHDsYqHcpSLhjow57MHBLpZMM1unPUSbPDYb5yTa"
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

func TestGetAgentsIfChanged_UpdateNeeded(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)

	regClient := mock_registry.NewMockClient(ctrl)
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
	testConfig := config.Config{
		ChainID: 123,
	}
	rs := &registryStore{
		rc:                   regClient,
		cfg:                  testConfig,
		bms:                  NewBotManifestStore(manifestClient),
		lastCompletedVersion: "",
		lastUpdate:           time.Now().Add(-2 * time.Hour), // test forced timeout
	}
	assignmentList := []*registry.Assignment{
		{
			AgentID:          "test-bot-1",
			AgentManifest:    testManifest1,
			AssignedScanners: 6,
			ScannerIndex:     1, // initial order of the scanner
		},
	}

	// it should load for the first time

	regClient.EXPECT().GetAssignmentHash(testScannerID).Return(&registry.AssignmentHash{}, nil)
	regClient.EXPECT().GetAssignmentList(
		gomock.Any(), big.NewInt(int64(testConfig.ChainID)), testScannerID,
	).Return(assignmentList, nil)
	regClient.EXPECT().PegLatestBlock().Return(nil)
	regClient.EXPECT().ResetOpts()
	manifestClient.EXPECT().GetAgentManifest(gomock.Any(), gomock.Any()).Return(
		testManifest, nil,
	)

	agents, update, err := rs.GetAgentsIfChanged(testScannerID)

	r.NoError(err)
	r.True(update)
	r.Len(agents, len(assignmentList))
	firstShardID := agents[0].ShardConfig.ShardID
	r.Equal(assignmentList[0].ScannerIndex, int(firstShardID))

	// it should update the shard id when ordering changes

	updatedAssignmentList := []*registry.Assignment{
		{
			AgentID:          "test-bot-1",
			AgentManifest:    testManifest1,
			AssignedScanners: 6,
			ScannerIndex:     2, // scanner order change: 1 => 2
		},
	}

	regClient.EXPECT().GetAssignmentHash(testScannerID).Return(&registry.AssignmentHash{
		Hash: "some-changed-version-hash",
	}, nil)
	regClient.EXPECT().GetAssignmentList(
		gomock.Any(), big.NewInt(int64(testConfig.ChainID)), testScannerID,
	).Return(updatedAssignmentList, nil)
	regClient.EXPECT().PegLatestBlock().Return(nil)
	regClient.EXPECT().ResetOpts()
	manifestClient.EXPECT().GetAgentManifest(gomock.Any(), gomock.Any()).Return(
		testManifest, nil,
	).Times(0) // expect no calls to this: it should use the cache

	agents, update, err = rs.GetAgentsIfChanged(testScannerID)

	r.NoError(err)
	r.True(update)
	r.Len(agents, len(assignmentList))
	secondShardID := agents[0].ShardConfig.ShardID
	r.Equal(updatedAssignmentList[0].ScannerIndex, int(secondShardID))

	r.NotEqual(firstShardID, secondShardID)
}

func TestGetAgentsIfChanged_NoUpdateNeeded(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)

	regClient := mock_registry.NewMockClient(ctrl)
	manifestClient := mock_manifest.NewMockClient(ctrl)

	testConfig := config.Config{
		ChainID: 123,
	}
	rs := &registryStore{
		rc:                   regClient,
		cfg:                  testConfig,
		bms:                  NewBotManifestStore(manifestClient),
		lastCompletedVersion: "version-hash-1",
		lastUpdate:           time.Now(),
	}

	regClient.EXPECT().GetAssignmentHash(testScannerID).Return(&registry.AssignmentHash{
		Hash: "version-hash-2", // test hash difference
	}, nil)
	regClient.EXPECT().PegLatestBlock().Return(nil)
	regClient.EXPECT().ResetOpts()
	regClient.EXPECT().GetAssignmentList(
		gomock.Any(), big.NewInt(int64(testConfig.ChainID)), testScannerID,
	).Return(nil, errors.New("failed to get assignment list"))

	agents, update, err := rs.GetAgentsIfChanged(testScannerID)

	r.Error(err)
	r.False(update)
	r.Nil(agents)
}
