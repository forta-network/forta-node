package lifecycle

import (
	"time"
)

// Timeouts
const (
	readCooldown        = time.Minute * 5
	inactivityThreshold = time.Minute * 15
)

// BotTracker tracks activity time of a bot.
type BotTracker struct {
	botID        string
	lastActivity time.Time
	lastRead     time.Time
}

// NewBotTracker creates new.
func NewBotTracker(botID string) *BotTracker {
	return &BotTracker{
		botID:        botID,
		lastActivity: time.Now(),
	}
}

// IsInactive tells if the bot is inactive.
func (bt *BotTracker) IsInactive() bool {
	// return positive result if we shouldn't read yet
	// this lets us read every once in a while and not detect inactivities
	// too often and make aggressive decisions
	if time.Since(bt.lastRead) < readCooldown {
		return false
	}

	// set the read timestamp so we hit the read cooldown next time
	bt.lastRead = time.Now()

	return time.Since(bt.lastActivity) > inactivityThreshold
}

// SaveActivity saves the activity timestamp when called at the time of an activity.
func (bt *BotTracker) SaveActivity() {
	bt.lastActivity = time.Now()
}

// BotID returns the ID of the bot that is tracked.
func (bt *BotTracker) BotID() string {
	return bt.botID
}
