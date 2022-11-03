package publisher

import (
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchData_AppendPrivateAlert_PerFinding(t *testing.T) {
	bd := BatchData{}
	alert := &protocol.SignedAlert{
		Alert: &protocol.Alert{Id: "alertId", Finding: &protocol.Finding{
			Private: true,
		}},
	}
	nr := &protocol.NotifyRequest{
		SignedAlert:    alert,
		EvalTxRequest:  &protocol.EvaluateTxRequest{},
		EvalTxResponse: &protocol.EvaluateTxResponse{},
		AgentInfo: &protocol.AgentInfo{
			Manifest: "agentInfo",
		},
	}

	assert.Len(t, bd.PrivateAlerts, 0)
	bd.AppendAlert(nr)
	assert.Len(t, bd.PrivateAlerts, 1)
	assert.Equal(t, nr.AgentInfo.Manifest, bd.PrivateAlerts[0].AgentManifest)
	assert.Len(t, bd.PrivateAlerts[0].Alerts, 1)
	assert.EqualValues(t, alert, bd.PrivateAlerts[0].Alerts[0])
}

func TestBatchData_AppendPrivateAlert_Tx(t *testing.T) {
	bd := BatchData{}
	alert := &protocol.SignedAlert{
		Alert: &protocol.Alert{Id: "alertId", Finding: &protocol.Finding{}},
	}
	nr := &protocol.NotifyRequest{
		SignedAlert:   alert,
		EvalTxRequest: &protocol.EvaluateTxRequest{},
		EvalTxResponse: &protocol.EvaluateTxResponse{
			Private: true,
		},
		AgentInfo: &protocol.AgentInfo{
			Manifest: "agentInfo",
		},
	}

	assert.Len(t, bd.PrivateAlerts, 0)
	bd.AppendAlert(nr)
	assert.Len(t, bd.PrivateAlerts, 1)
	assert.Equal(t, nr.AgentInfo.Manifest, bd.PrivateAlerts[0].AgentManifest)
	assert.Len(t, bd.PrivateAlerts[0].Alerts, 1)
	assert.EqualValues(t, alert, bd.PrivateAlerts[0].Alerts[0])
}

func TestBatchData_AppendPrivateAlert_Block(t *testing.T) {
	bd := BatchData{}
	alert := &protocol.SignedAlert{
		Alert: &protocol.Alert{Id: "alertId", Finding: &protocol.Finding{}},
	}
	nr := &protocol.NotifyRequest{
		SignedAlert:      alert,
		EvalBlockRequest: &protocol.EvaluateBlockRequest{},
		EvalBlockResponse: &protocol.EvaluateBlockResponse{
			Private: true,
		},
		AgentInfo: &protocol.AgentInfo{
			Manifest: "agentInfo",
		},
	}

	assert.Len(t, bd.PrivateAlerts, 0)
	bd.AppendAlert(nr)
	assert.Len(t, bd.PrivateAlerts, 1)
	assert.Equal(t, nr.AgentInfo.Manifest, bd.PrivateAlerts[0].AgentManifest)
	assert.Len(t, bd.PrivateAlerts[0].Alerts, 1)
	assert.EqualValues(t, alert, bd.PrivateAlerts[0].Alerts[0])
}

func TestBatchData_AppendPrivateAlert_Combination(t *testing.T) {
	bd := BatchData{}
	alert := &protocol.SignedAlert{
		Alert: &protocol.Alert{Id: "alertId", Finding: &protocol.Finding{}},
	}
	nr := &protocol.NotifyRequest{
		SignedAlert:      alert,
		EvalCombinationRequest: &protocol.EvaluateCombinationRequest{},
		EvalCombinationResponse: &protocol.EvaluateCombinationResponse{
			Private: true,
		},
		AgentInfo: &protocol.AgentInfo{
			Manifest: "agentInfo",
		},
	}

	assert.Len(t, bd.PrivateAlerts, 0)
	bd.AppendAlert(nr)
	assert.Len(t, bd.PrivateAlerts, 1)
	assert.Equal(t, nr.AgentInfo.Manifest, bd.PrivateAlerts[0].AgentManifest)
	assert.Len(t, bd.PrivateAlerts[0].Alerts, 1)
	assert.EqualValues(t, alert, bd.PrivateAlerts[0].Alerts[0])
}

func TestShouldSkipPublishing(t *testing.T) {
	veryRecently := time.Now().Add(-time.Second * 2)

	testCases := []struct {
		name                string
		publisher           *Publisher
		batch               *protocol.AlertBatch
		expectedMsgContains string
		expectedSkipValue   bool
	}{
		{
			name: "has alerts",
			publisher: &Publisher{
				lastBatchSendAttempt: veryRecently,
			},
			batch: &protocol.AlertBatch{
				AlertCount: 2,
			},
			expectedSkipValue: false,
		},
		{
			name: "has metrics and running bots",
			publisher: &Publisher{
				lastBatchSendAttempt: veryRecently,
				botConfigs:           []config.AgentConfig{{}},
			},
			batch: &protocol.AlertBatch{
				Metrics: []*protocol.AgentMetrics{{}},
			},
			expectedSkipValue: false,
		},
		{
			name: "no metrics, running bots, too early",
			publisher: &Publisher{
				lastBatchSendAttempt: veryRecently,
				botConfigs:           []config.AgentConfig{{}},
			},
			batch:               &protocol.AlertBatch{},
			expectedSkipValue:   true,
			expectedMsgContains: "fast report deadline",
		},
		{
			name: "no metrics, running bots, not early",
			publisher: &Publisher{
				lastBatchSendAttempt: time.Now().Add(-fastReportInterval),
				botConfigs:           []config.AgentConfig{{}},
			},
			batch:             &protocol.AlertBatch{},
			expectedSkipValue: false,
		},
		{
			name: "no bots, too early",
			publisher: &Publisher{
				lastBatchSendAttempt: veryRecently,
			},
			batch:               &protocol.AlertBatch{},
			expectedSkipValue:   true,
			expectedMsgContains: "slow report deadline",
		},
		{
			name: "no bots, not early",
			publisher: &Publisher{
				lastBatchSendAttempt: time.Now().Add(-slowReportInterval), // not recently
			},
			batch:             &protocol.AlertBatch{},
			expectedSkipValue: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := require.New(t)

			msg, skip := testCase.publisher.shouldSkipPublishing(testCase.batch)
			r.Equal(testCase.expectedSkipValue, skip)
			r.Contains(msg, testCase.expectedMsgContains)
		})
	}
}
