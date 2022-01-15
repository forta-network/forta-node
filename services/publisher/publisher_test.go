package publisher

import (
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBatchData_AppendPrivateAlert_Tx(t *testing.T) {
	bd := BatchData{}
	alert := &protocol.SignedAlert{
		Alert: &protocol.Alert{Id: "alertId"},
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
		Alert: &protocol.Alert{Id: "alertId"},
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
