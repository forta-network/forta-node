package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/forta-network/forta-core-go/protocol"
	mock_agentgrpc "github.com/forta-network/forta-node/clients/agentgrpc/mocks"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/botio"
	"github.com/forta-network/forta-node/services/components/botio/botreq"
	mock_containers "github.com/forta-network/forta-node/services/components/containers/mocks"
	mock_lifecycle "github.com/forta-network/forta-node/services/components/lifecycle/mocks"
	mock_metrics "github.com/forta-network/forta-node/services/components/metrics/mocks"
	mock_registry "github.com/forta-network/forta-node/services/components/registry/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testBotID1       = "0x0100000000000000000000000000000000000000000000000000000000000000"
	testBotID2       = "0x0200000000000000000000000000000000000000000000000000000000000000"
	testBotID3       = "0x0300000000000000000000000000000000000000000000000000000000000000"
	testImageRef1    = "test-image-ref-1"
	testImageRef2    = "test-image-ref-2"
	testImageRef3    = "test-image-ref-3"
	testImageRef     = testImageRef1
	testContainerID  = "test-container-id"
	testContainerID1 = "test-container-id-1"
	testContainerID2 = "test-container-id-2"
	testContainerID3 = "test-container-id-3"
)

// LifecycleTestSuite composes type botLifecycleManager with a concrete type botPool
// and verifies that the bots will be managed as expected after assignment.
//
// This is different from bot pool and bot manager unit tests and acts as a component test
// which combines the two and verifies acceptance criteria in Given-When-Then style.
//
// The bot manager and the bot pool are expected to run in separate docker containers
// and stay connected via a mediator (see package mediator). This test avoids that complexity
// by making the bot manager call the concrete bot pool directly.
type LifecycleTestSuite struct {
	r *require.Assertions

	msgClient        *mock_clients.MockMessageClient
	lifecycleMetrics *mock_metrics.MockLifecycle
	botGrpc          *mock_agentgrpc.MockClient
	botRegistry      *mock_registry.MockBotRegistry
	botContainers    *mock_containers.MockBotClient
	dialer           *mock_agentgrpc.MockBotDialer
	botMonitor       *mock_lifecycle.MockBotMonitor

	resultChannels botreq.SendReceiveChannels

	botPool    *botPool
	botManager *botLifecycleManager

	suite.Suite
}

func (s *LifecycleTestSuite) DisableHeartbeatBot() {
	s.botManager.setLastHeartbeatTime(time.Now().UTC().Add(-10 * time.Minute))
}

func (s *LifecycleTestSuite) EnableHeartbeatBot() {
	s.botManager.setLastHeartbeatTime(time.Time{})
}

func TestLifecycleTestSuite(t *testing.T) {
	suite.Run(t, &LifecycleTestSuite{})
}

func (s *LifecycleTestSuite) SetupTest() {
	s.r = s.Require()
	botRemoveTimeout = 0

	ctrl := gomock.NewController(s.T())
	s.msgClient = mock_clients.NewMockMessageClient(ctrl)
	s.lifecycleMetrics = mock_metrics.NewMockLifecycle(ctrl)
	s.botGrpc = mock_agentgrpc.NewMockClient(ctrl)
	s.botRegistry = mock_registry.NewMockBotRegistry(ctrl)
	s.botContainers = mock_containers.NewMockBotClient(ctrl)
	s.dialer = mock_agentgrpc.NewMockBotDialer(ctrl)
	s.resultChannels = botreq.MakeResultChannels()
	s.botMonitor = mock_lifecycle.NewMockBotMonitor(ctrl)

	// expecting any amount of health calls as we do not focus on testing it in this suite
	s.lifecycleMetrics.EXPECT().HealthCheckAttempt(gomock.Any()).AnyTimes()
	s.botGrpc.EXPECT().DoHealthCheck(gomock.Any()).AnyTimes()
	s.lifecycleMetrics.EXPECT().HealthCheckSuccess(gomock.Any()).AnyTimes()

	botClientFactory := botio.NewBotClientFactory(s.resultChannels.SendOnly(), s.msgClient, s.lifecycleMetrics, s.dialer)
	s.botPool = NewBotPool(context.Background(), s.lifecycleMetrics, botClientFactory, 0)
	s.botPool.waitInit = true // hack to make testing synchronous
	s.botManager = NewManager(s.botRegistry, s.botContainers, s.botPool, s.lifecycleMetrics, s.botMonitor)
}

