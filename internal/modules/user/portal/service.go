package portal

import (
	"context"
	"io"
	"path"
	"strings"

	"gorm.io/gorm"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/storage"
	"hei-gin/internal/modules/shared"
)

// Service é—¨æˆ·ç«¯ç”¨æˆ·ä¸­å¿ƒæœåŠ¡ã€‚
//
// Author: Charlie
type Service struct {
	repo    *Repo
	storage *storage.Manager
}

// NewService æž„é€ æœåŠ¡ã€‚
func NewService(db *gorm.DB, storage *storage.Manager) *Service {
	return &Service{repo: NewRepo(db), storage: storage}
}

// New æž„å»º user.portal æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB, d.Storage)
	return module.Module{
		Name:   "user.portal",
		Order:  31,
		Models: []any{&Profile{}},
		Routes: []module.RouteRegistrar{s.registerRoutes},
	}
}

// Me å½“å‰ç”¨æˆ·æ¦‚è§ˆã€‚
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

// UpdateProfile æ›´æ–°èµ„æ–™ã€‚
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
	return s.repo.UpdateProfile(ctx, accountID, updates)
}

// UploadAvatar ä¸Šä¼ å¤´åƒå¹¶æ›´æ–°èµ„æ–™ã€‚
func (s *Service) UploadAvatar(ctx context.Context, accountID string, filename string, r io.Reader, size int64, contentType string) (string, error) {
	ext := path.Ext(filename)
	object := storage.ObjectKey("avatar/portal", idgen.Next()+ext)
	url, err := s.storage.Provider().Put(ctx, object, r, size, contentType)
	if err != nil {
		return "", err
	}
	_ = s.getOrCreate(ctx, accountID)
	_ = s.repo.UpdateProfile(ctx, accountID, map[string]any{"avatar": url, "updated_by": accountID})
	return url, nil
}

// UpdatePassword æ›´æ–°å¯†ç ã€‚
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

// UpdatePhone æ›´æ–°æ‰‹æœºå·ã€‚
func (s *Service) UpdatePhone(ctx context.Context, accountID string, req PhoneUpdateParam) error {
	if err := s.verifyPassword(ctx, accountID, req.Password); err != nil {
		return err
	}
	_ = s.getOrCreate(ctx, accountID)
	return s.repo.UpdateProfile(ctx, accountID, map[string]any{"phone": req.Phone, "updated_by": accountID})
}

// UpdateEmail æ›´æ–°é‚®ç®±ã€‚
func (s *Service) UpdateEmail(ctx context.Context, accountID string, req EmailUpdateParam) error {
	if err := s.verifyPassword(ctx, accountID, req.Password); err != nil {
		return err
	}
	_ = s.getOrCreate(ctx, accountID)
	return s.repo.UpdateProfile(ctx, accountID, map[string]any{
		"email": strings.TrimSpace(req.Email), "updated_by": accountID,
	})
}

func (s *Service) verifyPassword(ctx context.Context, accountID, password string) error {
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

// SessionFromContext ä»Ž context å–ä¼šè¯ã€‚
func SessionFromContext(ctx context.Context) *security.SessionPayload {
	return contextx.Session(ctx)
}

// AccountIDFromContext ä»Ž context å–è´¦å· IDã€‚
func AccountIDFromContext(ctx context.Context) string {
	return contextx.AccountID(ctx)
}
