package cooldown

import (
	"sync"
	"time"
)

// Cooldown keeps track of operations by id and tells when an operation
// should be cooling down.
type Cooldown interface {
	ShouldCoolDown(id string) bool
}

type cooldownCounter struct {
	count          int
	cooldownEndsAt time.Time
}

type cooldown struct {
	threshold        int
	cooldownDuration time.Duration
	counters         map[string]*cooldownCounter
	mu               sync.Mutex
}

// New creates a new cooldown.
func New(threshold int, cooldownDuration time.Duration) *cooldown {
	cd := &cooldown{
		threshold:        threshold,
		cooldownDuration: cooldownDuration,
		counters:         make(map[string]*cooldownCounter),
	}
	go cd.autoCleanup()
	return cd
}

// ShouldCoolDown tells if the operation with given id should cool down and
// increments the cooldownCounter.
func (cd *cooldown) ShouldCoolDown(id string) bool {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	counter := cd.counters[id]
	if counter == nil {
		counter = &cooldownCounter{count: 1}
		cd.counters[id] = counter
		return false
	}
	if time.Now().Before(counter.cooldownEndsAt) {
		return true
	}
	if counter.count >= cd.threshold {
		counter.cooldownEndsAt = time.Now().Add(cd.cooldownDuration)
		counter.count = 0
		return true
	}
	counter.count++
	return false
}

// deallocate inactive cooldown counters
func (cd *cooldown) autoCleanup() {
	ticker := time.NewTicker(time.Hour)
	for range ticker.C {
		cd.mu.Lock()
		for id, cooldownCounter := range cd.counters {
			if time.Since(cooldownCounter.cooldownEndsAt) > time.Hour {
				delete(cd.counters, id)
			}
		}
		cd.mu.Unlock()
	}
}
