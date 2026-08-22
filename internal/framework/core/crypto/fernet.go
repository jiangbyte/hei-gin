// Package crypto 提供配置值 Fernet 兼容加解密（github.com/fernet/fernet-go）。
//
// 若 Config.Crypto.VaultAddr 已设置，仍仅从环境变量 / 配置读取 FernetKey
// （HEI_CRYPTO_FERNET_KEY 或 crypto.fernet_key），不实现完整 Vault 客户端。
//
// Author: Charlie
package crypto

import (
	"errors"
	"fmt"
	"os"
	"strings"

	fernet "github.com/fernet/fernet-go"
)

// Codec Fernet 兼容编解码器（Python cryptography.fernet.Fernet）。
//
// Author: Charlie
type Codec struct {
	key *fernet.Key
}

// NewFernet 从 url-safe base64 编码的 32 字节密钥创建 Codec。
func NewFernet(key string) (*Codec, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, errors.New("crypto: empty fernet key")
	}
	k, err := fernet.DecodeKey(key)
	if err != nil {
		return nil, fmt.Errorf("crypto: decode fernet key: %w", err)
	}
	return &Codec{key: k}, nil
}

// NewFernetFromConfig 优先用 cfgKey；若 VaultAddr 非空且 cfgKey 为空，尝试环境变量 HEI_CRYPTO_FERNET_KEY。
func NewFernetFromConfig(fernetKey, vaultAddr string) (*Codec, error) {
	key := strings.TrimSpace(fernetKey)
	if key == "" && strings.TrimSpace(vaultAddr) != "" {
		key = strings.TrimSpace(os.Getenv("HEI_CRYPTO_FERNET_KEY"))
	}
	if key == "" {
		return nil, nil
	}
	return NewFernet(key)
}

// Encrypt 加密明文，返回 url-safe base64 token。
func (c *Codec) Encrypt(plaintext string) (string, error) {
	if c == nil || c.key == nil {
		return plaintext, nil
	}
	token, err := fernet.EncryptAndSign([]byte(plaintext), c.key)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

// Decrypt 解密 Fernet token；失败返回错误。
func (c *Codec) Decrypt(token string) (string, error) {
	if c == nil || c.key == nil {
		return token, nil
	}
	plain := fernet.VerifyAndDecrypt([]byte(token), 0, []*fernet.Key{c.key})
	if plain == nil {
		return "", errors.New("crypto: invalid or expired fernet token")
	}
	return string(plain), nil
}

// LooksEncrypted 启发式判断字符串是否像 Fernet token。
func LooksEncrypted(v string) bool {
	v = strings.TrimSpace(v)
	if len(v) < 60 {
		return false
	}
	return strings.HasPrefix(v, "gAAAAA")
}
