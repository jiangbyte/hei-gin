package middleware

import (
	"net/http"
	"strings"

	"hei-gin/sdk/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	base := cors.New(cors.Config{
		AllowOrigins:     config.C.CORS.AllowOrigins,
		AllowMethods:     config.C.CORS.AllowMethods,
		AllowHeaders:     config.C.CORS.AllowHeaders,
		AllowCredentials: config.C.CORS.AllowCredentials,
	})

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" || strings.EqualFold(origin, "null") {
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
