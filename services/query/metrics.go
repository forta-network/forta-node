package query

import (
	"math"
	"time"

	"github.com/forta-network/forta-node/protocol"
)

// Constants
const (
	DefaultFlushInterval              = time.Minute
	DefaultMetricsThresholdDurationMs = 50
)

// AgentMetricsAggregator aggregates agent metrics and returns them when the available data
// hits up to a buffer limit.
type AgentMetricsAggregator struct {
	allMetrics    []*metricsContainer
	flushInterval time.Duration
	lastFlush     time.Time
	thresholdMs   int
}

type metricsContainer struct {
	FindingCount   int64
	NoFindingCount int64
	*protocol.AgentMetrics
}

// NewAgentMetricsAggregator creates a new agent metrics aggregator.
func NewMetricsAggregator(flushIntervalSeconds, thresholdMs int) *AgentMetricsAggregator {
	flushInterval := DefaultFlushInterval
	if flushIntervalSeconds > 0 {
		flushInterval = (time.Duration)(flushIntervalSeconds) * time.Second
	}

	threshold := DefaultMetricsThresholdDurationMs
	if thresholdMs > 0 {
		threshold = thresholdMs
	}

	return &AgentMetricsAggregator{
		flushInterval: flushInterval,
		lastFlush:     time.Now(), // avoid flushing immediately
		thresholdMs:   threshold,
	}
}

// PutTxProcessingData puts tx processing metric data of an agent.
func (ama *AgentMetricsAggregator) PutTxProcessingData(agentID string, data *protocol.MetricData) {
	if data == nil {
		return
	}

	for _, agentMetrics := range ama.allMetrics {
		if agentMetrics.AgentMetrics.AgentId == agentID {
			agentMetrics.AgentMetrics.TxProcessing.Data = append(agentMetrics.AgentMetrics.TxProcessing.Data, data)
			return
		}
	}
	ama.allMetrics = append(ama.allMetrics, &metricsContainer{
		AgentMetrics: &protocol.AgentMetrics{
			AgentId:      agentID,
			TxProcessing: &protocol.MetricContainer{Data: []*protocol.MetricData{data}},
		},
	})
}

// PutBlockProcessingData puts block processing metric data of an agent.
func (ama *AgentMetricsAggregator) PutBlockProcessingData(agentID string, data *protocol.MetricData) {
	if data == nil {
		return
	}

	for _, agentMetrics := range ama.allMetrics {
		if agentMetrics.AgentMetrics.AgentId == agentID {
			agentMetrics.AgentMetrics.BlockProcessing.Data = append(agentMetrics.AgentMetrics.BlockProcessing.Data, data)
			return
		}
	}
	ama.allMetrics = append(ama.allMetrics, &metricsContainer{
		AgentMetrics: &protocol.AgentMetrics{
			AgentId:         agentID,
			BlockProcessing: &protocol.MetricContainer{Data: []*protocol.MetricData{data}},
		},
	})
}

// CountFinding increases the right counter depending on the nil or existing alert.
func (ama *AgentMetricsAggregator) CountFinding(agentID string, hasAlert bool) {
	for _, agentMetrics := range ama.allMetrics {
		if agentMetrics.AgentId != agentID {
			continue
		}
		if hasAlert {
			agentMetrics.FindingCount++
		} else {
			agentMetrics.NoFindingCount++
		}
		return
	}
}

// TryFlush checks the flushing condition(s) an returns metrics accordingly.
func (ama *AgentMetricsAggregator) TryFlush() []*protocol.AgentMetrics {
	now := time.Now()
	if now.Sub(ama.lastFlush) < ama.flushInterval {
		return nil
	}

	(allAgentMetrics)(ama.allMetrics).Fix(ama.thresholdMs)

	ama.lastFlush = now
	allContainers := ama.allMetrics
	ama.allMetrics = make([]*metricsContainer, 0)

	var allMetrics []*protocol.AgentMetrics
	for _, container := range allContainers {
		allMetrics = append(allMetrics, container.AgentMetrics)
	}

	return allMetrics
}

// allAgentMetrics is an alias type for post-processing aggregated in-memory metrics
// before we publish them.
type allAgentMetrics []*metricsContainer

func (allMetrics allAgentMetrics) Fix(thresholdMs int) {
	allMetrics.CalculateAverages()
	allMetrics.RemoveLowValues(thresholdMs)
}

func (allMetrics allAgentMetrics) CalculateAverages() {
	for _, agentMetrics := range allMetrics {
		agentMetrics.TxProcessing.Average = avgMetricArray(agentMetrics.TxProcessing.Data)
		agentMetrics.BlockProcessing.Average = avgMetricArray(agentMetrics.BlockProcessing.Data)
	}
}

func avgMetricArray(data []*protocol.MetricData) int64 {
	var sum int64
	for _, dataPoint := range data {
		sum += dataPoint.Number
	}
	return sum / int64(len(data))
}

func (allMetrics allAgentMetrics) RemoveLowValues(thresholdMs int) {
	for _, agentMetrics := range allMetrics {
		agentMetrics.TxProcessing.Data = reduceMetricArray(agentMetrics.TxProcessing.Data, thresholdMs)
		agentMetrics.BlockProcessing.Data = reduceMetricArray(agentMetrics.BlockProcessing.Data, thresholdMs)
	}
}

func reduceMetricArray(oldData []*protocol.MetricData, threshold int) (newData []*protocol.MetricData) {
	for _, dataPoint := range oldData {
		if dataPoint.Number >= int64(threshold) {
			newData = append(newData, dataPoint)
		}
	}
	return
}

func (allMetrics allAgentMetrics) CalculateFindingRates() {
	for _, agentMetrics := range allMetrics {
		agentMetrics.FindingRatePct = float32(calculateRatePct(agentMetrics.FindingCount, agentMetrics.FindingCount+agentMetrics.NoFindingCount))
	}
}

func calculateRatePct(dividend, divisor int64) float64 {
	const decimals = 100
	const toPct = 100
	res := float64(dividend*decimals*toPct) / float64(divisor)
	return math.Round(res) / float64(decimals)
}