func (s *LifecycleTestSuite) TestDownloadTimeout() {
	s.T().Log("should redownload a bot if downloading times out")

	s.DisableHeartbeatBot()
	assigned := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}

	// given that there is a new bot assignment
	// and no new assignments after the first time
	s.botRegistry.EXPECT().LoadAssignedBots().Return(assigned, nil).Times(2)
	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "1").Times(2)

	// then the bot should be redownloaded, launched, dialed and initialized
	// upon download timeouts for the first time

	err := errors.New("download timeout")
	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), assigned).
		Return([]error{err}).Times(1)
	s.lifecycleMetrics.EXPECT().FailurePull(err, assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusRunning().Times(1) // not bots running due to download failure

	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), assigned).
		Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), assigned[0]).Return(nil).Times(1)
	s.lifecycleMetrics.EXPECT().StatusRunning(assigned[0]).Times(1) // bot is running

	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusAttached(assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusInitialized(assigned[0]).Times(1)
	s.dialer.EXPECT().DialBot(assigned[0]).Return(s.botGrpc, nil).Times(1)
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{}, nil).
		Times(1)

	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(nil)).Times(1)
	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(assigned)).Times(1)

	// when the bot manager manages the assigned bots over time
	s.r.NoError(s.botManager.ManageBots(context.Background()))
	s.r.NoError(s.botManager.ManageBots(context.Background()))
}

func (s *LifecycleTestSuite) TestLaunchFailure() {
	s.T().Log("should relaunch a bot if launching fails")

	s.DisableHeartbeatBot()
	assigned := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}

	// given that there is a new bot assignment
	// and no new assignments after the first time
	s.botRegistry.EXPECT().LoadAssignedBots().Return(assigned, nil).Times(2)
	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "1").Times(2)

	// then the bot should be relaunched, dialed and initialized
	// upon launch failure for the first time

	err := errors.New("failed to launch")
	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), assigned).
		Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), assigned[0]).Return(err).Times(1)
	s.lifecycleMetrics.EXPECT().FailureLaunch(err, assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusRunning().Times(1) // not bots running due to download failure

	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), assigned).
		Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), assigned[0]).Return(nil).Times(1)
	s.lifecycleMetrics.EXPECT().StatusRunning(assigned[0]).Times(1) // bot is running

	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusAttached(assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusInitialized(assigned[0]).Times(1)
	s.dialer.EXPECT().DialBot(assigned[0]).Return(s.botGrpc, nil).Times(1)
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{}, nil).
		Times(1)

	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(assigned)).Times(1)
	s.botMonitor.EXPECT().MonitorBots(nil).Times(1)

	// when the bot manager manages the assigned bots over time
	s.r.NoError(s.botManager.ManageBots(context.Background()))
	s.r.NoError(s.botManager.ManageBots(context.Background()))
}

func (s *LifecycleTestSuite) TestDialFailure() {
	s.T().Log("should not reload or redial a bot if dialing finally fails")

	s.DisableHeartbeatBot()
	assigned := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}

	// given that there is a new bot assignment
	// and no new assignments after the first time
	s.botRegistry.EXPECT().LoadAssignedBots().Return(assigned, nil).Times(2)
	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "1").Times(2)

	// then there should be no reloading and redialing upon dialing failures
	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), assigned).Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), assigned[0]).Return(nil).Times(1)
	s.lifecycleMetrics.EXPECT().StatusRunning(assigned[0]).Times(2)
	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0]).Times(1)
	s.dialer.EXPECT().DialBot(assigned[0]).Return(nil, errors.New("failed to dial")).Times(1)

	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(assigned)).Times(2)

	// when the bot manager manages the assigned bots over time
	s.r.NoError(s.botManager.ManageBots(context.Background()))
	s.r.NoError(s.botManager.ManageBots(context.Background()))
}

