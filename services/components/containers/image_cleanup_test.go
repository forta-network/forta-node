package containers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/forta-network/forta-node/clients/docker"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	mock_registry "github.com/forta-network/forta-node/services/components/registry/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testCleanupImage1     = "bafybei-testcleanupimage1"
	testCleanupImage2     = "bafybei-testcleanupimage2"
	testCleanupImage3     = "bafybei-testcleanupimage3"
	testHeartbeatBotImage = "bafybei-heartbeatbot"
)

type ImageCleanupTestSuite struct {
	r *require.Assertions

	client      *mock_clients.MockDockerClient
	botRegistry *mock_registry.MockBotRegistry

	imageCleanup *imageCleanup

	suite.Suite
}

func TestImageCleanupTestSuite(t *testing.T) {
	suite.Run(t, &ImageCleanupTestSuite{})
}

func (s *ImageCleanupTestSuite) SetupTest() {
	logrus.SetLevel(logrus.DebugLevel)

	s.r = s.Require()

	ctrl := gomock.NewController(s.T())
	s.client = mock_clients.NewMockDockerClient(ctrl)
	s.botRegistry = mock_registry.NewMockBotRegistry(ctrl)

	s.imageCleanup = NewImageCleanup(s.client, s.botRegistry)
}

func (s *ImageCleanupTestSuite) TestIntervalSkip() {
	s.imageCleanup.lastCleanup = time.Now() // very close cleanup time

	// no calls expected

	s.r.NoError(s.imageCleanup.Do(context.Background()))
}

func (s *ImageCleanupTestSuite) TestContainerListError() {
	s.client.EXPECT().GetContainers(gomock.Any()).Return(nil, errors.New("test error"))

	s.r.Error(s.imageCleanup.Do(context.Background()))
}

func (s *ImageCleanupTestSuite) TestImagesListError() {
	s.client.EXPECT().GetContainers(gomock.Any()).Return(docker.ContainerList{}, nil)
	s.client.EXPECT().ListDigestReferences(gomock.Any()).Return(nil, errors.New("test error"))

	s.r.Error(s.imageCleanup.Do(context.Background()))
}

func (s *ImageCleanupTestSuite) TestHeartbeatBotError() {
	s.client.EXPECT().GetContainers(gomock.Any()).Return(docker.ContainerList{}, nil)
	s.client.EXPECT().ListDigestReferences(gomock.Any()).Return([]string{testCleanupImage1}, nil)
	s.botRegistry.EXPECT().LoadHeartbeatBot().Return(nil, errors.New("test error"))

	s.r.Error(s.imageCleanup.Do(context.Background()))
}

func (s *ImageCleanupTestSuite) TestRemoveImageError() {
	initialLastCleanup := s.imageCleanup.lastCleanup

	s.client.EXPECT().GetContainers(gomock.Any()).Return(docker.ContainerList{}, nil)
	s.client.EXPECT().ListDigestReferences(gomock.Any()).Return([]string{testCleanupImage1}, nil)
	s.botRegistry.EXPECT().LoadHeartbeatBot().Return(&config.AgentConfig{Image: testHeartbeatBotImage}, nil)
	s.client.EXPECT().RemoveImage(gomock.Any(), testCleanupImage1).Return(errors.New("test error"))

	// no error and mutated last cleanup timestamp: removal errors do not affect this
	s.r.NoError(s.imageCleanup.Do(context.Background()))
	s.r.NotEqual(initialLastCleanup, s.imageCleanup.lastCleanup)
}

func (s *ImageCleanupTestSuite) TestCleanupSuccess() {
	initialLastCleanup := s.imageCleanup.lastCleanup

	s.imageCleanup.exclusionList = []string{testCleanupImage1} // excluding image
	s.client.EXPECT().GetContainers(gomock.Any()).Return(docker.ContainerList{
		{
			Image: testCleanupImage2, // image in use by container
		},
	}, nil)
	s.client.EXPECT().ListDigestReferences(gomock.Any()).Return(
		[]string{testCleanupImage1, testCleanupImage2, testCleanupImage3, testHeartbeatBotImage}, nil,
	)
	s.botRegistry.EXPECT().LoadHeartbeatBot().Return(&config.AgentConfig{Image: testHeartbeatBotImage}, nil)
	s.client.EXPECT().RemoveImage(gomock.Any(), testCleanupImage3).Return(nil) // only removes image 3

	s.r.NoError(s.imageCleanup.Do(context.Background()))
	s.r.NotEqual(initialLastCleanup, s.imageCleanup.lastCleanup)
}
