// internal/modules/profile/service.go 用户中心共享服务（admin/portal 一套实现，按表区分）。
//
// Author: Charlie

package profile

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/notify"
	"hei-gin/internal/framework/platform/runtimecfg"
	"hei-gin/internal/framework/platform/storage"
	"hei-gin/internal/modules/shared"
)

// otpBindPrefix 绑定验证码 Redis 前缀。
const otpBindPrefix = "profile:otp:bind:"

// Service 用户中心服务。
//
// Author: Charlie
type Service struct {
	repo           *Repo
	rdb            *redis.Client
	notify         *notify.Facade
	storage        *storage.Manager
	runtime        *runtimecfg.Settings
	passwordPolicy *shared.PasswordPolicy
	avatarPrefix   string
	accountType    security.AccountType
}

// NewService 构造按账户类型绑定的用户中心服务。
func NewService(db *gorm.DB, rdb *redis.Client, nf *notify.Facade, storage *storage.Manager,
	rt *runtimecfg.Settings, accountType security.AccountType, table, avatarPrefix string) *Service {
	return &Service{
		repo:           NewRepo(db, table),
		rdb:            rdb,
		notify:         nf,
		storage:        storage,
		runtime:        rt,
		passwordPolicy: shared.NewPasswordPolicy(db, rt),
		avatarPrefix:   avatarPrefix,
		accountType:    accountType,
	}
}

// Me 当前用户概览。
func (s *Service) Me(ctx context.Context, sess *security.SessionPayload) (*MeResult, error) {
	if sess == nil {
		return nil, errUnauthorized
	}
	profile := s.getOrCreate(ctx, sess.AccountID)
	account, _ := s.repo.GetAccountIdentifier(ctx, sess.AccountID)
	return &MeResult{
		AccountID:       sess.AccountID,
		AccountType:     sess.AccountType,
		Account:         account,
		Name:            profile.Name,
		Nickname:        profile.Nickname,
		Avatar:          profile.Avatar,
		RoleIDs:         sess.RoleIDs,
		DeptIDs:         sess.DeptIDs,
		GroupIDs:        sess.GroupIDs,
		RoleIDNames:     s.repo.LoadIDNames(ctx, "sys_role", sess.RoleIDs),
		DeptIDNames:     s.repo.LoadIDNames(ctx, "sys_dept", sess.DeptIDs),
		GroupIDNames:    s.repo.LoadIDNames(ctx, "sys_group", sess.GroupIDs),
		PermissionKeys:  sess.PermissionKeys,
		Profile:         profile,
		PasswordExpired: sess.PasswordExpired,
		ForceBindEmail:  s.forceBind(ctx, "EMAIL"),
		ForceBindPhone:  s.forceBind(ctx, "PHONE"),
	}, nil
}

// UpdateProfile 更新资料（remark 仅管理端表存在）。
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
	if req.Remark != nil && s.repo.Table() == ProfileTableAdmin {
		updates["remark"] = *req.Remark
	}
	return s.repo.UpdateProfile(ctx, accountID, updates)
}

// UploadAvatar 上传头像并更新资料。
func (s *Service) UploadAvatar(ctx context.Context, accountID string, filename string, r io.Reader, size int64, contentType string) (string, error) {
	ext := path.Ext(filename)
	object := storage.ObjectKey("avatar/"+s.avatarPrefix, idgen.Next()+ext)
	url, err := s.storage.Provider().Put(ctx, object, r, size, contentType)
	if err != nil {
		return "", err
	}
	_ = s.getOrCreate(ctx, accountID)
	_ = s.repo.UpdateProfile(ctx, accountID, map[string]any{"avatar": url, "updated_by": accountID})
	return url, nil
}

