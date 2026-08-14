// internal/modules/auth/captcha.go 验证码。
//
// Author: Charlie

package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html"
	"strings"
)

// 验证码与密码传输、缓存键约定。
const (
	captchaAlphabet   = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	captchaKeyPrefix  = "captcha:"
	passwordKeyPrefix = "password:crypto:"
)

func captchaSVGBase64(value string) string {
	escaped := html.EscapeString(value)
	var noise strings.Builder
	for i := 0; i < 6; i++ {
		noise.WriteString(fmt.Sprintf(
			`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#94a3b8" stroke-width="1" opacity="0.45" />`,
			randBelow(140), randBelow(44), randBelow(140), randBelow(44),
		))
	}
	var texts strings.Builder
	for i, ch := range escaped {
		rot := randBelow(21) - 10
		x := 22 + i*26
		y := 29 + randBelow(5)
		texts.WriteString(fmt.Sprintf(
			`<text x="%d" y="%d" font-size="24" font-family="Arial, sans-serif" font-weight="700" fill="#0f172a" transform="rotate(%d %d 25)">%s</text>`,
			x, y, rot, x, string(ch),
		))
	}
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="140" height="44" viewBox="0 0 140 44">` +
		`<rect width="140" height="44" rx="6" fill="#f8fafc"/>` +
		noise.String() + texts.String() + `</svg>`
	return base64.StdEncoding.EncodeToString([]byte(svg))
}

func randBelow(n int) int {
	if n <= 0 {
		return 0
	}
	b := make([]byte, 1)
	_, _ = rand.Read(b)
	return int(b[0]) % n
}
