package user

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type cacheOrg struct {
	ID       string
	Name     string
	ParentID *string
}

type cacheGroup struct {
	ID       string
	Name     string
	ParentID *string
}

type positionRow struct {
	ID   string
	Name string
}

type roleCodeRow struct {
	Code string
}

type rawResource struct {
	ID            string
	ParentID      *string
	Code          string
	Name          string
	Category      string
	Type          string
	RoutePath     *string
	ComponentPath *string
	RedirectPath  *string
	Icon          *string
	Color         *string
	IsVisible     string
	IsCache       string
	IsAffix       string
	IsBreadcrumb  string
	ExternalURL   *string
	Description   *string
	SortCode      int
	Status        string
}

func (r *repository) Page(ctx context.Context, p *UserPageParam) ([]SysUser, int64) {
	q := r.db.WithContext(ctx).Model(&SysUser{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?", like, like, like, like)
	}
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}
	var total int64
	q.Count(&total)
	var rows []SysUser
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) FindByID(ctx context.Context, id string) (*SysUser, error) {
	var e SysUser
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) CountByUsername(ctx context.Context, username, excludeID string) int64 {
	var count int64
	q := r.db.WithContext(ctx).Model(&SysUser{}).Where("username = ?", username)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	q.Count(&count)
	return count
}

func (r *repository) Create(ctx context.Context, entity *SysUser) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&SysUser{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) DeleteUsersCascade(ctx context.Context, ids []string) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Where("user_id IN ?", ids).Delete(&RelUserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("user_id IN ?", ids).Delete(&RelUserPermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("user_id IN ?", ids).Delete(&SysQuickAction{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id IN ?", ids).Delete(&SysUser{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) ReplaceUserRoles(ctx context.Context, userID string, roles []RelUserRole) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Where("user_id = ?", userID).Delete(&RelUserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(roles) > 0 {
		if err := tx.Create(&roles).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *repository) ReplaceUserPermissions(ctx context.Context, userID string, perms []RelUserPermission) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Where("user_id = ?", userID).Delete(&RelUserPermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(perms) > 0 {
		if err := tx.Create(&perms).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *repository) UpdateStatusByIDs(ctx context.Context, ids []string, status string) error {
	return r.db.WithContext(ctx).Model(&SysUser{}).Where("id IN ?", ids).Updates(map[string]interface{}{"status": status}).Error
}

func (r *repository) FindRoleRelsByUserIDs(ctx context.Context, userIDs []string) []RelUserRole {
	var rows []RelUserRole
	r.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&rows)
	return rows
}

func (r *repository) FindRoleRelsByUserID(ctx context.Context, userID string) []RelUserRole {
	var rows []RelUserRole
	r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows)
	return rows
}

func (r *repository) FindPermissionRelsByUserID(ctx context.Context, userID string) []RelUserPermission {
	var rows []RelUserPermission
	r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows)
	return rows
}

func (r *repository) FindDistinctUserPermissionCodes(ctx context.Context, userID string) []RelUserPermission {
	var rows []RelUserPermission
	r.db.WithContext(ctx).Where("user_id = ?", userID).Select("DISTINCT permission_code").Find(&rows)
	return rows
}

func (r *repository) FindDistinctRolePermissionCodes(ctx context.Context, roleIDs []string) []RelRolePermission {
	var rows []RelRolePermission
	r.db.WithContext(ctx).Where("role_id IN ?", roleIDs).Select("DISTINCT permission_code").Find(&rows)
	return rows
}

func (r *repository) FindOrgCacheRows(ctx context.Context) []cacheOrg {
	var rows []cacheOrg
	r.db.WithContext(ctx).Table("sys_org").Select("id,name,parent_id").Find(&rows)
	return rows
}

func (r *repository) FindGroupCacheRows(ctx context.Context) []cacheGroup {
	var rows []cacheGroup
	r.db.WithContext(ctx).Table("sys_group").Select("id,name,parent_id").Find(&rows)
	return rows
}

func (r *repository) FindPositionRows(ctx context.Context, ids []string) []positionRow {
	var rows []positionRow
	r.db.WithContext(ctx).Table("sys_position").Where("id IN ?", ids).Find(&rows)
	return rows
}

func (r *repository) UpdateAvatar(ctx context.Context, entity *SysUser, avatar string) error {
	return r.db.WithContext(ctx).Model(entity).Update("avatar", avatar).Error
}

func (r *repository) UpdatePassword(ctx context.Context, userID, password string) error {
	return r.db.WithContext(ctx).Model(&SysUser{}).Where("id = ?", userID).Update("password", password).Error
}

func (r *repository) FindRoleCodesByIDs(ctx context.Context, ids []string) []roleCodeRow {
	var rows []roleCodeRow
	r.db.WithContext(ctx).Table("sys_role").Where("id IN ?", ids).Find(&rows)
	return rows
}

func (r *repository) FindEnabledResources(ctx context.Context) []rawResource {
	var resources []rawResource
	r.db.WithContext(ctx).Table("sys_resource").Where("status = ?", "ENABLED").Order("sort_code ASC").Find(&resources)
	return resources
}

func (r *repository) FindRoleResourcesByRoleIDs(ctx context.Context, roleIDs []string) []RelRoleResource {
	var rows []RelRoleResource
	r.db.WithContext(ctx).Where("role_id IN ?", roleIDs).Find(&rows)
	return rows
}

func (r *repository) FindEnabledResourcesByIDs(ctx context.Context, ids []string) []rawResource {
	var resources []rawResource
	r.db.WithContext(ctx).Table("sys_resource").Where("id IN ? AND status = ?", ids, "ENABLED").Order("sort_code ASC").Find(&resources)
	return resources
}
