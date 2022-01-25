package health

import (
	"sync"
	"time"
)

// TimeTracker is useful for tracking activity time in implementations.
type TimeTracker struct {
	ts *time.Time
	mu sync.RWMutex
}

// Set sets the timer.
func (tt *TimeTracker) Set() {
	tt.mu.Lock()
	ts := time.Now()
	tt.ts = &ts
	tt.mu.Unlock()
}

// Check checks the time.
func (tt *TimeTracker) Check(timeout time.Duration) (formatted string, status Status) {
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	if tt.ts == nil {
		return "", StatusUnknown
	}

	isLate := tt.ts.Add(timeout).Before(time.Now())
	if isLate {
		return tt.string(), StatusLagging
	}
	return tt.string(), StatusOK
}

// String implements the fmt.Stringer interface.
func (tt *TimeTracker) String() string {
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	return tt.string()
}

func (tt *TimeTracker) string() string {
	if tt.ts == nil {
		return ""
	}
	return tt.ts.Format(time.RFC3339)
}

// GetReport constructs and returns a report from check results.
func (tt *TimeTracker) GetReport(name string) *Report {
	var report Report
	report.Name = name
	report.Details, report.Status = tt.Check(time.Minute * 5)
	return &report
}

// ErrorTracker is useful for tracking the latest error.
type ErrorTracker struct {
	err error
	mu  sync.RWMutex
}

// Set sets the tracker.
func (et *ErrorTracker) Set(err error) {
	et.mu.Lock()
	et.err = err
	et.mu.Unlock()
}

// GetReport constructs and returns a report.
func (et *ErrorTracker) GetReport(name string) *Report {
	et.mu.RLock()
	defer et.mu.RUnlock()
	var report Report
	report.Name = name
	report.Status = StatusOK
	if et.err != nil {
		report.Status = StatusFailing
		report.Details = et.err.Error()
	}
	return &report
}

// MessageTracker is useful for tracking the latest message about something.
type MessageTracker struct {
	msg string
	mu  sync.RWMutex
}

// Set sets the tracker.
func (mt *MessageTracker) Set(msg string) {
	mt.mu.Lock()
	mt.msg = msg
	mt.mu.Unlock()
}

// GetReport constructs and returns a report.
func (mt *MessageTracker) GetReport(name string) *Report {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	var report Report
	report.Name = name
	report.Status = StatusInfo
	report.Details = mt.msg
	return &report
}
