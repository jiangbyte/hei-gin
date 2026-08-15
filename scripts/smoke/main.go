// Package main 后端启动冒烟测试（临时）：验证码 -> RSA 加密 -> 登录 -> /me。
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const base = "http://127.0.0.1:8000/api/v1/admin"

func getJSON(url string, out any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("json %s: %w (%s)", url, err, string(body))
	}
	return nil
}

func postJSON(url string, payload any, out any) error {
	raw, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewReader(raw))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("json %s: %w (%s)", url, err, string(body))
		}
	}
	return nil
}

func main() {
	// 1. captcha
	var cap struct {
		Data struct {
			CaptchaID   string `json:"captcha_id"`
			ImageBase64 string `json:"image_base64"`
		} `json:"data"`
	}
	if err := getJSON(base+"/captcha", &cap); err != nil {
		panic(err)
	}
	svg, err := base64.StdEncoding.DecodeString(cap.Data.ImageBase64)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`<text[^>]*>([^<]+)</text>`)
	var answer string
	for _, m := range re.FindAllStringSubmatch(string(svg), -1) {
		answer += m[1]
	}
	fmt.Println("captcha_id:", cap.Data.CaptchaID, "answer:", answer)

	// 2. password key
	var pk struct {
		Data struct {
			KeyID     string `json:"key_id"`
			PublicKey string `json:"public_key"`
		} `json:"data"`
	}
	if err := getJSON(base+"/password-key", &pk); err != nil {
		panic(err)
	}
	der, err := base64.StdEncoding.DecodeString(pk.Data.PublicKey)
	if err != nil {
		panic(err)
	}
	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		panic(err)
	}
	rsaPub := pub.(*rsa.PublicKey)
	enc, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, []byte("123456"), nil)
	if err != nil {
		panic(err)
	}
	password := base64.StdEncoding.EncodeToString(enc)

	// 3. login
	var login struct {
		Code string `json:"code"`
		Data struct {
			Token           string `json:"token"`
			AccountID       string `json:"account_id"`
			ForceBind       string `json:"force_bind_email"`
			PasswordExpired string `json:"password_expired"`
		} `json:"data"`
		Message string `json:"message"`
	}
	err = postJSON(base+"/login", map[string]any{
		"account":         "superadmin",
		"password":        password,
		"password_key_id": pk.Data.KeyID,
		"captcha_id":      cap.Data.CaptchaID,
		"captcha_value":   answer,
		"remember_me":     "true",
		"login_mode":      "PASSWORD",
	}, &login)
	if err != nil {
		panic(err)
	}
	fmt.Println("login code:", login.Code, "message:", login.Message, "token_len:", len(login.Data.Token), "password_expired:", login.Data.PasswordExpired)

	// 4. /me (authenticated, reads DB profile + id names)
	req, _ := http.NewRequest("GET", base+"/me", nil)
	req.Header.Set("Authorization", login.Data.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var me struct {
		Data struct {
			Account         string                      `json:"account"`
			AccountType     string                      `json:"account_type"`
			Name            *string                     `json:"name"`
			Nickname        *string                     `json:"nickname"`
			Avatar          *string                     `json:"avatar"`
			RoleIDNames     []struct{ ID, Name string } `json:"role_id_names"`
			DeptIDNames     []struct{ ID, Name string } `json:"dept_id_names"`
			GroupIDNames    []struct{ ID, Name string } `json:"group_id_names"`
			PermissionKeys  []string                    `json:"permission_keys"`
			PasswordExpired string                      `json:"password_expired"`
			ForceBindEmail  string                      `json:"force_bind_email"`
			Profile         map[string]any              `json:"profile"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &me); err != nil {
		panic(err)
	}
	fmt.Println("/me message:", me.Message, "account:", me.Data.Account, "type:", me.Data.AccountType,
		"name:", ptrStr(me.Data.Name), "nickname:", ptrStr(me.Data.Nickname),
		"roles:", len(me.Data.RoleIDNames), "depts:", len(me.Data.DeptIDNames), "groups:", len(me.Data.GroupIDNames),
		"perms:", len(me.Data.PermissionKeys), "password_expired:", me.Data.PasswordExpired,
		"force_bind_email:", me.Data.ForceBindEmail, "profile_keys:", len(me.Data.Profile))
	if len(me.Data.PermissionKeys) > 0 {
		fmt.Println("  sample perms:", strings.Join(me.Data.PermissionKeys[:min(6, len(me.Data.PermissionKeys))], ", "))
		authGET(login.Data.Token, "/dashboard/overview")
		authGET(login.Data.Token, "/sys/config/list?category=AUTH_PASSWORD")
		authGET(login.Data.Token, "/sys/notices/my-page?current=1&size=3")
		authGET(login.Data.Token, "/sys/banners/page?current=1&size=10")
		authGET(login.Data.Token, "/sys/banners/list?position=HOME_TOP")
		authGET(login.Data.Token, "/sys/banners/detail?id=7491889345134235648")
		authGET(login.Data.Token, "/sys/roles/own-resource?id=1&account_type=ADMIN")
		authGET(login.Data.Token, "/sys/roles/own-client-resource?id=1&account_type=ADMIN")
		authGET(login.Data.Token, "/sys/resources/current")
		authGET(login.Data.Token, "/sys/client-resources/tree?account_type=ADMIN")
		authGETFile(login.Data.Token)
	}
}

func ptrStr(p *string) string {
	if p == nil {
		return "<nil>"
	}
	return *p
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// authGET 带 token 调用受保护接口。
func authGET(token, path string) {
	req, _ := http.NewRequest("GET", base+path, nil)
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("GET", path, "err:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("GET", path, "status:", resp.StatusCode, "body:", truncate(string(body), 180))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// authGETFile 文件分页验证（per-provider URL + created_name）。
func authGETFile(token string) {
	req, _ := http.NewRequest("GET", base+"/sys/file/page?current=1&size=20", nil)
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("GET file page err:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var r struct {
		Code string `json:"code"`
		Data struct {
			Records []struct {
				ID              string  `json:"id"`
				ObjectName      string  `json:"object_name"`
				StorageProvider string  `json:"storage_provider"`
				URL             string  `json:"url"`
				CreatedName     *string `json:"created_name"`
				UpdatedName     *string `json:"updated_name"`
			} `json:"records"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &r); err != nil {
		fmt.Println("file page parse err:", err, string(body[:min(len(body), 120)]))
		return
	}
	fmt.Println("file page code:", r.Code, "records:", len(r.Data.Records))
	for i, rec := range r.Data.Records {
		if i >= 4 {
			break
		}
		fmt.Printf("  #%d provider=%s object=%s url=%s created_name=%s\n", i, rec.StorageProvider, rec.ObjectName, truncate(rec.URL, 80), ptrStr(rec.CreatedName))
	}
}
