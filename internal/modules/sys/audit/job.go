// internal/modules/sys/audit/job.go 定时任务。
//
// Author: Charlie

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

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/notify"
)

// AlertLog 对应 sys_alert_log。
//
// Author: Charlie
type AlertLog struct {
	ID          string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	RuleName    string         `gorm:"column:rule_name;size:64;not null" json:"rule_name"`
	Severity    string         `gorm:"column:severity;size:16;not null" json:"severity"`
	Summary     string         `gorm:"column:summary;size:255;not null" json:"summary"`
	Details     datatypes.JSON `gorm:"column:details;type:json" json:"details"`
	NotifiedVia *string        `gorm:"column:notified_via;size:64" json:"notified_via"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
}

// TableName 返回告警日志表名。
func (AlertLog) TableName() string { return "sys_alert_log" }

// 告警规则名（与 hei-boot AuditAlertJob 对齐）。
const (
	ruleBruteForce   = "audit_volume"
	ruleUnusualHours = "unusual_hours"
	ruleSensitiveOps = "sensitive_ops"
	ruleBulkDelete   = "bulk_delete"
	ruleIPAnomaly    = "ip_anomaly"
)

// sensitiveActions 凌晨非常时段 / 短窗口敏感操作判定动作集。
var sensitiveActions = []string{"role_create", "role_grant", "permission_change", "permission_grant"}

// withJobs 附加任务处理器（gojob 调度器收集）。
func (s *Service) withJobs(m module.Module, nf *notify.Facade, cfg *config.Config) module.Module {
	m.Jobs = append(m.Jobs, module.Job{
		Name: "sys_audit_alert",
		Run: func(ctx context.Context, _ string) (string, error) {
			if err := s.runAuditAlertJob(ctx, nf); err != nil {
				return "", err
			}
			return "done", nil
		},
	})
	m.Jobs = append(m.Jobs, module.Job{
		Name: "sys_audit_log_cleanup",
		Run: func(ctx context.Context, paramJSON string) (string, error) {
			return s.runAuditLogCleanupJob(ctx, paramJSON, cfg)
		},
	})
	m.Models = append(m.Models, &AlertLog{})
	return m
}

// runAuditLogCleanupJob 按保留天数批量清理过期审计日志（对齐 hei-boot AuditLogCleanupJob）。
func (s *Service) runAuditLogCleanupJob(ctx context.Context, paramJSON string, cfg *config.Config) (string, error) {
	loginRetention := 180
	operationRetention := 365
	batchSize := 1000
	if cfg != nil {
		loginRetention = cfg.Audit.LoginRetentionDays
		operationRetention = cfg.Audit.OperationRetentionDays
		if cfg.Audit.CleanupBatchSize > 0 {
			batchSize = cfg.Audit.CleanupBatchSize
		}
	}
	if strings.TrimSpace(paramJSON) != "" && paramJSON != "null" {
		var m map[string]any
		if err := json.Unmarshal([]byte(paramJSON), &m); err == nil {
			if v, ok := asInt(m["loginRetentionDays"]); ok {
				loginRetention = v
			}
			if v, ok := asInt(m["operationRetentionDays"]); ok {
				operationRetention = v
			}
			if v, ok := asInt(m["batchSize"]); ok && v > 0 {
				batchSize = v
			}
		}
	}

	deletedLogin := 0
	if loginRetention > 0 {
		n, err := s.repo.CleanupExpiredLoginLogs(ctx, loginRetention, batchSize)
		if err != nil {
			return "", err
		}
		deletedLogin = n
	}

	deletedOperation := 0
	if operationRetention > 0 {
		n, err := s.repo.CleanupExpiredOperationLogs(ctx, operationRetention, batchSize)
		if err != nil {
			return "", err
		}
		deletedOperation = n
	}

	return fmt.Sprintf(
		"deletedLogin=%d,deletedOperation=%d,loginRetentionDays=%d,operationRetentionDays=%d,batchSize=%d",
		deletedLogin, deletedOperation, loginRetention, operationRetention, batchSize,
	), nil
}

func asInt(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	case json.Number:
		i, err := n.Int64()
		return int(i), err == nil
	default:
		return 0, false
	}
}

// runAuditAlertJob 按配置规则扫描并发送告警（对齐 hei-boot AuditAlertJob 五项规则）。
func (s *Service) runAuditAlertJob(ctx context.Context, nf *notify.Facade) error {
	if !runtimeBool(ctx, s.repo.DB(), nf, "AUDIT_ALERT_ENABLED", true) {
		return nil
	}
	fired := 0
	if runtimeBool(ctx, s.repo.DB(), nf, "AUDIT_ALERT_RULE_BRUTE_FORCE", true) {
		if s.evaluateBruteForce(ctx, nf) {
			fired++
		}
	}
	if runtimeBool(ctx, s.repo.DB(), nf, "AUDIT_ALERT_RULE_UNUSUAL_HOURS", true) {
		if s.evaluateUnusualHours(ctx, nf) {
			fired++
		}
	}
	if runtimeBool(ctx, s.repo.DB(), nf, "AUDIT_ALERT_RULE_SENSITIVE_OPS", true) {
		if s.evaluateSensitiveOps(ctx, nf) {
			fired++
		}
	}
	if runtimeBool(ctx, s.repo.DB(), nf, "AUDIT_ALERT_RULE_BULK_DELETE", true) {
		if s.evaluateBulkDelete(ctx, nf) {
			fired++
		}
	}
	if runtimeBool(ctx, s.repo.DB(), nf, "AUDIT_ALERT_RULE_IP_ANOMALY", true) {
		if s.evaluateIPAnomaly(ctx, nf) {
			fired++
		}
	}
	return nil
}

// evaluateBruteForce 分析窗口内审计日志总量超过阈值则告警。
func (s *Service) evaluateBruteForce(ctx context.Context, nf *notify.Facade) bool {
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
		return false
	}
	if volume < threshold {
		return false
	}
	summary := fmt.Sprintf("Audit log volume %d exceeded threshold %d in last %d seconds", volume, threshold, windowSeconds)
	details := map[string]any{
		"volume": volume, "threshold": threshold,
		"window_seconds": windowSeconds, "window_minutes": max(1, windowSeconds/60),
		"since": since.Format(time.RFC3339),
	}
	cooldown := runtimeInt(ctx, s.repo.DB(), nf, "AUDIT_ALERT_ALERT_COOLDOWN_SECONDS", 1800)
	if cooldown < windowSeconds {
		cooldown = windowSeconds
	}
	return s.fireAlert(ctx, nf, ruleBruteForce, "WARNING", summary, details, cooldown)
}

// evaluateUnusualHours 凌晨 0-6 点出现角色/权限变更等敏感操作。
func (s *Service) evaluateUnusualHours(ctx context.Context, nf *notify.Facade) bool {
	now := time.Now()
	if now.Hour() > 5 {
		return false
	}
	since := now.Add(-time.Hour)
	var logs []OperationLog
	if err := s.repo.DB().WithContext(ctx).Where("created_at >= ? AND action IN ?", since, sensitiveActions).
		Find(&logs).Error; err != nil {
		return false
	}
	if len(logs) == 0 {
		return false
	}
	seen := map[string]struct{}{}
	actions := []string{}
	for _, l := range logs {
		if l.Action == "" {
			continue
		}
		if _, ok := seen[l.Action]; ok {
			continue
		}
		seen[l.Action] = struct{}{}
		actions = append(actions, l.Action)
	}
	summary := fmt.Sprintf("凌晨 %d 时检测到 %d 次敏感操作", now.Hour(), len(logs))
	return s.fireAlert(ctx, nf, ruleUnusualHours, "WARNING", summary,
		map[string]any{"count": len(logs), "actions": actions},
		alertCooldown(ctx, nf))
}

// evaluateSensitiveOps 5 分钟内角色授权/权限变更等敏感操作，按账户聚合告警。
func (s *Service) evaluateSensitiveOps(ctx context.Context, nf *notify.Facade) bool {
	since := time.Now().Add(-5 * time.Minute)
	rows := s.countByAccount(ctx, "created_at >= ? AND action IN ?", since, []string{"role_grant", "permission_change", "permission_grant"})
	fired := false
	for _, r := range rows {
		summary := fmt.Sprintf("账户 %s 执行了敏感操作 (%d 次)", r.AccountID, r.Cnt)
		fired = s.fireAlert(ctx, nf, ruleSensitiveOps, "WARNING", summary,
			map[string]any{"account_id": r.AccountID, "count": r.Cnt},
			alertCooldown(ctx, nf)) || fired
	}
	return fired
}

// evaluateBulkDelete 同账户 5 分钟内删除操作达到阈值。
func (s *Service) evaluateBulkDelete(ctx context.Context, nf *notify.Facade) bool {
	threshold := int64(runtimeInt(ctx, s.repo.DB(), nf, "AUDIT_ALERT_BULK_DELETE_THRESHOLD", 20))
	if threshold < 1 {
		threshold = 1
	}
	since := time.Now().Add(-5 * time.Minute)
	rows := s.countByAccount(ctx, "created_at >= ? AND action = ?", since, "delete")
	fired := false
	for _, r := range rows {
		if r.Cnt < threshold {
			continue
		}
		summary := fmt.Sprintf("账户 %s 在 5 分钟内删除了 %d 条记录", r.AccountID, r.Cnt)
		fired = s.fireAlert(ctx, nf, ruleBulkDelete, "WARNING", summary,
			map[string]any{"account_id": r.AccountID, "count": r.Cnt, "threshold": threshold},
			alertCooldown(ctx, nf)) || fired
	}
	return fired
}

// evaluateIPAnomaly 同账户 15 分钟内从多个不同 IP 成功登录达到阈值。
func (s *Service) evaluateIPAnomaly(ctx context.Context, nf *notify.Facade) bool {
	threshold := runtimeInt(ctx, s.repo.DB(), nf, "AUDIT_ALERT_IP_ANOMALY_THRESHOLD", 3)
	if threshold < 1 {
		threshold = 1
	}
	since := time.Now().Add(-15 * time.Minute)
	var rows []struct {
		AccountID string `gorm:"column:account_id"`
		IP        string `gorm:"column:ip"`
	}
	if err := s.repo.DB().WithContext(ctx).Table("sys_operation_audit_log").
		Select("account_id", "ip").
		Where("created_at >= ? AND action = ? AND success = ? AND account_id IS NOT NULL", since, "login", true).
		Find(&rows).Error; err != nil {
		return false
	}
	ipsByAccount := map[string]map[string]struct{}{}
	for _, r := range rows {
		acc := r.AccountID
		if acc == "" {
			continue
		}
		if ipsByAccount[acc] == nil {
			ipsByAccount[acc] = map[string]struct{}{}
		}
		ipsByAccount[acc][r.IP] = struct{}{}
	}
	fired := false
	for acc, ips := range ipsByAccount {
		if len(ips) < threshold {
			continue
		}
		summary := fmt.Sprintf("账户 %s 在 15 分钟内从 %d 个不同 IP 登录", acc, len(ips))
		fired = s.fireAlert(ctx, nf, ruleIPAnomaly, "WARNING", summary,
			map[string]any{"account_id": acc, "ip_count": len(ips), "threshold": threshold},
			alertCooldown(ctx, nf)) || fired
	}
	return fired
}

type accountCount struct {
	AccountID string `gorm:"column:account_id"`
	Cnt       int64  `gorm:"column:cnt"`
}

// countByAccount 按账户聚合审计计数（WHERE 条件为预编译片段）。
func (s *Service) countByAccount(ctx context.Context, where string, args ...any) []accountCount {
	var rows []accountCount
	q := s.repo.DB().WithContext(ctx).Table("sys_operation_audit_log").
		Select("account_id, count(*) AS cnt").
		Where(where, args...).
		Where("account_id IS NOT NULL AND account_id <> ''").
		Group("account_id")
	if err := q.Scan(&rows).Error; err != nil {
		return nil
	}
	return rows
}

// alertCooldown 告警冷却（秒，至少 60）。
func alertCooldown(ctx context.Context, nf *notify.Facade) int {
	c := runtimeInt(ctx, nil, nf, "AUDIT_ALERT_ALERT_COOLDOWN_SECONDS", 1800)
	if c < 60 {
		c = 60
	}
	return c
}

// fireAlert 公共告警出口：冷却期抑制 → 通知渠道 → 写入 sys_alert_log。
func (s *Service) fireAlert(ctx context.Context, nf *notify.Facade, ruleName, severity, summary string, details map[string]any, cooldownSeconds int) bool {
	cooldownSince := time.Now().Add(-time.Duration(cooldownSeconds) * time.Second)
	var recent int64
	if err := s.repo.DB().WithContext(ctx).Model(&AlertLog{}).
		Where("rule_name = ? AND created_at >= ?", ruleName, cooldownSince).
		Count(&recent).Error; err != nil {
		return false
	}
	if recent > 0 {
		return false
	}
	notified := notifyAuditChannels(ctx, nf, s.repo.DB(), summary)
	detailsJSON, _ := json.Marshal(details)
	via := strings.Join(notified, ",")
	if via == "" {
		via = "sys_alert_log"
	}
	row := AlertLog{
		ID:        idgen.Next(),
		RuleName:  ruleName,
		Severity:  severity,
		Summary:   summary,
		Details:   datatypes.JSON(detailsJSON),
		CreatedAt: time.Now().UTC(),
	}
	row.NotifiedVia = &via
	if err := s.repo.DB().WithContext(ctx).Create(&row).Error; err != nil {
		return false
	}
	return true
}

func notifyAuditChannels(ctx context.Context, nf *notify.Facade, db *gorm.DB, summary string) []string {
	if nf == nil {
		return nil
	}
	var out []string
	if runtimeBool(ctx, db, nf, "AUDIT_ALERT_NOTIFY_PUSH", true) {
		if err := nf.SendPush(ctx, "审计告警", summary); err == nil {
			out = append(out, "push")
		}
	}
	if runtimeBool(ctx, db, nf, "AUDIT_ALERT_NOTIFY_EMAIL", true) {
		to := runtimeString(ctx, db, nf, "AUDIT_ALERT_NOTIFY_EMAIL_TO", "")
		if to != "" {
			if err := nf.SendMail(ctx, to, "审计告警", summary); err == nil {
				out = append(out, "email")
			}
		}
	}
	if runtimeBool(ctx, db, nf, "AUDIT_ALERT_NOTIFY_CUSTOM_WEBHOOK", false) {
		url := runtimeString(ctx, db, nf, "AUDIT_ALERT_WEBHOOK_URL", "")
		secret := runtimeString(ctx, db, nf, "AUDIT_ALERT_WEBHOOK_SECRET", "")
		if url != "" {
			payload := fmt.Sprintf(`{"title":"审计告警","body":%q}`, summary)
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
