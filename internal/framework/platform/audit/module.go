// internal/framework/platform/audit/module.go 审计模块名推导（对齐 hei-boot AuditServiceImpl.buildModule）。
//
// Author: Charlie

package audit

// BuildModule 对齐 hei-boot AuditServiceImpl.buildModule：
// resourceType=="resources" → "resource"，其余一律 "iam"。
//
// Author: Charlie
func BuildModule(resourceType string) string {
	if resourceType == "resources" {
		return "resource"
	}
	return "iam"
}
