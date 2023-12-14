package sharding

import (
	"sort"
	"strconv"

	"github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-node/config"
)

// CalculateShardConfig calculates the sharding parameters and returns a config.
// It looks at the assignments across all chains and deterministically finds out which chain
// and which shard.
func CalculateShardConfigV2(
	assignment *registry.Assignment, agentManifest *manifest.SignedAgentManifest,
) (*config.ShardConfig, bool) {
	// we think of each chain in ascending order
	chainIDs := agentManifest.Manifest.ChainIDs
	sort.Slice(chainIDs, func(i, j int) bool {
		return chainIDs[i] < chainIDs[j]
	})

	index := assignment.AllChainsScannerIndex
	for _, chainID := range chainIDs {
		var target, shards uint

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

		if target == 0 && shards == 0 {
			shards = 1
			target = defaultTargetCount
		}

		// check if this scanner's index is not within the current chain's portion of assignments
		// then we cannot find the shard id from this chain
		size := int(target * shards)
		if index >= size {
			// fix the index as we are moving onto the next portion
			index -= size
			continue
		}

		shardID := CalculateShardID(target, uint(index))
		return CreateShardConfig(shardID, shards, target, chainID), true
	}
	// invalid shard config
	return nil, false
}
