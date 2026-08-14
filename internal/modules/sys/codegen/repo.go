package codegen

import (
	"context"

	"gorm.io/gorm"
)

// Repo 代码生成持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建方案。
func (r *Repo) Create(ctx context.Context, row *Plan) error {
	return r.with(ctx).Create(row).Error
}

// Update 更新方案。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Plan{}).Where("id = ?", id).Updates(updates).Error
}

// CountByName 按名称统计（排除指定 ID，用于重名校验）。
func (r *Repo) CountByName(ctx context.Context, name, excludeID string) (int64, error) {
	db := r.with(ctx).Model(&Plan{}).Where("name = ?", name)
	if excludeID != "" {
		db = db.Where("id <> ?", excludeID)
	}
	var total int64
	err := db.Count(&total).Error
	return total, err
}

// DeleteByIDs 事务删除方案及字段。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plan_id IN ?", ids).Delete(&Field{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", ids).Delete(&Plan{}).Error
	})
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*Plan, error) {
	var row Plan
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Page 分页查询。
func (r *Repo) Page(ctx context.Context, q PageParam) (rows []Plan, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&Plan{})
	if q.Name != "" {
		db = db.Where("name ILIKE ?", "%"+q.Name+"%")
	}
	if q.GenType != "" {
		db = db.Where("gen_type = ?", q.GenType)
	}
	if q.MainTable != "" {
		db = db.Where("main_table ILIKE ?", "%"+q.MainTable+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("updated_at desc, id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// ListFields 查询方案字段（可选 table_role 过滤，按 role/sort/id 排序）。
func (r *Repo) ListFields(ctx context.Context, planID, tableRole string) ([]Field, error) {
	db := r.with(ctx).Where("plan_id = ?", planID)
	if tableRole != "" {
		db = db.Where("table_role = ?", tableRole)
	}
	var rows []Field
	err := db.Order("table_role asc, sort asc, id asc").Find(&rows).Error
	return rows, err
}

// ReplaceFields 事务重建方案字段（先删后插）。
func (r *Repo) ReplaceFields(ctx context.Context, planID string, fields []Field) error {
	return r.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plan_id = ?", planID).Delete(&Field{}).Error; err != nil {
			return err
		}
		if len(fields) == 0 {
			return nil
		}
		return tx.Create(&fields).Error
	})
}

// TableRow information_schema 表摘要。
//
// Author: Charlie
type TableRow struct {
	TableName    string `gorm:"column:table_name"`
	TableComment string `gorm:"column:table_comment"`
}

// ColumnRow information_schema 列摘要。
//
// Author: Charlie
type ColumnRow struct {
	ColumnName    string `gorm:"column:column_name"`
	ColumnComment string `gorm:"column:column_comment"`
	DataType      string `gorm:"column:data_type"`
	UDTName       string `gorm:"column:udt_name"`
	IsNullable    string `gorm:"column:is_nullable"`
	MaxLength     *int   `gorm:"column:max_length"`
	Sort          int    `gorm:"column:sort"`
	IsPrimaryKey  bool   `gorm:"column:is_primary_key"`
}

// ListTables 列出 public schema 下的基表（排除代码生成自身表）。
func (r *Repo) ListTables(ctx context.Context) ([]TableRow, error) {
	var rows []TableRow
	err := r.with(ctx).Raw(`
SELECT c.relname AS table_name,
       COALESCE(obj_description(c.oid), '') AS table_comment
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = current_schema()
  AND c.relkind = 'r'
  AND c.relname NOT IN ('sys_codegen_plan', 'sys_codegen_field')
ORDER BY c.relname`).Scan(&rows).Error
	return rows, err
}

// ListColumns 列出表列元数据（含主键标记）。
func (r *Repo) ListColumns(ctx context.Context, tableName string) ([]ColumnRow, error) {
	var rows []ColumnRow
	err := r.with(ctx).Raw(`
SELECT c.column_name,
       COALESCE(pgd.description, '') AS column_comment,
       c.data_type,
       c.udt_name,
       c.character_maximum_length AS max_length,
       c.is_nullable,
       c.ordinal_position AS sort,
       EXISTS (
         SELECT 1
         FROM information_schema.table_constraints tc
         JOIN information_schema.key_column_usage kcu
           ON tc.constraint_name = kcu.constraint_name
          AND tc.table_schema = kcu.table_schema
         WHERE tc.table_schema = c.table_schema
           AND tc.table_name = c.table_name
           AND tc.constraint_type = 'PRIMARY KEY'
           AND kcu.column_name = c.column_name
       ) AS is_primary_key
FROM information_schema.columns c
LEFT JOIN pg_catalog.pg_statio_all_tables st
  ON st.schemaname = c.table_schema AND st.relname = c.table_name
LEFT JOIN pg_catalog.pg_description pgd
  ON pgd.objoid = st.relid AND pgd.objsubid = c.ordinal_position
WHERE c.table_schema = current_schema()
  AND c.table_name = ?
ORDER BY c.ordinal_position`, tableName).Scan(&rows).Error
	return rows, err
}
