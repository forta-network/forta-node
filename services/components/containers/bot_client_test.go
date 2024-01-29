package containers

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/errdefs"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/clients/messaging"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testBotID1       = "0x0100000000000000000000000000000000000000000000000000000000000000"
	testBotID2       = "0x0200000000000000000000000000000000000000000000000000000000000000"
	testBotID3       = "0x0300000000000000000000000000000000000000000000000000000000000000"
	testImageRef     = "bafybeielvnt5apaxbk6chthc4dc3p6vscpx3ai4uvti7gwh253j7facsxu@sha256:e0e9efb6699b02750f6a9668084d37314f1de3a80da7e19c1d40da73ee57dd45"
	testContainerID  = "test-container-id"
	testContainerID1 = "test-container-id-1"
	testContainerID2 = "test-container-id-2"
	testBotNetworkID = "test-bot-network-id"
)

type BotClientTestSuite struct {
	r *require.Assertions

	client         *mock_clients.MockDockerClient
	botImageClient *mock_clients.MockDockerClient
	msgClient      *mock_clients.MockMessageClient

	botClient *botClient

	suite.Suite
}

func TestBotClientTestSuite(t *testing.T) {
	suite.Run(t, &BotClientTestSuite{})
}

func (s *BotClientTestSuite) SetupTest() {
	s.r = s.Require()

	ctrl := gomock.NewController(s.T())
	s.client = mock_clients.NewMockDockerClient(ctrl)
	s.botImageClient = mock_clients.NewMockDockerClient(ctrl)
	s.msgClient = mock_clients.NewMockMessageClient(ctrl)

	s.botImageClient.EXPECT().SetImagePullCooldown(ImagePullCooldownThreshold, ImagePullCooldownDuration)

	s.botClient = NewBotClient(config.LogConfig{}, config.ResourcesConfig{}, "", s.client, s.botImageClient, s.msgClient)
}

func (s *BotClientTestSuite) TestEnsureBotImages() {
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
	expectedImagePulls := []docker.ImagePull{
		{
			Name: testBotID1,
			Ref:  testImageRef,
		},
		{
			Name: testBotID2,
			Ref:  testImageRef,
		},
	}

	retErrs := []error{
		errors.New("err1"), errors.New("err2"),
	}
	s.botImageClient.EXPECT().EnsureLocalImages(gomock.Any(), BotPullTimeout, expectedImagePulls).Return(retErrs)

	s.r.Equal(retErrs, s.botClient.EnsureBotImages(context.Background(), botConfigs))
}

func (s *BotClientTestSuite) TestLaunchBot_Exists() {
	botConfig := config.AgentConfig{
		ID:    testBotID1,
		Image: testImageRef,
	}

	s.client.EXPECT().EnsurePublicNetwork(gomock.Any(), botConfig.ContainerName()).Return(testBotNetworkID, nil)
	s.client.EXPECT().GetContainerByName(gomock.Any(), botConfig.ContainerName()).Return(nil, nil)
	for _, serviceContainerName := range getServiceContainerNames() {
		s.client.EXPECT().GetContainerByName(gomock.Any(), serviceContainerName).Return(&types.Container{
			ID: testContainerID,
		}, nil)
		s.client.EXPECT().AttachNetwork(gomock.Any(), testContainerID, testBotNetworkID).Return(nil)
	}

	resources := &docker.ContainerResources{
		CPUStats: docker.CPUStats{
			CPUUsage: docker.CPUUsage{
				TotalUsage: 33,
			},
		},
		MemoryStats: docker.MemoryStats{
			Usage: 100,
		},
		NetworkStats: map[string]docker.NetworkStats{
			"eth0": {
				RxBytes: 123,
				TxBytes: 456,
			},
		},
	}
	executed := make(chan bool)
	s.client.EXPECT().ContainerStats(gomock.Any(), botConfig.ContainerName()).Return(resources, nil).Times(1)
	s.msgClient.EXPECT().PublishProto(messaging.SubjectMetricAgent, gomock.Any()).Do(func(v1, v2 interface{}) {
		metrics := v2.(*protocol.AgentMetricList)
		assert.Len(s.T(), metrics.Metrics, 4)

		// CPU metric
		assert.Equal(s.T(), botConfig.ID, metrics.Metrics[0].AgentId)
		assert.Equal(s.T(), domain.MetricDockerResourcesCPU, metrics.Metrics[0].Name)
		assert.Equal(s.T(), float64(33), metrics.Metrics[0].Value)

		// Memory metric
		assert.Equal(s.T(), botConfig.ID, metrics.Metrics[1].AgentId)
		assert.Equal(s.T(), domain.MetricDockerResourcesMemory, metrics.Metrics[1].Name)
		assert.Equal(s.T(), float64(100), metrics.Metrics[1].Value)

		// Network bytes received metric
		assert.Equal(s.T(), botConfig.ID, metrics.Metrics[3].AgentId)
		assert.Equal(s.T(), domain.MetricDockerResourcesNetworkReceive, metrics.Metrics[3].Name)
		assert.Equal(s.T(), float64(123), metrics.Metrics[3].Value)

		// Network bytes sent metric
		assert.Equal(s.T(), botConfig.ID, metrics.Metrics[2].AgentId)
		assert.Equal(s.T(), domain.MetricDockerResourcesNetworkSent, metrics.Metrics[2].Name)
		assert.Equal(s.T(), float64(456), metrics.Metrics[2].Value)

		close(executed)
	})

	s.r.NoError(s.botClient.LaunchBot(context.Background(), botConfig))
	<-executed
}

