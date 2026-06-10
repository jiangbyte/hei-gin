package resource

import "hei-gin/sdk/utils"

// SysModuleToModuleVO 将 resource.SysModule 映射到 resource.ModuleVO
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

	// *time.Time → string manual conversion
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)

	return dst
}

// ModuleVOToSysModule 将 resource.ModuleVO 映射到 resource.SysModule
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

	// string → *time.Time manual conversion
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)

	return dst
}

// SysResourceToResourceVO 将 resource.SysResource 映射到 resource.ResourceVO
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

	// *time.Time → string manual conversion
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)

	return dst
}

// ResourceVOToSysResource 将 resource.ResourceVO 映射到 resource.SysResource
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

	// string → *time.Time manual conversion
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)

	return dst
}
