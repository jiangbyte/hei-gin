// Package stringly 提供全局 stringly JSON 编解码：
// 标量 bool/数值在线上为 JSON 字符串；对象与数组保持结构。
//
// Author: Charlie
package stringly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Marshal 将 v 编码为 stringly JSON（bool/数字 → 字符串）。
func Marshal(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var tree any
	if err := json.Unmarshal(b, &tree); err != nil {
		return nil, err
	}
	return json.Marshal(stringify(tree))
}

// Unmarshal 将 stringly JSON 解码到 dest（指针）；按目标类型宽松解析字符串标量。
func Unmarshal(data []byte, dest any) error {
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("stringly: Unmarshal destination must be a non-nil pointer")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var raw any
	if err := dec.Decode(&raw); err != nil {
		return err
	}
	return assign(rv.Elem(), raw)
}

func stringify(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, val := range x {
			out[k] = stringify(val)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, val := range x {
			out[i] = stringify(val)
		}
		return out
	case bool:
		if x {
			return "true"
		}
		return "false"
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case json.Number:
		return x.String()
	default:
		return v
	}
}

func assign(dst reflect.Value, raw any) error {
	if !dst.IsValid() {
		return nil
	}
	if raw == nil {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	for dst.Kind() == reflect.Pointer {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()
	}

	if dst.Type() == reflect.TypeOf(time.Time{}) {
		s, err := asString(raw)
		if err != nil {
			return err
		}
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			t, err = time.Parse(time.RFC3339, s)
			if err != nil {
				return fmt.Errorf("stringly: invalid time %q", s)
			}
		}
		dst.Set(reflect.ValueOf(t))
		return nil
	}

	switch dst.Kind() {
	case reflect.Struct:
		m, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("stringly: expected object for %s", dst.Type())
		}
		for name, index := range structFields(dst.Type()) {
			val, ok := m[name]
			if !ok {
				continue
			}
			fv := dst.FieldByIndex(index)
			if err := assign(fv, val); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
		return nil
	case reflect.Slice:
		arr, ok := raw.([]any)
		if !ok {
			return fmt.Errorf("stringly: expected array for %s", dst.Type())
		}
		slice := reflect.MakeSlice(dst.Type(), len(arr), len(arr))
		for i := range arr {
			if err := assign(slice.Index(i), arr[i]); err != nil {
				return err
			}
		}
		dst.Set(slice)
		return nil
	case reflect.Array:
		arr, ok := raw.([]any)
		if !ok {
			return fmt.Errorf("stringly: expected array for %s", dst.Type())
		}
		if len(arr) != dst.Len() {
			return fmt.Errorf("stringly: array length mismatch")
		}
		for i := range arr {
			if err := assign(dst.Index(i), arr[i]); err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		m, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("stringly: expected object for map")
		}
		if dst.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("stringly: only map[string]T supported")
		}
		if dst.IsNil() {
			dst.Set(reflect.MakeMap(dst.Type()))
		}
		for k, val := range m {
			mv := reflect.New(dst.Type().Elem()).Elem()
			if err := assign(mv, val); err != nil {
				return err
			}
			dst.SetMapIndex(reflect.ValueOf(k), mv)
		}
		return nil
	case reflect.Bool:
		b, err := asBool(raw)
		if err != nil {
			return err
		}
		dst.SetBool(b)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := asInt(raw)
		if err != nil {
			return err
		}
		dst.SetInt(n)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := asInt(raw)
		if err != nil {
			return err
		}
		if n < 0 {
			return fmt.Errorf("stringly: negative value for unsigned")
		}
		dst.SetUint(uint64(n))
		return nil
	case reflect.Float32, reflect.Float64:
		f, err := asFloat(raw)
		if err != nil {
			return err
		}
		dst.SetFloat(f)
		return nil
	case reflect.String:
		s, err := asString(raw)
		if err != nil {
			return err
		}
		dst.SetString(s)
		return nil
	case reflect.Interface:
		dst.Set(reflect.ValueOf(raw))
		return nil
	default:
		return fmt.Errorf("stringly: unsupported kind %s", dst.Kind())
	}
}

func structFields(t reflect.Type) map[string][]int {
	out := make(map[string][]int)
	var walk func(t reflect.Type, prefix []int)
	walk = func(t reflect.Type, prefix []int) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			index := append(append([]int{}, prefix...), i)
			if f.Anonymous {
				ft := f.Type
				if ft.Kind() == reflect.Pointer {
					ft = ft.Elem()
				}
				if ft.Kind() == reflect.Struct && ft != reflect.TypeOf(time.Time{}) {
					walk(ft, index)
					continue
				}
			}
			if f.PkgPath != "" {
				continue
			}
			name := jsonName(f)
			if name == "-" || name == "" {
				continue
			}
			if _, exists := out[name]; !exists {
				out[name] = index
			}
		}
	}
	walk(t, nil)
	return out
}

func jsonName(f reflect.StructField) string {
	tag := f.Tag.Get("json")
	if tag == "" {
		return f.Name
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "" {
		return f.Name
	}
	return name
}

func asBool(raw any) (bool, error) {
	switch v := raw.(type) {
	case bool:
		return v, nil
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return false, fmt.Errorf("stringly: empty boolean")
		}
		switch strings.ToLower(s) {
		case "true", "1":
			return true, nil
		case "false", "0":
			return false, nil
		default:
			return false, fmt.Errorf("stringly: expected \"true\" or \"false\", got %q", s)
		}
	case json.Number:
		s := v.String()
		if s == "1" {
			return true, nil
		}
		if s == "0" {
			return false, nil
		}
		return false, fmt.Errorf("stringly: expected \"true\" or \"false\", got %q", s)
	case float64:
		if v == 1 {
			return true, nil
		}
		if v == 0 {
			return false, nil
		}
		return false, fmt.Errorf("stringly: invalid boolean number")
	default:
		return false, fmt.Errorf("stringly: invalid boolean")
	}
}

func asInt(raw any) (int64, error) {
	switch v := raw.(type) {
	case json.Number:
		return v.Int64()
	case float64:
		return int64(v), nil
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return 0, nil
		}
		return strconv.ParseInt(s, 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("stringly: invalid integer")
	}
}

func asFloat(raw any) (float64, error) {
	switch v := raw.(type) {
	case json.Number:
		return v.Float64()
	case float64:
		return v, nil
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return 0, nil
		}
		return strconv.ParseFloat(s, 64)
	default:
		return 0, fmt.Errorf("stringly: invalid float")
	}
}

func asString(raw any) (string, error) {
	switch v := raw.(type) {
	case string:
		return v, nil
	case json.Number:
		return v.String(), nil
	case float64:
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10), nil
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	default:
		return "", fmt.Errorf("stringly: invalid string")
	}
}
