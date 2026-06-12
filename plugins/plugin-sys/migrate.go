package plugin_sys

import (
	"context"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"hei-gin/plugins/plugin-sys/banner"
	"hei-gin/plugins/plugin-sys/config"
	"hei-gin/plugins/plugin-sys/dict"
	"hei-gin/plugins/plugin-sys/file"
	"hei-gin/plugins/plugin-sys/group"
	"hei-gin/plugins/plugin-sys/home"
	sysLog "hei-gin/plugins/plugin-sys/log"
	"hei-gin/plugins/plugin-sys/notice"
	"hei-gin/plugins/plugin-sys/org"
	"hei-gin/plugins/plugin-sys/position"
	"hei-gin/plugins/plugin-sys/resource"
	"hei-gin/plugins/plugin-sys/role"
	"hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/utils"
)

func init() {
	db.RegisterModel(&user.SysUser{})
	db.RegisterModel(&user.RelUserRole{})
	db.RegisterModel(&user.RelUserPermission{})
	db.RegisterModel(&user.RelRolePermission{})
	db.RegisterModel(&user.RelRoleResource{})
	db.RegisterModel(&role.SysRole{})
	db.RegisterModel(&org.SysOrg{})
	db.RegisterModel(&group.SysGroup{})
	db.RegisterModel(&position.SysPosition{})
	db.RegisterModel(&dict.SysDict{})
	db.RegisterModel(&config.SysConfig{})
	db.RegisterModel(&banner.SysBanner{})
	db.RegisterModel(&home.SysQuickAction{})
	db.RegisterModel(&sysLog.SysLog{})
	db.RegisterModel(&notice.SysNotice{})
	db.RegisterModel(&file.SysFile{})
	db.RegisterModel(&resource.SysResource{})
	db.RegisterModel(&resource.SysModule{})

	db.RegisterSeed("admin user", seedAdminUser)
}

func seedAdminUser() error {
	var count int64
	db.DB.Model(&user.SysUser{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		log.Println("[Seed] Admin user already exists, skipped")
		return nil
	}

	now := time.Now()
	hashed, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := user.SysUser{
		ID:         utils.GenerateID(),
		Username:   strPtr("admin"),
		Password:   strPtr(string(hashed)),
		Nickname:   strPtr("超管"),
		Status:     string(enums.UserStatusActive),
		LoginCount: 0,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
	if err := db.DB.WithContext(context.TODO()).Create(&admin).Error; err != nil {
		return err
	}
	log.Println("[Seed] Admin user created (username: admin, password: 123456)")
	return nil
}

func strPtr(s string) *string { return &s }
