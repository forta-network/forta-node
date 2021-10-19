package query_test

import (
	"encoding/json"
	"testing"
	"time"

	mtx "github.com/forta-protocol/forta-node/metrics"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/services/query"
	"github.com/forta-protocol/forta-node/utils"
	"github.com/stretchr/testify/require"
)

const (
	testAgentID1 = "agent-id-1"
	testAgentID2 = "agent-id-2"
)

var (
	testNow         = time.Now()
	testTime1       time.Time
	testBucketTime1 time.Time
	testTime2       time.Time
	testBucketTime2 time.Time
	testTime3       time.Time
	testBucketTime3 time.Time
	testTime4       time.Time
	testBucketTime4 time.Time
)

func init() {
	query.DefaultBucketInterval = time.Second

	testTime1 = testNow
	testBucketTime1 = query.FindClosestBucketTime(testTime1)
	testTime2 = testTime1.Add(time.Nanosecond)
	testBucketTime2 = query.FindClosestBucketTime(testTime2)
	testTime3 = testTime2.Add(time.Nanosecond)
	testBucketTime3 = query.FindClosestBucketTime(testTime3)
	testTime4 = testNow.Add(query.DefaultBucketInterval)
	testBucketTime4 = query.FindClosestBucketTime(testTime4)
}

type Metrics protocol.AgentMetrics

func (metrics *Metrics) GetMetric(name string) *protocol.MetricSummary {
	for _, summary := range metrics.Metrics {
		if summary.Name == name {
			return summary
		}
	}
	return nil
}

