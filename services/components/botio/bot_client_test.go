package botio

import (
	"context"
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	mock_agentgrpc "github.com/forta-network/forta-node/clients/agentgrpc/mocks"
	"github.com/forta-network/forta-node/clients/messaging"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/botio/botreq"
	mock_metrics "github.com/forta-network/forta-node/services/components/metrics/mocks"
	"google.golang.org/grpc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testBotID             = "test-bot"
	testRequestID         = "test-request-id"
	testResponseID        = "test-response-id"
	testCombinerSourceBot = "0x1d646c4045189991fdfd24a66b192a294158b839a6ec121d740474bdacb3abcd"
)

// BotClientSuite is a test suite to test the agent pool.
type BotClientSuite struct {
	r *require.Assertions

	alertConfig *protocol.AlertConfig

	msgClient        *mock_clients.MockMessageClient
	botGrpc          *mock_agentgrpc.MockClient
	lifecycleMetrics *mock_metrics.MockLifecycle
	botDialer        *mock_agentgrpc.MockBotDialer
	resultChannels   botreq.SendReceiveChannels

	botClient *botClient

	suite.Suite
}

// TestBotClientSuite runs the test suite.
func TestBotClientSuite(t *testing.T) {
	suite.Run(t, &BotClientSuite{})
}

// SetupTest sets up the test.
func (s *BotClientSuite) SetupTest() {
	s.r = require.New(s.T())

	ctrl := gomock.NewController(s.T())
	s.msgClient = mock_clients.NewMockMessageClient(ctrl)
	s.botGrpc = mock_agentgrpc.NewMockClient(ctrl)
	s.lifecycleMetrics = mock_metrics.NewMockLifecycle(ctrl)
	s.botDialer = mock_agentgrpc.NewMockBotDialer(ctrl)
	s.resultChannels = botreq.MakeResultChannels()

	s.botDialer.EXPECT().DialBot(gomock.Any()).Return(s.botGrpc, nil).AnyTimes()

	s.alertConfig = &protocol.AlertConfig{
		Subscriptions: []*protocol.CombinerBotSubscription{
			{
				BotId: testCombinerSourceBot,
			},
		},
	}

	s.botClient = NewBotClient(context.Background(), config.AgentConfig{
		ID: testBotID,
	}, s.msgClient, s.lifecycleMetrics, s.botDialer, s.resultChannels.SendOnly())
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{
		AlertConfig: s.alertConfig,
	}, nil).AnyTimes()
}

// TestStartProcessStop tests the starting, processing and stopping flow for a bot.
func (s *BotClientSuite) TestStartProcessStop() {
	combinerSubscriptions := MakeCombinerBotSubscriptions(s.alertConfig.Subscriptions, s.botClient.Config())

	s.lifecycleMetrics.EXPECT().Start(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusAttached(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusInitialized(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().ActionSubscribe(combinerSubscriptions)
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsAlertSubscribe, combinerSubscriptions)

	s.botClient.StartProcessing()
	s.botClient.Initialize()

	<-s.botClient.Initialized()

	txReq := &protocol.EvaluateTxRequest{
		Event: &protocol.TransactionEvent{
			Block: &protocol.TransactionEvent_EthBlock{BlockNumber: "123123"},
			Transaction: &protocol.TransactionEvent_EthTransaction{
				Hash: "0x0",
			},
		},
	}
	encodedTxReq, err := agentgrpc.EncodeMessage(txReq)
	s.r.NoError(err)
	txResp := &protocol.EvaluateTxResponse{Metadata: map[string]string{"imageHash": ""}}

	blockReq := &protocol.EvaluateBlockRequest{Event: &protocol.BlockEvent{BlockNumber: "123123"}}
	encodedBlockReq, err := agentgrpc.EncodeMessage(blockReq)
	s.r.NoError(err)
	blockResp := &protocol.EvaluateBlockResponse{Metadata: map[string]string{"imageHash": ""}}

	combinerReq := &protocol.EvaluateAlertRequest{
		TargetBotId: testBotID,
		Event: &protocol.AlertEvent{
			Alert: &protocol.AlertEvent_Alert{
				Hash:      "123123",
				Source:    &protocol.AlertEvent_Alert_Source{Bot: &protocol.AlertEvent_Alert_Bot{Id: testCombinerSourceBot}},
				CreatedAt: time.Now().Format(time.RFC3339Nano),
			},
		},
	}
	encodedCombinerReq, err := agentgrpc.EncodeMessage(combinerReq)
	s.r.NoError(err)
	combinerResp := &protocol.EvaluateAlertResponse{Metadata: map[string]string{"imageHash": ""}}

	// test tx handling
	s.botGrpc.EXPECT().Invoke(
		gomock.Any(), agentgrpc.MethodEvaluateTx,
		gomock.AssignableToTypeOf(&grpc.PreparedMsg{}), gomock.AssignableToTypeOf(&protocol.EvaluateTxResponse{}),
	).Return(nil)
	s.botClient.TxRequestCh() <- &botreq.TxRequest{
		Encoded:  encodedTxReq,
		Original: txReq,
	}
	txResult := <-s.resultChannels.Tx
	txResp.Timestamp = txResult.Response.Timestamp // bypass - hard to match

	// test block handling
	s.botGrpc.EXPECT().Invoke(
		gomock.Any(), agentgrpc.MethodEvaluateBlock,
		gomock.AssignableToTypeOf(&grpc.PreparedMsg{}), gomock.AssignableToTypeOf(&protocol.EvaluateBlockResponse{}),
	).Return(nil)
	s.botClient.BlockRequestCh() <- &botreq.BlockRequest{
		Encoded:  encodedBlockReq,
		Original: blockReq,
	}
	blockResult := <-s.resultChannels.Block
	blockResp.Timestamp = blockResult.Response.Timestamp // bypass - hard to match

	// test combine alert handling
	s.botGrpc.EXPECT().Invoke(
		gomock.Any(), agentgrpc.MethodEvaluateAlert,
		gomock.AssignableToTypeOf(&grpc.PreparedMsg{}), gomock.AssignableToTypeOf(&protocol.EvaluateAlertResponse{}),
	).Return(nil)
	s.botClient.CombinationRequestCh() <- &botreq.CombinationRequest{
		Encoded:  encodedCombinerReq,
		Original: combinerReq,
	}
	alertResult := <-s.resultChannels.CombinationAlert
	combinerResp.Timestamp = alertResult.Response.Timestamp // bypass - hard to match

	s.r.Equal(txReq, txResult.Request)
	s.r.Equal(txResp, txResult.Response)
	s.r.Equal(blockReq, blockResult.Request)
	s.r.Equal(blockResp, blockResult.Response)
	s.r.Equal(combinerReq, alertResult.Request)
	s.r.Equal(combinerResp, alertResult.Response)

	s.botGrpc.EXPECT().Close().AnyTimes()
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsAlertUnsubscribe, combinerSubscriptions)
	s.lifecycleMetrics.EXPECT().ActionUnsubscribe(combinerSubscriptions)

	s.r.NoError(s.botClient.Close())
}
