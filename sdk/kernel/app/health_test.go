package app

import (
	"context"
	"testing"

	"hei-gin/sdk/config"
)

func TestReadinessReportMarksUninitializedDependenciesUnready(t *testing.T) {
	config.C = &config.Config{
		App: config.AppConfig{Name: "hei-gin", Version: "1.0.0"},
	}

	report := readinessReport(context.Background())
	if report.Ready {
		t.Fatal("readiness should be false when db and redis are uninitialized")
	}
	if len(report.Components) != 2 {
		t.Fatalf("components len = %d, want 2", len(report.Components))
	}
}
