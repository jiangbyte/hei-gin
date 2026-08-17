// internal/modules/sys/feedback/repo.go 持久化仓储。
//
// Author: Charlie

package feedback

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/db/dialect"
)

// Repo 反馈持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建反馈。
func (r *Repo) Create(ctx context.Context, row *Feedback) error {
	return r.with(ctx).Create(row).Error
}

// Update 按 ID 更新。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Feedback{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByIDs 批量删除。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Feedback{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*Feedback, error) {
	var row Feedback
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// GetBySubmitter 按提交者查询。
func (r *Repo) GetBySubmitter(ctx context.Context, id, accountID, accountType string) (*Feedback, error) {
	var row Feedback
	if err := r.with(ctx).First(&row,
		"id = ? AND submitter_account_id = ? AND submitter_account_type = ?",
		id, accountID, accountType,
	).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageAdmin 管理端分页。
func (r *Repo) PageAdmin(ctx context.Context, q PageParam) (rows []Feedback, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&Feedback{})
	if q.Title != "" {
		db = db.Where(dialect.ILike(db, "title"), "%"+q.Title+"%")
	}
	if q.Category != "" {
		db = db.Where("category = ?", q.Category)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if q.SubmitterAccountType != "" {
		db = db.Where("submitter_account_type = ?", q.SubmitterAccountType)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("created_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// PageBySubmitter 提交者分页。
func (r *Repo) PageBySubmitter(ctx context.Context, accountID, accountType string, cur, size int) (rows []Feedback, total int64, err error) {
	db := r.with(ctx).Model(&Feedback{}).
		Where("submitter_account_id = ? AND submitter_account_type = ?", accountID, accountType)
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("created_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// UpdateReply 更新回复信息。
func (r *Repo) UpdateReply(ctx context.Context, id string, status string, reply *string, meta ReplyMeta, repliedAt time.Time) error {
	return r.Update(ctx, id, map[string]any{
		"status": status, "reply": reply, "replied_by": meta.RepliedBy,
		"replied_at": repliedAt, "updated_by": meta.UpdatedBy,
	})
}

// ListFilesByObjectNames 按对象名集合查询 sys_file 元数据。
func (r *Repo) ListFilesByObjectNames(ctx context.Context, objectNames []string) ([]FileRow, error) {
	if len(objectNames) == 0 {
		return nil, nil
	}
	var rows []FileRow
	err := r.with(ctx).Table("sys_file").
		Select("id", "object_name", "original_name", "content_type", "size").
		Where("object_name IN ?", objectNames).
		Find(&rows).Error
	return rows, err
}

// FileMapByObjectNames 按对象名批量加载文件元数据（供附件回填）。
func (r *Repo) FileMapByObjectNames(ctx context.Context, objectNames []string) map[string]FileRow {
	out := make(map[string]FileRow)
	if len(objectNames) == 0 {
		return out
	}
	rows, err := r.ListFilesByObjectNames(ctx, objectNames)
	if err != nil {
		return out
	}
	for _, row := range rows {
		out[row.ObjectName] = row
	}
	return out
}

// ProfileBrief 账号资料摘要（昵称/头像）。
type ProfileBrief struct {
	Nickname string `gorm:"column:nickname"`
	Avatar   string `gorm:"column:avatar"`
}

// ProfileNames 按账号类型查询资料表昵称/头像（admin/portal 表结构一致）。
func (r *Repo) ProfileNames(ctx context.Context, accountType string, accountIDs []string) map[string]ProfileBrief {
	out := make(map[string]ProfileBrief)
	if len(accountIDs) == 0 {
		return out
	}
	table := "profile_user_admin"
	if strings.EqualFold(accountType, "PORTAL") {
		table = "profile_user_portal"
	}
	var rows []struct {
		AccountID string `gorm:"column:account_id"`
		ProfileBrief
	}
	if err := r.with(ctx).Table(table).
		Select("account_id", "nickname", "avatar").
		Where("account_id IN ?", accountIDs).
		Find(&rows).Error; err != nil {
		return out
	}
	for _, row := range rows {
		out[row.AccountID] = row.ProfileBrief
	}
	return out
}
