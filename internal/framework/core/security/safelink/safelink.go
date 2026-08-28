// Package safelink 校验面向浏览器的链接（Banner 等），拒绝危险 scheme。
//
// Author: Charlie
package safelink

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var dangerousScheme = regexp.MustCompile(`(?i)^\s*(javascript|data|vbscript|blob)\s*:`)

// ValidateBannerLink 按 linkType 校验 url。
// URL: 允许 http(s) 绝对地址或以 / 开头的相对路径；拒绝 javascript/data 等。
// ROUTE: 必须为相对路径（以 / 开头，无 scheme）。
// NONE: 允许空 url。
func ValidateBannerLink(linkType, rawURL string) error {
	lt := strings.ToUpper(strings.TrimSpace(linkType))
	if lt == "" {
		lt = "URL"
	}
	u := strings.TrimSpace(rawURL)
	switch lt {
	case "NONE":
		return nil
	case "ROUTE":
		if u == "" {
			return fmt.Errorf("route link requires path")
		}
		return validateRelativePath(u)
	case "URL":
		if u == "" {
			return nil
		}
		return ValidatePublicHref(u)
	default:
		return fmt.Errorf("unsupported link_type: %s", linkType)
	}
}

// ValidatePublicHref 允许 http(s) 或相对 / 路径，拒绝危险 scheme。
func ValidatePublicHref(raw string) error {
	u := strings.TrimSpace(raw)
	if u == "" {
		return fmt.Errorf("empty url")
	}
	if dangerousScheme.MatchString(u) {
		return fmt.Errorf("dangerous url scheme")
	}
	if strings.HasPrefix(u, "/") {
		return validateRelativePath(u)
	}
	parsed, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("scheme not allowed: %s", parsed.Scheme)
	}
	if parsed.Host == "" {
		return fmt.Errorf("missing host")
	}
	return nil
}

func validateRelativePath(u string) error {
	if !strings.HasPrefix(u, "/") {
		return fmt.Errorf("path must start with /")
	}
	if strings.HasPrefix(u, "//") {
		return fmt.Errorf("protocol-relative url not allowed")
	}
	if dangerousScheme.MatchString(u) {
		return fmt.Errorf("dangerous url scheme")
	}
	if strings.ContainsAny(u, "\r\n\x00") {
		return fmt.Errorf("invalid path characters")
	}
	return nil
}
