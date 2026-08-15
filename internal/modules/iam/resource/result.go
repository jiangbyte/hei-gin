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

// GrantModule 资源授权模块选项（对齐 hei-boot SysResourceGrantModuleOptionResult：模块 → 菜单 → 按钮权限树）。
//
// Author: Charlie
type GrantModule struct {
	ID    string            `json:"id"`
	Title string            `json:"title"`
	Menu  []GrantMenuOption `json:"menu"`
}

// GrantMenuOption 资源授权菜单选项节点（对齐 hei-boot SysResourceGrantMenuOptionResult）。
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

// PermissionOption 资源可绑定权限选项（对齐 hei-boot SysResourcePermissionOptionResult）。
//
// Author: Charlie
type PermissionOption struct {
	ID            string `json:"id"`
	PermissionKey string `json:"permission_key"`
	Title         string `json:"title"`
	DataScope     string `json:"data_scope"`
}

// PermissionItem 已注册权限项（对齐 hei-boot SysRegisteredPermissionResult）。
//
// Author: Charlie
type PermissionItem struct {
	PermissionKey string `json:"permission_key"`
	Name          string `json:"name"`
	ModuleCode    string `json:"module_code"`
	ResourceCode  string `json:"resource_code"`
	Action        string `json:"action"`
}

// ButtonResult 按钮资源行（对齐 hei-boot SysResourceButtonResult：含权限绑定字段）。
//
// Author: Charlie
type ButtonResult struct {
	ID                    string   `json:"id"`
	ParentID              *string  `json:"parent_id"`
	ParentIDName          *string  `json:"parent_id_name"`
	Code                  string   `json:"code"`
	Name                  string   `json:"name"`
	PermissionKey         string   `json:"permission_key"`
	PermissionDescription *string  `json:"permission_description"`
	DataScope             string   `json:"data_scope"`
	CustomScopeDeptIDs    []string `json:"custom_scope_dept_ids"`
	ModuleID              *string  `json:"module_id"`
	ModuleIDName          *string  `json:"module_id_name"`
	ModuleClient          *string  `json:"module_client"`
	Sort                  int      `json:"sort"`
	Status                string   `json:"status"`
	Description           *string  `json:"description"`
	CreatedAt             string   `json:"created_at"`
	UpdatedAt             string   `json:"updated_at"`
}
