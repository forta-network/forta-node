package query_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/forta-network/forta-node/protocol"
	"github.com/forta-network/forta-node/services/query"
	"github.com/forta-network/forta-node/utils"
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

		txProcessing2Avg   float64 = 45
		txProcessing2Count uint32  = 1
		txProcessing2Max   float64 = 45
		txProcessing2P95   float64 = 45
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

	metrics1 := metrics[0] // Agent 1, bucket 1
	r.Equal(testAgentID1, metrics1.AgentId)
	r.Equal(uint32(1), metrics1.ErrorCount)
	r.Equal(uint32(3), metrics1.ResponseCount)
	r.Equal(uint32(3), metrics1.FindingCount)

	r.Equal(txProcessingAvg, metrics1.TxProcessingLatencyMs.Average)
	r.Equal(txProcessingCount, metrics1.TxProcessingLatencyMs.Count)
	r.Equal(txProcessingMax, metrics1.TxProcessingLatencyMs.Max)
	r.Equal(txProcessingP95, metrics1.TxProcessingLatencyMs.P95)
	r.Equal(utils.FormatTime(testBucketTime1), metrics1.Timestamp)

	r.Nil(metrics1.BlockProcessingLatencyMs)

	metrics2 := metrics[1] // Agent 2, bucket 1
	r.Equal(testAgentID2, metrics2.AgentId)
	r.Equal(uint32(1), metrics2.ErrorCount)
	r.Equal(uint32(3), metrics2.ResponseCount)
	r.Equal(uint32(3), metrics2.FindingCount)

	r.Equal(txProcessingAvg, metrics2.TxProcessingLatencyMs.Average)
	r.Equal(txProcessingCount, metrics2.TxProcessingLatencyMs.Count)
	r.Equal(txProcessingMax, metrics2.TxProcessingLatencyMs.Max)
	r.Equal(txProcessingP95, metrics2.TxProcessingLatencyMs.P95)
	r.Equal(utils.FormatTime(testBucketTime1), metrics2.Timestamp)

	r.Nil(metrics2.BlockProcessingLatencyMs)

	metrics3 := metrics[2] // Agent 1, bucket 2
	r.Equal(testAgentID1, metrics3.AgentId)
	r.Equal(uint32(0), metrics3.ErrorCount)
	r.Equal(uint32(2), metrics3.ResponseCount)
	r.Equal(uint32(2), metrics3.FindingCount)

	r.Equal(txProcessing2Avg, metrics3.TxProcessingLatencyMs.Average)
	r.Equal(txProcessing2Count, metrics3.TxProcessingLatencyMs.Count)
	r.Equal(txProcessing2Max, metrics3.TxProcessingLatencyMs.Max)
	r.Equal(txProcessing2P95, metrics3.TxProcessingLatencyMs.P95)
	r.Equal(utils.FormatTime(testBucketTime4), metrics3.Timestamp)

	r.Equal(blockProcessingAvg, metrics3.BlockProcessingLatencyMs.Average)
	r.Equal(blockProcessingCount, metrics3.BlockProcessingLatencyMs.Count)
	r.Equal(blockProcessingMax, metrics3.BlockProcessingLatencyMs.Max)
	r.Equal(blockProcessingP95, metrics3.BlockProcessingLatencyMs.P95)
	r.Equal(utils.FormatTime(testBucketTime4), metrics3.Timestamp)

	metrics4 := metrics[3] // Agent 2, bucket 2
	r.Equal(testAgentID2, metrics4.AgentId)
	r.Equal(uint32(0), metrics4.ErrorCount)
	r.Equal(uint32(2), metrics4.ResponseCount)
	r.Equal(uint32(2), metrics4.FindingCount)

	r.Equal(txProcessing2Avg, metrics4.TxProcessingLatencyMs.Average)
	r.Equal(txProcessing2Count, metrics4.TxProcessingLatencyMs.Count)
	r.Equal(txProcessing2Max, metrics4.TxProcessingLatencyMs.Max)
	r.Equal(txProcessing2P95, metrics4.TxProcessingLatencyMs.P95)
	r.Equal(utils.FormatTime(testBucketTime4), metrics4.Timestamp)

	r.Equal(blockProcessingAvg, metrics4.BlockProcessingLatencyMs.Average)
	r.Equal(blockProcessingCount, metrics4.BlockProcessingLatencyMs.Count)
	r.Equal(blockProcessingMax, metrics4.BlockProcessingLatencyMs.Max)
	r.Equal(blockProcessingP95, metrics4.BlockProcessingLatencyMs.P95)
	r.Equal(utils.FormatTime(testBucketTime4), metrics4.Timestamp)
}
