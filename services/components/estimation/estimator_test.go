package estimation

import (
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/stretchr/testify/require"
)

type testBlockTimeline struct {
	score float64
	delay time.Duration
	ok    bool
}

func (t *testBlockTimeline) EstimateBlockScore() (float64, bool) {
	return t.score, t.ok
}

func (t *testBlockTimeline) GetDelay() (time.Duration, bool) {
	return t.delay, t.ok
}

func TestEstimator(t *testing.T) {
	r := require.New(t)

	blockTimeline := &testBlockTimeline{
		ok: false,
	}
	estimator := NewEstimator(blockTimeline)

	r.Equal(health.StatusUnknown, estimator.Health()[0].Status)
	r.Equal(health.StatusUnknown, estimator.Health()[1].Status)

	blockTimeline.ok = true
	blockTimeline.delay = time.Second
	blockTimeline.score = 0.12

	r.Equal(health.StatusInfo, estimator.Health()[0].Status)
	r.Equal(health.StatusInfo, estimator.Health()[1].Status)
}
