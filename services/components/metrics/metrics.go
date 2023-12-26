package metrics

import (
	"time"

	"github.com/forta-network/forta-core-go/domain"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
)

const (
	MetricFinding       = "finding"
	MetricTxRequest     = "tx.request"
	MetricTxLatency     = "tx.latency"
	MetricTxError       = "tx.error"
	MetricTxSuccess     = "tx.success"
	MetricTxDrop        = "tx.drop"
	MetricTxBlockAge    = "tx.block.age"
	MetricTxEventAge    = "tx.event.age"
	MetricBlockBlockAge = "block.block.age"
	MetricBlockEventAge = "block.event.age"
	MetricBlockRequest  = "block.request"
	MetricBlockLatency  = "block.latency"
	MetricBlockError    = "block.error"
	MetricBlockSuccess  = "block.success"
	MetricBlockDrop     = "block.drop"

	MetricJSONRPCLatency          = "jsonrpc.latency"
	MetricJSONRPCRequest          = "jsonrpc.request"
	MetricJSONRPCSuccess          = "jsonrpc.success"
	MetricJSONRPCThrottled        = "jsonrpc.throttled"
	MetricPublicAPIProxyLatency   = "publicapi.latency"
	MetricPublicAPIProxyRequest   = "publicapi.request"
	MetricPublicAPIProxySuccess   = "publicapi.success"
	MetricPublicAPIProxyThrottled = "publicapi.throttled"
	MetricFindingsDropped         = "findings.dropped"
	MetricCombinerRequest         = "combiner.request"
	MetricCombinerLatency         = "combiner.latency"
	MetricCombinerError           = "combiner.error"
	MetricCombinerSuccess         = "combiner.success"
	MetricCombinerDrop            = "combiner.drop"
)

func SendAgentMetrics(client clients.MessageClient, ms []*protocol.AgentMetric) {
	if len(ms) > 0 {
		client.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{
			Metrics: ms,
		})
	}
}

func CreateAgentMetric(agt config.AgentConfig, metric string, value float64) *protocol.AgentMetric {
	return &protocol.AgentMetric{
		AgentId:   agt.ID,
		Timestamp: time.Now().Format(time.RFC3339),
		Name:      metric,
		Value:     value,
		ShardId:   agt.ShardID(),
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

func createMetrics(agt config.AgentConfig, timestamp string, metricMap map[string]float64, details string) []*protocol.AgentMetric {
	var res []*protocol.AgentMetric

	for name, value := range metricMap {
		res = append(res, &protocol.AgentMetric{
			AgentId:   agt.ID,
			Timestamp: timestamp,
			Name:      name,
			Value:     value,
			ShardId:   agt.ShardID(),
			Details:   details,
		})
	}
	return res
}

func durationMs(from time.Time, to time.Time) float64 {
	return float64(to.Sub(from).Milliseconds())
}

func GetBlockMetrics(agt config.AgentConfig, resp *protocol.EvaluateBlockResponse, times *domain.TrackingTimestamps) []*protocol.AgentMetric {
	metrics := make(map[string]float64)

	metrics[MetricBlockRequest] = 1
	metrics[MetricFinding] = float64(len(resp.Findings))
	metrics[MetricBlockLatency] = float64(resp.LatencyMs)
	metrics[MetricBlockBlockAge] = durationMs(times.Block, times.BotRequest)
	metrics[MetricBlockEventAge] = durationMs(times.Feed, times.BotRequest)

	if resp.Status == protocol.ResponseStatus_ERROR {
		metrics[MetricBlockError] = 1
	} else if resp.Status == protocol.ResponseStatus_SUCCESS {
		metrics[MetricBlockSuccess] = 1
	}

	return createMetrics(agt, resp.Timestamp, metrics, "")
}

func GetTxMetrics(agt config.AgentConfig, resp *protocol.EvaluateTxResponse, times *domain.TrackingTimestamps) []*protocol.AgentMetric {
	metrics := make(map[string]float64)

	metrics[MetricTxRequest] = 1
	metrics[MetricFinding] = float64(len(resp.Findings))
	metrics[MetricTxLatency] = float64(resp.LatencyMs)
	metrics[MetricTxBlockAge] = durationMs(times.Block, times.BotRequest)
	metrics[MetricTxEventAge] = durationMs(times.Feed, times.BotRequest)

	if resp.Status == protocol.ResponseStatus_ERROR {
		metrics[MetricTxError] = 1
	} else if resp.Status == protocol.ResponseStatus_SUCCESS {
		metrics[MetricTxSuccess] = 1
	}

	return createMetrics(agt, resp.Timestamp, metrics, "")
}

func GetCombinerMetrics(agt config.AgentConfig, resp *protocol.EvaluateAlertResponse, times *domain.TrackingTimestamps) []*protocol.AgentMetric {
	metrics := make(map[string]float64)

	metrics[MetricCombinerRequest] = 1
	metrics[MetricFinding] = float64(len(resp.Findings))
	metrics[MetricCombinerLatency] = float64(resp.LatencyMs)

	if resp.Status == protocol.ResponseStatus_ERROR {
		metrics[MetricCombinerError] = 1
	} else if resp.Status == protocol.ResponseStatus_SUCCESS {
		metrics[MetricCombinerSuccess] = 1
	}

	return createMetrics(agt, resp.Timestamp, metrics, "")
}

func GetJSONRPCMetrics(agt config.AgentConfig, at time.Time, success, throttled int, latencyMs time.Duration, method string) []*protocol.AgentMetric {
	values := make(map[string]float64)
	if latencyMs > 0 {
		values[MetricJSONRPCLatency] = float64(latencyMs.Milliseconds())
	}
	if success > 0 {
		values[MetricJSONRPCSuccess] = float64(success)
		values[MetricJSONRPCRequest] += float64(success)
	}
	if throttled > 0 {
		values[MetricJSONRPCThrottled] = float64(throttled)
		values[MetricJSONRPCRequest] += float64(throttled)
	}
	return createMetrics(agt, at.Format(time.RFC3339), values, method)
}

func GetPublicAPIMetrics(botID string, at time.Time, success, throttled int, latencyMs time.Duration) []*protocol.AgentMetric {
	values := make(map[string]float64)
	if latencyMs > 0 {
		values[MetricPublicAPIProxyLatency] = float64(latencyMs.Milliseconds())
	}
	if success > 0 {
		values[MetricPublicAPIProxySuccess] = float64(success)
		values[MetricPublicAPIProxyRequest] += float64(success)
	}
	if throttled > 0 {
		values[MetricPublicAPIProxyThrottled] = float64(throttled)
		values[MetricPublicAPIProxyRequest] += float64(throttled)
	}
	//TODO: get the shardID into this eventually
	return createMetrics(config.AgentConfig{ID: botID}, at.Format(time.RFC3339), values, "")
}
