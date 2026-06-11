package middleware

import (
	"strings"

	"hei-gin/sdk/auth"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/config"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/result"

	"github.com/gin-gonic/gin"
)

var wsPathSuffixes = []string{"/ws"}

// AuthCheck returns a global middleware that:
//  1. Skips auth for configured public paths
//  2. Requires CONSUMER login for /api/v{n}/c/... routes
//  3. Requires BUSINESS login for all other /api/v{n}/... routes
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// OPTIONS – no auth
		if method == "OPTIONS" {
			c.Next()
			return
		}

		// WebSocket – let WS handler do its own token-based auth
		for _, suffix := range wsPathSuffixes {
			if strings.HasSuffix(path, suffix) {
				c.Next()
				return
			}
		}

		// Public paths from config – no auth
		for _, pp := range config.C.Auth.PublicPaths {
			if pathPrefixMatch(path, pp) {
				attachOptionalLogin(c)
				c.Next()
				return
			}
		}

		// /api/v{n}/c/... → CONSUMER auth
		if strings.HasPrefix(path, "/api/v") && len(path) > len("/api/v") {
			afterV := path[7:] // after "/api/v"
			slash := strings.IndexByte(afterV, '/')
			if slash >= 0 {
				afterVer := afterV[slash+1:] // e.g. "c/..." or "sys/..."
				if strings.HasPrefix(afterVer, "c/") {
					uid := auth.Consumer.GetLoginIDDefaultNull(c)
					if uid == "" {
						c.Abort()
						result.Failure(c, "未授权/未登录", 401)
						return
					}
					authMW.AttachLoginContext(c, string(enums.LoginTypeConsumer), uid)
					c.Next()
					return
				}

				// All other /api/v{n}/... routes → BUSINESS auth
				uid := auth.GetLoginIDDefaultNull(c)
				if uid == "" {
					c.Abort()
					result.Failure(c, "未授权/未登录", 401)
					return
				}
				authMW.AttachLoginContext(c, string(enums.LoginTypeBusiness), uid)
				c.Next()
				return
			}
		}

		attachOptionalLogin(c)
		c.Next()
	}
}

func attachOptionalLogin(c *gin.Context) {
	if uid := auth.GetLoginIDDefaultNull(c); uid != "" {
		authMW.AttachLoginContext(c, string(enums.LoginTypeBusiness), uid)
	} else if uid := auth.Consumer.GetLoginIDDefaultNull(c); uid != "" {
		authMW.AttachLoginContext(c, string(enums.LoginTypeConsumer), uid)
	}
}

func pathPrefixMatch(path, prefix string) bool {
	if prefix == "" {
		return false
	}
	if prefix == "/" {
		return path == "/"
	}
	if strings.HasSuffix(prefix, "/") {
		return strings.HasPrefix(path, prefix)
	}
	return path == prefix || strings.HasPrefix(path, prefix+"/")
}
