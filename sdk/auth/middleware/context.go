package middleware

import (
	"context"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/infra/db"

	"github.com/gin-gonic/gin"
)

func AttachLoginContext(c *gin.Context, realm *auth.Realm, uid string) {
	if uid == "" {
		return
	}
	if realm == nil {
		realm = auth.Business
	}

	c.Set("login_id", uid)
	c.Set("login_realm_id", string(realm.ID))
	c.Set("login_realm", realm)

	if username := realm.GetExtra(c, "username"); username != nil {
		if u, ok := username.(string); ok && u != "" {
			c.Set("loginUser", u)
		}
	}

	ctx := context.WithValue(c.Request.Context(), db.CtxKeyLoginID{}, uid)
	c.Request = c.Request.WithContext(ctx)
}
