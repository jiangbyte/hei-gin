// Package otel 提供可选 OpenTelemetry 初始化桩。
//
// Author: Charlie
package otel

import (
	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/logger"

	"go.uber.org/zap"
)

// Init 在 otel.enabled=false 时 no-op；启用时记录 endpoint 并返回 nil（避免强依赖 OTLP SDK）。
// 后续可在此接入 otlptracehttp exporter。
//
// Author: Charlie
func Init(cfg config.OTelConfig) error {
	if !cfg.Enabled {
		return nil
	}
	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = "localhost:4317"
	}
	logger.L.Info("otel enabled (stub; no exporter wired)", zap.String("endpoint", endpoint))
	return nil
}
