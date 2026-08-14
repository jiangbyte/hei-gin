// internal/modules/iam/resource/param.go 入参定义。
//
// Author: Charlie

package resource

// ResourceAddParam 创建资源入参。
//
// Author: Charlie
type ResourceAddParam struct {
	ParentID     *string `json:"parent_id"`
	Code         string  `json:"code" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	ResourceType string  `json:"resource_type" binding:"required"`
	ModuleID     *string `json:"module_id"`
	Path         *string `json:"path"`
	Component    *string `json:"component"`
	Redirect     *string `json:"redirect"`
	Icon         *string `json:"icon"`
	Color        *string `json:"color"`
	Href         *string `json:"href"`
	Sort         int     `json:"sort"`
	IsVisible    *bool   `json:"is_visible"`
	IsCache      bool    `json:"is_cache"`
	IsAffix      bool    `json:"is_affix"`
	Status       string  `json:"status"`
	Description  *string `json:"description"`
	Layout       *string `json:"layout"`
}

// ResourceEditParam 更新资源入参。
//
// Author: Charlie
type ResourceEditParam struct {
	ID string `json:"id" binding:"required"`
	ResourceAddParam
}

// ResourcePageParam 资源分页查询。
//
// Author: Charlie
type ResourcePageParam struct {
	Current  int    `form:"current" json:"current"`
	Size     int    `form:"size" json:"size"`
	Name     string `form:"name" json:"name"`
	Code     string `form:"code" json:"code"`
	ModuleID string `form:"module_id" json:"module_id"`
	Status   string `form:"status" json:"status"`
}

// Normalize 分页规范化。
func (q ResourcePageParam) Normalize() (current, size int) {
	current, size = q.Current, q.Size
	if current < 1 {
		current = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return current, size
}

// ModuleAddParam 创建资源模块入参。
//
// Author: Charlie
type ModuleAddParam struct {
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Client      string  `json:"client"`
	Icon        *string `json:"icon"`
	Color       *string `json:"color"`
	Sort        int     `json:"sort"`
	Status      string  `json:"status"`
	Description *string `json:"description"`
}

// ModuleEditParam 更新资源模块入参。
//
// Author: Charlie
type ModuleEditParam struct {
	ID string `json:"id" binding:"required"`
	ModuleAddParam
}

// ModulePageParam 资源模块分页查询。
//
// Author: Charlie
type ModulePageParam struct {
	Current int    `form:"current" json:"current"`
	Size    int    `form:"size" json:"size"`
	Name    string `form:"name" json:"name"`
	Client  string `form:"client" json:"client"`
	Status  string `form:"status" json:"status"`
}

// Normalize 分页规范化。
func (q ModulePageParam) Normalize() (current, size int) {
	current, size = q.Current, q.Size
	if current < 1 {
		current = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return current, size
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// ResourcePermissionBindParam 资源绑定权限入参（管理端/客户端）。
//
// Author: Charlie
type ResourcePermissionBindParam struct {
	ResourceID     string   `json:"resource_id" binding:"required"`
	PermissionKeys []string `json:"permission_keys"`
	AccountType    string   `json:"account_type"`
}

// ButtonAddParam 创建按钮资源入参。
//
// Author: Charlie
type ButtonAddParam struct {
	ResourceID  string  `json:"resource_id" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Sort        int     `json:"sort"`
	Status      string  `json:"status"`
	Description *string `json:"description"`
}

// ButtonEditParam 更新按钮资源入参。
//
// Author: Charlie
type ButtonEditParam struct {
	ID string `json:"id" binding:"required"`
	ButtonAddParam
}

// ButtonPageParam 按钮资源分页查询。
//
// Author: Charlie
type ButtonPageParam struct {
	Current    int    `form:"current" json:"current"`
	Size       int    `form:"size" json:"size"`
	ResourceID string `form:"resource_id" json:"resource_id"`
	Code       string `form:"code" json:"code"`
	Name       string `form:"name" json:"name"`
	Status     string `form:"status" json:"status"`
}

// Normalize 分页规范化。
func (q ButtonPageParam) Normalize() (current, size int) {
	current, size = q.Current, q.Size
	if current < 1 {
		current = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return current, size
}
