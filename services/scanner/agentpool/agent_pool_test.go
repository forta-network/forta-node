package agentpool

import (
	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/protocol"
	"testing"

	mock_clients "OpenZeppelin/fortify-node/clients/mocks"
	"OpenZeppelin/fortify-node/services/scanner"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testAgentName  = "test-agent"
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
		dialer: func(agentCfg config.AgentConfig) clients.AgentClient {
			return s.agentClient
		},
	}
}

// TestStartProcessStop tests the starting, processing and stopping flow for an agent.
func (s *Suite) TestStartProcessStop() {
	agentConfig := config.AgentConfig{
		Name: testAgentName,
	}
	agentPayload := messaging.AgentPayload{
		agentConfig,
	}
	emptyPayload := messaging.AgentPayload{}

	// Given that there are no agents running
	// When the latest list is received,
	// Then a "run" action should be published
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsActionRun, gomock.Any())
	s.r.NoError(s.ap.handleAgentVersionsUpdate(agentPayload))
	// And the processing state must be "paused".
	s.r.Equal(int64(1), processingState.paused)

	// Given that the agent is known to the pool but it is not ready yet
	s.r.Equal(1, len(s.ap.agents))
	s.r.False(s.ap.agents[0].ready)
	// When the agent pool receives a message saying that the agent started to run
	s.r.NoError(s.ap.handleStatusRunning(agentPayload))
	// Then the processing state must not be "paused"
	s.r.Equal(int64(0), processingState.paused)

	// Given that the agent is running and the processing state is not paused
	// When an evaluate requests are received
	// Then the agent should process them

	txReq := &protocol.EvaluateTxRequest{Event: &protocol.TransactionEvent{Block: &protocol.TransactionEvent_EthBlock{BlockNumber: "123123"}}}
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
