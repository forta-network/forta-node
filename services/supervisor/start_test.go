package supervisor

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/forta-network/forta-core-go/release"

	"github.com/docker/docker/api/types"

	mrelease "github.com/forta-network/forta-core-go/release/mocks"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testImageRef               = "some.docker.registry.io/foobar@sha256:cdd4ddccf5e9c740eb4144bcc68e3ea3a056789ec7453e94a6416dcfc80937a4"
	testNodeNetworkID          = "node-network-id"
	testNatsNetworkID          = "nats-network-id"
	testPublicAPINetworkID     = "public-api-network-id"
	testGenericContainerID     = "test-generic-container-id"
	testInspectorContainerID   = "test-inspector-container-id"
	testScannerContainerID     = "test-scanner-container-id"
	testProxyContainerID       = "test-proxy-container-id"
	testPublicAPIContainerID   = "test-public-api-container-id"
	testSupervisorContainerID  = "test-supervisor-container-id"
	testAgentID                = "test-agent"
	testAgentContainerName     = "forta-agent-test-age-cdd4" // This is a result
	testAgentNetworkID         = "test-agent-network-id"
	testAgentContainerID       = "test-agent-container-id"
	testJWTProviderContainerID = "test-jwt-provider-container-id"
)

// TestSuite runs the test suite.
func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

// Suite is a test suite to test the tx node runner implementation.
type Suite struct {
	r *require.Assertions

	dockerClient     *mock_clients.MockDockerClient
	globalClient     *mock_clients.MockDockerClient
	agentImageClient *mock_clients.MockDockerClient
	releaseClient    *mrelease.MockClient

	msgClient *mock_clients.MockMessageClient

	service *SupervisorService

	suite.Suite
}

// configMatcher is a wrapper to implement the Matcher interface.
type configMatcher clients.DockerContainerConfig

// Matches implements the gomock.Matcher interface.
func (m configMatcher) Matches(x interface{}) bool {
	c1, ok := x.(clients.DockerContainerConfig)
	if !ok {
		return false
	}
	c2 := m

	if c2.Env != nil && c1.Env == nil {
		return false
	}

	for k2, v2 := range c2.Env {
		if v1, ok := c1.Env[k2]; !ok {
			return false
		} else {
			if v1 != v2 {
				return false
			}
		}

	}

	return c1.Name == c2.Name
}

// String implements the gomock.Matcher interface.
func (m configMatcher) String() string {
	return fmt.Sprintf("%+v", (clients.DockerContainerConfig)(m))
}

