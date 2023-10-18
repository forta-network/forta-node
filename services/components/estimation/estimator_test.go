package estimation

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/feeds/timeline"
	"github.com/stretchr/testify/require"
)

// TODO: No need to test the estimation score here. That should be moved to the timeline.
func TestEstimator(t *testing.T) {
	r := require.New(t)

	blockTimeline := &timeline.BlockTimeline{}
	threshold := 10
	estimator := NewEstimator(blockTimeline, threshold)

	currMin := time.Now().UTC().Truncate(time.Minute)
	min1 := currMin.Add(time.Minute * -2)
	min2 := currMin.Add(time.Minute * -1)
	min3 := currMin

	min1Ts := hexutil.EncodeUint64(uint64(min1.Unix()))
	min2Ts := hexutil.EncodeUint64(uint64(min2.Unix()))
	min3Ts := hexutil.EncodeUint64(uint64(min3.Unix()))

	// no blocks handled yet: should give unknown estimation result
	r.Equal(health.StatusUnknown, estimator.estimate(min1)[0].Status)

	// add first minute block number: should give unknown estimation result
	blockTimeline.HandleBlock(blockForTimestamp(min1Ts, "0x100"))
	result := estimator.estimate(min1)[0]
	r.Equal(health.StatusUnknown, result.Status)

	// add second minute block number: should give unknown estimation result
	blockTimeline.HandleBlock(blockForTimestamp(min2Ts, "0x200"))
	result = estimator.estimate(min2)[0]
	r.Equal(health.StatusUnknown, result.Status)

	// add third minute block number: should give an estimation result
	blockTimeline.HandleBlock(blockForTimestamp(min3Ts, "0x300"))
	result = estimator.estimate(min3)[0]
	r.Equal(health.StatusInfo, result.Status)
	r.Equal("1.00", result.Details)
}

func blockForTimestamp(ts, blockNumber string) *domain.BlockEvent {
	return &domain.BlockEvent{
		Block: &domain.Block{
			Timestamp: ts,
			Number:    blockNumber,
		},
	}
}
