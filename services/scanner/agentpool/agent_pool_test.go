package agentpool

import (
	"context"
	"github.com/forta-network/forta-node/services/scanner/agentpool/poolagent"
	"testing"
	"time"

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
	testCombinerSourceBot = "0x1d646c4045189991fdfd24a66b192a294158b839a6ec121d740474bdacb3as12"
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
		botChanges:              make(chan []*poolagent.Agent),
		msgClient:               s.msgClient,
		dialer: func(agentCfg config.AgentConfig) (clients.AgentClient, error) {
			return s.agentClient, nil
		},
	}
	go s.ap.applyBotChanges()
}

// TestStartProcessStop tests the starting, processing and stopping flow for an agent.
func (s *Suite) TestStartProcessStop() {
	agentConfig := config.AgentConfig{
		ID: testAgentID,
		AlertConfig: &protocol.AlertConfig{
			Subscriptions: []*protocol.CombinerBotSubscription{
				{
					BotId: testCombinerSourceBot,
				},
			},
		},
	}

	agentPayload := messaging.AgentPayload{
		agentConfig,
	}
	emptyPayload := messaging.AgentPayload{}

	// Prior to invoking initialize method, agent.start metric should be emitted.
	s.msgClient.EXPECT().PublishProto(messaging.SubjectMetricAgent, gomock.Any())
	s.agentClient.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	// Given that there are no agents running
	// When the latest list is received,
	// Then a "run" action should be published
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsStatusAttached, gomock.Any())
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsActionRun, gomock.Any())
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsAlertSubscribe, gomock.Any())
	s.msgClient.EXPECT().PublishProto(messaging.SubjectMetricAgent, gomock.Any())
	s.r.NoError(s.ap.handleAgentVersionsUpdate(agentPayload))

	// wait for length to be 1 (async via channel)
	start := time.Now()
	for time.Since(start) < (5*time.Second) && len(s.ap.agents) != 1 {
		time.Sleep(50 * time.Millisecond)
	}

	// Given that the agent is known to the pool but it is not ready yet
	s.r.Equal(1, len(s.ap.agents))
	s.r.False(s.ap.agents[0].IsReady())
	// When the agent pool receives a message saying that the agent started to run
	s.msgClient.EXPECT().PublishProto(messaging.SubjectMetricAgent, gomock.Any()).Times(2)
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
		TargetBotId: testAgentID,
		Event: &protocol.AlertEvent{
			Alert: &protocol.AlertEvent_Alert{
				Hash:      "123123",
				Source:    &protocol.AlertEvent_Alert_Source{Bot: &protocol.AlertEvent_Alert_Bot{Id: testCombinerSourceBot}},
				CreatedAt: time.Now().Format(time.RFC3339Nano),
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