// SetupTest sets up the test.
func (s *Suite) SetupTest() {
	s.r = require.New(s.T())
	os.Setenv(config.EnvHostFortaDir, "/tmp/forta")
	s.dockerClient = mock_clients.NewMockDockerClient(gomock.NewController(s.T()))
	s.globalClient = mock_clients.NewMockDockerClient(gomock.NewController(s.T()))
	s.agentImageClient = mock_clients.NewMockDockerClient(gomock.NewController(s.T()))
	s.releaseClient = mrelease.NewMockClient(gomock.NewController(s.T()))

	s.msgClient = mock_clients.NewMockMessageClient(gomock.NewController(s.T()))
	service := &SupervisorService{
		ctx:              context.Background(),
		client:           s.dockerClient,
		globalClient:     s.globalClient,
		msgClient:        s.msgClient,
		releaseClient:    s.releaseClient,
		agentImageClient: s.agentImageClient,
	}
	service.config.Config.TelemetryConfig.Disable = true
	service.config.Config.Log.Level = "debug"
	service.config.Config.ChainID = 1
	service.config.Config.AdvancedConfig.IPFSExperiment = true
	s.service = service

	s.releaseClient.EXPECT().GetReleaseManifest(gomock.Any()).Return(&release.ReleaseManifest{}, nil).AnyTimes()

	s.initialContainerCheck()
	s.dockerClient.EXPECT().EnsureLocalImage(service.ctx, gomock.Any(), gomock.Any()).Times(2) // needs to get nats and ipfs
	s.dockerClient.EXPECT().CreatePublicNetwork(service.ctx, gomock.Any()).Return(testNodeNetworkID, nil)
	s.dockerClient.EXPECT().CreateInternalNetwork(service.ctx, gomock.Any()).Return(testNatsNetworkID, nil) // for nats
	s.dockerClient.EXPECT().StartContainer(
		service.ctx, (configMatcher)(
			clients.DockerContainerConfig{
				Name: config.DockerIpfsContainerName,
			},
		),
	).Return(&clients.DockerContainer{}, nil)
	s.dockerClient.EXPECT().StartContainer(
		service.ctx, (configMatcher)(
			clients.DockerContainerConfig{
				Name: config.DockerStorageContainerName,
			},
		),
	).Return(&clients.DockerContainer{}, nil)
	s.dockerClient.EXPECT().StartContainer(
		service.ctx, (configMatcher)(
			clients.DockerContainerConfig{
				Name: config.DockerNatsContainerName,
			},
		),
	).Return(&clients.DockerContainer{}, nil)
	s.dockerClient.EXPECT().StartContainer(
		service.ctx, (configMatcher)(
			clients.DockerContainerConfig{
				Name: config.DockerJSONRPCProxyContainerName,
			},
		),
	).Return(&clients.DockerContainer{ID: testProxyContainerID}, nil)
	s.dockerClient.EXPECT().StartContainer(
		service.ctx, (configMatcher)(
			clients.DockerContainerConfig{
				Name: config.DockerPublicAPIProxyContainerName,
			},
		),
	).Return(&clients.DockerContainer{ID: testPublicAPIContainerID}, nil)
	s.dockerClient.EXPECT().StartContainer(
		service.ctx, (configMatcher)(
			clients.DockerContainerConfig{
				Name: config.DockerScannerContainerName,
			},
		),
	).Return(&clients.DockerContainer{ID: testScannerContainerID}, nil)
	s.dockerClient.EXPECT().StartContainer(
		service.ctx, (configMatcher)(
			clients.DockerContainerConfig{
				Name: config.DockerJWTProviderContainerName,
			},
		),
	).Return(&clients.DockerContainer{ID: testJWTProviderContainerID}, nil)
	s.dockerClient.EXPECT().StartContainer(
		service.ctx, (configMatcher)(
			clients.DockerContainerConfig{
				Name: config.DockerInspectorContainerName,
			},
		),
	).Return(&clients.DockerContainer{ID: testProxyContainerID}, nil)
	s.dockerClient.EXPECT().HasLocalImage(service.ctx, gomock.Any()).Return(true, nil).AnyTimes()
	s.globalClient.EXPECT().GetContainerByName(service.ctx, config.DockerSupervisorContainerName).Return(&types.Container{ID: testSupervisorContainerID}, nil).AnyTimes()
	s.dockerClient.EXPECT().AttachNetwork(service.ctx, testSupervisorContainerID, testNodeNetworkID)
	s.dockerClient.EXPECT().AttachNetwork(service.ctx, testSupervisorContainerID, testNatsNetworkID)
	s.dockerClient.EXPECT().GetContainerByName(service.ctx, config.DockerJSONRPCProxyContainerName).Return(&types.Container{ID: testProxyContainerID}, nil).AnyTimes()
	s.dockerClient.EXPECT().GetContainerByName(service.ctx, config.DockerInspectorContainerName).Return(&types.Container{ID: testInspectorContainerID}, nil).AnyTimes()
	s.dockerClient.EXPECT().GetContainerByName(service.ctx, config.DockerScannerContainerName).Return(&types.Container{ID: testScannerContainerID}, nil).AnyTimes()
	s.dockerClient.EXPECT().GetContainerByName(service.ctx, config.DockerScannerContainerName).Return(&types.Container{ID: testScannerContainerID}, nil).AnyTimes()
	s.dockerClient.EXPECT().GetContainerByName(
		service.ctx,
		config.DockerJWTProviderContainerName,
	).Return(&types.Container{ID: testJWTProviderContainerID}, nil).AnyTimes()
	s.dockerClient.EXPECT().WaitContainerStart(service.ctx, gomock.Any()).Return(nil).AnyTimes()
	s.msgClient.EXPECT().Subscribe(messaging.SubjectAgentsActionRun, gomock.Any())
	s.msgClient.EXPECT().Subscribe(messaging.SubjectAgentsActionStop, gomock.Any())

	s.r.NoError(service.start())
}