func (s *LifecycleTestSuite) TestInitializeFailure() {
	s.T().Log("should reconnect to a bot if initialization finally fails")

	s.DisableHeartbeatBot()
	assigned := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}

	// given that there is a new bot assignment
	// and no new assignments after the first time
	s.botRegistry.EXPECT().LoadAssignedBots().Return(assigned, nil).Times(2)
	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "1").Times(2)

	err := errors.New("failed to init")

	// then there should be no reloading and redialing upon initialization failures

	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), assigned).Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), assigned[0]).Return(nil).Times(1)

	s.lifecycleMetrics.EXPECT().StatusRunning(assigned[0])
	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0])
	s.dialer.EXPECT().DialBot(assigned[0]).Return(s.botGrpc, nil)
	s.lifecycleMetrics.EXPECT().StatusAttached(assigned[0])
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(nil, err)
	s.lifecycleMetrics.EXPECT().FailureInitialize(err, assigned[0])
	s.botGrpc.EXPECT().Close()
	s.lifecycleMetrics.EXPECT().ClientClose(assigned[0])

	s.lifecycleMetrics.EXPECT().StatusRunning(assigned[0])
	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0])
	s.dialer.EXPECT().DialBot(assigned[0]).Return(s.botGrpc, nil)
	s.lifecycleMetrics.EXPECT().StatusAttached(assigned[0])
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{}, nil)
	s.lifecycleMetrics.EXPECT().StatusInitialized(assigned[0])

	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(assigned)).Times(2)

	// when the bot manager manages the assigned bots over time
	s.r.NoError(s.botManager.ManageBots(context.Background()))
	s.r.NoError(s.botManager.ManageBots(context.Background()))
}

func (s *LifecycleTestSuite) TestExitedRestarted() {
	s.T().Log("should restart and reconnect to exited bots")

	s.DisableHeartbeatBot()
	assigned := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}

	// given that there is a new bot assignment
	// and no new assignments after the first time
	s.botRegistry.EXPECT().LoadAssignedBots().Return(assigned, nil).Times(2)
	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "1").Times(2)

	// then there should be restart and reinitialization

	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), assigned).Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), assigned[0]).Return(nil).Times(1)

	s.lifecycleMetrics.EXPECT().StatusRunning(assigned[0])
	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0])
	s.dialer.EXPECT().DialBot(assigned[0]).Return(s.botGrpc, nil)
	s.lifecycleMetrics.EXPECT().StatusAttached(assigned[0])
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{}, nil)
	s.lifecycleMetrics.EXPECT().StatusInitialized(assigned[0])

	s.lifecycleMetrics.EXPECT().StatusRunning(assigned[0])
	s.botGrpc.EXPECT().Close()
	s.lifecycleMetrics.EXPECT().ClientClose(assigned[0])
	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0])
	s.dialer.EXPECT().DialBot(assigned[0]).Return(s.botGrpc, nil)
	s.lifecycleMetrics.EXPECT().StatusAttached(assigned[0])
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{}, nil)
	s.lifecycleMetrics.EXPECT().StatusInitialized(assigned[0])

	dockerContainerName := fmt.Sprintf("/%s", assigned[0].ContainerName())

	s.botContainers.EXPECT().LoadBotContainers(gomock.Any()).Return([]types.Container{
		{
			ID:    testContainerID,
			Names: []string{dockerContainerName},
			State: "running",
		},
	}, nil).Times(1)
	s.botContainers.EXPECT().LoadBotContainers(gomock.Any()).Return([]types.Container{
		{
			ID:    testContainerID,
			Names: []string{dockerContainerName},
			State: "exited",
		},
	}, nil).Times(1)

	s.lifecycleMetrics.EXPECT().ActionRestart(assigned[0])
	s.botContainers.EXPECT().StartWaitBotContainer(gomock.Any(), testContainerID).Return(nil)
	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(assigned)).Times(2)

	// when the bot manager manages the assigned bots over time
	s.r.NoError(s.botManager.ManageBots(context.Background()))
	s.r.NoError(s.botManager.RestartExitedBots(context.Background()))
	s.r.NoError(s.botManager.ManageBots(context.Background()))
	s.r.NoError(s.botManager.RestartExitedBots(context.Background()))
}

