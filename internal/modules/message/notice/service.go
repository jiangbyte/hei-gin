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

// Service 通知业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造通知服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 message.notice 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "message.notice",
		Order:  40,
		Models: []any{&Notice{}, &NoticeRead{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建通知。
func (s *Service) Create(ctx context.Context, req CreateParam, createdBy, updatedBy *string) error {
	row := fromCreate(req)
	row.ID = idgen.Next()
	row.CreatedBy = createdBy
	row.UpdatedBy = updatedBy
	return s.repo.Create(ctx, &row)
}

// Update 更新通知。
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

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	if err := s.repo.DeleteByIDs(ctx, ids); err != nil {
		return err
	}
	_ = s.repo.DeleteReadsByNoticeIDs(ctx, ids)
	return nil
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*Notice, error) {
	return s.repo.GetByID(ctx, id)
}

// PageAdmin 管理端分页。
func (s *Service) PageAdmin(ctx context.Context, q PageParam) (rows []Notice, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PageAdmin(ctx, q)
	return rows, total, current, size, err
}

// Publish 发布通知。
func (s *Service) Publish(ctx context.Context, ids []string, p PublishParam) error {
	return s.repo.UpdateByIDs(ctx, ids, map[string]any{
		"status": p.Status, "publish_at": p.PublishAt,
		"sender_account_id": p.SenderAccountID, "sender_account_type": p.SenderAccountType,
		"updated_by": p.UpdatedBy,
	})
}

// Revoke 撤回通知。
func (s *Service) Revoke(ctx context.Context, ids []string, p RevokeParam) error {
	return s.repo.UpdateByIDs(ctx, ids, map[string]any{
		"status": p.Status, "revoked_at": p.RevokedAt,
	})
}

// Pin 置顶/取消置顶。
func (s *Service) Pin(ctx context.Context, req PinParam) error {
	return s.repo.Update(ctx, req.ID, map[string]any{
		"is_pinned": req.IsPinned, "pinned_until": req.PinnedUntil,
	})
}

// PagePublished 已发布通知分页。
func (s *Service) PagePublished(ctx context.Context, q PageParam) (rows []Notice, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PagePublished(ctx, q)
	return rows, total, current, size, err
}

// MyDetail 用户端详情（增加浏览量）。
func (s *Service) MyDetail(ctx context.Context, id string) (*Notice, error) {
	row, err := s.repo.GetPublishedByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.repo.IncrViewCount(ctx, row.ID, row.ViewCount+1)
	return row, nil
}

// UnreadCount 未读数。
func (s *Service) UnreadCount(ctx context.Context, accountType, accountID string) (int64, error) {
	return s.repo.CountUnread(ctx, accountType, accountID)
}

// MarkRead 标记已读。
func (s *Service) MarkRead(ctx context.Context, rec ReadRecord) error {
	row := NoticeRead{
		ID: idgen.Next(), NoticeID: rec.NoticeID, AccountType: rec.AccountType,
		AccountID: rec.AccountID, ReadAt: rec.ReadAt,
	}
	return s.repo.FirstOrCreateRead(ctx, row)
}

// MarkAllRead 全部标记已读。
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
