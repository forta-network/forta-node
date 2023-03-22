package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/forta-network/forta-node/clients/ratelimiter"
	"github.com/stretchr/testify/require"
)

const testClientID = "1"

func TestRateLimiting(t *testing.T) {
	r := require.New(t)
	rateLimiter := ratelimiter.NewRateLimiter(0.5, 1) // replenish every 2s (1/0.5)
	reachedLimit := rateLimiter.ExceedsLimit(testClientID)
	r.False(reachedLimit)
	reachedLimit = rateLimiter.ExceedsLimit(testClientID)
	r.True(reachedLimit)

	time.Sleep(time.Second * 5) // way larger than 2s
	reachedLimit = rateLimiter.ExceedsLimit(testClientID)
	r.False(reachedLimit)
}
