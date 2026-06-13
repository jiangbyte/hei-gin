package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/config"
)

func TestSetupRoutersExposeMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.C = &config.Config{
		App:     config.AppConfig{Name: "hei-gin", Version: "1.0.0"},
		Swagger: config.SwaggerConfig{Enabled: false},
	}

	r := gin.New()
	SetupRouters(r)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "hei_http_inflight_requests") {
		t.Fatalf("metrics body missing expected counter: %s", resp.Body.String())
	}
}

func TestSetupRoutersExposeRegistrySnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.C = &config.Config{
		App:     config.AppConfig{Name: "hei-gin", Version: "1.0.0", Debug: true},
		Swagger: config.SwaggerConfig{Enabled: false},
	}

	r := gin.New()
	SetupRouters(r)

	req := httptest.NewRequest(http.MethodGet, "/debug/registry", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "\"plugins\"") {
		t.Fatalf("registry snapshot body missing plugins field: %s", resp.Body.String())
	}
}

func TestSetupRoutersHideRegistrySnapshotWhenDebugDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.C = &config.Config{
		App:     config.AppConfig{Name: "hei-gin", Version: "1.0.0", Debug: false},
		Swagger: config.SwaggerConfig{Enabled: false},
	}

	r := gin.New()
	SetupRouters(r)

	req := httptest.NewRequest(http.MethodGet, "/debug/registry", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", resp.Code)
	}
}
