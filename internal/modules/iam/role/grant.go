// internal/modules/iam/role/grant.go 授权逻辑。
//
// Author: Charlie

package role

import (
	"context"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/view"
)

// OwnUsers 角色已关联账号成员。
func (s *Service) OwnUsers(ctx context.Context, id string) (*view.OwnUserResult, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	accountIDs, err := s.rel.ListTargetIDs(ctx, relation.SubjectRole, id, relation.RelRoleUser, "")
	if err != nil {
		return nil, err
	}
	users, err := view.LoadAccountViews(ctx, s.db, accountIDs)
	if err != nil {
		return nil, err
	}
	return &view.OwnUserResult{ID: id, Users: users, AccountIDs: accountIDs}, nil
}

// GrantUsers 全量替换角色成员账号。
func (s *Service) GrantUsers(ctx context.Context, req GrantUserParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	accountTypes, err := view.LoadAccountTypes(ctx, s.db, req.AccountIDs)
	if err != nil {
		return err
	}
	if err := s.rel.ReplaceSubjectAccounts(ctx, relation.SubjectRole, req.ID, relation.RelRoleUser, req.AccountIDs, accountTypes); err != nil {
		return err
	}
	s.invalidateAccounts(ctx, req.AccountIDs)
	return nil
}

// OwnResources 角色已拥有管理端资源授权。
func (s *Service) OwnResources(ctx context.Context, id, accountType string) (*OwnResourceResult, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	typ := orAdmin(accountType)
	modules, err := s.resources.ListGrantModules(ctx, typ)
	if err != nil {
		return nil, err
	}
	grants, err := s.rel.ListResourceGrants(ctx, relation.SubjectRole, id, relation.RelRoleResource, relation.TargetResource, typ)
	if err != nil {
		return nil, err
	}
	return &OwnResourceResult{ID: id, Modules: modules, GrantInfoList: grants}, nil
}

// GrantResources 全量替换角色管理端资源授权。
func (s *Service) GrantResources(ctx context.Context, req GrantResourceParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	// 先取受影响成员（旧成员+新成员），授权变更后强制下线
	affected, err := s.membersOf(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := s.rel.ReplaceResourceGrants(ctx, relation.SubjectRole, req.ID, relation.RelRoleResource, relation.TargetResource, orAdmin(req.AccountType), req.GrantInfoList); err != nil {
		return err
	}
	s.invalidateAccounts(ctx, affected)
	return nil
}

// membersOf 列出角色当前成员账号 ID。
func (s *Service) membersOf(ctx context.Context, roleID string) ([]string, error) {
	return s.rel.ListTargetIDs(ctx, relation.SubjectRole, roleID, relation.RelRoleUser, "")
}

// OwnClientResources 角色已拥有客户端资源授权。
func (s *Service) OwnClientResources(ctx context.Context, id, accountType string) (*OwnClientResourceResult, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	typ := orAdmin(accountType)
	modules, err := s.clients.ListGrantModules(ctx, typ)
	if err != nil {
		return nil, err
	}
	grants, err := s.rel.ListResourceGrants(ctx, relation.SubjectRole, id, relation.RelRoleClientResource, relation.TargetClientResource, typ)
	if err != nil {
		return nil, err
	}
	return &OwnClientResourceResult{ID: id, Modules: modules, GrantInfoList: grants}, nil
}

// GrantClientResources 全量替换角色客户端资源授权。
func (s *Service) GrantClientResources(ctx context.Context, req GrantResourceParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	affected, err := s.membersOf(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := s.rel.ReplaceResourceGrants(ctx, relation.SubjectRole, req.ID, relation.RelRoleClientResource, relation.TargetClientResource, orAdmin(req.AccountType), req.GrantInfoList); err != nil {
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
