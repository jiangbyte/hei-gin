package client

// TreeNode 客户端资源树节点。
//
// Author: Charlie
type TreeNode struct {
	ClientResource
	Children []TreeNode `json:"children"`
}

// ModuleOption 客户端模块选择项。
//
// Author: Charlie
type ModuleOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// GrantModule 客户端资源授权模块选项（含模块下启用资源）。
//
// Author: Charlie
type GrantModule struct {
	ModuleID  string           `json:"module_id"`
	Name      string           `json:"name"`
	Resources []ClientResource `json:"resources"`
}
