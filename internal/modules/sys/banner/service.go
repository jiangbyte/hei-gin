// internal/modules/sys/banner/service.go 业务服务。
//
// Author: Charlie

package banner

import (
	"context"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service Banner 业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造 Banner 服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 sys.banner 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.banner",
		Models: []any{&Banner{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
		Jobs: []module.Job{{
			Name: "bannerStatusJob",
			Run:  s.bannerStatusJobHandler,
		}},
	}
}

// Create 创建 Banner。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	lt := req.LinkType
	if lt == "" {
		lt = "URL"
	}
	targets := req.TargetAccountTypes
	if len(targets) == 0 {
		targets = datatypes.JSON([]byte("[]"))
	}
	row := Banner{
		ID: idgen.Next(), Title: req.Title, Image: req.Image, URL: req.URL, LinkType: lt,
		Summary: req.Summary, Description: req.Description, Category: req.Category, Type: req.Type,
		Position: req.Position, TargetAccountTypes: targets, Sort: req.Sort, Status: statusOr(req.Status),
		StartAt: req.StartAt, EndAt: req.EndAt,
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新 Banner。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	lt := req.LinkType
	if lt == "" {
		lt = "URL"
	}
	updates := map[string]any{
		"title": req.Title, "image": req.Image, "url": req.URL, "link_type": lt,
		"summary": req.Summary, "description": req.Description, "category": req.Category,
		"type": req.Type, "position": req.Position, "sort": req.Sort, "status": statusOr(req.Status),
		"start_at": req.StartAt, "end_at": req.EndAt,
	}
	if len(req.TargetAccountTypes) > 0 {
		updates["target_account_types"] = req.TargetAccountTypes
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*Banner, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []Banner, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// List 启用 Banner 列表。
func (s *Service) List(ctx context.Context, q ListParam) ([]Banner, error) {
	return s.repo.List(ctx, q.Position, security.StatusEnabled)
}

// Interaction 互动上报：找到 Banner 行并将互动计数 +1。
func (s *Service) Interaction(ctx context.Context, id string) error {
	n, err := s.repo.IncrementInteraction(ctx, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("banner not found")
	}
	return nil
}

// PortalList 门户端有效 Banner 列表。
func (s *Service) PortalList(ctx context.Context, q PortalListParam) ([]Banner, error) {
	return s.repo.ListPortal(ctx, q, security.StatusEnabled)
}

func statusOr(st string) string {
	if st == "" {
		return security.StatusEnabled
	}
	return st
}
