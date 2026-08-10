package admin

import (
	"context"

	"gorm.io/gorm"
)

// Repo 管理端用户资料持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// GetProfile 按账号 ID 查资料。
func (r *Repo) GetProfile(ctx context.Context, accountID string) (*Profile, error) {
	var p Profile
	if err := r.with(ctx).First(&p, "account_id = ?", accountID).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// CreateProfile 创建资料。
func (r *Repo) CreateProfile(ctx context.Context, p *Profile) error {
	return r.with(ctx).Create(p).Error
}

// UpdateProfile 更新资料字段。
func (r *Repo) UpdateProfile(ctx context.Context, accountID string, updates map[string]any) error {
	return r.with(ctx).Model(&Profile{}).Where("account_id = ?", accountID).Updates(updates).Error
}

// UpsertProfile 按 account_id 插入或更新资料（供 iam/account 跨模块调用，对齐 boot ProfileApi）。
func (r *Repo) UpsertProfile(ctx context.Context, p *Profile) error {
	var n int64
	if err := r.with(ctx).Model(&Profile{}).Where("account_id = ?", p.AccountID).Count(&n).Error; err != nil {
		return err
	}
	if n == 0 {
		return r.CreateProfile(ctx, p)
	}
	return r.UpdateProfile(ctx, p.AccountID, map[string]any{
		"name": p.Name, "nickname": p.Nickname, "avatar": p.Avatar, "signature": p.Signature,
		"phone": p.Phone, "email": p.Email, "remark": p.Remark,
	})
}

// DeleteByAccountIDs 按账号 ID 批量删除资料。
func (r *Repo) DeleteByAccountIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.with(ctx).Where("account_id IN ?", ids).Delete(&Profile{}).Error
}

// GetAccountPassword 查账号密码哈希。
func (r *Repo) GetAccountPassword(ctx context.Context, accountID string) (*AccountPassword, error) {
	var acc AccountPassword
	if err := r.with(ctx).First(&acc, "id = ?", accountID).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}

// UpdateAccountPassword 更新密码哈希。
func (r *Repo) UpdateAccountPassword(ctx context.Context, accountID, hash string) error {
	return r.with(ctx).Model(&AccountPassword{}).Where("id = ?", accountID).
		Update("password_hash", hash).Error
}
