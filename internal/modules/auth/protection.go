// internal/modules/auth/protection.go 登录保护。
//
// Author: Charlie

package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"hei-gin/internal/framework/core/security"
)

const (
	failAcctPrefix = "auth:fail:acct:"
	failIPPrefix   = "auth:fail:ip:"
	lockAcctPrefix = "auth:lock:acct:"
	lockIPPrefix   = "auth:lock:ip:"
	otpLoginPrefix = "auth:otp:login:"
	resetTokPrefix     = "password:reset:"
	resetPwdOtpPrefix  = "auth:otp:reset-password:"
)

// EnsureLoginAllowed 登录前检查账号/IP 是否锁定。
func (r *Repo) EnsureLoginAllowed(ctx context.Context, accountType security.AccountType, account, clientIP string) error {
	acct := normalizeAccount(account)
	typeName := string(accountType)
	if ok, _ := r.rdb.Exists(ctx, lockAcctPrefix+typeName+":"+acct).Result(); ok > 0 {
		return errAccountLocked
	}
	if clientIP != "" {
		if ok, _ := r.rdb.Exists(ctx, lockIPPrefix+clientIP).Result(); ok > 0 {
			return errIPLocked
		}
	}
	return nil
}

// RecordLoginFailure 记录失败并在超限时锁定。
func (r *Repo) RecordLoginFailure(ctx context.Context, cfg loginProtectCfg, accountType security.AccountType, account, clientIP string) {
	typeName := string(accountType)
	acct := normalizeAccount(account)
	r.bumpFailure(ctx,
		failAcctPrefix+typeName+":"+acct,
		lockAcctPrefix+typeName+":"+acct,
		cfg.AccountMax, cfg.WindowSeconds, cfg.LockSeconds,
	)
	if clientIP != "" {
		r.bumpFailure(ctx,
			failIPPrefix+clientIP,
			lockIPPrefix+clientIP,
			cfg.IPMax, cfg.WindowSeconds, cfg.LockSeconds,
		)
	}
}

// ClearLoginFailures 登录成功后清除失败计数。
func (r *Repo) ClearLoginFailures(ctx context.Context, accountType security.AccountType, account, clientIP string) {
	typeName := string(accountType)
	acct := normalizeAccount(account)
	_ = r.rdb.Del(ctx, failAcctPrefix+typeName+":"+acct).Err()
	if clientIP != "" {
		_ = r.rdb.Del(ctx, failIPPrefix+clientIP).Err()
	}
}

func (r *Repo) bumpFailure(ctx context.Context, failKey, lockKey string, maxFailures, windowSeconds, lockSeconds int) {
	if maxFailures <= 0 {
		return
	}
	n, err := r.rdb.Incr(ctx, failKey).Result()
	if err != nil {
		return
	}
	if n == 1 && windowSeconds > 0 {
		_ = r.rdb.Expire(ctx, failKey, time.Duration(windowSeconds)*time.Second).Err()
	}
	if int(n) >= maxFailures {
		ttl := time.Duration(lockSeconds) * time.Second
		if ttl <= 0 {
			ttl = 15 * time.Minute
		}
		_ = r.rdb.Set(ctx, lockKey, "1", ttl).Err()
		_ = r.rdb.Del(ctx, failKey).Err()
	}
}

type loginProtectCfg struct {
	WindowSeconds int
	AccountMax    int
	IPMax         int
	LockSeconds   int
}

// StoreLoginOTP 缓存登录 OTP（5 分钟）。
func (r *Repo) StoreLoginOTP(ctx context.Context, accountType, channel, target, code string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	key := otpLoginPrefix + accountType + ":" + channel + ":" + target
	return r.rdb.Set(ctx, key, code, ttl).Err()
}

// ConsumeLoginOTP 校验并消费登录 OTP。
func (r *Repo) ConsumeLoginOTP(ctx context.Context, accountType, channel, target, code string) bool {
	key := otpLoginPrefix + accountType + ":" + channel + ":" + target
	stored, err := r.rdb.Get(ctx, key).Result()
	_ = r.rdb.Del(ctx, key)
	if err == redis.Nil || err != nil || stored == "" {
		return false
	}
	return stored == strings.TrimSpace(code)
}

// StoreResetToken 缓存密码重置令牌。
func (r *Repo) StoreResetToken(ctx context.Context, token, accountID string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return r.rdb.Set(ctx, resetTokPrefix+token, accountID, ttl).Err()
}

// StoreResetPasswordOTP 缓存手机重置密码 OTP。
func (r *Repo) StoreResetPasswordOTP(ctx context.Context, accountType, phone, code string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	key := resetPwdOtpPrefix + accountType + ":PHONE:" + phone
	return r.rdb.Set(ctx, key, code, ttl).Err()
}

// ConsumeResetPasswordOTP 校验并消费手机重置密码 OTP。
func (r *Repo) ConsumeResetPasswordOTP(ctx context.Context, accountType, phone, code string) bool {
	key := resetPwdOtpPrefix + accountType + ":PHONE:" + phone
	stored, err := r.rdb.Get(ctx, key).Result()
	_ = r.rdb.Del(ctx, key)
	if err == redis.Nil || err != nil || stored == "" {
		return false
	}
	return stored == strings.TrimSpace(code)
}

// ConsumeResetToken 消费重置令牌，返回账号 ID。
func (r *Repo) ConsumeResetToken(ctx context.Context, token string) (string, error) {
	key := resetTokPrefix + token
	id, err := r.rdb.Get(ctx, key).Result()
	_ = r.rdb.Del(ctx, key)
	if err == redis.Nil || id == "" {
		return "", fmt.Errorf("重置令牌无效或已过期")
	}
	if err != nil {
		return "", err
	}
	return id, nil
}

func sixDigitCode() (string, error) {
	var n uint32
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	n = uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return fmt.Sprintf("%06d", n%1000000), nil
}

func normalizeAccount(account string) string {
	return strings.ToLower(strings.TrimSpace(account))
}

func normalizePhone(phone string) string {
	return strings.TrimSpace(phone)
}

func newResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	const hexdigits = "0123456789abcdef"
	out := make([]byte, 64)
	for i, v := range b {
		out[i*2] = hexdigits[v>>4]
		out[i*2+1] = hexdigits[v&0x0f]
	}
	return string(out), nil
}
