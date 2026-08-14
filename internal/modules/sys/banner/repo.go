// internal/modules/sys/banner/repo.go 持久化仓储。
//
// Author: Charlie

package banner

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Repo Banner 持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

// DB 返回底层 gorm.DB。
func (r *Repo) DB() *gorm.DB { return r.db }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建 Banner。
func (r *Repo) Create(ctx context.Context, row *Banner) error {
	return r.with(ctx).Create(row).Error
}

// Update 按 ID 更新字段。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Banner{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByIDs 批量删除。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Banner{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*Banner, error) {
	var row Banner
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Page 分页查询。
func (r *Repo) Page(ctx context.Context, q PageParam) (rows []Banner, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&Banner{})
	if q.Title != "" {
		db = db.Where("title ILIKE ?", "%"+q.Title+"%")
	}
	if q.Position != "" {
		db = db.Where("position = ?", q.Position)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("sort asc, id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// List 按位置列出启用 Banner。
func (r *Repo) List(ctx context.Context, position, status string) ([]Banner, error) {
	db := r.with(ctx).Model(&Banner{}).Where("status = ?", status)
	if position != "" {
		db = db.Where("position = ?", position)
	}
	var rows []Banner
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// IncrementInteraction 互动计数 +1，返回受影响行数（0 表示 Banner 不存在）。
func (r *Repo) IncrementInteraction(ctx context.Context, id string) (int64, error) {
	res := r.with(ctx).Model(&Banner{}).Where("id = ?", id).
		UpdateColumn("interaction_count", gorm.Expr("interaction_count + 1"))
	return res.RowsAffected, res.Error
}

// ListPortal 门户端有效 Banner 列表：状态启用且在展示窗口内，按 sort 排序。
func (r *Repo) ListPortal(ctx context.Context, q PortalListParam, status string) ([]Banner, error) {
	now := time.Now()
	db := r.with(ctx).Model(&Banner{}).Where("status = ?", status)
	if q.Position != "" {
		db = db.Where("position = ?", q.Position)
	}
	if q.Category != "" {
		db = db.Where("category = ?", q.Category)
	}
	if q.Type != "" {
		db = db.Where("type = ?", q.Type)
	}
	db = db.Where("(start_at IS NULL OR start_at <= ?) AND (end_at IS NULL OR end_at >= ?)", now, now)
	var rows []Banner
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}
