// internal/modules/profile/identity/crypto.go 实名敏感字段加解密、哈希与脱敏。
//
// Author: Charlie
package identity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"

	"hei-gin/internal/framework/core/crypto"
	"hei-gin/internal/framework/platform/runtimecfg"
)

// FieldCrypto 实名敏感字段工具。
type FieldCrypto struct {
	codec *crypto.Codec
}

// NewFieldCrypto 构造实名加解密器（优先 PROFILE_IDENTITY_CRYPTO_KEY，回退全局 Fernet）。
func NewFieldCrypto(rt *runtimecfg.Settings, fallbackKey string) (*FieldCrypto, error) {
	key := ""
	if rt != nil {
		key = strings.TrimSpace(rt.GetString(context.Background(), "PROFILE_IDENTITY_CRYPTO_KEY", ""))
	}
	if key == "" {
		key = strings.TrimSpace(fallbackKey)
	}
	if key == "" {
		return nil, fmt.Errorf("identity crypto is not configured; set PROFILE_IDENTITY_CRYPTO_KEY")
	}
	codec, err := crypto.NewFernet(key)
	if err != nil {
		return nil, err
	}
	return &FieldCrypto{codec: codec}, nil
}

// Encrypt 加密明文。
func (c *FieldCrypto) Encrypt(plaintext string) (string, error) {
	if c == nil || c.codec == nil {
		return "", fmt.Errorf("identity crypto is not configured; set PROFILE_IDENTITY_CRYPTO_KEY")
	}
	plain := strings.TrimSpace(plaintext)
	if plain == "" {
		return "", nil
	}
	return c.codec.Encrypt(plain)
}

// Decrypt 解密密文。
func (c *FieldCrypto) Decrypt(ciphertext string) (string, error) {
	if c == nil || c.codec == nil {
		return "", fmt.Errorf("identity crypto is not configured; set PROFILE_IDENTITY_CRYPTO_KEY")
	}
	cipher := strings.TrimSpace(ciphertext)
	if cipher == "" {
		return "", nil
	}
	return c.codec.Decrypt(cipher)
}

// HashDocumentNo SHA256(upper(type)|upper(no)) hex。
func (c *FieldCrypto) HashDocumentNo(documentType, documentNo string) string {
	_ = c
	if strings.TrimSpace(documentNo) == "" {
		return ""
	}
	normalizedType := strings.ToUpper(strings.TrimSpace(documentType))
	normalizedNo := strings.ToUpper(strings.TrimSpace(documentNo))
	sum := sha256.Sum256([]byte(normalizedType + "|" + normalizedNo))
	return hex.EncodeToString(sum[:])
}

// MaskRealName 保留首字，其余 *。
func (c *FieldCrypto) MaskRealName(realName string) *string {
	_ = c
	trimmed := strings.TrimSpace(realName)
	if trimmed == "" {
		return nil
	}
	runes := []rune(trimmed)
	if len(runes) <= 1 {
		out := "*"
		return &out
	}
	out := string(runes[0]) + strings.Repeat("*", len(runes)-1)
	return &out
}

// MaskDocumentNo 保留前 3 + 后 4，中间 *。
func (c *FieldCrypto) MaskDocumentNo(documentNo string) *string {
	_ = c
	trimmed := strings.TrimSpace(documentNo)
	if trimmed == "" {
		return nil
	}
	n := utf8.RuneCountInString(trimmed)
	if n <= 7 {
		out := strings.Repeat("*", n)
		return &out
	}
	runes := []rune(trimmed)
	keepPrefix := 3
	if keepPrefix > len(runes) {
		keepPrefix = len(runes)
	}
	keepSuffix := 4
	if keepSuffix > len(runes)-keepPrefix {
		keepSuffix = len(runes) - keepPrefix
	}
	out := string(runes[:keepPrefix]) +
		strings.Repeat("*", len(runes)-keepPrefix-keepSuffix) +
		string(runes[len(runes)-keepSuffix:])
	return &out
}
