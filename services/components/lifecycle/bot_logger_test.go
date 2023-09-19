package lifecycle

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/clients/agentlogs"
	"github.com/forta-network/forta-core-go/security"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	mock_containers "github.com/forta-network/forta-node/services/components/containers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSendBotLogsSuite(t *testing.T) {
	suite.Run(t, &BotLoggerSuite{})
}

type BotLoggerSuite struct {
	r *require.Assertions

	botLogger    *botLogger
	botClient    *mock_containers.MockBotClient
	dockerClient *mock_clients.MockDockerClient
	key          *keystore.Key
	suite.Suite
}

func (s *BotLoggerSuite) SetupTest() {
	t := s.T()
	ctrl := gomock.NewController(s.T())
	r := s.Require()

	botClient := mock_containers.NewMockBotClient(ctrl)
	dockerClient := mock_clients.NewMockDockerClient(ctrl)

	dir := t.TempDir()
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := ks.NewAccount("Forta123")
	r.NoError(err)

	key, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	r.NoError(err)

	s.botClient = botClient
	s.dockerClient = dockerClient
	s.key = key
	s.r = r
}

func (s *BotLoggerSuite) TestSendBotLogs() {
	botLogger := NewBotLogger(
		s.botClient, s.dockerClient, s.key,
		func(agents agentlogs.Agents, authToken string) error {
			s.r.Equal(2, len(agents))
			s.r.Equal("bot1", agents[0].ID)
			s.r.Equal("bot2", agents[1].ID)
			// TODO: test sendLogs = 1 and sendLogs = 0 conditions
			return nil
		},
	)
	ctx := context.Background()

	mockContainers := []types.Container{
		{
			ID:    "bot1",
			Image: "forta/bot:latest",
		},
		{
			ID:    "bot2",
			Image: "forta/bot:latest",
		},
	}
	s.botClient.EXPECT().LoadBotContainers(ctx).Return(mockContainers, nil)

	botLogger.SendBotLogs(ctx)
}
