// internal/modules/iam/account/authorization.go 账户授权聚合（对齐 hei-boot AccountApi.getAuthorization）。
//
// Author: Charlie

package account

import (
	"context"
	"strings"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/modules/iam/relation"
)

// SessionAuthorization 登录会话授权快照。
type SessionAuthorization = security.AuthorizationSnapshot

// GetSessionAuthorization 聚合账号角色/组织/资源/权限（供 auth issueSession 与 /me 使用）。
func (s *Service) GetSessionAuthorization(ctx context.Context, accountID string) (*security.AuthorizationSnapshot, error) {
	acc, err := s.repo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	roleIDs, err := s.repo.ListRoleIDs(ctx, accountID)
	if err != nil {
		return nil, err
	}
	groupIDs, err := s.repo.ListGroupIDs(ctx, accountID)
	if err != nil {
		return nil, err
	}
	groupRoleIDs, err := s.repo.ListRoleIDsByGroups(ctx, groupIDs, acc.AccountType)
	if err != nil {
		return nil, err
	}
	roleIDSet := map[string]struct{}{}
	for _, id := range roleIDs {
		roleIDSet[id] = struct{}{}
	}
	for _, id := range groupRoleIDs {
		if _, ok := roleIDSet[id]; !ok {
			roleIDSet[id] = struct{}{}
			roleIDs = append(roleIDs, id)
		}
	}
	deptIDs, err := s.repo.ListDeptIDs(ctx, accountID, acc.AccountType)
	if err != nil {
		return nil, err
	}
	roleCodes, err := s.repo.ListRoleCodes(ctx, roleIDs)
	if err != nil {
		return nil, err
	}
	resourceIDs, err := s.resources.ListGrantedResourceIDs(ctx, accountID, groupIDs, roleIDs, acc.AccountType)
	if err != nil {
		return nil, err
	}
	clientResourceIDs, err := s.repo.ListGrantedClientResourceIDs(ctx, accountID, groupIDs, roleIDs, acc.AccountType)
	if err != nil {
		return nil, err
	}
	keys, grants, err := s.EnsureSuperPermissions(ctx, accountID)
	if err != nil {
		return nil, err
	}
	clientKeys, err := s.repo.ListClientPermissionKeys(ctx, accountID, roleIDs, groupIDs, acc.AccountType)
	if err != nil {
		return nil, err
	}
	return &security.AuthorizationSnapshot{
		RoleIDs:              roleIDs,
		RoleCodes:            roleCodes,
		DeptIDs:              deptIDs,
		GroupIDs:             groupIDs,
		ResourceIDs:          resourceIDs,
		ClientResourceIDs:    clientResourceIDs,
		PermissionKeys:       keys,
		ClientPermissionKeys: clientKeys,
		PermissionGrants:     grants,
	}, nil
}

// AssignRegisterDefaults 注册后分配默认角色/部门（对齐 hei-boot assignRole/assignPrimaryDept）。
func (s *Service) AssignRegisterDefaults(ctx context.Context, accountID string, accountType security.AccountType) error {
	typeName := strings.ToUpper(string(accountType))
	roleID := strings.TrimSpace(s.runtime.GetString(ctx, "AUTH_REGISTER_"+typeName+"_DEFAULT_ROLE_ID", ""))
	deptID := strings.TrimSpace(s.runtime.GetString(ctx, "AUTH_REGISTER_"+typeName+"_DEFAULT_DEPT_ID", ""))
	if roleID != "" {
		if err := s.rel.ReplaceTargetIDs(ctx, relation.SubjectAccount, accountID, relation.RelAccountRole, relation.TargetRole, string(accountType), []string{roleID}); err != nil {
			return err
		}
	}
	if deptID != "" {
		if err := s.rel.ReplaceDeptGrants(ctx, accountID, string(accountType), []relation.DeptGrantInfo{{DeptID: deptID, IsPrimary: true}}); err != nil {
			return err
		}
	}
	return nil
}
