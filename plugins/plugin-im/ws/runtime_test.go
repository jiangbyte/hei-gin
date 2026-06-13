package ws

import (
	"net/http"
	"sync"
	"testing"

	"hei-gin/sdk/config"
)

func TestCheckOriginHonorsAllowedOrigins(t *testing.T) {
	config.C = &config.Config{
		Raw: map[string]any{
			"ws": map[string]any{
				"allowed_origins": []string{"https://allowed.example.com"},
			},
		},
	}
	resetConfigCacheForTest()

	req, _ := http.NewRequest(http.MethodGet, "http://localhost/ws", nil)
	req.Header.Set("Origin", "https://allowed.example.com")
	if !checkOrigin(req) {
		t.Fatal("allowed origin should pass")
	}

	req.Header.Set("Origin", "https://blocked.example.com")
	if checkOrigin(req) {
		t.Fatal("blocked origin should fail")
	}
}

func TestGetClientIPIgnoresForwardHeadersForUntrustedProxy(t *testing.T) {
	config.C = &config.Config{Raw: map[string]any{}}
	resetConfigCacheForTest()

	req, _ := http.NewRequest(http.MethodGet, "http://localhost/ws", nil)
	req.RemoteAddr = "10.0.0.8:9000"
	req.Header.Set("X-Forwarded-For", "1.1.1.1")
	if got := getClientIP(req); got != "10.0.0.8" {
		t.Fatalf("client ip = %q, want remote addr ip", got)
	}
}

func resetConfigCacheForTest() {
	_once = sync.Once{}
	_cfg = WSConfig{}
	_upgrader = nil
}
