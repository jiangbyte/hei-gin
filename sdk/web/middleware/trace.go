package middleware

import (
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("trace_id")
		if traceID == "" {
			traceID = utils.GenerateTraceID()
		}
		c.Set("trace_id", traceID)
		c.Next()
	}
}
