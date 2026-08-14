// Package notify æä¾›é‚®ä»¶ / çŸ­ä¿¡ / æŽ¨é€é€šçŸ¥é—¨é¢ã€‚
//
// Author: Charlie
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

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/logger"
)

// Facade é€šçŸ¥å‘é€é—¨é¢ã€‚
//
// Author: Charlie
type Facade struct {
	cfg config.NotifyConfig
	db  *gorm.DB
}

// NewFacade æž„é€ é€šçŸ¥é—¨é¢ã€‚
func NewFacade(cfg config.NotifyConfig, db *gorm.DB) *Facade {
	return &Facade{cfg: cfg, db: db}
}

// SendMail å‘é€çº¯æ–‡æœ¬é‚®ä»¶ã€‚
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

// SendSMS å‘é€çŸ­ä¿¡ï¼ˆæœªå¯ç”¨æ—¶ä»…æ‰“æ—¥å¿—ï¼‰ã€‚
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

// SendTemplated æŒ‰æ¨¡æ¿é”®å‘é€ï¼ˆå†…ç½®è‹¥å¹²è®¤è¯ç›¸å…³æ¨¡æ¿ï¼‰ã€‚
func (f *Facade) SendTemplated(ctx context.Context, template, to string, vars map[string]any) error {
	subject, body := renderTemplate(template, vars)
	channel := strings.ToUpper(strings.TrimSpace(template))
	if strings.Contains(to, "@") || channel == "RESET_PASSWORD_CODE" || channel == "LOGIN_CODE_MAIL" {
		return f.SendMail(ctx, to, subject, body)
	}
	return f.SendSMS(ctx, to, body)
}

// SendPush é€šç”¨ HTTP æŽ¨é€ï¼ˆæœªå¯ç”¨æ—¶æ‰“æ—¥å¿—ï¼‰ã€‚
func (f *Facade) SendPush(ctx context.Context, title, body string) error {
	_ = ctx
	if !f.cfg.Push.Enabled || f.cfg.Push.URL == "" {
		logger.L.Info("push disabled, skip send", zap.String("title", title), zap.String("body", body))
		return nil
	}
	logger.L.Info("push send (stub)", zap.String("url", f.cfg.Push.URL), zap.String("title", title))
	return nil
}

// SendWebhook è‡ªå®šä¹‰ Webhookï¼ˆå¸¦å¯é€‰ secret æ—¥å¿—ï¼‰ã€‚
func (f *Facade) SendWebhook(ctx context.Context, url, secret, payload string) error {
	_ = ctx
	_ = secret
	if url == "" {
		return fmt.Errorf("notify: empty webhook url")
	}
	logger.L.Info("webhook send (stub)", zap.String("url", url), zap.Int("payload_len", len(payload)))
	return nil
}

// GetRuntimeString ä»Ž sys_config è¯»å–è¿è¡Œæ—¶å­—ç¬¦ä¸²ï¼›æ— åˆ™è¿”å›ž defã€‚
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
		return app + " ç™»å½•éªŒè¯ç ", fmt.Sprintf("æ‚¨çš„ç™»å½•éªŒè¯ç ä¸º %sï¼Œ%s åˆ†é’Ÿå†…æœ‰æ•ˆã€‚", get("code", ""), get("expire_minutes", "5"))
	case "RESET_PASSWORD_CODE":
		link := get("reset_link", "")
		if link != "" {
			return app + " é‡ç½®å¯†ç ", fmt.Sprintf("è¯·ç‚¹å‡»é“¾æŽ¥é‡ç½®å¯†ç ï¼ˆ%s åˆ†é’Ÿå†…æœ‰æ•ˆï¼‰ï¼š\n%s", get("expire_minutes", "10"), link)
		}
		return app + " é‡ç½®å¯†ç ", fmt.Sprintf("æ‚¨çš„å¯†ç é‡ç½®ä»¤ç‰Œä¸º %sï¼Œ%s åˆ†é’Ÿå†…æœ‰æ•ˆã€‚", get("token", get("code", "")), get("expire_minutes", "10"))
	default:
		return app + " é€šçŸ¥", fmt.Sprintf("%v", vars)
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
