package ratelimiter_test

import (
	"testing"
	"time"

	json_rpc "github.com/forta-network/forta-node/services/json-rpc"
	"github.com/stretchr/testify/require"
)

const testClientID = "1"

func TestRateLimiting(t *testing.T) {
	r := require.New(t)
	rateLimiter := json_rpc.NewRateLimiter(0.5, 1) // replenish every 2s (1/0.5)
	reachedLimit := rateLimiter.ExceedsLimit(testClientID)
	r.False(reachedLimit)
	reachedLimit = rateLimiter.ExceedsLimit(testClientID)
	r.True(reachedLimit)

	time.Sleep(time.Second * 5) // way larger than 2s
	reachedLimit = rateLimiter.ExceedsLimit(testClientID)
	r.False(reachedLimit)
}
