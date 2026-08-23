// internal/modules/auth/register.go 门户注册三通道（对齐 hei-boot registerPortal）。
//
// Author: Charlie

package auth

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var nonAccountCharRe = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

func (s *Service) consumeRegisterCode(ctx context.Context, channel, target, code string) error {
	if strings.TrimSpace(code) == "" {
		return fmt.Errorf("验证码不能为空")
	}
	ch := strings.ToUpper(strings.TrimSpace(channel))
	normTarget := normalizeRegisterTarget(ch, target)
	if !s.repo.ConsumeRegisterOTP(ctx, ch, normTarget, code) {
		return fmt.Errorf("验证码无效或已过期")
	}
	return nil
}

func normalizeRegisterTarget(channel, target string) string {
	t := strings.TrimSpace(target)
	if channel == "EMAIL" {
		return strings.ToLower(t)
	}
	return t
}

func (s *Service) identityExists(ctx context.Context, identityType, identifier string) bool {
	if s.accounts == nil {
		return false
	}
	type checker interface {
		IdentityExists(ctx context.Context, identityType, identifier string) bool
	}
	c, ok := s.accounts.(checker)
	if !ok {
		return false
	}
	return c.IdentityExists(ctx, identityType, identifier)
}

func (s *Service) allocateAccountLogin(ctx context.Context, base string) (string, error) {
	if s.accounts == nil {
		return "", errAccountFinder
	}
	type allocator interface {
		AllocateUniqueAccount(ctx context.Context, base string) (string, error)
	}
	a, ok := s.accounts.(allocator)
	if !ok {
		return "", errPortalRegistrar
	}
	return a.AllocateUniqueAccount(ctx, base)
}

func allocateBaseFromEmail(email string) string {
	local := email
	if i := strings.Index(email, "@"); i > 0 {
		local = email[:i]
	}
	return sanitizeAccountBase(local)
}

func allocateBaseFromPhone(phone string) string {
	digits := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, phone)
	if len(digits) > 8 {
		digits = digits[len(digits)-8:]
	}
	return "u" + digits
}

func sanitizeAccountBase(base string) string {
	base = strings.ToLower(strings.TrimSpace(base))
	base = nonAccountCharRe.ReplaceAllString(base, "_")
	base = strings.Trim(base, "_")
	if base == "" {
		base = "user"
	}
	if len(base) > 48 {
		base = base[:48]
	}
	return base
}

func registerPortalInput(accountLogin, passwordHash string, email, phone, nickname *string) PortalRegisterInput {
	in := PortalRegisterInput{
		AccountLogin: accountLogin,
		PasswordHash: passwordHash,
		Nickname:     nickname,
		Email:        email,
		Phone:        phone,
	}
	if email != nil && *email != "" {
		in.EmailEnabled = true
		in.EmailVerified = true
	}
	if phone != nil && *phone != "" {
		in.PhoneEnabled = true
		in.PhoneVerified = true
	}
	return in
}
