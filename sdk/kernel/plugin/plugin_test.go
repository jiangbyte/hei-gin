package plugin

import "testing"

type stubPlugin struct {
	name string
	NoopPlugin
}

func (s *stubPlugin) Name() string { return s.name }

func TestRegisterRejectsDuplicatePlugin(t *testing.T) {
	ResetForTest()

	Register(&stubPlugin{name: "dup"})

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected duplicate plugin registration panic")
		}
	}()
	Register(&stubPlugin{name: "dup"})
}

func TestRegisterRejectsAfterFreeze(t *testing.T) {
	ResetForTest()
	Freeze()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected frozen registry panic")
		}
	}()
	Register(&stubPlugin{name: "late"})
}

func TestSnapshotReflectsRegisteredPlugin(t *testing.T) {
	ResetForTest()
	Register(&stubPlugin{name: "p1"})

	snapshot := Snapshot()
	if len(snapshot) != 1 || snapshot[0].Name != "p1" {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
}
