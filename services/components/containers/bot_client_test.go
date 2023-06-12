package containers

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/forta-network/forta-node/clients/docker"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/golang/mock/gomock"
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

	s.botImageClient.EXPECT().SetImagePullCooldown(ImagePullCooldownThreshold, ImagePullCooldownDuration)

	s.botClient = NewBotClient(config.LogConfig{}, config.ResourcesConfig{}, s.client, s.botImageClient)
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

func (s *BotClientTestSuite) TestEnsureBotImage() {
	botConfig := config.AgentConfig{
		ID:    testBotID1,
		Image: testImageRef,
	}

	err := errors.New("err1")

	s.botImageClient.EXPECT().EnsureLocalImage(gomock.Any(), botConfig.ID, botConfig.Image).Return(err)

	s.r.Error(err, s.botClient.EnsureBotImage(context.Background(), botConfig))
}

func (s *BotClientTestSuite) TestLaunchBot_Exists() {
	botConfig := config.AgentConfig{
		ID:    testBotID1,
		Image: testImageRef,
	}

	s.client.EXPECT().GetContainerByName(gomock.Any(), botConfig.ContainerName()).Return(nil, nil)

	s.r.NoError(s.botClient.LaunchBot(context.Background(), botConfig))
}

func (s *BotClientTestSuite) TestLaunchBot_GetContainerError() {
	botConfig := config.AgentConfig{
		ID:    testBotID1,
		Image: testImageRef,
	}

	err := errors.New("unexpected error")
	s.client.EXPECT().GetContainerByName(gomock.Any(), botConfig.ContainerName()).Return(nil, err)

	s.r.ErrorContains(s.botClient.LaunchBot(context.Background(), botConfig), err.Error())
}

func (s *BotClientTestSuite) TestLaunchBot_DoesNotExist() {
	botConfig := config.AgentConfig{
		ID:    testBotID1,
		Image: testImageRef,
	}

	s.client.EXPECT().GetContainerByName(gomock.Any(), botConfig.ContainerName()).Return(nil, docker.ErrContainerNotFound)
	s.client.EXPECT().CreatePublicNetwork(gomock.Any(), botConfig.ContainerName()).Return(testBotNetworkID, nil)
	botContainerCfg := NewBotContainerConfig(testBotNetworkID, botConfig, config.LogConfig{}, config.ResourcesConfig{})
	s.client.EXPECT().StartContainer(gomock.Any(), botContainerCfg).Return(nil, nil)
	for _, serviceContainerName := range getServiceContainerNames() {
		s.client.EXPECT().GetContainerByName(gomock.Any(), serviceContainerName).Return(&types.Container{
			ID: testContainerID,
		}, nil)
		s.client.EXPECT().AttachNetwork(gomock.Any(), testContainerID, testBotNetworkID).Return(nil)
	}

	s.r.NoError(s.botClient.LaunchBot(context.Background(), botConfig))
}

func (s *BotClientTestSuite) TestTearDownBot() {
	botConfig := config.AgentConfig{
		ID:    testBotID1,
		Image: testImageRef,
	}

	s.client.EXPECT().GetContainerByName(gomock.Any(), botConfig.ContainerName()).Return(&types.Container{
		ID:    testContainerID2,
		Image: testImageRef,
	}, nil)
	for _, serviceContainerName := range getServiceContainerNames() {
		s.client.EXPECT().GetContainerByName(gomock.Any(), serviceContainerName).Return(&types.Container{
			ID: testContainerID,
		}, nil)
		s.client.EXPECT().DetachNetwork(gomock.Any(), testContainerID, botConfig.ContainerName()).Return(nil)
	}
	s.client.EXPECT().RemoveContainer(gomock.Any(), testContainerID2)
	s.client.EXPECT().RemoveNetworkByName(gomock.Any(), botConfig.ContainerName())
	s.client.EXPECT().RemoveImage(gomock.Any(), testImageRef)

	s.r.NoError(s.botClient.TearDownBot(context.Background(), botConfig.ContainerName(), true))
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
