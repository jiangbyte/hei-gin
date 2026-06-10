package utils

import (
	"reflect"
	"strings"

	"hei-gin/sdk/constants"
)

// StripSystemFields removes system audit fields from a map and returns a new map.
// System fields include: id, created_at, created_by, updated_at, updated_by.
// Additional fields can be specified via extraFields.
func StripSystemFields(data map[string]any, extraFields ...string) map[string]any {
	result := make(map[string]any, len(data))
	for k, v := range data {
		if constants.BASE_SYSTEM_FIELDS[k] {
			continue
		}
		skip := false
		for _, ef := range extraFields {
			if k == ef {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		result[k] = v
	}
	return result
}

// ApplyUpdate applies updateData fields to an entity, skipping system-protected fields.
// extraProtected specifies additional fields to protect from updates.
func ApplyUpdate(entity any, updateData map[string]any, extraProtected ...string) {
	protected := make(map[string]bool, len(constants.BASE_SYSTEM_FIELDS)+len(extraProtected))
	for k := range constants.BASE_SYSTEM_FIELDS {
		protected[k] = true
	}
	for _, f := range extraProtected {
		protected[f] = true
	}

	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for key, value := range updateData {
		if protected[key] {
			continue
		}
		// Match by JSON tag first, fallback to field name
		field := findFieldByNameOrTag(typ, key)
		if field == nil || !field.IsExported() {
			continue
		}
		fv := val.FieldByIndex(field.Index)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}
		rv := reflect.ValueOf(value)
		if rv.Type().AssignableTo(fv.Type()) {
			fv.Set(rv)
		}
	}
}

// findFieldByNameOrTag finds a struct field by name or json tag.
func findFieldByNameOrTag(typ reflect.Type, name string) *reflect.StructField {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Name == name {
			return &f
		}
		tag := f.Tag.Get("json")
		if tag != "" {
			tagName := strings.SplitN(tag, ",", 2)[0]
			if tagName == name {
				return &f
			}
		}
	}
	return nil
}

// SafeStrPtr returns the string value from a string pointer, or "" if nil.
func SafeStrPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
