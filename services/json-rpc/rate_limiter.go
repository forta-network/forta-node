package json_rpc

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// RateLimiter rate limits requests.
type RateLimiter struct {
	bufferSize     int
	cooldown       time.Duration
	clientLimiters map[string]*reqLimiter
	mu             sync.Mutex
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(reqs int, every time.Duration) *RateLimiter {
	if reqs <= 0 {
		log.Panicf("invalid req count '%d' provided to rate limiter", reqs)
	}
	rl := &RateLimiter{
		bufferSize:     reqs,
		cooldown:       (time.Duration)(float64(every) / float64(reqs)),
		clientLimiters: make(map[string]*reqLimiter),
	}
	go rl.autoCleanup()
	return rl
}

// CheckLimit tries adding a request to the limiting channel and returns boolean to signal
// if we hit the rate limit.
func (rl *RateLimiter) CheckLimit(clientID string) bool {
	select {
	case rl.getChannel(clientID) <- struct{}{}:
		return false // not blocked by insert == not hit the limit yet
	default:
		return true
	}
}

func (rl *RateLimiter) getChannel(clientID string) chan struct{} {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	reqLimiter := rl.clientLimiters[clientID]
	if reqLimiter != nil {
		return reqLimiter.limitCh
	}
	reqLimiter = newReqLimiter(clientID, rl.bufferSize, rl.cooldown)
	rl.clientLimiters[clientID] = reqLimiter
	return reqLimiter.limitCh
}

// this keeps allocation under control with little effect to overall functionality
func (rl *RateLimiter) autoCleanup() {
	ticker := time.NewTicker(time.Hour)
	for range ticker.C {
		rl.mu.Lock()
		for _, limiter := range rl.clientLimiters {
			limiter.Close()
		}
		rl.clientLimiters = make(map[string]*reqLimiter)
		rl.mu.Unlock()
	}
}

type reqLimiter struct {
	ctx        context.Context
	cancelFunc func()

	clientID string
	limitCh  chan struct{}
	ticker   *time.Ticker
}

func newReqLimiter(clientID string, size int, cooldown time.Duration) *reqLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	reqLimiter := &reqLimiter{
		ctx:        ctx,
		cancelFunc: cancel,
		clientID:   clientID,
		limitCh:    make(chan struct{}, size),
		ticker:     time.NewTicker(cooldown),
	}
	go reqLimiter.cooldownRequests()
	return reqLimiter
}

func (rl *reqLimiter) cooldownRequests() {
	for {
		select {
		case <-rl.ctx.Done():
			return
		case <-rl.ticker.C:
			<-rl.limitCh
		}
	}
}

func (rl *reqLimiter) Close() error {
	rl.ticker.Stop()
	rl.cancelFunc() // let cooldown goroutine exit
	// do not close the limit channel, it will be deallocated
	return nil
}
