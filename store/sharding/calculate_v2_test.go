package sharding

import (
	"testing"

	"github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-node/config"
	"github.com/stretchr/testify/require"
)

func TestCalculateShardConfigV2(t *testing.T) {
	type args struct {
		assignedScanners int
		scannerIndex     int
		chainSettings    map[string]manifest.AgentChainSettings
		chainIDs         []int64
	}
	tests := []struct {
		name                string
		args                args
		expectedShardConfig *config.ShardConfig
		invalid             bool
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
				chainIDs: []int64{1},
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 2,
				Target:  2,
				Shards:  3,
				ChainID: 1,
			},
		},
		{
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
				chainIDs: []int64{1},
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 2,
				Target:  2,
				Shards:  3,
				ChainID: 1,
			},
		},
		{
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
				chainIDs: []int64{1},
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 0,
				Target:  6,
				Shards:  1,
				ChainID: 1,
			},
		},
		{
			name: "should return valid config if not sharded",
			args: args{
				assignedScanners: 6,
				scannerIndex:     2,
				chainSettings:    nil,
				chainIDs:         []int64{1},
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 0,
				Target:  3,
				Shards:  1,
				ChainID: 1,
			},
		},
		{
			name: "should find the second chain id",
			args: args{
				assignedScanners: 6,
				scannerIndex:     4,
				chainSettings:    nil,
				chainIDs:         []int64{1, 2},
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 0,
				Target:  3,
				Shards:  1,
				ChainID: 2,
			},
		},
		{
			name: "should find the shard id for the second chain id",
			args: args{
				assignedScanners: 15, // this should be 12 but let's assume that there's an assignment mistake
				scannerIndex:     11,
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
				chainIDs: []int64{1, 2},
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 2,
				Target:  2,
				Shards:  3,
				ChainID: 2,
			},
		},
		{
			name: "should be invalid if index too large",
			args: args{
				// chain settings are nil, there are two chains,
				// there must be 2 chains × 3 target = 6 assignments
				assignedScanners: 12, // unexpectedly many assignments
				scannerIndex:     8,  // unexpectedly large index
				chainSettings:    nil,
				chainIDs:         []int64{1, 2},
			},
			expectedShardConfig: nil,
			invalid:             true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				r := require.New(t)

				assignment := &registry.Assignment{
					ScannerIndices: registry.ScannerIndices{
						AllChainsAssignedScanners: tt.args.assignedScanners,
						AllChainsScannerIndex:     tt.args.scannerIndex,
					},
				}

				agentManifest := &manifest.SignedAgentManifest{
					Manifest: &manifest.AgentManifest{
						ChainIDs:      tt.args.chainIDs,
						ChainSettings: tt.args.chainSettings,
					},
				}

				shardConfig, ok := CalculateShardConfigV2(assignment, agentManifest)
				r.Equal(!tt.invalid, ok)
				r.Equal(tt.expectedShardConfig, shardConfig)
			},
		)
	}
}
