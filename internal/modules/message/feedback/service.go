// internal/modules/message/feedback/service.go 业务服务。
//
// Author: Charlie

package feedback

import (
	"context"
	"time"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service åé¦ˆä¸šåŠ¡æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService æž„é€ åé¦ˆæœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New æž„å»º message.feedback æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "message.feedback",
		Order:  41,
		Models: []any{&Feedback{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Submit æäº¤åé¦ˆã€‚
func (s *Service) Submit(ctx context.Context, req CreateParam, meta SubmitMeta) error {
	row := Feedback{
		ID: idgen.Next(), Title: req.Title, Content: req.Content, Category: req.Category,
		Contact: req.Contact, AttachObjectNames: jsonList(req.AttachObjectNames),
		Status: "PENDING", SubmitterAccountType: meta.AccountType, SubmitterAccountID: meta.AccountID,
		CreatedBy: &meta.CreatedBy, UpdatedBy: &meta.CreatedBy,
	}
	return s.repo.Create(ctx, &row)
}

// Update å›žå¤/æ›´æ–°åé¦ˆã€‚
func (s *Service) Update(ctx context.Context, req UpdateParam, meta ReplyMeta) error {
	return s.repo.UpdateReply(ctx, req.ID, req.Status, req.Reply, meta, time.Now().UTC())
}

// Delete æ‰¹é‡åˆ é™¤ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*Feedback, error) {
	return s.repo.GetByID(ctx, id)
}

// PageAdmin ç®¡ç†ç«¯åˆ†é¡µã€‚
func (s *Service) PageAdmin(ctx context.Context, q PageParam) (rows []Feedback, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PageAdmin(ctx, q)
	return rows, total, current, size, err
}

// MyPage æˆ‘çš„åé¦ˆåˆ†é¡µã€‚
func (s *Service) MyPage(ctx context.Context, q schema.PageQuery, accountID, accountType string) (rows []Feedback, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PageBySubmitter(ctx, accountID, accountType, current, size)
	return rows, total, current, size, err
}

// MyDetail æˆ‘çš„åé¦ˆè¯¦æƒ…ã€‚
func (s *Service) MyDetail(ctx context.Context, id, accountID, accountType string) (*Feedback, error) {
	return s.repo.GetBySubmitter(ctx, id, accountID, accountType)
}
