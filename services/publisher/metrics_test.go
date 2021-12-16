package publisher_test

import (
	"time"

	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/services/publisher"
	"github.com/forta-protocol/forta-node/utils"
	"github.com/stretchr/testify/assert"

	"testing"
)

var (
	testNow = time.Now()
)

func init() {
	publisher.DefaultBucketInterval = time.Millisecond
}

type MetricsMathTest struct {
	metrics  []float64
	expected *protocol.MetricSummary
}

func TestAgentMetricsAggregator_math(t *testing.T) {

	tests := []*MetricsMathTest{
		{
			metrics: []float64{1, 2, 3, 4, 5},
			expected: &protocol.MetricSummary{
				Name:    "test.metric",
				Count:   5,
				Max:     5,
				Average: 3,
				Sum:     15,
				P95:     4,
			},
		},
		{
			metrics: []float64{1, 10, 34},
			expected: &protocol.MetricSummary{
				Name:    "test.metric",
				Count:   3,
				Max:     34,
				Average: 15,
				Sum:     45,
				P95:     10,
			},
		},
		{
			metrics: []float64{45},
			expected: &protocol.MetricSummary{
				Name:    "test.metric",
				Count:   1,
				Max:     45,
				Average: 45,
				Sum:     45,
				P95:     45,
			},
		},
	}

	for _, test := range tests {
		testTime1 := testNow

		var metrics []*protocol.AgentMetric
		for _, val := range test.metrics {
			metrics = append(metrics, &protocol.AgentMetric{
				AgentId:   "agentID",
				Timestamp: utils.FormatTime(testTime1),
				Name:      "test.metric",
				Value:     val,
			})
		}

		aggregator := publisher.NewMetricsAggregator()
		err := aggregator.AddAgentMetrics(&protocol.AgentMetricList{Metrics: metrics})
		assert.NoError(t, err)
		time.Sleep(publisher.DefaultBucketInterval * 2)

		res := aggregator.TryFlush()

		assert.Len(t, res, 1)
		assert.Len(t, res[0].Metrics, 1)
		assert.Equal(t, res[0].Metrics[0], test.expected)
	}

}
