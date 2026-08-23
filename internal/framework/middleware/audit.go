// internal/framework/middleware/audit.go 操作审计中间件（对齐 hei-boot OperationAuditAspect）。
//
// Author: Charlie

package middleware

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/platform/audit"
)

var auditResourceIDRe = regexp.MustCompile(`[0-9a-fA-F]{8,}`)

// AuditPublisher 发布操作审计事件的最小接口（*audit.Queue 满足，测试可注入伪实现）。
//
// Author: Charlie
type AuditPublisher interface {
	Publish(audit.Event)
}

// OperationAudit 按路由声明 resourceType/action，请求结束后发布操作审计事件。
func OperationAudit(pub AuditPublisher, resourceType, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if ShouldSkipAudit(resourceType, action) {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		if pub == nil {
			return
		}
		status := c.Writer.Status()
		success := status < http.StatusBadRequest
		ctx := c.Request.Context()
		ip := contextx.ClientIP(ctx)
		if ip == "" {
			ip = c.ClientIP()
		}
		durationMs := int(time.Since(start).Milliseconds())
		if durationMs < 0 {
			durationMs = 0
		}
		ev := audit.Event{
			Module:       audit.BuildModule(resourceType),
			ResourceType: resourceType,
			Action:       action,
			ResourceID:   extractAuditResourceID(c.Request.URL.Path),
			Summary:      c.Request.Method + " " + c.Request.URL.Path,
			RequestID:    contextx.RequestID(ctx),
			IP:           ip,
			UserAgent:    c.Request.UserAgent(),
			Success:      success,
			OccurredAt:   time.Now().UTC(),
			Extra: map[string]any{
				"duration_ms": durationMs,
			},
		}
		if !success {
			if msg := strings.TrimSpace(c.Errors.String()); msg != "" {
				ev.ErrorMessage = msg
			} else {
				ev.ErrorMessage = http.StatusText(status)
			}
		}
		if sess := contextx.Session(ctx); sess != nil {
			ev.AccountID = sess.AccountID
			ev.AccountType = strings.ToLower(string(sess.AccountType))
		}
		pub.Publish(ev)
	}
}

func extractAuditResourceID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	last := parts[len(parts)-1]
	if auditResourceIDRe.MatchString(last) {
		return last
	}
	return ""
}
