package user

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) Page(ctx context.Context, p *ClientUserPageParam) ([]ClientUser, int64) {
	q := r.db.WithContext(ctx).Model(&ClientUser{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?", like, like, like, like)
	}
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}
	var total int64
	q.Count(&total)
	var rows []ClientUser
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) FindByID(ctx context.Context, id string) (*ClientUser, error) {
	var e ClientUser
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) CountByUsername(ctx context.Context, username string, excludeID string) int64 {
	var count int64
	q := r.db.WithContext(ctx).Model(&ClientUser{}).Where("username = ?", username)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	q.Count(&count)
	return count
}

func (r *repository) Create(ctx context.Context, entity *ClientUser) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&ClientUser{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) UpdateAvatar(ctx context.Context, entity *ClientUser, avatar string) error {
	return r.db.WithContext(ctx).Model(entity).Update("avatar", avatar).Error
}

func (r *repository) UpdatePassword(ctx context.Context, id string, password string) error {
	return r.db.WithContext(ctx).Model(&ClientUser{}).Where("id = ?", id).Update("password", password).Error
}

func (r *repository) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&ClientUser{}).Error
}
