// internal/modules/iam/resource/result.go 出参定义。
//
// Author: Charlie

package resource

// TreeNode 资源树节点。
//
// Author: Charlie
type TreeNode struct {
	Resource
	Children []TreeNode `json:"children"`
}

// ModuleOption 资源模块选择项。
//
// Author: Charlie
type ModuleOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// GrantModule 资源授权模块选项（含模块下启用资源）。
//
// Author: Charlie
type GrantModule struct {
	ModuleID  string     `json:"module_id"`
	Name      string     `json:"name"`
	Resources []Resource `json:"resources"`
}

// PermissionItem 已注册权限项。
//
// Author: Charlie
type PermissionItem struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
