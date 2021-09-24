package utils

import "time"

// FormatTime formats given time to as string format.
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTime parses the formatted string time.
func ParseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
