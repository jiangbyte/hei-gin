// internal/modules/sys/config/service.go 业务服务。
//
// Author: Charlie

package config

import (
	"context"
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

// æ•æ„Ÿé…ç½®é”®ï¼ˆä¸Ž FastAPI SENSITIVE_CONFIG_KEYS å¯¹é½çš„å­é›†ï¼‰ã€‚
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

// Service ç³»ç»Ÿé…ç½®ä¸šåŠ¡æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct {
	repo   *Repo
	fernet *crypto.Codec
	notify *notify.Facade
}

// NewService æž„é€ é…ç½®æœåŠ¡ã€‚
func NewService(db *gorm.DB, fernet *crypto.Codec, nf *notify.Facade) *Service {
	return &Service{repo: NewRepo(db), fernet: fernet, notify: nf}
}

// New æž„å»º sys.config æ¨¡å—ã€‚
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

// Create åˆ›å»ºé…ç½®ã€‚
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

// Update æ›´æ–°é…ç½®ã€‚
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

// Delete æ‰¹é‡åˆ é™¤ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail è¯¦æƒ…ï¼ˆæ•æ„Ÿå€¼è§£å¯†åŽè¿”å›žï¼‰ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*Config, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.reveal(row)
	return row, nil
}

// Page åˆ†é¡µã€‚
func (s *Service) Page(ctx context.Context, q PageParam) (rows []Config, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	for i := range rows {
		s.reveal(&rows[i])
	}
	return rows, total, current, size, err
}

// List åˆ—è¡¨ã€‚
func (s *Service) List(ctx context.Context, q ListParam) ([]Config, error) {
	rows, err := s.repo.List(ctx, q)
	for i := range rows {
		s.reveal(&rows[i])
	}
	return rows, err
}

// BatchSave æ‰¹é‡ä¿å­˜é…ç½®ï¼šæŒ‰ config_key å­˜åœ¨åˆ™æ›´æ–°ï¼Œä¸å­˜åœ¨åˆ™æ–°å»ºã€‚
func (s *Service) BatchSave(ctx context.Context, req BatchSaveParam) error {
	items := make([]BatchItemParam, 0, len(req.Items))
	for _, it := range req.Items {
		if isSensitive(it.ConfigKey, it.Category) && (it.ConfigValue == nil || *it.ConfigValue == "") {
			// æ•æ„Ÿé…ç½®ç©ºå€¼è·³è¿‡ï¼Œé¿å…è¦†ç›–å·²æœ‰å¯†æ–‡
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
		if !ok {
			nr := Config{
				ID:          idgen.Next(),
				ConfigKey:   it.ConfigKey,
				ConfigValue: s.maybeEncrypt(it.ConfigKey, it.Category, it.ConfigValue),
				Category:    it.Category,
				Remark:      pickRemark(it.Description, it.Remark),
				ValueType:   "STRING",
				ExtJSON:     datatypes.JSON([]byte("{}")),
			}
			if err := s.repo.Create(ctx, &nr); err != nil {
				return err
			}
			continue
		}
		updates := map[string]any{
			"config_value": s.maybeEncrypt(it.ConfigKey, it.Category, it.ConfigValue),
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

// TestWebhook å‘é€å®¡è®¡å‘Šè­¦æµ‹è¯• Webhookã€‚
func (s *Service) TestWebhook(ctx context.Context, url, secret string) error {
	if s.notify == nil {
		return errors.New("notify facade unavailable")
	}
	payload := `{"msg_type":"text","content":{"text":"å®¡è®¡å‘Šè­¦ç³»ç»Ÿæµ‹è¯•æ¶ˆæ¯"}}`
	return s.notify.SendWebhook(ctx, url, secret, payload)
}

// TestPush å‘é€å®¡è®¡å‘Šè­¦æµ‹è¯•æŽ¨é€ã€‚
func (s *Service) TestPush(ctx context.Context) error {
	if s.notify == nil {
		return errors.New("notify facade unavailable")
	}
	return s.notify.SendPush(ctx, "å®¡è®¡å‘Šè­¦", "å®¡è®¡å‘Šè­¦ç³»ç»Ÿæµ‹è¯•æ¶ˆæ¯")
}
