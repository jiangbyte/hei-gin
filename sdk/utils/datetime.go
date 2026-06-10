package utils

import "time"

var datetimeFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02",
}

// ParseDateTime attempts to parse a string as a datetime using multiple formats.
// Supported formats: "2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02".
// NOTE: Parsed as UTC — use ParseDateTimeLocal for local-timezone strings.
func ParseDateTime(v string) (time.Time, error) {
	for _, fmt := range datetimeFormats {
		if t, err := time.Parse(fmt, v); err == nil {
			return t, nil
		}
	}
	return time.Time{}, &time.ParseError{Layout: datetimeFormats[0], Value: v}
}

// ParseDateTimeLocal parses a datetime string in the system's Local timezone.
// This is the correct choice for parsing cursor values that come from
// FormatDateTime / FormatDateTimePtr output, since those format time.Time
// using the Local timezone.
func ParseDateTimeLocal(v string) (time.Time, error) {
	for _, fmt := range datetimeFormats {
		if t, err := time.ParseInLocation(fmt, v, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, &time.ParseError{Layout: datetimeFormats[0], Value: v}
}

// FormatDateTime formats a time.Time as "2006-01-02 15:04:05".
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDate formats a time.Time as "2006-01-02".
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTimePtr formats a *time.Time as "2006-01-02 15:04:05", returning "" if nil.
func FormatDateTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatDatePtr formats a *time.Time as "2006-01-02", returning "" if nil.
func FormatDatePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}
