package sharding

import (
	"strconv"

	"github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-node/config"
)

// CalculateShardConfig calculates the sharding parameters and returns a config.
// It looks at the same chain as the current scanner when selecting assignment values.
func CalculateShardConfig(
	assignment *registry.Assignment, agentManifest *manifest.SignedAgentManifest, chainID int,
) *config.ShardConfig {
	assignedScanners := assignment.SameChainAssignedScanners
	scannerIndex := assignment.SameChainScannerIndex

	var (
		target, shards uint
	)

	// check if there is a default chain setting
	chainSetting, ok := agentManifest.Manifest.ChainSettings[keyDefaultChainSetting]
	// if not a sharded bot, shard is always 0
	if ok {
		target = chainSetting.Target
		shards = chainSetting.Shards
	}

	// check if there is a chain setting for the scanner's chain
	chainIDStr := strconv.FormatInt(int64(chainID), 10)
	chainSetting, ok = agentManifest.Manifest.ChainSettings[chainIDStr]
	// if not a sharded bot, shard is always 0
	if ok {
		target = chainSetting.Target
		shards = chainSetting.Shards
	}

	// if no sharding specified, shard count is 1 and target is total assigns
	if shards == 0 {
		target = uint(assignment.SameChainAssignedScanners)

		return CreateShardConfig(defaultShardID, minShardCount, target, 0)
	}

	// fallback for target, calculate it from shard to assign ratio.
	// target defaults to total assigns / shards
	if target == 0 && shards != 0 {
		target = uint(uint64(assignedScanners) / uint64(shards))
	}

	shardID := CalculateShardID(target, uint(scannerIndex))

	return CreateShardConfig(shardID, shards, target, 0)
}
