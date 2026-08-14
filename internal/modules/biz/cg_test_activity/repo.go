// internal/modules/biz/cg_test_activity/repo.go 持久化仓储。
//
// Author: Charlie

package cg_test_activity

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/core/security/datascope"
)

// Repo æ´»åŠ¨æŒä¹…åŒ–ã€‚
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo æž„é€  Repoã€‚
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create åˆ›å»ºæ´»åŠ¨ã€‚
func (r *Repo) Create(ctx context.Context, row *Activity) error {
	return r.with(ctx).Create(row).Error
}

// Update æ›´æ–°æ´»åŠ¨ã€‚
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Activity{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByIDs æ‰¹é‡åˆ é™¤ã€‚
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Activity{}).Error
}

// GetByID æŒ‰ä¸»é”®æŸ¥è¯¢ã€‚
func (r *Repo) GetByID(ctx context.Context, id string) (*Activity, error) {
	var row Activity
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Page åˆ†é¡µæŸ¥è¯¢ï¼›sess éžç©ºæ—¶æŒ‰ owner_dept_id æ•°æ®èŒƒå›´è¿‡æ»¤ã€‚
func (r *Repo) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []Activity, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Activity{})
	if sess != nil {
		db = datascope.Apply(db, sess, "owner_dept_id")
	}
	if p.Code != "" {
		db = db.Where("code ILIKE ?", "%"+p.Code+"%")
	}
	if p.Name != "" {
		db = db.Where("name ILIKE ?", "%"+p.Name+"%")
	}
	if p.Category != "" {
		db = db.Where("category = ?", p.Category)
	}
	if p.Type != "" {
		db = db.Where("type = ?", p.Type)
	}
	if p.Status != "" {
		db = db.Where("status = ?", p.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("created_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
