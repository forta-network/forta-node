package json_rpc

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// RateLimiter rate limits requests.
type RateLimiter struct {
	rate           int
	burst          int
	clientLimiters map[string]*rate.Limiter
	mu             sync.Mutex
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(rateN, burst int) *RateLimiter {
	if rateN <= 0 || burst <= 0 {
		log.Panic("non-positive rate limiter arg")
	}
	rl := &RateLimiter{
		rate:           rateN,
		burst:          burst,
		clientLimiters: make(map[string]*rate.Limiter),
	}
	go rl.autoCleanup()
	return rl
}

// CheckLimit tries adding a request to the limiting channel and returns boolean to signal
// if we hit the rate limit.
func (rl *RateLimiter) CheckLimit(clientID string) bool {
	return rl.reserveClient(clientID).Delay() > 0
}

func (rl *RateLimiter) reserveClient(clientID string) *rate.Reservation {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	limiter := rl.clientLimiters[clientID]
	if limiter == nil {
		limiter = rate.NewLimiter(rate.Limit(rl.rate), rl.burst)
		rl.clientLimiters[clientID] = limiter
	}
	return limiter.Reserve()
}

// this keeps allocation under control with little effect to overall functionality
func (rl *RateLimiter) autoCleanup() {
	ticker := time.NewTicker(time.Hour)
	for range ticker.C {
		rl.mu.Lock()
		for clientID, limiter := range rl.clientLimiters {
			// if it allows max burst now, then it makes sense to deallocate it
			if isNotActive := limiter.AllowN(time.Now(), rl.burst); isNotActive {
				rl.clientLimiters[clientID] = nil
			}
		}
		rl.mu.Unlock()
	}
}
