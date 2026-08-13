package account

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
	adminuser "hei-gin/modules/user/admin"
	portaluser "hei-gin/modules/user/portal"
)

// Lookup 供 auth 按身份或 ID 解析账号。
//
// Author: Charlie
type Lookup interface {
	FindByIdentity(ctx context.Context, identityType, identifier string) (*Account, *Identity, error)
	GetByID(ctx context.Context, id string) (*Account, error)
}

// Service 账号服务（资料经 user 模块 Repo）。
//
// Author: Charlie
type Service struct {
	repo   *Repo
	admin  *adminuser.Repo
	portal *portaluser.Repo
}

// NewService 构造账号服务。
func NewService(db *gorm.DB) *Service {
	return &Service{
		repo:   NewRepo(db),
		admin:  adminuser.NewRepo(db),
		portal: portaluser.NewRepo(db),
	}
}

// New 构建 iam.account 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return s.withJobs(module.Module{
		Name:   "iam.account",
		Models: []any{&Account{}, &Identity{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	})
}

// AsLookup 返回 auth 查找接口。
func (s *Service) AsLookup() Lookup { return s }

// FindByIdentity 按身份类型与标识查找账号。
func (s *Service) FindByIdentity(ctx context.Context, identityType, identifier string) (*Account, *Identity, error) {
	ident, err := s.repo.FindIdentity(ctx, identityType, identifier)
	if err != nil {
		return nil, nil, err
	}
	acc, err := s.repo.GetByID(ctx, ident.AccountID)
	if err != nil {
		return nil, nil, err
	}
	return acc, ident, nil
}

// GetByID 按主键查询账号。
func (s *Service) GetByID(ctx context.Context, id string) (*Account, error) {
	return s.repo.GetByID(ctx, id)
}

// FindEnabledByIdentity 解析已启用账号的登录身份（实现 AccountFinder）。
func (s *Service) FindEnabledByIdentity(ctx context.Context, accountType security.AccountType, identityType, identifier string) (accountID, passwordHash string, err error) {
	acc, _, err := s.FindByIdentity(ctx, identityType, identifier)
	if err != nil {
		return "", "", err
	}
	if acc.AccountType != string(accountType) || acc.AccountStatus != security.AccountStatusEnabled {
		return "", "", gorm.ErrRecordNotFound
	}
	return acc.ID, acc.PasswordHash, nil
}

// EnsureSuperPermissions 从 sys_iam_relation 解析角色权限键与授权。
func (s *Service) EnsureSuperPermissions(ctx context.Context, accountID string) (keys []string, grants []security.PermissionGrant, err error) {
	roleIDs, err := s.repo.ListRoleIDs(ctx, accountID)
	if err != nil {
		return nil, nil, err
	}
	rows, err := s.repo.ListRolePermissions(ctx, roleIDs)
	if err != nil {
		return nil, nil, err
	}
	seen := map[string]struct{}{}
	for _, r := range rows {
		if _, ok := seen[r.TargetKey]; ok {
			continue
		}
		seen[r.TargetKey] = struct{}{}
		keys = append(keys, r.TargetKey)
		grants = append(grants, security.PermissionGrant{
			PermissionKey: r.TargetKey,
			DataScope:     security.DataScope(r.DataScope),
			SourceType:    "ROLE",
			SourceID:      r.SourceID,
		})
	}
	if keys == nil {
		keys = []string{}
	}
	return keys, grants, nil
}

// GetEnabledAccount 返回已启用账号类型。
func (s *Service) GetEnabledAccount(ctx context.Context, accountID string) (security.AccountType, error) {
	acc, err := s.repo.GetByID(ctx, accountID)
	if err != nil {
		return "", err
	}
	if acc.AccountStatus != security.AccountStatusEnabled {
		return "", gorm.ErrRecordNotFound
	}
	return security.AccountType(acc.AccountType), nil
}

// UpdatePasswordHash 更新密码哈希。
func (s *Service) UpdatePasswordHash(ctx context.Context, accountID, passwordHash string) error {
	return s.repo.UpdatePasswordHash(ctx, accountID, passwordHash)
}

