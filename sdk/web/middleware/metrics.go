package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/observability"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		observability.IncHTTPInflight()
		defer observability.DecHTTPInflight()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		observability.ObserveHTTPRequest(c.Request.Method, route, c.Writer.Status(), time.Since(start).Seconds())
	}
}
