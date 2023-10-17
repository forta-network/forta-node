package estimation

import (
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/feeds/timeline"
)

// Estimator does performance estimations.
type Estimator struct {
	blockTimeline    *timeline.BlockTimeline
	blockThreshold   int
	expectedDistance int
}

func NewEstimator(blockTimeline *timeline.BlockTimeline, blockThreshold, expectedDistance int) *Estimator {
	return &Estimator{
		blockTimeline:  blockTimeline,
		blockThreshold: blockThreshold,
	}
}

// Name implements health.Reporter.
func (e *Estimator) Name() string {
	return "estimator"
}

// Health implements health.Reporter.
func (e *Estimator) Health() health.Reports {
	// only look at the previous minute
	return e.estimate(time.Now().Add(time.Minute * -1))
}

func (e *Estimator) estimate(atTime time.Time) health.Reports {
	// need at least two minutes in the time line to start calculating the lag
	// 1st min: unreliable numbers
	// 2nd min: reliable numbers
	// 3rd min: means the 2nd min is over and we should look at that
	tooEarly := e.blockTimeline.Size() < 3
	lag, ok := e.blockTimeline.CalculateLag(atTime)
	if !ok || tooEarly {
		return health.Reports{
			{
				Name:   "json-rpc-performance",
				Status: health.StatusUnknown,
			},
			{
				Name:   "json-rpc-delay",
				Status: health.StatusUnknown,
			},
		}
	}

	// offset the lag by expected distance
	lag = lag - int64(e.expectedDistance)

	jsonRpcPerformance := (float64(e.blockThreshold) - float64(lag)) / float64(e.blockThreshold)
	if jsonRpcPerformance < 0 {
		jsonRpcPerformance = 0
	}
	if jsonRpcPerformance > 1 {
		jsonRpcPerformance = 1
	}

	jsonRpcPerformanceStr := fmt.Sprintf("%.2f", jsonRpcPerformance)
	delay, _ := e.blockTimeline.GetDelay()
	return health.Reports{
		{
			Name:    "json-rpc-performance",
			Status:  health.StatusInfo,
			Details: jsonRpcPerformanceStr,
		},
		{
			Name:    "json-rpc-delay",
			Status:  health.StatusInfo,
			Details: delay.String(),
		},
	}
}