func (s *LifecycleTestSuite) TestInactiveRestarted() {
	s.T().Log("should restart and reconnect to the inactive and exited bots")

	s.DisableHeartbeatBot()
	assigned := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}

	// given that there is a new bot assignment
	// and no new assignments after the first time
	s.botRegistry.EXPECT().LoadAssignedBots().Return(assigned, nil).Times(1)
	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "1")

	// then there should be restart and reinitialization

	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), assigned).Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), assigned[0]).Return(nil).Times(1)

	s.lifecycleMetrics.EXPECT().StatusRunning(assigned[0])
	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0])
	s.dialer.EXPECT().DialBot(assigned[0]).Return(s.botGrpc, nil)
	s.lifecycleMetrics.EXPECT().StatusAttached(assigned[0])
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{}, nil)
	s.lifecycleMetrics.EXPECT().StatusInitialized(assigned[0])

	s.botGrpc.EXPECT().Close()
	s.lifecycleMetrics.EXPECT().ClientClose(assigned[0])
	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0])
	s.lifecycleMetrics.EXPECT().StatusInactive(assigned[0])
	s.dialer.EXPECT().DialBot(assigned[0]).Return(s.botGrpc, nil)
	s.lifecycleMetrics.EXPECT().StatusAttached(assigned[0])
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{}, nil)
	s.lifecycleMetrics.EXPECT().StatusInitialized(assigned[0])

	s.botMonitor.EXPECT().GetInactiveBots().Return([]string{testBotID1})
	s.botContainers.EXPECT().StopBot(gomock.Any(), assigned[0])

	dockerContainerName := fmt.Sprintf("/%s", assigned[0].ContainerName())

	s.botContainers.EXPECT().LoadBotContainers(gomock.Any()).Return([]types.Container{
		{
			ID:    testContainerID,
			Names: []string{dockerContainerName},
			State: "exited",
		},
	}, nil).Times(1)

	s.lifecycleMetrics.EXPECT().ActionRestart(assigned[0])
	s.botContainers.EXPECT().StartWaitBotContainer(gomock.Any(), testContainerID).Return(nil)

	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(assigned))

	// when the bot manager manages the assigned bots over time
	s.r.NoError(s.botManager.ManageBots(context.Background()))
	s.r.NoError(s.botManager.ExitInactiveBots(context.Background()))
	s.r.NoError(s.botManager.RestartExitedBots(context.Background()))
}

func (s *LifecycleTestSuite) TestUnassigned() {
	s.T().Log("should tear down unassigned bots")

	s.DisableHeartbeatBot()
	assigned := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}

	// given that there is a new bot assignment
	s.botRegistry.EXPECT().LoadAssignedBots().Return(assigned, nil).Times(1)
	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "1")
	// and the assignment is removed shortly
	s.botRegistry.EXPECT().LoadAssignedBots().Return(nil, nil).Times(1)
	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "0")

	// then the bot should be started
	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), assigned).Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), assigned[0]).Return(nil).Times(1)
	s.lifecycleMetrics.EXPECT().StatusRunning(assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusAttached(assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusInitialized(assigned[0]).Times(1)
	s.dialer.EXPECT().DialBot(assigned[0]).Return(s.botGrpc, nil).Times(1)
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{}, nil).
		Times(1)

	// and should shortly be torn down
	s.lifecycleMetrics.EXPECT().StatusStopping(assigned[0])
	s.botContainers.EXPECT().TearDownBot(gomock.Any(), assigned[0].ContainerName(), true).Return(nil)
	s.lifecycleMetrics.EXPECT().StatusRunning().Times(1)
	s.lifecycleMetrics.EXPECT().ClientClose(assigned[0])
	s.botGrpc.EXPECT().Close().AnyTimes()

	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(assigned)).Times(1)
	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(nil)).Times(1)

	// when the bot manager manages the assigned bots over time
	s.r.NoError(s.botManager.ManageBots(context.Background()))
	createdBotClient := s.botPool.GetCurrentBotClients()[0]
	s.r.NoError(s.botManager.ManageBots(context.Background()))
	<-createdBotClient.Closed()
}

