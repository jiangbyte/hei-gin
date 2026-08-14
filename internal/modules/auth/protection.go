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
	resetTokPrefix = "password:reset:"
)

// EnsureLoginAllowed ç™»å½•å‰æ£€æŸ¥è´¦å·/IP æ˜¯å¦é”å®šã€‚
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

// RecordLoginFailure è®°å½•å¤±è´¥å¹¶åœ¨è¶…é™æ—¶é”å®šã€‚
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

// ClearLoginFailures ç™»å½•æˆåŠŸåŽæ¸…é™¤å¤±è´¥è®¡æ•°ã€‚
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

// StoreLoginOTP ç¼“å­˜ç™»å½• OTPï¼ˆ5 åˆ†é’Ÿï¼‰ã€‚
func (r *Repo) StoreLoginOTP(ctx context.Context, accountType, channel, target, code string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	key := otpLoginPrefix + accountType + ":" + channel + ":" + target
	return r.rdb.Set(ctx, key, code, ttl).Err()
}

// ConsumeLoginOTP æ ¡éªŒå¹¶æ¶ˆè´¹ç™»å½• OTPã€‚
func (r *Repo) ConsumeLoginOTP(ctx context.Context, accountType, channel, target, code string) bool {
	key := otpLoginPrefix + accountType + ":" + channel + ":" + target
	stored, err := r.rdb.Get(ctx, key).Result()
	_ = r.rdb.Del(ctx, key)
	if err == redis.Nil || err != nil || stored == "" {
		return false
	}
	return stored == strings.TrimSpace(code)
}

// StoreResetToken ç¼“å­˜å¯†ç é‡ç½®ä»¤ç‰Œã€‚
func (r *Repo) StoreResetToken(ctx context.Context, token, accountID string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return r.rdb.Set(ctx, resetTokPrefix+token, accountID, ttl).Err()
}

// ConsumeResetToken æ¶ˆè´¹é‡ç½®ä»¤ç‰Œï¼Œè¿”å›žè´¦å· IDã€‚
func (r *Repo) ConsumeResetToken(ctx context.Context, token string) (string, error) {
	key := resetTokPrefix + token
	id, err := r.rdb.Get(ctx, key).Result()
	_ = r.rdb.Del(ctx, key)
	if err == redis.Nil || id == "" {
		return "", fmt.Errorf("é‡ç½®ä»¤ç‰Œæ— æ•ˆæˆ–å·²è¿‡æœŸ")
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
