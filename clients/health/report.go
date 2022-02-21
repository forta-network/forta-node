package health

import (
	"fmt"
	"strings"
	"time"

	"github.com/forta-protocol/forta-node/utils"
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

// Time tries parsing details as time.
func (report *Report) Time() (*time.Time, bool) {
	t, err := time.Parse(report.Details, time.RFC3339)
	if err != nil {
		return nil, false
	}
	return &t, true
}

// Reports can be marshaled into CSV format.
type Reports []*Report

// NameContains finds a report by checking the name string.
func (reports Reports) NameContains(name string) (*Report, bool) {
	for _, report := range reports {
		if strings.Contains(report.Name, name) {
			return report, true
		}
	}
	return nil, false
}

// GetByName finds a report by checking if there is an exact match.
func (reports Reports) GetByName(name string) (*Report, bool) {
	for _, report := range reports {
		if report.Name == name {
			return report, true
		}
	}
	return nil, false
}

// ObfuscateDetails obfuscates details in each report.
func (reports Reports) ObfuscateDetails() {
	for _, report := range reports {
		report.Details = utils.ObfuscateURLs(report.Details)
	}
}

// SummaryReport implements some methods to help construct summary `Reports` easily.
type SummaryReport struct {
	report   Report
	messages []string
}

// NewSummary creates a new summary report.
func NewSummary() *SummaryReport {
	return &SummaryReport{
		report: Report{
			Name:    "summary",
			Status:  StatusOK,
			Details: "all services are healthy",
		},
	}
}

// Add adds a summary message.
func (sr *SummaryReport) Add(msg string) *SummaryReport {
	if len(sr.messages) > 0 {
		msg = " " + msg
	}
	sr.messages = append(sr.messages, msg)
	return sr
}

// Addf adds a summary message.
func (sr *SummaryReport) Addf(msg string, args ...interface{}) *SummaryReport {
	return sr.Add(fmt.Sprintf(msg, args...))
}

// Punc puts a punctuation mark.
func (sr *SummaryReport) Punc(punc string) *SummaryReport {
	lastMsg := sr.lastMsg()
	if lastMsg == "" || (lastMsg[len(lastMsg)-1] == '.') {
		return sr
	}
	sr.messages = append(sr.messages, punc)
	return sr
}

func (sr *SummaryReport) lastMsg() string {
	length := len(sr.messages)
	if length == 0 {
		return ""
	}
	return sr.messages[length-1]
}

// Status sets the report status.
func (sr *SummaryReport) Status(status Status) *SummaryReport {
	sr.report.Status = status
	return sr
}

// Finish constructs the summary message and returns the constructed report.
func (sr *SummaryReport) Finish() *Report {
	sr.report.Details = strings.Join(sr.messages, "")
	return &sr.report
}

// Fail returns a negative summary.
func (sr *SummaryReport) Fail() *Report {
	sr.report.Status = StatusFailing
	sr.report.Details = "failed to summarize statuses"
	return &sr.report
}
