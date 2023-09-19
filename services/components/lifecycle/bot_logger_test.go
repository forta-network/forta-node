package lifecycle

import (
	"testing"

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

	botLogger *botLogger
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

	sendLogsMockFn := func(agents agentlogs.Agents, authToken string) error {
		return nil
	}

	botLogger := NewBotLogger(botClient, dockerClient, key, sendLogsMockFn)

	s.botLogger = botLogger
	s.r = r
}

func (s *BotLoggerSuite) TestSendBotLogs() {
	s.r.NotNil(s.botLogger)
}
