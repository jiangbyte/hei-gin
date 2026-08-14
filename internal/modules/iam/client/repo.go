// internal/modules/iam/client/repo.go 持久化仓储。
//
// Author: Charlie

package client

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

// Repo å®¢æˆ·ç«¯èµ„æºæŒä¹…åŒ–ã€‚
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo æž„é€  Repoã€‚
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// CreateModule åˆ›å»ºå®¢æˆ·ç«¯æ¨¡å—ã€‚
func (r *Repo) CreateModule(ctx context.Context, row *ClientModule) error {
	return r.with(ctx).Create(row).Error
}

// UpdateModule æ›´æ–°å®¢æˆ·ç«¯æ¨¡å—ã€‚
func (r *Repo) UpdateModule(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&ClientModule{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteModules æ‰¹é‡åˆ é™¤å®¢æˆ·ç«¯æ¨¡å—ã€‚
func (r *Repo) DeleteModules(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&ClientModule{}).Error
}

// GetModuleByID æŒ‰ä¸»é”®æŸ¥è¯¢å®¢æˆ·ç«¯æ¨¡å—ã€‚
func (r *Repo) GetModuleByID(ctx context.Context, id string) (*ClientModule, error) {
	var row ClientModule
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageModules å®¢æˆ·ç«¯æ¨¡å—åˆ†é¡µã€‚
func (r *Repo) PageModules(ctx context.Context, p ModulePageParam) (rows []ClientModule, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&ClientModule{})
	if p.Name != "" {
		db = db.Where("name ILIKE ?", "%"+p.Name+"%")
	}
	if p.AccountType != "" {
		db = db.Where("account_type = ?", p.AccountType)
	}
	if p.Status != "" {
		db = db.Where("status = ?", p.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("sort asc, id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// ListEnabledModules åˆ—å‡ºå¯ç”¨æ¨¡å—ï¼ˆå¯é€‰è´¦å·ç±»åž‹è¿‡æ»¤ï¼‰ã€‚
func (r *Repo) ListEnabledModules(ctx context.Context, accountType string) ([]ClientModule, error) {
	db := r.with(ctx).Model(&ClientModule{}).Where("status = ?", security.StatusEnabled)
	if accountType != "" {
		db = db.Where("account_type = ?", accountType)
	}
	var rows []ClientModule
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// CreateResource åˆ›å»ºå®¢æˆ·ç«¯èµ„æºã€‚
func (r *Repo) CreateResource(ctx context.Context, row *ClientResource) error {
	return r.with(ctx).Create(row).Error
}

// UpdateResource æ›´æ–°å®¢æˆ·ç«¯èµ„æºã€‚
func (r *Repo) UpdateResource(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&ClientResource{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteResources æ‰¹é‡åˆ é™¤å®¢æˆ·ç«¯èµ„æºã€‚
func (r *Repo) DeleteResources(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&ClientResource{}).Error
}

// GetResourceByID æŒ‰ä¸»é”®æŸ¥è¯¢å®¢æˆ·ç«¯èµ„æºã€‚
func (r *Repo) GetResourceByID(ctx context.Context, id string) (*ClientResource, error) {
	var row ClientResource
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageResources å®¢æˆ·ç«¯èµ„æºåˆ†é¡µã€‚
func (r *Repo) PageResources(ctx context.Context, p ResourcePageParam) (rows []ClientResource, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&ClientResource{})
	if p.Name != "" {
		db = db.Where("name ILIKE ?", "%"+p.Name+"%")
	}
	if p.ModuleID != "" {
		db = db.Where("module_id = ?", p.ModuleID)
	}
	if p.Status != "" {
		db = db.Where("status = ?", p.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("sort asc, id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// ListGrantResources åˆ—å‡ºæŒ‡å®šè´¦å·ç±»åž‹æ¨¡å—ä¸‹çš„å¯ç”¨å®¢æˆ·ç«¯èµ„æºï¼ˆæŽˆæƒæ ‘ç”¨ï¼‰ã€‚
func (r *Repo) ListGrantResources(ctx context.Context, accountType string) ([]ClientResource, error) {
	var rows []ClientResource
	err := r.with(ctx).
		Where("status = ? AND module_id IN (SELECT id FROM sys_client_module WHERE account_type = ? AND status = ?)",
			security.StatusEnabled, accountType, security.StatusEnabled).
		Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// ListResources åˆ—å‡ºå®¢æˆ·ç«¯èµ„æºï¼ˆå¯é€‰æ¨¡å—è¿‡æ»¤ï¼‰ã€‚
func (r *Repo) ListResources(ctx context.Context, moduleID string) ([]ClientResource, error) {
	db := r.with(ctx).Model(&ClientResource{})
	if moduleID != "" {
		db = db.Where("module_id = ?", moduleID)
	}
	var rows []ClientResource
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}