func (s *LifecycleTestSuite) TestHearbeatBotLoads() {
	s.T().Log("should load heartbeat bot")
	s.EnableHeartbeatBot()

	botsToRun := []config.AgentConfig{
		*heartbeatBot,
	}

	// given that there is a new bot assignment
	// and no new assignments after the first time
	s.botRegistry.EXPECT().LoadAssignedBots().Return(nil, nil).Times(1)
	s.botRegistry.EXPECT().LoadHeartbeatBot().Return(heartbeatBot, nil).Times(1)

	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "0").Times(1)

	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), botsToRun).
		Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), botsToRun[0]).Return(nil).Times(1)
	s.lifecycleMetrics.EXPECT().StatusRunning(botsToRun[0]).Times(1) // bot is running

	s.lifecycleMetrics.EXPECT().ClientDial(botsToRun[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusAttached(botsToRun[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusInitialized(botsToRun[0]).Times(1)
	s.dialer.EXPECT().DialBot(botsToRun[0]).Return(s.botGrpc, nil).Times(1)
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{}, nil).
		Times(1)

	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(botsToRun)).Times(1)

	// when the bot manager manages the assigned bots over time
	s.r.NoError(s.botManager.ManageBots(context.Background()))
}

func (s *LifecycleTestSuite) TestConfigUpdated() {
	s.T().Log("should update bot config without tearing down")

	s.DisableHeartbeatBot()
	assigned := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
		},
	}
	updated := []config.AgentConfig{
		{
			ID:    testBotID1,
			Image: testImageRef,
			ShardConfig: &config.ShardConfig{
				ShardID: 1,
			},
		},
	}

	// given that there is a new bot assignment
	s.botRegistry.EXPECT().LoadAssignedBots().Return(assigned, nil).Times(1)
	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "1")
	// and the assigned bot's shard config is updated shortly
	s.botRegistry.EXPECT().LoadAssignedBots().Return(updated, nil).Times(1)
	s.lifecycleMetrics.EXPECT().SystemStatus("load.assigned.bots", "1")

	// then the config of the bot should be updated

	s.botContainers.EXPECT().EnsureBotImages(gomock.Any(), assigned).Return([]error{nil}).Times(1)
	s.botContainers.EXPECT().LaunchBot(gomock.Any(), assigned[0]).Return(nil).Times(1)
	s.lifecycleMetrics.EXPECT().StatusRunning(assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().ClientDial(assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusAttached(assigned[0]).Times(1)
	s.lifecycleMetrics.EXPECT().StatusInitialized(assigned[0]).Times(1)
	s.dialer.EXPECT().DialBot(assigned[0]).Return(s.botGrpc, nil).Times(1)
	s.botGrpc.EXPECT().Initialize(gomock.Any(), gomock.Any()).Return(&protocol.InitializeResponse{}, nil).
		Times(1)

	s.lifecycleMetrics.EXPECT().StatusRunning(updated[0]).Times(1)
	s.lifecycleMetrics.EXPECT().ActionUpdate(updated[0])

	s.botMonitor.EXPECT().MonitorBots(GetBotIDs(assigned)).Times(2)

	// when the bot manager manages the assigned bots over time
	s.r.NoError(s.botManager.ManageBots(context.Background()))
	s.r.NoError(s.botManager.ManageBots(context.Background()))

	s.r.Equal(updated[0], s.botPool.botClients[0].Config())
}
