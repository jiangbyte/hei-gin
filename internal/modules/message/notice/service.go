// internal/modules/message/notice/service.go 业务服务。
//
// Author: Charlie

package notice

import (
	"context"
	"time"

	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service é€šçŸ¥ä¸šåŠ¡æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService æž„é€ é€šçŸ¥æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New æž„å»º message.notice æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "message.notice",
		Order:  40,
		Models: []any{&Notice{}, &NoticeRead{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create åˆ›å»ºé€šçŸ¥ã€‚
func (s *Service) Create(ctx context.Context, req CreateParam, createdBy, updatedBy *string) error {
	row := fromCreate(req)
	row.ID = idgen.Next()
	row.CreatedBy = createdBy
	row.UpdatedBy = updatedBy
	return s.repo.Create(ctx, &row)
}

// Update æ›´æ–°é€šçŸ¥ã€‚
func (s *Service) Update(ctx context.Context, req UpdateParam, updatedBy *string) error {
	row := fromCreate(req.CreateParam)
	updates := map[string]any{
		"kind": row.Kind, "title": row.Title, "content": row.Content, "content_type": row.ContentType,
		"category": row.Category, "severity": row.Severity, "target_scope": row.TargetScope,
		"target_account_types": row.TargetAccountTypes, "target_account_ids": row.TargetAccountIDs,
		"target_dept_ids": row.TargetDeptIDs, "target_role_ids": row.TargetRoleIDs,
		"publish_locations": row.PublishLocations, "is_pinned": row.IsPinned, "pinned_until": row.PinnedUntil,
		"status": row.Status, "publish_at": row.PublishAt, "expire_at": row.ExpireAt, "extra": row.Extra,
	}
	if updatedBy != nil {
		updates["updated_by"] = *updatedBy
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete æ‰¹é‡åˆ é™¤ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	if err := s.repo.DeleteByIDs(ctx, ids); err != nil {
		return err
	}
	_ = s.repo.DeleteReadsByNoticeIDs(ctx, ids)
	return nil
}

// Detail è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*Notice, error) {
	return s.repo.GetByID(ctx, id)
}

// PageAdmin ç®¡ç†ç«¯åˆ†é¡µã€‚
func (s *Service) PageAdmin(ctx context.Context, q PageParam) (rows []Notice, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PageAdmin(ctx, q)
	return rows, total, current, size, err
}

// Publish å‘å¸ƒé€šçŸ¥ã€‚
func (s *Service) Publish(ctx context.Context, ids []string, p PublishParam) error {
	return s.repo.UpdateByIDs(ctx, ids, map[string]any{
		"status": p.Status, "publish_at": p.PublishAt,
		"sender_account_id": p.SenderAccountID, "sender_account_type": p.SenderAccountType,
		"updated_by": p.UpdatedBy,
	})
}

// Revoke æ’¤å›žé€šçŸ¥ã€‚
func (s *Service) Revoke(ctx context.Context, ids []string, p RevokeParam) error {
	return s.repo.UpdateByIDs(ctx, ids, map[string]any{
		"status": p.Status, "revoked_at": p.RevokedAt,
	})
}

// Pin ç½®é¡¶/å–æ¶ˆç½®é¡¶ã€‚
func (s *Service) Pin(ctx context.Context, req PinParam) error {
	return s.repo.Update(ctx, req.ID, map[string]any{
		"is_pinned": req.IsPinned, "pinned_until": req.PinnedUntil,
	})
}

// PagePublished å·²å‘å¸ƒé€šçŸ¥åˆ†é¡µã€‚
func (s *Service) PagePublished(ctx context.Context, q PageParam) (rows []Notice, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PagePublished(ctx, q)
	return rows, total, current, size, err
}

// MyDetail ç”¨æˆ·ç«¯è¯¦æƒ…ï¼ˆå¢žåŠ æµè§ˆé‡ï¼‰ã€‚
func (s *Service) MyDetail(ctx context.Context, id string) (*Notice, error) {
	row, err := s.repo.GetPublishedByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.repo.IncrViewCount(ctx, row.ID, row.ViewCount+1)
	return row, nil
}

// UnreadCount æœªè¯»æ•°ã€‚
func (s *Service) UnreadCount(ctx context.Context, accountType, accountID string) (int64, error) {
	return s.repo.CountUnread(ctx, accountType, accountID)
}

// MarkRead æ ‡è®°å·²è¯»ã€‚
func (s *Service) MarkRead(ctx context.Context, rec ReadRecord) error {
	row := NoticeRead{
		ID: idgen.Next(), NoticeID: rec.NoticeID, AccountType: rec.AccountType,
		AccountID: rec.AccountID, ReadAt: rec.ReadAt,
	}
	return s.repo.FirstOrCreateRead(ctx, row)
}

// MarkAllRead å…¨éƒ¨æ ‡è®°å·²è¯»ã€‚
func (s *Service) MarkAllRead(ctx context.Context, accountType, accountID string, readAt time.Time) error {
	ids, err := s.repo.ListUnreadIDs(ctx, accountType, accountID)
	if err != nil {
		return err
	}
	for _, id := range ids {
		row := NoticeRead{
			ID: idgen.Next(), NoticeID: id, AccountType: accountType,
			AccountID: accountID, ReadAt: readAt,
		}
		_ = s.repo.CreateRead(ctx, &row)
	}
	return nil
}

func fromCreate(req CreateParam) Notice {
	return Notice{
		Kind: req.Kind, Title: req.Title, Content: req.Content, ContentType: req.ContentType,
		Category: req.Category, Severity: req.Severity, TargetScope: req.TargetScope,
		TargetAccountTypes: jsonList(req.TargetAccountTypes), TargetAccountIDs: jsonList(req.TargetAccountIDs),
		TargetDeptIDs: jsonList(req.TargetDeptIDs), TargetRoleIDs: jsonList(req.TargetRoleIDs),
		PublishLocations: jsonObj(req.PublishLocations), IsPinned: req.IsPinned, PinnedUntil: req.PinnedUntil,
		Status: req.Status, PublishAt: req.PublishAt, ExpireAt: req.ExpireAt, Extra: jsonObj(req.Extra),
	}
}
