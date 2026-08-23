// Package security 提供账号登录名规则（字母/数字/下划线，3-64）。
//
// Author: Charlie
package security

import (
	"fmt"
	"regexp"
	"strings"
)

var accountLoginPattern = regexp.MustCompile(`^[a-zA-Z0-9_]{3,64}$`)

// RequireAccountLogin 校验并返回规范化后的账号登录名。
func RequireAccountLogin(account string) (string, error) {
	value := strings.TrimSpace(account)
	if value == "" {
		return "", fmt.Errorf("请输入用户名")
	}
	if !accountLoginPattern.MatchString(value) {
		return "", fmt.Errorf("账号仅允许字母、数字和下划线，长度 3-64")
	}
	return value, nil
}

// SanitizeAccountBase 从邮箱本地段等来源生成合法账号前缀。
func SanitizeAccountBase(raw string) string {
	var b strings.Builder
	for _, ch := range raw {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' {
			b.WriteRune(ch)
		}
	}
	cleaned := strings.ToLower(b.String())
	if cleaned == "" {
		cleaned = "user"
	}
	if len(cleaned) < 3 {
		cleaned += strings.Repeat("0", 3-len(cleaned))
	}
	if len(cleaned) > 48 {
		cleaned = cleaned[:48]
	}
	return cleaned
}
