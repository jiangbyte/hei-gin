// internal/modules/sys/audit/service.go 业务服务。
//
// Author: Charlie

package audit

import (
	"context"

	"gorm.io/gorm"

	auditpkg "hei-gin/internal/framework/platform/audit"
	"hei-gin/internal/framework/platform/module"
)

// Service 审计日志业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造审计服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 sys.audit 模块（含 auditAlertJob）。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	m := module.Module{
		Name:   "sys.audit",
		Models: []any{&OperationLog{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
	return s.withJobs(m, d.Notify, d.Cfg)
}

// Page 分页查询。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []OperationLog, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	for i := range rows {
		enrichOperationLog(&rows[i])
	}
	enrichOperatorNames(ctx, s.repo.db, rows)
	return rows, total, current, size, err
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*OperationLog, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	enrichOperationLog(row)
	enrichOperatorNames(ctx, s.repo.db, []OperationLog{*row})
	return row, nil
}

func enrichOperationLog(row *OperationLog) {
	if row == nil {
		return
	}
	rt := ""
	if row.ResourceType != nil {
		rt = *row.ResourceType
	}
	auditpkg.EnrichActivityLabels(row.Module, rt, row.Action, &row.ActionName, &row.ActionType, &row.ModuleLabel)
}

// MyPage 当前用户本人审计日志分页。
func (s *Service) MyPage(ctx context.Context, accountID string, q PageParam) (rows []OperationLog, total int64, current, size int, err error) {
	q.AccountID = accountID
	return s.Page(ctx, q)
}

// MyDetail 当前用户本人审计详情。
func (s *Service) MyDetail(ctx context.Context, accountID, id string) (*OperationLog, error) {
	row, err := s.Detail(ctx, id)
	if err != nil {
		return nil, err
	}
	if row.AccountID == nil || *row.AccountID != accountID {
		return nil, gorm.ErrRecordNotFound
	}
	return row, nil
}
