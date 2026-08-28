// Package safeurl 校验出站 HTTP(S) URL，缓解 SSRF（禁止私网/元数据地址等）。
//
// Author: Charlie
package safeurl

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// Options 控制允许的协议（默认仅 https）。
type Options struct {
	AllowHTTP bool
}

// Validate 解析并校验 rawURL：无 userinfo、协议受限、主机 DNS 解析后不得为私网/环回/链路本地/未指定地址。
func Validate(rawURL string, opts Options) error {
	raw := strings.TrimSpace(rawURL)
	if raw == "" {
		return fmt.Errorf("empty url")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case "https":
		// ok
	case "http":
		if !opts.AllowHTTP {
			return fmt.Errorf("http scheme not allowed")
		}
	default:
		return fmt.Errorf("scheme not allowed: %s", u.Scheme)
	}
	if u.User != nil {
		return fmt.Errorf("url userinfo not allowed")
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("missing host")
	}
	if isBlockedHostLiteral(host) {
		return fmt.Errorf("blocked host: %s", host)
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("dns lookup failed: %w", err)
	}
	if len(ips) == 0 {
		return fmt.Errorf("dns lookup returned no addresses")
	}
	for _, ip := range ips {
		if isBlockedIP(ip) {
			return fmt.Errorf("blocked address: %s", ip.String())
		}
	}
	return nil
}

func isBlockedHostLiteral(host string) bool {
	h := strings.ToLower(strings.TrimSpace(host))
	switch h {
	case "localhost", "metadata.google.internal", "metadata":
		return true
	}
	return false
}

func isBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified() || ip.IsMulticast() {
		return true
	}
	// AWS/GCP/Azure 等云元数据常见链路本地；IsLinkLocalUnicast 已覆盖 169.254.0.0/16。
	// 额外拒绝 CGNAT / 文档网段等常被滥用地址。
	if ip4 := ip.To4(); ip4 != nil {
		if ip4[0] == 100 && ip4[1] >= 64 && ip4[1] <= 127 { // 100.64.0.0/10
			return true
		}
		if ip4[0] == 192 && ip4[1] == 0 && ip4[2] == 0 { // 192.0.0.0/24
			return true
		}
		if ip4[0] == 192 && ip4[1] == 0 && ip4[2] == 2 { // 192.0.2.0/24 TEST-NET-1
			return true
		}
		if ip4[0] == 198 && (ip4[1] == 18 || ip4[1] == 19) { // 198.18.0.0/15
			return true
		}
		if ip4[0] == 198 && ip4[1] == 51 && ip4[2] == 100 { // TEST-NET-2
			return true
		}
		if ip4[0] == 203 && ip4[1] == 0 && ip4[2] == 113 { // TEST-NET-3
			return true
		}
	}
	if ip.To4() == nil {
		// IPv6 ULA fc00::/7 — IsPrivate 在 Go 1.17+ 已覆盖；再拦 unique local 边缘情况。
		if len(ip) == net.IPv6len && (ip[0]&0xfe) == 0xfc {
			return true
		}
	}
	return false
}
