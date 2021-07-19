package utils

// String converts string ptr to string.
func String(str *string) string {
	if str != nil {
		return *str
	}
	return ""
}

// StringPtr converts string to string ptr.
func StringPtr(str string) *string {
	return &str
}
