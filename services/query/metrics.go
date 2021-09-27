package query

import (
	"sort"
	"time"

	"github.com/forta-network/forta-node/protocol"
	"github.com/forta-network/forta-node/utils"
	"github.com/shopspring/decimal"
)

// Metric fields
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
	Time            time.Time
	FindingCount    uint32
	TxProcessing    agentResponseMetrics
	BlockProcessing agentResponseMetrics
	protocol.AgentMetrics
}

func (mb *metricsBucket) CreateAndGetSummary(name string) *protocol.MetricSummary {
	for _, summary := range mb.Metrics {
		if summary.Name == name {
			return summary
		}
	}
	summary := &protocol.MetricSummary{Name: name}
	mb.Metrics = append(mb.Metrics, summary)
	return summary
}

type agentResponseMetrics struct {
	Request uint32
	Success uint32
	Error   uint32
	Latency []uint32
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
	bucket.TxProcessing.Latency = append(bucket.TxProcessing.Latency, resp.LatencyMs)
	bucket.TxProcessing.Request++
	if resp.Status == protocol.ResponseStatus_ERROR {
		bucket.TxProcessing.Error++
	} else {
		bucket.TxProcessing.Success++
	}
	bucket.FindingCount += uint32(len(resp.Findings))
}

// AggregateFromBlockResponse aggregates metrics from a block response.
func (ama *AgentMetricsAggregator) AggregateFromBlockResponse(agentID string, resp *protocol.EvaluateBlockResponse) {
	t, _ := time.Parse(time.RFC3339, resp.Timestamp)
	bucket := ama.findBucket(agentID, t)
	bucket.BlockProcessing.Latency = append(bucket.BlockProcessing.Latency, resp.LatencyMs)
	bucket.BlockProcessing.Request++
	if resp.Status == protocol.ResponseStatus_ERROR {
		bucket.BlockProcessing.Error++
	} else {
		bucket.BlockProcessing.Success++
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
	allMetrics.PrepareLatencyMetrics()
	allMetrics.PrepareCountMetrics()
}

func (allMetrics allAgentMetrics) PrepareLatencyMetrics() {
	for _, agentMetrics := range allMetrics {
		if len(agentMetrics.TxProcessing.Latency) > 0 {
			latencyNums := agentMetrics.TxProcessing.Latency
			txLatency := agentMetrics.CreateAndGetSummary(MetricTxLatency)
			txLatency.Count = uint32(len(latencyNums))
			txLatency.Average = avgMetricArray(latencyNums)
			txLatency.Max = maxDataPoint(latencyNums)
			txLatency.P95 = calcP95(latencyNums)
			txLatency.Sum = sumNums(latencyNums)
		}
		if len(agentMetrics.BlockProcessing.Latency) > 0 {
			latencyNums := agentMetrics.BlockProcessing.Latency
			blockLatency := agentMetrics.CreateAndGetSummary(MetricBlockLatency)
			blockLatency.Count = uint32(len(latencyNums))
			blockLatency.Average = avgMetricArray(latencyNums)
			blockLatency.Max = maxDataPoint(latencyNums)
			blockLatency.P95 = calcP95(latencyNums)
			blockLatency.Sum = sumNums(latencyNums)
		}
	}
}

func (allMetrics allAgentMetrics) PrepareCountMetrics() {
	for _, agentMetrics := range allMetrics {
		finding := agentMetrics.CreateAndGetSummary(MetricFinding)
		setCountMetric(finding, agentMetrics.FindingCount)

		if len(agentMetrics.TxProcessing.Latency) > 0 {
			request := agentMetrics.CreateAndGetSummary(MetricTxRequest)
			setCountMetric(request, agentMetrics.TxProcessing.Request)
			success := agentMetrics.CreateAndGetSummary(MetricTxSuccess)
			setCountMetric(success, agentMetrics.TxProcessing.Success)
			errorM := agentMetrics.CreateAndGetSummary(MetricTxError)
			setCountMetric(errorM, agentMetrics.TxProcessing.Error)
		}

		if len(agentMetrics.BlockProcessing.Latency) > 0 {
			request := agentMetrics.CreateAndGetSummary(MetricBlockRequest)
			setCountMetric(request, agentMetrics.BlockProcessing.Request)
			success := agentMetrics.CreateAndGetSummary(MetricBlockSuccess)
			setCountMetric(success, agentMetrics.BlockProcessing.Success)
			errorM := agentMetrics.CreateAndGetSummary(MetricBlockError)
			setCountMetric(errorM, agentMetrics.BlockProcessing.Error)
		}
	}
}

func setCountMetric(summary *protocol.MetricSummary, count uint32) {
	summary.Count = count
	if count > 0 {
		summary.Average = 1
		summary.Max = 1
		summary.P95 = 1
		summary.Sum = float64(count)
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

func maxDataPoint(data []uint32) float64 {
	var max float64
	for _, dataPoint := range data {
		if float64(dataPoint) > max {
			max = float64(dataPoint)
		}
	}
	return max
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

func sumNums(data []uint32) (n float64) {
	for _, d := range data {
		n += float64(d)
	}
	return
}
