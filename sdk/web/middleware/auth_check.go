package middleware

import (
	"strings"

	"hei-gin/sdk/auth"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/config"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

var wsPathSuffixes = []string{"/ws"}

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		if method == "OPTIONS" {
			c.Next()
			return
		}

		for _, suffix := range wsPathSuffixes {
			if strings.HasSuffix(path, suffix) {
				c.Next()
				return
			}
		}

		for _, pp := range config.C.Auth.PublicPaths {
			if pathPrefixMatch(path, pp) {
				attachOptionalLogin(c)
				c.Next()
				return
			}
		}

		if strings.HasPrefix(path, "/api/v") && len(path) > len("/api/v") {
			afterV := path[7:]
			slash := strings.IndexByte(afterV, '/')
			if slash >= 0 {
				afterVer := afterV[slash+1:]
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
