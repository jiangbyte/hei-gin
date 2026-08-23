// internal/modules/iam/relation/service.go 业务服务。
//
// Author: Charlie

package relation

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
)

// New 构建 iam.relation 模块。
func New(_ *module.Deps) module.Module {
	return module.Module{
		Name:   "iam.relation",
		Models: []any{&Relation{}},
	}
}

// Service 关系服务：主体-目标关系查询与全量替换授权（先删后插，事务）。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造关系服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// ListSubjectIDsByTarget 按目标反查主体 ID（如组/角色成员账号）。
func (s *Service) ListSubjectIDsByTarget(ctx context.Context, relationType, targetType, targetID string) ([]string, error) {
	return s.repo.ListSubjectIDsByTarget(ctx, relationType, targetType, targetID)
}

// ReplaceGroupUsers 全量替换用户组成员（对齐 hei-boot replaceGroupUsers：ACCOUNT_GROUP + SUBJECT_ACCOUNT）。
func (s *Service) ReplaceGroupUsers(ctx context.Context, groupID string, accountIDs []string, accountTypes map[string]string) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteByTarget(tx, RelAccountGroup, TargetGroup, groupID); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(accountIDs))
		seen := map[string]struct{}{}
		for _, id := range accountIDs {
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			accType := accountTypes[id]
			if accType == "" {
				return gorm.ErrRecordNotFound
			}
			rows = append(rows, newRelation(SubjectAccount, id, accType, RelAccountGroup, TargetGroup, groupID))
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// ReplaceRoleUsers 全量替换角色成员（对齐 hei-boot replaceRoleUsers：ACCOUNT_ROLE + SUBJECT_ACCOUNT）。
func (s *Service) ReplaceRoleUsers(ctx context.Context, roleID string, accountIDs []string, accountTypes map[string]string) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteByTarget(tx, RelAccountRole, TargetRole, roleID); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(accountIDs))
		seen := map[string]struct{}{}
		for _, id := range accountIDs {
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			accType := accountTypes[id]
			if accType == "" {
				return gorm.ErrRecordNotFound
			}
			rows = append(rows, newRelation(SubjectAccount, id, accType, RelAccountRole, TargetRole, roleID))
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// ListTargetIDs 列出主体已关联目标 ID（accountType 为空不过滤）。
func (s *Service) ListTargetIDs(ctx context.Context, subjectType, subjectID, relationType, accountType string) ([]string, error) {
	rows, err := s.repo.ListRelations(ctx, subjectType, subjectID, relationType, accountType)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	seen := map[string]struct{}{}
	for _, r := range rows {
		if r.TargetID == "" {
			continue
		}
		if _, ok := seen[r.TargetID]; ok {
			continue
		}
		seen[r.TargetID] = struct{}{}
		out = append(out, r.TargetID)
	}
	return out, nil
}

// ReplaceTargetIDs 先删后插全量替换主体-目标简单关系（角色/用户组等）。
func (s *Service) ReplaceTargetIDs(ctx context.Context, subjectType, subjectID, relationType, targetType, accountType string, targetIDs []string) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteSubjectRelations(tx, subjectType, subjectID, relationType, accountType); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(targetIDs))
		seen := map[string]struct{}{}
		for _, id := range targetIDs {
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			rows = append(rows, newRelation(subjectType, subjectID, accountType, relationType, targetType, id))
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// ReplaceDeptGrants 先删后插全量替换账号-部门授予（is_primary 落库）。
func (s *Service) ReplaceDeptGrants(ctx context.Context, accountID, accountType string, grants []DeptGrantInfo) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteSubjectRelations(tx, SubjectAccount, accountID, RelAccountDept, accountType); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(grants))
		for _, g := range grants {
			if g.DeptID == "" {
				continue
			}
			rel := newRelation(SubjectAccount, accountID, accountType, RelAccountDept, TargetDept, g.DeptID)
			rel.IsPrimary = g.IsPrimary
			rows = append(rows, rel)
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// ReplaceResourceGrants 先删后插全量替换主体-资源授予（权限键展开为按钮资源，对齐 hei-boot replaceSubjectResourceGrants）。
func (s *Service) ReplaceResourceGrants(ctx context.Context, subjectType, subjectID, relationType, targetType, accountType string, grants []ResourceGrantInfo) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteSubjectRelations(tx, subjectType, subjectID, relationType, accountType); err != nil {
			return err
		}
		originalIDs := make([]string, 0, len(grants))
		permissionKeys := make([]string, 0)
		for _, g := range grants {
			if g.ResourceID != "" {
				originalIDs = append(originalIDs, g.ResourceID)
			}
			for _, key := range g.PermissionKeys {
				if key != "" {
					permissionKeys = append(permissionKeys, key)
				}
			}
		}
		originalIDs = uniqueStrings(originalIDs)
		permissionKeys = uniqueStrings(permissionKeys)
		if len(originalIDs) > 0 {
			if err := s.repo.AssertResourcesExist(tx, targetType, originalIDs); err != nil {
				return err
			}
		}
		resourceIDs := append([]string{}, originalIDs...)
		if len(permissionKeys) > 0 {
			expanded, err := s.repo.ResolvePermissionResourceIDs(tx, targetType, accountType, permissionKeys)
			if err != nil {
				return err
			}
			resourceIDs = append(resourceIDs, expanded...)
		}
		resourceIDs = uniqueStrings(resourceIDs)
		rows := make([]Relation, 0, len(resourceIDs))
		origSet := map[string]struct{}{}
		for _, id := range originalIDs {
			origSet[id] = struct{}{}
		}
		for _, id := range resourceIDs {
			rel := newRelation(subjectType, subjectID, accountType, relationType, targetType, id)
			if _, direct := origSet[id]; direct {
				rel.GrantMode = GrantDirect
			} else {
				rel.GrantMode = GrantCascade
			}
			rel.DataScope = string(security.DataScopeAll)
			rows = append(rows, rel)
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// ListDeptGrants 账号已拥有部门授予明细。
func (s *Service) ListDeptGrants(ctx context.Context, accountID, accountType string) ([]DeptGrantInfo, error) {
	rows, err := s.repo.ListRelations(ctx, SubjectAccount, accountID, RelAccountDept, accountType)
	if err != nil {
		return nil, err
	}
	out := make([]DeptGrantInfo, 0, len(rows))
	for _, r := range rows {
		out = append(out, DeptGrantInfo{DeptID: r.TargetID, IsPrimary: r.IsPrimary})
	}
	return out, nil
}

// ListResourceGrants 主体已拥有资源授予明细（按键分组回填 permission_keys）。
func (s *Service) ListResourceGrants(ctx context.Context, subjectType, subjectID, relationType, targetType, accountType string) ([]ResourceGrantInfo, error) {
	rows, err := s.repo.ListRelations(ctx, subjectType, subjectID, relationType, accountType)
	if err != nil {
		return nil, err
	}
	byResource := map[string][]string{}
	order := []string{}
	for _, r := range rows {
		if r.TargetID == "" {
			continue
		}
		if _, ok := byResource[r.TargetID]; !ok {
			order = append(order, r.TargetID)
		}
		if r.TargetKey != "" {
			byResource[r.TargetID] = append(byResource[r.TargetID], r.TargetKey)
		}
	}
	out := make([]ResourceGrantInfo, 0, len(order))
	for _, id := range order {
		out = append(out, ResourceGrantInfo{ResourceID: id, PermissionKeys: byResource[id]})
	}
	return out, nil
}

// ReplaceSubjectAccounts 先删后插全量替换主体-账号成员（GROUP_USER/ROLE_USER）。
func (s *Service) ReplaceSubjectAccounts(ctx context.Context, subjectType, subjectID, relationType string, accountIDs []string, accountTypes map[string]string) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteSubjectRelations(tx, subjectType, subjectID, relationType, ""); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(accountIDs))
		seen := map[string]struct{}{}
		for _, id := range accountIDs {
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			accType := accountTypes[id]
			if accType == "" {
				return gorm.ErrRecordNotFound
			}
			rows = append(rows, newRelation(subjectType, subjectID, accType, relationType, TargetAccount, id))
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// BindResourcePermissions 先删后插为资源绑定权限键（RESOURCE_PERMISSION/CLIENT_RESOURCE_PERMISSION）。
func (s *Service) BindResourcePermissions(ctx context.Context, subjectType, subjectID, relationType, accountType string, permissionKeys []string) error {
	return s.BindResourcePermissionDetail(ctx, subjectType, subjectID, relationType, accountType, permissionKeys, "", nil, 0, nil)
}

// BindResourcePermissionDetail 先删后插为资源绑定权限（带数据范围/自定义部门/排序/描述；对齐 hei-boot bindPermission）。
func (s *Service) BindResourcePermissionDetail(ctx context.Context, subjectType, subjectID, relationType, accountType string,
	permissionKeys []string, dataScope string, customScopeDeptIDs []string, sort int, description *string) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		keys := uniqueStrings(permissionKeys)
		if err := s.repo.deleteSubjectRelationsByKeys(tx, subjectType, subjectID, relationType, accountType, keys); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(permissionKeys))
		for _, key := range permissionKeys {
			if key == "" {
				continue
			}
			rel := newRelation(subjectType, subjectID, accountType, relationType, TargetPermission, "")
			rel.TargetKey = key
			rel.GrantMode = GrantCascade
			rel.DataScope = dataScope
			if rel.DataScope == "" {
				rel.DataScope = string(security.DataScopeAll)
			}
			if len(customScopeDeptIDs) > 0 {
				b, _ := json.Marshal(customScopeDeptIDs)
				rel.CustomScopeDeptIDs = b
			}
			rel.Sort = sort
			rel.Description = description
			rows = append(rows, rel)
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// DeleteBySubjectIDs 按主体 id 集合删除指定关系类型的关系。
func (s *Service) DeleteBySubjectIDs(ctx context.Context, subjectType string, subjectIDs []string, relationType string) error {
	return s.repo.DeleteBySubjectIDs(ctx, subjectType, subjectIDs, relationType)
}

// DeleteByTargetIDs 按目标 id 集合删除指定关系类型的关系。
func (s *Service) DeleteByTargetIDs(ctx context.Context, targetType string, targetIDs []string, relationType string) error {
	return s.repo.DeleteByTargetIDs(ctx, targetType, targetIDs, relationType)
}

// newRelation 构造默认启用关系行。
func newRelation(subjectType, subjectID, accountType, relationType, targetType, targetID string) Relation {
	return Relation{
		ID:                 idgen.Next(),
		SubjectType:        subjectType,
		SubjectID:          subjectID,
		AccountType:        accountType,
		RelationType:       relationType,
		TargetType:         targetType,
		TargetID:           targetID,
		GrantMode:          GrantCascade,
		DataScope:          string(security.DataScopeAll),
		CustomScopeDeptIDs: datatypes.JSON([]byte("[]")),
		Status:             security.StatusEnabled,
		Extra:              datatypes.JSON([]byte("{}")),
	}
}
