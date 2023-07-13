package runner

import (
	"testing"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/stretchr/testify/require"
)

func TestParseReportValue(t *testing.T) {
	r := require.New(t)

	report := &health.Report{
		Name:    "foo",
		Status:  health.StatusOK,
		Details: "1420687622144",
	}

	value := parseReportValue(report)
	r.Equal(float64(1420687622144), value)
}
