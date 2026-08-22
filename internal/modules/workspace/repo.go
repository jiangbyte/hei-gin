// internal/modules/workspace/repo.go 持久化仓储。
//
// Author: Charlie

package workspace

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

const activityLimit = 10

// Repo 工作台持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// ListShortcutsByAccount 按账号查询快捷应用（按 sort、id 升序）。
func (r *Repo) ListShortcutsByAccount(ctx context.Context, accountID string) ([]WorkspaceShortcut, error) {
	var rows []WorkspaceShortcut
	err := r.with(ctx).Where("account_id = ?", accountID).
		Order("sort asc, id asc").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []WorkspaceShortcut{}
	}
	return rows, nil
}

// DeleteShortcutsByAccount 删除账号下全部快捷应用。
func (r *Repo) DeleteShortcutsByAccount(ctx context.Context, accountID string) error {
	return r.with(ctx).Where("account_id = ?", accountID).Delete(&WorkspaceShortcut{}).Error
}

// CreateShortcuts 批量创建快捷应用。
func (r *Repo) CreateShortcuts(ctx context.Context, rows []WorkspaceShortcut) error {
	if len(rows) == 0 {
		return nil
	}
	return r.with(ctx).Create(&rows).Error
}

// ListMenusByIDs 按 ID 查询启用 MENU 资源。
func (r *Repo) ListMenusByIDs(ctx context.Context, ids []string) ([]MenuResource, error) {
	if len(ids) == 0 {
		return []MenuResource{}, nil
	}
	var rows []MenuResource
	err := r.with(ctx).
		Select("id, code, name, resource_type, path, icon, status").
		Where("id IN ? AND status = ? AND resource_type = ?", ids, security.StatusEnabled, "MENU").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []MenuResource{}
	}
	return rows, nil
}

// GetMenuByID 按主键查询菜单资源。
func (r *Repo) GetMenuByID(ctx context.Context, id string) (*MenuResource, error) {
	var row MenuResource
	if err := r.with(ctx).
		Select("id, code, name, resource_type, path, icon, status").
		First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ListRecentOperations 查询本人近期操作日志（排除 login）。
func (r *Repo) ListRecentOperations(ctx context.Context, accountID string) ([]AuditActivity, error) {
	var rows []AuditActivity
	err := r.with(ctx).
		Where("account_id = ? AND (action <> ? OR action IS NULL)", accountID, "login").
		Order("created_at desc").
		Limit(activityLimit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []AuditActivity{}
	}
	return rows, nil
}

// ListRecentLogins 查询本人近期登录日志。
func (r *Repo) ListRecentLogins(ctx context.Context, accountID string) ([]AuditActivity, error) {
	var rows []AuditActivity
	err := r.with(ctx).
		Where("account_id = ? AND action = ?", accountID, "login").
		Order("created_at desc").
		Limit(activityLimit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []AuditActivity{}
	}
	return rows, nil
}

// ListRoleIDs 查账号已启用角色 ID。
func (r *Repo) ListRoleIDs(ctx context.Context, accountID string) ([]string, error) {
	var ids []string
	err := r.with(ctx).Table("sys_iam_relation").
		Select("target_id").
		Where("subject_type = ? AND subject_id = ? AND relation_type = ? AND status = ?",
			"ACCOUNT", accountID, "ACCOUNT_ROLE", security.StatusEnabled).
		Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, nil
}

// ListGroupIDs 查账号已加入的用户组 ID。
func (r *Repo) ListGroupIDs(ctx context.Context, accountID string) ([]string, error) {
	var ids []string
	err := r.with(ctx).Table("sys_iam_relation").
		Select("target_id").
		Where("subject_type = ? AND subject_id = ? AND relation_type = ? AND status = ?",
			"ACCOUNT", accountID, "ACCOUNT_GROUP", security.StatusEnabled).
		Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, nil
}

// ListGrantedResourceIDs 列出账号/角色/用户组主体已授予的资源 ID。
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
	fullArgs = append(fullArgs, []string{"ACCOUNT_RESOURCE", "GROUP_RESOURCE", "ROLE_RESOURCE"}, "RESOURCE", security.StatusEnabled)
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
