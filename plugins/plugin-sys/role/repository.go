package role

import (
	"context"

	userModel "hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/db"
	"gorm.io/gorm"
)

type resourceExtraRow struct {
	ID    string
	Extra *string
}

func Page(ctx context.Context, p *RolePageParam) ([]SysRole, int64) {
	q := db.DB.WithContext(ctx).Model(&SysRole{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("name LIKE ? OR code LIKE ?", like, like)
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	var total int64
	q.Count(&total)
	var rows []SysRole
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current-1)*p.Size).Find(&rows)
	return rows, total
}

func FindByID(ctx context.Context, id string) (*SysRole, error) {
	var e SysRole
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func Create(ctx context.Context, entity *SysRole) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysRole{}).Where("id = ?", id).Updates(up).Error
}

func CountUserRoleRefs(ctx context.Context, ids []string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&userModel.RelUserRole{}).Where("role_id IN ?", ids).Count(&count)
	return count
}

func DeleteCascade(ctx context.Context, ids []string) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("role_id IN ?", ids).Delete(&userModel.RelRolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("role_id IN ?", ids).Delete(&userModel.RelRoleResource{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("role_id IN ?", ids).Delete(&userModel.RelUserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id IN ?", ids).Delete(&SysRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func FindRolePermissions(ctx context.Context, roleID string) []userModel.RelRolePermission {
	var rows []userModel.RelRolePermission
	db.DB.WithContext(ctx).Where("role_id = ?", roleID).Find(&rows)
	return rows
}

func FindRolePermissionCodes(ctx context.Context, roleID string) []userModel.RelRolePermission {
	var rows []userModel.RelRolePermission
	db.DB.WithContext(ctx).Where("role_id = ?", roleID).Select("permission_code").Find(&rows)
	return rows
}

func FindRoleResources(ctx context.Context, roleID string) []userModel.RelRoleResource {
	var rows []userModel.RelRoleResource
	db.DB.WithContext(ctx).Where("role_id = ?", roleID).Select("resource_id").Find(&rows)
	return rows
}

func ReplaceRolePermissions(ctx context.Context, roleID string, perms []userModel.RelRolePermission) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("role_id = ?", roleID).Delete(&userModel.RelRolePermission{}).Error; err != nil {
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

func ReplaceRoleResourcesAndAppendPerms(ctx context.Context, roleID string, resources []userModel.RelRoleResource, resourceIDs []string, perms []userModel.RelRolePermission) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("role_id = ?", roleID).Delete(&userModel.RelRoleResource{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(resources) > 0 {
		if err := tx.Create(&resources).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(resourceIDs) > 0 && len(perms) > 0 {
		if err := tx.Create(&perms).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func FindResourceExtras(ctx context.Context, ids []string) []resourceExtraRow {
	var rows []resourceExtraRow
	db.DB.WithContext(ctx).Table("sys_resource").Where("id IN ? AND extra IS NOT NULL AND extra != ''", ids).Find(&rows)
	return rows
}

func FindExistingPermissionCodesWithTx(tx *gorm.DB, roleID string) []userModel.RelRolePermission {
	var rows []userModel.RelRolePermission
	tx.Where("role_id = ?", roleID).Select("permission_code").Find(&rows)
	return rows
}
