// internal/modules/iam/account/service.go 业务服务。
//
// Author: Charlie

package account

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/runtimecfg"
	"hei-gin/internal/framework/platform/storage"
	"hei-gin/internal/modules/auth/oauth"
	"hei-gin/internal/modules/iam/client"
	"hei-gin/internal/modules/iam/group"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/resource"
	"hei-gin/internal/modules/iam/role"
	"hei-gin/internal/modules/profile"
	"hei-gin/internal/modules/profile/identity"
)

// Service 账号服务（资料经 user 模块 Repo，授权经 relation 模块）。
//
// Author: Charlie
type Service struct {
	repo      *Repo
	admin     *profile.Repo
	portal    *profile.Repo
	rel       *relation.Service
	roles     *role.Repo
	groups    *group.Repo
	resources *resource.Service
	clients   *client.Service
	runtime   *runtimecfg.Settings
	sessions  *security.SessionStore
	storage   *storage.Manager
	identity  *identity.Service
}

// NewService 构造账号服务。
func NewService(db *gorm.DB, rdb *redis.Client, rt *runtimecfg.Settings, sto *storage.Manager, idSvc *identity.Service) *Service {
	return &Service{
		repo:      NewRepo(db, rdb),
		admin:     profile.AdminRepo(db),
		portal:    profile.PortalRepo(db),
		rel:       relation.NewService(db),
		roles:     role.NewRepo(db),
		groups:    group.NewRepo(db),
		resources: resource.NewService(db),
		clients:   client.NewService(db),
		runtime:   rt,
		storage:   sto,
		identity:  idSvc,
	}
}

