package lifecycle

import (
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/services/components/metrics"
	mock_metrics "github.com/forta-network/forta-node/services/components/metrics/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	testTrackerBotID1 = "test-tracker-bot-id-1"
	testTrackerBotID2 = "test-tracker-bot-id-2"
	testTrackerBotID3 = "test-tracker-bot-id-3"
)

func TestBotMonitor(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	lifecycleMetrics := mock_metrics.NewMockLifecycle(ctrl)

	// we have a monitor with three trackers
	botMonitor := NewBotMonitor(lifecycleMetrics)
	botMonitor.trackers = []*BotTracker{
		{
			botID:        testTrackerBotID1,
			lastRead:     time.Time{}, // no reads yet
			lastActivity: time.Now().Add(-inactivityThreshold - 1),
		},
		{
			botID:        testTrackerBotID2,
			lastRead:     time.Time{}, // no reads yet
			lastActivity: time.Now().Add(-inactivityThreshold - 1),
		},
		{
			botID:        testTrackerBotID3,
			lastRead:     time.Time{}, // no reads yet
			lastActivity: time.Now().Add(-expiryThreshold - 1),
		},
	}

	// the bot monitor handles a nil payload struct
	r.NoError(botMonitor.UpdateWithMetrics(nil))
	// the bot monitor records the activity of the second bot
	r.NoError(botMonitor.UpdateWithMetrics(&protocol.AgentMetricList{
		Metrics: []*protocol.AgentMetric{
			{
				Name:    metrics.MetricStatusActive,
				AgentId: testTrackerBotID2,
			},
		},
	}))

	// the latest inactivity list is requested: bot 1 is inactive, bot 3 is stale (tracker dropped)
	lifecycleMetrics.EXPECT().StatusInactive([]string{testTrackerBotID1})
	inactiveBots := botMonitor.GetInactiveBots()
	r.Len(inactiveBots, 1)
	r.Equal(testTrackerBotID1, inactiveBots[0])
	r.Len(botMonitor.trackers, 2) // shrinked down because stale bot 3 tracker is dropped!

	// a second request to get the latest inactivity will hit the read cooldown
	// which is helpful in not detecting inactive bots too often but every once in a while
	inactiveBots = botMonitor.GetInactiveBots()
	r.Len(inactiveBots, 0)
	r.Len(botMonitor.trackers, 2) // should keep the previous trackers
}
