package health

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HealthClient makes health check requests.
type HealthClient interface {
	CheckHealth(name, port string) Reports
}

type healthClient struct{}

// NewClient creates a new client.
func NewClient() *healthClient {
	return &healthClient{}
}

func containerURL(port string) string {
	return fmt.Sprintf("http://localhost:%s/health", port)
}

func singleReport(name string, status Status, details string) Reports {
	return Reports{
		&Report{
			Name:    name,
			Status:  status,
			Details: details,
		},
	}
}

// shortens to 64 bytes
func shortenResponse(b []byte) []byte {
	if len(b) > 61 {
		return append(b[:61], '.', '.', '.')
	}
	return b
}

type errorResponse struct {
	Error string `json:"error"`
}

func (hc *healthClient) CheckHealth(name, port string) (reports Reports) {
	rawurl := containerURL(port)
	apiName := fmt.Sprintf("%s.health-api", name)
	resp, err := http.Get(rawurl)
	if err != nil {
		return singleReport(apiName, StatusDown, fmt.Sprintf("request failed: %v", err))
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return singleReport(apiName, StatusFailing, fmt.Sprintf("failed to read: %v", err))
	}
	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.Unmarshal(b, &errResp); err != nil {
			return singleReport(apiName, StatusFailing, fmt.Sprintf("bad error response: %v: %s", err, string(b)))
		}
		return singleReport(apiName, StatusFailing, fmt.Sprintf("responded with error: %s", errResp.Error))
	}
	if err := json.Unmarshal(b, &reports); err != nil {
		return singleReport(apiName, StatusFailing, fmt.Sprintf("bad response: %v: %s", err, string(b)))
	}
	// prefix report names
	for _, report := range reports {
		report.Name = fmt.Sprintf("%s.%s", name, report.Name)
	}
	return reports
}
