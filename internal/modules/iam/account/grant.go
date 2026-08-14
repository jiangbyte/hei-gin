// internal/modules/iam/account/grant.go 授权逻辑。
//
// Author: Charlie

package account

import (
	"context"

	"hei-gin/internal/modules/iam/group"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/role"
)

// OwnRoles 账号已拥有角色。
func (s *Service) OwnRoles(ctx context.Context, id string) (*OwnRoleResult, error) {
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	roleIDs, err := s.rel.ListTargetIDs(ctx, relation.SubjectAccount, id, relation.RelAccountRole, acc.AccountType)
	if err != nil {
		return nil, err
	}
	roles, err := s.loadRoles(ctx, roleIDs)
	if err != nil {
		return nil, err
	}
	return &OwnRoleResult{ID: id, Roles: roles, RoleIDs: roleIDs}, nil
}

// GrantRoles 全量替换账号角色授权。
func (s *Service) GrantRoles(ctx context.Context, req GrantRoleParam) error {
	acc, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	return s.rel.ReplaceTargetIDs(ctx, relation.SubjectAccount, req.ID, relation.RelAccountRole, relation.TargetRole, acc.AccountType, req.RoleIDs)
}

// OwnGroups 账号已拥有用户组。
func (s *Service) OwnGroups(ctx context.Context, id string) (*OwnGroupResult, error) {
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	groupIDs, err := s.rel.ListTargetIDs(ctx, relation.SubjectAccount, id, relation.RelAccountGroup, acc.AccountType)
	if err != nil {
		return nil, err
	}
	groups, err := s.loadGroups(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	return &OwnGroupResult{ID: id, Groups: groups, GroupIDs: groupIDs}, nil
}

// GrantGroups 全量替换账号用户组授权。
func (s *Service) GrantGroups(ctx context.Context, req GrantGroupParam) error {
	acc, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	return s.rel.ReplaceTargetIDs(ctx, relation.SubjectAccount, req.ID, relation.RelAccountGroup, relation.TargetGroup, acc.AccountType, req.GroupIDs)
}

// OwnDepts 账号已拥有部门授权。
func (s *Service) OwnDepts(ctx context.Context, id string) (*OwnDeptResult, error) {
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	grants, err := s.rel.ListDeptGrants(ctx, id, acc.AccountType)
	if err != nil {
		return nil, err
	}
	return &OwnDeptResult{ID: id, GrantInfoList: grants}, nil
}

// GrantDepts 全量替换账号部门授权。
func (s *Service) GrantDepts(ctx context.Context, req GrantDeptParam) error {
	acc, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	return s.rel.ReplaceDeptGrants(ctx, req.ID, acc.AccountType, req.GrantInfoList)
}

// OwnResources 账号已拥有管理端资源授权。
func (s *Service) OwnResources(ctx context.Context, id, accountType string) (*OwnResourceResult, error) {
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	typ := accountType
	if typ == "" {
		typ = acc.AccountType
	}
	modules, err := s.resources.ListGrantModules(ctx, typ)
	if err != nil {
		return nil, err
	}
	grants, err := s.rel.ListResourceGrants(ctx, relation.SubjectAccount, id, relation.RelAccountResource, relation.TargetResource, typ)
	if err != nil {
		return nil, err
	}
	return &OwnResourceResult{ID: id, Modules: modules, GrantInfoList: grants}, nil
}

// GrantResources 全量替换账号管理端资源授权。
func (s *Service) GrantResources(ctx context.Context, req GrantResourceParam) error {
	acc, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	typ := req.AccountType
	if typ == "" {
		typ = acc.AccountType
	}
	return s.rel.ReplaceResourceGrants(ctx, relation.SubjectAccount, req.ID, relation.RelAccountResource, relation.TargetResource, typ, req.GrantInfoList)
}

// OwnClientResources 账号已拥有客户端资源授权。
func (s *Service) OwnClientResources(ctx context.Context, id, accountType string) (*OwnClientResourceResult, error) {
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	typ := accountType
	if typ == "" {
		typ = acc.AccountType
	}
	modules, err := s.clients.ListGrantModules(ctx, typ)
	if err != nil {
		return nil, err
	}
	grants, err := s.rel.ListResourceGrants(ctx, relation.SubjectAccount, id, relation.RelAccountClientResource, relation.TargetClientResource, typ)
	if err != nil {
		return nil, err
	}
	return &OwnClientResourceResult{ID: id, Modules: modules, GrantInfoList: grants}, nil
}

// GrantClientResources 全量替换账号客户端资源授权。
func (s *Service) GrantClientResources(ctx context.Context, req GrantResourceParam) error {
	acc, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	typ := req.AccountType
	if typ == "" {
		typ = acc.AccountType
	}
	return s.rel.ReplaceResourceGrants(ctx, relation.SubjectAccount, req.ID, relation.RelAccountClientResource, relation.TargetClientResource, typ, req.GrantInfoList)
}

func (s *Service) loadRoles(ctx context.Context, ids []string) ([]role.Role, error) {
	rows, err := s.roles.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []role.Role{}
	}
	return rows, nil
}

func (s *Service) loadGroups(ctx context.Context, ids []string) ([]group.Group, error) {
	rows, err := s.groups.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []group.Group{}
	}
	return rows, nil
}