// UpdatePassword 更新密码（旧密码或 OTP 校验 + 密码策略校验 + 历史记录）。
func (s *Service) UpdatePassword(ctx context.Context, accountID string, req PasswordUpdateParam) error {
	if err := s.verify(ctx, accountID, req.OldPassword, req.OTPCode, "PASSWORD_CHANGE"); err != nil {
		return err
	}
	accountName, _ := s.repo.GetAccountIdentifier(ctx, accountID)
	email, phone := "", ""
	if p, err := s.repo.GetProfile(ctx, accountID); err == nil {
		if p.Email != nil {
			email = *p.Email
		}
		if p.Phone != nil {
			phone = *p.Phone
		}
	}
	if err := s.passwordPolicy.Validate(ctx, req.NewPassword, accountID, accountName, email, phone); err != nil {
		return err
	}
	hash, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateAccountPassword(ctx, accountID, hash); err != nil {
		return err
	}
	return s.passwordPolicy.RecordHistory(ctx, accountID, req.NewPassword, accountID, "self_update")
}

// UpdatePhone 更新手机号（密码/OTP 校验 + 维护 PHONE 登录身份）。
func (s *Service) UpdatePhone(ctx context.Context, accountID string, req PhoneUpdateParam) error {
	if err := s.verify(ctx, accountID, req.Password, req.OTPCode, "PHONE"); err != nil {
		return err
	}
	_ = s.getOrCreate(ctx, accountID)
	updates := map[string]any{"updated_by": accountID}
	if req.Phone != nil {
		updates["phone"] = strings.TrimSpace(*req.Phone)
	}
	if err := s.repo.UpdateProfile(ctx, accountID, updates); err != nil {
		return err
	}
	enabled := req.PhoneLoginEnabled != nil && *req.PhoneLoginEnabled
	return s.syncIdentity(ctx, accountID, "PHONE", req.Phone, enabled)
}

// UpdateEmail 更新邮箱（密码/OTP 校验 + 维护 EMAIL 登录身份）。
func (s *Service) UpdateEmail(ctx context.Context, accountID string, req EmailUpdateParam) error {
	if err := s.verify(ctx, accountID, req.Password, req.OTPCode, "EMAIL"); err != nil {
		return err
	}
	_ = s.getOrCreate(ctx, accountID)
	updates := map[string]any{"updated_by": accountID}
	if req.Email != nil {
		updates["email"] = strings.TrimSpace(*req.Email)
	}
	if err := s.repo.UpdateProfile(ctx, accountID, updates); err != nil {
		return err
	}
	enabled := req.EmailLoginEnabled != nil && *req.EmailLoginEnabled
	return s.syncIdentity(ctx, accountID, "EMAIL", req.Email, enabled)
}

// OrgInfo 组织关联信息。
func (s *Service) OrgInfo(sess *security.SessionPayload) OrgInfoResult {
	return OrgInfoResult{
		RoleIDs:  sess.RoleIDs,
		DeptIDs:  sess.DeptIDs,
		GroupIDs: sess.GroupIDs,
	}
}

// SendBindCode 发送绑定/改密验证码（通道未配置时仅写入 OTP 缓存并返回成功）。
// scene 区分绑定（BIND_PHONE_CODE/BIND_EMAIL_CODE）与改密（CHANGE_PASSWORD_CODE）。
func (s *Service) SendBindCode(ctx context.Context, accountID, scene, channel, target string) error {
	code, err := sixDigitCode()
	if err != nil {
		return err
	}
	ttl := 5 * time.Minute
	if err := s.storeOTP(ctx, channel, target, code, ttl); err != nil {
		return err
	}
	if s.notify != nil {
		vars := map[string]any{"app_name": "HEI", "code": code, "expire_minutes": 5}
		_ = s.notify.SendTemplated(ctx, scene, target, vars)
	}
	_ = accountID
	return nil
}

// SendPasswordChangeCode 发送改密验证码到已绑定手机/邮箱。
func (s *Service) SendPasswordChangeCode(ctx context.Context, accountID string) error {
	p, err := s.repo.GetProfile(ctx, accountID)
	if err != nil {
		return err
	}
	if p.Phone != nil && *p.Phone != "" {
		return s.SendBindCode(ctx, accountID, "CHANGE_PASSWORD_CODE", "PHONE", *p.Phone)
	}
	if p.Email != nil && *p.Email != "" {
		return s.SendBindCode(ctx, accountID, "CHANGE_PASSWORD_CODE", "EMAIL", *p.Email)
	}
	return nil
}

