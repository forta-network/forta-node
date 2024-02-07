package lifecycle

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/clients/agentlogs"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/clients/docker"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	mock_containers "github.com/forta-network/forta-node/services/components/containers/mocks"
	mock_registry "github.com/forta-network/forta-node/services/components/registry/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSendBotLogsSuite(t *testing.T) {
	suite.Run(t, &BotLoggerSuite{})
}

type BotLoggerSuite struct {
	r *require.Assertions

	botClient     *mock_containers.MockBotClient
	dockerClient  *mock_clients.MockDockerClient
	agentRegistry *mock_registry.MockBotRegistry
	key           *keystore.Key
	suite.Suite
}

func (s *BotLoggerSuite) SetupTest() {
	t := s.T()
	ctrl := gomock.NewController(s.T())
	r := s.Require()

	botClient := mock_containers.NewMockBotClient(ctrl)
	dockerClient := mock_clients.NewMockDockerClient(ctrl)
	agentRegistry := mock_registry.NewMockBotRegistry(ctrl)

	dir := t.TempDir()
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := ks.NewAccount("Forta123")
	r.NoError(err)

	key, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	r.NoError(err)

	s.botClient = botClient
	s.dockerClient = dockerClient
	s.agentRegistry = agentRegistry
	s.key = key
	s.r = r
}

func (s *BotLoggerSuite) TestSendBotLogs() {
	botLogger := NewBotLogger(
		s.botClient, s.dockerClient, s.agentRegistry, s.key,
		func(agents agentlogs.Agents, authToken string) error {
			s.r.Equal(2, len(agents))

			s.r.Equal("bot1ID", agents[0].ID)
			s.r.EqualValues(1, agents[0].ChainID)

			s.r.Equal("bot2ID", agents[1].ID)
			s.r.EqualValues(2, agents[1].ChainID)
			return nil
		},
	)
	ctx := context.Background()

	mockContainers := []types.Container{
		{
			ID:    "bot1",
			Image: "forta/bot:latest",
			Labels: map[string]string{
				docker.LabelFortaSettingsAgentLogsEnable: "true",
				docker.LabelFortaBotID:                   "bot1ID",
			},
		},
		{
			ID:    "bot2",
			Image: "forta/bot:latest",
			Labels: map[string]string{
				docker.LabelFortaSettingsAgentLogsEnable: "true",
				docker.LabelFortaBotID:                   "bot2ID",
			},
		},
	}
	s.dockerClient.EXPECT().GetContainerLogs(
		ctx, "bot1",
		"60s",
		defaultAgentLogTailLines,
	).Return("some log", nil).Times(1)

	s.dockerClient.EXPECT().GetContainerLogs(
		ctx, "bot2",
		"60s",
		defaultAgentLogTailLines,
	).Return("some log", nil).Times(1)

	s.agentRegistry.EXPECT().GetConfigByID("bot1ID").Return(&config.AgentConfig{ChainID: 1}, nil).Times(1)
	s.agentRegistry.EXPECT().GetConfigByID("bot2ID").Return(&config.AgentConfig{ChainID: 2}, nil).Times(1)

	s.botClient.EXPECT().LoadBotContainers(ctx).Return(mockContainers, nil)
	s.r.NoError(botLogger.SendBotLogs(ctx, time.Minute))
}

// should fail if there is an error loading
// bot containers
func (s *BotLoggerSuite) TestLoadBotContainersError() {
	botLogger := NewBotLogger(
		s.botClient, s.dockerClient, s.agentRegistry, s.key,
		func(agents agentlogs.Agents, authToken string) error {
			return nil
		},
	)
	ctx := context.Background()

	mockContainers := []types.Container{}

	s.botClient.EXPECT().LoadBotContainers(ctx).Return(mockContainers, errors.New("test"))
	s.r.EqualError(botLogger.SendBotLogs(ctx, time.Minute), "failed to load the bot containers: test")
}

// Should not send agent logs if fails
// to get container logs but continue processing
func (s *BotLoggerSuite) TestGetContainerLogsError() {
	botLogger := NewBotLogger(
		s.botClient, s.dockerClient, s.agentRegistry, s.key,
		func(agents agentlogs.Agents, authToken string) error {
			s.r.Equal(1, len(agents))
			s.r.Equal("bot2ID", agents[0].ID)
			s.r.Equal("some log", agents[0].Logs)
			return nil
		},
	)
	ctx := context.Background()

	mockContainers := []types.Container{
		{
			ID:    "bot1",
			Image: "forta/bot:latest",
			Labels: map[string]string{
				docker.LabelFortaSettingsAgentLogsEnable: "true",
			},
		},
		{
			ID:    "bot2",
			Image: "forta/bot:latest",
			Labels: map[string]string{
				docker.LabelFortaSettingsAgentLogsEnable: "true",
				docker.LabelFortaBotID:                   "bot2ID",
			},
		},
	}

	s.botClient.EXPECT().LoadBotContainers(ctx).Return(mockContainers, nil)

	s.dockerClient.EXPECT().GetContainerLogs(
		ctx, "bot1",
		"60s",
		defaultAgentLogTailLines,
	).Return("", errors.New("test")).Times(1)

	s.dockerClient.EXPECT().GetContainerLogs(
		ctx, "bot2",
		"60s",
		defaultAgentLogTailLines,
	).Return("some log", nil).Times(1)

	s.agentRegistry.EXPECT().GetConfigByID("bot2ID").Return(&config.AgentConfig{ChainID: 2}, nil).Times(1)

	s.r.NoError(botLogger.SendBotLogs(ctx, time.Minute))
}

// Fails sending agent logs
func (s *BotLoggerSuite) TestFailsToSendLogs() {
	botLogger := NewBotLogger(
		s.botClient, s.dockerClient, s.agentRegistry, s.key,
		func(agents agentlogs.Agents, authToken string) error {
			return errors.New("test")
		},
	)
	ctx := context.Background()

	mockContainers := []types.Container{
		{
			ID:    "bot1",
			Image: "forta/bot:latest",
			Labels: map[string]string{
				docker.LabelFortaSettingsAgentLogsEnable: "true",
				docker.LabelFortaBotID:                   "bot1ID",
			},
		},
	}

	s.botClient.EXPECT().LoadBotContainers(ctx).Return(mockContainers, nil)

	s.dockerClient.EXPECT().GetContainerLogs(
		ctx, "bot1",
		"60s",
		defaultAgentLogTailLines,
	).Return("some log", nil).Times(1)

	s.agentRegistry.EXPECT().GetConfigByID("bot1ID").Return(&config.AgentConfig{ChainID: 1}, nil).Times(1)

	s.r.EqualError(botLogger.SendBotLogs(ctx, time.Minute), "failed to send agent logs: test")
}
