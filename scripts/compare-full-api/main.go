// 全量 API 契约对比：hei-boot OpenAPI vs hei-gin debug/routes。
//
// 用法：
//
//	go run ./scripts/compare-full-api --boot http://127.0.0.1:8000 --gin http://127.0.0.1:8001
package main

import (
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
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type openapiDoc struct {
	Paths map[string]map[string]json.RawMessage `json:"paths"`
}

type routeInfo struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type apiResp struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

func main() {
	boot := flag.String("boot", "http://127.0.0.1:8000", "hei-boot base url")
	gin := flag.String("gin", "http://127.0.0.1:8001", "hei-gin base url")
	openAPI := flag.String("openapi", "scripts/boot-openapi.json", "boot openapi json file")
	out := flag.String("out", "scripts/full-api-diff.json", "report output")
	flag.Parse()

	oa, err := loadOpenAPI(*openAPI)
	if err != nil {
		fatal(err)
	}
	bootRoutes := openapiRoutes(oa)
	ginRoutes, err := fetchGinRoutes(*gin)
	if err != nil {
		fatal(err)
	}

	onlyBoot := diff(bootRoutes, ginRoutes)
	onlyGin := diff(ginRoutes, bootRoutes)

	report := map[string]any{
		"boot_base":    *boot,
		"gin_base":     *gin,
		"boot_count":   len(bootRoutes),
		"gin_count":    len(ginRoutes),
		"common_count": len(intersect(bootRoutes, ginRoutes)),
		"only_boot":    formatRoutes(onlyBoot),
		"only_gin":     formatRoutes(onlyGin),
		"generated_at": time.Now().Format(time.RFC3339),
	}
	raw, _ := json.MarshalIndent(report, "", "  ")
	if err := os.WriteFile(*out, raw, 0644); err != nil {
		fatal(err)
	}

	fmt.Printf("boot=%d gin=%d common=%d only_boot=%d only_gin=%d\n",
		len(bootRoutes), len(ginRoutes), len(intersect(bootRoutes, ginRoutes)), len(onlyBoot), len(onlyGin))
	fmt.Println("only_boot:")
	for _, r := range onlyBoot {
		fmt.Println(" ", r.method, r.path)
	}
	fmt.Println("only_gin:")
	for _, r := range onlyGin {
		fmt.Println(" ", r.method, r.path)
	}
	fmt.Println("written", *out)

	if len(onlyBoot) > 0 || len(onlyGin) > 0 {
		os.Exit(1)
	}
}

type routeKey struct{ method, path string }

func loadOpenAPI(path string) (*openapiDoc, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc openapiDoc
	return &doc, json.Unmarshal(b, &doc)
}

var paramRe = regexp.MustCompile(`\{[^}]+\}`)
var ginParamRe = regexp.MustCompile(`:[^/]+`)

func normPath(p string) string {
	p = paramRe.ReplaceAllString(p, ":id")
	p = ginParamRe.ReplaceAllString(p, ":id")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}

func openapiRoutes(doc *openapiDoc) map[routeKey]bool {
	out := map[routeKey]bool{}
	for path, methods := range doc.Paths {
		if !strings.HasPrefix(path, "/api/") &&
			path != "/" &&
			!strings.HasPrefix(path, "/easyTrans/") &&
			path != "/v3/api-docs" &&
			path != "/doc.html" &&
			!strings.HasPrefix(path, "/swagger-ui") {
			continue
		}
		if strings.HasPrefix(path, "/api/v1/admin/dashboard/") {
			continue
		}
		for m := range methods {
			mu := strings.ToUpper(m)
			if mu == "GET" || mu == "POST" || mu == "PUT" || mu == "PATCH" || mu == "DELETE" {
				out[routeKey{mu, normPath(path)}] = true
			}
		}
	}
	return out
}

func fetchGinRoutes(base string) (map[routeKey]bool, error) {
	tok, err := loginAdmin(base)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("GET", base+"/api/v1/internal/debug/routes", nil)
	req.Header.Set("Authorization", tok)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var ar struct {
		Data []routeInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &ar); err != nil {
		return nil, err
	}
	out := map[routeKey]bool{}
	for _, r := range ar.Data {
		out[routeKey{strings.ToUpper(r.Method), normPath(r.Path)}] = true
	}
	return out, nil
}

func loginAdmin(base string) (string, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "123456", DB: 4})
	defer rdb.Close()

	capID, err := getCaptchaID(base + "/api/v1/admin/captcha")
	if err != nil {
		return "", err
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err := rdb.Set(ctx, "captcha:"+capID, string(hash), 5*time.Minute).Err(); err != nil {
		return "", err
	}

	keyID, pubKey, err := getPasswordKey(base + "/api/v1/admin/password-key")
	if err != nil {
		return "", err
	}
	encPwd, err := encryptPasswordOAEP(pubKey, "123456")
	if err != nil {
		return "", err
	}

	loginBody := map[string]any{
		"account": "superadmin", "password": encPwd, "password_key_id": keyID,
		"login_mode": "PASSWORD", "captcha_id": capID, "captcha_value": "test", "remember_me": true,
	}
	b, _ := json.Marshal(loginBody)
	resp, err := http.Post(base+"/api/v1/admin/login", "application/json", strings.NewReader(string(b)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var ar struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &ar); err != nil {
		return "", err
	}
	if ar.Data.Token == "" {
		return "", fmt.Errorf("login failed: %s", string(raw))
	}
	return ar.Data.Token, nil
}

func getCaptchaID(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var ar struct {
		Data struct {
			CaptchaID string `json:"captcha_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &ar); err != nil {
		return "", err
	}
	return ar.Data.CaptchaID, nil
}

func getPasswordKey(url string) (keyID, pubB64 string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var ar struct {
		Data struct {
			KeyID     string `json:"key_id"`
			PublicKey string `json:"public_key"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &ar); err != nil {
		return "", "", err
	}
	return ar.Data.KeyID, ar.Data.PublicKey, nil
}

func diff(a, b map[routeKey]bool) []routeKey {
	var out []routeKey
	for k := range a {
		if !b[k] {
			out = append(out, k)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].path == out[j].path {
			return out[i].method < out[j].method
		}
		return out[i].path < out[j].path
	})
	return out
}

func intersect(a, b map[routeKey]bool) map[routeKey]bool {
	out := map[routeKey]bool{}
	for k := range a {
		if b[k] {
			out[k] = true
		}
	}
	return out
}

func formatRoutes(rs []routeKey) []string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.method+" "+r.path)
	}
	return out
}

func encryptPasswordOAEP(pubB64, password string) (string, error) {
	der, err := base64.StdEncoding.DecodeString(pubB64)
	if err != nil {
		return "", err
	}
	pubAny, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return "", err
	}
	enc, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubAny.(*rsa.PublicKey), []byte(password), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enc), nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
