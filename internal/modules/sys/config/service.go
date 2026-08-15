// internal/modules/sys/config/service.go 业务服务。
//
// Author: Charlie

package config

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/crypto"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/notify"
	"hei-gin/internal/modules/shared"
)

// 敏感配置键（与 FastAPI SENSITIVE_CONFIG_KEYS 对齐的子集）。
var sensitiveKeys = map[string]struct{}{
	"AUTH_DEFAULT_PASSWORD":         {},
	"AUDIT_ALERT_WEBHOOK_SECRET":    {},
	"MAIL_LOCAL_PASSWORD":           {},
	"MAIL_ALIYUN_ACCESS_KEY_SECRET": {},
	"MAIL_TENCENT_SECRET_KEY":       {},
	"SMS_ALIYUN_ACCESS_KEY_SECRET":  {},
	"SMS_TENCENT_SECRET_KEY":        {},
	"PUSH_DINGTALK_SECRET":          {},
	"PUSH_LARK_SECRET":              {},
	"STORAGE_MINIO_ACCESS_KEY":      {},
	"STORAGE_MINIO_SECRET_KEY":      {},
	"STORAGE_ALIYUN_ACCESS_KEY":     {},
	"STORAGE_ALIYUN_SECRET_KEY":     {},
	"STORAGE_TENCENT_ACCESS_KEY":    {},
	"STORAGE_TENCENT_SECRET_KEY":    {},
}

// Service 系统配置业务服务。
//
// Author: Charlie
type Service struct {
	repo   *Repo
	fernet *crypto.Codec
	notify *notify.Facade
}

// NewService 构造配置服务。
func NewService(db *gorm.DB, fernet *crypto.Codec, nf *notify.Facade) *Service {
	return &Service{repo: NewRepo(db), fernet: fernet, notify: nf}
}

