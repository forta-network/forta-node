package estimation

import (
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
)

// BlockTimeline has block creation and processing knowledge
// and can make estimations and calculations based on that.
type BlockTimeline interface {
	EstimateBlockScore() (float64, bool)
	GetDelay() (time.Duration, bool)
}

// Estimator does performance estimations.
type Estimator struct {
	blockTimeline BlockTimeline
}

func NewEstimator(blockTimeline BlockTimeline) *Estimator {
	return &Estimator{
		blockTimeline: blockTimeline,
	}
}

// Name implements health.Reporter.
func (e *Estimator) Name() string {
	return "estimator"
}

// Health implements health.Reporter.
func (e *Estimator) Health() health.Reports {
	estimate, ok := e.blockTimeline.EstimateBlockScore()
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
	jsonRpcPerformanceStr := fmt.Sprintf("%.2f", estimate)
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
