// internal/modules/iam/group/grant.go 授权逻辑。
//
// Author: Charlie

package group

import (
	"context"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/result"
	"hei-gin/internal/modules/iam/role"
)

// OwnUsers ç”¨æˆ·ç»„å·²å…³è”è´¦å·æˆå‘˜ã€‚
func (s *Service) OwnUsers(ctx context.Context, id string) (*result.OwnUserResult, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	accountIDs, err := s.rel.ListTargetIDs(ctx, relation.SubjectGroup, id, relation.RelGroupUser, "")
	if err != nil {
		return nil, err
	}
	users, err := result.LoadAccountViews(ctx, s.db, accountIDs)
	if err != nil {
		return nil, err
	}
	return &result.OwnUserResult{ID: id, Users: users, AccountIDs: accountIDs}, nil
}

// GrantUsers å…¨é‡æ›¿æ¢ç”¨æˆ·ç»„æˆå‘˜è´¦å·ã€‚
func (s *Service) GrantUsers(ctx context.Context, req GrantUserParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	accountTypes, err := result.LoadAccountTypes(ctx, s.db, req.AccountIDs)
	if err != nil {
		return err
	}
	return s.rel.ReplaceSubjectAccounts(ctx, relation.SubjectGroup, req.ID, relation.RelGroupUser, req.AccountIDs, accountTypes)
}

// OwnRoles ç”¨æˆ·ç»„å·²æ‹¥æœ‰è§’è‰²ã€‚
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

// GrantRoles å…¨é‡æ›¿æ¢ç”¨æˆ·ç»„è§’è‰²æŽˆæƒã€‚
func (s *Service) GrantRoles(ctx context.Context, req GrantRoleParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	return s.rel.ReplaceTargetIDs(ctx, relation.SubjectGroup, req.ID, relation.RelGroupRole, relation.TargetRole, orAdmin(req.AccountType), req.RoleIDs)
}

// OwnResources ç”¨æˆ·ç»„å·²æ‹¥æœ‰ç®¡ç†ç«¯èµ„æºæŽˆæƒã€‚
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

// GrantResources å…¨é‡æ›¿æ¢ç”¨æˆ·ç»„ç®¡ç†ç«¯èµ„æºæŽˆæƒã€‚
func (s *Service) GrantResources(ctx context.Context, req GrantResourceParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	return s.rel.ReplaceResourceGrants(ctx, relation.SubjectGroup, req.ID, relation.RelGroupResource, relation.TargetResource, orAdmin(req.AccountType), req.GrantInfoList)
}

// OwnClientResources ç”¨æˆ·ç»„å·²æ‹¥æœ‰å®¢æˆ·ç«¯èµ„æºæŽˆæƒã€‚
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

// GrantClientResources å…¨é‡æ›¿æ¢ç”¨æˆ·ç»„å®¢æˆ·ç«¯èµ„æºæŽˆæƒã€‚
func (s *Service) GrantClientResources(ctx context.Context, req GrantResourceParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return err
	}
	return s.rel.ReplaceResourceGrants(ctx, relation.SubjectGroup, req.ID, relation.RelGroupClientResource, relation.TargetClientResource, orAdmin(req.AccountType), req.GrantInfoList)
}

func orAdmin(t string) string {
	if t == "" {
		return string(security.AccountAdmin)
	}
	return t
}