// New 构建 iam.account 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB, d.Redis, d.Runtime, d.Storage, identity.FromDeps(d))
	s.sessions = d.Sessions
	d.Provide(AccountFinderKey, s)
	return s.withJobs(module.Module{
		Name:   "iam.account",
		Models: []any{&Account{}, &Identity{}, &security.PasswordHistory{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	})
}

// invalidateSessions 授权变更后强制该账号下线（对齐 hei-boot LoginHelper.logoutAccount）。
func (s *Service) invalidateSessions(ctx context.Context, accountID string) {
	if s.sessions == nil || accountID == "" {
		return
	}
	_ = s.sessions.DeleteAllForAccountAnyType(ctx, accountID)
}

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

// EnsureSuperPermissions 从 sys_iam_relation 解析角色权限键与授权，并展开资源授权（按钮）权限键。
// 超管（内置 superadmin 账号或持有 SUPER_ADMIN 角色）补通配 *:*:*，对齐 hei-boot IamRelationServiceImpl。
func (s *Service) EnsureSuperPermissions(ctx context.Context, accountID string) (keys []string, grants []security.PermissionGrant, err error) {
	roleIDs, err := s.repo.ListRoleIDs(ctx, accountID)
	if err != nil {
		return nil, nil, err
	}
	groupIDs, err := s.repo.ListGroupIDs(ctx, accountID)
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
	acc, err := s.repo.GetByID(ctx, accountID)
	if err != nil {
		return nil, nil, err
	}
	expanded, err := s.repo.ListExpandedPermissionGrants(ctx, accountID, roleIDs, groupIDs, acc.AccountType)
	if err != nil {
		return nil, nil, err
	}
	for _, r := range expanded {
		if _, ok := seen[r.TargetKey]; ok {
			continue
		}
		seen[r.TargetKey] = struct{}{}
		keys = append(keys, r.TargetKey)
		grants = append(grants, security.PermissionGrant{
			PermissionKey:      r.TargetKey,
			DataScope:          security.DataScope(r.DataScope),
			CustomScopeDeptIDs: r.CustomScopeDeptIDs,
			SourceType:         r.SourceType,
			SourceID:           r.SourceID,
		})
	}
	if s.isSuperAdmin(ctx, accountID, roleIDs) {
		if _, ok := seen["*:*:*"]; !ok {
			seen["*:*:*"] = struct{}{}
			keys = append(keys, "*:*:*")
			grants = append(grants, security.PermissionGrant{
				PermissionKey: "*:*:*",
				DataScope:     security.DataScopeAll,
				SourceType:    "SUPER_ADMIN",
				SourceID:      accountID,
			})
		}
	}
	if keys == nil {
		keys = []string{}
	}
	return keys, grants, nil
}

// isSuperAdmin 判定超管：内置 superadmin 账号（identity=superadmin）或持有 SUPER_ADMIN 角色。
func (s *Service) isSuperAdmin(ctx context.Context, accountID string, roleIDs []string) bool {
	var n int64
	if err := s.repo.DB().WithContext(ctx).Table("sys_account_identity").
		Where("account_id = ? AND identity_type = ? AND identifier = ?", accountID, "ACCOUNT", "superadmin").
		Count(&n).Error; err == nil && n > 0 {
		return true
	}
	if len(roleIDs) == 0 {
		return false
	}
	var m int64
	if err := s.repo.DB().WithContext(ctx).Table("sys_role").
		Where("id IN ? AND code = ?", roleIDs, "SUPER_ADMIN").Count(&m).Error; err == nil {
		return m > 0
	}
	return false
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

// HasBoundIdentity 是否已绑定 EMAIL/PHONE 身份。
func (s *Service) HasBoundIdentity(ctx context.Context, accountID, identityType string) bool {
	idents, err := s.repo.FindIdentities(ctx, accountID)
	if err != nil {
		return false
	}
	for _, id := range idents {
		if id.IdentityType == identityType && id.BindStatus == BindBound && strings.TrimSpace(id.Identifier) != "" {
			return true
		}
	}
	return false
}

// Create 创建账号（对齐 hei-boot AccountServiceImpl.create：RSA 解密密码、校验、全量身份）。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	accountLogin, err := security.RequireAccountLogin(req.Account)
	if err != nil {
		return err
	}
	req.Account = accountLogin
	if strings.EqualFold(req.AccountStatus, "CANCELLED") {
		return fmt.Errorf("注销状态不允许通过管理端设置")
	}
	accountType := strings.ToUpper(strings.TrimSpace(req.AccountType))
	if accountType != "ADMIN" && accountType != "PORTAL" {
		return fmt.Errorf("unsupported account type: %s", req.AccountType)
	}
	if _, err := s.repo.FindIdentity(ctx, IdentityAccount, req.Account); err == nil {
		return fmt.Errorf("account identifier already exists")
	}
	rawPassword, err := s.resolveCreatePassword(ctx, req.Password, req.PasswordKeyID)
	if err != nil {
		return err
	}
	hash, err := security.HashPassword(rawPassword)
	if err != nil {
		return err
	}
	st := req.AccountStatus
	if st == "" {
		st = security.AccountStatusEnabled
	}
	accID := idgen.Next()
	acc := Account{ID: accID, PasswordHash: hash, AccountType: accountType, AccountStatus: st}
	ident := Identity{
		ID: idgen.Next(), AccountID: accID, IdentityType: IdentityAccount, Identifier: req.Account,
		Verified: true, IsPrimary: true, BindStatus: BindBound,
	}
	if err := s.repo.CreateAccount(ctx, acc, ident); err != nil {
		return err
	}
	_ = s.recordHistory(ctx, accID, rawPassword, accID, "admin_reset")
	if accountType == string(security.AccountAdmin) {
		return s.admin.UpsertProfile(ctx, &profile.Profile{
			AccountID: accID, Nickname: req.Nickname, Avatar: req.Avatar,
			Signature: req.Signature, Phone: req.Phone, Email: req.Email, Remark: req.Remark,
		})
	}
	return s.portal.UpsertProfile(ctx, &profile.Profile{
		AccountID: accID, Nickname: req.Nickname, Avatar: req.Avatar,
		Signature: req.Signature, Phone: req.Phone, Email: req.Email,
	})
}

// Update 更新账号与资料（对齐 hei-boot AccountServiceImpl.update：可选改密 + 全量替换身份）。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	accountLogin, err := security.RequireAccountLogin(req.Account)
	if err != nil {
		return err
	}
	req.Account = accountLogin
	acc, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("account not found")
	}
	if err := s.assertAccountAccessible(ctx, req.ID); err != nil {
		return err
	}
	if strings.EqualFold(acc.AccountStatus, "CANCELLED") {
		return fmt.Errorf("已注销账号不允许通过管理端修改")
	}
	if strings.EqualFold(req.AccountStatus, "CANCELLED") {
		return fmt.Errorf("注销状态不允许通过管理端设置")
	}
	accountType := strings.ToUpper(strings.TrimSpace(req.AccountType))
	if accountType != "ADMIN" && accountType != "PORTAL" {
		return fmt.Errorf("unsupported account type: %s", req.AccountType)
	}
	if existing, err := s.repo.FindIdentity(ctx, IdentityAccount, req.Account); err == nil && existing.AccountID != req.ID {
		return fmt.Errorf("account identifier already exists")
	}
	st := req.AccountStatus
	if st == "" {
		st = security.AccountStatusEnabled
	}
	updates := map[string]any{"account_type": accountType, "account_status": st}
	if req.Password != nil && *req.Password != "" {
		rawPassword, derr := s.repo.DecryptPassword(ctx, req.PasswordKeyID, *req.Password)
		if derr != nil {
			return derr
		}
		hash, herr := security.HashPassword(rawPassword)
		if herr != nil {
			return herr
		}
		updates["password_hash"] = hash
		_ = s.recordHistory(ctx, req.ID, rawPassword, req.ID, "admin_reset")
	}
	if err := s.repo.UpdateAccount(ctx, req.ID, updates, req.Account); err != nil {
		return err
	}
	if accountType == string(security.AccountAdmin) {
		return s.admin.UpsertProfile(ctx, &profile.Profile{
			AccountID: req.ID, Nickname: req.Nickname, Avatar: req.Avatar,
			Signature: req.Signature, Phone: req.Phone, Email: req.Email, Remark: req.Remark,
		})
	}
	return s.portal.UpsertProfile(ctx, &profile.Profile{
		AccountID: req.ID, Nickname: req.Nickname, Avatar: req.Avatar,
		Signature: req.Signature, Phone: req.Phone, Email: req.Email,
	})
}

