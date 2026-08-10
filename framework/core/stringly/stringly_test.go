package stringly_test

import (
	"strings"
	"testing"

	"hei-gin/framework/core/stringly"
)

type sample struct {
	Code    int    `json:"code"`
	Enabled bool   `json:"enabled"`
	Count   int64  `json:"count"`
	Name    string `json:"name"`
	Tags    []int  `json:"tags"`
}

func TestMarshalScalarsAsJSONStrings(t *testing.T) {
	b, err := stringly.Marshal(sample{Code: 200, Enabled: true, Count: 12, Name: "x", Tags: []int{1, 2}})
	if err != nil {
		t.Fatal(err)
	}
	json := string(b)
	for _, want := range []string{`"code":"200"`, `"enabled":"true"`, `"count":"12"`, `"tags":["1","2"]`} {
		if !strings.Contains(json, want) {
			t.Fatalf("missing %s in %s", want, json)
		}
	}
	for _, bad := range []string{`"code":200`, `"enabled":true`} {
		if strings.Contains(json, bad) {
			t.Fatalf("unexpected native scalar %s in %s", bad, json)
		}
	}
}

func TestUnmarshalScalarsFromJSONStrings(t *testing.T) {
	var s sample
	err := stringly.Unmarshal([]byte(`{"code":"401","enabled":"false","count":"99","name":"n","tags":["3","4"]}`), &s)
	if err != nil {
		t.Fatal(err)
	}
	if s.Code != 401 || s.Enabled || s.Count != 99 || s.Name != "n" {
		t.Fatalf("got %+v", s)
	}
	if len(s.Tags) != 2 || s.Tags[0] != 3 || s.Tags[1] != 4 {
		t.Fatalf("tags %+v", s.Tags)
	}
}

func TestUnmarshalRejectsWeirdBool(t *testing.T) {
	var s sample
	err := stringly.Unmarshal([]byte(`{"enabled":"yes"}`), &s)
	if err == nil {
		t.Fatal("expected error")
	}
}
