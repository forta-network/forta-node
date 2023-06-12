package botio_test

import (
	"context"
	"testing"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients/messaging"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/botio"
	"github.com/forta-network/forta-node/services/components/botio/botreq"
	mock_botio "github.com/forta-network/forta-node/services/components/botio/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SenderTestSuite struct {
	r *require.Assertions

	botPool   *mock_botio.MockBotPool
	botClient *mock_botio.MockBotClient
	msgClient *mock_clients.MockMessageClient

	sender botio.Sender

	suite.Suite
}

func TestSender(t *testing.T) {
	suite.Run(t, &SenderTestSuite{})
}

func (s *SenderTestSuite) SetupTest() {
	s.r = s.Require()

	ctrl := gomock.NewController(s.T())
	s.botPool = mock_botio.NewMockBotPool(ctrl)
	s.botClient = mock_botio.NewMockBotClient(ctrl)
	s.msgClient = mock_clients.NewMockMessageClient(ctrl)

	s.botPool.EXPECT().GetCurrentBotClients().Return([]botio.BotClient{s.botClient}).AnyTimes()

	s.sender = botio.NewSender(context.Background(), s.msgClient, s.botPool)
}

func (s *SenderTestSuite) TestHealth() {
	s.botClient.EXPECT().TxBufferIsFull().Return(false)
	reports := s.sender.Health()
	s.r.Equal("agents.total", reports[0].Name)
	s.r.Equal("agents.lagging", reports[1].Name)
}

func (s *SenderTestSuite) TestSendEvaluateTxRequest() {
	s.botPool.EXPECT().WaitForAll().Times(1)
	s.botClient.EXPECT().ShouldProcessBlock(gomock.Any()).Return(true)
	s.botClient.EXPECT().Config().Return(config.AgentConfig{})
	s.botClient.EXPECT().Closed().Return(make(chan struct{}))
	s.botClient.EXPECT().TxRequestCh().Return(make(chan *botreq.TxRequest, 1))

	s.sender.SendEvaluateTxRequest(&protocol.EvaluateTxRequest{
		Event: &protocol.TransactionEvent{
			Transaction: &protocol.TransactionEvent_EthTransaction{
				Hash: "0x1",
			},
			Block: &protocol.TransactionEvent_EthBlock{
				BlockNumber: "0x1",
			},
		},
	})
}

func (s *SenderTestSuite) TestSendEvaluateBlockRequest() {
	s.botPool.EXPECT().WaitForAll().Times(1)
	s.botClient.EXPECT().ShouldProcessBlock(gomock.Any()).Return(true)
	s.botClient.EXPECT().Config().Return(config.AgentConfig{})
	s.botClient.EXPECT().Closed().Return(make(chan struct{}))
	s.botClient.EXPECT().BlockRequestCh().Return(make(chan *botreq.BlockRequest, 1))
	s.msgClient.EXPECT().Publish(messaging.SubjectScannerBlock, gomock.Any())

	s.sender.SendEvaluateBlockRequest(&protocol.EvaluateBlockRequest{
		Event: &protocol.BlockEvent{
			Block: &protocol.BlockEvent_EthBlock{
				Number: "0x1",
			},
		},
	})
}

func (s *SenderTestSuite) TestSendEvaluateAlertRequest() {
	s.botPool.EXPECT().WaitForAll().Times(1)
	s.botClient.EXPECT().ShouldProcessAlert(gomock.Any()).Return(true)
	s.botClient.EXPECT().Config().Return(config.AgentConfig{}).Times(2)
	s.botClient.EXPECT().Closed().Return(make(chan struct{}))
	s.botClient.EXPECT().CombinationRequestCh().Return(make(chan *botreq.CombinationRequest, 1))
	s.msgClient.EXPECT().Publish(messaging.SubjectScannerAlert, gomock.Any())

	s.sender.SendEvaluateAlertRequest(&protocol.EvaluateAlertRequest{
		Event: &protocol.AlertEvent{
			Alert: &protocol.AlertEvent_Alert{
				Source: &protocol.AlertEvent_Alert_Source{
					Bot: &protocol.AlertEvent_Alert_Bot{},
				},
			},
		},
	})
}