func TestAgentMetricsAggregator(t *testing.T) {
	r := require.New(t)

	aggregator := query.NewMetricsAggregator()

	// Bucket 1: 3 findings, 1 error response
	// Bucket 2: No findings, no error responses
	txResponses := []*protocol.EvaluateTxResponse{
		{
			Status:    protocol.ResponseStatus_SUCCESS,
			Findings:  []*protocol.Finding{{}, {}},
			Timestamp: utils.FormatTime(testTime1),
			LatencyMs: 1,
		},
		{
			Status:    protocol.ResponseStatus_SUCCESS,
			Findings:  []*protocol.Finding{{}},
			Timestamp: utils.FormatTime(testTime2),
			LatencyMs: 10,
		},
		{
			Status:    protocol.ResponseStatus_ERROR,
			Findings:  []*protocol.Finding{},
			Timestamp: utils.FormatTime(testTime3),
			LatencyMs: 34,
		},
		{
			Status:    protocol.ResponseStatus_SUCCESS,
			Findings:  []*protocol.Finding{},
			Timestamp: utils.FormatTime(testTime4),
			LatencyMs: 45,
		},
	}
	var (
		txProcessingAvg   float64 = 15 // 1 + 10 + 34 = 45 => 45 / 3 = 15
		txProcessingCount uint32  = 3
		txProcessingMax   float64 = 34
		txProcessingP95   float64 = 10
		txProcessingSum   float64 = 45

		txProcessing2Avg   float64 = 45
		txProcessing2Count uint32  = 1
		txProcessing2Max   float64 = 45
		txProcessing2P95   float64 = 45
		txProcessing2Sum   float64 = 45
	)

	// Bucket 2: 2 findings, no error responses
	blockResponses := []*protocol.EvaluateBlockResponse{
		{
			Status:    protocol.ResponseStatus_SUCCESS,
			Findings:  []*protocol.Finding{{}, {}},
			Timestamp: utils.FormatTime(testTime4),
			LatencyMs: 20,
		},
	}
	var (
		blockProcessingAvg   float64 = 20 // 20 / 1 = 20
		blockProcessingCount uint32  = 1
		blockProcessingMax   float64 = 20
		blockProcessingP95   float64 = 20
		blockProcessingSum   float64 = 20
	)

	// Agent 1: First the blocks and then the txs
	for _, blockResp := range blockResponses {
		aggregator.AggregateFromBlockResponse(testAgentID1, blockResp)
	}
	for _, txResp := range txResponses {
		aggregator.AggregateFromTxResponse(testAgentID1, txResp)
	}

	// Agent 2: First the txs and then the blocks
	for _, txResp := range txResponses {
		aggregator.AggregateFromTxResponse(testAgentID2, txResp)
	}
	for _, blockResp := range blockResponses {
		aggregator.AggregateFromBlockResponse(testAgentID2, blockResp)
	}

	_ = aggregator.AddAgentMetric(&protocol.AgentMetric{
		AgentId:   testAgentID2,
		Timestamp: utils.FormatTime(testTime1),
		Name:      mtx.MetricStop,
		Value:     1,
	})

	// Ensure that we have waited long enough until the flush interval.
	time.Sleep(query.DefaultBucketInterval * 2)

	metrics := aggregator.TryFlush()
	r.Len(metrics, 4) // 2 agents × 2 buckets
	for i, bucket := range metrics {
		if i+1 == len(metrics) {
			continue
		}
		// Current bucket's timestamp should not be after the next bucket's timestamp.
		r.False(utils.ParseTime(bucket.Timestamp).After(utils.ParseTime(metrics[i+1].Timestamp)))
	}

	b, _ := json.MarshalIndent(metrics, "", "  ")
	t.Log("flushed metrics:", string(b))

	metrics1 := (*Metrics)(metrics[0]) // Agent 1, bucket 1
	r.Equal(testAgentID1, metrics1.AgentId)
	r.Equal(uint32(1), metrics1.GetMetric(mtx.MetricTxError).Count)
	r.Equal(uint32(3), metrics1.GetMetric(mtx.MetricFinding).Count)
	r.Equal(uint32(3), metrics1.GetMetric(mtx.MetricTxRequest).Count)
	r.Equal(uint32(2), metrics1.GetMetric(mtx.MetricTxSuccess).Count)
	r.Equal(uint32(0), metrics1.GetMetric(mtx.MetricStop).Count)

	txLatencyMetric1 := metrics1.GetMetric(mtx.MetricTxLatency)
	r.Equal(txProcessingAvg, txLatencyMetric1.Average)
	r.Equal(txProcessingCount, txLatencyMetric1.Count)
	r.Equal(txProcessingMax, txLatencyMetric1.Max)
	r.Equal(txProcessingP95, txLatencyMetric1.P95)
	r.Equal(txProcessingSum, txLatencyMetric1.Sum)

	r.Equal(utils.FormatTime(testBucketTime1), metrics1.Timestamp)

	r.Nil(metrics1.GetMetric(mtx.MetricBlockLatency))
	r.Nil(metrics1.GetMetric(mtx.MetricBlockError))
	r.Nil(metrics1.GetMetric(mtx.MetricBlockSuccess))
	r.Nil(metrics1.GetMetric(mtx.MetricBlockRequest))

	metrics2 := (*Metrics)(metrics[1]) // Agent 2, bucket 1
	r.Equal(testAgentID2, metrics2.AgentId)
	r.Equal(utils.FormatTime(testBucketTime1), metrics2.Timestamp)

	r.Equal(uint32(1), metrics2.GetMetric(mtx.MetricTxError).Count)
	r.Equal(uint32(3), metrics2.GetMetric(mtx.MetricFinding).Count)
	r.Equal(uint32(3), metrics2.GetMetric(mtx.MetricTxRequest).Count)
	r.Equal(uint32(2), metrics2.GetMetric(mtx.MetricTxSuccess).Count)
	r.Equal(uint32(1), metrics2.GetMetric(mtx.MetricStop).Count)

	txLatencyMetric2 := metrics2.GetMetric(mtx.MetricTxLatency)
	r.Equal(txProcessingAvg, txLatencyMetric2.Average)
	r.Equal(txProcessingCount, txLatencyMetric2.Count)
	r.Equal(txProcessingMax, txLatencyMetric2.Max)
	r.Equal(txProcessingP95, txLatencyMetric2.P95)
	r.Equal(txProcessingSum, txLatencyMetric2.Sum)

	r.Nil(metrics2.GetMetric(mtx.MetricBlockLatency))
	r.Nil(metrics2.GetMetric(mtx.MetricBlockError))
	r.Nil(metrics2.GetMetric(mtx.MetricBlockSuccess))
	r.Nil(metrics2.GetMetric(mtx.MetricBlockRequest))

	metrics3 := (*Metrics)(metrics[2]) // Agent 1, bucket 2
	r.Equal(testAgentID1, metrics3.AgentId)
	r.Equal(utils.FormatTime(testBucketTime4), metrics3.Timestamp)

	r.Equal(uint32(2), metrics3.GetMetric(mtx.MetricFinding).Count)
	r.Equal(uint32(0), metrics3.GetMetric(mtx.MetricTxError).Count)
	r.Equal(uint32(1), metrics3.GetMetric(mtx.MetricTxRequest).Count)
	r.Equal(uint32(1), metrics3.GetMetric(mtx.MetricTxSuccess).Count)
	r.Equal(uint32(0), metrics3.GetMetric(mtx.MetricBlockError).Count)
	r.Equal(uint32(1), metrics3.GetMetric(mtx.MetricBlockRequest).Count)
	r.Equal(uint32(1), metrics3.GetMetric(mtx.MetricBlockSuccess).Count)

	txLatencyMetric3 := metrics3.GetMetric(mtx.MetricTxLatency)
	r.Equal(txProcessing2Avg, txLatencyMetric3.Average)
	r.Equal(txProcessing2Count, txLatencyMetric3.Count)
	r.Equal(txProcessing2Max, txLatencyMetric3.Max)
	r.Equal(txProcessing2P95, txLatencyMetric3.P95)
	r.Equal(txProcessing2Sum, txLatencyMetric3.Sum)

	blockLatencyMetric3 := metrics3.GetMetric(mtx.MetricBlockLatency)
	r.Equal(blockProcessingAvg, blockLatencyMetric3.Average)
	r.Equal(blockProcessingCount, blockLatencyMetric3.Count)
	r.Equal(blockProcessingMax, blockLatencyMetric3.Max)
	r.Equal(blockProcessingP95, blockLatencyMetric3.P95)
	r.Equal(blockProcessingSum, blockLatencyMetric3.Sum)

	metrics4 := (*Metrics)(metrics[3]) // Agent 1, bucket 2
	r.Equal(testAgentID2, metrics4.AgentId)
	r.Equal(utils.FormatTime(testBucketTime4), metrics4.Timestamp)

	r.Equal(uint32(2), metrics4.GetMetric(mtx.MetricFinding).Count)
	r.Equal(uint32(0), metrics4.GetMetric(mtx.MetricTxError).Count)
	r.Equal(uint32(1), metrics4.GetMetric(mtx.MetricTxRequest).Count)
	r.Equal(uint32(1), metrics4.GetMetric(mtx.MetricTxSuccess).Count)
	r.Equal(uint32(0), metrics4.GetMetric(mtx.MetricBlockError).Count)
	r.Equal(uint32(1), metrics4.GetMetric(mtx.MetricBlockRequest).Count)
	r.Equal(uint32(1), metrics4.GetMetric(mtx.MetricBlockSuccess).Count)

	txLatencyMetric4 := metrics4.GetMetric(mtx.MetricTxLatency)
	r.Equal(txProcessing2Avg, txLatencyMetric4.Average)
	r.Equal(txProcessing2Count, txLatencyMetric4.Count)
	r.Equal(txProcessing2Max, txLatencyMetric4.Max)
	r.Equal(txProcessing2P95, txLatencyMetric4.P95)
	r.Equal(txProcessing2Sum, txLatencyMetric4.Sum)

	blockLatencyMetric4 := metrics4.GetMetric(mtx.MetricBlockLatency)
	r.Equal(blockProcessingAvg, blockLatencyMetric4.Average)
	r.Equal(blockProcessingCount, blockLatencyMetric4.Count)
	r.Equal(blockProcessingMax, blockLatencyMetric4.Max)
	r.Equal(blockProcessingP95, blockLatencyMetric4.P95)
	r.Equal(blockProcessingSum, blockLatencyMetric4.Sum)
}
