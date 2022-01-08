package encoding_test

import (
	"github.com/forta-protocol/forta-node/encoding"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testMap = map[string]string{
	"key1": "val1",
	"key2": "val2",
}
var testBatch = &protocol.AlertBatch{
	ChainId:     1,
	BlockStart:  2,
	BlockEnd:    3,
	AlertCount:  4,
	MaxSeverity: 5,
	Results: []*protocol.BlockResults{
		{
			Block: &protocol.Block{
				BlockHash:      "hash",
				BlockNumber:    12345,
				BlockTimestamp: "timestamp",
			},
			Results: []*protocol.AgentAlerts{
				{
					AgentManifest: "manifest",
					Alerts: []*protocol.SignedAlert{
						{
							Alert: &protocol.Alert{
								Id:   "alertID",
								Type: 0,
								Finding: &protocol.Finding{
									Protocol:    "ethereum",
									Severity:    1,
									Metadata:    testMap,
									Type:        1,
									AlertId:     "0xalertId",
									Name:        "name",
									Description: "description",
								},
								Timestamp: "timestamp",
								Metadata:  testMap,
								Agent: &protocol.AgentInfo{
									Image:     "image",
									ImageHash: "imageHash",
									Id:        "id",
									Manifest:  "manifest",
								},
								Tags: testMap,
								Scanner: &protocol.ScannerInfo{
									Address: "scanner",
								},
							},
							Signature: &protocol.Signature{
								Signer:    "signer",
								Signature: "signature",
							},
							ChainId:         "0x1",
							BlockNumber:     "0x2",
							PublishedWithTx: "0x3",
						},
					},
				},
			},
			Transactions: []*protocol.TransactionResults{
				{
					Transaction: &protocol.TransactionEvent{
						Type: 2,
						Transaction: &protocol.TransactionEvent_EthTransaction{
							Type:     "a",
							Nonce:    "b",
							GasPrice: "c",
							Gas:      "d",
							Value:    "e",
							Input:    "f",
							V:        "g",
							R:        "h",
							S:        "i",
							To:       "j",
							Hash:     "k",
							From:     "l",
						},
						Receipt: &protocol.TransactionEvent_EthReceipt{
							Root:              "a",
							Status:            "b",
							CumulativeGasUsed: "c",
							LogsBloom:         "d",
							Logs:              nil,
							TransactionHash:   "e",
							ContractAddress:   "f",
							GasUsed:           "g",
							BlockHash:         "h",
							BlockNumber:       "i",
							TransactionIndex:  "j",
						},
						Network: &protocol.TransactionEvent_Network{ChainId: "0x1"},
					},
					Results: nil,
				},
			},
		},
	},
	Agents: []*protocol.BatchAgent{
		{
			Info: &protocol.AgentInfo{
				Image:     "image",
				ImageHash: "imageHash",
				Id:        "id",
				Manifest:  "manifest",
			},
			Blocks:       nil,
			Transactions: nil,
		},
	},
	Metrics: []*protocol.AgentMetrics{
		{
			AgentId:   "agentID",
			Timestamp: "timestamp",
			Metrics: []*protocol.MetricSummary{
				{
					Name:    "name",
					Count:   1,
					Max:     2,
					Average: 3,
					Sum:     4,
					P95:     5,
				},
			},
		},
	},
	ScannerVersion: &protocol.ScannerVersion{
		Commit: "commit", Ipfs: "ipfs",
	},
	Parent: "parent",
}

func TestEncodeAndDecodeBatch(t *testing.T) {
	res, err := encoding.EncodeBatch(testBatch)
	assert.NoError(t, err)

	decoded, err := encoding.DecodeBatch(res)
	assert.NoError(t, err)

	assert.True(t, proto.Equal(decoded, testBatch))
}
