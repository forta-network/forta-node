package query_test

import (
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/services/query"
	"github.com/forta-protocol/forta-node/utils"
	"github.com/stretchr/testify/assert"
	"time"

	"testing"
)

var (
	testNow = time.Now()
)

func init() {
	query.DefaultBucketInterval = time.Second
}

func TestAgentMetricsAggregator(t *testing.T) {
	aggregator := query.NewMetricsAggregator()

	testTime1 := testNow
	//testBucketTime1 := query.FindClosestBucketTime(testTime1)

	err := aggregator.AddAgentMetrics(&protocol.AgentMetricList{Metrics: []*protocol.AgentMetric{
		{
			AgentId:   "agentID",
			Timestamp: utils.FormatTime(testTime1),
			Name:      "test.metric",
			Value:     1,
		},
		{
			AgentId:   "agentID",
			Timestamp: utils.FormatTime(testTime1),
			Name:      "test.metric",
			Value:     2,
		},
		{
			AgentId:   "agentID",
			Timestamp: utils.FormatTime(testTime1),
			Name:      "test.metric",
			Value:     3,
		},
		{
			AgentId:   "agentID",
			Timestamp: utils.FormatTime(testTime1),
			Name:      "test.metric",
			Value:     4,
		},
		{
			AgentId:   "agentID",
			Timestamp: utils.FormatTime(testTime1),
			Name:      "test.metric",
			Value:     5,
		},
	}})
	assert.NoError(t, err)

	time.Sleep(query.DefaultBucketInterval * 2)

	res := aggregator.TryFlush()

	assert.Len(t, res, 1)
	assert.Len(t, res[0].Metrics, 1)

	assert.Equal(t, &protocol.MetricSummary{
		Name:    "test.metric",
		Count:   5,
		Max:     5,
		Average: 3,
		Sum:     15,
		P95:     4,
	}, res[0].Metrics[0])

}
