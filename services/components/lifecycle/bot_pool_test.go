package lifecycle

import (
	"context"
	"testing"

	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/botio"
	mock_botio "github.com/forta-network/forta-node/services/components/botio/mocks"
	mock_metrics "github.com/forta-network/forta-node/services/components/metrics/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// BotPoolTestSuite has unit tests for the bot client pool.
type BotPoolTestSuite struct {
	r *require.Assertions

	lifecycleMetrics *mock_metrics.MockLifecycle
	botClientFactory *mock_botio.MockBotClientFactory
	botClient1       *mock_botio.MockBotClient
	botClient2       *mock_botio.MockBotClient

	botPool *botPool

	suite.Suite
}

func TestBotPoolTestSuite(t *testing.T) {
	suite.Run(t, &BotPoolTestSuite{})
}

func (s *BotPoolTestSuite) SetupTest() {
	s.r = s.Require()
	botRemoveTimeout = 0

	ctrl := gomock.NewController(s.T())
	s.lifecycleMetrics = mock_metrics.NewMockLifecycle(ctrl)
	s.botClientFactory = mock_botio.NewMockBotClientFactory(ctrl)
	s.botClient1 = mock_botio.NewMockBotClient(ctrl)
	s.botClient2 = mock_botio.NewMockBotClient(ctrl)

	s.botPool = NewBotPool(context.Background(), s.lifecycleMetrics, s.botClientFactory, 0)
	s.botPool.waitInit = true
}

func (s *BotPoolTestSuite) TestAddUpdate() {
	assigned := []config.AgentConfig{
		{
			ID:    testBotID2,
			Image: testImageRef,
		},
	}
	updated := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
		{
			ID:    testBotID2,
			Image: testImageRef,
			ShardConfig: &config.ShardConfig{
				ShardID: 1,
			},
		},
	}

	s.botPool.botClients = []botio.BotClient{s.botClient2}

	s.botClient2.EXPECT().Config().Return(assigned[0]).Times(3)
	s.botClientFactory.EXPECT().NewBotClient(gomock.Any(), updated[0]).Return(s.botClient1)
	s.botClient1.EXPECT().Initialize()
	s.botClient1.EXPECT().StartProcessing()

	s.botClient1.EXPECT().Config().Return(updated[0]).Times(1)
	s.botClient2.EXPECT().Config().Return(assigned[0]).Times(1)
	s.botClient2.EXPECT().SetConfig(updated[1])
	//s.botClient2.EXPECT().Config().Return(updated[1]).Times(2)
	s.lifecycleMetrics.EXPECT().ActionUpdate(updated[1])

	s.botPool.UpdateBotsWithLatestConfigs(updated)

	s.r.Len(s.botPool.botClients, 2)
	s.r.Equal(s.botPool.botClients[0], s.botClient2)
	s.r.Equal(s.botPool.botClients[1], s.botClient1)
}

func (s *BotPoolTestSuite) TestRemove() {
	assigned := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
		{
			ID:    testBotID2,
			Image: testImageRef,
		},
	}
	removed := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}

	s.botPool.botClients = []botio.BotClient{s.botClient1, s.botClient2}
	s.botClient1.EXPECT().Config().Return(assigned[0]).AnyTimes()
	s.botClient2.EXPECT().Config().Return(assigned[1]).AnyTimes()
	s.botClient1.EXPECT().Close().AnyTimes()

	s.botPool.RemoveBotsWithConfigs(removed)

	s.r.Len(s.botPool.botClients, 1)
	s.r.Equal(s.botPool.botClients[0], s.botClient1)
}

func (s *BotPoolTestSuite) TestReinit() {
	assigned := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}

	s.botPool.botClients = []botio.BotClient{s.botClient1}
	s.botClient1.EXPECT().Config().Return(assigned[0]).AnyTimes()
	s.botClient1.EXPECT().Initialize()

	s.botPool.ReinitBotsWithConfigs(assigned)

	s.r.Len(s.botPool.botClients, 1)
	s.r.Equal(s.botPool.botClients[0], s.botClient1)
}

func (s *BotPoolTestSuite) TestWaitForAll() {
	latest := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}
	botPool := NewBotPool(context.Background(), s.lifecycleMetrics, s.botClientFactory, len(latest))
	botPool.waitInit = true

	s.botClientFactory.EXPECT().NewBotClient(gomock.Any(), latest[0]).Return(s.botClient1)
	s.botClient1.EXPECT().Config().Return(latest[0]).Times(1)
	s.botClient1.EXPECT().LogStatus().AnyTimes()
	s.botClient1.EXPECT().Initialize()
	s.botClient1.EXPECT().StartProcessing()

	botPool.UpdateBotsWithLatestConfigs(latest)
	botPool.WaitForAll() // should be non-blocking

	s.r.Len(botPool.botClients, 1)
	s.r.Equal(botPool.botClients[0], s.botClient1)
}
