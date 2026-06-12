package provider

import (
	"hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/infra/db"
)

type UserProvider struct{}

func (p *UserProvider) GetUserNameByID(id string) string {
	var u user.SysUser
	if err := db.DB.First(&u, "id = ?", id).Error; err != nil {
		return id
	}
	if u.Nickname == nil {
		return id
	}
	return *u.Nickname
}
