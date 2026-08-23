package schema_test

import (
	"encoding/json"
	"testing"

	"hei-gin/internal/framework/core/schema"
)

func TestJSONOrNullEmpty(t *testing.T) {
	var v schema.JSONOrNull
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "null" {
		t.Fatalf("want null, got %s", raw)
	}
}

func TestJSONOrNullObject(t *testing.T) {
	v := schema.JSONOrNullFromBytes([]byte(`{"a":1}`))
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"a":1}` {
		t.Fatalf("unexpected %s", raw)
	}
}

func TestStringPtrBlank(t *testing.T) {
	empty := ""
	if schema.StringPtr(&empty) != nil {
		t.Fatal("blank string should be nil")
	}
}
