// internal/framework/middleware/audit.go 操作审计中间件（对齐 hei-boot OperationAuditAspect：仅成功入库，success 恒 true）。
//
// Author: Charlie

package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/platform/audit"
)

// AuditPublisher 发布操作审计事件的最小接口（*audit.Queue 满足，测试可注入伪实现）。
//
// Author: Charlie
type AuditPublisher interface {
	Publish(audit.Event)
}

// OperationAudit 按路由声明 resourceType/action，请求成功后发布操作审计事件。
//
// 用法与 RequirePermission 同级挂在路由上（对齐 hei-boot @OperationAudit 传参）：
//
//	api.POST("/v1/admin/sys/accounts/create",
//	  admin,
//	  middleware.RequirePermission(d.Perms, "iam:account:create", "账户创建"),
//	  middleware.OperationAudit(d.Audit, "iam_account", "create"),
//	  s.create,
//	)
//
// 语义对齐 hei-boot @OperationAudit AOP：
//   - 仅当处理器正常返回（status < 400）才记录，失败/越权请求不产生审计行；
//   - success 恒为 true；module 按 BuildModule 规则；
//   - summary 为 "METHOD 完整路径"；account_type 存小写（admin/portal）。
func OperationAudit(pub AuditPublisher, resourceType, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if pub == nil {
			return
		}
		if c.Writer.Status() >= http.StatusBadRequest {
			return
		}
		ctx := c.Request.Context()
		ip := contextx.ClientIP(ctx)
		if ip == "" {
			ip = c.ClientIP()
		}
		ev := audit.Event{
			Module:       audit.BuildModule(resourceType),
			ResourceType: resourceType,
			Action:       action,
			Summary:      c.Request.Method + " " + c.Request.URL.Path,
			RequestID:    contextx.RequestID(ctx),
			IP:           ip,
			UserAgent:    c.Request.UserAgent(),
			Success:      true,
			OccurredAt:   time.Now().UTC(),
		}
		if sess := contextx.Session(ctx); sess != nil {
			ev.AccountID = sess.AccountID
			ev.AccountType = strings.ToLower(string(sess.AccountType))
		}
		pub.Publish(ev)
	}
}
