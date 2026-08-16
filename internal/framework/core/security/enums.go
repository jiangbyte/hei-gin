// internal/framework/core/security/enums.go 枚举常量。
//
// Author: Charlie

package security

import "strings"

// AccountType 划分 API 面（/admin|/portal）。
//
// Author: Charlie
type AccountType string

const (
	AccountAdmin  AccountType = "ADMIN"
	AccountPortal AccountType = "PORTAL"
)

// URLSegment 返回路径段 admin 或 portal。
func (t AccountType) URLSegment() string {
	switch t {
	case AccountAdmin:
		return "admin"
	case AccountPortal:
		return "portal"
	default:
		return string(t)
	}
}

// DataScope 是权限授予上的行级数据范围。
//
// Author: Charlie
type DataScope string

const (
	DataScopeAll          DataScope = "ALL"
	DataScopeDeptAndChild DataScope = "DEPT_AND_CHILD"
	DataScopeDept         DataScope = "DEPT"
	DataScopeSelf         DataScope = "SELF"
	DataScopeCustom       DataScope = "CUSTOM"
)

const (
	StatusEnabled  = "ENABLED"
	StatusDisabled = "DISABLED"
)

const (
	AccountStatusEnabled   = "ENABLED"
	AccountStatusDisabled  = "DISABLED"
	AccountStatusCancelled = "CANCELLED"
)

// DeviceLabelFromUserAgent 根据 User-Agent 粗略推断设备标签（对齐 boot/fastapi）。
func DeviceLabelFromUserAgent(userAgent string) *string {
	if userAgent == "" {
		return nil
	}
	value := strings.ToLower(userAgent)
	var label string
	switch {
	case strings.Contains(value, "mobile"), strings.Contains(value, "android"), strings.Contains(value, "iphone"):
		label = "Mobile"
	case strings.Contains(value, "ipad"), strings.Contains(value, "tablet"):
		label = "Tablet"
	default:
		label = "Desktop"
	}
	return &label
}