func (s *Service) verify(ctx context.Context, accountID, password, otpCode, channel string) error {
	if strings.TrimSpace(otpCode) != "" {
		target := ""
		p, err := s.repo.GetProfile(ctx, accountID)
		if err == nil {
			switch channel {
			case "PHONE":
				if p.Phone != nil {
					target = *p.Phone
				}
			case "EMAIL":
				if p.Email != nil {
					target = *p.Email
				}
			}
		}
		if target == "" {
			return errOTPInvalid
		}
		if !s.consumeOTP(ctx, channel, target, otpCode) {
			return errOTPInvalid
		}
		return nil
	}
	if password == "" {
		return errPasswordRequired
	}
	acc, err := s.repo.GetAccountPassword(ctx, accountID)
	if err != nil {
		return err
	}
	if !security.CheckPassword(acc.PasswordHash, password) {
		return errPassword
	}
	return nil
}

// syncIdentity 维护 sys_account_identity 登录身份（PHONE/EMAIL），enabled=false 时删除。
func (s *Service) syncIdentity(ctx context.Context, accountID, identityType string, value *string, enabled bool) error {
	db := s.repo.db.WithContext(ctx)
	if !enabled || value == nil || strings.TrimSpace(*value) == "" {
		// 关闭登录通道或置空：删除该类型身份
		return db.Table("sys_account_identity").
			Where("account_id = ? AND identity_type = ?", accountID, identityType).
			Delete(nil).Error
	}
	identifier := strings.ToLower(strings.TrimSpace(*value))
	var n int64
	if err := db.Table("sys_account_identity").
		Where("account_id = ? AND identity_type = ?", accountID, identityType).
		Count(&n).Error; err != nil {
		return err
	}
	if n == 0 {
		row := map[string]any{
			"id": idgen.Next(), "account_id": accountID, "identity_type": identityType,
			"identifier": identifier, "verified": true, "is_primary": false,
			"bind_status": "BOUND", "created_by": accountID, "updated_by": accountID,
		}
		return db.Table("sys_account_identity").Create(row).Error
	}
	return db.Table("sys_account_identity").
		Where("account_id = ? AND identity_type = ?", accountID, identityType).
		Update("identifier", identifier).Error
}

func (s *Service) storeOTP(ctx context.Context, channel, target, code string, ttl time.Duration) error {
	if s.rdb == nil {
		return nil
	}
	key := otpBindPrefix + string(s.accountType) + ":" + channel + ":" + strings.ToLower(strings.TrimSpace(target))
	return s.rdb.Set(ctx, key, code, ttl).Err()
}

func (s *Service) consumeOTP(ctx context.Context, channel, target, code string) bool {
	if s.rdb == nil {
		return false
	}
	key := otpBindPrefix + string(s.accountType) + ":" + channel + ":" + strings.ToLower(strings.TrimSpace(target))
	stored, err := s.rdb.Get(ctx, key).Result()
	_ = s.rdb.Del(ctx, key)
	if err == redis.Nil || err != nil || stored == "" {
		return false
	}
	return stored == strings.TrimSpace(code)
}

// forceBind 读取 AUTH_FORCE_BIND_{TYPE}_{CHANNEL} 运行时配置。
func (s *Service) forceBind(ctx context.Context, channel string) bool {
	key := "AUTH_FORCE_BIND_" + strings.ToUpper(string(s.accountType)) + "_" + channel
	if s.runtime != nil {
		return s.runtime.GetBool(ctx, key, false)
	}
	return false
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

func sixDigitCode() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	n := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return fmt.Sprintf("%06d", n%1000000), nil
}

var (
	errUnauthorized     = &passErr{msg: "unauthorized"}
	errPassword         = &passErr{msg: "password incorrect"}
	errPasswordRequired = &passErr{msg: "password or otp_code required"}
	errOTPInvalid       = &passErr{msg: "otp code invalid or expired"}
)

type passErr struct{ msg string }

func (e *passErr) Error() string { return e.msg }

// SessionFromContext 从 context 取会话。
func SessionFromContext(ctx context.Context) *security.SessionPayload {
	return contextx.Session(ctx)
}

// AccountIDFromContext 从 context 取账号 ID。
func AccountIDFromContext(ctx context.Context) string {
	return contextx.AccountID(ctx)
}