// UpdateLoginIdentity 更新邮箱/手机号登录身份。
func (s *Service) UpdateLoginIdentity(ctx context.Context, req UpdateLoginIdentityParam) error {
	acc, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("account not found")
	}
	if err := s.assertAccountAccessible(ctx, req.ID); err != nil {
		return err
	}
	if strings.EqualFold(acc.AccountStatus, "CANCELLED") {
		return fmt.Errorf("已注销账号不允许通过管理端修改")
	}
	if boolOf(req.EmailLoginEnabled) && !hasText(req.Email) {
		return fmt.Errorf("email login requires an email")
	}
	if boolOf(req.PhoneLoginEnabled) && !hasText(req.Phone) {
		return fmt.Errorf("phone login requires a phone")
	}
	return s.replaceAccountIdentities(ctx, req.ID,
		boolOf(req.EmailLoginEnabled), req.Email, req.Email,
		true, BindBound,
		boolOf(req.PhoneLoginEnabled), req.Phone, req.Phone,
		true, BindBound)
}

// resolveCreatePassword 管理端建号密码：RSA 解密，空则回退 AUTH_DEFAULT_PASSWORD（对齐 hei-boot resolveCreatePassword）。
func (s *Service) resolveCreatePassword(ctx context.Context, password, passwordKeyID string) (string, error) {
	raw := ""
	if password != "" {
		dec, err := s.repo.DecryptPassword(ctx, passwordKeyID, password)
		if err != nil {
			return "", err
		}
		raw = dec
	}
	if raw == "" && s.runtime != nil {
		raw = s.runtime.GetString(ctx, "AUTH_DEFAULT_PASSWORD", "")
	}
	if strings.TrimSpace(raw) == "" {
		return "", fmt.Errorf("password is required")
	}
	return raw, nil
}

