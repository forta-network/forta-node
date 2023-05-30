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
	r.Equal(false, botTracker.IsInactive())
}

func TestActive_ShortCircuit(t *testing.T) {
	r := require.New(t)

	botTracker := NewBotTracker(testBotID)
	r.Equal(false, botTracker.IsInactive())
	r.Equal(false, botTracker.IsInactive())
}

func TestInactive(t *testing.T) {
	r := require.New(t)

	botTracker := NewBotTracker(testBotID)
	botTracker.lastActivity = time.Now().Add(-inactivityThreshold - 1)
	r.Equal(true, botTracker.IsInactive())

	// should say "active" for the second time to avoid quick reads
	r.Equal(false, botTracker.IsInactive())
}

func TestSaveActivity(t *testing.T) {
	r := require.New(t)

	botTracker := NewBotTracker(testBotID)
	botTracker.lastActivity = time.Now().Add(-inactivityThreshold - 1)
	botTracker.SaveActivity()

	// the status is "active" because SaveActivity() updated the time
	r.Equal(false, botTracker.IsInactive())
}

func TestGetBotID(t *testing.T) {
	r := require.New(t)

	botTracker := NewBotTracker(testBotID)
	r.Equal(testBotID, botTracker.BotID())
}
