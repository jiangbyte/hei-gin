package codegen

// DatabaseTableResult 数据库表摘要（代码生成表列表）。
//
// Author: Charlie
type DatabaseTableResult struct {
	TableName    string `json:"table_name"`
	TableComment string `json:"table_comment"`
}

// DatabaseColumnResult 数据库列元数据。
//
// Author: Charlie
type DatabaseColumnResult struct {
	ColumnName       string `json:"column_name"`
	ColumnComment    string `json:"column_comment"`
	DBType           string `json:"db_type"`
	PythonType       string `json:"python_type"`
	TypescriptType   string `json:"typescript_type"`
	IsPrimaryKey     bool   `json:"is_primary_key"`
	IsNullable       bool   `json:"is_nullable"`
	MaxLength        *int   `json:"max_length"`
}

// PreviewFileResult 代码生成预览单文件结果。
//
// Author: Charlie
type PreviewFileResult struct {
	Path     string `json:"path"`
	Language string `json:"language"`
	Content  string `json:"content"`
}

// PreviewResult 代码生成预览结果。
//
// Author: Charlie
type PreviewResult struct {
	Files []PreviewFileResult `json:"files"`
}
