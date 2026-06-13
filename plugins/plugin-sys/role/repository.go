package role

import (
	"context"

	"gorm.io/gorm"
	userModel "hei-gin/plugins/plugin-sys/user"
)

type repository struct {
	db *gorm.DB
}

type resourceExtraRow struct {
	ID    string
	Extra *string
}

func (r *repository) Page(ctx context.Context, p *RolePageParam) ([]SysRole, int64) {
	q := r.db.WithContext(ctx).Model(&SysRole{})
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
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) FindByID(ctx context.Context, id string) (*SysRole, error) {
	var e SysRole
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Create(ctx context.Context, entity *SysRole) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&SysRole{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) CountUserRoleRefs(ctx context.Context, ids []string) int64 {
	var count int64
	r.db.WithContext(ctx).Model(&userModel.RelUserRole{}).Where("role_id IN ?", ids).Count(&count)
	return count
}

func (r *repository) DeleteCascade(ctx context.Context, ids []string) error {
	tx := r.db.WithContext(ctx).Begin()
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

func (r *repository) FindRolePermissions(ctx context.Context, roleID string) []userModel.RelRolePermission {
	var rows []userModel.RelRolePermission
	r.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&rows)
	return rows
}

func (r *repository) FindRolePermissionCodes(ctx context.Context, roleID string) []userModel.RelRolePermission {
	var rows []userModel.RelRolePermission
	r.db.WithContext(ctx).Where("role_id = ?", roleID).Select("permission_code").Find(&rows)
	return rows
}

func (r *repository) FindRoleResources(ctx context.Context, roleID string) []userModel.RelRoleResource {
	var rows []userModel.RelRoleResource
	r.db.WithContext(ctx).Where("role_id = ?", roleID).Select("resource_id").Find(&rows)
	return rows
}

func (r *repository) FindUserIDsByRoleID(ctx context.Context, roleID string) []string {
	var rows []userModel.RelUserRole
	r.db.WithContext(ctx).Where("role_id = ?", roleID).Select("user_id").Find(&rows)
	userIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.UserID == "" {
			continue
		}
		userIDs = append(userIDs, row.UserID)
	}
	return userIDs
}

func (r *repository) ReplaceRolePermissions(ctx context.Context, roleID string, perms []userModel.RelRolePermission) error {
	tx := r.db.WithContext(ctx).Begin()
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

func (r *repository) ReplaceRoleResourcesAndAppendPerms(ctx context.Context, roleID string, resources []userModel.RelRoleResource, resourceIDs []string, perms []userModel.RelRolePermission) error {
	tx := r.db.WithContext(ctx).Begin()
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

func (r *repository) FindResourceExtras(ctx context.Context, ids []string) []resourceExtraRow {
	var rows []resourceExtraRow
	r.db.WithContext(ctx).Table("sys_resource").Where("id IN ? AND extra IS NOT NULL AND extra != ''", ids).Find(&rows)
	return rows
}

func FindExistingPermissionCodesWithTx(tx *gorm.DB, roleID string) []userModel.RelRolePermission {
	var rows []userModel.RelRolePermission
	tx.Where("role_id = ?", roleID).Select("permission_code").Find(&rows)
	return rows
}
