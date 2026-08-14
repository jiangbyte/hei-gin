package role

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/iam/client"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/resource"
	"hei-gin/internal/modules/shared"
)

// Service è§’è‰²æœåŠ¡ï¼ˆæŽˆæƒç» relation æ¨¡å—ï¼Œæˆå‘˜è´¦å·è§†å›¾ç» result æ¨¡å—ï¼‰ã€‚
//
// Author: Charlie
type Service struct {
	db        *gorm.DB
	repo      *Repo
	rel       *relation.Service
	resources *resource.Service
	clients   *client.Service
}

// NewService æž„é€ è§’è‰²æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:        db,
		repo:      NewRepo(db),
		rel:       relation.NewService(db),
		resources: resource.NewService(db),
		clients:   client.NewService(db),
	}
}

// New æž„å»º iam.role æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.role",
		Models: []any{&Role{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create åˆ›å»ºè§’è‰²ã€‚
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Role{
		ID: idgen.Next(), Code: req.Code, Name: req.Name,
		Category: orDef(req.Category, "SYS"), ScopeType: orDef(req.ScopeType, "PLATFORM"),
		OwnerDeptID: req.OwnerDeptID, Sort: orSort(req.Sort), Status: orStatus(req.Status),
		Description: req.Description, Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.Create(ctx, &row)
}

// Update æ›´æ–°è§’è‰²ã€‚
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{
		"code": req.Code, "name": req.Name, "category": orDef(req.Category, "SYS"),
		"scope_type": orDef(req.ScopeType, "PLATFORM"), "owner_dept_id": req.OwnerDeptID,
		"sort": orSort(req.Sort), "status": orStatus(req.Status), "description": req.Description,
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete æ‰¹é‡åˆ é™¤ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail è§’è‰²è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*Role, error) {
	return s.repo.GetByID(ctx, id)
}

// Page åˆ†é¡µã€‚
func (s *Service) Page(ctx context.Context, p PageParam) (rows []Role, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p)
	return rows, total, current, size, err
}

func orStatus(st string) string {
	if st == "" {
		return security.StatusEnabled
	}
	return st
}

func orDef(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

func orSort(n int) int {
	if n == 0 {
		return 99
	}
	return n
}
