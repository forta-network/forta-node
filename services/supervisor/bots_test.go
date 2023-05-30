package supervisor

import (
	"errors"
	"testing"

	mock_lifecycle "github.com/forta-network/forta-node/services/components/lifecycle/mocks"
	"github.com/golang/mock/gomock"
)

func TestBotManagement(t *testing.T) {
	supervisor := &SupervisorService{}

	botManager := mock_lifecycle.NewMockBotLifecycleManager(gomock.NewController(t))
	supervisor.botLifecycle.BotManager = botManager

	testErr := errors.New("test error - ignore")

	// both methods should be executed in order and the errors should not short circuit
	gomock.InOrder(
		botManager.EXPECT().ManageBots(gomock.Any()).Return(testErr),
		botManager.EXPECT().RestartExitedBots(gomock.Any()).Return(testErr),
		botManager.EXPECT().ExitInactiveBots(gomock.Any()).Return(testErr),
		botManager.EXPECT().ManageBots(gomock.Any()).Return(testErr),
		botManager.EXPECT().RestartExitedBots(gomock.Any()).Return(testErr),
		botManager.EXPECT().ExitInactiveBots(gomock.Any()).Return(testErr),
	)

	supervisor.doRefreshBotContainers()
	supervisor.doRefreshBotContainers()
}
