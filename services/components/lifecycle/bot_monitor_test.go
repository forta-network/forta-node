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
	testTrackerBotID4 = "test-tracker-bot-id-4"
	testTrackerBotID5 = "test-tracker-bot-id-5"
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
			lastActivity: time.Now().Add(-inactivityThreshold * 2),
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

	// the latest inactivity list is requested: bot 2 is active, the rest are inactive
	inactiveBots := botMonitor.GetInactiveBots()
	r.Len(inactiveBots, 2)
	r.Equal(testTrackerBotID1, inactiveBots[0])
	r.Equal(testTrackerBotID3, inactiveBots[1])
	r.Len(botMonitor.trackers, 3)

	// a second request to get the latest inactivity will hit the read cooldown
	// which is helpful in not detecting inactive bots too often but every once in a while
	inactiveBots = botMonitor.GetInactiveBots()
	r.Len(inactiveBots, 0)
	r.Len(botMonitor.trackers, 3) // should keep the previous trackers

	// it should add a new tracker and drop a stale one (drops bot 3, adds bot 4)
	botMonitor.MonitorBots([]string{testTrackerBotID1, testTrackerBotID2, testTrackerBotID4})
	r.Len(botMonitor.trackers, 3)
	r.Equal(testTrackerBotID1, botMonitor.trackers[0].BotID())
	r.Equal(testTrackerBotID2, botMonitor.trackers[1].BotID())
	r.Equal(testTrackerBotID4, botMonitor.trackers[2].BotID())

	// and should not add a duplicate one
	botMonitor.MonitorBots([]string{testTrackerBotID1, testTrackerBotID2, testTrackerBotID4})
	r.Len(botMonitor.trackers, 3)
	r.Equal(testTrackerBotID1, botMonitor.trackers[0].BotID())
	r.Equal(testTrackerBotID2, botMonitor.trackers[1].BotID())
	r.Equal(testTrackerBotID4, botMonitor.trackers[2].BotID())

	// and should be able to add an even newer one
	botMonitor.MonitorBots([]string{testTrackerBotID1, testTrackerBotID2, testTrackerBotID4, testTrackerBotID5})
	r.Len(botMonitor.trackers, 4)
	r.Equal(testTrackerBotID1, botMonitor.trackers[0].BotID())
	r.Equal(testTrackerBotID2, botMonitor.trackers[1].BotID())
	r.Equal(testTrackerBotID4, botMonitor.trackers[2].BotID())
	r.Equal(testTrackerBotID5, botMonitor.trackers[3].BotID())
}
