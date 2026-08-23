// internal/modules/iam/account/repo_permissions.go 资源授权权限键展开（对齐 hei-boot listPermissionGrantsByAccount）。
//
// Author: Charlie

package account

import (
	"context"
	"encoding/json"
	"strings"

	"gorm.io/datatypes"

	"hei-gin/internal/modules/iam/relation"
)

type resourceGrantRow struct {
	TargetID    string `gorm:"column:target_id"`
	GrantMode   string `gorm:"column:grant_mode"`
	SubjectType string `gorm:"column:subject_type"`
	SubjectID   string `gorm:"column:subject_id"`
	AccountType string `gorm:"column:account_type"`
}

type resourcePermissionRow struct {
	SubjectID          string         `gorm:"column:subject_id"`
	TargetKey          string         `gorm:"column:target_key"`
	DataScope          string         `gorm:"column:data_scope"`
	CustomScopeDeptIDs datatypes.JSON `gorm:"column:custom_scope_dept_ids"`
	AccountType        string         `gorm:"column:account_type"`
}

func expandsGrantMode(mode string) bool {
	mode = strings.TrimSpace(mode)
	return mode == "" || mode == relation.GrantCascade || mode == relation.GrantDirect
}

func matchesAccountType(accountType, relAccountType string) bool {
	relAccountType = strings.TrimSpace(relAccountType)
	if relAccountType == "" {
		return true
	}
	return strings.EqualFold(relAccountType, accountType)
}

func parseCustomDeptIDs(raw datatypes.JSON) []string {
	if len(raw) == 0 {
		return nil
	}
	var ids []string
	if err := json.Unmarshal(raw, &ids); err != nil {
		return nil
	}
	return ids
}

// ListExpandedPermissionGrants 展开 SUBJECT_RESOURCE_GRANT → RESOURCE_PERMISSION 权限授予。
func (r *Repo) ListExpandedPermissionGrants(ctx context.Context, accountID string, roleIDs, groupIDs []string, accountType string) ([]permRow, error) {
	var grants []resourceGrantRow
	db := r.with(ctx).Table("sys_iam_relation").
		Select("target_id, grant_mode, subject_type, subject_id, account_type").
		Where("relation_type = ? AND target_type = ? AND status = ?", relation.RelSubjectResourceGrant, relation.TargetResource, "ENABLED")
	cond := "((subject_type = ? AND subject_id = ?)"
	args := []any{relation.SubjectAccount, accountID}
	if len(groupIDs) > 0 {
		cond += " OR (subject_type = ? AND subject_id IN ?)"
		args = append(args, relation.SubjectGroup, groupIDs)
	}
	if len(roleIDs) > 0 {
		cond += " OR (subject_type = ? AND subject_id IN ?)"
		args = append(args, relation.SubjectRole, roleIDs)
	}
	cond += ")"
	if err := db.Where(cond, args...).Find(&grants).Error; err != nil {
		return nil, err
	}
	resourceIDs := map[string]struct{}{}
	type grantMeta struct {
		subjectType string
		subjectID   string
	}
	grantByResource := map[string][]grantMeta{}
	for _, g := range grants {
		if !expandsGrantMode(g.GrantMode) || !matchesAccountType(accountType, g.AccountType) {
			continue
		}
		if g.TargetID == "" {
			continue
		}
		resourceIDs[g.TargetID] = struct{}{}
		grantByResource[g.TargetID] = append(grantByResource[g.TargetID], grantMeta{
			subjectType: g.SubjectType,
			subjectID:   g.SubjectID,
		})
	}
	if len(resourceIDs) == 0 {
		return nil, nil
	}
	ids := make([]string, 0, len(resourceIDs))
	for id := range resourceIDs {
		ids = append(ids, id)
	}
	var perms []resourcePermissionRow
	if err := r.with(ctx).Table("sys_iam_relation").
		Select("subject_id, target_key, data_scope, custom_scope_dept_ids, account_type").
		Where("subject_type = ? AND relation_type = ? AND target_type = ? AND status = ?",
			relation.SubjectResource, relation.RelResourcePermission, relation.TargetPermission, "ENABLED").
		Where("subject_id IN ?", ids).
		Find(&perms).Error; err != nil {
		return nil, err
	}
	byKey := map[string]permRow{}
	for _, p := range perms {
		key := strings.TrimSpace(p.TargetKey)
		if key == "" || !matchesAccountType(accountType, p.AccountType) {
			continue
		}
		metaList := grantByResource[p.SubjectID]
		if len(metaList) == 0 {
			continue
		}
		if _, exists := byKey[key]; exists {
			continue
		}
		meta := metaList[0]
		byKey[key] = permRow{
			TargetKey:          key,
			DataScope:          p.DataScope,
			SourceID:           meta.subjectID,
			SourceType:         meta.subjectType,
			CustomScopeDeptIDs: parseCustomDeptIDs(p.CustomScopeDeptIDs),
		}
	}
	out := make([]permRow, 0, len(byKey))
	for _, row := range byKey {
		out = append(out, row)
	}
	return out, nil
}
