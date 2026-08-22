// internal/modules/sys/notice/service.go 业务服务。
//
// Author: Charlie

package notice

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
)

// Service 通知业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造通知服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 sys.notice 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.notice",
		Order:  40,
		Models: []any{&Notice{}, &NoticeRead{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建通知（对齐 hei-boot：规范化 kind/目标范围并校验）。
func (s *Service) Create(ctx context.Context, req CreateParam, createdBy, updatedBy *string) error {
	if err := s.normalizeAndValidate(&req); err != nil {
		return err
	}
	row := fromCreate(req)
	row.ID = idgen.Next()
	row.CreatedBy = createdBy
	row.UpdatedBy = updatedBy
	return s.repo.Create(ctx, &row)
}

// Update 更新通知。
func (s *Service) Update(ctx context.Context, req UpdateParam, updatedBy *string) error {
	if err := s.normalizeAndValidate(&req.CreateParam); err != nil {
		return err
	}
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
	for i := range rows {
		rows[i].IsRead = nil
	}
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

// Pin 置顶/取消置顶（仅公告允许置顶）。
func (s *Service) Pin(ctx context.Context, req PinParam) error {
	row, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if row.Kind != "ANNOUNCEMENT" {
		return errPinOnlyAnnouncement
	}
	return s.repo.Update(ctx, req.ID, map[string]any{
		"is_pinned": req.IsPinned, "pinned_until": req.PinnedUntil,
	})
}

// PagePublished 已发布通知分页（按当前用户可见性过滤）。
func (s *Service) PagePublished(ctx context.Context, q PageParam, accountType, accountID string) (rows []Notice, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PagePublished(ctx, q, accountType, accountID)
	return rows, total, current, size, err
}

// MyDetail 用户端详情：校验可见性、公告累加浏览并标记已读。
func (s *Service) MyDetail(ctx context.Context, id, accountType, accountID string) (*Notice, error) {
	row, err := s.repo.GetPublishedByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !visibleTo(row, accountType, accountID) {
		return nil, gorm.ErrRecordNotFound
	}
	if row.Kind == "ANNOUNCEMENT" {
		_ = s.repo.IncrViewCount(ctx, row.ID, row.ViewCount+1)
	}
	_ = s.MarkRead(ctx, ReadRecord{
		NoticeID: id, AccountType: accountType, AccountID: accountID, ReadAt: time.Now().UTC(),
	})
	row.IsRead = strPtr("true")
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

// MarkReads 批量标记已读（先查已存在再批次插入，避免 N 次 FirstOrCreate）。
func (s *Service) MarkReads(ctx context.Context, accountType, accountID string, noticeIDs []string, readAt time.Time) error {
	if len(noticeIDs) == 0 {
		return nil
	}
	existing, err := s.repo.ListExistingReadNoticeIDs(ctx, accountType, accountID, noticeIDs)
	if err != nil {
		return err
	}
	rows := make([]NoticeRead, 0, len(noticeIDs))
	for _, id := range noticeIDs {
		if id == "" {
			continue
		}
		if _, ok := existing[id]; ok {
			continue
		}
		rows = append(rows, NoticeRead{
			ID: idgen.Next(), NoticeID: id, AccountType: accountType,
			AccountID: accountID, ReadAt: readAt,
		})
	}
	return s.repo.CreateReadsInBatches(ctx, rows)
}

// MarkAllRead 全部标记已读。
func (s *Service) MarkAllRead(ctx context.Context, accountType, accountID string, readAt time.Time) error {
	ids, err := s.repo.ListUnreadIDs(ctx, accountType, accountID)
	if err != nil {
		return err
	}
	rows := make([]NoticeRead, 0, len(ids))
	for _, id := range ids {
		rows = append(rows, NoticeRead{
			ID: idgen.Next(), NoticeID: id, AccountType: accountType,
			AccountID: accountID, ReadAt: readAt,
		})
	}
	return s.repo.CreateReadsInBatches(ctx, rows)
}

// visibleTo 校验通知对当前用户可见（与 applyVisibility 逻辑一致）。
func visibleTo(n *Notice, accountType, accountID string) bool {
	scope := n.TargetScope
	if scope == "" {
		scope = "ALL"
	}
	var types []string
	_ = json.Unmarshal(n.TargetAccountTypes, &types)
	typeMatch := len(types) == 0
	for _, t := range types {
		if t == accountType {
			typeMatch = true
			break
		}
	}
	if scope == "ALL" || scope == "ACCOUNT_TYPE" {
		return typeMatch
	}
	if scope == "SPECIFIC" {
		var ids []string
		_ = json.Unmarshal(n.TargetAccountIDs, &ids)
		for _, id := range ids {
			if id == accountID {
				return true
			}
		}
		return false
	}
	return false
}

// normalizeAndValidate 规范化 kind/目标范围并按类型校验（对齐 hei-boot）。
func (s *Service) normalizeAndValidate(req *CreateParam) error {
	kind := strings.ToUpper(strings.TrimSpace(req.Kind))
	if kind != "NOTIFICATION" && kind != "ANNOUNCEMENT" {
		return errInvalidKind
	}
	req.Kind = kind

	scope := strings.ToUpper(strings.TrimSpace(req.TargetScope))
	if scope == "" {
		scope = "ALL"
	}
	if scope != "ALL" && scope != "ACCOUNT_TYPE" && scope != "SPECIFIC" {
		return errInvalidScope
	}
	req.TargetScope = scope
	if len(req.TargetAccountTypes) == 0 {
		return errTargetTypesRequired
	}
	if scope == "SPECIFIC" && len(req.TargetAccountIDs) == 0 {
		return errTargetAccountsRequired
	}

	if kind == "ANNOUNCEMENT" {
		if !hasEnabledLocation(req.PublishLocations) {
			return errAnnouncementLocationRequired
		}
	} else {
		if req.Category == nil || strings.TrimSpace(*req.Category) == "" {
			return errNotificationCategoryRequired
		}
		req.IsPinned = false
		req.PinnedUntil = nil
		req.ExpireAt = nil
	}

	// 状态规范化：ENABLED/ENABLE 视为草稿。
	req.Status = resolveStatus(req.Status)
	if req.Status == "PUBLISHED" && req.PublishAt == nil {
		now := time.Now().UTC()
		req.PublishAt = &now
	}
	return nil
}

func resolveStatus(status string) string {
	normalized := strings.ToUpper(strings.TrimSpace(status))
	if normalized == "" || normalized == "ENABLED" || normalized == "ENABLE" {
		return "DRAFT"
	}
	if normalized == "DRAFT" || normalized == "PUBLISHED" || normalized == "REVOKED" {
		return normalized
	}
	return "DRAFT"
}

func hasEnabledLocation(locations map[string]any) bool {
	if len(locations) == 0 {
		return false
	}
	for _, v := range locations {
		if b, ok := v.(bool); ok && b {
			return true
		}
		if s, ok := v.(string); ok && strings.EqualFold(s, "true") {
			return true
		}
	}
	return false
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

var (
	errPinOnlyAnnouncement          = &noticeErr{msg: "仅公告支持置顶"}
	errInvalidKind                  = &noticeErr{msg: "kind 必须是 NOTIFICATION 或 ANNOUNCEMENT"}
	errInvalidScope                 = &noticeErr{msg: "目标范围仅支持全部 / 按账户类型 / 指定用户"}
	errTargetTypesRequired          = &noticeErr{msg: "必须选择目标账户类型"}
	errTargetAccountsRequired       = &noticeErr{msg: "指定用户时必须选择目标用户"}
	errAnnouncementLocationRequired = &noticeErr{msg: "公告必须选择至少一个发布位置"}
	errNotificationCategoryRequired = &noticeErr{msg: "通知必须选择分类"}
)

type noticeErr struct{ msg string }

func (e *noticeErr) Error() string { return e.msg }
