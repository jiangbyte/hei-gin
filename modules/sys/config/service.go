package config

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 系统配置业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造配置服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 sys.config 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.config",
		Models: []any{&Config{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建配置。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	vt := req.ValueType
	if vt == "" {
		vt = "STRING"
	}
	row := Config{
		ID: idgen.Next(), ConfigKey: req.ConfigKey, ConfigValue: req.ConfigValue, Category: req.Category,
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
	updates := map[string]any{
		"config_key": req.ConfigKey, "config_value": req.ConfigValue, "category": req.Category,
		"remark": req.Remark, "sort_code": req.SortCode, "value_type": vt, "label": req.Label,
		"scope": req.Scope, "scene": req.Scene,
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*Config, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []Config, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// List 列表。
func (s *Service) List(ctx context.Context, q ListParam) ([]Config, error) {
	return s.repo.List(ctx, q)
}
