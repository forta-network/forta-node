package query

import (
	"sort"
	"time"

	"github.com/forta-network/forta-node/protocol"
	"github.com/forta-network/forta-node/utils"
	"github.com/shopspring/decimal"
)

// Adjustable package settings
var (
	DefaultBucketInterval = time.Minute
)

// AgentMetricsAggregator aggregates agents' metrics and produces a list of summary of them when flushed.
type AgentMetricsAggregator struct {
	buckets   []*metricsBucket
	lastFlush time.Time
}

type metricsBucket struct {
	Time                   time.Time
	TxProcessingLatency    []uint32
	BlockProcessingLatency []uint32
	protocol.AgentMetrics
}

// NewAgentMetricsAggregator creates a new agent metrics aggregator.
func NewMetricsAggregator() *AgentMetricsAggregator {
	return &AgentMetricsAggregator{
		lastFlush: time.Now(), // avoid flushing immediately
	}
}

func (ama *AgentMetricsAggregator) findBucket(agentID string, t time.Time) *metricsBucket {
	bucketTime := FindClosestBucketTime(t)
	for _, bucket := range ama.buckets {
		if bucket.AgentId != agentID {
			continue
		}
		if !bucket.Time.Equal(bucketTime) {
			continue
		}
		return bucket
	}
	bucket := &metricsBucket{Time: bucketTime}
	bucket.AgentId = agentID
	bucket.Timestamp = utils.FormatTime(bucketTime)
	ama.buckets = append(ama.buckets, bucket)
	return bucket
}

// FindClosestBucketTime finds the closest bucket time. If it is per minute and the time is 15:15:15,
// then the closest is 15:15:00.
func FindClosestBucketTime(t time.Time) time.Time {
	ts := t.UnixNano()
	rem := ts % int64(DefaultBucketInterval)
	return time.Unix(0, ts-rem)
}

type agentResponse protocol.EvaluateTxResponse

// AggregateFromTxResponse aggregates metrics from a tx response.
func (ama *AgentMetricsAggregator) AggregateFromTxResponse(agentID string, resp *protocol.EvaluateTxResponse) {
	t, _ := time.Parse(time.RFC3339, resp.Timestamp)
	bucket := ama.findBucket(agentID, t)
	bucket.TxProcessingLatency = append(bucket.TxProcessingLatency, resp.LatencyMs)
	if bucket.TxProcessingLatencyMs == nil {
		bucket.TxProcessingLatencyMs = &protocol.MetricSummary{}
	}
	bucket.TxProcessingLatencyMs.Count++
	ama.aggregateFromAgentResponse(bucket, (*agentResponse)(resp))
}

// AggregateFromBlockResponse aggregates metrics from a block response.
func (ama *AgentMetricsAggregator) AggregateFromBlockResponse(agentID string, resp *protocol.EvaluateBlockResponse) {
	t, _ := time.Parse(time.RFC3339, resp.Timestamp)
	bucket := ama.findBucket(agentID, t)
	bucket.BlockProcessingLatency = append(bucket.BlockProcessingLatency, resp.LatencyMs)
	if bucket.BlockProcessingLatencyMs == nil {
		bucket.BlockProcessingLatencyMs = &protocol.MetricSummary{}
	}
	bucket.BlockProcessingLatencyMs.Count++
	ama.aggregateFromAgentResponse(bucket, (*agentResponse)(resp))
}

func (ama *AgentMetricsAggregator) aggregateFromAgentResponse(bucket *metricsBucket, resp *agentResponse) {
	bucket.ResponseCount++
	if resp.Status == protocol.ResponseStatus_ERROR {
		bucket.ErrorCount++
	}
	bucket.FindingCount += uint32(len(resp.Findings))
}

// TryFlush checks the flushing condition(s) an returns metrics accordingly.
func (ama *AgentMetricsAggregator) TryFlush() []*protocol.AgentMetrics {
	now := time.Now()
	if now.Sub(ama.lastFlush) < DefaultBucketInterval {
		return nil
	}

	ama.lastFlush = now
	buckets := ama.buckets
	ama.buckets = nil

	(allAgentMetrics)(buckets).Fix()

	var allMetrics []*protocol.AgentMetrics
	for _, bucket := range buckets {
		allMetrics = append(allMetrics, &bucket.AgentMetrics)
	}

	return allMetrics
}

// allAgentMetrics is an alias type for post-processing aggregated in-memory metrics
// before we publish them.
type allAgentMetrics []*metricsBucket

func (allMetrics allAgentMetrics) Fix() {
	sort.Slice(allMetrics, func(i, j int) bool {
		return allMetrics[i].Time.Before(allMetrics[j].Time)
	})
	allMetrics.CalculateAverages()
	allMetrics.FindMaxValues()
	allMetrics.CalculateP95()
}

func (allMetrics allAgentMetrics) CalculateAverages() {
	for _, agentMetrics := range allMetrics {
		if agentMetrics.TxProcessingLatency != nil {
			agentMetrics.AgentMetrics.TxProcessingLatencyMs.Average = avgMetricArray(agentMetrics.TxProcessingLatency)
		}
		if agentMetrics.BlockProcessingLatency != nil {
			agentMetrics.AgentMetrics.BlockProcessingLatencyMs.Average = avgMetricArray(agentMetrics.BlockProcessingLatency)
		}
	}
}

func avgMetricArray(data []uint32) float64 {
	sum := decimal.NewFromInt(0)
	for _, dataPoint := range data {
		sum = sum.Add(decimal.NewFromInt32(int32(dataPoint)))
	}
	f, _ := sum.Div(decimal.NewFromInt32(int32(len(data)))).Round(2).Float64()
	return f
}

func (allMetrics allAgentMetrics) FindMaxValues() {
	for _, agentMetrics := range allMetrics {
		findMetricsMax(agentMetrics)
	}
}

func findMetricsMax(agentMetrics *metricsBucket) {
	if agentMetrics.TxProcessingLatency != nil {
		agentMetrics.AgentMetrics.TxProcessingLatencyMs.Max = maxDataPoint(agentMetrics.TxProcessingLatency)
	}
	if agentMetrics.BlockProcessingLatency != nil {
		agentMetrics.AgentMetrics.BlockProcessingLatencyMs.Max = maxDataPoint(agentMetrics.BlockProcessingLatency)
	}
}

func maxDataPoint(data []uint32) float64 {
	var max float64
	for _, dataPoint := range data {
		if float64(dataPoint) > max {
			max = float64(dataPoint)
		}
	}
	return max
}

func (allMetrics allAgentMetrics) CalculateP95() {
	for _, agentMetrics := range allMetrics {
		if agentMetrics.TxProcessingLatency != nil {
			agentMetrics.AgentMetrics.TxProcessingLatencyMs.P95 = calcP95(agentMetrics.TxProcessingLatency)
		}
		if agentMetrics.BlockProcessingLatency != nil {
			agentMetrics.AgentMetrics.BlockProcessingLatencyMs.P95 = calcP95(agentMetrics.BlockProcessingLatency)
		}
	}
}

func calcP95(data []uint32) float64 {
	switch len(data) {
	case 0:
		return 0
	case 1:
		return float64(data[0])
	}

	k := len(data)
	k95, _ := decimal.NewFromInt32(int32(k)).Mul(decimal.NewFromFloat32(0.95)).Floor().BigFloat().Int64()
	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})
	return float64(data[k95-1])
}
