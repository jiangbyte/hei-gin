// internal/modules/iam/group/service.go 业务服务。
//
// Author: Charlie

package group

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/core/security/datascope"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/iam/client"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/resource"
	"hei-gin/internal/modules/iam/role"
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
func New(d *module.Deps) module.Module {
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
			_ = s.sessions.DeleteAllForAccountAnyType(ctx, id)
		}
	}
}

// membersOf 列出用户组当前成员账号 ID。
func (s *Service) membersOf(ctx context.Context, groupID string) ([]string, error) {
	return s.rel.ListSubjectIDsByTarget(ctx, relation.RelAccountGroup, relation.TargetGroup, groupID)
}

// Create 创建用户组。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Group{
		ID: idgen.Next(), Name: req.Name, OwnerDeptID: req.OwnerDeptID,
		Description: req.Description, Status: orStatus(req.Status), Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新用户组（数据范围校验；对齐 hei-boot assertOwnerOrDeptAccessible）。
func (s *Service) Update(ctx context.Context, req EditParam, sess *security.SessionPayload) error {
	cur, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := s.assertScope(sess, cur); err != nil {
		return err
	}
	updates := map[string]any{
		"name": req.Name, "owner_dept_id": req.OwnerDeptID,
		"description": req.Description, "status": orStatus(req.Status),
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除（先校验数据范围、清组关联，再删组；对齐 hei-boot GroupServiceImpl.delete）。
func (s *Service) Delete(ctx context.Context, ids []string, sess *security.SessionPayload) error {
	rows, err := s.repo.GetByIDs(ctx, ids)
	if err != nil {
		return err
	}
	for i := range rows {
		if err := s.assertScope(sess, &rows[i]); err != nil {
			return err
		}
	}
	_ = s.rel.DeleteBySubjectIDs(ctx, relation.SubjectGroup, ids, "")
	_ = s.rel.DeleteByTargetIDs(ctx, relation.TargetGroup, ids, "")
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 用户组详情（数据范围校验）。
func (s *Service) Detail(ctx context.Context, id string, sess *security.SessionPayload) (*Group, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.assertScope(sess, row); err != nil {
		return nil, err
	}
	return row, nil
}

// Page 分页（数据范围过滤）。
func (s *Service) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []Group, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p, sess)
	return rows, total, current, size, err
}

// assertScope 数据范围断言：ALL 放行；SELF 比创建人；部门类要求 owner_dept_id 落在可见部门内。
func (s *Service) assertScope(sess *security.SessionPayload, row *Group) error {
	if sess == nil {
		return datascope.ErrDenied
	}
	var ownerDept string
	if row.OwnerDeptID != nil {
		ownerDept = *row.OwnerDeptID
	}
	var ownerAccount string
	if row.CreatedBy != nil {
		ownerAccount = *row.CreatedBy
	}
	return datascope.AssertKey(sess, "iam:group:page", ownerDept, ownerAccount)
}

func orStatus(st string) string {
	if st == "" {
		return security.StatusEnabled
	}
	return st
}
