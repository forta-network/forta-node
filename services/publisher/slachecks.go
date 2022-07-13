package publisher

import (
	"sort"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
)

// SLAChecksAggregator aggregates sla checks and produces a list of summary of them when flushed.
type SLAChecksAggregator struct {
	buckets        []*slaChecksBucket
	bucketInterval time.Duration
	lastFlush      time.Time
}

type slaChecksBucket struct {
	Time          time.Time
	CheckCounters map[string][]uint32
	protocol.SLAChecks
}

func (mb *slaChecksBucket) CreateAndGetSummary(name string) *protocol.MetricSummary {
	for _, summary := range mb.Metrics {
		if summary.Name == name {
			return summary
		}
	}
	summary := &protocol.MetricSummary{Name: name}
	mb.Metrics = append(mb.Metrics, summary)

	return summary
}

// NewSLAChecksAggregator creates a new sla check aggregator.
func NewSLAChecksAggregator(bucketInterval time.Duration) *SLAChecksAggregator {
	return &SLAChecksAggregator{
		bucketInterval: bucketInterval,
		lastFlush:      time.Now(), // avoid flushing immediately
	}
}

func (ama *SLAChecksAggregator) findBucket(t time.Time) *slaChecksBucket {
	bucketTime := ama.FindClosestBucketTime(t)
	for _, bucket := range ama.buckets {
		if !bucket.Time.Equal(bucketTime) {
			continue
		}
		return bucket
	}

	bucket := &slaChecksBucket{
		Time:          bucketTime,
		CheckCounters: make(map[string][]uint32),
	}
	bucket.Timestamp = utils.FormatTime(bucketTime)
	ama.buckets = append(ama.buckets, bucket)
	return bucket
}

// FindClosestBucketTime finds the closest bucket time. If it is per minute and the time is 15:15:15,
// then the closest is 15:15:00.
func (ama *SLAChecksAggregator) FindClosestBucketTime(t time.Time) time.Time {
	ts := t.UnixNano()
	rem := ts % int64(ama.bucketInterval)
	return time.Unix(0, ts-rem)
}

func (ama *SLAChecksAggregator) AddSLACheck(ms *protocol.SLACheckList) error {
	for _, m := range ms.Checks {
		t, _ := time.Parse(time.RFC3339, m.Timestamp)
		bucket := ama.findBucket(t)
		bucket.CheckCounters[m.Name] = append(bucket.CheckCounters[m.Name], uint32(m.Value))
	}
	return nil
}

// ForceFlush flushes without asking questions
func (ama *SLAChecksAggregator) ForceFlush() []*protocol.SLAChecks {
	now := time.Now()

	ama.lastFlush = now
	buckets := ama.buckets
	ama.buckets = nil

	(allSLAChecks)(buckets).Fix()

	var allSLAChecks []*protocol.SLAChecks
	for _, bucket := range buckets {
		allSLAChecks = append(allSLAChecks, &bucket.SLAChecks)
	}

	return allSLAChecks
}

// TryFlush checks the flushing condition(s) and returns sla checks accordingly.
func (ama *SLAChecksAggregator) TryFlush() ([]*protocol.SLAChecks, bool) {
	now := time.Now()
	if now.Sub(ama.lastFlush) < ama.bucketInterval {
		return nil, false
	}

	ama.lastFlush = now
	buckets := ama.buckets
	ama.buckets = nil

	(allSLAChecks)(buckets).Fix()

	var allSlaChecks []*protocol.SLAChecks
	for _, bucket := range buckets {
		allSlaChecks = append(allSlaChecks, &bucket.SLAChecks)
	}

	return allSlaChecks, true
}

// allSLAChecks is an alias type for post-processing aggregated in-memory sla checks
// before we publish them.
type allSLAChecks []*slaChecksBucket

func (allChecks allSLAChecks) Fix() {
	sort.Slice(
		allChecks, func(i, j int) bool {
			return allChecks[i].Time.Before(allChecks[j].Time)
		},
	)
	allChecks.PrepareSLAChecks()
}

func (allChecks allSLAChecks) PrepareSLAChecks() {
	for _, slaChecks := range allChecks {
		for metricName, list := range slaChecks.CheckCounters {
			if len(list) > 0 {
				summary := slaChecks.CreateAndGetSummary(metricName)
				summary.Count = uint32(len(list))
				summary.Average = avgMetricArray(list)
				summary.Max = maxDataPoint(list)
				summary.P95 = calcP95(list)
				summary.Sum = sumNums(list)
			}
		}
	}
}
