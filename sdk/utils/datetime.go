package utils

import "time"

// FormatDateTime formats a time.Time as "2006-01-02 15:04:05".
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDateTimePtr formats a *time.Time as "2006-01-02 15:04:05", returning "" if nil.
func FormatDateTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// ParseDateTime parses a "2006-01-02 15:04:05" string into time.Time using local timezone.
// Also accepts "2006-01-02T15:04:05" or "2006-01-02" for flexibility.
func ParseDateTime(v string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.ParseInLocation(f, v, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, &time.ParseError{Layout: formats[0], Value: v}
}

// ParseDateTimePtr parses a *string in "2006-01-02 15:04:05" format into *time.Time.
// Returns nil if v is nil, empty, or unparseable.
func ParseDateTimePtr(v *string) *time.Time {
	if v == nil || *v == "" {
		return nil
	}
	t, err := ParseDateTime(*v)
	if err != nil {
		return nil
	}
	return &t
}
