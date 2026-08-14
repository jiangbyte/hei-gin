package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/notify"
	"hei-gin/internal/modules/shared"
)

// AlertLog å¯¹åº” sys_alert_logã€‚
//
// Author: Charlie
type AlertLog struct {
	ID          string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	RuleName    string         `gorm:"column:rule_name;size:64;not null" json:"rule_name"`
	Severity    string         `gorm:"column:severity;size:16;not null" json:"severity"`
	Summary     string         `gorm:"column:summary;size:255;not null" json:"summary"`
	Details     datatypes.JSON `gorm:"column:details;type:jsonb" json:"details"`
	NotifiedVia *string        `gorm:"column:notified_via;size:64" json:"notified_via"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
}

// TableName è¿”å›žå‘Šè­¦æ—¥å¿—è¡¨åã€‚
func (AlertLog) TableName() string { return "sys_alert_log" }

// withJobs é™„åŠ  SnailJob handlersã€‚
func (s *Service) withJobs(m module.Module, nf *notify.Facade) module.Module {
	m.Jobs = append(m.Jobs, module.Job{
		Name: "auditAlertJob",
		Run: func(ctx context.Context, _ string) error {
			return s.runAuditAlertJob(ctx, nf)
		},
	})
	m.Models = append(m.Models, &AlertLog{})
	return m
}

func (s *Service) runAuditAlertJob(ctx context.Context, nf *notify.Facade) error {
	if !runtimeBool(ctx, s.repo.DB(), nf, "AUDIT_ALERT_ENABLED", true) {
		return nil
	}
	if !runtimeBool(ctx, s.repo.DB(), nf, "AUDIT_ALERT_RULE_BRUTE_FORCE", true) {
		return nil
	}
	windowSeconds := runtimeInt(ctx, s.repo.DB(), nf, "AUDIT_ALERT_ANALYSIS_INTERVAL_SECONDS", 60)
	if windowSeconds < 60 {
		windowSeconds = 60
	}
	threshold := int64(runtimeInt(ctx, s.repo.DB(), nf, "AUDIT_ALERT_BRUTE_FORCE_THRESHOLD", 10))
	if threshold < 1 {
		threshold = 1
	}
	since := time.Now().Add(-time.Duration(windowSeconds) * time.Second)
	var volume int64
	if err := s.repo.DB().WithContext(ctx).Model(&OperationLog{}).
		Where("created_at >= ?", since).Count(&volume).Error; err != nil {
		return err
	}
	if volume < threshold {
		return nil
	}
	cooldown := runtimeInt(ctx, s.repo.DB(), nf, "AUDIT_ALERT_ALERT_COOLDOWN_SECONDS", 1800)
	if cooldown < windowSeconds {
		cooldown = windowSeconds
	}
	cooldownSince := time.Now().Add(-time.Duration(cooldown) * time.Second)
	var recent int64
	if err := s.repo.DB().WithContext(ctx).Model(&AlertLog{}).
		Where("rule_name = ? AND created_at >= ?", "audit_volume", cooldownSince).
		Count(&recent).Error; err != nil {
		return err
	}
	if recent > 0 {
		return nil
	}
	summary := fmt.Sprintf("Audit log volume %d exceeded threshold %d in last %d seconds", volume, threshold, windowSeconds)
	notified := notifyAuditChannels(ctx, nf, s.repo.DB(), summary)
	details, _ := json.Marshal(map[string]any{
		"volume":         volume,
		"threshold":      threshold,
		"window_seconds": windowSeconds,
		"since":          since.Format(time.RFC3339),
	})
	via := strings.Join(notified, ",")
	if via == "" {
		via = "sys_alert_log"
	}
	row := AlertLog{
		ID:        idgen.Next(),
		RuleName:  "audit_volume",
		Severity:  "WARNING",
		Summary:   summary,
		Details:   datatypes.JSON(details),
		CreatedAt: time.Now().UTC(),
	}
	row.NotifiedVia = &via
	return s.repo.DB().WithContext(ctx).Create(&row).Error
}

func notifyAuditChannels(ctx context.Context, nf *notify.Facade, db *gorm.DB, summary string) []string {
	if nf == nil {
		return nil
	}
	var out []string
	if runtimeBool(ctx, db, nf, "AUDIT_ALERT_NOTIFY_PUSH", true) {
		if err := nf.SendPush(ctx, "å®¡è®¡å‘Šè­¦", summary); err == nil {
			out = append(out, "push")
		}
	}
	if runtimeBool(ctx, db, nf, "AUDIT_ALERT_NOTIFY_EMAIL", true) {
		to := runtimeString(ctx, db, nf, "AUDIT_ALERT_NOTIFY_EMAIL_TO", "")
		if to != "" {
			if err := nf.SendMail(ctx, to, "å®¡è®¡å‘Šè­¦", summary); err == nil {
				out = append(out, "email")
			}
		}
	}
	if runtimeBool(ctx, db, nf, "AUDIT_ALERT_NOTIFY_CUSTOM_WEBHOOK", false) {
		url := runtimeString(ctx, db, nf, "AUDIT_ALERT_WEBHOOK_URL", "")
		secret := runtimeString(ctx, db, nf, "AUDIT_ALERT_WEBHOOK_SECRET", "")
		if url != "" {
			payload := fmt.Sprintf(`{"title":"å®¡è®¡å‘Šè­¦","body":%q}`, summary)
			if err := nf.SendWebhook(ctx, url, secret, payload); err == nil {
				out = append(out, "webhook")
			}
		}
	}
	return out
}

func runtimeString(ctx context.Context, db *gorm.DB, nf *notify.Facade, key, def string) string {
	if nf != nil {
		if v := nf.GetRuntimeString(ctx, key, ""); v != "" {
			return v
		}
	}
	if db != nil {
		var v string
		if err := db.WithContext(ctx).Table("sys_config").Select("config_value").
			Where("config_key = ?", key).Limit(1).Scan(&v).Error; err == nil && v != "" {
			return v
		}
	}
	return def
}

func runtimeBool(ctx context.Context, db *gorm.DB, nf *notify.Facade, key string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(runtimeString(ctx, db, nf, key, "")))
	if v == "" {
		return def
	}
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func runtimeInt(ctx context.Context, db *gorm.DB, nf *notify.Facade, key string, def int) int {
	v := strings.TrimSpace(runtimeString(ctx, db, nf, key, ""))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// New æž„å»º sys.audit æ¨¡å—ï¼ˆå« auditAlertJobï¼‰ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	m := module.Module{
		Name:   "sys.audit",
		Models: []any{&OperationLog{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
	return s.withJobs(m, d.Notify)
}
