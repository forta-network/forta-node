package metrics

import (
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/domain"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
)

func SendAgentMetrics(client clients.MessageClient, ms []*protocol.AgentMetric) {
	if len(ms) > 0 {
		client.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{
			Metrics: ms,
		})
	}
}

func CreateAgentMetricV1(agt config.AgentConfig, metric string, value float64) *protocol.AgentMetric {
	return &protocol.AgentMetric{
		AgentId:   agt.ID,
		Timestamp: time.Now().Format(time.RFC3339),
		Name:      metric,
		Value:     value,
		ShardId:   agt.ShardID(),
		ChainId:   int64(agt.ChainID),
	}
}

func CreateAgentMetricV2(agt config.AgentConfig, metric string, value float64, chainID int64) *protocol.AgentMetric {
	return &protocol.AgentMetric{
		AgentId:   agt.ID,
		Timestamp: time.Now().Format(time.RFC3339),
		Name:      metric,
		Value:     value,
		ShardId:   agt.ShardID(),
		ChainId:   chainID,
	}
}

func CreateDetailedAgentMetricV2(agt config.AgentConfig, metric string, value float64, details string, chainID int64) *protocol.AgentMetric {
	return &protocol.AgentMetric{
		AgentId:   agt.ID,
		Timestamp: time.Now().Format(time.RFC3339),
		Name:      metric,
		Value:     value,
		ShardId:   agt.ShardID(),
		ChainId:   chainID,
		Details:   details,
	}
}

func CreateEventMetric(t time.Time, id string, metric string, details string) *protocol.AgentMetric {
	return &protocol.AgentMetric{
		AgentId:   id,
		Timestamp: t.Format(time.RFC3339),
		Name:      metric,
		Value:     1,
		Details:   details,
	}
}

func CreateSystemMetric(metric string, value float64, details string) *protocol.AgentMetric {
	return &protocol.AgentMetric{
		AgentId:   "system",
		Timestamp: time.Now().Format(time.RFC3339),
		Name:      metric,
		Value:     value,
		Details:   details,
	}
}

func CreateAgentResourcesMetric(agt config.AgentConfig, t time.Time, metric string, value float64) *protocol.AgentMetric {
	return &protocol.AgentMetric{
		AgentId:   agt.ID,
		ShardId:   agt.ShardID(),
		ChainId:   int64(agt.ChainID),
		Timestamp: t.Format(time.RFC3339),
		Name:      metric,
		Value:     value,
	}
}

func createMetrics(agt config.AgentConfig, timestamp string, metricMap map[string]float64) []*protocol.AgentMetric {
	var res []*protocol.AgentMetric

	for name, value := range metricMap {
		res = append(res, &protocol.AgentMetric{
			AgentId:   agt.ID,
			Timestamp: timestamp,
			Name:      name,
			Value:     value,
			ShardId:   agt.ShardID(),
			ChainId:   int64(agt.ChainID),
		})
	}
	return res
}

func durationMs(from time.Time, to time.Time) float64 {
	return float64(to.Sub(from).Milliseconds())
}

func GetBlockMetrics(agt config.AgentConfig, resp *protocol.EvaluateBlockResponse, times *domain.TrackingTimestamps) []*protocol.AgentMetric {
	metrics := make(map[string]float64)

	metrics[domain.MetricBlockRequest] = 1
	metrics[domain.MetricFinding] = float64(len(resp.Findings))
	metrics[domain.MetricBlockLatency] = float64(resp.LatencyMs)
	metrics[domain.MetricBlockBlockAge] = durationMs(times.Block, times.BotRequest)
	metrics[domain.MetricBlockEventAge] = durationMs(times.Feed, times.BotRequest)

	if resp.Status == protocol.ResponseStatus_ERROR {
		metrics[domain.MetricBlockError] = 1
	} else if resp.Status == protocol.ResponseStatus_SUCCESS {
		metrics[domain.MetricBlockSuccess] = 1
	}

	return createMetrics(agt, resp.Timestamp, metrics)
}

