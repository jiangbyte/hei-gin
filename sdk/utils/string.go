package utils

// SafeStrPtr returns the string value from a string pointer, or "" if nil.
func SafeStrPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
