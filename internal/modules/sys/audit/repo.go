// internal/modules/sys/audit/repo.go 持久化仓储。
//
// Author: Charlie

package audit

import (
	"context"

	"gorm.io/gorm"
)

// Repo 审计日志持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Page 分页查询审计日志。
func (r *Repo) Page(ctx context.Context, q PageParam) (rows []OperationLog, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&OperationLog{})
	if q.Module != "" {
		db = db.Where("module = ?", q.Module)
	}
	if q.Action != "" {
		db = db.Where("action = ?", q.Action)
	}
	if q.AccountID != "" {
		db = db.Where("account_id = ?", q.AccountID)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("created_at desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*OperationLog, error) {
	var row OperationLog
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// DB 返回底层 *gorm.DB（供 Job 直查）。
func (r *Repo) DB() *gorm.DB { return r.db }
