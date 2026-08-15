// internal/modules/iam/relation/param.go 入参定义。
//
// Author: Charlie

package relation

// DeptGrantInfo 账号-部门授予明细（对齐 hei-boot SysDeptGrantResult：dept_id + is_primary）。
//
// Author: Charlie
type DeptGrantInfo struct {
	DeptID    string `json:"dept_id" binding:"required"`
	IsPrimary bool   `json:"is_primary"`
}

// ResourceGrantInfo 主体-资源授予明细（对齐 hei-boot SysResourceGrantResult：resource_id + permission_keys）。
//
// Author: Charlie
type ResourceGrantInfo struct {
	ResourceID     string   `json:"resource_id" binding:"required"`
	PermissionKeys []string `json:"permission_keys"`
}
