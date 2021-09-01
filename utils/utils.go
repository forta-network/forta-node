package utils

import "time"

// ShortenString shortens the string if the string is longer than given length.
func ShortenString(str string, maxLength int) string {
	if len(str) <= maxLength {
		return str
	}
	return str[:maxLength]
}

// TryTimes will try an action up to `times` times with a delay of `delay` between each attempt
func TryTimes(handler func(attempt int) error, times int, delay time.Duration) error {
	var result error
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for i := 0; i < times; i++ {
		if err := handler(i); err == nil {
			return nil
		} else {
			result = err
		}
		<-ticker.C
	}
	return result
}

func MapKeys(m map[string]bool) []string {
	var res []string
	for k := range m {
		res = append(res, k)
	}
	return res
}
