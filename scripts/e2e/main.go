// Package main 双库 API 兼容扫路由：Redis 植入验证码 → 登录 → 全路由扫描 → 方言硬检查。
//
// Usage:
//
//	go run ./scripts/e2e -base http://127.0.0.1:8000 -redis redis://:123456@127.0.0.1:6379/4 -out scripts/e2e/report-mysql.json
package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"hei-gin/internal/framework/core/stringly"
)

type apiResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type routeInfo struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type callResult struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	URL        string `json:"url"`
	Status     int    `json:"status"`
	BizCode    int    `json:"biz_code,omitempty"`
	Body       string `json:"body"`
	Error      string `json:"error,omitempty"`
	HardCheck  string `json:"hard_check,omitempty"`
	HardOK     *bool  `json:"hard_ok,omitempty"`
	Is5xx      bool   `json:"is_5xx"`
	SQLSuspect bool   `json:"sql_suspect"`
}

type report struct {
	BaseURL     string       `json:"base_url"`
	GeneratedAt string       `json:"generated_at"`
	AdminOK     bool         `json:"admin_login_ok"`
	PortalOK    bool         `json:"portal_login_ok"`
	ReadCases   caseBucket   `json:"read_cases"`
	CRUDCases   caseBucket   `json:"crud_cases"`
	Skipped     []caseResult `json:"skipped"`
	Total       int          `json:"total"`
	OK2xx4xx    int          `json:"ok_2xx_4xx"`
	Fail5xx     int          `json:"fail_5xx"`
	HardPass    int          `json:"hard_pass"`
	HardFail    int          `json:"hard_fail"`
	Results     []callResult `json:"results"`
	Fail5xxList []callResult `json:"fail_5xx_list"`
	HardFails   []callResult `json:"hard_fails"`
}

func main() {
	base := flag.String("base", "http://127.0.0.1:8000", "API base URL")
	redisURL := flag.String("redis", "redis://:123456@127.0.0.1:6379/4", "Redis URL for captcha plant")
	outPath := flag.String("out", "scripts/e2e/report.json", "report JSON path")
	mode := flag.String("mode", "full", "full|sweep|shape (full=shape+crud+sweep)")
	flag.Parse()

	ctx := context.Background()
	rdb, err := openRedis(*redisURL)
	if err != nil {
		fatalf("redis: %v", err)
	}
	defer rdb.Close()

	adminTok, err := login(ctx, rdb, *base, "/api/v1/admin", "superadmin")
	adminOK := err == nil
	if err != nil {
		fmt.Println("admin login FAILED:", err)
	} else {
		fmt.Println("admin login OK, token_len=", len(adminTok))
	}
	portalTok, err := login(ctx, rdb, *base, "/api/v1/portal", "user")
	portalOK := err == nil
	if err != nil {
		fmt.Println("portal login FAILED:", err)
	} else {
		fmt.Println("portal login OK, token_len=", len(portalTok))
	}

	rep := report{
		BaseURL:     *base,
		GeneratedAt: time.Now().Format(time.RFC3339),
		AdminOK:     adminOK,
		PortalOK:    portalOK,
	}

	doShape := *mode == "full" || *mode == "shape"
	doSweep := *mode == "full" || *mode == "sweep"

	routes, err := fetchRoutes(*base, adminTok)
	if err != nil {
		fatalf("fetch routes: %v", err)
	}
	fmt.Println("routes:", len(routes))

	if doShape && adminOK {
		fmt.Println("--- UNIVERSE (ALL APIs shape + CRUD) ---")
		runUniverse(*base, adminTok, portalTok, routes, &rep.ReadCases, &rep.CRUDCases, &rep.Skipped)
	}

	if doSweep {
		// Hard checks first (must succeed business-wise).
		for _, hc := range hardChecks() {
			tok := adminTok
			if strings.Contains(hc.path, "/portal/") {
				tok = portalTok
			}
			cr := doCall(hc.method, *base+hc.path, tok, hc.body)
			cr.HardCheck = hc.name
			ok := cr.Status >= 200 && cr.Status < 300 && (cr.BizCode == 0 || cr.BizCode == 200) && !cr.Is5xx
			cr.HardOK = &ok
			if ok {
				rep.HardPass++
				fmt.Println("HARD PASS", hc.name, cr.Status)
			} else {
				rep.HardFail++
				rep.HardFails = append(rep.HardFails, cr)
				fmt.Println("HARD FAIL", hc.name, cr.Status, cr.BizCode, truncate(cr.Body, 160))
			}
			rep.Results = append(rep.Results, cr)
		}

		seen := map[string]bool{}
		for _, r := range routes {
			key := r.Method + " " + r.Path
			if seen[key] {
				continue
			}
			seen[key] = true
			if reason, skip := skipRouteReason(r); skip {
				rep.Skipped = append(rep.Skipped, caseResult{Name: key, Error: reason, OK: true})
				continue
			}
			path := materializePath(r.Path)
			url := *base + path
			if strings.Contains(path, "/page") || strings.HasSuffix(path, "/list") || strings.HasSuffix(path, "/tree") {
				if !strings.Contains(url, "?") {
					url += "?current=1&size=10"
				} else {
					url += "&current=1&size=10"
				}
				if strings.Contains(path, "/roles/") || strings.Contains(path, "/accounts") {
					url += "&name=a"
				}
				if strings.Contains(path, "/banners") {
					url += "&title=a"
				}
				if strings.Contains(path, "/config") {
					url += "&config_key=a"
				}
				if strings.Contains(path, "/dicts") {
					url += "&code=a"
				}
			}
			tok := pickToken(r.Path, adminTok, portalTok)
			body := ""
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
				body = "{}"
			}
			cr := doCall(r.Method, url, tok, body)
			rep.Total++
			if cr.Is5xx {
				rep.Fail5xx++
				rep.Fail5xxList = append(rep.Fail5xxList, cr)
				fmt.Println("5xx", r.Method, path, cr.Status, truncate(cr.Body, 120))
			} else {
				rep.OK2xx4xx++
			}
			rep.Results = append(rep.Results, cr)
		}
	}

	raw, _ := json.MarshalIndent(rep, "", "  ")
	if err := os.WriteFile(*outPath, raw, 0o644); err != nil {
		fatalf("write report: %v", err)
	}
	fmt.Printf("\n=== SUMMARY ===\nadmin=%v portal=%v get_shape=%d/%d crud+writes=%d/%d skipped=%d sweep=%d ok=%d fail5xx=%d hard=%d/%d\nout=%s\n",
		rep.AdminOK, rep.PortalOK,
		rep.ReadCases.Pass, rep.ReadCases.Total,
		rep.CRUDCases.Pass, rep.CRUDCases.Total,
		len(rep.Skipped),
		rep.Total, rep.OK2xx4xx, rep.Fail5xx, rep.HardPass, rep.HardPass+rep.HardFail, *outPath)
	fail := !rep.AdminOK || !rep.PortalOK || rep.Fail5xx > 0 || rep.HardFail > 0 ||
		len(rep.ReadCases.Fail) > 0 || len(rep.CRUDCases.Fail) > 0
	if fail {
		os.Exit(1)
	}
}

