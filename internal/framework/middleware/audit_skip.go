// internal/framework/middleware/audit_skip.go 高频操作审计跳过表（对齐 hei-boot AuditSkipCatalog）。
//
// Author: Charlie

package middleware

import "strings"

var auditSkipKeys = map[string]struct{}{
	"auth:refresh":                      {},
	"auth:send_login_code":              {},
	"auth:send_register_code":           {},
	"sys_file:upload":                   {},
	"profile_center:upload_avatar":      {},
	"profile_center:send_password_code": {},
	"profile_center:send_bind_phone_code": {},
	"profile_center:send_bind_email_code": {},
	"sys_notice:read":                   {},
	"sys_notice:read_all":               {},
	"sys_banner:interaction":            {},
	"real_name_case:callback":           {},
}

// ShouldSkipAudit 是否跳过操作审计入库。
func ShouldSkipAudit(resourceType, action string) bool {
	key := normalizeAuditKey(resourceType) + ":" + normalizeAuditKey(action)
	_, ok := auditSkipKeys[key]
	return ok
}

func normalizeAuditKey(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	return strings.ReplaceAll(v, "-", "_")
}
