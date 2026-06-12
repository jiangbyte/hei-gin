package resource

import (
	"context"

	"hei-gin/sdk/db"
	"hei-gin/sdk/utils"
	"gorm.io/gorm"
)

func ModulePageQuery(ctx context.Context, current, size int) ([]SysModule, int64) {
	var total int64
	db.DB.WithContext(ctx).Model(&SysModule{}).Count(&total)
	var rows []SysModule
	db.DB.WithContext(ctx).Order("created_at DESC").Limit(size).Offset((current-1)*size).Find(&rows)
	return rows, total
}

func FindModuleByID(ctx context.Context, id string) (*SysModule, error) {
	var e SysModule
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func CreateModule(ctx context.Context, entity *SysModule) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateModuleByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysModule{}).Where("id = ?", id).Updates(up).Error
}

func DeleteModules(ctx context.Context, ids []string) error {
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysModule{}).Error
}

func ResourcePageQuery(ctx context.Context, current, size int) ([]SysResource, int64) {
	var total int64
	db.DB.WithContext(ctx).Model(&SysResource{}).Count(&total)
	var rows []SysResource
	db.DB.WithContext(ctx).Order("sort_code ASC").Limit(size).Offset((current-1)*size).Find(&rows)
	return rows, total
}

func ListResources(ctx context.Context, category string) ([]SysResource, error) {
	q := db.DB.WithContext(ctx).Model(&SysResource{}).Order("sort_code ASC")
	if category != "" {
		q = q.Where("category = ?", category)
	}
	var all []SysResource
	if err := q.Find(&all).Error; err != nil {
		return nil, err
	}
	return all, nil
}

func ListAllResources(ctx context.Context) ([]SysResource, error) {
	var all []SysResource
	if err := db.DB.WithContext(ctx).Model(&SysResource{}).Order("sort_code ASC").Find(&all).Error; err != nil {
		return nil, err
	}
	return all, nil
}

func FindResourceByID(ctx context.Context, id string) (*SysResource, error) {
	var e SysResource
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func CreateResource(ctx context.Context, entity *SysResource) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateResourceByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysResource{}).Where("id = ?", id).Updates(up).Error
}

func DeleteResourcesCascade(ctx context.Context, ids []string) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Table("rel_role_resource").Where("resource_id IN ?", ids).Delete(nil).Error; err != nil {
		tx.Rollback()
		return err
	}
	subQuery := tx.Table("rel_role_resource").Select("role_id").Where("resource_id IN ?", ids)
	if err := tx.Table("rel_role_permission").Where("role_id IN (?)", subQuery).Delete(nil).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id IN ?", ids).Delete(&SysResource{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func CollectResourceDescendants(ctx context.Context, ids []string) []string {
	m := make(map[string]bool)
	for _, id := range ids {
		m[id] = true
	}

	var all []SysResource
	if err := db.DB.WithContext(ctx).Find(&all).Error; err != nil {
		return ids
	}
	cm := make(map[string][]string)
	for _, r := range all {
		if r.ParentID != nil && *r.ParentID != "" {
			cm[*r.ParentID] = append(cm[*r.ParentID], r.ID)
		}
	}

	q := make([]string, len(ids))
	copy(q, ids)
	for len(q) > 0 {
		pid := q[len(q)-1]
		q = q[:len(q)-1]
		for _, cid := range cm[pid] {
			if !m[cid] {
				m[cid] = true
				q = append(q, cid)
			}
		}
	}
	r := make([]string, 0, len(m))
	for id := range m {
		r = append(r, id)
	}
	return r
}

func SyncPermissions(ctx context.Context, resourceID string, oldCode, newCode string) error {
	if oldCode == newCode {
		return nil
	}

	tx := db.DB.WithContext(ctx).Begin()

	var roleResources []relRoleResource
	if err := tx.Table("rel_role_resource").Where("resource_id = ?", resourceID).Find(&roleResources).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(roleResources) == 0 {
		tx.Rollback()
		return nil
	}

	roleIDs := make([]string, len(roleResources))
	for i, rr := range roleResources {
		roleIDs[i] = rr.RoleID
	}

	if oldCode != "" {
		if err := tx.Table("rel_role_permission").Where("role_id IN ? AND permission_code = ?", roleIDs, oldCode).Delete(nil).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if newCode != "" {
		var existing []struct{ RoleID string }
		if err := tx.Table("rel_role_permission").Where("role_id IN ? AND permission_code = ?", roleIDs, newCode).Select("role_id").Find(&existing).Error; err != nil {
			tx.Rollback()
			return err
		}
		existingMap := make(map[string]bool)
		for _, e := range existing {
			existingMap[e.RoleID] = true
		}
		permBatch := make([]struct {
			ID             string
			RoleID         string
			PermissionCode string
		}, 0)
		for _, rid := range roleIDs {
			if existingMap[rid] {
				continue
			}
			permBatch = append(permBatch, struct {
				ID             string
				RoleID         string
				PermissionCode string
			}{utils.GenerateID(), rid, newCode})
		}
		if len(permBatch) > 0 {
			if err := tx.Table("rel_role_permission").CreateInBatches(permBatch, 100).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

var _ *gorm.DB
