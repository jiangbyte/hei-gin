// internal/modules/iam/group/grant.go 授权逻辑。
//
// Author: Charlie

package group

import (
	"context"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/view"
	"hei-gin/internal/modules/iam/role"
)

// OwnUsers 用户组已关联账号成员。
func (s *Service) OwnUsers(ctx context.Context, id string) (*view.OwnUserResult, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	accountIDs, err := s.rel.ListTargetIDs(ctx, relation.SubjectGroup, id, relation.RelGroupUser, "")
	if err != nil {
		return nil, err
	}
	users, err := view.LoadAccountViews(ctx, s.db, accountIDs)
	if err != nil {
		return nil, err
	}
	return &view.OwnUserResult{ID: id, Users: users, AccountIDs: accountIDs}, nil
}

// GrantUsers 全量替换用户组成员账号。
func (s *Service) GrantUsers(ctx context.Context, req GrantUserParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	accountTypes, err := view.LoadAccountTypes(ctx, s.db, req.AccountIDs)
	if err != nil {
		return err
	}
	if err := s.rel.ReplaceSubjectAccounts(ctx, relation.SubjectGroup, req.ID, relation.RelGroupUser, req.AccountIDs, accountTypes); err != nil {
		return err
	}
	s.invalidateAccounts(ctx, req.AccountIDs)
	return nil
}

// OwnRoles 用户组已拥有角色。
func (s *Service) OwnRoles(ctx context.Context, id, accountType string) (*OwnRoleResult, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	roleIDs, err := s.rel.ListTargetIDs(ctx, relation.SubjectGroup, id, relation.RelGroupRole, orAdmin(accountType))
	if err != nil {
		return nil, err
	}
	roles, err := s.roles.GetByIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}
	if roles == nil {
		roles = []role.Role{}
	}
	return &OwnRoleResult{ID: id, Roles: roles, RoleIDs: roleIDs}, nil
}

// GrantRoles 全量替换用户组角色授权。
func (s *Service) GrantRoles(ctx context.Context, req GrantRoleParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	affected, err := s.membersOf(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := s.rel.ReplaceTargetIDs(ctx, relation.SubjectGroup, req.ID, relation.RelGroupRole, relation.TargetRole, orAdmin(req.AccountType), req.RoleIDs); err != nil {
		return err
	}
	s.invalidateAccounts(ctx, affected)
	return nil
}

// OwnResources 用户组已拥有管理端资源授权。
func (s *Service) OwnResources(ctx context.Context, id, accountType string) (*OwnResourceResult, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	typ := orAdmin(accountType)
	modules, err := s.resources.ListGrantModules(ctx, typ)
	if err != nil {
		return nil, err
	}
	grants, err := s.rel.ListResourceGrants(ctx, relation.SubjectGroup, id, relation.RelGroupResource, relation.TargetResource, typ)
	if err != nil {
		return nil, err
	}
	return &OwnResourceResult{ID: id, Modules: modules, GrantInfoList: grants}, nil
}

// GrantResources 全量替换用户组管理端资源授权。
func (s *Service) GrantResources(ctx context.Context, req GrantResourceParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	affected, err := s.membersOf(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := s.rel.ReplaceResourceGrants(ctx, relation.SubjectGroup, req.ID, relation.RelGroupResource, relation.TargetResource, orAdmin(req.AccountType), req.GrantInfoList); err != nil {
		return err
	}
	s.invalidateAccounts(ctx, affected)
	return nil
}

// OwnClientResources 用户组已拥有客户端资源授权。
func (s *Service) OwnClientResources(ctx context.Context, id, accountType string) (*OwnClientResourceResult, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	typ := orAdmin(accountType)
	modules, err := s.clients.ListGrantModules(ctx, typ)
	if err != nil {
		return nil, err
	}
	grants, err := s.rel.ListResourceGrants(ctx, relation.SubjectGroup, id, relation.RelGroupClientResource, relation.TargetClientResource, typ)
	if err != nil {
		return nil, err
	}
	return &OwnClientResourceResult{ID: id, Modules: modules, GrantInfoList: grants}, nil
}

// GrantClientResources 全量替换用户组客户端资源授权。
func (s *Service) GrantClientResources(ctx context.Context, req GrantResourceParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	affected, err := s.membersOf(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := s.rel.ReplaceResourceGrants(ctx, relation.SubjectGroup, req.ID, relation.RelGroupClientResource, relation.TargetClientResource, orAdmin(req.AccountType), req.GrantInfoList); err != nil {
		return err
	}
	s.invalidateAccounts(ctx, affected)
	return nil
}

func orAdmin(t string) string {
	if t == "" {
		return string(security.AccountAdmin)
	}
	return t
}
