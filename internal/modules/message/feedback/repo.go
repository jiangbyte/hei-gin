// internal/modules/message/feedback/repo.go 持久化仓储。
//
// Author: Charlie

package feedback

import (
	"context"
	"time"

	"gorm.io/gorm"
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
		db = db.Where("title ILIKE ?", "%"+q.Title+"%")
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
