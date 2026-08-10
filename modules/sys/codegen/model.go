// Package codegen 提供低代码/代码生成方案与字段配置。
package codegen

import "time"

// Plan 代码生成方案实体，对应表 sys_codegen_plan。
//
// Author: Charlie
type Plan struct {
	ID               string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	Name             string    `gorm:"column:name;size:128;uniqueIndex;not null" json:"name"`
	GenType          string    `gorm:"column:gen_type;size:32;not null" json:"gen_type"`
	Author           string    `gorm:"column:author;size:64;not null" json:"author"`
	Description      *string   `gorm:"column:description" json:"description"`
	MainTable        string    `gorm:"column:main_table;size:128;not null" json:"main_table"`
	MainPK           string    `gorm:"column:main_pk;size:128;not null;default:id" json:"main_pk"`
	MainEntityName   string    `gorm:"column:main_entity_name;size:128;not null" json:"main_entity_name"`
	MainModulePath   string    `gorm:"column:main_module_path;size:255;not null" json:"main_module_path"`
	MainBusinessName string    `gorm:"column:main_business_name;size:128;not null" json:"main_business_name"`
	APIPrefix        string    `gorm:"column:api_prefix;size:255;not null" json:"api_prefix"`
	PermissionPrefix string    `gorm:"column:permission_prefix;size:128;not null" json:"permission_prefix"`
	ResourceModuleID *string   `gorm:"column:resource_module_id;size:64" json:"resource_module_id"`
	ParentResourceID *string   `gorm:"column:parent_resource_id;size:64" json:"parent_resource_id"`
	MenuName         string    `gorm:"column:menu_name;size:64;not null" json:"menu_name"`
	MenuPath         string    `gorm:"column:menu_path;size:255;not null" json:"menu_path"`
	ComponentPath    string    `gorm:"column:component_path;size:255;not null" json:"component_path"`
	Icon             *string   `gorm:"column:icon;size:255" json:"icon"`
	Sort             int       `gorm:"column:sort;not null;default:99" json:"sort"`
	TreeParentField  *string   `gorm:"column:tree_parent_field;size:128" json:"tree_parent_field"`
	TreeLabelField   *string   `gorm:"column:tree_label_field;size:128" json:"tree_label_field"`
	SubTable         *string   `gorm:"column:sub_table;size:128" json:"sub_table"`
	SubPK            *string   `gorm:"column:sub_pk;size:128" json:"sub_pk"`
	SubForeignKey    *string   `gorm:"column:sub_foreign_key;size:128" json:"sub_foreign_key"`
	SubEntityName    *string   `gorm:"column:sub_entity_name;size:128" json:"sub_entity_name"`
	SubBusinessName  *string   `gorm:"column:sub_business_name;size:128" json:"sub_business_name"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy        *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy        *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 Plan 对应的数据库表名。
func (Plan) TableName() string { return "sys_codegen_plan" }

// Field 代码生成字段实体，对应表 sys_codegen_field。
//
// Author: Charlie
type Field struct {
	ID             string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	PlanID         string    `gorm:"column:plan_id;size:64;not null;index" json:"plan_id"`
	TableRole      string    `gorm:"column:table_role;size:16;not null;default:MAIN" json:"table_role"`
	ColumnName     string    `gorm:"column:column_name;size:128;not null" json:"column_name"`
	ColumnComment  *string   `gorm:"column:column_comment;size:255" json:"column_comment"`
	DBType         string    `gorm:"column:db_type;size:128;not null" json:"db_type"`
	PythonType     string    `gorm:"column:python_type;size:64;not null;default:str" json:"python_type"`
	TypescriptType string    `gorm:"column:typescript_type;size:64;not null;default:string" json:"typescript_type"`
	FormWidget     string    `gorm:"column:form_widget;size:32;not null;default:input" json:"form_widget"`
	DictCode       *string   `gorm:"column:dict_code;size:128" json:"dict_code"`
	QueryOperator  *string   `gorm:"column:query_operator;size:32" json:"query_operator"`
	ShowInTable    bool      `gorm:"column:show_in_table;not null;default:true" json:"show_in_table"`
	ShowInForm     bool      `gorm:"column:show_in_form;not null;default:true" json:"show_in_form"`
	ShowInDetail   bool      `gorm:"column:show_in_detail;not null;default:true" json:"show_in_detail"`
	ShowInQuery    bool      `gorm:"column:show_in_query;not null;default:false" json:"show_in_query"`
	IsPrimaryKey   bool      `gorm:"column:is_primary_key;not null;default:false" json:"is_primary_key"`
	IsRequired     bool      `gorm:"column:is_required;not null;default:false" json:"is_required"`
	IsUnique       bool      `gorm:"column:is_unique;not null;default:false" json:"is_unique"`
	IsNullable     bool      `gorm:"column:is_nullable;not null;default:true" json:"is_nullable"`
	MaxLength      *int      `gorm:"column:max_length" json:"max_length"`
	Sort           int       `gorm:"column:sort;not null;default:99" json:"sort"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy      *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy      *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 Field 对应的数据库表名。
func (Field) TableName() string { return "sys_codegen_field" }
