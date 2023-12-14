package sharding

import (
	"testing"

	"github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-node/config"
	"github.com/stretchr/testify/require"
)

func TestCalculateShardConfig(t *testing.T) {
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
				chainID: 1,
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 2,
				Target:  2,
				Shards:  3,
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
				chainID: 1,
			},
			expectedShardConfig: &config.ShardConfig{
				ShardID: 0,
				Target:  6,
				Shards:  1,
			},
		},
		{
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
				r := require.New(t)

				assignment := &registry.Assignment{
					ScannerIndices: registry.ScannerIndices{
						SameChainAssignedScanners: tt.args.assignedScanners,
						SameChainScannerIndex:     tt.args.scannerIndex,
					},
				}

				agentManifest := &manifest.SignedAgentManifest{
					Manifest: &manifest.AgentManifest{
						ChainSettings: tt.args.chainSettings,
					},
				}

				shardConfig := CalculateShardConfig(assignment, agentManifest, tt.args.chainID)
				r.Equal(tt.expectedShardConfig, shardConfig)
			},
		)
	}
}
