// Package htmlsafe 对存储型 HTML 做轻量消毒（剥离 script / 事件属性 / 危险 URL）。
//
// Author: Charlie
package htmlsafe

import (
	"regexp"
	"strings"
)

var (
	reScript     = regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script\s*>`)
	reStyle      = regexp.MustCompile(`(?is)<style\b[^>]*>.*?</style\s*>`)
	reIframe     = regexp.MustCompile(`(?is)<iframe\b[^>]*>.*?</iframe\s*>`)
	reObject     = regexp.MustCompile(`(?is)<object\b[^>]*>.*?</object\s*>`)
	reEmbed      = regexp.MustCompile(`(?is)<embed\b[^>]*/?>`)
	reEventAttr  = regexp.MustCompile(`(?i)\s+on[a-z]+\s*=\s*("[^"]*"|'[^']*'|[^\s>]+)`)
	reDangerHref = regexp.MustCompile(`(?i)\s(href|src|xlink:href|action|formaction)\s*=\s*("|')\s*(javascript|data|vbscript)\s*:`)
)

// Sanitize 若 contentType 为 HTML 则消毒；其它类型原样返回。
func Sanitize(contentType, content string) string {
	if !strings.EqualFold(strings.TrimSpace(contentType), "HTML") {
		return content
	}
	out := content
	out = reScript.ReplaceAllString(out, "")
	out = reStyle.ReplaceAllString(out, "")
	out = reIframe.ReplaceAllString(out, "")
	out = reObject.ReplaceAllString(out, "")
	out = reEmbed.ReplaceAllString(out, "")
	out = reEventAttr.ReplaceAllString(out, "")
	out = reDangerHref.ReplaceAllString(out, ` $1=$2#blocked`)
	return out
}
