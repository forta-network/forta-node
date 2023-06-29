package lifecycle

import (
	"sync"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/services/components/metrics"
)

// BotMonitorUpdater updates the bot monitor.
type BotMonitorUpdater interface {
	UpdateWithMetrics(*protocol.AgentMetricList) error
}

// BotMonitorState reads the bot monitor state.
type BotMonitorState interface {
	MonitorBots([]string)
	GetInactiveBots() []string
}

// BotMonitor monitors the statuses of the bots using the incoming metrics.
type BotMonitor interface {
	BotMonitorUpdater
	BotMonitorState
}

type botMonitor struct {
	lifecycleMetrics metrics.Lifecycle
	trackers         []*BotTracker
	mu               sync.Mutex
}

var _ BotMonitor = &botMonitor{}

// NewBotMonitor creates a new bot monitor.
func NewBotMonitor(lifecycleMetrics metrics.Lifecycle) *botMonitor {
	return &botMonitor{
		lifecycleMetrics: lifecycleMetrics,
	}
}

// UpdateWithMetrics updates the trackers with metrics.
func (bm *botMonitor) UpdateWithMetrics(botMetrics *protocol.AgentMetricList) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if botMetrics == nil {
		return nil
	}

	for _, botMetric := range botMetrics.Metrics {
		if botMetric.Name == metrics.MetricStatusActive {
			bm.saveBotActivity(botMetric.AgentId)
		}
	}

	return nil
}

func (bm *botMonitor) findTrackerAndDo(botID string, do func(*BotTracker)) {
	for _, tracker := range bm.trackers {
		if tracker.BotID() == botID {
			do(tracker)
			return
		}
	}
}

func (bm *botMonitor) missTrackerAndDo(botID string, do func()) {
	for _, tracker := range bm.trackers {
		if tracker.BotID() == botID {
			return
		}
	}
	do()
}

func (bm *botMonitor) saveBotActivity(botID string) {
	bm.findTrackerAndDo(botID, func(tracker *BotTracker) {
		tracker.SaveActivity()
	})
}

func (bm *botMonitor) ensureTrackerExists(botID string) {
	bm.missTrackerAndDo(botID, func() {
		bm.trackers = append(bm.trackers, NewBotTracker(botID))
	})
}

func (bm *botMonitor) dropStaleTrackers(botIDs []string) {
	var preservedTrackers []*BotTracker
	for _, tracker := range bm.trackers {
		for _, botID := range botIDs {
			if tracker.BotID() == botID {
				preservedTrackers = append(preservedTrackers, tracker)
				break
			}
		}
	}
	bm.trackers = preservedTrackers
}

// MonitorBots makes sure that the bots with given IDs are monitored.
func (bm *botMonitor) MonitorBots(botIDs []string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	for _, botID := range botIDs {
		bm.ensureTrackerExists(botID)
	}
	bm.dropStaleTrackers(botIDs)
}

// GetInactiveBots returns the list of the inactive bot IDs.
func (bm *botMonitor) GetInactiveBots() (inactive []string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	for _, tracker := range bm.trackers {
		if tracker.IsInactive() {
			inactive = append(inactive, tracker.BotID())
		}
	}

	return
}
