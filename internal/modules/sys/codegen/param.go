// internal/modules/sys/codegen/param.go 入参定义（对齐 hei-boot）。
//
// Author: Charlie

package codegen

import "hei-gin/internal/framework/core/schema"

// AddParam 创建代码生成方案入参。
//
// Author: Charlie
type AddParam struct {
	Name             string  `json:"name" binding:"required"`
	GenType          string  `json:"gen_type" binding:"required"`
	Author           string  `json:"author" binding:"required"`
	Description      *string `json:"description"`
	Table            string  `json:"table_name" binding:"required"`
	PKColumn         string  `json:"pk_column"`
	EntityName       string  `json:"entity_name" binding:"required"`
	ModulePath       string  `json:"module_path" binding:"required"`
	BusinessName     string  `json:"business_name" binding:"required"`
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

// EditParam 更新代码生成方案入参（ID + AddParam）。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam 代码生成方案分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Name      string `form:"name"`
	GenType   string `form:"gen_type"`
	Table string `form:"table_name"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// FieldUpdateItemParam 单条代码生成字段更新项。
//
// Author: Charlie
type FieldUpdateItemParam struct {
	ID            string  `json:"id"`
	TableRole     string  `json:"table_role" binding:"required"`
	ColumnName    string  `json:"column_name" binding:"required"`
	Label         *string `json:"label"`
	DBType        string  `json:"db_type" binding:"required"`
	ValueType     string  `json:"value_type"`
	UIType        string  `json:"ui_type"`
	Widget        string  `json:"widget"`
	DictCode      *string `json:"dict_code"`
	QueryOperator *string `json:"query_operator"`
	InTable       bool    `json:"in_table"`
	InForm        bool    `json:"in_form"`
	InDetail      bool    `json:"in_detail"`
	InQuery       bool    `json:"in_query"`
	PrimaryKey    bool    `json:"primary_key"`
	Required      bool    `json:"required"`
	UniqueFlag    bool    `json:"unique_flag"`
	Nullable      bool    `json:"nullable"`
	MaxLength     *int    `json:"max_length"`
	Sort          int     `json:"sort"`
}

// FieldsUpdateBatchParam 批量更新代码生成字段配置。
//
// Author: Charlie
type FieldsUpdateBatchParam struct {
	PlanID string                 `json:"plan_id" binding:"required"`
	Fields []FieldUpdateItemParam `json:"fields" binding:"required"`
}

// FieldQuery 字段查询入参。
//
// Author: Charlie
type FieldQuery struct {
	PlanID    string `form:"plan_id" binding:"required"`
	TableRole string `form:"table_role"`
}

// TableColumnsQuery 表字段元数据查询。
//
// Author: Charlie
type TableColumnsQuery struct {
	TableName string `form:"table_name" binding:"required"`
}
