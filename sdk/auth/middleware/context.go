package middleware

import (
	"context"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/db"

	"github.com/gin-gonic/gin"
)

func AttachLoginContext(c *gin.Context, loginType, uid string) {
	if uid == "" {
		return
	}

	c.Set("login_id", uid)
	c.Set("login_type", loginType)

	if username := auth.GetExtraByType(c, loginType, "username"); username != nil {
		if u, ok := username.(string); ok && u != "" {
			c.Set("loginUser", u)
		}
	}

	ctx := context.WithValue(c.Request.Context(), db.CtxKeyLoginID{}, uid)
	c.Request = c.Request.WithContext(ctx)
}