func (s *Suite) initialContainerCheck() {
	for _, containerName := range knownServiceContainerNames {
		s.dockerClient.EXPECT().GetContainerByName(s.service.ctx, containerName).Return(&types.Container{ID: testGenericContainerID}, nil)
	}

	s.dockerClient.EXPECT().GetContainers(s.service.ctx).Return(
		[]types.Container{
			{
				Names: []string{"/forta-agent-name"},
				ID:    testGenericContainerID,
				Labels: map[string]string{
					clients.DockerLabelFortaSupervisorStrategyVersion: SupervisorStrategyVersion,
				},
			},
			{
				Names: []string{"/forta-agent-name"},
				ID:    testGenericContainerID,
				Labels: map[string]string{
					clients.DockerLabelFortaSupervisorStrategyVersion: "old",
				},
			},
		}, nil,
	)

	// service containers + 1 old agent
	for i := 0; i < len(knownServiceContainerNames)+1; i++ {
		s.dockerClient.EXPECT().RemoveContainer(s.service.ctx, testGenericContainerID).Return(nil)
		s.dockerClient.EXPECT().WaitContainerPrune(s.service.ctx, testGenericContainerID).Return(nil)
	}
	for i := 0; i < len(knownServiceContainerNames)+1; i++ {
		s.dockerClient.EXPECT().RemoveNetworkByName(s.service.ctx, gomock.Any()).Return(nil)
	}
}

func testAgentData() (config.AgentConfig, messaging.AgentPayload) {
	agentConfig := config.AgentConfig{
		ID:      testAgentID,
		Image:   testImageRef,
		ChainID: 1,
	}
	return agentConfig, messaging.AgentPayload{
		agentConfig,
	}
}

// TestAgentRun tests running the agent.
func (s *Suite) TestAgentRun() {
	agentConfig, agentPayload := testAgentData()
	// Creates the agent network, starts the agent container, attaches the scanner and the proxy to the
	// agent network, publishes a "running" message.
	s.agentImageClient.EXPECT().EnsureLocalImage(gomock.Any(), "agent test-agent", agentConfig.Image).Return(nil)
	s.dockerClient.EXPECT().CreatePublicNetwork(gomock.Any(), testAgentContainerName).Return(testAgentNetworkID, nil)
	s.dockerClient.EXPECT().StartContainer(
		gomock.Any(), (configMatcher)(
			clients.DockerContainerConfig{
				Name: agentConfig.ContainerName(),
				Env: map[string]string{
					config.EnvFortaChainID: "1",
				},
			},
		),
	).Return(&clients.DockerContainer{Name: agentConfig.ContainerName(), ID: testAgentContainerID}, nil)

	s.dockerClient.EXPECT().AttachNetwork(gomock.Any(), testScannerContainerID, testAgentNetworkID)
	s.dockerClient.EXPECT().AttachNetwork(gomock.Any(), testProxyContainerID, testAgentNetworkID)
	s.dockerClient.EXPECT().AttachNetwork(gomock.Any(), testJWTProviderContainerID, testAgentNetworkID)
	s.dockerClient.EXPECT().AttachNetwork(gomock.Any(), testPublicAPIContainerID, testAgentNetworkID)
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsStatusRunning, agentPayload)
	s.msgClient.EXPECT().PublishProto(messaging.SubjectMetricAgent, gomock.Any()).AnyTimes()

	s.r.NoError(s.service.handleAgentRunWithContext(s.service.ctx, agentPayload))
}

// TestAgentRunAgain tests running an agent twice.
func (s *Suite) TestAgentRunAgain() {
	s.TestAgentRun()

	agentConfig, agentPayload := testAgentData()

	// Expect it to only publish a message again to ensure the subscribers that
	// the agent is running.
	s.agentImageClient.EXPECT().EnsureLocalImage(gomock.Any(), "agent test-agent", agentConfig.Image).Return(nil)
	s.dockerClient.EXPECT().CreatePublicNetwork(gomock.Any(), gomock.Any()).Return(testNodeNetworkID, nil)
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsStatusRunning, agentPayload)
	s.msgClient.EXPECT().PublishProto(messaging.SubjectMetricAgent, gomock.Any()).AnyTimes()

	s.r.NoError(s.service.handleAgentRunWithContext(s.service.ctx, agentPayload))
}

// TestAgentStop tests stopping an agent.
func (s *Suite) TestAgentStopOne() {
	s.TestAgentRun()

	_, agentPayload := testAgentData()
	// Stops the agent container and publishes a "stopped" message.
	s.dockerClient.EXPECT().StopContainer(s.service.ctx, testAgentContainerID)
	s.dockerClient.EXPECT().RemoveNetworkByName(s.service.ctx, testAgentContainerID)
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsStatusStopped, agentPayload)

	s.r.NoError(s.service.handleAgentStop(agentPayload))
}

// TestAgentStopNone tests stopping when there are no agents.
func (s *Suite) TestAgentStopNone() {
	s.TestAgentRun()

	agentPayload := messaging.AgentPayload{}

	s.r.NoError(s.service.handleAgentStop(agentPayload))
}
