// content_disposition.go Content-Disposition 构建（ASCII fallback + RFC5987 filename*）。
//
// Author: Charlie
package file

import (
	"net/url"
	"strings"
	"unicode"
)

// contentDispositionAttachment 构建 attachment Content-Disposition 头（对齐 boot/fastapi）。
func contentDispositionAttachment(originalName string) string {
	name := strings.TrimSpace(originalName)
	if name == "" {
		name = "download"
	}
	var ascii strings.Builder
	for _, r := range name {
		if r == '"' || r < 0x20 || r > 0x7E || !unicode.IsPrint(r) {
			ascii.WriteByte('_')
			continue
		}
		ascii.WriteRune(r)
	}
	asciiName := strings.Trim(ascii.String(), "_")
	if asciiName == "" {
		asciiName = "download"
	}
	encoded := url.PathEscape(name)
	return `attachment; filename="` + asciiName + `"; filename*=UTF-8''` + encoded
}
