// internal/modules/iam/position/service.go 业务服务。
//
// Author: Charlie

package position

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/core/security/datascope"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
)

// Service 职位服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造职位服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 iam.position 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.position",
		Models: []any{&Position{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建职位。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Position{
		ID: idgen.Next(), Name: req.Name, Category: req.Category, OwnerDeptID: req.OwnerDeptID,
		Sort: orSort(req.Sort), IsVirtual: req.IsVirtual, Status: orStatus(req.Status),
		Description: req.Description, Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新职位（数据范围校验；对齐 hei-boot assertOwnerOrDeptAccessible）。
func (s *Service) Update(ctx context.Context, req EditParam, sess *security.SessionPayload) error {
	cur, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := s.assertScope(sess, cur); err != nil {
		return err
	}
	updates := map[string]any{
		"name": req.Name, "category": req.Category, "owner_dept_id": req.OwnerDeptID,
		"sort": orSort(req.Sort), "is_virtual": req.IsVirtual, "status": orStatus(req.Status),
		"description": req.Description,
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除（数据范围校验）。
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
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 职位详情（数据范围校验）。
func (s *Service) Detail(ctx context.Context, id string, sess *security.SessionPayload) (*Position, error) {
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
func (s *Service) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []Position, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p, sess)
	return rows, total, current, size, err
}

// assertScope 数据范围断言：ALL 放行；SELF 比创建人；部门类要求 owner_dept_id 落在可见部门内。
func (s *Service) assertScope(sess *security.SessionPayload, row *Position) error {
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
	return datascope.AssertKey(sess, "iam:position:page", ownerDept, ownerAccount)
}

func orStatus(st string) string {
	if st == "" {
		return security.StatusEnabled
	}
	return st
}

func orSort(n int) int {
	if n == 0 {
		return 99
	}
	return n
}
