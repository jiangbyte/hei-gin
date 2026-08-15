// internal/modules/iam/role/service.go 业务服务。
//
// Author: Charlie

package role

import (
	"context"
	"fmt"

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

// Service 角色服务（授权经 relation 模块，成员账号视图经 result 模块）。
//
// Author: Charlie
type Service struct {
	db        *gorm.DB
	repo      *Repo
	rel       *relation.Service
	resources *resource.Service
	clients   *client.Service
	sessions  *security.SessionStore
}

// NewService 构造角色服务。
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:        db,
		repo:      NewRepo(db),
		rel:       relation.NewService(db),
		resources: resource.NewService(db),
		clients:   client.NewService(db),
	}
}

// New 构建 iam.role 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	s.sessions = d.Sessions
	return module.Module{
		Name:   "iam.role",
		Models: []any{&Role{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// invalidateAccounts 授权变更后强制受影响账号下线（对齐 hei-boot logoutAccounts）。
func (s *Service) invalidateAccounts(ctx context.Context, accountIDs []string) {
	if s.sessions == nil || len(accountIDs) == 0 {
		return
	}
	for _, id := range accountIDs {
		if id != "" {
			_ = s.sessions.DeleteAllForAccount(ctx, id)
		}
	}
}

// Create 创建角色（code 唯一校验；对齐 hei-boot RoleServiceImpl.create）。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	if _, err := s.repo.FindByCode(ctx, req.Code); err == nil {
		return fmt.Errorf("角色编码已存在")
	}
	ext := req.Extra
	if len(ext) == 0 {
		ext = datatypes.JSON([]byte("{}"))
	}
	row := Role{
		ID: idgen.Next(), Code: req.Code, Name: req.Name,
		Category: orDef(req.Category, "SYS"), ScopeType: orDef(req.ScopeType, "PLATFORM"),
		OwnerDeptID: req.OwnerDeptID, Sort: orSort(req.Sort), Status: orStatus(req.Status),
		IsBuiltin: boolOr(req.IsBuiltin), Description: req.Description, Extra: ext,
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新角色。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{
		"code": req.Code, "name": req.Name, "category": orDef(req.Category, "SYS"),
		"scope_type": orDef(req.ScopeType, "PLATFORM"), "owner_dept_id": req.OwnerDeptID,
		"sort": orSort(req.Sort), "status": orStatus(req.Status), "description": req.Description,
	}
	if req.IsBuiltin != nil {
		updates["is_builtin"] = *req.IsBuiltin
	}
	if len(req.Extra) > 0 {
		updates["extra"] = req.Extra
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除（先清角色关联，再删角色；对齐 hei-boot RoleServiceImpl.delete）。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	_ = s.rel.DeleteBySubjectIDs(ctx, relation.SubjectRole, ids, "")
	_ = s.rel.DeleteByTargetIDs(ctx, relation.TargetRole, ids, "")
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 角色详情。
func (s *Service) Detail(ctx context.Context, id string) (*Role, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, p PageParam) (rows []Role, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p)
	return rows, total, current, size, err
}

func boolOr(p *bool) bool { return p != nil && *p }

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
