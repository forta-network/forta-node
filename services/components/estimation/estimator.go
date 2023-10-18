package estimation

import (
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/feeds/timeline"
)

// Estimator does performance estimations.
type Estimator struct {
	blockTimeline  *timeline.BlockTimeline
	blockThreshold int
}

func NewEstimator(blockTimeline *timeline.BlockTimeline, blockThreshold int) *Estimator {
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
	lag, ok := e.blockTimeline.CalculateLag()
	if !ok {
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
