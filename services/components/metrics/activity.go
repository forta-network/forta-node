package metrics

import "github.com/forta-network/forta-core-go/protocol"

// FindActiveBotsFromMetrics finds the active bots from given bot metrics.
func FindActiveBotsFromMetrics(allBotMetrics []*protocol.AgentMetrics) (found []string) {
	for _, botMetrics := range allBotMetrics {
		botID := botMetrics.AgentId
		for _, botMetric := range botMetrics.Metrics {
			if botMetric.Name == MetricTxLatency ||
				botMetric.Name == MetricBlockLatency ||
				botMetric.Name == MetricCombinerLatency {
				found = append(found, botID)
				break
			}
		}
	}
	return
}
