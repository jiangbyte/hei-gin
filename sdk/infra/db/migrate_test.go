package db

import "testing"

type testModel struct{}

func TestRegisterModelRejectsDuplicate(t *testing.T) {
	ResetForTest()
	RegisterModel(&testModel{})

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected duplicate model registration panic")
		}
	}()
	RegisterModel(&testModel{})
}

func TestRegisterSeedRejectsAfterFreeze(t *testing.T) {
	ResetForTest()
	Freeze()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected frozen seed registry panic")
		}
	}()
	RegisterSeed("s1", func() error { return nil })
}

func TestSnapshotContainsModelAndSeed(t *testing.T) {
	ResetForTest()
	RegisterModel(&testModel{})
	RegisterSeed("seed-1", func() error { return nil })

	snapshot := Snapshot()
	if len(snapshot.Models) != 1 || len(snapshot.Seeds) != 1 {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
}
