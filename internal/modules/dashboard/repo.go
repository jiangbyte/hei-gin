// internal/modules/dashboard/repo.go 持久化仓储。
//
// Author: Charlie

package dashboard

import (
	"context"

	"gorm.io/gorm"
)

// Repo 仪表盘统计持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Count 统计表行数。
func (r *Repo) Count(ctx context.Context, table string) int64 {
	var n int64
	_ = r.with(ctx).Table(table).Count(&n).Error
	return n
}

// CountWhere 按条件统计表行数。
func (r *Repo) CountWhere(ctx context.Context, table, where string, args ...any) int64 {
	var n int64
	_ = r.with(ctx).Table(table).Where(where, args...).Count(&n).Error
	return n
}