// replaceAccountIdentities 全量替换 EMAIL/PHONE 登录身份（对齐 hei-boot replaceAccountIdentities）。
func (s *Service) replaceAccountIdentities(ctx context.Context, accountID string,
	emailLoginEnabled bool, emailIdentity, email *string, emailVerified bool, emailBindStatus string,
	phoneLoginEnabled bool, phoneIdentity, phone *string, phoneVerified bool, phoneBindStatus string) error {
	if err := s.repo.DeleteIdentitiesExcept(ctx, accountID, IdentityAccount); err != nil {
		return err
	}
	emailID := firstNonNil(emailIdentity, email)
	if emailLoginEnabled && hasText(emailID) {
		row := Identity{
			ID: idgen.Next(), AccountID: accountID, IdentityType: IdentityEmail,
			Identifier: strings.ToLower(strings.TrimSpace(*emailID)), Verified: emailVerified,
			IsPrimary: false, BindStatus: emailBindStatus,
		}
		if err := s.repo.CreateIdentity(ctx, &row); err != nil {
			return err
		}
	}
	phoneID := firstNonNil(phoneIdentity, phone)
	if phoneLoginEnabled && hasText(phoneID) {
		row := Identity{
			ID: idgen.Next(), AccountID: accountID, IdentityType: IdentityPhone,
			Identifier: strings.TrimSpace(*phoneID), Verified: phoneVerified,
			IsPrimary: false, BindStatus: phoneBindStatus,
		}
		if err := s.repo.CreateIdentity(ctx, &row); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) recordHistory(ctx context.Context, accountID, rawPassword, changedBy, reason string) error {
	policy := security.NewPasswordPolicy(s.repo.DB(), s.runtime)
	return policy.RecordHistory(ctx, accountID, rawPassword, changedBy, reason)
}

func boolOf(p *bool) bool { return p != nil && *p }

func strOr(p *string, def string) string {
	if p == nil || *p == "" {
		return def
	}
	return *p
}

func firstNonNil(a, b *string) *string {
	if a != nil && strings.TrimSpace(*a) != "" {
		return a
	}
	if b != nil && strings.TrimSpace(*b) != "" {
		return b
	}
	return nil
}

func hasText(p *string) bool { return p != nil && strings.TrimSpace(*p) != "" }

// Delete 先删关联授权关系，再删双端资料与身份、账号（对齐 hei-boot cleanupAccountSideData）。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := s.assertAccountAccessible(ctx, id); err != nil {
			return err
		}
	}
	_ = s.rel.DeleteBySubjectIDs(ctx, "ACCOUNT", ids, "")
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
	if err := s.assertAccountAccessible(ctx, id); err != nil {
		return nil, err
	}
	return s.loadDetail(ctx, id)
}

// Page 分页；sess 可选，传入时按数据范围过滤。批量加载身份与资料，避免 N+1。
func (s *Service) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (records []AccountListResult, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err := s.repo.PageAccounts(ctx, p, sess)
	if err != nil {
		return nil, 0, current, size, err
	}
	if len(rows) == 0 {
		return []AccountListResult{}, total, current, size, nil
	}
	ids := make([]string, 0, len(rows))
	for _, a := range rows {
		ids = append(ids, a.ID)
	}
	idents, _ := s.repo.FindAccountIdentities(ctx, ids)
	adminProfiles, _ := s.admin.ListByAccountIDs(ctx, ids)
	portalProfiles, _ := s.portal.ListByAccountIDs(ctx, ids)

	records = make([]AccountListResult, 0, len(rows))
	for i := range rows {
		a := &rows[i]
		vo := AccountListResult{
			ID: a.ID, AccountType: a.AccountType, AccountStatus: a.AccountStatus,
			LatestLoginTime: a.LatestLoginTime, UpdatedAt: a.UpdatedAt,
		}
		vo.Account = idents[a.ID]
		var avatar *string
		if a.AccountType == string(security.AccountAdmin) {
			if adminP := adminProfiles[a.ID]; adminP != nil {
				vo.Nickname, avatar, vo.Phone, vo.Email, vo.Remark =
					adminP.Nickname, adminP.Avatar, adminP.Phone, adminP.Email, adminP.Remark
			}
		} else if portalP := portalProfiles[a.ID]; portalP != nil {
			vo.Nickname, avatar, vo.Phone, vo.Email =
				portalP.Nickname, portalP.Avatar, portalP.Phone, portalP.Email
		}
		vo.Avatar = s.resolveAvatar(ctx, avatar)
		records = append(records, vo)
	}
	return records, total, current, size, nil
}

