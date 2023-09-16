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
)

func sendLogsMock(agents agentlogs.Agents, authToken string) error {
	return nil
}

func TestNewBotLogger(t *testing.T) {
	ctrl := gomock.NewController(t)

	r := require.New(t)

	botClient := mock_containers.NewMockBotClient(ctrl)
	dockerClient := mock_clients.NewMockDockerClient(ctrl)

	dir := t.TempDir()
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := ks.NewAccount("Forta123")
	r.NoError(err)

	key, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	r.NoError(err)

	botLogger := NewBotLogger(botClient, dockerClient, key, sendLogsMock)
	r.NotNil(botLogger)
}
