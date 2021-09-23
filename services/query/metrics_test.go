package query_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/forta-network/forta-node/protocol"
	"github.com/forta-network/forta-node/services/query"
	"github.com/stretchr/testify/require"
)

const (
	testFlushIntervalSeconds = 1
	testAgentID1             = "agent-id-1"
	testAgentID2             = "agent-id-2"
)

var (
	testStartTimestamp = time.Now().Add(time.Hour).Format(time.RFC3339)
	testTimestamp      = time.Now().Add(time.Hour * 2).Format(time.RFC3339)
	testEndTimestamp   = time.Now().Add(time.Hour * 3).Format(time.RFC3339)
)

func TestAgentMetricsAggregator(t *testing.T) {
	r := require.New(t)

	aggregator := query.NewMetricsAggregator(testFlushIntervalSeconds)

	txProcessingData := []*protocol.MetricData{
		{
			Timestamp: testStartTimestamp,
			Value:     1,
		},
		{
			Timestamp: testTimestamp,
			Value:     10,
		},
		{
			Timestamp: testEndTimestamp,
			Value:     34,
		},
	}
	var (
		txProcessingAvg   float64 = 15 // 1 + 10 + 34 = 45 => 45 / 3 = 15
		txProcessingCount int32   = 3
		txProcessingMax   float64 = 34
		txProcessingP95   float64 = 10
	)

	blockProcessingData := []*protocol.MetricData{
		{
			Timestamp: testTimestamp,
			Value:     20,
		},
	}
	var (
		blockProcessingAvg   float64 = 20 // 20 / 1 = 20
		blockProcessingCount int32   = 1
		blockProcessingMax   float64 = 20
		blockProcessingP95   float64 = 20
	)

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

	b, _ := json.MarshalIndent(metrics, "", "  ")
	t.Log("flushed metrics:", string(b))

	metrics1 := metrics[0]
	r.Equal(testAgentID1, metrics1.AgentId)
	r.Equal(expectedFindingRatePct1, metrics1.FindingRatePct)

	r.Equal(txProcessingAvg, metrics1.TxProcessing.Average)
	r.Equal(txProcessingCount, metrics1.TxProcessing.Count)
	r.Equal(txProcessingMax, metrics1.TxProcessing.Max)
	r.Equal(txProcessingP95, metrics1.TxProcessing.P95)
	r.Equal(testStartTimestamp, metrics1.TxProcessing.StartTimestamp)
	r.Equal(testEndTimestamp, metrics1.TxProcessing.EndTimestamp)

	r.Equal(blockProcessingAvg, metrics1.BlockProcessing.Average)
	r.Equal(blockProcessingCount, metrics1.BlockProcessing.Count)
	r.Equal(blockProcessingMax, metrics1.BlockProcessing.Max)
	r.Equal(blockProcessingP95, metrics1.BlockProcessing.P95)
	r.Equal(testTimestamp, metrics1.BlockProcessing.StartTimestamp)
	r.Equal(testTimestamp, metrics1.BlockProcessing.EndTimestamp)

	metrics2 := metrics[1]
	r.Equal(testAgentID2, metrics2.AgentId)
	r.Equal(expectedFindingRatePct2, metrics2.FindingRatePct)

	r.Equal(txProcessingAvg, metrics2.TxProcessing.Average)
	r.Equal(txProcessingCount, metrics2.TxProcessing.Count)
	r.Equal(txProcessingMax, metrics2.TxProcessing.Max)
	r.Equal(txProcessingP95, metrics2.TxProcessing.P95)
	r.Equal(testStartTimestamp, metrics2.TxProcessing.StartTimestamp)
	r.Equal(testEndTimestamp, metrics2.TxProcessing.EndTimestamp)

	r.Equal(blockProcessingAvg, metrics2.BlockProcessing.Average)
	r.Equal(blockProcessingCount, metrics2.BlockProcessing.Count)
	r.Equal(blockProcessingMax, metrics2.BlockProcessing.Max)
	r.Equal(blockProcessingP95, metrics2.BlockProcessing.P95)
	r.Equal(testTimestamp, metrics2.BlockProcessing.StartTimestamp)
	r.Equal(testTimestamp, metrics2.BlockProcessing.EndTimestamp)
}
