package sharding

import "github.com/forta-network/forta-node/config"

const (
	keyDefaultChainSetting = "default"
	defaultShardID         = 0
	minShardCount          = 1
	defaultTargetCount     = 3
)

// CreateShardConfig creates a new shard config by using given arguments.
func CreateShardConfig(shardID, shards, target uint, chainID int64) *config.ShardConfig {
	return &config.ShardConfig{
		ShardID: shardID,
		Target:  target,
		Shards:  shards,
		ChainID: chainID,
	}
}

// CalculateShardID returns shard ID for an index in an assignment array.
// The shard IDs are considered to be distributed evenly in an increased order.
// Example:
// Target: 6, Shards: 3 (6 × 3 = 18 assignments)
// The assignment array by shard IDs should be [0,0,0,0,0,0,1,1,1,1,1,1,2,2,2,2,2,2]
func CalculateShardID(target, idx uint) uint {
	if target == 0 {
		return defaultShardID
	}
	return idx / target
}
