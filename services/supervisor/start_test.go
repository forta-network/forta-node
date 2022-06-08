package supervisor

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/release"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"

	mrelease "github.com/forta-network/forta-core-go/release/mocks"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	netmgmt "github.com/forta-network/forta-node/services/network"
	mock_network "github.com/forta-network/forta-node/services/network/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testImageRef                = "some.docker.registry.io/foobar@sha256:cdd4ddccf5e9c740eb4144bcc68e3ea3a056789ec7453e94a6416dcfc80937a4"
	testNodeNetworkID           = "test-node-network-id"
	testGenericContainerID      = "test-generic-container-id"
	testScannerContainerID      = "test-scanner-container-id"
	testProxyContainerID        = "test-proxy-container-id"
	testSupervisorContainerID   = "test-supervisor-container-id"
	testHostnetContainerID      = "test-hostnet-container-id"
	testAgentID                 = "test-agent"
	testAgentContainerName      = "forta-agent-test-age-cdd4" // This is a result
	testAgentAdminContainerName = "forta-agent-admin-test-age-cdd4"
	testAgentContainerID        = "test-agent-container-id"
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
	botManager       *mock_network.MockBotManager
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

	return c1.Name == c2.Name && reflect.DeepEqual(c1.LinkNetworkIDs, c2.LinkNetworkIDs) &&
		c1.NetworkID == c2.NetworkID
}

// String implements the gomock.Matcher interface.
func (m configMatcher) String() string {
	return fmt.Sprintf("%+v", (clients.DockerContainerConfig)(m))
}

// SetupTest sets up the test.
func (s *Suite) SetupTest() {
	s.r = require.New(s.T())
	os.Setenv(config.EnvHostFortaDir, "/tmp/forta")
	disableSocketDirCheck = true // testing socket dir cleanup is very tricky - disable
	s.dockerClient = mock_clients.NewMockDockerClient(gomock.NewController(s.T()))
	s.globalClient = mock_clients.NewMockDockerClient(gomock.NewController(s.T()))
	s.agentImageClient = mock_clients.NewMockDockerClient(gomock.NewController(s.T()))
	s.botManager = mock_network.NewMockBotManager(gomock.NewController(s.T()))
	s.releaseClient = mrelease.NewMockClient(gomock.NewController(s.T()))

	s.msgClient = mock_clients.NewMockMessageClient(gomock.NewController(s.T()))
	service := &SupervisorService{
		ctx:              context.Background(),
		client:           s.dockerClient,
		globalClient:     s.globalClient,
		msgClient:        s.msgClient,
		releaseClient:    s.releaseClient,
		agentImageClient: s.agentImageClient,
		botManager:       s.botManager,
	}
	service.config.Config.TelemetryConfig.Disable = true
	service.config.Config.Log.Level = "debug"
	s.service = service

	s.releaseClient.EXPECT().GetReleaseManifest(gomock.Any(), gomock.Any()).Return(&release.ReleaseManifest{}, nil).AnyTimes()

	s.initialContainerCheck()
	s.dockerClient.EXPECT().EnsureLocalImage(service.ctx, gomock.Any(), gomock.Any()).Times(2) // needs to get nats and ipfs
	s.dockerClient.EXPECT().HasLocalImage(service.ctx, gomock.Any()).Return(true).AnyTimes()

	// should get the supervisor container (self) once to find out the node image
	s.globalClient.EXPECT().GetContainerByName(service.ctx, config.DockerSupervisorContainerName).Return(&types.Container{ID: testSupervisorContainerID}, nil).AnyTimes()

	// should remove old and run new host network detection container
	s.dockerClient.EXPECT().RemoveContainer(service.ctx, config.DockerHostNetContainerName).Return(nil)
	s.dockerClient.EXPECT().StartContainer(service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name:      config.DockerHostNetContainerName,
		NetworkID: "host",
	}), false).Return(&clients.DockerContainer{ID: testHostnetContainerID}, nil)
	s.dockerClient.EXPECT().GetContainerByID(gomock.Any(), testHostnetContainerID).Return(&types.Container{
		ID:    testHostnetContainerID,
		State: "exited",
	}, nil).AnyTimes()
	s.dockerClient.EXPECT().GetContainerLogs(service.ctx, testHostnetContainerID, "", -1).Return(
		time.Now().Format(time.RFC3339)+" "+netmgmt.MarshalHostNetworking(&netmgmt.Host{
			DefaultInterfaceName: "eth0",
			DefaultSubnet:        "192.168.0.0/24",
			DefaultGateway:       "192.168.0.1",
			Docker0Subnet:        "10.99.0.0/24",
		}), nil,
	)
	s.botManager.EXPECT().Init(gomock.Any(), gomock.Any())

	// should create node network
	s.dockerClient.EXPECT().CreatePublicNetwork(service.ctx, config.DockerNodeNetworkName).Return(testNodeNetworkID, nil)
	s.dockerClient.EXPECT().GetNetworkByID(service.ctx, testNodeNetworkID).Return(types.NetworkResource{
		IPAM: network.IPAM{
			Config: []network.IPAMConfig{
				{
					Subnet: "100.100.0.0/24",
				},
			},
		},
	}, nil)

	// should attach supervisor to the service network
	s.dockerClient.EXPECT().AttachNetwork(service.ctx, testSupervisorContainerID, testNodeNetworkID)

	s.dockerClient.EXPECT().StartContainer(service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name:      config.DockerIpfsContainerName,
		NetworkID: testNodeNetworkID,
	})).Return(&clients.DockerContainer{}, nil)

	s.dockerClient.EXPECT().StartContainer(service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name:      config.DockerNatsContainerName,
		NetworkID: testNodeNetworkID,
	})).Return(&clients.DockerContainer{}, nil)

	s.dockerClient.EXPECT().StartContainer(service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name:      config.DockerJSONRPCProxyContainerName,
		NetworkID: testNodeNetworkID,
	})).Return(&clients.DockerContainer{ID: testProxyContainerID}, nil)

	s.dockerClient.EXPECT().StartContainer(service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name:      config.DockerScannerContainerName,
		NetworkID: testNodeNetworkID,
	})).Return(&clients.DockerContainer{ID: testScannerContainerID}, nil)

	s.dockerClient.EXPECT().WaitContainerStart(service.ctx, gomock.Any()).Return(nil).AnyTimes()
	s.msgClient.EXPECT().Subscribe(messaging.SubjectAgentsActionRun, gomock.Any())
	s.msgClient.EXPECT().Subscribe(messaging.SubjectAgentsActionStop, gomock.Any())

	s.r.NoError(service.start())
}

