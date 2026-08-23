// Package schema wire helpers — JSON scalars as strings (align hei-boot StringlyTypedJacksonModule).
//
// Author: Charlie
package schema

import (
	"encoding/json"
	"strconv"
	"strings"
)

// WireBool serializes as "true"/"false" strings.
type WireBool bool

// MarshalJSON implements json.Marshaler.
func (b WireBool) MarshalJSON() ([]byte, error) {
	if b {
		return []byte(`"true"`), nil
	}
	return []byte(`"false"`), nil
}

// WireBoolPtr converts *bool to *WireBool.
func WireBoolPtr(v *bool) *WireBool {
	if v == nil {
		return nil
	}
	b := WireBool(*v)
	return &b
}

// WireBoolValue converts bool to WireBool.
func WireBoolValue(v bool) WireBool { return WireBool(v) }

// IntStringPtr formats *int as decimal string pointer (nil stays nil).
func IntStringPtr(v *int) *string {
	if v == nil {
		return nil
	}
	s := strconv.Itoa(*v)
	return &s
}

// IntString formats int as decimal string.
func IntString(v int) string { return strconv.Itoa(v) }

// StringPtr returns nil for blank strings (boot nullable String fields).
func StringPtr(v *string) *string {
	if v == nil {
		return nil
	}
	if strings.TrimSpace(*v) == "" {
		return nil
	}
	return v
}

// JSONOrNull marshals empty JSON as null (align hei-boot Map/JSON null).
type JSONOrNull json.RawMessage

// MarshalJSON implements json.Marshaler.
func (j JSONOrNull) MarshalJSON() ([]byte, error) {
	raw := []byte(j)
	if len(raw) == 0 || string(raw) == "null" || string(raw) == "{}" || string(raw) == "[]" {
		return []byte("null"), nil
	}
	return json.RawMessage(raw).MarshalJSON()
}

// JSONOrNullFromBytes builds JSONOrNull from DB bytes.
func JSONOrNullFromBytes(raw []byte) JSONOrNull {
	if len(raw) == 0 || string(raw) == "null" {
		return JSONOrNull(nil)
	}
	return JSONOrNull(raw)
}
