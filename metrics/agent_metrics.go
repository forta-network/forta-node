package metrics

import (
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/protocol"
)

const (
	MetricFinding      = "finding"
	MetricTxRequest    = "tx.request"
	MetricTxLatency    = "tx.latency"
	MetricTxError      = "tx.error"
	MetricTxSuccess    = "tx.success"
	MetricBlockRequest = "block.request"
	MetricBlockLatency = "block.latency"
	MetricBlockError   = "block.error"
	MetricBlockSuccess = "block.success"
	MetricStop         = "agent.stop"
)

func createMetrics(agentID, timestamp string, metricMap map[string]float64) []protocol.AgentMetric {
	var res []protocol.AgentMetric
	for name, value := range metricMap {
		res = append(res, protocol.AgentMetric{
			AgentId:   agentID,
			Timestamp: timestamp,
			Name:      name,
			Value:     value,
		})
	}
	return res
}

func GetBlockMetrics(agt config.AgentConfig, resp *protocol.EvaluateBlockResponse) []protocol.AgentMetric {
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

func GetTxMetrics(agt config.AgentConfig, resp *protocol.EvaluateTxResponse) []protocol.AgentMetric {
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
