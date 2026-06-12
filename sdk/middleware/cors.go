package middleware

import (
	"net/http"
	"strings"

	"hei-gin/sdk/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a Gin middleware that configures Cross-Origin Resource Sharing.
func CORS() gin.HandlerFunc {
	base := cors.New(cors.Config{
		AllowOrigins:     config.C.CORS.AllowOrigins,
		AllowMethods:     config.C.CORS.AllowMethods,
		AllowHeaders:     config.C.CORS.AllowHeaders,
		AllowCredentials: config.C.CORS.AllowCredentials,
	})

	// Fallback for requests with Origin: null (Electron apps, file://, etc.).
	// CORS spec does not allow the wildcard "*" to match a null origin,
	// so we must return Access-Control-Allow-Origin: null explicitly.
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" || strings.EqualFold(origin, "null") {
			// Handle preflight
			if c.Request.Method == http.MethodOptions {
				c.Header("Access-Control-Allow-Origin", "null")
				c.Header("Access-Control-Allow-Methods", strings.Join(config.C.CORS.AllowMethods, ", "))
				c.Header("Access-Control-Allow-Headers", strings.Join(config.C.CORS.AllowHeaders, ", "))
				c.Header("Access-Control-Max-Age", "86400")
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.Header("Access-Control-Allow-Origin", "null")
			c.Next()
			return
		}
		base(c)
	}
}
