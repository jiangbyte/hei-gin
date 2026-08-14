package group

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
	"hei-gin/internal/modules/iam/role"
	"hei-gin/internal/modules/shared"
)

// Service ç”¨æˆ·ç»„æœåŠ¡ï¼ˆæŽˆæƒç» relation æ¨¡å—ï¼Œæˆå‘˜è´¦å·è§†å›¾ç» result æ¨¡å—ï¼‰ã€‚
//
// Author: Charlie
type Service struct {
	db        *gorm.DB
	repo      *Repo
	rel       *relation.Service
	roles     *role.Repo
	resources *resource.Service
	clients   *client.Service
}

// NewService æž„é€ ç”¨æˆ·ç»„æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:        db,
		repo:      NewRepo(db),
		rel:       relation.NewService(db),
		roles:     role.NewRepo(db),
		resources: resource.NewService(db),
		clients:   client.NewService(db),
	}
}

// New æž„å»º iam.group æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.group",
		Models: []any{&Group{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create åˆ›å»ºç”¨æˆ·ç»„ã€‚
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Group{
		ID: idgen.Next(), Name: req.Name, OwnerDeptID: req.OwnerDeptID,
		Description: req.Description, Status: orStatus(req.Status), Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.Create(ctx, &row)
}

// Update æ›´æ–°ç”¨æˆ·ç»„ã€‚
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{
		"name": req.Name, "owner_dept_id": req.OwnerDeptID,
		"description": req.Description, "status": orStatus(req.Status),
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete æ‰¹é‡åˆ é™¤ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail ç”¨æˆ·ç»„è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*Group, error) {
	return s.repo.GetByID(ctx, id)
}

// Page åˆ†é¡µã€‚
func (s *Service) Page(ctx context.Context, p PageParam) (rows []Group, total int64, current, size int, err error) {
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
