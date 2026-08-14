// internal/framework/core/security/password.go 密码哈希。
//
// Author: Charlie

package security

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 对明文密码哈希。
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword 校验明文是否匹配 bcrypt 哈希。
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
