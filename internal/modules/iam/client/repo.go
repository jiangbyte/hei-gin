// internal/modules/iam/client/repo.go 持久化仓储。
//
// Author: Charlie

package client

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

// Repo 客户端资源持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// CreateModule 创建客户端模块。
func (r *Repo) CreateModule(ctx context.Context, row *ClientModule) error {
	return r.with(ctx).Create(row).Error
}

// UpdateModule 更新客户端模块。
func (r *Repo) UpdateModule(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&ClientModule{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteModules 批量删除客户端模块。
func (r *Repo) DeleteModules(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&ClientModule{}).Error
}

// GetModuleByID 按主键查询客户端模块。
func (r *Repo) GetModuleByID(ctx context.Context, id string) (*ClientModule, error) {
	var row ClientModule
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageModules 客户端模块分页。
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

// ListEnabledModules 列出启用模块（可选账号类型过滤）。
func (r *Repo) ListEnabledModules(ctx context.Context, accountType string) ([]ClientModule, error) {
	db := r.with(ctx).Model(&ClientModule{}).Where("status = ?", security.StatusEnabled)
	if accountType != "" {
		db = db.Where("account_type = ?", accountType)
	}
	var rows []ClientModule
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// CreateResource 创建客户端资源。
func (r *Repo) CreateResource(ctx context.Context, row *ClientResource) error {
	return r.with(ctx).Create(row).Error
}

// UpdateResource 更新客户端资源。
func (r *Repo) UpdateResource(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&ClientResource{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteResources 批量删除客户端资源。
func (r *Repo) DeleteResources(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&ClientResource{}).Error
}

// GetResourceByID 按主键查询客户端资源。
func (r *Repo) GetResourceByID(ctx context.Context, id string) (*ClientResource, error) {
	var row ClientResource
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageResources 客户端资源分页。
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

// ListGrantResources 列出指定账号类型模块下的启用客户端资源（授权树用）。
func (r *Repo) ListGrantResources(ctx context.Context, accountType string) ([]ClientResource, error) {
	var rows []ClientResource
	err := r.with(ctx).
		Where("status = ? AND module_id IN (SELECT id FROM sys_client_module WHERE account_type = ? AND status = ?)",
			security.StatusEnabled, accountType, security.StatusEnabled).
		Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}

// ListResources 列出客户端资源（可选模块过滤）。
func (r *Repo) ListResources(ctx context.Context, moduleID string) ([]ClientResource, error) {
	db := r.with(ctx).Model(&ClientResource{})
	if moduleID != "" {
		db = db.Where("module_id = ?", moduleID)
	}
	var rows []ClientResource
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}
