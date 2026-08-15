// internal/modules/sys/notice/repo.go 持久化仓储。
//
// Author: Charlie

package notice

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Repo 通知持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建通知。
func (r *Repo) Create(ctx context.Context, row *Notice) error {
	return r.with(ctx).Create(row).Error
}

// Update 按 ID 更新字段。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Notice{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateByIDs 批量更新。
func (r *Repo) UpdateByIDs(ctx context.Context, ids []string, updates map[string]any) error {
	return r.with(ctx).Model(&Notice{}).Where("id IN ?", ids).Updates(updates).Error
}

// DeleteByIDs 批量删除通知。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Notice{}).Error
}

// DeleteReadsByNoticeIDs 删除通知已读记录。
func (r *Repo) DeleteReadsByNoticeIDs(ctx context.Context, noticeIDs []string) error {
	return r.with(ctx).Where("notice_id IN ?", noticeIDs).Delete(&NoticeRead{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*Notice, error) {
	var row Notice
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// GetPublishedByID 查询已发布通知。
func (r *Repo) GetPublishedByID(ctx context.Context, id string) (*Notice, error) {
	var row Notice
	if err := r.with(ctx).First(&row, "id = ? AND status = ?", id, "PUBLISHED").Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// IncrViewCount 增加浏览次数。
func (r *Repo) IncrViewCount(ctx context.Context, id string, count int) error {
	return r.with(ctx).Model(&Notice{}).Where("id = ?", id).UpdateColumn("view_count", count).Error
}

// PageAdmin 管理端分页。
func (r *Repo) PageAdmin(ctx context.Context, q PageParam) (rows []Notice, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&Notice{})
	if q.Title != "" {
		db = db.Where("title ILIKE ?", "%"+q.Title+"%")
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if q.Kind != "" {
		db = db.Where("kind = ?", q.Kind)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("created_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// PagePublished 已发布通知分页（含可见性过滤、过期过滤与已读回填）。
// accountID 为空时跳过已读回填（如匿名门户公告列表）。
func (r *Repo) PagePublished(ctx context.Context, q PageParam, accountType, accountID string) (rows []Notice, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&Notice{}).Where("status = ?", "PUBLISHED")
	db = applyVisibility(db, accountType, accountID)
	db = db.Where("kind <> ? OR expire_at IS NULL OR expire_at > ?", "ANNOUNCEMENT", time.Now())
	if q.Kind != "" {
		db = db.Where("kind = ?", q.Kind)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err = db.Order("is_pinned DESC, publish_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	if accountID != "" && len(rows) > 0 {
		r.backfillRead(ctx, accountType, accountID, rows)
	}
	return rows, total, nil
}

// CountUnread 统计当前用户可见的未读已发布通知数。
func (r *Repo) CountUnread(ctx context.Context, accountType, accountID string) (int64, error) {
	sub := r.db.Model(&NoticeRead{}).Select("notice_id").
		Where("account_type = ? AND account_id = ?", accountType, accountID)
	db := r.with(ctx).Model(&Notice{}).
		Where("status = ? AND id NOT IN (?)", "PUBLISHED", sub)
	db = applyVisibility(db, accountType, accountID)
	db = db.Where("kind <> ? OR expire_at IS NULL OR expire_at > ?", "ANNOUNCEMENT", time.Now())
	var total int64
	err := db.Count(&total).Error
	return total, err
}

// ListUnreadIDs 列出当前用户可见的未读已发布通知 ID。
func (r *Repo) ListUnreadIDs(ctx context.Context, accountType, accountID string) ([]string, error) {
	sub := r.db.Model(&NoticeRead{}).Select("notice_id").
		Where("account_type = ? AND account_id = ?", accountType, accountID)
	db := r.with(ctx).Model(&Notice{}).
		Where("status = ? AND id NOT IN (?)", "PUBLISHED", sub)
	db = applyVisibility(db, accountType, accountID)
	db = db.Where("kind <> ? OR expire_at IS NULL OR expire_at > ?", "ANNOUNCEMENT", time.Now())
	var ids []string
	err := db.Pluck("id", &ids).Error
	return ids, err
}

// backfillRead 批量回填当前页已读标记（避免逐行查已读的 N+1）。
func (r *Repo) backfillRead(ctx context.Context, accountType, accountID string, rows []Notice) {
	ids := make([]string, 0, len(rows))
	for i := range rows {
		ids = append(ids, rows[i].ID)
	}
	var reads []NoticeRead
	if err := r.db.WithContext(ctx).Where("account_type = ? AND account_id = ? AND notice_id IN ?",
		accountType, accountID, ids).Find(&reads).Error; err != nil {
		return
	}
	readSet := make(map[string]struct{}, len(reads))
	for _, rd := range reads {
		readSet[rd.NoticeID] = struct{}{}
	}
	for i := range rows {
		_, ok := readSet[rows[i].ID]
		rows[i].IsRead = ok
	}
}

// applyVisibility 按目标范围过滤已发布通知（ALL/ACCOUNT_TYPE 匹配类型 JSON，SPECIFIC 匹配账号 JSON）。
func applyVisibility(db *gorm.DB, accountType, accountID string) *gorm.DB {
	cond := "(target_scope IN (''ALL'',''ACCOUNT_TYPE'') AND (target_account_types IS NULL OR jsonb_array_length(COALESCE(target_account_types::jsonb,''[]''::jsonb)) = 0 OR jsonb_exists((target_account_types)::jsonb, ?)))"
	args := []any{accountType}
	if accountID != "" {
		cond += " OR (target_scope = ''SPECIFIC'' AND jsonb_exists((target_account_ids)::jsonb, ?))"
		args = append(args, accountID)
	}
	return db.Where(cond, args...)
}

// FirstOrCreateRead 标记已读（幂等）。
func (r *Repo) FirstOrCreateRead(ctx context.Context, row NoticeRead) error {
	return r.with(ctx).
		Where("notice_id = ? AND account_type = ? AND account_id = ?", row.NoticeID, row.AccountType, row.AccountID).
		FirstOrCreate(&row).Error
}

// CreateRead 创建已读记录。
func (r *Repo) CreateRead(ctx context.Context, row *NoticeRead) error {
	return r.with(ctx).Create(row).Error
}
