package agentpool

import (
	"testing"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	mock_clients "github.com/forta-protocol/forta-node/clients/mocks"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/services/scanner"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testAgentID    = "test-agent"
	testRequestID  = "test-request-id"
	testResponseID = "test-response-id"
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
		txResults:    make(chan *scanner.TxResult, DefaultBufferSize),
		blockResults: make(chan *scanner.BlockResult, DefaultBufferSize),
		msgClient:    s.msgClient,
		dialer: func(agentCfg config.AgentConfig) (clients.AgentClient, error) {
			return s.agentClient, nil
		},
		activeAgents: make(chan int, 1000),
	}
}

// TestStartProcessStop tests the starting, processing and stopping flow for an agent.
func (s *Suite) TestStartProcessStop() {
	agentConfig := config.AgentConfig{
		ID: testAgentID,
	}
	agentPayload := messaging.AgentPayload{
		agentConfig,
	}
	emptyPayload := messaging.AgentPayload{}

	// Given that there are no agents running
	// When the latest list is received,
	// Then a "run" action should be published
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsStatusAttached, gomock.Any())
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsActionRun, gomock.Any())
	s.msgClient.EXPECT().Publish(messaging.SubjectMetricAgent, gomock.Any())

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

	txReq := &protocol.EvaluateTxRequest{Event: &protocol.TransactionEvent{Block: &protocol.TransactionEvent_EthBlock{BlockNumber: "123123"}, Transaction: &protocol.TransactionEvent_EthTransaction{Hash: "0x0"}}}
	txResp := &protocol.EvaluateTxResponse{Metadata: map[string]string{}}
	blockReq := &protocol.EvaluateBlockRequest{Event: &protocol.BlockEvent{BlockNumber: "123123"}}
	blockResp := &protocol.EvaluateBlockResponse{Metadata: map[string]string{}}
	s.agentClient.EXPECT().EvaluateTx(gomock.Any(), txReq).Return(txResp, nil)
	s.agentClient.EXPECT().EvaluateBlock(gomock.Any(), blockReq).Return(blockResp, nil)
	s.ap.SendEvaluateTxRequest(txReq)
	s.ap.SendEvaluateBlockRequest(blockReq)
	txResult := <-s.ap.TxResults()
	blockResult := <-s.ap.BlockResults()
	s.r.Equal(txReq, txResult.Request)
	s.r.Equal(txResp, txResult.Response)
	s.r.Equal(blockReq, blockResult.Request)
	s.r.Equal(blockResp, blockResult.Response)

	// Given that the agent is running
	// When an empty agent list is received
	// Then a "stop" action should be published
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsActionStop, gomock.Any())
	// And the agent must be closed
	s.agentClient.EXPECT().Close()
	s.r.NoError(s.ap.handleAgentVersionsUpdate(emptyPayload))
}
