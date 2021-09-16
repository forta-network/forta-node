package query

import (
	"time"

	"github.com/forta-network/forta-node/protocol"
)

// Constants
const (
	DefaultFlushInterval = time.Minute
)

// AgentMetricsAggregator aggregates agent metrics and returns them when the available data
// hits up to a buffer limit.
type AgentMetricsAggregator struct {
	allMetrics    []*protocol.AgentMetrics
	flushInterval time.Duration
	lastFlush     time.Time
}

// NewAgentMetricsAggregator creates a new agent metrics aggregator.
func NewMetricsAggregator(flushIntervalSeconds int) *AgentMetricsAggregator {
	flushInterval := DefaultFlushInterval
	if flushIntervalSeconds > 0 {
		flushInterval = (time.Duration)(flushIntervalSeconds) * time.Second
	}

	return &AgentMetricsAggregator{
		flushInterval: flushInterval,
		lastFlush:     time.Now(), // avoid flushing immediately
	}
}

// PutTxProcessingData puts tx processing metric data of an agent.
func (ama *AgentMetricsAggregator) PutTxProcessingData(agentID string, data *protocol.MetricData) {
	for _, agentMetrics := range ama.allMetrics {
		if agentMetrics.AgentId == agentID {
			agentMetrics.TxProcessing = append(agentMetrics.TxProcessing, data)
			return
		}
	}
	ama.allMetrics = append(ama.allMetrics, &protocol.AgentMetrics{
		AgentId:      agentID,
		TxProcessing: []*protocol.MetricData{data},
	})
}

// PutBlockProcessingData puts block processing metric data of an agent.
func (ama *AgentMetricsAggregator) PutBlockProcessingData(agentID string, data *protocol.MetricData) {
	for _, agentMetrics := range ama.allMetrics {
		if agentMetrics.AgentId == agentID {
			agentMetrics.BlockProcessing = append(agentMetrics.BlockProcessing, data)
			return
		}
	}
	ama.allMetrics = append(ama.allMetrics, &protocol.AgentMetrics{
		AgentId:         agentID,
		BlockProcessing: []*protocol.MetricData{data},
	})
}

// TryFlush checks the flushing condition(s) an returns metrics accordingly.
func (ama *AgentMetricsAggregator) TryFlush() []*protocol.AgentMetrics {
	now := time.Now()
	if now.Sub(ama.lastFlush) < ama.flushInterval {
		return nil
	}
	ama.lastFlush = now
	allMetrics := ama.allMetrics
	ama.allMetrics = make([]*protocol.AgentMetrics, 0)
	return allMetrics
}
