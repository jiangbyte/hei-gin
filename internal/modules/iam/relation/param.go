package relation

// DeptGrantInfo 账号-部门授予明细（授权入参与已拥有结果共用）。
//
// Author: Charlie
type DeptGrantInfo struct {
	DeptID             string   `json:"dept_id" binding:"required"`
	DataScope          string   `json:"data_scope"`
	CustomScopeDeptIDs []string `json:"custom_scope_dept_ids"`
}

// ResourceGrantInfo 主体-资源授予明细（授权入参与已拥有结果共用）。
//
// Author: Charlie
type ResourceGrantInfo struct {
	ResourceID string `json:"resource_id" binding:"required"`
	GrantMode  string `json:"grant_mode"`
	DataScope  string `json:"data_scope"`
}