func (s *BotClientTestSuite) TestLaunchBot_GetContainerError() {
	botConfig := config.AgentConfig{
		ID:    testBotID1,
		Image: testImageRef,
	}

	s.client.EXPECT().EnsurePublicNetwork(gomock.Any(), botConfig.ContainerName()).Return(testBotNetworkID, nil)
	err := errors.New("unexpected error")
	s.client.EXPECT().GetContainerByName(gomock.Any(), botConfig.ContainerName()).Return(nil, err)

	s.r.ErrorContains(s.botClient.LaunchBot(context.Background(), botConfig), err.Error())
}

func (s *BotClientTestSuite) TestLaunchBot_DoesNotExist() {
	botConfig := config.AgentConfig{
		ID:    testBotID1,
		Image: testImageRef,
	}

	s.client.EXPECT().EnsurePublicNetwork(gomock.Any(), botConfig.ContainerName()).Return(testBotNetworkID, nil)
	s.client.EXPECT().GetContainerByName(gomock.Any(), botConfig.ContainerName()).Return(nil, docker.ErrContainerNotFound)
	botContainerCfg := NewBotContainerConfig(testBotNetworkID, botConfig, config.LogConfig{}, config.ResourcesConfig{}, "")
	s.client.EXPECT().StartContainer(gomock.Any(), botContainerCfg).Return(nil, nil)
	for _, serviceContainerName := range getServiceContainerNames() {
		s.client.EXPECT().GetContainerByName(gomock.Any(), serviceContainerName).Return(&types.Container{
			ID: testContainerID,
		}, nil)
		s.client.EXPECT().AttachNetwork(gomock.Any(), testContainerID, testBotNetworkID).Return(nil)
	}

	s.client.EXPECT().ContainerStats(gomock.Any(), botConfig.ContainerName()).Return(nil, errdefs.NotFound(errors.New(""))).MaxTimes(1)

	s.r.NoError(s.botClient.LaunchBot(context.Background(), botConfig))
}

func (s *BotClientTestSuite) TestTearDownBot() {
	botConfig := config.AgentConfig{
		ID:    testBotID1,
		Image: testImageRef,
	}

	testErr := errors.New("some error")
	s.client.EXPECT().GetContainerByName(gomock.Any(), botConfig.ContainerName()).Return(&types.Container{
		ID:    testContainerID2,
		Image: testImageRef,
	}, nil)
	for _, serviceContainerName := range getServiceContainerNames() {
		s.client.EXPECT().GetContainerByName(gomock.Any(), serviceContainerName).Return(&types.Container{
			ID: testContainerID,
		}, nil)
		s.client.EXPECT().DetachNetwork(gomock.Any(), testContainerID, botConfig.ContainerName()).Return(testErr)
	}
	timeout := BotShutdownTimeout
	s.client.EXPECT().ShutdownContainer(gomock.Any(), testContainerID2, &timeout).Return(testErr)
	s.client.EXPECT().RemoveContainer(gomock.Any(), testContainerID2).Return(testErr)
	s.client.EXPECT().RemoveNetworkByName(gomock.Any(), botConfig.ContainerName()).Return(testErr)

	s.r.NoError(s.botClient.TearDownBot(context.Background(), botConfig.ContainerName()))
}

func (s *BotClientTestSuite) TestStopBot() {
	botConfig := config.AgentConfig{
		ID:    testBotID1,
		Image: testImageRef,
	}

	s.client.EXPECT().GetContainerByName(gomock.Any(), botConfig.ContainerName()).Return(&types.Container{
		ID:    testContainerID2,
		Image: testImageRef,
	}, nil)
	s.client.EXPECT().StopContainer(gomock.Any(), testContainerID2)

	s.r.NoError(s.botClient.StopBot(context.Background(), botConfig))
}

func (s *BotClientTestSuite) TestLoadBotContainers() {
	expectedContainers := docker.ContainerList{{}}
	s.client.EXPECT().GetContainersByLabel(gomock.Any(), docker.LabelFortaIsBot, LabelValueFortaIsBot).Return(expectedContainers, nil)

	containers, err := s.botClient.LoadBotContainers(context.Background())
	s.r.NoError(err)
	s.r.Equal(([]types.Container)(expectedContainers), containers)
}

func (s *BotClientTestSuite) TestStartWaitBotContainer() {
	s.client.EXPECT().StartContainerWithID(gomock.Any(), testContainerID).Return(nil)
	s.client.EXPECT().WaitContainerStart(gomock.Any(), testContainerID).Return(nil)

	s.r.NoError(s.botClient.StartWaitBotContainer(context.Background(), testContainerID))
}