type hardCheck struct {
	name, method, path, body string
}

func hardChecks() []hardCheck {
	return []hardCheck{
		{name: "health_live", method: "GET", path: "/api/v1/internal/health/live"},
		{name: "health_ready", method: "GET", path: "/api/v1/internal/health/ready"},
		{name: "dashboard_overview", method: "GET", path: "/api/v1/admin/dashboard/overview"},
		{name: "roles_page_ilike", method: "GET", path: "/api/v1/admin/sys/roles/page?current=1&size=5&name=admin"},
		{name: "banners_page_json", method: "GET", path: "/api/v1/admin/sys/banners/page?current=1&size=5&title=a"},
		{name: "banners_list", method: "GET", path: "/api/v1/admin/sys/banners/list?position=HOME_TOP"},
		{name: "notices_page", method: "GET", path: "/api/v1/admin/sys/notices/page?current=1&size=5&title=a"},
		{name: "notices_my_page", method: "GET", path: "/api/v1/admin/sys/notices/my-page?current=1&size=5"},
		{name: "accounts_page_ilike", method: "GET", path: "/api/v1/admin/sys/accounts/page?current=1&size=5&account=super"},
		{name: "codegen_tables", method: "GET", path: "/api/v1/admin/sys/codegen/tables"},
		{name: "portal_banners_list", method: "GET", path: "/api/v1/portal/sys/banners/list"},
		{name: "portal_dicts_tree", method: "GET", path: "/api/v1/portal/sys/dicts/tree"},
		{name: "portal_notices_list", method: "GET", path: "/api/v1/portal/sys/notices/list"},
		{name: "admin_me", method: "GET", path: "/api/v1/admin/me"},
	}
}

func skipRoute(r routeInfo) bool {
	_, skip := skipRouteReason(r)
	return skip
}

func skipRouteReason(r routeInfo) (string, bool) {
	if r.Method == "HEAD" || r.Method == "OPTIONS" || r.Method == "CONNECT" {
		return "method", true
	}
	p := r.Path
	switch {
	case p == "/api/v1/internal/debug/routes":
		return "debug", true
	case p == "/metrics":
		return "metrics", true
	case strings.Contains(p, "/oauth/"):
		return "oauth", true
	case strings.HasSuffix(p, "/login"):
		return "auth-bootstrap", true
	case strings.HasSuffix(p, "/captcha"):
		return "auth-bootstrap", true
	case strings.HasSuffix(p, "/password-key"):
		return "auth-bootstrap", true
	case strings.Contains(p, "/cancel"):
		return "session-destructive", true
	case strings.Contains(p, "/logout"):
		return "session-destructive", true
	case strings.Contains(p, "/upload"):
		return "storage", true
	case strings.Contains(p, "/avatar"):
		return "storage", true
	case strings.Contains(p, "/download"):
		return "binary", true
	}
	return "", false
}

