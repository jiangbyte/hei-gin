// internal/modules/iam/resource/repo.go 持久化仓储。
//
// Author: Charlie

package resource

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

// Repo 资源持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// CreateResource 创建资源。
func (r *Repo) CreateResource(ctx context.Context, row *Resource) error {
	return r.with(ctx).Create(row).Error
}

// UpdateResource 更新资源。
func (r *Repo) UpdateResource(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Resource{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteResources 批量删除资源。
func (r *Repo) DeleteResources(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Resource{}).Error
}

// GetResourceByID 按主键查询资源。
func (r *Repo) GetResourceByID(ctx context.Context, id string) (*Resource, error) {
	var row Resource
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageResources 资源分页。
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

// ListResourcesByClient 按客户端列出启用资源。
func (r *Repo) ListResourcesByClient(ctx context.Context, client string) ([]Resource, error) {
	var rows []Resource
	err := r.with(ctx).
		Where("status = ? AND module_id IN (SELECT id FROM sys_resource_module WHERE client = ? AND status = ?)",
			security.StatusEnabled, client, security.StatusEnabled).
		Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// GetResourcesByIDs 按 ID 集合查询资源。
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

// PageButtons 按钮资源分页。
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

// ListGrantResources 列出指定客户端模块下的启用资源（授权树用）。
func (r *Repo) ListGrantResources(ctx context.Context, client string) ([]Resource, error) {
	var rows []Resource
	err := r.with(ctx).
		Where("status = ? AND module_id IN (SELECT id FROM sys_resource_module WHERE client = ? AND status = ?)",
			security.StatusEnabled, client, security.StatusEnabled).
		Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// ListResources 列出资源（可选模块过滤）。
func (r *Repo) ListResources(ctx context.Context, moduleID string) ([]Resource, error) {
	db := r.with(ctx).Model(&Resource{})
	if moduleID != "" {
		db = db.Where("module_id = ?", moduleID)
	}
	var rows []Resource
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// CreateModule 创建资源模块。
func (r *Repo) CreateModule(ctx context.Context, row *ResourceModule) error {
	return r.with(ctx).Create(row).Error
}

// UpdateModule 更新资源模块。
func (r *Repo) UpdateModule(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&ResourceModule{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteModules 批量删除资源模块。
func (r *Repo) DeleteModules(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&ResourceModule{}).Error
}

// GetModuleByID 按主键查询资源模块。
func (r *Repo) GetModuleByID(ctx context.Context, id string) (*ResourceModule, error) {
	var row ResourceModule
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageModules 资源模块分页。
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

// ListEnabledModules 列出启用模块（可选客户端过滤）。
func (r *Repo) ListEnabledModules(ctx context.Context, client string) ([]ResourceModule, error) {
	db := r.with(ctx).Model(&ResourceModule{}).Where("status = ?", security.StatusEnabled)
	if client != "" {
		db = db.Where("client = ?", client)
	}
	var rows []ResourceModule
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}
