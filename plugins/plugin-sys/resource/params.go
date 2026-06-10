package resource

import (
	"hei-gin/sdk/utils"
)


// ModuleVO represents a system module view object.
type ModuleVO struct {
	ID          string     `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Category    string     `json:"category"`
	Icon        *string    `json:"icon"`
	Color       *string    `json:"color"`
	Description *string    `json:"description"`
	IsVisible   string     `json:"is_visible"`
	Status      string     `json:"status"`
	SortCode    int        `json:"sort_code"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   *string    `json:"created_by"`
	UpdatedAt   string `json:"updated_at"`
	UpdatedBy   *string    `json:"updated_by"`
}

// ModulePageParam holds pagination parameters for module page queries.
type ModulePageParam struct {
	Current int `json:"current" form:"current"`
	Size    int `json:"size" form:"size"`
}

// ResourceVO represents a system resource view object with optional children for tree rendering.
type ResourceVO struct {
	ID            string        `json:"id"`
	Code          string        `json:"code"`
	Name          string        `json:"name"`
	Category      string        `json:"category"`
	Type          string        `json:"type"`
	Description   *string       `json:"description"`
	ParentID      *string       `json:"parent_id"`
	RoutePath     *string       `json:"route_path"`
	ComponentPath *string       `json:"component_path"`
	RedirectPath  *string       `json:"redirect_path"`
	Icon          *string       `json:"icon"`
	Color         *string       `json:"color"`
	IsVisible     string        `json:"is_visible"`
	IsCache       string        `json:"is_cache"`
	IsAffix       string        `json:"is_affix"`
	IsBreadcrumb  string        `json:"is_breadcrumb"`
	ExternalURL   *string       `json:"external_url"`
	Extra         *string       `json:"extra"`
	Status        string        `json:"status"`
	SortCode      int           `json:"sort_code"`
	CreatedAt     string    `json:"created_at"`
	CreatedBy     *string       `json:"created_by"`
	UpdatedAt     string    `json:"updated_at"`
	UpdatedBy     *string       `json:"updated_by"`
	Children      []*ResourceVO `json:"children"`
}

// ResourcePageParam holds pagination parameters for resource page queries.
type ResourcePageParam struct {
	Current int `json:"current" form:"current"`
	Size    int `json:"size" form:"size"`
}

// ===== Internal types =====

type relRoleResource struct {
	ID         string
	RoleID     string
	ResourceID string
}

type relRolePermission struct {
	ID             string
	RoleID         string
	PermissionCode string
}

// ===== Converters =====

func toResourceVO(r *SysResource) *ResourceVO {
	return &ResourceVO{
		ID: r.ID, Code: r.Code, Name: r.Name, Category: r.Category, Type: r.Type,
		Description: r.Description, ParentID: r.ParentID, RoutePath: r.RoutePath,
		ComponentPath: r.ComponentPath, RedirectPath: r.RedirectPath, Icon: r.Icon,
		Color: r.Color, IsVisible: r.IsVisible, IsCache: r.IsCache, IsAffix: r.IsAffix,
		IsBreadcrumb: r.IsBreadcrumb, ExternalURL: r.ExternalURL, Extra: r.Extra,
		Status: r.Status, SortCode: r.SortCode, CreatedAt: utils.FormatDateTimePtr(r.CreatedAt), CreatedBy: r.CreatedBy,
		UpdatedAt: utils.FormatDateTimePtr(r.UpdatedAt), UpdatedBy: r.UpdatedBy,
	}
}