func (s *Suite) initialContainerCheck() {
	for _, containerName := range []string{
		config.DockerScannerContainerName,
		config.DockerJSONRPCProxyContainerName,
		config.DockerNatsContainerName,
		config.DockerIpfsContainerName,
	} {
		s.dockerClient.EXPECT().GetContainerByName(s.service.ctx, containerName).Return(&types.Container{ID: testGenericContainerID}, nil)
	}

	s.dockerClient.EXPECT().GetContainers(s.service.ctx).Return([]types.Container{
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
	}, nil)

	// service containers + 1 old agent
	expectedContainerCount := config.DockerSupervisorManagedContainers + 1
	for i := 0; i < expectedContainerCount; i++ {
		s.dockerClient.EXPECT().RemoveContainer(s.service.ctx, testGenericContainerID).Return(nil)
		s.dockerClient.EXPECT().WaitContainerPrune(s.service.ctx, testGenericContainerID).Return(nil)
	}
	// expected container count + 1 old network
	expectedNetworkCount := expectedContainerCount + 1
	for i := 0; i < expectedNetworkCount; i++ {
		s.dockerClient.EXPECT().RemoveNetworkByName(s.service.ctx, gomock.Any()).Return(nil)
	}
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
	// Starts the agent admin container, sets networking rules, attaches the agent container networking
	// to admin container's interfaces, publishes a "running" message.

	s.agentImageClient.EXPECT().EnsureLocalImage(s.service.ctx, "agent test-agent", agentConfig.Image).Return(nil)

	s.dockerClient.EXPECT().StartContainer(s.service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name:      agentConfig.AdminContainerName(),
		NetworkID: testNodeNetworkID,
	}), true).Return(&clients.DockerContainer{Name: agentConfig.AdminContainerName(), ID: testAgentAdminContainerName}, nil)

	s.botManager.EXPECT().SetBotAdminRules(agentConfig.AdminContainerName())

	s.dockerClient.EXPECT().StartContainer(s.service.ctx, (configMatcher)(clients.DockerContainerConfig{
		Name:      agentConfig.ContainerName(),
		NetworkID: fmt.Sprintf("container:%s", agentConfig.AdminContainerName()),
	})).Return(&clients.DockerContainer{Name: agentConfig.ContainerName(), ID: testAgentContainerID}, nil)

	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsStatusRunning, agentPayload)

	s.r.NoError(s.service.handleAgentRun(agentPayload))
}

// TestAgentRunAgain tests running an agent twice.
func (s *Suite) TestAgentRunAgain() {
	s.TestAgentRun()

	agentConfig, agentPayload := testAgentData()

	// Expect it to only publish a message again to ensure the subscribers that
	// the agent is running.
	s.agentImageClient.EXPECT().EnsureLocalImage(s.service.ctx, "agent test-agent", agentConfig.Image).Return(nil)
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
