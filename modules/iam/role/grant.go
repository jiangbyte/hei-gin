package role

import (
	"context"

	"hei-gin/framework/core/security"
	"hei-gin/modules/iam/relation"
	"hei-gin/modules/iam/result"
)

// OwnUsers 角色已关联账号成员。
func (s *Service) OwnUsers(ctx context.Context, id string) (*result.OwnUserResult, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	accountIDs, err := s.rel.ListTargetIDs(ctx, relation.SubjectRole, id, relation.RelRoleUser, "")
	if err != nil {
		return nil, err
	}
	users, err := result.LoadAccountViews(ctx, s.db, accountIDs)
	if err != nil {
		return nil, err
	}
	return &result.OwnUserResult{ID: id, Users: users, AccountIDs: accountIDs}, nil
}

// GrantUsers 全量替换角色成员账号。
func (s *Service) GrantUsers(ctx context.Context, req GrantUserParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	accountTypes, err := result.LoadAccountTypes(ctx, s.db, req.AccountIDs)
	if err != nil {
		return err
	}
	return s.rel.ReplaceSubjectAccounts(ctx, relation.SubjectRole, req.ID, relation.RelRoleUser, req.AccountIDs, accountTypes)
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
	return s.rel.ReplaceResourceGrants(ctx, relation.SubjectRole, req.ID, relation.RelRoleResource, relation.TargetResource, orAdmin(req.AccountType), req.GrantInfoList)
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
	return s.rel.ReplaceResourceGrants(ctx, relation.SubjectRole, req.ID, relation.RelRoleClientResource, relation.TargetClientResource, orAdmin(req.AccountType), req.GrantInfoList)
}

func orAdmin(t string) string {
	if t == "" {
		return string(security.AccountAdmin)
	}
	return t
}