// Create 创建账号并写入对应端资料。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return err
	}
	st := req.AccountStatus
	if st == "" {
		st = security.AccountStatusEnabled
	}
	accID := idgen.Next()
	acc := Account{ID: accID, PasswordHash: hash, AccountType: req.AccountType, AccountStatus: st}
	ident := Identity{
		ID: idgen.Next(), AccountID: accID, IdentityType: IdentityAccount, Identifier: req.Account,
		Verified: true, IsPrimary: true, BindStatus: BindBound,
	}
	if err := s.repo.CreateAccount(ctx, acc, ident); err != nil {
		return err
	}
	if req.AccountType == string(security.AccountAdmin) {
		return s.admin.UpsertProfile(ctx, &adminuser.Profile{
			AccountID: accID, Name: req.Name, Nickname: req.Nickname, Avatar: req.Avatar,
			Signature: req.Signature, Phone: req.Phone, Email: req.Email, Remark: req.Remark,
		})
	}
	return s.portal.UpsertProfile(ctx, &portaluser.Profile{
		AccountID: accID, Name: req.Name, Nickname: req.Nickname, Avatar: req.Avatar,
		Signature: req.Signature, Phone: req.Phone, Email: req.Email,
	})
}

// Update 更新账号与资料。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	st := req.AccountStatus
	if st == "" {
		st = security.AccountStatusEnabled
	}
	updates := map[string]any{"account_type": req.AccountType, "account_status": st}
	if req.Password != nil && *req.Password != "" {
		hash, err := security.HashPassword(*req.Password)
		if err != nil {
			return err
		}
		updates["password_hash"] = hash
	}
	if err := s.repo.UpdateAccount(ctx, req.ID, updates, req.Account); err != nil {
		return err
	}
	if req.AccountType == string(security.AccountAdmin) {
		return s.admin.UpsertProfile(ctx, &adminuser.Profile{
			AccountID: req.ID, Name: req.Name, Nickname: req.Nickname, Avatar: req.Avatar,
			Signature: req.Signature, Phone: req.Phone, Email: req.Email, Remark: req.Remark,
		})
	}
	return s.portal.UpsertProfile(ctx, &portaluser.Profile{
		AccountID: req.ID, Name: req.Name, Nickname: req.Nickname, Avatar: req.Avatar,
		Signature: req.Signature, Phone: req.Phone, Email: req.Email,
	})
}

// Delete 先删双端资料，再删身份与账号。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	if err := s.admin.DeleteByAccountIDs(ctx, ids); err != nil {
		return err
	}
	if err := s.portal.DeleteByAccountIDs(ctx, ids); err != nil {
		return err
	}
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 账号详情。
func (s *Service) Detail(ctx context.Context, id string) (*AccountResult, error) {
	return s.loadDetail(ctx, id)
}

// Page 分页；sess 可选，传入时按数据范围过滤。
func (s *Service) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (records []AccountResult, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err := s.repo.PageAccounts(ctx, p, sess)
	if err != nil {
		return nil, 0, current, size, err
	}
	records = make([]AccountResult, 0, len(rows))
	for _, a := range rows {
		vo, err := s.loadDetail(ctx, a.ID)
		if err != nil {
			continue
		}
		records = append(records, *vo)
	}
	return records, total, current, size, nil
}

func (s *Service) loadDetail(ctx context.Context, id string) (*AccountResult, error) {
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	vo := &AccountResult{
		ID: acc.ID, AccountType: acc.AccountType, AccountStatus: acc.AccountStatus,
		CancelledAt: acc.CancelledAt, CancelledBy: acc.CancelledBy, CancelReason: acc.CancelReason,
		LastLoginIP: acc.LastLoginIP, LastLoginAddress: acc.LastLoginAddress, LastLoginTime: acc.LastLoginTime,
		LastLoginDevice: acc.LastLoginDevice, LatestLoginIP: acc.LatestLoginIP, LatestLoginAddress: acc.LatestLoginAddress,
		LatestLoginTime: acc.LatestLoginTime, LatestLoginDevice: acc.LatestLoginDevice,
		CreatedAt: acc.CreatedAt, CreatedBy: acc.CreatedBy, UpdatedAt: acc.UpdatedAt, UpdatedBy: acc.UpdatedBy,
	}
	if ident, err := s.repo.FindAccountIdentity(ctx, id); err == nil {
		vo.Account = ident.Identifier
	}
	if acc.AccountType == string(security.AccountAdmin) {
		if p, err := s.admin.GetProfile(ctx, id); err == nil {
			vo.Name, vo.Nickname, vo.Avatar, vo.Signature, vo.Phone, vo.Email, vo.Remark =
				p.Name, p.Nickname, p.Avatar, p.Signature, p.Phone, p.Email, p.Remark
		}
	} else if p, err := s.portal.GetProfile(ctx, id); err == nil {
		vo.Name, vo.Nickname, vo.Avatar, vo.Signature, vo.Phone, vo.Email =
			p.Name, p.Nickname, p.Avatar, p.Signature, p.Phone, p.Email
	}
	return vo, nil
}
