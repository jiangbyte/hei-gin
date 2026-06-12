package user

import (
	"context"

	"hei-gin/sdk/db"
	"gorm.io/gorm"
)

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

func Page(ctx context.Context, p *UserPageParam) ([]SysUser, int64) {
	q := db.DB.WithContext(ctx).Model(&SysUser{})
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
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current-1)*p.Size).Find(&rows)
	return rows, total
}

func FindByID(ctx context.Context, id string) (*SysUser, error) {
	var e SysUser
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func CountByUsername(ctx context.Context, username, excludeID string) int64 {
	var count int64
	q := db.DB.WithContext(ctx).Model(&SysUser{}).Where("username = ?", username)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	q.Count(&count)
	return count
}

func Create(ctx context.Context, entity *SysUser) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysUser{}).Where("id = ?", id).Updates(up).Error
}

func DeleteUsersCascade(ctx context.Context, ids []string) error {
	tx := db.DB.WithContext(ctx).Begin()
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

func ReplaceUserRoles(ctx context.Context, userID string, roles []RelUserRole) error {
	tx := db.DB.WithContext(ctx).Begin()
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

func ReplaceUserPermissions(ctx context.Context, userID string, perms []RelUserPermission) error {
	tx := db.DB.WithContext(ctx).Begin()
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

func UpdateStatusByIDs(ctx context.Context, ids []string, status string) error {
	return db.DB.WithContext(ctx).Model(&SysUser{}).Where("id IN ?", ids).Updates(map[string]interface{}{"status": status}).Error
}

func FindRoleRelsByUserIDs(ctx context.Context, userIDs []string) []RelUserRole {
	var rows []RelUserRole
	db.DB.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&rows)
	return rows
}

func FindRoleRelsByUserID(ctx context.Context, userID string) []RelUserRole {
	var rows []RelUserRole
	db.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&rows)
	return rows
}

func FindPermissionRelsByUserID(ctx context.Context, userID string) []RelUserPermission {
	var rows []RelUserPermission
	db.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&rows)
	return rows
}

func FindDistinctUserPermissionCodes(ctx context.Context, userID string) []RelUserPermission {
	var rows []RelUserPermission
	db.DB.WithContext(ctx).Where("user_id = ?", userID).Select("DISTINCT permission_code").Find(&rows)
	return rows
}

func FindDistinctRolePermissionCodes(ctx context.Context, roleIDs []string) []RelRolePermission {
	var rows []RelRolePermission
	db.DB.WithContext(ctx).Where("role_id IN ?", roleIDs).Select("DISTINCT permission_code").Find(&rows)
	return rows
}

func FindOrgCacheRows(ctx context.Context) []cacheOrg {
	var rows []cacheOrg
	db.DB.WithContext(ctx).Table("sys_org").Select("id,name,parent_id").Find(&rows)
	return rows
}

func FindGroupCacheRows(ctx context.Context) []cacheGroup {
	var rows []cacheGroup
	db.DB.WithContext(ctx).Table("sys_group").Select("id,name,parent_id").Find(&rows)
	return rows
}

func FindPositionRows(ctx context.Context, ids []string) []positionRow {
	var rows []positionRow
	db.DB.WithContext(ctx).Table("sys_position").Where("id IN ?", ids).Find(&rows)
	return rows
}

func UpdateAvatar(ctx context.Context, entity *SysUser, avatar string) error {
	return db.DB.WithContext(ctx).Model(entity).Update("avatar", avatar).Error
}

func UpdatePassword(ctx context.Context, userID, password string) error {
	return db.DB.WithContext(ctx).Model(&SysUser{}).Where("id = ?", userID).Update("password", password).Error
}

func FindRoleCodesByIDs(ctx context.Context, ids []string) []roleCodeRow {
	var rows []roleCodeRow
	db.DB.WithContext(ctx).Table("sys_role").Where("id IN ?", ids).Find(&rows)
	return rows
}

func FindEnabledResources(ctx context.Context) []rawResource {
	var resources []rawResource
	db.DB.WithContext(ctx).Table("sys_resource").Where("status = ?", "ENABLED").Order("sort_code ASC").Find(&resources)
	return resources
}

func FindRoleResourcesByRoleIDs(ctx context.Context, roleIDs []string) []RelRoleResource {
	var rows []RelRoleResource
	db.DB.WithContext(ctx).Where("role_id IN ?", roleIDs).Find(&rows)
	return rows
}

func FindEnabledResourcesByIDs(ctx context.Context, ids []string) []rawResource {
	var resources []rawResource
	db.DB.WithContext(ctx).Table("sys_resource").Where("id IN ? AND status = ?", ids, "ENABLED").Order("sort_code ASC").Find(&resources)
	return resources
}

var _ *gorm.DB
