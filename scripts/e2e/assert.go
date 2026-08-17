package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type caseResult struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Error   string `json:"error,omitempty"`
	URL     string `json:"url,omitempty"`
	Status  int    `json:"status,omitempty"`
	BizCode int    `json:"biz_code,omitempty"`
	Body    string `json:"body,omitempty"`
}

type caseBucket struct {
	Total int          `json:"total"`
	Pass  int          `json:"pass"`
	Fail  []caseResult `json:"fail"`
}

func (b *caseBucket) add(cr caseResult) {
	b.Total++
	if cr.OK {
		b.Pass++
		fmt.Println("PASS", cr.Name)
	} else {
		b.Fail = append(b.Fail, cr)
		fmt.Println("FAIL", cr.Name, cr.Error)
	}
}

// looseEnvelope 用标准 json 解析 stringly 响应（code 为字符串；data 保为 RawMessage）。
type looseEnvelope struct {
	Code    json.RawMessage `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func parseLoose(body []byte) (ar apiResp, data any, err error) {
	var le looseEnvelope
	if err = json.Unmarshal(body, &le); err != nil {
		return ar, nil, fmt.Errorf("envelope json: %w", err)
	}
	ar.Message = le.Message
	ar.Data = le.Data
	ar.Code, _ = parseCode(le.Code)
	if len(le.Data) == 0 || string(le.Data) == "null" {
		return ar, nil, nil
	}
	if err = json.Unmarshal(le.Data, &data); err != nil {
		return ar, nil, fmt.Errorf("data json: %w", err)
	}
	return ar, data, nil
}

func parseCode(raw json.RawMessage) (int, error) {
	if len(raw) == 0 {
		return 0, nil
	}
	var n int
	if err := json.Unmarshal(raw, &n); err == nil {
		return n, nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func parseEnvelope(body []byte) (apiResp, map[string]any, error) {
	ar, data, err := parseLoose(body)
	if err != nil {
		return ar, nil, err
	}
	if data == nil {
		return ar, nil, nil
	}
	if m, ok := data.(map[string]any); ok {
		return ar, m, nil
	}
	return ar, nil, nil
}

func parseDataArray(ar apiResp) ([]any, error) {
	if ar.Data == nil || string(ar.Data) == "null" {
		return nil, fmt.Errorf("data null")
	}
	var data any
	if err := json.Unmarshal(ar.Data, &data); err != nil {
		return nil, err
	}
	arr, ok := data.([]any)
	if !ok {
		return nil, fmt.Errorf("data not array")
	}
	return arr, nil
}

func parseDataArrayMaps(ar apiResp) ([]map[string]any, error) {
	arr, err := parseDataArray(ar)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(arr))
	for _, it := range arr {
		if m, ok := it.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out, nil
}

func assertBizOK(status, code int) error {
	if status < 200 || status >= 300 {
		return fmt.Errorf("http status %d", status)
	}
	if code != 200 && code != 0 {
		return fmt.Errorf("biz code %d", code)
	}
	return nil
}

func assertKeys(m map[string]any, keys ...string) error {
	if m == nil {
		return fmt.Errorf("data is nil")
	}
	var missing []string
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing keys: %s", strings.Join(missing, ","))
	}
	return nil
}

func assertPage(m map[string]any) ([]map[string]any, error) {
	if err := assertKeys(m, "size", "current", "total", "pages", "records"); err != nil {
		return nil, err
	}
	raw, ok := m["records"]
	if !ok {
		return nil, fmt.Errorf("records missing")
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("records not array")
	}
	out := make([]map[string]any, 0, len(arr))
	for _, it := range arr {
		if mm, ok := it.(map[string]any); ok {
			out = append(out, mm)
		}
	}
	return out, nil
}

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		return fmt.Sprintf("%.0f", t)
	default:
		return fmt.Sprint(v)
	}
}

func firstID(records []map[string]any) string {
	if len(records) == 0 {
		return ""
	}
	return asString(records[0]["id"])
}

func findIDByField(records []map[string]any, field, want string) string {
	for _, r := range records {
		if asString(r[field]) == want {
			return asString(r["id"])
		}
	}
	return ""
}
