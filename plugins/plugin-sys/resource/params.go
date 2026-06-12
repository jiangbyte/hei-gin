package resource

import "hei-gin/sdk/utils"

// ModuleVO represents a system module view object.
type ModuleVO struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Icon        *string `json:"icon"`
	Color       *string `json:"color"`
	Description *string `json:"description"`
	IsVisible   string  `json:"is_visible"`
	Status      string  `json:"status"`
	SortCode    int     `json:"sort_code"`
	CreatedAt   string  `json:"created_at"`
	CreatedBy   *string `json:"created_by"`
	UpdatedAt   string  `json:"updated_at"`
	UpdatedBy   *string `json:"updated_by"`
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
	CreatedAt     string        `json:"created_at"`
	CreatedBy     *string       `json:"created_by"`
	UpdatedAt     string        `json:"updated_at"`
	UpdatedBy     *string       `json:"updated_by"`
	Children      []*ResourceVO `json:"children"`
}

// ResourcePageParam holds pagination parameters for resource page queries.
type ResourcePageParam struct {
	Current int `json:"current" form:"current"`
	Size    int `json:"size" form:"size"`
}

func SysModuleToModuleVO(src *SysModule) *ModuleVO {
	if src == nil {
		return nil
	}

	dst := &ModuleVO{}
	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.Icon = src.Icon
	dst.Color = src.Color
	dst.Description = src.Description
	dst.IsVisible = src.IsVisible
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
	return dst
}

func ModuleVOToSysModule(src *ModuleVO) *SysModule {
	if src == nil {
		return nil
	}

	dst := &SysModule{}
	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.Icon = src.Icon
	dst.Color = src.Color
	dst.Description = src.Description
	dst.IsVisible = src.IsVisible
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)
	return dst
}

func SysResourceToResourceVO(src *SysResource) *ResourceVO {
	if src == nil {
		return nil
	}

	dst := &ResourceVO{}
	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.Type = src.Type
	dst.Description = src.Description
	dst.ParentID = src.ParentID
	dst.RoutePath = src.RoutePath
	dst.ComponentPath = src.ComponentPath
	dst.RedirectPath = src.RedirectPath
	dst.Icon = src.Icon
	dst.Color = src.Color
	dst.IsVisible = src.IsVisible
	dst.IsCache = src.IsCache
	dst.IsAffix = src.IsAffix
	dst.IsBreadcrumb = src.IsBreadcrumb
	dst.ExternalURL = src.ExternalURL
	dst.Extra = src.Extra
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
	return dst
}

func ResourceVOToSysResource(src *ResourceVO) *SysResource {
	if src == nil {
		return nil
	}

	dst := &SysResource{}
	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.Type = src.Type
	dst.Description = src.Description
	dst.ParentID = src.ParentID
	dst.RoutePath = src.RoutePath
	dst.ComponentPath = src.ComponentPath
	dst.RedirectPath = src.RedirectPath
	dst.Icon = src.Icon
	dst.Color = src.Color
	dst.IsVisible = src.IsVisible
	dst.IsCache = src.IsCache
	dst.IsAffix = src.IsAffix
	dst.IsBreadcrumb = src.IsBreadcrumb
	dst.ExternalURL = src.ExternalURL
	dst.Extra = src.Extra
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)
	return dst
}
