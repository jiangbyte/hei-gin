// Package otel 提供可选 OpenTelemetry 初始化占位（尚未接入真实 OTLP exporter）。
//
// Author: Charlie
package otel

import (
	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/logger"

	"go.uber.org/zap"
)

// Init 占位初始化：otel.enabled=false 时 no-op；启用时仅打日志，不导出 traces/metrics。
// 真正接入 otlptracehttp / otlpmetric 时在此替换实现，勿在业务模块直接依赖本包之外的 OTel SDK。
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
	logger.L.Info("otel enabled (placeholder; no exporter wired)", zap.String("endpoint", endpoint))
	return nil
}
