package prometheus

import (
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/stretchr/testify/require"
)

func TestParseReportValue(t *testing.T) {
	testCases := []struct {
		report         *health.Report
		expectedResult float64
		approximate    bool
	}{
		{
			report: &health.Report{
				Name:    "foo.service.bar.m1",
				Status:  health.StatusOK,
				Details: "1420687622144",
			},
			expectedResult: 1420687622144,
		},
		{
			report: &health.Report{
				Name:    "foo.service.bar.m2",
				Status:  health.StatusOK,
				Details: "true",
			},
			expectedResult: 1,
		},
		{
			report: &health.Report{
				Name:    "foo.service.bar.m3",
				Status:  health.StatusOK,
				Details: time.Now().Add(time.Second * -10).Format(time.RFC3339),
			},
			expectedResult: 10,
			approximate:    true,
		},
		{
			report: &health.Report{
				Name:    "foo.service.bar.m4",
				Status:  health.StatusOK,
				Details: "message 1",
			},
			expectedResult: 0,
		},
		{
			report: &health.Report{
				Name:    "foo.service.bar.m5",
				Status:  health.StatusFailing,
				Details: "message 2",
			},
			expectedResult: 1,
		},
		{
			report: &health.Report{
				Name:    "skipped-invalid-name-pattern",
				Status:  health.StatusFailing,
				Details: "message X",
			},
			expectedResult: 1,
		},
		{
			report: &health.Report{
				Name:    "foo.service.skipped-single-name",
				Status:  health.StatusFailing,
				Details: "message Y",
			},
			expectedResult: 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.report.Name, func(t *testing.T) {
			r := require.New(t)

			result := parseReportValue(testCase.report)
			if testCase.approximate {
				r.InEpsilon(testCase.expectedResult, result, 0.5)
			} else {
				r.Equal(testCase.expectedResult, result)
			}
		})
	}
}
