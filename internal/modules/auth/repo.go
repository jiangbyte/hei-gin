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

// Repo 认证 Redis 持久化（验证码、密码密钥）。
//
// Author: Charlie
type Repo struct{ rdb *redis.Client }

// NewRepo 构造 Repo。
func NewRepo(rdb *redis.Client) *Repo { return &Repo{rdb: rdb} }

func newID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CreateCaptcha 生成验证码并存入 Redis。
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

// VerifyCaptcha 校验验证码。
func (r *Repo) VerifyCaptcha(ctx context.Context, captchaID, value string) error {
	key := captchaKeyPrefix + captchaID
	stored, err := r.rdb.Get(ctx, key).Result()
	_ = r.rdb.Del(ctx, key)
	if err == redis.Nil || stored == "" {
		return fmt.Errorf("无效或过期的验证码")
	}
	if err != nil {
		return err
	}
	if !security.CheckPassword(stored, strings.ToLower(strings.TrimSpace(value))) {
		return fmt.Errorf("无效或过期的验证码")
	}
	return nil
}

// CreatePasswordKey 生成 RSA 密钥对，私钥存 Redis。
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

// DecryptPassword 用 Redis 私钥解密密码。
func (r *Repo) DecryptPassword(ctx context.Context, keyID, encryptedValue string) (string, error) {
	if keyID == "" {
		return "", fmt.Errorf("缺少 password_key_id")
	}
	key := passwordKeyPrefix + keyID
	raw, err := r.rdb.Get(ctx, key).Result()
	_ = r.rdb.Del(ctx, key)
	if err == redis.Nil || raw == "" {
		return "", fmt.Errorf("无效或过期的密码加密密钥")
	}
	if err != nil {
		return "", err
	}
	if encryptedValue == "" {
		return "", nil
	}
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return "", fmt.Errorf("无效或过期的密码加密密钥")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("无效或过期的密码加密密钥")
	}
	rsaKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("无效或过期的密码加密密钥")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("无效的加密密码")
	}
	plain, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaKey, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("无效的加密密码")
	}
	return string(plain), nil
}
