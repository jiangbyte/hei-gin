package account

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/iam/client"
	"hei-gin/internal/modules/iam/group"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/resource"
	"hei-gin/internal/modules/iam/role"
	"hei-gin/internal/modules/shared"
	adminuser "hei-gin/internal/modules/user/admin"
	portaluser "hei-gin/internal/modules/user/portal"
)

// Lookup ä¾› auth æŒ‰èº«ä»½æˆ– ID è§£æžè´¦å·ã€‚
//
// Author: Charlie
type Lookup interface {
	FindByIdentity(ctx context.Context, identityType, identifier string) (*Account, *Identity, error)
	GetByID(ctx context.Context, id string) (*Account, error)
}

// Service è´¦å·æœåŠ¡ï¼ˆèµ„æ–™ç» user æ¨¡å— Repoï¼ŒæŽˆæƒç» relation æ¨¡å—ï¼‰ã€‚
//
// Author: Charlie
type Service struct {
	repo      *Repo
	admin     *adminuser.Repo
	portal    *portaluser.Repo
	rel       *relation.Service
	roles     *role.Repo
	groups    *group.Repo
	resources *resource.Service
	clients   *client.Service
}

// NewService æž„é€ è´¦å·æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service {
	return &Service{
		repo:      NewRepo(db),
		admin:     adminuser.NewRepo(db),
		portal:    portaluser.NewRepo(db),
		rel:       relation.NewService(db),
		roles:     role.NewRepo(db),
		groups:    group.NewRepo(db),
		resources: resource.NewService(db),
		clients:   client.NewService(db),
	}
}

// New æž„å»º iam.account æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return s.withJobs(module.Module{
		Name:   "iam.account",
		Models: []any{&Account{}, &Identity{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	})
}

// AsLookup è¿”å›ž auth æŸ¥æ‰¾æŽ¥å£ã€‚
func (s *Service) AsLookup() Lookup { return s }

// FindByIdentity æŒ‰èº«ä»½ç±»åž‹ä¸Žæ ‡è¯†æŸ¥æ‰¾è´¦å·ã€‚
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

// GetByID æŒ‰ä¸»é”®æŸ¥è¯¢è´¦å·ã€‚
func (s *Service) GetByID(ctx context.Context, id string) (*Account, error) {
	return s.repo.GetByID(ctx, id)
}

// FindEnabledByIdentity è§£æžå·²å¯ç”¨è´¦å·çš„ç™»å½•èº«ä»½ï¼ˆå®žçŽ° AccountFinderï¼‰ã€‚
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

// EnsureSuperPermissions ä»Ž sys_iam_relation è§£æžè§’è‰²æƒé™é”®ä¸ŽæŽˆæƒã€‚
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

// GetEnabledAccount è¿”å›žå·²å¯ç”¨è´¦å·ç±»åž‹ã€‚
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

// UpdatePasswordHash æ›´æ–°å¯†ç å“ˆå¸Œã€‚
func (s *Service) UpdatePasswordHash(ctx context.Context, accountID, passwordHash string) error {
	return s.repo.UpdatePasswordHash(ctx, accountID, passwordHash)
}

// Create åˆ›å»ºè´¦å·å¹¶å†™å…¥å¯¹åº”ç«¯èµ„æ–™ã€‚
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

// Update æ›´æ–°è´¦å·ä¸Žèµ„æ–™ã€‚
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

// Delete å…ˆåˆ åŒç«¯èµ„æ–™ï¼Œå†åˆ èº«ä»½ä¸Žè´¦å·ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	if err := s.admin.DeleteByAccountIDs(ctx, ids); err != nil {
		return err
	}
	if err := s.portal.DeleteByAccountIDs(ctx, ids); err != nil {
		return err
	}
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail è´¦å·è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*AccountResult, error) {
	return s.loadDetail(ctx, id)
}

// Page åˆ†é¡µï¼›sess å¯é€‰ï¼Œä¼ å…¥æ—¶æŒ‰æ•°æ®èŒƒå›´è¿‡æ»¤ã€‚
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
