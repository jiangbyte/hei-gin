// internal/framework/middleware/metrics.go 指标。
//
// Author: Charlie

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusHandler 暴露 Prometheus 默认注册表（含 Go 运行时指标）。
//
// Author: Charlie
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