// New 构建 sys.config 模块。
func New(d *shared.Deps) module.Module {
	var codec *crypto.Codec
	if d.Cfg != nil {
		codec, _ = crypto.NewFernetFromConfig(d.Cfg.Crypto.FernetKey, d.Cfg.Crypto.VaultAddr)
	}
	s := NewService(d.DB, codec, d.Notify)
	return module.Module{
		Name:   "sys.config",
		Models: []any{&Config{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

func pickRemark(description, remark *string) *string {
	if remark != nil {
		return remark
	}
	return description
}

func (s *Service) maybeEncrypt(key string, category *string, value *string) *string {
	if value == nil || *value == "" || s.fernet == nil {
		return value
	}
	if crypto.LooksEncrypted(*value) {
		return value
	}
	if !isSensitive(key, category) {
		return value
	}
	enc, err := s.fernet.Encrypt(*value)
	if err != nil {
		return value
	}
	return &enc
}

func (s *Service) maybeDecrypt(key string, category *string, value *string) *string {
	if value == nil || *value == "" || s.fernet == nil {
		return value
	}
	if !crypto.LooksEncrypted(*value) && !isSensitive(key, category) {
		return value
	}
	plain, err := s.fernet.Decrypt(*value)
	if err != nil {
		return value
	}
	return &plain
}

func isSensitive(key string, category *string) bool {
	if _, ok := sensitiveKeys[key]; ok {
		return true
	}
	if category == nil {
		return false
	}
	c := strings.ToUpper(strings.TrimSpace(*category))
	return c == "SECRET" || c == "SENSITIVE" || c == "CREDENTIAL"
}

func (s *Service) reveal(row *Config) {
	if row == nil {
		return
	}
	row.ConfigValue = s.maybeDecrypt(row.ConfigKey, row.Category, row.ConfigValue)
}

// Create 创建配置。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	vt := req.ValueType
	if vt == "" {
		vt = "STRING"
	}
	val := s.maybeEncrypt(req.ConfigKey, req.Category, req.ConfigValue)
	row := Config{
		ID: idgen.Next(), ConfigKey: req.ConfigKey, ConfigValue: val, Category: req.Category,
		Remark: req.Remark, SortCode: req.SortCode, ValueType: vt, Label: req.Label, Scope: req.Scope,
		Scene: req.Scene, ExtJSON: datatypes.JSON([]byte("{}")),
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新配置。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	vt := req.ValueType
	if vt == "" {
		vt = "STRING"
	}
	val := s.maybeEncrypt(req.ConfigKey, req.Category, req.ConfigValue)
	updates := map[string]any{
		"config_key": req.ConfigKey, "config_value": val, "category": req.Category,
		"remark": req.Remark, "sort_code": req.SortCode, "value_type": vt, "label": req.Label,
		"scope": req.Scope, "scene": req.Scene,
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 详情（敏感值解密后返回）。
func (s *Service) Detail(ctx context.Context, id string) (*Config, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.reveal(row)
	return row, nil
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []Config, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	for i := range rows {
		s.reveal(&rows[i])
	}
	return rows, total, current, size, err
}

// List 列表。
func (s *Service) List(ctx context.Context, q ListParam) ([]Config, error) {
	rows, err := s.repo.List(ctx, q)
	for i := range rows {
		s.reveal(&rows[i])
	}
	return rows, err
}

// BatchSave 批量保存配置：按 config_key 存在则更新，不存在则新建。
func (s *Service) BatchSave(ctx context.Context, req BatchSaveParam) error {
	items := make([]BatchItemParam, 0, len(req.Items))
	for _, it := range req.Items {
		if isSensitive(it.ConfigKey, it.Category) && (it.ConfigValue == nil || *it.ConfigValue == "") {
			// 敏感配置空值跳过，避免覆盖已有密文
			continue
		}
		items = append(items, it)
	}
	if len(items) == 0 {
		return nil
	}
	keys := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, it := range items {
		if _, ok := seen[it.ConfigKey]; ok {
			continue
		}
		seen[it.ConfigKey] = struct{}{}
		keys = append(keys, it.ConfigKey)
	}
	rows, err := s.repo.ListByKeys(ctx, keys)
	if err != nil {
		return err
	}
	existing := make(map[string]Config, len(rows))
	for _, row := range rows {
		existing[row.ConfigKey] = row
	}
	for _, it := range items {
		row, ok := existing[it.ConfigKey]
		// 敏感键落库后以 ext_json.is_set=true 标记「已配置」（对齐 web loadByCategory 契约）。
		ext := datatypes.JSON([]byte("{}"))
		if isSensitive(it.ConfigKey, it.Category) && it.ConfigValue != nil && *it.ConfigValue != "" {
			ext, _ = json.Marshal(map[string]bool{"is_set": true})
		}
		if !ok {
			nr := Config{
				ID:          idgen.Next(),
				ConfigKey:   it.ConfigKey,
				ConfigValue: s.maybeEncrypt(it.ConfigKey, it.Category, it.ConfigValue),
				Category:    it.Category,
				Remark:      pickRemark(it.Description, it.Remark),
				ValueType:   "STRING",
				ExtJSON:     ext,
			}
			if err := s.repo.Create(ctx, &nr); err != nil {
				return err
			}
			continue
		}
		updates := map[string]any{
			"config_value": s.maybeEncrypt(it.ConfigKey, it.Category, it.ConfigValue),
		}
		if string(ext) != "{}" {
			updates["ext_json"] = ext
		}
		if it.Category != nil {
			updates["category"] = it.Category
		}
		if it.Description != nil {
			updates["remark"] = pickRemark(it.Description, it.Remark)
		}
		if err := s.repo.Update(ctx, row.ID, updates); err != nil {
			return err
		}
	}
	return nil
}

// TestWebhook 发送审计告警测试 Webhook。
func (s *Service) TestWebhook(ctx context.Context, url, secret string) error {
	if s.notify == nil {
		return errors.New("notify facade unavailable")
	}
	payload := `{"msg_type":"text","content":{"text":"审计告警系统测试消息"}}`
	return s.notify.SendWebhook(ctx, url, secret, payload)
}

// TestPush 发送审计告警测试推送。
func (s *Service) TestPush(ctx context.Context) error {
	if s.notify == nil {
		return errors.New("notify facade unavailable")
	}
	return s.notify.SendPush(ctx, "审计告警", "审计告警系统测试消息")
}
