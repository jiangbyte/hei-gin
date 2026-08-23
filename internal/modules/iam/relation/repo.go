// internal/modules/iam/relation/repo.go 持久化仓储。
//
// Author: Charlie

package relation

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

// Repo 关系持久化（sys_iam_relation）。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// ListRelations 列出主体指定关系类型的关系行（accountType 为空不过滤）。
func (r *Repo) ListRelations(ctx context.Context, subjectType, subjectID, relationType, accountType string) ([]Relation, error) {
	db := r.with(ctx).Where("subject_type = ? AND subject_id = ? AND relation_type = ? AND status = ?",
		subjectType, subjectID, relationType, security.StatusEnabled)
	if accountType != "" {
		db = db.Where("account_type = ?", accountType)
	}
	var rows []Relation
	if err := db.Order("sort asc, id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// deleteSubjectRelations 删除主体指定关系类型的关系（accountType 为空删全部类型，供事务内调用）。
func (r *Repo) deleteSubjectRelations(db *gorm.DB, subjectType, subjectID, relationType, accountType string) error {
	q := db.Where("subject_type = ? AND subject_id = ? AND relation_type = ?", subjectType, subjectID, relationType)
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	return q.Delete(&Relation{}).Error
}

// deleteSubjectRelationsByKeys 按权限键删除主体关系（对齐 hei-boot bindPermission 单键覆盖）。
func (r *Repo) deleteSubjectRelationsByKeys(db *gorm.DB, subjectType, subjectID, relationType, accountType string, targetKeys []string) error {
	if len(targetKeys) == 0 {
		return nil
	}
	q := db.Where("subject_type = ? AND subject_id = ? AND relation_type = ? AND target_key IN ?",
		subjectType, subjectID, relationType, targetKeys)
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	return q.Delete(&Relation{}).Error
}

// deleteByTarget 按目标删除指定关系类型（组/角色成员替换用）。
func (r *Repo) deleteByTarget(db *gorm.DB, relationType, targetType, targetID string) error {
	return db.Where("relation_type = ? AND target_type = ? AND target_id = ?",
		relationType, targetType, targetID).Delete(&Relation{}).Error
}

// ListSubjectIDsByTarget 按目标反查主体 ID（如用户组成员账号）。
func (r *Repo) ListSubjectIDsByTarget(ctx context.Context, relationType, targetType, targetID string) ([]string, error) {
	var ids []string
	err := r.with(ctx).Model(&Relation{}).
		Select("subject_id").
		Where("relation_type = ? AND target_type = ? AND target_id = ? AND status = ?",
			relationType, targetType, targetID, security.StatusEnabled).
		Order("sort asc, id asc").
		Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, nil
}

// DeleteBySubjectIDs 按主体 id 集合删除指定关系类型的关系（批量清按钮权限绑定用）。
func (r *Repo) DeleteBySubjectIDs(ctx context.Context, subjectType string, subjectIDs []string, relationType string) error {
	if len(subjectIDs) == 0 {
		return nil
	}
	q := r.with(ctx).Where("subject_type = ? AND subject_id IN ?", subjectType, subjectIDs)
	if relationType != "" {
		q = q.Where("relation_type = ?", relationType)
	}
	return q.Delete(&Relation{}).Error
}

// DeleteByTargetIDs 按目标 id 集合删除指定关系类型的关系（清理部门/角色/组被引用关系）。
func (r *Repo) DeleteByTargetIDs(ctx context.Context, targetType string, targetIDs []string, relationType string) error {
	if len(targetIDs) == 0 {
		return nil
	}
	q := r.with(ctx).Where("target_type = ? AND target_id IN ?", targetType, targetIDs)
	if relationType != "" {
		q = q.Where("relation_type = ?", relationType)
	}
	return q.Delete(&Relation{}).Error
}

// CreateInBatches 批量插入关系行（供事务内调用）。
func (r *Repo) CreateInBatches(db *gorm.DB, rows []Relation) error {
	return db.CreateInBatches(rows, 200).Error
}

// AssertResourcesExist 校验资源 ID 集合在指定资源表中全部存在（targetType 决定表）。
func (r *Repo) AssertResourcesExist(db *gorm.DB, targetType string, resourceIDs []string) error {
	if len(resourceIDs) == 0 {
		return nil
	}
	var count int64
	table := "sys_resource"
	if targetType == TargetClientResource {
		table = "sys_client_resource"
	}
	if err := db.Table(table).Where("id IN ?", resourceIDs).Count(&count).Error; err != nil {
		return err
	}
	if int(count) != len(uniqueStrings(resourceIDs)) {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ResolvePermissionResourceIDs 将权限键解析为 BUTTON/ACTION 资源 ID（对齐 hei-boot resolvePermissionResourceIds）。
func (r *Repo) ResolvePermissionResourceIDs(db *gorm.DB, targetType, accountType string, permissionKeys []string) ([]string, error) {
	if len(permissionKeys) == 0 {
		return []string{}, nil
	}
	permTable := "sys_resource"
	if targetType == TargetClientResource {
		permTable = "sys_client_resource"
	}
	var fromPerm []string
	q := db.Table("sys_iam_relation").
		Select("DISTINCT subject_id").
		Where("subject_type = ? AND relation_type = ? AND target_type = ? AND status = ? AND target_key IN ?",
			SubjectResource, RelResourcePermission, TargetPermission, security.StatusEnabled, permissionKeys)
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	if targetType == TargetClientResource {
		q = db.Table("sys_iam_relation").
			Select("DISTINCT subject_id").
			Where("subject_type = ? AND relation_type = ? AND target_type = ? AND status = ? AND target_key IN ?",
				SubjectClientResource, RelClientResourcePermission, TargetPermission, security.StatusEnabled, permissionKeys)
		if accountType != "" {
			q = q.Where("account_type = ?", accountType)
		}
	}
	if err := q.Scan(&fromPerm).Error; err != nil {
		return nil, err
	}
	var fromCode []string
	if err := db.Table(permTable).
		Select("id").
		Where("code IN ? AND resource_type IN ?", permissionKeys, []string{"BUTTON", "ACTION"}).
		Scan(&fromCode).Error; err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(fromPerm)+len(fromCode))
	for _, id := range append(fromPerm, fromCode...) {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	if len(out) < len(uniqueStrings(permissionKeys)) && len(out) == 0 && len(permissionKeys) > 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return out, nil
}

func uniqueStrings(items []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, s := range items {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
