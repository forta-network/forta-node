package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	mock_agentgrpc "github.com/forta-network/forta-node/clients/agentgrpc/mocks"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	mock_containers "github.com/forta-network/forta-node/services/components/containers/mocks"
	mock_lifecycle "github.com/forta-network/forta-node/services/components/lifecycle/mocks"
	mock_metrics "github.com/forta-network/forta-node/services/components/metrics/mocks"
	mock_registry "github.com/forta-network/forta-node/services/components/registry/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// BotLifecycleManagerTestSuite has unit tests for the bot lifecycle manager.
type BotLifecycleManagerTestSuite struct {
	r *require.Assertions

	msgClient        *mock_clients.MockMessageClient
	lifecycleMetrics *mock_metrics.MockLifecycle
	botGrpc          *mock_agentgrpc.MockClient
	botRegistry      *mock_registry.MockBotRegistry
	botContainers    *mock_containers.MockBotClient
	botPool          *mock_lifecycle.MockBotPoolUpdater

	botManager *botLifecycleManager

	suite.Suite
}

func TestBotLifecycleManagerTestSuite(t *testing.T) {
	suite.Run(t, &BotLifecycleManagerTestSuite{})
}

func (s *BotLifecycleManagerTestSuite) SetupTest() {
	s.r = s.Require()
	botRemoveTimeout = 0

	ctrl := gomock.NewController(s.T())
	s.msgClient = mock_clients.NewMockMessageClient(ctrl)
	s.lifecycleMetrics = mock_metrics.NewMockLifecycle(ctrl)
	s.botGrpc = mock_agentgrpc.NewMockClient(ctrl)
	s.botRegistry = mock_registry.NewMockBotRegistry(ctrl)
	s.botContainers = mock_containers.NewMockBotClient(ctrl)
	s.botPool = mock_lifecycle.NewMockBotPoolUpdater(ctrl)

	s.botManager = NewManager(s.botRegistry, s.botContainers, s.botPool, s.lifecycleMetrics)
}

func (s *BotLifecycleManagerTestSuite) TestAddUpdateRemove() {
	alreadyRunning := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
		{
			ID:    testBotID2,
			Image: testImageRef,
		},
	}
	latestAssigned := []config.AgentConfig{
		{
			ID:    testBotID3,
			Image: testImageRef,
		},
		{
			ID:    testBotID1,
			Image: testImageRef,
			ShardConfig: &config.ShardConfig{
				ShardID: 1,
			},
		},
	}
	addedBot := latestAssigned[0]
	removedBot := alreadyRunning[1]

	s.botManager.runningBots = alreadyRunning

	s.botRegistry.EXPECT().LoadAssignedBots().Return(latestAssigned, nil).Times(1)

	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), []config.AgentConfig{addedBot}).Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), addedBot).Return(nil).Times(1)

	s.botPool.EXPECT().RemoveBotsWithConfigs([]config.AgentConfig{removedBot})
	s.lifecycleMetrics.EXPECT().StatusStopping([]config.AgentConfig{removedBot})
	s.botContainers.EXPECT().TearDownBot(gomock.Any(), removedBot)

	s.lifecycleMetrics.EXPECT().StatusRunning(latestAssigned).Times(1)
	s.botPool.EXPECT().UpdateBotsWithLatestConfigs(latestAssigned)

	s.r.NoError(s.botManager.ManageBots(context.Background()))
}

func (s *BotLifecycleManagerTestSuite) TestRestart() {
	botConfigs := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
		{
			ID:    testBotID2,
			Image: testImageRef,
		},
	}

	s.botManager.runningBots = botConfigs

	dockerContainerName1 := fmt.Sprintf("/%s", botConfigs[0].ContainerName())
	dockerContainerName2 := fmt.Sprintf("/%s", botConfigs[1].ContainerName())

	s.botContainers.EXPECT().LoadBotContainers(gomock.Any()).Return([]types.Container{
		{
			ID:    testContainerID1,
			Names: []string{dockerContainerName1},
			State: "exited",
		},
		{
			ID:    testContainerID2,
			Names: []string{dockerContainerName2},
			State: "exited",
		},
	}, nil).Times(1)

	s.lifecycleMetrics.EXPECT().ActionRestart(botConfigs[0])
	s.botContainers.EXPECT().StartWaitBotContainer(gomock.Any(), testContainerID1).Return(nil)

	s.lifecycleMetrics.EXPECT().ActionRestart(botConfigs[1])
	s.botContainers.EXPECT().StartWaitBotContainer(gomock.Any(), testContainerID2).Return(errors.New("failed to start"))

	// reinitialize only
	s.botPool.EXPECT().ReinitBotsWithConfigs([]config.AgentConfig{botConfigs[0]})

	s.r.NoError(s.botManager.RestartExitedBots(context.Background()))
}
