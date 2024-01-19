package publisher

import (
	"sort"
	"sync"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/shopspring/decimal"
)

// AgentMetricsAggregator aggregates agents' metrics and produces a list of summary of them when flushed.
type AgentMetricsAggregator struct {
	// int64 is a chain id
	bucketsByChainID map[int64][]*metricsBucket
	bucketInterval   time.Duration
	lastFlush        time.Time
	mu               sync.RWMutex
}

type metricsBucket struct {
	Time        time.Time
	MetricsData map[string]metricsData
	protocol.AgentMetrics
}

type metricsData struct {
	Counters []uint32
	Details  string
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
		mu:               sync.RWMutex{},
		bucketInterval:   bucketInterval,
		lastFlush:        time.Now(), // avoid flushing immediately
		bucketsByChainID: make(map[int64][]*metricsBucket),
	}
}

func (ama *AgentMetricsAggregator) findBucket(agentID string, chainID int64, t time.Time) *metricsBucket {
	bucketTime := ama.FindClosestBucketTime(t)
	buckets, ok := ama.bucketsByChainID[chainID]
	if !ok {
		ama.bucketsByChainID[chainID] = make([]*metricsBucket, 0)
	}

	for _, bucket := range buckets {
		if bucket.AgentId != agentID {
			continue
		}
		if !bucket.Time.Equal(bucketTime) {
			continue
		}
		return bucket
	}

	bucket := &metricsBucket{
		Time:        bucketTime,
		MetricsData: make(map[string]metricsData),
	}

	bucket.AgentId = agentID
	bucket.Timestamp = utils.FormatTime(bucketTime)
	ama.bucketsByChainID[chainID] = append(ama.bucketsByChainID[chainID], bucket)
	return bucket
}

// FindClosestBucketTime finds the closest bucket time. If it is per minute and the time is 15:15:15,
// then the closest is 15:15:00.
func (ama *AgentMetricsAggregator) FindClosestBucketTime(t time.Time) time.Time {
	ts := t.UnixNano()
	rem := ts % int64(ama.bucketInterval)
	return time.Unix(0, ts-rem)
}

func (ama *AgentMetricsAggregator) AddAgentMetrics(ms *protocol.AgentMetricList) error {
	ama.mu.Lock()
	defer ama.mu.Unlock()

	for _, m := range ms.Metrics {
		t, _ := time.Parse(time.RFC3339, m.Timestamp)
		bucket := ama.findBucket(m.AgentId, m.ChainId, t)
		bucket.MetricsData[m.Name] = metricsData{
			Counters: append(bucket.MetricsData[m.Name].Counters, uint32(m.Value)),
			Details:  m.Details,
		}
	}

	return nil
}

// ForceFlush flushes without asking questions
func (ama *AgentMetricsAggregator) ForceFlush() []*protocol.AgentMetrics {
	ama.mu.Lock()
	defer ama.mu.Unlock()

	now := time.Now()

	ama.lastFlush = now
	buckets := ama.bucketsByChainID
	ama.bucketsByChainID = make(map[int64][]*metricsBucket)

	(allAgentMetrics)(buckets).Fix()

	var allMetrics []*protocol.AgentMetrics
	for _, metricsBuckets := range buckets {
		for _, bucket := range metricsBuckets {
			allMetrics = append(allMetrics, &bucket.AgentMetrics)
		}
	}

	return allMetrics
}

// TryFlush checks the flushing condition(s) an returns metrics accordingly.
func (ama *AgentMetricsAggregator) TryFlush() ([]*protocol.AgentMetrics, bool) {
	ama.mu.Lock()
	defer ama.mu.Unlock()

	now := time.Now()
	if now.Sub(ama.lastFlush) < ama.bucketInterval {
		return nil, false
	}

	ama.lastFlush = now
	buckets := ama.bucketsByChainID
	ama.bucketsByChainID = make(map[int64][]*metricsBucket)

	(allAgentMetrics)(buckets).Fix()

	var allMetrics []*protocol.AgentMetrics
	for _, metricsBuckets := range buckets {
		for _, bucket := range metricsBuckets {
			allMetrics = append(allMetrics, &bucket.AgentMetrics)
		}
	}

	return allMetrics, true
}

// allAgentMetrics is an alias type for post-processing aggregated in-memory metrics
// before we publish them.
type allAgentMetrics map[int64][]*metricsBucket

func (allMetrics allAgentMetrics) Fix() {
	for _, metricsBuckets := range allMetrics {
		sort.Slice(metricsBuckets, func(i, j int) bool {
			return metricsBuckets[i].Time.Before(metricsBuckets[j].Time)
		})
	}
	allMetrics.PrepareMetrics()
}

func (allMetrics allAgentMetrics) PrepareMetrics() {
	for chainID, metricsBuckets := range allMetrics {
		for _, agentMetrics := range metricsBuckets {
			for metricName, data := range agentMetrics.MetricsData {
				if len(data.Counters) > 0 {
					summary := agentMetrics.CreateAndGetSummary(metricName)
					summary.Count = uint32(len(data.Counters))
					summary.Average = avgMetricArray(data.Counters)
					summary.Max = maxDataPoint(data.Counters)
					summary.P95 = calcP95(data.Counters)
					summary.Sum = sumNums(data.Counters)
					summary.Details = data.Details
					summary.ChainId = chainID
				}
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
