// internal/modules/sys/audit/service.go 业务服务。
//
// Author: Charlie

package audit

import (
	"context"

	"gorm.io/gorm"
)

// Service 审计日志业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造审计服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// Page 分页查询。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []OperationLog, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*OperationLog, error) {
	return s.repo.GetByID(ctx, id)
}
