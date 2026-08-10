package security

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

// ParseAccountType 解析账号类型字符串，非法时返回 false。
func ParseAccountType(s string) (AccountType, bool) {
	switch AccountType(s) {
	case AccountAdmin, AccountPortal:
		return AccountType(s), true
	}
	return "", false
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
