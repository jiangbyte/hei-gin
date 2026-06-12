package user

import (
	"context"

	"hei-gin/sdk/infra/db"
)

func Page(ctx context.Context, p *ClientUserPageParam) ([]ClientUser, int64) {
	q := db.DB.WithContext(ctx).Model(&ClientUser{})
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

func FindByID(ctx context.Context, id string) (*ClientUser, error) {
	var e ClientUser
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func CountByUsername(ctx context.Context, username string, excludeID string) int64 {
	var count int64
	q := db.DB.WithContext(ctx).Model(&ClientUser{}).Where("username = ?", username)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	q.Count(&count)
	return count
}

func Create(ctx context.Context, entity *ClientUser) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&ClientUser{}).Where("id = ?", id).Updates(up).Error
}

func UpdateAvatar(ctx context.Context, entity *ClientUser, avatar string) error {
	return db.DB.WithContext(ctx).Model(entity).Update("avatar", avatar).Error
}

func UpdatePassword(ctx context.Context, id string, password string) error {
	return db.DB.WithContext(ctx).Model(&ClientUser{}).Where("id = ?", id).Update("password", password).Error
}

func DeleteByIDs(ctx context.Context, ids []string) error {
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&ClientUser{}).Error
}
