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
	err = db.Order("sort asc, id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
