// internal/framework/platform/gojob/gojob_test.go
package gojob

import (
	"testing"
	"time"
)

func TestValidateAndNextCRON(t *testing.T) {
	if err := ValidateTrigger(TypeCRON, "0 0 2 * * *"); err != nil {
		t.Fatal(err)
	}
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	next, err := ComputeNextRunTime(TypeCRON, "0 0 2 * * *", from)
	if err != nil {
		t.Fatal(err)
	}
	if next.Hour() != 2 || next.Minute() != 0 {
		t.Fatalf("next=%v", next)
	}
}

func TestValidateAndNextFIXED(t *testing.T) {
	if err := ValidateTrigger(TypeFIXED, "60"); err != nil {
		t.Fatal(err)
	}
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	next, err := ComputeNextRunTime(TypeFIXED, "60", from)
	if err != nil {
		t.Fatal(err)
	}
	if !next.Equal(from.Add(60 * time.Second)) {
		t.Fatalf("next=%v", next)
	}
	if err := ValidateTrigger(TypeFIXED, "0"); err == nil {
		t.Fatal("expected error")
	}
}
