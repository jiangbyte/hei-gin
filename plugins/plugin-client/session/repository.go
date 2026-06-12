package session

import (
	"context"

	cliUser "hei-gin/plugins/plugin-client/user"
	"hei-gin/sdk/db"
)

func LoadUsers(ctx context.Context, userIDs []string) map[string]*cliUser.ClientUser {
	if len(userIDs) == 0 {
		return map[string]*cliUser.ClientUser{}
	}
	var users []cliUser.ClientUser
	db.DB.WithContext(ctx).Where("id IN ?", userIDs).Find(&users)
	result := make(map[string]*cliUser.ClientUser, len(users))
	for i := range users {
		user := users[i]
		result[user.ID] = &user
	}
	return result
}

func FindUserIDs(ctx context.Context, keyword string, limit int) []string {
	like := keyword + "%"
	var ids []string
	db.DB.WithContext(ctx).Model(&cliUser.ClientUser{}).
		Select("id").
		Where("id = ? OR id LIKE ? OR username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?",
			keyword, like, like, like, like, like).
		Order("last_login_at DESC, id ASC").
		Limit(limit).
		Pluck("id", &ids)
	return ids
}
