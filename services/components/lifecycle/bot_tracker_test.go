package lifecycle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testBotID = "test-bot-id"
)

func TestActive(t *testing.T) {
	r := require.New(t)

	botTracker := NewBotTracker(testBotID)
	r.Equal(TrackerStatusActive, botTracker.Status())
}

func TestActive_ShortCircuit(t *testing.T) {
	r := require.New(t)

	botTracker := NewBotTracker(testBotID)
	r.Equal(TrackerStatusActive, botTracker.Status())
	r.Equal(TrackerStatusActive, botTracker.Status())
}

func TestInactive(t *testing.T) {
	r := require.New(t)

	botTracker := NewBotTracker(testBotID)
	botTracker.lastActivity = time.Now().Add(-inactivityThreshold - 1)
	r.Equal(TrackerStatusInactive, botTracker.Status())

	// should say "active" for the second time to avoid quick reads
	r.Equal(TrackerStatusActive, botTracker.Status())
}

func TestStale(t *testing.T) {
	r := require.New(t)

	botTracker := NewBotTracker(testBotID)
	botTracker.lastActivity = time.Now().Add(-expiryThreshold - 1)
	r.Equal(TrackerStatusStale, botTracker.Status())

	// should say "active" for the second time to avoid quick reads
	r.Equal(TrackerStatusActive, botTracker.Status())
}

func TestSaveActivity(t *testing.T) {
	r := require.New(t)

	botTracker := NewBotTracker(testBotID)
	botTracker.lastActivity = time.Now().Add(-inactivityThreshold - 1)
	botTracker.SaveActivity()

	// the status is "active" because SaveActivity() updated the time
	r.Equal(TrackerStatusActive, botTracker.Status())
}

func TestGetBotID(t *testing.T) {
	r := require.New(t)

	botTracker := NewBotTracker(testBotID)
	r.Equal(testBotID, botTracker.BotID())
}
