package publisher

import (
	"sort"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/shopspring/decimal"
)

// AgentMetricsAggregator aggregates agents' metrics and produces a list of summary of them when flushed.
type AgentMetricsAggregator struct {
	buckets        []*metricsBucket
	bucketInterval time.Duration
	lastFlush      time.Time
}

type metricsBucket struct {
	Time           time.Time
	MetricCounters map[string][]uint32
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

// NewAgentMetricsAggregator creates a new agent metrics aggregator.
func NewMetricsAggregator(bucketInterval time.Duration) *AgentMetricsAggregator {
	return &AgentMetricsAggregator{
		bucketInterval: bucketInterval,
		lastFlush:      time.Now(), // avoid flushing immediately
	}
}

func (ama *AgentMetricsAggregator) findBucket(agentID string, t time.Time) *metricsBucket {
	bucketTime := ama.FindClosestBucketTime(t)
	for _, bucket := range ama.buckets {
		if bucket.AgentId != agentID {
			continue
		}
		if !bucket.Time.Equal(bucketTime) {
			continue
		}
		return bucket
	}
	bucket := &metricsBucket{
		Time:           bucketTime,
		MetricCounters: make(map[string][]uint32),
	}
	bucket.AgentId = agentID
	bucket.Timestamp = utils.FormatTime(bucketTime)
	ama.buckets = append(ama.buckets, bucket)
	return bucket
}

// FindClosestBucketTime finds the closest bucket time. If it is per minute and the time is 15:15:15,
// then the closest is 15:15:00.
func (ama *AgentMetricsAggregator) FindClosestBucketTime(t time.Time) time.Time {
	ts := t.UnixNano()
	rem := ts % int64(ama.bucketInterval)
	return time.Unix(0, ts-rem)
}

type agentResponse protocol.EvaluateTxResponse

func (ama *AgentMetricsAggregator) AddAgentMetrics(ms *protocol.AgentMetricList) error {
	for _, m := range ms.Metrics {
		t, _ := time.Parse(time.RFC3339, m.Timestamp)
		bucket := ama.findBucket(m.AgentId, t)
		bucket.MetricCounters[m.Name] = append(bucket.MetricCounters[m.Name], uint32(m.Value))
	}
	return nil
}

// ForceFlush flushes without asking questions
func (ama *AgentMetricsAggregator) ForceFlush() []*protocol.AgentMetrics {
	now := time.Now()

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

// TryFlush checks the flushing condition(s) an returns metrics accordingly.
func (ama *AgentMetricsAggregator) TryFlush() ([]*protocol.AgentMetrics, bool) {
	now := time.Now()
	if now.Sub(ama.lastFlush) < ama.bucketInterval {
		return nil, false
	}

	ama.lastFlush = now
	buckets := ama.buckets
	ama.buckets = nil

	(allAgentMetrics)(buckets).Fix()

	var allMetrics []*protocol.AgentMetrics
	for _, bucket := range buckets {
		allMetrics = append(allMetrics, &bucket.AgentMetrics)
	}

	return allMetrics, true
}

// allAgentMetrics is an alias type for post-processing aggregated in-memory metrics
// before we publish them.
type allAgentMetrics []*metricsBucket

func (allMetrics allAgentMetrics) Fix() {
	sort.Slice(allMetrics, func(i, j int) bool {
		return allMetrics[i].Time.Before(allMetrics[j].Time)
	})
	allMetrics.PrepareMetrics()
}

func (allMetrics allAgentMetrics) PrepareMetrics() {
	for _, agentMetrics := range allMetrics {
		for metricName, list := range agentMetrics.MetricCounters {
			if len(list) > 0 {
				summary := agentMetrics.CreateAndGetSummary(metricName)
				summary.Count = uint32(len(list))
				summary.Average = avgMetricArray(list)
				summary.Max = maxDataPoint(list)
				summary.P95 = calcP95(list)
				summary.Sum = sumNums(list)
			}
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
