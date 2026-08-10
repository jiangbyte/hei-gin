package feedback

import (
	"context"
	"time"

	"gorm.io/gorm"

	"hei-gin/framework/core/schema"
	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 反馈业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造反馈服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 message.feedback 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "message.feedback",
		Order:  41,
		Models: []any{&Feedback{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Submit 提交反馈。
func (s *Service) Submit(ctx context.Context, req CreateParam, meta SubmitMeta) error {
	row := Feedback{
		ID: idgen.Next(), Title: req.Title, Content: req.Content, Category: req.Category,
		Contact: req.Contact, AttachObjectNames: jsonList(req.AttachObjectNames),
		Status: "PENDING", SubmitterAccountType: meta.AccountType, SubmitterAccountID: meta.AccountID,
		CreatedBy: &meta.CreatedBy, UpdatedBy: &meta.CreatedBy,
	}
	return s.repo.Create(ctx, &row)
}

// Update 回复/更新反馈。
func (s *Service) Update(ctx context.Context, req UpdateParam, meta ReplyMeta) error {
	return s.repo.UpdateReply(ctx, req.ID, req.Status, req.Reply, meta, time.Now().UTC())
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*Feedback, error) {
	return s.repo.GetByID(ctx, id)
}

// PageAdmin 管理端分页。
func (s *Service) PageAdmin(ctx context.Context, q PageParam) (rows []Feedback, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PageAdmin(ctx, q)
	return rows, total, current, size, err
}

// MyPage 我的反馈分页。
func (s *Service) MyPage(ctx context.Context, q schema.PageQuery, accountID, accountType string) (rows []Feedback, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PageBySubmitter(ctx, accountID, accountType, current, size)
	return rows, total, current, size, err
}

// MyDetail 我的反馈详情。
func (s *Service) MyDetail(ctx context.Context, id, accountID, accountType string) (*Feedback, error) {
	return s.repo.GetBySubmitter(ctx, id, accountID, accountType)
}
