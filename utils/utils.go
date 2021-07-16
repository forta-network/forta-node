package utils

// ShortenString shortens the string if the string is longer than given length.
func ShortenString(str string, maxLength int) string {
	if len(str) <= maxLength {
		return str
	}
	return str[:maxLength]
}
