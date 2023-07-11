package botio

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	mock_agentgrpc "github.com/forta-network/forta-node/clients/agentgrpc/mocks"
	"github.com/forta-network/forta-node/clients/messaging"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/botio/botreq"
	mock_metrics "github.com/forta-network/forta-node/services/components/metrics/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
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
	lg        *logrus.Entry
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
	s.lg = logrus.WithField("component", "bot-client")

	s.botClient = NewBotClient(
		context.Background(), config.AgentConfig{
			ID: testBotID,
		}, s.msgClient, s.lifecycleMetrics, s.botDialer, s.resultChannels.SendOnly(),
	)
}

// TestStartProcessStop tests the starting, processing and stopping flow for a bot.
func (s *BotClientSuite) TestStartProcessStop() {
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(
		&protocol.InitializeResponse{
			AlertConfig: s.alertConfig,
		}, nil,
	).AnyTimes()

	combinerSubscriptions := MakeCombinerBotSubscriptions(s.alertConfig.Subscriptions, s.botClient.Config())

	s.lifecycleMetrics.EXPECT().ClientDial(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusAttached(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusInitialized(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().ActionSubscribe(combinerSubscriptions)

	// test health checks
	s.botGrpc.EXPECT().DoHealthCheck(gomock.Any())
	s.lifecycleMetrics.EXPECT().HealthCheckAttempt(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().HealthCheckSuccess(s.botClient.configUnsafe)

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
	txResp := &protocol.EvaluateTxResponse{Metadata: map[string]string{"imageHash": ""}}

	blockReq := &protocol.EvaluateBlockRequest{Event: &protocol.BlockEvent{BlockNumber: "123123"}}
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
	combinerResp := &protocol.EvaluateAlertResponse{Metadata: map[string]string{"imageHash": ""}}

	// test tx handling
	s.botGrpc.EXPECT().Invoke(
		gomock.Any(), agentgrpc.MethodEvaluateTx,
		gomock.AssignableToTypeOf(&protocol.EvaluateTxRequest{}), gomock.AssignableToTypeOf(&protocol.EvaluateTxResponse{}),
	).Return(nil)
	s.botClient.TxRequestCh() <- &botreq.TxRequest{
		Original: txReq,
	}
	txResult := <-s.resultChannels.Tx
	txResp.Timestamp = txResult.Response.Timestamp // bypass - hard to match
	txResp.LatencyMs = txResult.Response.LatencyMs // bypass - hard to match

	// test block handling
	s.botGrpc.EXPECT().Invoke(
		gomock.Any(), agentgrpc.MethodEvaluateBlock,
		gomock.AssignableToTypeOf(&protocol.EvaluateBlockRequest{}), gomock.AssignableToTypeOf(&protocol.EvaluateBlockResponse{}),
	).Return(nil)
	s.botClient.BlockRequestCh() <- &botreq.BlockRequest{
		Original: blockReq,
	}
	blockResult := <-s.resultChannels.Block
	blockResp.Timestamp = blockResult.Response.Timestamp // bypass - hard to match
	blockResp.LatencyMs = blockResult.Response.LatencyMs // bypass - hard to match

	// test combine alert handling
	s.botGrpc.EXPECT().Invoke(
		gomock.Any(), agentgrpc.MethodEvaluateAlert,
		gomock.AssignableToTypeOf(&protocol.EvaluateAlertRequest{}), gomock.AssignableToTypeOf(&protocol.EvaluateAlertResponse{}),
	).Return(nil)
	s.botClient.CombinationRequestCh() <- &botreq.CombinationRequest{
		Original: combinerReq,
	}
	alertResult := <-s.resultChannels.CombinationAlert
	combinerResp.Timestamp = alertResult.Response.Timestamp // bypass - hard to match
	combinerResp.LatencyMs = alertResult.Response.LatencyMs // bypass - hard to match

	// test error while handling
	invokeErr := fmt.Errorf("failed to invoke")
	s.botGrpc.EXPECT().Invoke(
		gomock.Any(), agentgrpc.MethodEvaluateAlert,
		gomock.AssignableToTypeOf(&protocol.EvaluateAlertRequest{}), gomock.AssignableToTypeOf(&protocol.EvaluateAlertResponse{}),
	).Return(invokeErr)
	s.lifecycleMetrics.EXPECT().BotError("combiner.invoke", invokeErr, s.botClient.configUnsafe)
	s.botClient.CombinationRequestCh() <- &botreq.CombinationRequest{
		Original: combinerReq,
	}
	<-s.resultChannels.CombinationAlert

	s.r.Equal(txReq, txResult.Request)
	s.r.Equal(txResp, txResult.Response)
	s.r.Equal(blockReq, blockResult.Request)
	s.r.Equal(blockResp, blockResult.Response)
	s.r.Equal(combinerReq, alertResult.Request)
	s.r.Equal(combinerResp, alertResult.Response)

	s.botGrpc.EXPECT().Close().AnyTimes()
	s.lifecycleMetrics.EXPECT().ClientClose(s.botClient.configUnsafe).AnyTimes()
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsAlertUnsubscribe, combinerSubscriptions)
	s.lifecycleMetrics.EXPECT().ActionUnsubscribe(combinerSubscriptions)

	s.r.NoError(s.botClient.Close())
}

func (s *BotClientSuite) TestCombinerBotSubscriptions() {
	s.botClient.SetAlertConfig(s.alertConfig)
	s.Equal(
		[]domain.CombinerBotSubscription{
			{

				Subscription: s.botClient.alertConfigUnsafe.Subscriptions[0],
				Subscriber: &domain.Subscriber{
					BotID:    s.botClient.configUnsafe.ID,
					BotOwner: s.botClient.configUnsafe.Owner,
					BotImage: s.botClient.configUnsafe.Image,
				},
			},
		},
		s.botClient.CombinerBotSubscriptions(),
	)
}

func (s *BotClientSuite) TestInitialize_Success() {
	s.lifecycleMetrics.EXPECT().ClientDial(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusAttached(s.botClient.configUnsafe)
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{
		Status:      protocol.ResponseStatus_SUCCESS,
		AlertConfig: s.alertConfig,
	}, nil).Times(1)
	s.lifecycleMetrics.EXPECT().StatusInitialized(s.botClient.configUnsafe)
	subs := MakeCombinerBotSubscriptions(s.alertConfig.Subscriptions, s.botClient.Config())
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsAlertSubscribe, subs)
	s.lifecycleMetrics.EXPECT().ActionSubscribe(subs)

	s.botClient.Initialize()
}

func (s *BotClientSuite) TestInitialize_Error() {
	s.lifecycleMetrics.EXPECT().ClientDial(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusAttached(s.botClient.configUnsafe)
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error")).Times(1)
	s.lifecycleMetrics.EXPECT().FailureInitialize(gomock.Any(), s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().ClientClose(s.botClient.configUnsafe)
	s.botGrpc.EXPECT().Close()

	s.botClient.Initialize()
}

func (s *BotClientSuite) TestInitialize_ResponseError() {
	s.lifecycleMetrics.EXPECT().ClientDial(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusAttached(s.botClient.configUnsafe)
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{
		Status: protocol.ResponseStatus_ERROR,
	}, nil).Times(1)
	s.lifecycleMetrics.EXPECT().FailureInitializeResponse(gomock.Any(), s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().ClientClose(s.botClient.configUnsafe)
	s.botGrpc.EXPECT().Close()

	s.botClient.Initialize()
}

func (s *BotClientSuite) TestInitialize_ValidationError() {
	s.lifecycleMetrics.EXPECT().ClientDial(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusAttached(s.botClient.configUnsafe)
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	s.lifecycleMetrics.EXPECT().FailureInitializeValidate(gomock.Any(), s.botClient.configUnsafe)

	s.botClient.Initialize()
}

func (s *BotClientSuite) TestHealthCheck() {
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(
		&protocol.InitializeResponse{}, nil,
	).AnyTimes()

	s.lifecycleMetrics.EXPECT().ClientDial(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusAttached(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusInitialized(s.botClient.configUnsafe)

	s.botClient.Initialize()

	<-s.botClient.Initialized()

	ctx := context.Background()

	// Mock HealthCheckAttempt() call
	s.lifecycleMetrics.EXPECT().HealthCheckAttempt(gomock.Any())

	s.botGrpc.EXPECT().DoHealthCheck(ctx).Return(nil)

	// Mock HealthCheckSuccess() call
	s.lifecycleMetrics.EXPECT().HealthCheckSuccess(gomock.Any())

	// Execute the method
	result := s.botClient.doHealthCheck(ctx, s.lg)

	s.r.False(result, "Expected healthCheck to return false")
}

func (s *BotClientSuite) TestHealthCheck_WithError() {
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(
		&protocol.InitializeResponse{}, nil,
	).AnyTimes()

	s.lifecycleMetrics.EXPECT().ClientDial(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusAttached(s.botClient.configUnsafe)
	s.lifecycleMetrics.EXPECT().StatusInitialized(s.botClient.configUnsafe)

	s.botClient.Initialize()

	<-s.botClient.Initialized()

	ctx := context.Background()

	// Mock HealthCheckAttempt() call
	s.lifecycleMetrics.EXPECT().HealthCheckAttempt(gomock.Any())

	err := fmt.Errorf("health check error")
	// Use Do() to modify the request parameter
	s.botGrpc.EXPECT().DoHealthCheck(ctx).Return(err)

	// Mock HealthCheckError() call
	s.lifecycleMetrics.EXPECT().HealthCheckError(gomock.Any(), gomock.Any())

	// Execute the method
	result := s.botClient.doHealthCheck(ctx, s.lg)

	s.r.False(result, "Expected healthCheck to return false")
}