func GetTxMetrics(agt config.AgentConfig, resp *protocol.EvaluateTxResponse, times *domain.TrackingTimestamps) []*protocol.AgentMetric {
	metrics := make(map[string]float64)

	metrics[domain.MetricTxRequest] = 1
	metrics[domain.MetricFinding] = float64(len(resp.Findings))
	metrics[domain.MetricTxLatency] = float64(resp.LatencyMs)
	metrics[domain.MetricTxBlockAge] = durationMs(times.Block, times.BotRequest)
	metrics[domain.MetricTxEventAge] = durationMs(times.Feed, times.BotRequest)

	if resp.Status == protocol.ResponseStatus_ERROR {
		metrics[domain.MetricTxError] = 1
	} else if resp.Status == protocol.ResponseStatus_SUCCESS {
		metrics[domain.MetricTxSuccess] = 1
	}

	return createMetrics(agt, resp.Timestamp, metrics)
}

func GetCombinerMetrics(agt config.AgentConfig, resp *protocol.EvaluateAlertResponse, times *domain.TrackingTimestamps) []*protocol.AgentMetric {
	metrics := make(map[string]float64)

	metrics[domain.MetricCombinerRequest] = 1
	metrics[domain.MetricFinding] = float64(len(resp.Findings))
	metrics[domain.MetricCombinerLatency] = float64(resp.LatencyMs)

	if resp.Status == protocol.ResponseStatus_ERROR {
		metrics[domain.MetricCombinerError] = 1
	} else if resp.Status == protocol.ResponseStatus_SUCCESS {
		metrics[domain.MetricCombinerSuccess] = 1
	}

	return createMetrics(agt, resp.Timestamp, metrics)
}

func GetJSONRPCMetrics(agt config.AgentConfig, at time.Time, success, throttled int, latencyMs time.Duration, method string) []*protocol.AgentMetric {
	values := make(map[string]float64)
	if latencyMs > 0 {
		values[domain.MetricJSONRPCLatency] = float64(latencyMs.Milliseconds())
	}
	if success > 0 {
		values[domain.MetricJSONRPCSuccess] = float64(success)
		values[domain.MetricJSONRPCRequest] += float64(success)
	}
	if throttled > 0 {
		values[domain.MetricJSONRPCThrottled] = float64(throttled)
		values[domain.MetricJSONRPCRequest] += float64(throttled)
	}
	return createJsonRpcMetrics(agt, at.Format(time.RFC3339), values, method)
}

func createJsonRpcMetrics(agt config.AgentConfig, timestamp string, metricMap map[string]float64, method string) []*protocol.AgentMetric {
	var res []*protocol.AgentMetric

	for name, value := range metricMap {
		res = append(res, &protocol.AgentMetric{
			AgentId:   agt.ID,
			Timestamp: timestamp,
			Name:      fmt.Sprintf("%s.%s", name, method),
			Value:     value,
			ShardId:   agt.ShardID(),
			ChainId:   int64(agt.ChainID),
		})
	}
	return res
}

func GetPublicAPIMetrics(botID string, at time.Time, success, throttled int, latency time.Duration) []*protocol.AgentMetric {
	values := make(map[string]float64)
	if latency > 0 {
		values[domain.MetricPublicAPIProxyLatency] = float64(latency.Milliseconds())
	}
	if success > 0 {
		values[domain.MetricPublicAPIProxySuccess] = float64(success)
		values[domain.MetricPublicAPIProxyRequest] += float64(success)
	}
	if throttled > 0 {
		values[domain.MetricPublicAPIProxyThrottled] = float64(throttled)
		values[domain.MetricPublicAPIProxyRequest] += float64(throttled)
	}
	//TODO: get the shardID into this eventually
	return createMetrics(config.AgentConfig{ID: botID}, at.Format(time.RFC3339), values)
}
