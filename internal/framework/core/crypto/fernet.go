// Package crypto 提供配置值 Fernet 兼容加解密。
//
// 若 Config.Crypto.VaultAddr 已设置，仍仅从环境变量 / 配置读取 FernetKey
// （HEI_CRYPTO_FERNET_KEY 或 crypto.fernet_key），不实现完整 Vault 客户端。
//
// Author: Charlie
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	fernetVersion byte = 0x80
	fernetTTL          = 0 // 0 = 不校验过期
)

// Codec Fernet 兼容编解码器（Python cryptography.fernet.Fernet）。
//
// Author: Charlie
type Codec struct {
	signingKey []byte
	encryptKey []byte
}

// NewFernet 从 url-safe base64 编码的 32 字节密钥创建 Codec。
func NewFernet(key string) (*Codec, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, errors.New("crypto: empty fernet key")
	}
	raw, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		// 兼容无 padding
		raw, err = base64.RawURLEncoding.DecodeString(key)
		if err != nil {
			return nil, fmt.Errorf("crypto: decode fernet key: %w", err)
		}
	}
	if len(raw) != 32 {
		return nil, fmt.Errorf("crypto: fernet key must decode to 32 bytes, got %d", len(raw))
	}
	return &Codec{
		signingKey: raw[:16],
		encryptKey: raw[16:],
	}, nil
}

// NewFernetFromConfig 优先用 cfgKey；若 VaultAddr 非空且 cfgKey 为空，尝试环境变量 HEI_CRYPTO_FERNET_KEY。
func NewFernetFromConfig(fernetKey, vaultAddr string) (*Codec, error) {
	key := strings.TrimSpace(fernetKey)
	if key == "" && strings.TrimSpace(vaultAddr) != "" {
		// Vault 完整客户端未接入：仅回退读环境变量中的 Fernet 密钥。
		key = strings.TrimSpace(os.Getenv("HEI_CRYPTO_FERNET_KEY"))
	}
	if key == "" {
		return nil, nil
	}
	return NewFernet(key)
}

// Encrypt 加密明文，返回 url-safe base64 token。
func (c *Codec) Encrypt(plaintext string) (string, error) {
	if c == nil {
		return plaintext, nil
	}
	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	ts := uint64(time.Now().Unix())
	padded := pkcs7Pad([]byte(plaintext), aes.BlockSize)
	block, err := aes.NewCipher(c.encryptKey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)

	token := make([]byte, 1+8+16+len(ciphertext))
	token[0] = fernetVersion
	binary.BigEndian.PutUint64(token[1:9], ts)
	copy(token[9:25], iv)
	copy(token[25:], ciphertext)

	mac := hmac.New(sha256.New, c.signingKey)
	mac.Write(token)
	sig := mac.Sum(nil)
	return base64.URLEncoding.EncodeToString(append(token, sig...)), nil
}

// Decrypt 解密 Fernet token；失败返回错误。
func (c *Codec) Decrypt(token string) (string, error) {
	if c == nil {
		return token, nil
	}
	raw, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		raw, err = base64.RawURLEncoding.DecodeString(token)
		if err != nil {
			return "", fmt.Errorf("crypto: decode token: %w", err)
		}
	}
	if len(raw) < 1+8+16+16+32 {
		return "", errors.New("crypto: token too short")
	}
	if raw[0] != fernetVersion {
		return "", errors.New("crypto: bad fernet version")
	}
	body, sig := raw[:len(raw)-32], raw[len(raw)-32:]
	mac := hmac.New(sha256.New, c.signingKey)
	mac.Write(body)
	if !hmac.Equal(mac.Sum(nil), sig) {
		return "", errors.New("crypto: invalid mac")
	}
	if fernetTTL > 0 {
		ts := int64(binary.BigEndian.Uint64(body[1:9]))
		if time.Now().Unix()-ts > fernetTTL {
			return "", errors.New("crypto: token expired")
		}
	}
	iv := body[9:25]
	ct := body[25:]
	if len(ct)%aes.BlockSize != 0 {
		return "", errors.New("crypto: bad ciphertext length")
	}
	block, err := aes.NewCipher(c.encryptKey)
	if err != nil {
		return "", err
	}
	pt := make([]byte, len(ct))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(pt, ct)
	pt, err = pkcs7Unpad(pt, aes.BlockSize)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

// LooksEncrypted 启发式判断字符串是否像 Fernet token。
func LooksEncrypted(v string) bool {
	v = strings.TrimSpace(v)
	if len(v) < 60 {
		return false
	}
	raw, err := base64.URLEncoding.DecodeString(v)
	if err != nil {
		raw, err = base64.RawURLEncoding.DecodeString(v)
		if err != nil {
			return false
		}
	}
	return len(raw) >= 57 && raw[0] == fernetVersion
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - len(data)%blockSize
	out := make([]byte, len(data)+pad)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(pad)
	}
	return out
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("crypto: invalid padding size")
	}
	pad := int(data[len(data)-1])
	if pad == 0 || pad > blockSize || pad > len(data) {
		return nil, errors.New("crypto: invalid padding")
	}
	for i := len(data) - pad; i < len(data); i++ {
		if data[i] != byte(pad) {
			return nil, errors.New("crypto: invalid padding bytes")
		}
	}
	return data[:len(data)-pad], nil
}
