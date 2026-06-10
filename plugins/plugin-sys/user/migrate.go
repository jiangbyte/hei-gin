package user

import (
	"context"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
)

func init() {
	db.RegisterModel(&SysUser{})
	db.RegisterModel(&RelUserRole{})
	db.RegisterModel(&RelUserPermission{})
	db.RegisterModel(&RelRolePermission{})
	db.RegisterModel(&RelRoleResource{})

	db.RegisterSeed("admin user", seedAdminUser)
}

func seedAdminUser() error {
	var count int64
	db.DB.Model(&SysUser{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		log.Println("[Seed] Admin user already exists, skipped")
		return nil
	}

	now := time.Now()
	hashed, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := SysUser{
		ID:         utils.GenerateID(),
		Username:   strPtr("admin"),
		Password:   strPtr(string(hashed)),
		Nickname:   strPtr("超管"),
		Status:     string(enums.UserStatusActive),
		LoginCount: 0,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
	if err := db.DB.WithContext(context.TODO()).Create(&user).Error; err != nil {
		return err
	}
	log.Println("[Seed] Admin user created (username: admin, password: 123456)")
	return nil
}

func strPtr(s string) *string { return &s }
