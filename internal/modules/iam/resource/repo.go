// internal/modules/iam/resource/repo.go 持久化仓储。
//
// Author: Charlie

package resource

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

// Repo èµ„æºæŒä¹…åŒ–ã€‚
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo æž„é€  Repoã€‚
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// CreateResource åˆ›å»ºèµ„æºã€‚
func (r *Repo) CreateResource(ctx context.Context, row *Resource) error {
	return r.with(ctx).Create(row).Error
}

// UpdateResource æ›´æ–°èµ„æºã€‚
func (r *Repo) UpdateResource(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Resource{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteResources æ‰¹é‡åˆ é™¤èµ„æºã€‚
func (r *Repo) DeleteResources(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Resource{}).Error
}

// GetResourceByID æŒ‰ä¸»é”®æŸ¥è¯¢èµ„æºã€‚
func (r *Repo) GetResourceByID(ctx context.Context, id string) (*Resource, error) {
	var row Resource
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageResources èµ„æºåˆ†é¡µã€‚
func (r *Repo) PageResources(ctx context.Context, p ResourcePageParam) (rows []Resource, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Resource{})
	if p.Name != "" {
		db = db.Where("name ILIKE ?", "%"+p.Name+"%")
	}
	if p.Code != "" {
		db = db.Where("code ILIKE ?", "%"+p.Code+"%")
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

// ListResourcesByClient æŒ‰å®¢æˆ·ç«¯åˆ—å‡ºå¯ç”¨èµ„æºã€‚
func (r *Repo) ListResourcesByClient(ctx context.Context, client string) ([]Resource, error) {
	var rows []Resource
	err := r.with(ctx).
		Where("status = ? AND module_id IN (SELECT id FROM sys_resource_module WHERE client = ? AND status = ?)",
			security.StatusEnabled, client, security.StatusEnabled).
		Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// GetResourcesByIDs æŒ‰ ID é›†åˆæŸ¥è¯¢èµ„æºã€‚
func (r *Repo) GetResourcesByIDs(ctx context.Context, ids []string) ([]Resource, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []Resource
	if err := r.with(ctx).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// PageButtons æŒ‰é’®èµ„æºåˆ†é¡µã€‚
func (r *Repo) PageButtons(ctx context.Context, p ButtonPageParam) (rows []Resource, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Resource{}).Where("resource_type = ?", ResourceTypeButton)
	if p.ResourceID != "" {
		db = db.Where("parent_id = ?", p.ResourceID)
	}
	if p.Code != "" {
		db = db.Where("code ILIKE ?", "%"+p.Code+"%")
	}
	if p.Name != "" {
		db = db.Where("name ILIKE ?", "%"+p.Name+"%")
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

// ListGrantResources åˆ—å‡ºæŒ‡å®šå®¢æˆ·ç«¯æ¨¡å—ä¸‹çš„å¯ç”¨èµ„æºï¼ˆæŽˆæƒæ ‘ç”¨ï¼‰ã€‚
func (r *Repo) ListGrantResources(ctx context.Context, client string) ([]Resource, error) {
	var rows []Resource
	err := r.with(ctx).
		Where("status = ? AND module_id IN (SELECT id FROM sys_resource_module WHERE client = ? AND status = ?)",
			security.StatusEnabled, client, security.StatusEnabled).
		Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// ListResources åˆ—å‡ºèµ„æºï¼ˆå¯é€‰æ¨¡å—è¿‡æ»¤ï¼‰ã€‚
func (r *Repo) ListResources(ctx context.Context, moduleID string) ([]Resource, error) {
	db := r.with(ctx).Model(&Resource{})
	if moduleID != "" {
		db = db.Where("module_id = ?", moduleID)
	}
	var rows []Resource
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// CreateModule åˆ›å»ºèµ„æºæ¨¡å—ã€‚
func (r *Repo) CreateModule(ctx context.Context, row *ResourceModule) error {
	return r.with(ctx).Create(row).Error
}

// UpdateModule æ›´æ–°èµ„æºæ¨¡å—ã€‚
func (r *Repo) UpdateModule(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&ResourceModule{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteModules æ‰¹é‡åˆ é™¤èµ„æºæ¨¡å—ã€‚
func (r *Repo) DeleteModules(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&ResourceModule{}).Error
}

// GetModuleByID æŒ‰ä¸»é”®æŸ¥è¯¢èµ„æºæ¨¡å—ã€‚
func (r *Repo) GetModuleByID(ctx context.Context, id string) (*ResourceModule, error) {
	var row ResourceModule
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageModules èµ„æºæ¨¡å—åˆ†é¡µã€‚
func (r *Repo) PageModules(ctx context.Context, p ModulePageParam) (rows []ResourceModule, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&ResourceModule{})
	if p.Name != "" {
		db = db.Where("name ILIKE ?", "%"+p.Name+"%")
	}
	if p.Client != "" {
		db = db.Where("client = ?", p.Client)
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

// ListEnabledModules åˆ—å‡ºå¯ç”¨æ¨¡å—ï¼ˆå¯é€‰å®¢æˆ·ç«¯è¿‡æ»¤ï¼‰ã€‚
func (r *Repo) ListEnabledModules(ctx context.Context, client string) ([]ResourceModule, error) {
	db := r.with(ctx).Model(&ResourceModule{}).Where("status = ?", security.StatusEnabled)
	if client != "" {
		db = db.Where("client = ?", client)
	}
	var rows []ResourceModule
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}
