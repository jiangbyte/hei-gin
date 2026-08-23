// internal/modules/iam/group/repo.go 持久化仓储。
//
// Author: Charlie

package group

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/core/security/datascope"
	"hei-gin/internal/framework/platform/db/dialect"
)

// Repo 用户组持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建用户组。
func (r *Repo) Create(ctx context.Context, row *Group) error {
	return r.with(ctx).Create(row).Error
}

// Update 更新用户组。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Group{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByIDs 批量删除。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Group{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*Group, error) {
	var row Group
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// GetByIDs 按 ID 集合查询（保持入参顺序）。
func (r *Repo) GetByIDs(ctx context.Context, ids []string) ([]Group, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []Group
	if err := r.with(ctx).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	byID := make(map[string]Group, len(rows))
	for _, row := range rows {
		byID[row.ID] = row
	}
	out := make([]Group, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			out = append(out, row)
		}
	}
	return out, nil
}

// Page 分页查询（sess 非空时按数据范围过滤；对齐 hei-boot applyOwnerOrDept 以 owner_dept_id 为归属列）。
func (r *Repo) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []Group, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Group{})
	if sess != nil {
		db = datascope.ApplyKey(db, sess, "iam:group:page", "owner_dept_id", "created_by")
	}
	if p.Name != "" {
		db = db.Where(dialect.ILike(db, "name"), "%"+p.Name+"%")
	}
	if p.Status != "" {
		db = db.Where("status = ?", p.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
