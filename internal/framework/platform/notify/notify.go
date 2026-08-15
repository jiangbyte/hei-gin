// Package notify 提供邮件 / 短信 / 推送通知门面。
//
// Author: Charlie
package notify

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/logger"
)

// Facade 通知发送门面（运行时配置以 sys_config 为权威，回退 yaml）。
//
// Author: Charlie
type Facade struct {
	cfg config.NotifyConfig
	db  *gorm.DB
	hc  *http.Client
}

// NewFacade 构造通知门面。
func NewFacade(cfg config.NotifyConfig, db *gorm.DB) *Facade {
	return &Facade{
		cfg: cfg,
		db:  db,
		hc:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GetRuntimeString 从 sys_config 读取运行时字符串；无则返回 def。
func (f *Facade) GetRuntimeString(ctx context.Context, key, def string) string {
	if f == nil || f.db == nil || key == "" {
		return def
	}
	var v string
	if err := f.db.WithContext(ctx).Table("sys_config").Select("config_value").
		Where("config_key = ?", key).Limit(1).Scan(&v).Error; err != nil || v == "" {
		return def
	}
	return v
}

func (f *Facade) runtime(ctx context.Context, key, def string) string {
	return f.GetRuntimeString(ctx, key, def)
}

func (f *Facade) runtimeInt(ctx context.Context, key string, def int) int {
	v := strings.TrimSpace(f.runtime(ctx, key, ""))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func (f *Facade) runtimeBool(ctx context.Context, key string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(f.runtime(ctx, key, "")))
	if v == "" {
		return def
	}
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

// SendMail 发送纯文本邮件（LOCAL SMTP 用运行时 MAIL_LOCAL_* 配置；云引擎打桩）。
func (f *Facade) SendMail(ctx context.Context, to, subject, body string) error {
	_ = ctx
	engine := strings.ToUpper(strings.TrimSpace(f.runtime(ctx, "DEFAULT_EMAIL_ENGINE", "LOCAL")))
	if engine != "" && engine != "LOCAL" {
		logger.L.Info("mail engine stub",
			zap.String("engine", engine), zap.String("to", to), zap.String("subject", subject))
		return nil
	}
	host := f.runtime(ctx, "MAIL_LOCAL_HOST", f.cfg.Mail.Host)
	port := f.runtimeInt(ctx, "MAIL_LOCAL_PORT", f.cfg.Mail.Port)
	if port <= 0 {
		port = 587
	}
	username := f.runtime(ctx, "MAIL_LOCAL_USERNAME", f.cfg.Mail.Username)
	password := f.runtime(ctx, "MAIL_LOCAL_PASSWORD", f.cfg.Mail.Password)
	from := f.runtime(ctx, "MAIL_LOCAL_FROM_EMAIL", f.cfg.Mail.From)
	fromName := f.runtime(ctx, "MAIL_LOCAL_FROM_NAME", "")
	if from == "" {
		from = username
	}
	if host == "" {
		logger.L.Info("mail not configured, skip send",
			zap.String("to", to), zap.String("subject", subject))
		return nil
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	displayFrom := from
	if fromName != "" {
		displayFrom = fmt.Sprintf("%s <%s>", fromName, from)
	}
	crlf := "\r\n"
	msg := []byte("From: " + displayFrom + crlf +
		"To: " + to + crlf +
		"Subject: " + subject + crlf +
		"MIME-Version: 1.0" + crlf +
		"Content-Type: text/plain; charset=UTF-8" + crlf +
		crlf + body + crlf)
	authRequired := f.runtimeBool(ctx, "MAIL_LOCAL_AUTH_REQUIRED", true)
	var auth smtp.Auth
	if authRequired && username != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}
	useSSL := f.runtimeBool(ctx, "MAIL_LOCAL_USE_SSL", f.cfg.Mail.SSL)
	useStartTLS := f.runtimeBool(ctx, "MAIL_LOCAL_USE_STARTTLS", false)
	if useSSL {
		return sendMailTLS(addr, host, auth, from, []string{to}, msg)
	}
	if useStartTLS {
		return sendMailStartTLS(addr, host, auth, from, []string{to}, msg)
	}
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}

// SendSMS 发送短信（云厂商未接 SDK 时打桩记录，行为与 hei-boot 配置契约一致）。
func (f *Facade) SendSMS(ctx context.Context, phone, content string) error {
	_ = ctx
	engine := strings.ToUpper(strings.TrimSpace(f.runtime(ctx, "DEFAULT_SMS_ENGINE", "ALIYUN")))
	if !f.cfg.SMS.Enabled {
		logger.L.Info("sms disabled, skip send",
			zap.String("phone", phone), zap.String("content", content))
		return nil
	}
	logger.L.Info("sms send (provider stub)",
		zap.String("engine", engine),
		zap.String("phone", phone),
		zap.String("content", content))
	return nil
}

// SendTemplated 按场景发送（邮件/短信模板优先运行时 MAIL_TEMPLATE_* / SMS_TEMPLATE_*，回退内置）。
func (f *Facade) SendTemplated(ctx context.Context, template, to string, vars map[string]any) error {
	channel := strings.ToUpper(strings.TrimSpace(template))
	isMail := strings.Contains(to, "@") || channel == "LOGIN_CODE_MAIL" || channel == "RESET_PASSWORD_CODE" ||
		strings.HasSuffix(channel, "_MAIL")
	if isMail {
		subject, body := f.renderMail(ctx, template, vars)
		return f.SendMail(ctx, to, subject, body)
	}
	body := f.renderSMS(ctx, template, vars)
	return f.SendSMS(ctx, to, body)
}

// renderMail 邮件模板：MAIL_TEMPLATE_{SCENE} JSON {subject, body}，{{var}} 替换，回退内置。
func (f *Facade) renderMail(ctx context.Context, template string, vars map[string]any) (string, string) {
	scene := normalizeScene(template)
	raw := f.runtime(ctx, "MAIL_TEMPLATE_"+scene, "")
	if raw != "" {
		var t struct {
			Subject string `json:"subject"`
			Body    string `json:"body"`
		}
		if err := json.Unmarshal([]byte(raw), &t); err == nil && (t.Subject != "" || t.Body != "") {
			return renderVars(t.Subject, vars), renderVars(t.Body, vars)
		}
	}
	return renderBuiltin(template, vars, true)
}

// renderSMS 短信模板：SMS_TEMPLATE_{SCENE} JSON {content}，{{var}} 替换，回退内置。
func (f *Facade) renderSMS(ctx context.Context, template string, vars map[string]any) string {
	scene := normalizeScene(template)
	raw := f.runtime(ctx, "SMS_TEMPLATE_"+scene, "")
	if raw != "" {
		var t struct {
			Code    string `json:"code"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal([]byte(raw), &t); err == nil && t.Content != "" {
			return renderVars(t.Content, vars)
		}
	}
	_, body := renderBuiltin(template, vars, false)
	return body
}

func normalizeScene(template string) string {
	s := strings.ToUpper(strings.TrimSpace(template))
	s = strings.TrimSuffix(s, "_MAIL")
	return s
}

// renderVars 将 {{key}} 占位符替换为变量值。
func renderVars(tpl string, vars map[string]any) string {
	if !strings.Contains(tpl, "{{") {
		return tpl
	}
	var b strings.Builder
	b.Grow(len(tpl))
	rest := tpl
	for {
		start := strings.Index(rest, "{{")
		if start < 0 {
			b.WriteString(rest)
			break
		}
		b.WriteString(rest[:start])
		end := strings.Index(rest[start:], "}}")
		if end < 0 {
			b.WriteString(rest[start:])
			break
		}
		key := strings.TrimSpace(rest[start+5 : start+end])
		val := ""
		if v, ok := vars[key]; ok && v != nil {
			val = fmt.Sprint(v)
		}
		b.WriteString(val)
		rest = rest[start+end+6:]
	}
	return b.String()
}

func renderBuiltin(template string, vars map[string]any, mail bool) (subject, body string) {
	get := func(k, def string) string {
		if vars == nil {
			return def
		}
		if v, ok := vars[k]; ok && v != nil {
			return fmt.Sprint(v)
		}
		return def
	}
	app := get("app_name", "HEI")
	switch normalizeScene(template) {
	case "LOGIN_CODE", "LOGIN_CODE_MAIL":
		return app + " 登录验证码", fmt.Sprintf("您的登录验证码为 %s，%s 分钟内有效。", get("code", ""), get("expire_minutes", "5"))
	case "RESET_PASSWORD_CODE":
		link := get("reset_link", "")
		if link != "" {
			return app + " 密码重置", fmt.Sprintf("请点击链接重置密码（%s 分钟内有效）：\n%s", get("expire_minutes", "10"), link)
		}
		return app + " 密码重置", fmt.Sprintf("您的密码重置令牌为 %s，%s 分钟内有效。", get("token", get("code", "")), get("expire_minutes", "10"))
	case "BIND_PHONE_CODE":
		return app + " 绑定手机验证码", fmt.Sprintf("您的绑定验证码为 %s，%s 分钟内有效。", get("code", ""), get("expire_minutes", "5"))
	case "BIND_EMAIL_CODE":
		return app + " 绑定邮箱验证码", fmt.Sprintf("您的绑定验证码为 %s，%s 分钟内有效。", get("code", ""), get("expire_minutes", "5"))
	case "CHANGE_PASSWORD_CODE", "CHANGE_EMAIL_CODE", "CHANGE_PHONE_CODE":
		return app + " 安全验证码", fmt.Sprintf("您的验证码为 %s，%s 分钟内有效。", get("code", ""), get("expire_minutes", "5"))
	default:
		if mail {
			return app + " 通知", fmt.Sprintf("%v", vars)
		}
		return "", fmt.Sprintf("%v", vars)
	}
}

// SendPush 通用 HTTP 推送（DingTalk / Lark / WeCom Webhook，运行时 PUSH_* 配置）。
func (f *Facade) SendPush(ctx context.Context, title, body string) error {
	_ = ctx
	engine := strings.ToUpper(strings.TrimSpace(f.runtime(ctx, "DEFAULT_MESSAGE_PUSH_ENGINE", "DINGTALK")))
	text := body
	if title != "" {
		text = title + "\n" + body
	}
	switch engine {
	case "DINGTALK":
		return f.sendDingtalk(ctx, text)
	case "LARK", "FEISHU":
		return f.sendLark(ctx, text)
	case "WECOM", "WECHAT_WORK", "WECHATWORK":
		return f.sendWecom(ctx, text)
	default:
		logger.L.Info("push engine unsupported (stub)", zap.String("engine", engine))
		return nil
	}
}

func (f *Facade) sendDingtalk(ctx context.Context, text string) error {
	webhook := strings.TrimSpace(f.runtime(ctx, "PUSH_DINGTALK_WEBHOOK", ""))
	if webhook == "" {
		logger.L.Info("push not configured: PUSH_DINGTALK_WEBHOOK empty, skip")
		return nil
	}
	secret := strings.TrimSpace(f.runtime(ctx, "PUSH_DINGTALK_SECRET", ""))
	u := webhook
	if secret != "" {
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		sign := signDingtalk(timestamp, secret)
		sep := "?"
		if strings.Contains(u, "?") {
			sep = "&"
		}
		u = u + sep + "timestamp=" + timestamp + "&sign=" + url.QueryEscape(sign)
	}
	payload := map[string]any{"msgtype": "text", "text": map[string]any{"content": text}}
	return f.postJSON(ctx, u, payload, "钉钉")
}

func (f *Facade) sendLark(ctx context.Context, text string) error {
	webhook := strings.TrimSpace(f.runtime(ctx, "PUSH_LARK_WEBHOOK", ""))
	if webhook == "" {
		logger.L.Info("push not configured: PUSH_LARK_WEBHOOK empty, skip")
		return nil
	}
	payload := map[string]any{"msg_type": "text", "content": map[string]any{"text": text}}
	secret := strings.TrimSpace(f.runtime(ctx, "PUSH_LARK_SECRET", ""))
	if secret != "" {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		payload["timestamp"] = timestamp
		payload["sign"] = signFeishu(timestamp, secret)
	}
	return f.postJSON(ctx, webhook, payload, "飞书")
}

func (f *Facade) sendWecom(ctx context.Context, text string) error {
	webhook := strings.TrimSpace(f.runtime(ctx, "PUSH_WECHAT_WORK_WEBHOOK", ""))
	if webhook == "" {
		logger.L.Info("push not configured: PUSH_WECHAT_WORK_WEBHOOK empty, skip")
		return nil
	}
	payload := map[string]any{"msgtype": "text", "text": map[string]any{"content": text}}
	return f.postJSON(ctx, webhook, payload, "企业微信")
}

// SendWebhook 自定义 Webhook 发送（审计告警测试/通知用；URL 为空返回错误）。
func (f *Facade) SendWebhook(ctx context.Context, webhookURL, secret, payload string) error {
	if strings.TrimSpace(webhookURL) == "" {
		return fmt.Errorf("notify: empty webhook url")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if secret != "" {
		req.Header.Set("X-HEI-Signature", secret)
	}
	resp, err := f.hc.Do(req)
	if err != nil {
		logger.L.Warn("webhook send failed", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook send failed: HTTP %d", resp.StatusCode)
	}
	return nil
}

func (f *Facade) postJSON(ctx context.Context, u string, payload map[string]any, label string) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := f.hc.Do(req)
	if err != nil {
		logger.L.Warn(label+" push failed", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf(label+" push failed: HTTP %d", resp.StatusCode)
	}
	return nil
}

func signDingtalk(timestamp, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestamp + "\n" + secret))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func signFeishu(timestamp, secret string) string {
	mac := hmac.New(sha256.New, []byte(timestamp+"\n"+secret))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func sendMailTLS(addr, host string, auth smtp.Auth, from string, to []string, msg []byte) error {
	hostOnly, _, err := net.SplitHostPort(addr)
	if err != nil {
		hostOnly = host
	}
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: hostOnly})
	if err != nil {
		return err
	}
	defer conn.Close()
	c, err := smtp.NewClient(conn, hostOnly)
	if err != nil {
		return err
	}
	defer c.Close()
	if auth != nil {
		if err := c.Auth(auth); err != nil {
			return err
		}
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

func sendMailStartTLS(addr, host string, auth smtp.Auth, from string, to []string, msg []byte) error {
	hostOnly, _, err := net.SplitHostPort(addr)
	if err != nil {
		hostOnly = host
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	c, err := smtp.NewClient(conn, hostOnly)
	if err != nil {
		return err
	}
	defer c.Close()
	if err := c.StartTLS(&tls.Config{ServerName: hostOnly}); err != nil {
		return err
	}
	if auth != nil {
		if err := c.Auth(auth); err != nil {
			return err
		}
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}