func pickToken(path, admin, portal string) string {
	if strings.Contains(path, "/portal/") {
		return portal
	}
	if strings.Contains(path, "/internal/") {
		return ""
	}
	return admin
}

func materializePath(p string) string {
	// Gin param placeholders
	p = strings.ReplaceAll(p, ":provider", "github")
	p = strings.ReplaceAll(p, ":id", "0")
	return p
}

func fetchRoutes(base, token string) ([]routeInfo, error) {
	var wrap struct {
		Data []routeInfo `json:"data"`
	}
	req, err := http.NewRequest("GET", base+"/api/v1/internal/debug/routes", nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("routes status %d: %s", resp.StatusCode, truncate(string(b), 200))
	}
	if err := stringly.Unmarshal(b, &wrap); err != nil {
		return nil, err
	}
	return wrap.Data, nil
}

func login(ctx context.Context, rdb *redis.Client, base, prefix, account string) (string, error) {
	var cap struct {
		Data struct {
			CaptchaID string `json:"captcha_id"`
		} `json:"data"`
	}
	if err := getJSON(base+prefix+"/captcha", &cap); err != nil {
		return "", fmt.Errorf("captcha: %w", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	if err := rdb.Set(ctx, "captcha:"+cap.Data.CaptchaID, string(hash), 5*time.Minute).Err(); err != nil {
		return "", fmt.Errorf("plant captcha: %w", err)
	}

	var pk struct {
		Data struct {
			KeyID     string `json:"key_id"`
			PublicKey string `json:"public_key"`
		} `json:"data"`
	}
	if err := getJSON(base+prefix+"/password-key", &pk); err != nil {
		return "", fmt.Errorf("password-key: %w", err)
	}
	der, err := base64.StdEncoding.DecodeString(pk.Data.PublicKey)
	if err != nil {
		return "", err
	}
	pubAny, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return "", err
	}
	enc, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubAny.(*rsa.PublicKey), []byte("123456"), nil)
	if err != nil {
		return "", err
	}
	var login struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := postJSON(base+prefix+"/login", map[string]any{
		"account":         account,
		"password":        base64.StdEncoding.EncodeToString(enc),
		"password_key_id": pk.Data.KeyID,
		"captcha_id":      cap.Data.CaptchaID,
		"captcha_value":   "test",
		"remember_me":     true,
		"login_mode":      "PASSWORD",
	}, &login); err != nil {
		return "", err
	}
	if login.Data.Token == "" {
		return "", fmt.Errorf("login code=%d msg=%s", login.Code, login.Message)
	}
	return login.Data.Token, nil
}

func doCall(method, url, token, body string) callResult {
	cr := callResult{Method: method, Path: url, URL: url}
	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		cr.Error = err.Error()
		cr.Is5xx = true
		return cr
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		cr.Error = err.Error()
		cr.Is5xx = true
		return cr
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	cr.Status = resp.StatusCode
	cr.Body = truncate(string(b), 400)
	cr.Is5xx = resp.StatusCode >= 500
	var ar apiResp
	if parsed, _, err := parseLoose(b); err == nil {
		ar = parsed
		cr.BizCode = ar.Code
		if ar.Code >= 500 {
			cr.Is5xx = true
		}
	}
	low := strings.ToLower(cr.Body + cr.Error)
	if strings.Contains(low, "sql") || strings.Contains(low, "dialect") ||
		strings.Contains(low, "pq:") || strings.Contains(low, "mysql") ||
		strings.Contains(low, "ilike") || strings.Contains(low, "jsonb") ||
		strings.Contains(low, "syntax error") {
		cr.SQLSuspect = true
	}
	// extract path from url for report readability
	if i := strings.Index(url, "/api/"); i >= 0 {
		cr.Path = url[i:]
	}
	return cr
}

func openRedis(raw string) (*redis.Client, error) {
	opt, err := redis.ParseURL(raw)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}

func getJSON(url string, out any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 500 {
		return fmt.Errorf("GET %s status %d: %s", url, resp.StatusCode, truncate(string(b), 200))
	}
	return stringly.Unmarshal(b, out)
}

func postJSON(url string, payload, out any) error {
	raw, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewReader(raw))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if out != nil {
		return stringly.Unmarshal(b, out)
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func fatalf(f string, a ...any) {
	fmt.Fprintf(os.Stderr, f+"\n", a...)
	os.Exit(2)
}
