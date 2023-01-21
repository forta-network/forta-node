package agentpool

import (
	"context"
	"testing"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	"github.com/forta-network/forta-node/clients/messaging"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/scanner"
	"google.golang.org/grpc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testAgentID           = "test-agent"
	testRequestID         = "test-request-id"
	testResponseID        = "test-response-id"
	testCombinerSourceBot = "0xe5d1a9f6da9140c7762a2cd470e0f150db57550f3d51c5e56cba1941c1e1fa3e"
)

// TestSuite runs the test suite.
func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

// Suite is a test suite to test the agent pool.
type Suite struct {
	r *require.Assertions

	msgClient   *mock_clients.MockMessageClient
	agentClient *mock_clients.MockAgentClient

	ap *AgentPool

	suite.Suite
}

// SetupTest sets up the test.
func (s *Suite) SetupTest() {
	s.r = require.New(s.T())
	s.msgClient = mock_clients.NewMockMessageClient(gomock.NewController(s.T()))
	s.agentClient = mock_clients.NewMockAgentClient(gomock.NewController(s.T()))
	s.ap = &AgentPool{
		ctx:                     context.Background(),
		txResults:               make(chan *scanner.TxResult),
		blockResults:            make(chan *scanner.BlockResult),
		combinationAlertResults: make(chan *scanner.CombinationAlertResult),
		msgClient:               s.msgClient,
		dialer: func(agentCfg config.AgentConfig) (clients.AgentClient, error) {
			return s.agentClient, nil
		},
	}
}

// TestStartProcessStop tests the starting, processing and stopping flow for an agent.
func (s *Suite) TestStartProcessStop() {
	alertCfg := &protocol.AlertConfig{
		Subscriptions: []*protocol.CombinerBotSubscription{
			{
				BotId: testCombinerSourceBot,
			},
		},
	}
	agentConfig := config.AgentConfig{
		ID:          testAgentID,
		AlertConfig: alertCfg,
	}
	initResponse := &protocol.InitializeResponse{AlertConfig: alertCfg}

	agentPayload := messaging.AgentPayload{
		agentConfig,
	}
	emptyPayload := messaging.AgentPayload{}

	s.agentClient.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(initResponse, nil).AnyTimes()

	// Given that there are no agents running
	// When the latest list is received,
	// Then a "run" action should be published
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsStatusAttached, gomock.Any())
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsActionRun, gomock.Any())
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsAlertSubscribe,gomock.Any())
	s.r.NoError(s.ap.handleAgentVersionsUpdate(agentPayload))

	// Given that the agent is known to the pool but it is not ready yet
	s.r.Equal(1, len(s.ap.agents))
	s.r.False(s.ap.agents[0].IsReady())
	// When the agent pool receives a message saying that the agent started to run
	s.r.NoError(s.ap.handleStatusRunning(agentPayload))
	// Then the agent must be marked ready
	s.r.True(s.ap.agents[0].IsReady())

	// Given that the agent is running
	// When an evaluate requests are received
	// Then the agent should process them

	txReq := &protocol.EvaluateTxRequest{
		Event: &protocol.TransactionEvent{
			Block: &protocol.TransactionEvent_EthBlock{BlockNumber: "123123"},
			Transaction: &protocol.TransactionEvent_EthTransaction{
				Hash: "0x0",
			},
		},
	}
	txResp := &protocol.EvaluateTxResponse{Metadata: map[string]string{"imageHash": ""}}
	blockReq := &protocol.EvaluateBlockRequest{Event: &protocol.BlockEvent{BlockNumber: "123123"}}
	blockResp := &protocol.EvaluateBlockResponse{Metadata: map[string]string{"imageHash": ""}}
	combinerReq := &protocol.EvaluateAlertRequest{
		Event: &protocol.AlertEvent{
			Alert: &protocol.AlertEvent_Alert{
				Hash:   "123123",
				Source: &protocol.AlertEvent_Alert_Source{Bot: &protocol.AlertEvent_Alert_Bot{Id: testCombinerSourceBot}},
			},
		},
	}
	// save combiner subscription
	combinerResp := &protocol.EvaluateAlertResponse{Metadata: map[string]string{"imageHash": ""}}

	// test tx handling
	s.agentClient.EXPECT().Invoke(
		gomock.Any(), agentgrpc.MethodEvaluateTx,
		gomock.AssignableToTypeOf(&grpc.PreparedMsg{}), gomock.AssignableToTypeOf(&protocol.EvaluateTxResponse{}),
	).Return(nil)
	s.ap.SendEvaluateTxRequest(txReq)
	txResult := <-s.ap.TxResults()
	txResp.Timestamp = txResult.Response.Timestamp // bypass - hard to match

	// test block handling
	s.agentClient.EXPECT().Invoke(
		gomock.Any(), agentgrpc.MethodEvaluateBlock,
		gomock.AssignableToTypeOf(&grpc.PreparedMsg{}), gomock.AssignableToTypeOf(&protocol.EvaluateBlockResponse{}),
	).Return(nil)
	s.msgClient.EXPECT().Publish(messaging.SubjectScannerBlock, gomock.Any())
	s.ap.SendEvaluateBlockRequest(blockReq)
	blockResult := <-s.ap.BlockResults()
	blockResp.Timestamp = blockResult.Response.Timestamp // bypass - hard to match

	// test combine alert handling
	s.agentClient.EXPECT().Invoke(
		gomock.Any(), agentgrpc.MethodEvaluateAlert,
		gomock.AssignableToTypeOf(&grpc.PreparedMsg{}), gomock.AssignableToTypeOf(&protocol.EvaluateAlertResponse{}),
	).Return(nil)
	s.msgClient.EXPECT().Publish(messaging.SubjectScannerAlert, gomock.Any())
	s.ap.SendEvaluateAlertRequest(combinerReq)
	alertResult := <-s.ap.CombinationAlertResults()
	combinerResp.Timestamp = alertResult.Response.Timestamp // bypass - hard to match

	s.r.Equal(txReq, txResult.Request)
	s.r.Equal(txResp, txResult.Response)
	s.r.Equal(blockReq, blockResult.Request)
	s.r.Equal(blockResp, blockResult.Response)
	s.r.Equal(combinerReq, alertResult.Request)
	s.r.Equal(combinerResp, alertResult.Response)

	// Given that the agent is running
	// When an empty agent list is received
	// Then a "stop" action should be published
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsActionStop, gomock.Any())
	// And the agent must be closed
	s.agentClient.EXPECT().Close()
	s.r.NoError(s.ap.handleAgentVersionsUpdate(emptyPayload))
}
