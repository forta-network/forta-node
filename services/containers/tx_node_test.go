package containers

import (
	"context"
	"fmt"
	"testing"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	mock_clients "github.com/forta-protocol/forta-node/clients/mocks"
	"github.com/forta-protocol/forta-node/config"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testImageRef           = "some.docker.registry.io/foobar@sha256:cdd4ddccf5e9c740eb4144bcc68e3ea3a056789ec7453e94a6416dcfc80937a4"
	testNodeNetworkID      = "node-network-id"
	testScannerContainerID = "test-scanner-container-id"
	testProxyContainerID   = "test-proxy-container-id"
	testAgentID            = "test-agent"
	testAgentContainerName = "forta-agent-test-age-cdd4" // This is a result
	testAgentNetworkID     = "test-agent-network-id"
	testAgentContainerID   = "test-agent-container-id"
)

// TestSuite runs the test suite.
func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

// Suite is a test suite to test the tx node runner implementation.
type Suite struct {
	r *require.Assertions

	dockerClient     *mock_clients.MockDockerClient
	dockerAuthClient *mock_clients.MockDockerClient

	msgClient *mock_clients.MockMessageClient

	service *TxNodeService

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

	return c1.Name == c2.Name
}

// String implements the gomock.Matcher interface.
func (m configMatcher) String() string {
	return fmt.Sprintf("%+v", (clients.DockerContainerConfig)(m))
}

// SetupTest sets up the test.
func (s *Suite) SetupTest() {
	s.r = require.New(s.T())
	s.dockerClient = mock_clients.NewMockDockerClient(gomock.NewController(s.T()))
	s.dockerAuthClient = mock_clients.NewMockDockerClient(gomock.NewController(s.T()))

	s.msgClient = mock_clients.NewMockMessageClient(gomock.NewController(s.T()))
	service := &TxNodeService{
		ctx:        context.Background(),
		client:     s.dockerClient,
		authClient: s.dockerAuthClient,
		msgClient:  s.msgClient,
	}
	service.config.Config.Log.Level = "debug"

	s.dockerClient.EXPECT().Prune(service.ctx)
	s.dockerClient.EXPECT().CreatePublicNetwork(service.ctx, gomock.Any()).Return(testNodeNetworkID, nil)
	s.dockerClient.EXPECT().StartContainer(service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name: config.DockerNatsContainerName,
	})).Return(&clients.DockerContainer{}, nil)
	s.dockerClient.EXPECT().StartContainer(service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name: config.DockerQueryContainerName,
	})).Return(&clients.DockerContainer{}, nil)
	s.dockerClient.EXPECT().StartContainer(service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name: config.DockerJSONRPCProxyContainerName,
	})).Return(&clients.DockerContainer{ID: testProxyContainerID}, nil)
	s.dockerClient.EXPECT().StartContainer(service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name: config.DockerScannerContainerName,
	})).Return(&clients.DockerContainer{ID: testScannerContainerID}, nil)
	s.dockerClient.EXPECT().HasLocalImage(service.ctx, gomock.Any()).Return(true).AnyTimes()

	s.r.NoError(service.start())
	s.service = service
}

func testAgentData() (config.AgentConfig, messaging.AgentPayload) {
	agentConfig := config.AgentConfig{
		ID:    testAgentID,
		Image: testImageRef,
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
	s.dockerClient.EXPECT().CreatePublicNetwork(s.service.ctx, testAgentContainerName).Return(testAgentNetworkID, nil)
	s.dockerClient.EXPECT().StartContainer(s.service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name: agentConfig.ContainerName(),
	})).Return(&clients.DockerContainer{Name: agentConfig.ContainerName(), ID: testAgentContainerID}, nil)
	s.dockerClient.EXPECT().AttachNetwork(s.service.ctx, testScannerContainerID, testAgentNetworkID)
	s.dockerClient.EXPECT().AttachNetwork(s.service.ctx, testProxyContainerID, testAgentNetworkID)
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsStatusRunning, agentPayload)

	s.r.NoError(s.service.handleAgentRun(agentPayload))
}

// TestAgentRunAgain tests running an agent twice.
func (s *Suite) TestAgentRunAgain() {
	s.TestAgentRun()

	_, agentPayload := testAgentData()
	// Expect it to only publish a message again to ensure the subscribers that
	// the agent is running.
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsStatusRunning, agentPayload)

	s.r.NoError(s.service.handleAgentRun(agentPayload))
}

// TestAgentStop tests stopping an agent.
func (s *Suite) TestAgentStopOne() {
	s.TestAgentRun()

	_, agentPayload := testAgentData()
	// Stops the agent container and publishes a "stopped" message.
	s.dockerClient.EXPECT().StopContainer(s.service.ctx, testAgentContainerID)
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsStatusStopped, agentPayload)

	s.r.NoError(s.service.handleAgentStop(agentPayload))
}

// TestAgentStopNone tests stopping when there are no agents.
func (s *Suite) TestAgentStopNone() {
	s.TestAgentRun()

	agentPayload := messaging.AgentPayload{}

	s.r.NoError(s.service.handleAgentStop(agentPayload))
}
