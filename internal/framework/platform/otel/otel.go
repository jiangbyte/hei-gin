// Package otel æä¾›å¯é€‰ OpenTelemetry åˆå§‹åŒ–æ¡©ã€‚
//
// Author: Charlie
package otel

import (
	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/logger"

	"go.uber.org/zap"
)

// Init åœ¨ otel.enabled=false æ—¶ no-opï¼›å¯ç”¨æ—¶è®°å½• endpoint å¹¶è¿”å›ž nilï¼ˆé¿å…å¼ºä¾èµ– OTLP SDKï¼‰ã€‚
// åŽç»­å¯åœ¨æ­¤æŽ¥å…¥ otlptracehttp exporterã€‚
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
