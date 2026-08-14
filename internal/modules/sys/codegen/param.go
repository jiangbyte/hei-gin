package codegen

import "hei-gin/internal/framework/core/schema"

// AddParam åˆ›å»ºä»£ç ç”Ÿæˆæ–¹æ¡ˆå…¥å‚ã€‚
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
	TreeParentField  *string `json:"tree_parent_field"`
	TreeLabelField   *string `json:"tree_label_field"`
	SubTable         *string `json:"sub_table"`
	SubPK            *string `json:"sub_pk"`
	SubForeignKey    *string `json:"sub_foreign_key"`
	SubEntityName    *string `json:"sub_entity_name"`
	SubBusinessName  *string `json:"sub_business_name"`
}

// EditParam æ›´æ–°ä»£ç ç”Ÿæˆæ–¹æ¡ˆå…¥å‚ï¼ˆID + AddParamï¼‰ã€‚
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam ä»£ç ç”Ÿæˆæ–¹æ¡ˆåˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Name      string `form:"name"`
	GenType   string `form:"gen_type"`
	MainTable string `form:"main_table"`
}

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// FieldUpdateItemParam å•æ¡ä»£ç ç”Ÿæˆå­—æ®µæ›´æ–°é¡¹ã€‚
//
// Author: Charlie
type FieldUpdateItemParam struct {
	ID             string  `json:"id"`
	TableRole      string  `json:"table_role" binding:"required"`
	ColumnName     string  `json:"column_name" binding:"required"`
	ColumnComment  *string `json:"column_comment"`
	DBType         string  `json:"db_type" binding:"required"`
	PythonType     string  `json:"python_type"`
	TypescriptType string  `json:"typescript_type"`
	FormWidget     string  `json:"form_widget"`
	DictCode       *string `json:"dict_code"`
	QueryOperator  *string `json:"query_operator"`
	ShowInTable    bool    `json:"show_in_table"`
	ShowInForm     bool    `json:"show_in_form"`
	ShowInDetail   bool    `json:"show_in_detail"`
	ShowInQuery    bool    `json:"show_in_query"`
	IsPrimaryKey   bool    `json:"is_primary_key"`
	IsRequired     bool    `json:"is_required"`
	IsUnique       bool    `json:"is_unique"`
	IsNullable     bool    `json:"is_nullable"`
	MaxLength      *int    `json:"max_length"`
	Sort           int     `json:"sort"`
}

// FieldsUpdateBatchParam æ‰¹é‡æ›´æ–°ä»£ç ç”Ÿæˆå­—æ®µé…ç½®ã€‚
//
// Author: Charlie
type FieldsUpdateBatchParam struct {
	PlanID string                 `json:"plan_id" binding:"required"`
	Fields []FieldUpdateItemParam `json:"fields" binding:"required"`
}

// FieldQuery å­—æ®µæŸ¥è¯¢å…¥å‚ã€‚
//
// Author: Charlie
type FieldQuery struct {
	PlanID    string `form:"plan_id" binding:"required"`
	TableRole string `form:"table_role"`
}
