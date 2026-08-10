package codegen

import "hei-gin/framework/core/schema"

// AddParam 创建代码生成方案入参。
//
// Author: Charlie
type AddParam struct {
	Name             string  `json:"name" binding:"required"`
	GenType          string  `json:"gen_type" binding:"required"`
	Author           string  `json:"author" binding:"required"`
	Description      *string `json:"description"`
	MainTable        string  `json:"main_table" binding:"required"`
	MainPK           string  `json:"main_pk"`
	MainEntityName   string  `json:"main_entity_name" binding:"required"`
	MainModulePath   string  `json:"main_module_path" binding:"required"`
	MainBusinessName string  `json:"main_business_name" binding:"required"`
	APIPrefix        string  `json:"api_prefix" binding:"required"`
	PermissionPrefix string  `json:"permission_prefix" binding:"required"`
	ResourceModuleID *string `json:"resource_module_id"`
	ParentResourceID *string `json:"parent_resource_id"`
	MenuName         string  `json:"menu_name" binding:"required"`
	MenuPath         string  `json:"menu_path" binding:"required"`
	ComponentPath    string  `json:"component_path" binding:"required"`
	Icon             *string `json:"icon"`
	Sort             int     `json:"sort"`
}

// PageParam 代码生成方案分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Name      string `form:"name"`
	GenType   string `form:"gen_type"`
	MainTable string `form:"main_table"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
