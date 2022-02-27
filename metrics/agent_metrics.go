package metrics

import (
	"time"

	"github.com/forta-protocol/forta-core-go/protocol"
	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
)

const (
	MetricFinding        = "finding"
	MetricTxRequest      = "tx.request"
	MetricTxLatency      = "tx.latency"
	MetricTxError        = "tx.error"
	MetricTxSuccess      = "tx.success"
	MetricTxDrop         = "tx.drop"
	MetricBlockRequest   = "block.request"
	MetricBlockLatency   = "block.latency"
	MetricBlockError     = "block.error"
	MetricBlockSuccess   = "block.success"
	MetricBlockDrop      = "block.drop"
	MetricStop           = "agent.stop"
	MetricJSONRPCRequest = "jsonrpc.request"
	MetricJSONRPCLatency = "jsonrpc.latency"
)

func SendAgentMetrics(client clients.MessageClient, ms []*protocol.AgentMetric) {
	if len(ms) > 0 {
		client.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{
			Metrics: ms,
		})
	}
}

func CreateAgentMetric(agentID, metric string, value float64) *protocol.AgentMetric {
	return &protocol.AgentMetric{
		AgentId:   agentID,
		Timestamp: time.Now().Format(time.RFC3339),
		Name:      metric,
		Value:     value,
	}
}

func createMetrics(agentID, timestamp string, metricMap map[string]float64) []*protocol.AgentMetric {
	var res []*protocol.AgentMetric
	for name, value := range metricMap {
		res = append(res, &protocol.AgentMetric{
			AgentId:   agentID,
			Timestamp: timestamp,
			Name:      name,
			Value:     value,
		})
	}
	return res
}

func GetBlockMetrics(agt config.AgentConfig, resp *protocol.EvaluateBlockResponse) []*protocol.AgentMetric {
	metrics := make(map[string]float64)

	metrics[MetricBlockRequest] = 1
	metrics[MetricFinding] = float64(len(resp.Findings))
	metrics[MetricBlockLatency] = float64(resp.LatencyMs)

	if resp.Status == protocol.ResponseStatus_ERROR {
		metrics[MetricBlockError] = 1
	} else if resp.Status == protocol.ResponseStatus_SUCCESS {
		metrics[MetricBlockSuccess] = 1
	}

	return createMetrics(agt.ID, resp.Timestamp, metrics)
}

func GetTxMetrics(agt config.AgentConfig, resp *protocol.EvaluateTxResponse) []*protocol.AgentMetric {
	metrics := make(map[string]float64)

	metrics[MetricTxRequest] = 1
	metrics[MetricFinding] = float64(len(resp.Findings))
	metrics[MetricTxLatency] = float64(resp.LatencyMs)

	if resp.Status == protocol.ResponseStatus_ERROR {
		metrics[MetricTxError] = 1
	} else if resp.Status == protocol.ResponseStatus_SUCCESS {
		metrics[MetricTxSuccess] = 1
	}

	return createMetrics(agt.ID, resp.Timestamp, metrics)
}

func GetJSONRPCMetrics(agt config.AgentConfig, at time.Time, latencyMs time.Duration) []*protocol.AgentMetric {
	return createMetrics(agt.ID, at.Format(time.RFC3339), map[string]float64{
		MetricJSONRPCRequest: 1,
		MetricJSONRPCLatency: float64(latencyMs.Milliseconds()),
	})
}
