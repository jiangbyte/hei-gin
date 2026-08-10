package cg_test_activity

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 活动服务。
//
// Author: Charlie
type Service struct {
	repo *Repo
}

// NewService 构造服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 biz.cg_test_activity 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "biz.cg_test_activity",
		Order:  90,
		Models: []any{&Activity{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建活动。
func (s *Service) Create(ctx context.Context, accountID string, req AddParam) error {
	row := fromAddParam(req)
	row.ID = idgen.Next()
	row.CreatedBy = &accountID
	row.UpdatedBy = &accountID
	return s.repo.Create(ctx, &row)
}

// Update 更新活动。
func (s *Service) Update(ctx context.Context, accountID string, req EditParam) error {
	row := fromAddParam(req.AddParam)
	return s.repo.Update(ctx, req.ID, map[string]any{
		"code": row.Code, "name": row.Name, "category": row.Category, "type": row.Type,
		"status": row.Status, "cover_url": row.CoverURL, "description": row.Description,
		"start_at": row.StartAt, "end_at": row.EndAt, "max_participants": row.MaxParticipants,
		"price": row.Price, "is_public": row.IsPublic, "need_approval": row.NeedApproval,
		"rule_config": row.RuleConfig, "extra": row.Extra, "updated_by": accountID,
	})
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 活动详情。
func (s *Service) Detail(ctx context.Context, id string) (*Activity, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, p PageParam) (rows []Activity, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p)
	return rows, total, current, size, err
}

func fromAddParam(req AddParam) Activity {
	return Activity{
		Code: req.Code, Name: req.Name, Category: req.Category, Type: req.Type, Status: req.Status,
		CoverURL: req.CoverURL, Description: req.Description, StartAt: req.StartAt, EndAt: req.EndAt,
		MaxParticipants: req.MaxParticipants, Price: req.Price, IsPublic: req.IsPublic,
		NeedApproval: req.NeedApproval, RuleConfig: mustJSON(req.RuleConfig), Extra: mustJSON(req.Extra),
	}
}

func mustJSON(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte("{}"))
	}
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return b
}
