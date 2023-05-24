package lifecycle

import (
	"time"
)

// Timeouts
const (
	readCooldown        = time.Minute * 5
	inactivityThreshold = time.Minute * 15
)

// TrackerStatus is tracker status enum type.
type TrackerStatus int

// Activity statuses
const (
	TrackerStatusActive TrackerStatus = iota + 1
	TrackerStatusInactive
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

// Status checks the bot tracker status.
func (bt *BotTracker) Status() TrackerStatus {
	// return positive result if we shouldn't read yet
	// this lets us read every once in a while and not detect inactivities
	// too often and make aggressive decisions
	if time.Since(bt.lastRead) < readCooldown {
		return TrackerStatusActive
	}
	bt.lastRead = time.Now()
	if time.Since(bt.lastActivity) > inactivityThreshold {
		return TrackerStatusInactive
	}
	return TrackerStatusActive
}

// SaveActivity saves the activity timestamp when called at the time of an activity.
func (bt *BotTracker) SaveActivity() {
	bt.lastActivity = time.Now()
}

// BotID returns the ID of the bot that is tracked.
func (bt *BotTracker) BotID() string {
	return bt.botID
}
