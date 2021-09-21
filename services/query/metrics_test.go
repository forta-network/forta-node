package query_test

import (
	"testing"
	"time"

	"github.com/forta-network/forta-node/protocol"
	"github.com/forta-network/forta-node/services/query"
	"github.com/stretchr/testify/require"
)

const (
	testFlushIntervalSeconds = 1
	testThresholdMs          = 10
	testAgentID1             = "agent-id-1"
	testAgentID2             = "agent-id-2"
)

var (
	testTimestamp = time.Now().Format(time.RFC3339)
)

func TestAgentMetricsAggregator(t *testing.T) {
	r := require.New(t)

	aggregator := query.NewMetricsAggregator(testFlushIntervalSeconds, testThresholdMs)

	txProcessingData := []*protocol.MetricData{
		{
			Timestamp: testTimestamp,
			Value:     1,
		},
		{
			Timestamp: testTimestamp,
			Value:     10,
		},
		{
			Timestamp: testTimestamp,
			Value:     34,
		},
	}
	var txProcessingAvg int64 = 15 // 1 + 10 + 34 = 45 => 45 / 3 = 15

	blockProcessingData := []*protocol.MetricData{
		{
			Timestamp: testTimestamp,
			Value:     20,
		},
	}
	var blockProcessingAvg int64 = 20 // 20 / 1 = 20

	// Agent 1: 1 out of 3 alert notifs have a finding
	var expectedFindingRatePct1 float32 = 33.33
	aggregator.CountFinding(testAgentID1, false)
	aggregator.CountFinding(testAgentID1, true)
	aggregator.CountFinding(testAgentID1, false)

	// Agent 1: First the blocks and then the txs
	for _, blockData := range blockProcessingData {
		aggregator.PutBlockProcessingData(testAgentID1, blockData)
	}
	for _, txData := range txProcessingData {
		aggregator.PutTxProcessingData(testAgentID1, txData)
	}

	// Agent 2: First the txs and then the blocks
	for _, txData := range txProcessingData {
		aggregator.PutTxProcessingData(testAgentID2, txData)
	}
	for _, blockData := range blockProcessingData {
		aggregator.PutBlockProcessingData(testAgentID2, blockData)
	}

	// Agent 2: 1 out of 4 alert notifs have a finding
	var expectedFindingRatePct2 float32 = 25
	aggregator.CountFinding(testAgentID2, false)
	aggregator.CountFinding(testAgentID2, true)
	aggregator.CountFinding(testAgentID2, false)
	aggregator.CountFinding(testAgentID2, false)

	// Ensure that we have waited long enough until the flush interval.
	time.Sleep(testFlushIntervalSeconds * time.Second)

	metrics := aggregator.TryFlush()
	r.Len(metrics, 2)

	metrics1 := metrics[0]
	r.Equal(testAgentID1, metrics1.AgentId)
	r.Equal(int32(testThresholdMs), metrics1.ThresholdMs)
	r.Equal(expectedFindingRatePct1, metrics1.FindingRatePct)

	r.Equal(txProcessingAvg, metrics1.TxProcessing.Average)
	r.Len(metrics1.TxProcessing.Data, 2)
	r.Equal(float64(10), metrics1.TxProcessing.Data[0].Value)
	r.Equal(float64(34), metrics1.TxProcessing.Data[1].Value)

	r.Equal(blockProcessingAvg, metrics1.BlockProcessing.Average)
	r.Len(metrics1.BlockProcessing.Data, 1)
	r.Equal(float64(20), metrics1.BlockProcessing.Data[0].Value)

	metrics2 := metrics[1]
	r.Equal(testAgentID2, metrics2.AgentId)
	r.Equal(int32(testThresholdMs), metrics2.ThresholdMs)
	r.Equal(expectedFindingRatePct2, metrics2.FindingRatePct)

	r.Equal(txProcessingAvg, metrics2.TxProcessing.Average)
	r.Len(metrics2.TxProcessing.Data, 2)
	r.Equal(float64(10), metrics2.TxProcessing.Data[0].Value)
	r.Equal(float64(34), metrics2.TxProcessing.Data[1].Value)

	r.Equal(blockProcessingAvg, metrics2.BlockProcessing.Average)
	r.Len(metrics2.BlockProcessing.Data, 1)
	r.Equal(float64(20), metrics2.BlockProcessing.Data[0].Value)
}
