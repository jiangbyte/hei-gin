package admin

import (
	"context"
	"io"
	"path"
	"strings"

	contextx "hei-gin/framework/core/context"
	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/framework/platform/storage"
	"hei-gin/modules/shared"
	"gorm.io/gorm"
)

// Service 管理端用户中心服务。
//
// Author: Charlie
type Service struct {
	repo    *Repo
	storage *storage.Manager
}

// NewService 构造服务。
func NewService(db *gorm.DB, storage *storage.Manager) *Service {
	return &Service{repo: NewRepo(db), storage: storage}
}

// New 构建 user.admin 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB, d.Storage)
	return module.Module{
		Name:   "user.admin",
		Order:  30,
		Models: []any{&Profile{}},
		Routes: []module.RouteRegistrar{s.registerRoutes},
	}
}

// Me 当前用户概览。
func (s *Service) Me(ctx context.Context, sess *security.SessionPayload) (*MeResult, error) {
	if sess == nil {
		return nil, errUnauthorized
	}
	profile := s.getOrCreate(ctx, sess.AccountID)
	return &MeResult{
		AccountID:      sess.AccountID,
		AccountType:    sess.AccountType,
		Name:           profile.Name,
		Nickname:       profile.Nickname,
		Avatar:         profile.Avatar,
		RoleIDs:        sess.RoleIDs,
		DeptIDs:        sess.DeptIDs,
		GroupIDs:       sess.GroupIDs,
		PermissionKeys: sess.PermissionKeys,
		Profile:        profile,
	}, nil
}

// UpdateProfile 更新资料。
func (s *Service) UpdateProfile(ctx context.Context, accountID string, req ProfileUpdateParam) error {
	_ = s.getOrCreate(ctx, accountID)
	updates := map[string]any{"updated_by": accountID}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Signature != nil {
		updates["signature"] = *req.Signature
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}
	return s.repo.UpdateProfile(ctx, accountID, updates)
}

// UploadAvatar 上传头像并更新资料。
func (s *Service) UploadAvatar(ctx context.Context, accountID string, filename string, r io.Reader, size int64, contentType string) (string, error) {
	ext := path.Ext(filename)
	object := storage.ObjectKey("avatar/admin", idgen.Next()+ext)
	url, err := s.storage.Provider().Put(ctx, object, r, size, contentType)
	if err != nil {
		return "", err
	}
	_ = s.getOrCreate(ctx, accountID)
	_ = s.repo.UpdateProfile(ctx, accountID, map[string]any{"avatar": url, "updated_by": accountID})
	return url, nil
}

// UpdatePassword 更新密码。
func (s *Service) UpdatePassword(ctx context.Context, accountID string, req PasswordUpdateParam) error {
	acc, err := s.repo.GetAccountPassword(ctx, accountID)
	if err != nil {
		return err
	}
	if req.OldPassword != "" && !security.CheckPassword(acc.PasswordHash, req.OldPassword) {
		return errOldPassword
	}
	hash, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	return s.repo.UpdateAccountPassword(ctx, accountID, hash)
}

// UpdatePhone 更新手机号。
func (s *Service) UpdatePhone(ctx context.Context, accountID string, req PhoneUpdateParam) error {
	if err := s.verifyAccountPassword(ctx, accountID, req.Password); err != nil {
		return err
	}
	_ = s.getOrCreate(ctx, accountID)
	return s.repo.UpdateProfile(ctx, accountID, map[string]any{"phone": req.Phone, "updated_by": accountID})
}

// UpdateEmail 更新邮箱。
func (s *Service) UpdateEmail(ctx context.Context, accountID string, req EmailUpdateParam) error {
	if err := s.verifyAccountPassword(ctx, accountID, req.Password); err != nil {
		return err
	}
	_ = s.getOrCreate(ctx, accountID)
	return s.repo.UpdateProfile(ctx, accountID, map[string]any{
		"email": strings.TrimSpace(req.Email), "updated_by": accountID,
	})
}

// OrgInfo 组织关联信息。
func (s *Service) OrgInfo(sess *security.SessionPayload) OrgInfoResult {
	return OrgInfoResult{
		RoleIDs:  sess.RoleIDs,
		DeptIDs:  sess.DeptIDs,
		GroupIDs: sess.GroupIDs,
	}
}

func (s *Service) verifyAccountPassword(ctx context.Context, accountID, password string) error {
	acc, err := s.repo.GetAccountPassword(ctx, accountID)
	if err != nil {
		return err
	}
	if password != "" && !security.CheckPassword(acc.PasswordHash, password) {
		return errPassword
	}
	return nil
}

func (s *Service) getOrCreate(ctx context.Context, accountID string) *Profile {
	p, err := s.repo.GetProfile(ctx, accountID)
	if err == nil {
		return p
	}
	p = &Profile{AccountID: accountID, CreatedBy: &accountID, UpdatedBy: &accountID}
	_ = s.repo.CreateProfile(ctx, p)
	return p
}

var (
	errUnauthorized = &passErr{msg: "unauthorized"}
	errOldPassword  = &passErr{msg: "old password incorrect"}
	errPassword     = &passErr{msg: "password incorrect"}
)

type passErr struct{ msg string }

func (e *passErr) Error() string { return e.msg }

// SendCode 发送绑定验证码（短信/邮件通道；未配置时静默成功）。
func (s *Service) SendCode(ctx context.Context, channel, target string) error {
	return nil
}

// SessionFromContext 从 context 取会话（handler 用）。
func SessionFromContext(ctx context.Context) *security.SessionPayload {
	return contextx.Session(ctx)
}

// AccountIDFromContext 从 context 取账号 ID。
func AccountIDFromContext(ctx context.Context) string {
	return contextx.AccountID(ctx)
}