// Package dialect 提供 MySQL / PostgreSQL 方言 SQL 片段助手。
//
// Author: Charlie
package dialect

import (
	"fmt"

	"gorm.io/gorm"
)

const (
	MySQL    = "mysql"
	Postgres = "postgres"
)

// Name 返回 dialector 名称（mysql / postgres）。
func Name(db *gorm.DB) string {
	if db == nil || db.Dialector == nil {
		return ""
	}
	return db.Dialector.Name()
}

// IsMySQL 是否 MySQL。
func IsMySQL(db *gorm.DB) bool { return Name(db) == MySQL }

// IsPostgres 是否 PostgreSQL。
func IsPostgres(db *gorm.DB) bool { return Name(db) == Postgres }

// JSONContainsString 判断 JSON 数组/对象是否包含字符串标量（一个 ? 占位符）。
func JSONContainsString(db *gorm.DB, column string) string {
	col := MustIdent(column)
	if IsMySQL(db) {
		return fmt.Sprintf("JSON_CONTAINS(%s, JSON_QUOTE(?), '$')", col)
	}
	// 列类型为 json，查询时再转为 jsonb 以使用 jsonb_exists
	return fmt.Sprintf("jsonb_exists((%s)::jsonb, ?)", col)
}

// JSONArrayEmptyOrNull 判断 JSON 数组为空或 NULL（无占位符）。
func JSONArrayEmptyOrNull(db *gorm.DB, column string) string {
	col := MustIdent(column)
	if IsMySQL(db) {
		return fmt.Sprintf("(%s IS NULL OR JSON_LENGTH(IFNULL(%s, JSON_ARRAY())) = 0)", col, col)
	}
	return fmt.Sprintf("(%s IS NULL OR jsonb_array_length(COALESCE(%s::jsonb, '[]'::jsonb)) = 0)", col, col)
}

// DayBucket 返回按日聚合的 SELECT 表达式与 GROUP BY 表达式。
func DayBucket(db *gorm.DB, column string) (selectExpr, groupBy string) {
	col := MustIdent(column)
	if IsMySQL(db) {
		return fmt.Sprintf("DATE_FORMAT(DATE(%s), '%%Y-%%m-%%d')", col), fmt.Sprintf("DATE(%s)", col)
	}
	return fmt.Sprintf("TO_CHAR(%s::date, 'YYYY-MM-DD')", col), fmt.Sprintf("%s::date", col)
}
