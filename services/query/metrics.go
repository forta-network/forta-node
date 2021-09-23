package query

import (
	"math"
	"sort"
	"time"

	"github.com/forta-network/forta-node/protocol"
	"github.com/shopspring/decimal"
)

// Constants
const (
	DefaultFlushInterval = time.Minute
)

// AgentMetricsAggregator aggregates agent metrics and returns them when the available data
// hits up to a buffer limit.
type AgentMetricsAggregator struct {
	allMetrics    []*metricsContainer
	flushInterval time.Duration
	lastFlush     time.Time
}

type metricsContainer struct {
	FindingCount    int64
	NoFindingCount  int64
	TxProcessing    []*protocol.MetricData
	BlockProcessing []*protocol.MetricData
	protocol.AgentMetrics
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
	if data == nil {
		return
	}

	for _, agentMetrics := range ama.allMetrics {
		if agentMetrics.AgentId == agentID {
			agentMetrics.TxProcessing = append(agentMetrics.TxProcessing, data)
			return
		}
	}
	ama.allMetrics = append(ama.allMetrics, &metricsContainer{
		TxProcessing: []*protocol.MetricData{data},
		AgentMetrics: protocol.AgentMetrics{
			AgentId: agentID,
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
			agentMetrics.BlockProcessing = append(agentMetrics.BlockProcessing, data)
			return
		}
	}
	ama.allMetrics = append(ama.allMetrics, &metricsContainer{
		BlockProcessing: []*protocol.MetricData{data},
		AgentMetrics: protocol.AgentMetrics{
			AgentId: agentID,
		},
	})
}

// CountFinding increases the right counter depending on the nil or existing alert.
func (ama *AgentMetricsAggregator) CountFinding(agentID string, hasAlert bool) {
	for _, agentMetrics := range ama.allMetrics {
		if agentMetrics.AgentId == agentID {
			if hasAlert {
				agentMetrics.FindingCount++
			} else {
				agentMetrics.NoFindingCount++
			}
			return
		}
	}
	if hasAlert {
		ama.allMetrics = append(ama.allMetrics, &metricsContainer{
			FindingCount: 1,
			AgentMetrics: protocol.AgentMetrics{AgentId: agentID},
		})
	} else {
		ama.allMetrics = append(ama.allMetrics, &metricsContainer{
			NoFindingCount: 1,
			AgentMetrics:   protocol.AgentMetrics{AgentId: agentID},
		})
	}
}

// CountResponse counts the response by checking the response status.
func (ama *AgentMetricsAggregator) CountResponse(agentID string, status protocol.ResponseStatus) {
	isError := status == protocol.ResponseStatus_ERROR

	for _, agentMetrics := range ama.allMetrics {
		if agentMetrics.AgentId != agentID {
			continue
		}
		agentMetrics.ResponseCount++
		if isError {
			agentMetrics.ErrorCount++
		}
		return
	}

	if isError {
		ama.allMetrics = append(ama.allMetrics, &metricsContainer{
			AgentMetrics: protocol.AgentMetrics{AgentId: agentID, ResponseCount: 1, ErrorCount: 1},
		})
	} else {
		ama.allMetrics = append(ama.allMetrics, &metricsContainer{
			AgentMetrics: protocol.AgentMetrics{AgentId: agentID, ResponseCount: 1},
		})
	}
}

// TryFlush checks the flushing condition(s) an returns metrics accordingly.
func (ama *AgentMetricsAggregator) TryFlush() []*protocol.AgentMetrics {
	now := time.Now()
	if now.Sub(ama.lastFlush) < ama.flushInterval {
		return nil
	}

	(allAgentMetrics)(ama.allMetrics).Fix()

	ama.lastFlush = now
	allContainers := ama.allMetrics
	ama.allMetrics = make([]*metricsContainer, 0)

	var allMetrics []*protocol.AgentMetrics
	for _, container := range allContainers {
		allMetrics = append(allMetrics, &container.AgentMetrics)
	}

	return allMetrics
}

// allAgentMetrics is an alias type for post-processing aggregated in-memory metrics
// before we publish them.
type allAgentMetrics []*metricsContainer

func (allMetrics allAgentMetrics) Fix() {
	allMetrics.CreateSummaries()
	allMetrics.CalculateFindingRates()
	allMetrics.CalculateAverages()
	allMetrics.FindMaxValues()
	allMetrics.CalculateP95()
	allMetrics.PutCounts()
	allMetrics.FindStartEndTimes()
}

func (allMetrics allAgentMetrics) CreateSummaries() {
	for _, agentMetrics := range allMetrics {
		if agentMetrics.TxProcessing != nil {
			agentMetrics.AgentMetrics.TxProcessing = &protocol.MetricSummary{}
		}
		if agentMetrics.BlockProcessing != nil {
			agentMetrics.AgentMetrics.BlockProcessing = &protocol.MetricSummary{}
		}
	}
}

func (allMetrics allAgentMetrics) CalculateAverages() {
	for _, agentMetrics := range allMetrics {
		if agentMetrics.TxProcessing != nil {
			agentMetrics.AgentMetrics.TxProcessing.Average = avgMetricArray(agentMetrics.TxProcessing)
		}
		if agentMetrics.BlockProcessing != nil {
			agentMetrics.AgentMetrics.BlockProcessing.Average = avgMetricArray(agentMetrics.BlockProcessing)
		}
	}
}

func avgMetricArray(data []*protocol.MetricData) float64 {
	sum := decimal.NewFromInt(0)
	for _, dataPoint := range data {
		sum = sum.Add(decimal.NewFromFloat(dataPoint.Value))
	}
	f, _ := sum.Div(decimal.NewFromInt32(int32(len(data)))).Round(2).Float64()
	return f
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

func (allMetrics allAgentMetrics) FindMaxValues() {
	for _, agentMetrics := range allMetrics {
		findMetricsMax(agentMetrics)
	}
}

func findMetricsMax(agentMetrics *metricsContainer) {
	if agentMetrics.TxProcessing != nil {
		agentMetrics.AgentMetrics.TxProcessing.Max = maxDataPoint(agentMetrics.TxProcessing)
	}
	if agentMetrics.BlockProcessing != nil {
		agentMetrics.AgentMetrics.BlockProcessing.Max = maxDataPoint(agentMetrics.BlockProcessing)
	}
}

func maxDataPoint(data []*protocol.MetricData) float64 {
	var max float64
	for _, dataPoint := range data {
		if dataPoint.Value > max {
			max = dataPoint.Value
		}
	}
	return max
}

func (allMetrics allAgentMetrics) CalculateP95() {
	for _, agentMetrics := range allMetrics {
		if agentMetrics.TxProcessing != nil {
			agentMetrics.AgentMetrics.TxProcessing.P95 = calcP95(agentMetrics.TxProcessing)
		}
		if agentMetrics.BlockProcessing != nil {
			agentMetrics.AgentMetrics.BlockProcessing.P95 = calcP95(agentMetrics.BlockProcessing)
		}
	}
}

func calcP95(data []*protocol.MetricData) float64 {
	switch len(data) {
	case 0:
		return 0
	case 1:
		return data[0].Value
	}

	k := len(data)
	k95, _ := decimal.NewFromInt32(int32(k)).Mul(decimal.NewFromFloat32(0.95)).Floor().BigFloat().Int64()
	sort.Slice(data, func(i, j int) bool {
		return data[i].Value < data[j].Value
	})
	return data[k95-1].Value
}

func (allMetrics allAgentMetrics) PutCounts() {
	for _, agentMetrics := range allMetrics {
		if agentMetrics.TxProcessing != nil {
			agentMetrics.AgentMetrics.TxProcessing.Count = int32(len(agentMetrics.TxProcessing))
		}
		if agentMetrics.BlockProcessing != nil {
			agentMetrics.AgentMetrics.BlockProcessing.Count = int32(len(agentMetrics.BlockProcessing))
		}
	}
}

func (allMetrics allAgentMetrics) FindStartEndTimes() {
	for _, agentMetrics := range allMetrics {
		if agentMetrics.TxProcessing != nil {
			agentMetrics.AgentMetrics.TxProcessing.StartTimestamp, agentMetrics.AgentMetrics.TxProcessing.EndTimestamp = findStartEndTs(agentMetrics.TxProcessing)
		}
		if agentMetrics.BlockProcessing != nil {
			agentMetrics.AgentMetrics.BlockProcessing.StartTimestamp, agentMetrics.AgentMetrics.BlockProcessing.EndTimestamp = findStartEndTs(agentMetrics.BlockProcessing)
		}
	}
}

func findStartEndTs(data []*protocol.MetricData) (startTs, endTs string) {
	sort.Slice(data, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, data[i].Timestamp)
		t2, _ := time.Parse(time.RFC3339, data[j].Timestamp)
		return t1.Before(t2)
	})
	startTs = data[0].Timestamp
	endTs = data[len(data)-1].Timestamp
	return
}
