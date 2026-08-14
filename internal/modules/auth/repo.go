// internal/modules/auth/repo.go 持久化仓储。
//
// Author: Charlie

package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"hei-gin/internal/framework/core/security"
)

// Repo è®¤è¯ Redis æŒä¹…åŒ–ï¼ˆéªŒè¯ç ã€å¯†ç å¯†é’¥ï¼‰ã€‚
//
// Author: Charlie
type Repo struct{ rdb *redis.Client }

// NewRepo æž„é€  Repoã€‚
func NewRepo(rdb *redis.Client) *Repo { return &Repo{rdb: rdb} }

func newID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CreateCaptcha ç”ŸæˆéªŒè¯ç å¹¶å­˜å…¥ Redisã€‚
func (r *Repo) CreateCaptcha(ctx context.Context, ttl time.Duration) (*CaptchaResult, error) {
	var b strings.Builder
	for i := 0; i < 4; i++ {
		idx := make([]byte, 1)
		if _, err := rand.Read(idx); err != nil {
			return nil, err
		}
		b.WriteByte(captchaAlphabet[int(idx[0])%len(captchaAlphabet)])
	}
	value := b.String()
	id, err := newID()
	if err != nil {
		return nil, err
	}
	hash, err := security.HashPassword(strings.ToLower(value))
	if err != nil {
		return nil, err
	}
	if err := r.rdb.Set(ctx, captchaKeyPrefix+id, hash, ttl).Err(); err != nil {
		return nil, err
	}
	return &CaptchaResult{
		CaptchaID:   id,
		ImageBase64: captchaSVGBase64(value),
		ImageType:   "image/svg+xml",
	}, nil
}

// VerifyCaptcha æ ¡éªŒéªŒè¯ç ã€‚
func (r *Repo) VerifyCaptcha(ctx context.Context, captchaID, value string) error {
	key := captchaKeyPrefix + captchaID
	stored, err := r.rdb.Get(ctx, key).Result()
	_ = r.rdb.Del(ctx, key)
	if err == redis.Nil || stored == "" {
		return fmt.Errorf("æ— æ•ˆæˆ–è¿‡æœŸçš„éªŒè¯ç ")
	}
	if err != nil {
		return err
	}
	if !security.CheckPassword(stored, strings.ToLower(strings.TrimSpace(value))) {
		return fmt.Errorf("æ— æ•ˆæˆ–è¿‡æœŸçš„éªŒè¯ç ")
	}
	return nil
}

// CreatePasswordKey ç”Ÿæˆ RSA å¯†é’¥å¯¹ï¼Œç§é’¥å­˜ Redisã€‚
func (r *Repo) CreatePasswordKey(ctx context.Context, ttl time.Duration) (*PasswordKeyResult, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	der, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	pubDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	id, err := newID()
	if err != nil {
		return nil, err
	}
	if err := r.rdb.Set(ctx, passwordKeyPrefix+id, string(privPEM), ttl).Err(); err != nil {
		return nil, err
	}
	return &PasswordKeyResult{
		KeyID:     id,
		PublicKey: base64.StdEncoding.EncodeToString(pubDER),
	}, nil
}

// DecryptPassword ç”¨ Redis ç§é’¥è§£å¯†å¯†ç ã€‚
func (r *Repo) DecryptPassword(ctx context.Context, keyID, encryptedValue string) (string, error) {
	if keyID == "" {
		return "", fmt.Errorf("ç¼ºå°‘ password_key_id")
	}
	key := passwordKeyPrefix + keyID
	raw, err := r.rdb.Get(ctx, key).Result()
	_ = r.rdb.Del(ctx, key)
	if err == redis.Nil || raw == "" {
		return "", fmt.Errorf("æ— æ•ˆæˆ–è¿‡æœŸçš„å¯†ç åŠ å¯†å¯†é’¥")
	}
	if err != nil {
		return "", err
	}
	if encryptedValue == "" {
		return "", nil
	}
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return "", fmt.Errorf("æ— æ•ˆæˆ–è¿‡æœŸçš„å¯†ç åŠ å¯†å¯†é’¥")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("æ— æ•ˆæˆ–è¿‡æœŸçš„å¯†ç åŠ å¯†å¯†é’¥")
	}
	rsaKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("æ— æ•ˆæˆ–è¿‡æœŸçš„å¯†ç åŠ å¯†å¯†é’¥")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("æ— æ•ˆçš„åŠ å¯†å¯†ç ")
	}
	plain, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaKey, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("æ— æ•ˆçš„åŠ å¯†å¯†ç ")
	}
	return string(plain), nil
}
