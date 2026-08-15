// internal/modules/iam/client/result.go 出参定义。
//
// Author: Charlie

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

// GrantModule 客户端资源授权模块选项（对齐 hei-boot ClientResourceServiceImpl.listGrantModules：模块 → 菜单 → 按钮权限树）。
//
// Author: Charlie
type GrantModule struct {
	ID    string            `json:"id"`
	Title string            `json:"title"`
	Menu  []GrantMenuOption `json:"menu"`
}

// GrantMenuOption 客户端资源授权菜单选项节点（对齐 hei-boot SysResourceGrantMenuOptionResult）。
//
// Author: Charlie
type GrantMenuOption struct {
	ID           string             `json:"id"`
	ModuleID     string             `json:"module_id"`
	ParentID     *string            `json:"parent_id"`
	ParentIDName string             `json:"parent_id_name"`
	Title        string             `json:"title"`
	Button       []PermissionOption `json:"button"`
}

// PermissionOption 客户端资源可绑定权限选项（对齐 hei-boot SysResourcePermissionOptionResult）。
//
// Author: Charlie
type PermissionOption struct {
	ID            string `json:"id"`
	PermissionKey string `json:"permission_key"`
	Title         string `json:"title"`
	DataScope     string `json:"data_scope"`
}