func (s *Service) applyProfiles(vo *AccountResult, accountType string, adminP, portalP *profile.Profile) {
	if accountType == string(security.AccountAdmin) {
		if adminP != nil {
			vo.Nickname, vo.Avatar, vo.Signature, vo.Phone, vo.Email, vo.Remark =
				adminP.Nickname, adminP.Avatar, adminP.Signature, adminP.Phone, adminP.Email, adminP.Remark
		}
		return
	}
	if portalP != nil {
		vo.Nickname, vo.Avatar, vo.Signature, vo.Phone, vo.Email =
			portalP.Nickname, portalP.Avatar, portalP.Signature, portalP.Phone, portalP.Email
	}
}

func (s *Service) applyIdentities(vo *AccountResult, idents []Identity) {
	vo.Identities = []IdentityResult{}
	for _, it := range idents {
		vo.Identities = append(vo.Identities, IdentityResult{
			ID: it.ID, AccountID: it.AccountID, IdentityType: it.IdentityType,
			Identifier: it.Identifier, Verified: it.Verified, IsPrimary: it.IsPrimary,
			BindStatus: it.BindStatus, CreatedAt: it.CreatedAt, CreatedBy: it.CreatedBy,
			UpdatedAt: it.UpdatedAt, UpdatedBy: it.UpdatedBy,
		})
		switch it.IdentityType {
		case IdentityEmail:
			vo.EmailLoginEnabled = true
			vo.EmailIdentity = &it.Identifier
			vo.EmailIdentityVerified = it.Verified
			vo.EmailIdentityBindStatus = &it.BindStatus
		case IdentityPhone:
			vo.PhoneLoginEnabled = true
			vo.PhoneIdentity = &it.Identifier
			vo.PhoneIdentityVerified = it.Verified
			vo.PhoneIdentityBindStatus = &it.BindStatus
		}
	}
}

func (s *Service) applyOAuthBindings(vo *AccountResult, binds []oauth.AccountOAuthBinding) {
	vo.OAuthBindings = []OAuthBindingResult{}
	for _, b := range binds {
		vo.OAuthBindings = append(vo.OAuthBindings, OAuthBindingResult{
			ID: b.ID, Provider: b.Provider, OpenID: b.OpenID, UnionID: b.UnionID,
			Nickname: b.Nickname, Avatar: b.Avatar, BoundAt: &b.BoundAt,
		})
	}
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
	if idents, err := s.repo.FindIdentities(ctx, id); err == nil {
		s.applyIdentities(vo, idents)
	}
	if binds, err := s.repo.FindOAuthBindings(ctx, id); err == nil {
		s.applyOAuthBindings(vo, binds)
	}
	if acc.AccountType == string(security.AccountAdmin) {
		if p, err := s.admin.GetProfile(ctx, id); err == nil {
			s.applyProfiles(vo, acc.AccountType, p, nil)
		}
	} else if p, err := s.portal.GetProfile(ctx, id); err == nil {
		s.applyProfiles(vo, acc.AccountType, nil, p)
	}
	vo.Avatar = s.resolveAvatar(ctx, vo.Avatar)
	if s.identity != nil {
		if status, err := s.identity.GetStatus(ctx, id); err == nil {
			vo.IdentityStatus = status
			if status != nil && status.RealNameMasked != "" {
				name := status.RealNameMasked
				vo.Name = &name
			}
		}
	}
	return vo, nil
}

// resolveAvatar 把头像对象引用解析为可访问 URL（对齐 profile / fastapi resolve_access_url）。
func (s *Service) resolveAvatar(ctx context.Context, value *string) *string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	if s.storage == nil {
		return value
	}
	u := s.storage.ResolveURL(ctx, *value)
	if u == "" {
		return nil
	}
	return &u
}
