// internal/modules/iam/dept/service.go 业务服务。
//
// Author: Charlie

package dept

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/core/security/datascope"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/iam/relation"
)

// Service 部门服务。
//
// Author: Charlie
type Service struct {
	repo *Repo
	rel  *relation.Service
}

// NewService 构造部门服务。
func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewRepo(db), rel: relation.NewService(db)}
}

// New 构建 iam.dept 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.dept",
		Models: []any{&Dept{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建部门。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Dept{
		ID: idgen.Next(), ParentID: req.ParentID, MasterID: req.MasterID, DeputyMasterID: req.DeputyMasterID,
		Name: req.Name, Category: req.Category, Sort: req.Sort, IsVirtual: req.IsVirtual,
		Status: orStatus(req.Status), Extra: datatypes.JSON([]byte("{}")),
	}
	if row.Sort == 0 {
		row.Sort = 99
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新部门（数据范围校验；对齐 hei-boot assertOwnerOrDeptAccessible 以 dept.id 为归属）。
func (s *Service) Update(ctx context.Context, req EditParam, sess *security.SessionPayload) error {
	cur, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := s.assertScope(sess, cur); err != nil {
		return err
	}
	updates := map[string]any{
		"parent_id": req.ParentID, "master_id": req.MasterID, "deputy_master_id": req.DeputyMasterID,
		"name": req.Name, "category": req.Category, "sort": req.Sort, "is_virtual": req.IsVirtual,
		"status": orStatus(req.Status),
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除（先校验数据范围、清引用关系，再删部门；对齐 hei-boot DeptServiceImpl.delete）。
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
	_ = s.rel.DeleteByTargetIDs(ctx, relation.TargetDept, ids, "")
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 部门详情（数据范围校验）。
func (s *Service) Detail(ctx context.Context, id string, sess *security.SessionPayload) (*Dept, error) {
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
func (s *Service) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []Dept, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p, sess)
	return rows, total, current, size, err
}

// Tree 部门树（数据范围过滤）。
func (s *Service) Tree(ctx context.Context, sess *security.SessionPayload) ([]TreeNode, error) {
	rows, err := s.repo.ListAll(ctx, sess)
	if err != nil {
		return nil, err
	}
	return buildDeptTree(rows, nil), nil
}

// assertScope 数据范围断言：ALL 放行；SELF 比创建人；部门类要求当前部门 id 落在可见部门内。
func (s *Service) assertScope(sess *security.SessionPayload, row *Dept) error {
	if sess == nil {
		return datascope.ErrDenied
	}
	var ownerAccount string
	if row.CreatedBy != nil {
		ownerAccount = *row.CreatedBy
	}
	return datascope.AssertKey(sess, "iam:dept:page", row.ID, ownerAccount)
}

func buildDeptTree(rows []Dept, parent *string) []TreeNode {
	ids := make(map[string]struct{}, len(rows))
	for _, r := range rows {
		ids[r.ID] = struct{}{}
	}
	out := make([]TreeNode, 0)
	for _, r := range rows {
		p := r.ParentID
		if p != nil && *p != "" {
			if _, ok := ids[*p]; !ok {
				p = nil
			}
		}
		if eqPtr(p, parent) {
			n := TreeNode{Dept: r, Children: buildDeptTree(rows, &r.ID)}
			out = append(out, n)
		}
	}
	return out
}

func eqPtr(a *string, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func orStatus(st string) string {
	if st == "" {
		return security.StatusEnabled
	}
	return st
}
