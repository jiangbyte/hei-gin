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

// Audit 请求成功后按路由注册表发布操作审计事件。
//
// 语义对齐 hei-boot @OperationAudit AOP：
//   - 仅当处理器正常返回（status < 400）才记录，失败/越权请求不产生审计行；
//   - success 恒为 true；module 按 buildModule 规则（resources→resource，其余 iam）；
//   - summary 为 "METHOD 完整路径"；account_type 存小写（admin/portal）。
func Audit(reg *audit.Registry, pub AuditPublisher) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if pub == nil || reg == nil {
			return
		}
		if c.Writer.Status() >= http.StatusBadRequest {
			// 失败请求不入库（对齐 hei-boot proceed 抛错时不发布）
			return
		}
		spec, ok := reg.Match(c.Request.Method, c.Request.URL.Path)
		if !ok {
			return
		}
		ctx := c.Request.Context()
		ip := contextx.ClientIP(ctx)
		if ip == "" {
			ip = c.ClientIP()
		}
		ev := audit.Event{
			Module:       audit.BuildModule(spec.ResourceType),
			ResourceType: spec.ResourceType,
			Action:       spec.Action,
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
