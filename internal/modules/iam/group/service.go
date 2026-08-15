// internal/modules/iam/group/service.go 业务服务。
//
// Author: Charlie

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

// Service 用户组服务（授权经 relation 模块，成员账号视图经 result 模块）。
//
// Author: Charlie
type Service struct {
	db        *gorm.DB
	repo      *Repo
	rel       *relation.Service
	roles     *role.Repo
	resources *resource.Service
	clients   *client.Service
	sessions  *security.SessionStore
}

// NewService 构造用户组服务。
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

// New 构建 iam.group 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	s.sessions = d.Sessions
	return module.Module{
		Name:   "iam.group",
		Models: []any{&Group{}},
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

// membersOf 列出用户组当前成员账号 ID。
func (s *Service) membersOf(ctx context.Context, groupID string) ([]string, error) {
	return s.rel.ListTargetIDs(ctx, relation.SubjectGroup, groupID, relation.RelGroupUser, "")
}

// Create 创建用户组。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Group{
		ID: idgen.Next(), Name: req.Name, OwnerDeptID: req.OwnerDeptID,
		Description: req.Description, Status: orStatus(req.Status), Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新用户组。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{
		"name": req.Name, "owner_dept_id": req.OwnerDeptID,
		"description": req.Description, "status": orStatus(req.Status),
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除（先清组关联，再删组；对齐 hei-boot GroupServiceImpl.delete）。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	_ = s.rel.DeleteBySubjectIDs(ctx, relation.SubjectGroup, ids, "")
	_ = s.rel.DeleteByTargetIDs(ctx, relation.TargetGroup, ids, "")
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 用户组详情。
func (s *Service) Detail(ctx context.Context, id string) (*Group, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
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
