package session

import (
	"context"

	"gorm.io/gorm"

	userModel "hei-gin/plugins/plugin-sys/user"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) LoadUsers(ctx context.Context, userIDs []string) map[string]*userModel.SysUser {
	if len(userIDs) == 0 {
		return map[string]*userModel.SysUser{}
	}
	var users []userModel.SysUser
	r.db.WithContext(ctx).Where("id IN ?", userIDs).Find(&users)
	result := make(map[string]*userModel.SysUser, len(users))
	for i := range users {
		user := users[i]
		result[user.ID] = &user
	}
	return result
}

func (r *repository) FindUserIDs(ctx context.Context, keyword string, limit int) []string {
	like := keyword + "%"
	var ids []string
	r.db.WithContext(ctx).Model(&userModel.SysUser{}).
		Select("id").
		Where("id = ? OR id LIKE ? OR username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?",
			keyword, like, like, like, like, like).
		Order("last_login_at DESC, id ASC").
		Limit(limit).
		Pluck("id", &ids)
	return ids
}
