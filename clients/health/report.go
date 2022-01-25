package health

import (
	"strings"
)

// Status represents status of a component.
type Status string

// Statuses
const (
	StatusOK      Status = "ok"
	StatusDown    Status = "down"
	StatusFailing Status = "failing"
	StatusLagging Status = "lagging"
	StatusInfo    Status = "info"
	StatusUnknown Status = "unknown"
)

// Report contains health info from a service or component.
type Report struct {
	Name    string `json:"name"`
	Status  Status `json:"status"`
	Details string `json:"details"`
}

// ReportName constructs a new report name in publisher.alert_client.* format.
func ReportName(service, component string, properties ...string) string {
	var names []string
	if len(service) > 0 {
		names = append(names, service)
	}
	if len(component) > 0 {
		names = append(names, component)
	}
	for _, prop := range properties {
		if len(prop) > 0 {
			names = append(names, prop)
		}
	}
	return strings.Join(names, ".")
}

// Reports can be marshaled into CSV format.
type Reports []*Report
