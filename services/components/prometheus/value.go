package prometheus

import (
	"strconv"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
)

// parseReportValue converts to three types of data to float: timestamp, number, boolean.
// If the value is none of them, finally, it tries to convert the error messages like
// value=1.
func parseReportValue(report *health.Report) (value float64) {
	if n, err := strconv.ParseFloat(report.Details, 64); err == nil {
		value = n
		return
	}

	if b, err := strconv.ParseBool(report.Details); err == nil {
		if b {
			value = 1
		}
		return
	}

	if t, err := time.Parse(time.RFC3339, report.Details); err == nil {
		if t.IsZero() {
			return 0
		}
		since := time.Since(t)
		value = since.Round(time.Second).Seconds()
		return
	}

	// important note: the logic in here is used only if we are trying to use an error message

	switch report.Status {
	case health.StatusOK, health.StatusInfo, health.StatusUnknown:
		value = 0

	case health.StatusFailing, health.StatusLagging, health.StatusDown:
		value = 1
	}
	return
}

func sinceUnixTs(numStr string) float64 {
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}
	if num == 0 {
		return 0
	}
	since := time.Since(time.Unix(int64(num), 0))
	return since.Round(time.Second).Seconds()
}
