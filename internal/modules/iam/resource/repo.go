// internal/modules/iam/resource/repo.go 持久化仓储。
//
// Author: Charlie

package resource

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/db/dialect"
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
		db = db.Where(dialect.ILike(db, "name"), dialect.Contains(p.Name))
	}
	if p.Code != "" {
		db = db.Where(dialect.ILike(db, "code"), dialect.Contains(p.Code))
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

// ButtonPermissions 批量加载按钮权限绑定（RESOURCE_PERMISSION，subject=按钮 id）。
func (r *Repo) ButtonPermissions(ctx context.Context, buttonIDs []string) map[string]buttonPerm {
	out := make(map[string]buttonPerm, len(buttonIDs))
	if len(buttonIDs) == 0 {
		return out
	}
	var rows []struct {
		SubjectID          string         `gorm:"column:subject_id"`
		TargetKey          string         `gorm:"column:target_key"`
		DataScope          string         `gorm:"column:data_scope"`
		CustomScopeDeptIDs datatypes.JSON `gorm:"column:custom_scope_dept_ids"`
		Description        *string        `gorm:"column:description"`
	}
	if err := r.with(ctx).Table("sys_iam_relation").
		Select("subject_id", "target_key", "data_scope", "custom_scope_dept_ids", "description").
		Where("subject_type = ? AND relation_type = ? AND subject_id IN ? AND status = ?",
			"RESOURCE", "RESOURCE_PERMISSION", buttonIDs, "ENABLED").
		Find(&rows).Error; err != nil {
		return out
	}
	for _, row := range rows {
		perm := buttonPerm{PermissionKey: row.TargetKey, DataScope: row.DataScope, Description: row.Description}
		if len(row.CustomScopeDeptIDs) > 0 {
			_ = json.Unmarshal(row.CustomScopeDeptIDs, &perm.CustomScopeDeptIDs)
		}
		out[row.SubjectID] = perm
	}
	return out
}

// GrantPermissions 批量加载资源权限绑定选项（RESOURCE_PERMISSION，subject=资源 id，可按账号类型过滤；对齐 hei-boot listGrantModules 第 3 步）。
func (r *Repo) GrantPermissions(ctx context.Context, resourceIDs []string, accountType string) (map[string][]PermissionOption, error) {
	out := make(map[string][]PermissionOption)
	if len(resourceIDs) == 0 {
		return out, nil
	}
	var rows []struct {
		ID          string  `gorm:"column:id"`
		SubjectID   string  `gorm:"column:subject_id"`
		TargetKey   string  `gorm:"column:target_key"`
		DataScope   string  `gorm:"column:data_scope"`
		Description *string `gorm:"column:description"`
	}
	q := r.with(ctx).Table("sys_iam_relation").
		Select("id", "subject_id", "target_key", "data_scope", "description").
		Where("subject_type = ? AND relation_type = ? AND subject_id IN ? AND status = ?",
			"RESOURCE", "RESOURCE_PERMISSION", resourceIDs, "ENABLED")
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		title := row.TargetKey
		if row.Description != nil && *row.Description != "" {
			title = *row.Description
		}
		out[row.SubjectID] = append(out[row.SubjectID], PermissionOption{
			ID: row.ID, PermissionKey: row.TargetKey, Title: title, DataScope: row.DataScope,
		})
	}
	return out, nil
}

// ModulesByIDs 批量查资源模块。
func (r *Repo) ModulesByIDs(ctx context.Context, ids []string) ([]ResourceModule, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []ResourceModule
	err := r.with(ctx).Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

// PageButtons 按钮资源分页。
func (r *Repo) PageButtons(ctx context.Context, p ButtonPageParam) (rows []Resource, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Resource{}).Where("resource_type = ?", ResourceTypeButton)
	if p.ParentID != "" {
		db = db.Where("parent_id = ?", p.ParentID)
	}
	if p.Code != "" {
		db = db.Where(dialect.ILike(db, "code"), dialect.Contains(p.Code))
	}
	if p.Name != "" {
		db = db.Where(dialect.ILike(db, "name"), dialect.Contains(p.Name))
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
		db = db.Where(dialect.ILike(db, "name"), dialect.Contains(p.Name))
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

// ListGrantedResourceIDs 列出账号/角色/用户组主体已授予的资源 ID（ACCOUNT_RESOURCE/GROUP_RESOURCE/ROLE_RESOURCE，按账号类型过滤）。
func (r *Repo) ListGrantedResourceIDs(ctx context.Context, accountID string, groupIDs, roleIDs []string, accountType string) ([]string, error) {
	cond := "(subject_type = ? AND subject_id = ?"
	args := []any{"ACCOUNT", accountID}
	if len(groupIDs) > 0 {
		cond += " OR (subject_type = ? AND subject_id IN ?)"
		args = append(args, "GROUP", groupIDs)
	}
	if len(roleIDs) > 0 {
		cond += " OR (subject_type = ? AND subject_id IN ?)"
		args = append(args, "ROLE", roleIDs)
	}
	cond += ")"
	fullArgs := make([]any, 0, len(args)+3)
	fullArgs = append(fullArgs, []string{
		"SUBJECT_RESOURCE_GRANT", "ACCOUNT_RESOURCE", "GROUP_RESOURCE", "ROLE_RESOURCE",
	}, "RESOURCE", "ENABLED")
	fullArgs = append(fullArgs, args...)
	q := r.with(ctx).Table("sys_iam_relation").
		Select("DISTINCT target_id").
		Where("relation_type IN ? AND target_type = ? AND status = ? AND "+cond, fullArgs...)
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	var ids []string
	if err := q.Scan(&ids).Error; err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, nil
}

// ListResourcesByIDsWithParents 按授权 ID 补齐祖先链后返回（保持 sort 排序；对齐 hei-boot listResourcesByIdsWithParents）。
func (r *Repo) ListResourcesByIDsWithParents(ctx context.Context, resourceIDs []string, client string) ([]Resource, error) {
	all, err := r.ListResourcesByClient(ctx, client)
	if err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return []Resource{}, nil
	}
	byID := make(map[string]*Resource, len(all))
	for i := range all {
		byID[all[i].ID] = &all[i]
	}
	selected := map[string]struct{}{}
	for _, id := range resourceIDs {
		cur := byID[id]
		for cur != nil {
			if _, ok := selected[cur.ID]; ok {
				break
			}
			selected[cur.ID] = struct{}{}
			if cur.ParentID == nil || *cur.ParentID == "" {
				break
			}
			cur = byID[*cur.ParentID]
		}
	}
	out := make([]Resource, 0, len(selected))
	for i := range all {
		if _, ok := selected[all[i].ID]; ok {
			out = append(out, all[i])
		}
	}
	return out, nil
}
