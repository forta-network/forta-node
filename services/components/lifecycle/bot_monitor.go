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

func (bm *botMonitor) saveBotActivity(botID string) {
	for _, tracker := range bm.trackers {
		if tracker.BotID() == botID {
			tracker.SaveActivity()
			return
		}
	}
}

// GetInactiveBots returns the list of the inactive bot IDs.
func (bm *botMonitor) GetInactiveBots() (inactive []string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	var preservedTrackers []*BotTracker
	for _, tracker := range bm.trackers {
		switch tracker.Status() {
		case TrackerStatusActive:
			preservedTrackers = append(preservedTrackers, tracker)

		case TrackerStatusInactive:
			preservedTrackers = append(preservedTrackers, tracker)
			inactive = append(inactive, tracker.BotID())

		case TrackerStatusStale:
			// ignore so the tracker is dropped
		}
	}

	if len(inactive) > 0 {
		bm.lifecycleMetrics.StatusInactive(inactive)
	}

	bm.trackers = preservedTrackers
	return
}
