package metrics

import (
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"
)

// FindActiveBotsFromMetrics finds the active bots from given bot metrics.
func FindActiveBotsFromMetrics(allBotMetrics []*protocol.AgentMetrics) (found []config.AgentConfig) {
	for _, botMetrics := range allBotMetrics {
		botID := botMetrics.AgentId
		for _, botMetric := range botMetrics.Metrics {
			if botMetric.Name == MetricHealthCheckSuccess {
				// copy over shardID value so metric will indicate shard
				cfg := &config.AgentConfig{ID: botID}
				if botMetric.ShardId >= 0 {
					cfg.ShardConfig = &config.ShardConfig{ShardID: uint(botMetric.ShardId)}
				}

				found = append(found, *cfg)
				break
			}
		}
	}
	return
}
