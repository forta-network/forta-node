package supervisor

import (
	"errors"
	"testing"

	mock_containers "github.com/forta-network/forta-node/services/components/containers/mocks"
	mock_lifecycle "github.com/forta-network/forta-node/services/components/lifecycle/mocks"
	"github.com/golang/mock/gomock"
)

func TestBotManagement(t *testing.T) {
	supervisor := &SupervisorService{}

	ctrl := gomock.NewController(t)
	botManager := mock_lifecycle.NewMockBotLifecycleManager(ctrl)
	imageCleanup := mock_containers.NewMockImageCleanup(ctrl)
	supervisor.botLifecycle.BotManager = botManager
	supervisor.botLifecycle.ImageCleanup = imageCleanup

	testErr := errors.New("test error - ignore")

	// both methods should be executed in order and the errors should not short circuit
	gomock.InOrder(
		botManager.EXPECT().ManageBots(gomock.Any()).Return(testErr),
		botManager.EXPECT().CleanupUnusedBots(gomock.Any()).Return(testErr),
		imageCleanup.EXPECT().Do(gomock.Any()).Return(testErr),
		botManager.EXPECT().RestartExitedBots(gomock.Any()).Return(testErr),
		botManager.EXPECT().ExitInactiveBots(gomock.Any()).Return(testErr),
		botManager.EXPECT().ManageBots(gomock.Any()).Return(testErr),
		botManager.EXPECT().CleanupUnusedBots(gomock.Any()).Return(testErr),
		imageCleanup.EXPECT().Do(gomock.Any()).Return(testErr),
		botManager.EXPECT().RestartExitedBots(gomock.Any()).Return(testErr),
		botManager.EXPECT().ExitInactiveBots(gomock.Any()).Return(testErr),
	)

	supervisor.doRefreshBotContainers()
	supervisor.doRefreshBotContainers()
}
