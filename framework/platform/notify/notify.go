// Package notify 提供邮件 / 短信 / 推送通知门面。
package notify

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"hei-gin/framework/core/config"
	"hei-gin/framework/core/logger"
)

// Facade 通知发送门面。
//
// Author: Charlie
type Facade struct {
	cfg config.NotifyConfig
	db  *gorm.DB
}

// NewFacade 构造通知门面。
func NewFacade(cfg config.NotifyConfig, db *gorm.DB) *Facade {
	return &Facade{cfg: cfg, db: db}
}

// SendMail 发送纯文本邮件。
func (f *Facade) SendMail(ctx context.Context, to, subject, body string) error {
	_ = ctx
	if !f.cfg.Mail.Enabled {
		logger.L.Info("mail disabled, skip send",
			zap.String("to", to), zap.String("subject", subject))
		return nil
	}
	from := f.cfg.Mail.From
	if from == "" {
		from = f.cfg.Mail.Username
	}
	addr := fmt.Sprintf("%s:%d", f.cfg.Mail.Host, f.cfg.Mail.Port)
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body + "\r\n")
	auth := smtp.PlainAuth("", f.cfg.Mail.Username, f.cfg.Mail.Password, f.cfg.Mail.Host)
	if f.cfg.Mail.SSL {
		return sendMailTLS(addr, f.cfg.Mail.Host, auth, from, []string{to}, msg)
	}
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}

// SendSMS 发送短信（未启用时仅打日志）。
func (f *Facade) SendSMS(ctx context.Context, phone, content string) error {
	_ = ctx
	if !f.cfg.SMS.Enabled {
		logger.L.Info("sms disabled, skip send",
			zap.String("phone", phone), zap.String("content", content))
		return nil
	}
	logger.L.Info("sms send (provider stub)",
		zap.String("provider", f.cfg.SMS.Provider),
		zap.String("phone", phone),
		zap.String("content", content))
	return nil
}

// SendTemplated 按模板键发送（内置若干认证相关模板）。
func (f *Facade) SendTemplated(ctx context.Context, template, to string, vars map[string]any) error {
	subject, body := renderTemplate(template, vars)
	channel := strings.ToUpper(strings.TrimSpace(template))
	if strings.Contains(to, "@") || channel == "RESET_PASSWORD_CODE" || channel == "LOGIN_CODE_MAIL" {
		return f.SendMail(ctx, to, subject, body)
	}
	return f.SendSMS(ctx, to, body)
}

// SendPush 通用 HTTP 推送（未启用时打日志）。
func (f *Facade) SendPush(ctx context.Context, title, body string) error {
	_ = ctx
	if !f.cfg.Push.Enabled || f.cfg.Push.URL == "" {
		logger.L.Info("push disabled, skip send", zap.String("title", title), zap.String("body", body))
		return nil
	}
	logger.L.Info("push send (stub)", zap.String("url", f.cfg.Push.URL), zap.String("title", title))
	return nil
}

// SendWebhook 自定义 Webhook（带可选 secret 日志）。
func (f *Facade) SendWebhook(ctx context.Context, url, secret, payload string) error {
	_ = ctx
	_ = secret
	if url == "" {
		return fmt.Errorf("notify: empty webhook url")
	}
	logger.L.Info("webhook send (stub)", zap.String("url", url), zap.Int("payload_len", len(payload)))
	return nil
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

func renderTemplate(template string, vars map[string]any) (subject, body string) {
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
	switch strings.ToUpper(strings.TrimSpace(template)) {
	case "LOGIN_CODE", "LOGIN_CODE_MAIL":
		return app + " 登录验证码", fmt.Sprintf("您的登录验证码为 %s，%s 分钟内有效。", get("code", ""), get("expire_minutes", "5"))
	case "RESET_PASSWORD_CODE":
		link := get("reset_link", "")
		if link != "" {
			return app + " 重置密码", fmt.Sprintf("请点击链接重置密码（%s 分钟内有效）：\n%s", get("expire_minutes", "10"), link)
		}
		return app + " 重置密码", fmt.Sprintf("您的密码重置令牌为 %s，%s 分钟内有效。", get("token", get("code", "")), get("expire_minutes", "10"))
	default:
		return app + " 通知", fmt.Sprintf("%v", vars)
	}
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
