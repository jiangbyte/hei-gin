// Package db æä¾› GORM æ‰“å¼€ã€å®¡è®¡å­—æ®µæ··å…¥ä¸Žäº‹åŠ¡è¾…åŠ©ã€‚
package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"hei-gin/internal/framework/core/config"
	contextx "hei-gin/internal/framework/core/context"
)

// Model ä¸ºå…±äº«å®¡è®¡å­—æ®µæ··å…¥ã€‚
//
// Author: Charlie
type Model struct {
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// OwnerDept è¡Œçº§å½’å±žéƒ¨é—¨å­—æ®µæ··å…¥ã€‚
//
// Author: Charlie
type OwnerDept struct {
	OwnerDeptID *string `gorm:"column:owner_dept_id;size:64;index" json:"owner_dept_id"`
}

// Open æŒ‰é…ç½®æ‰“å¼€ Postgres æˆ– sqlite: å‰ç¼€çš„ SQLiteï¼Œå¹¶æ³¨å†Œå®¡è®¡å›žè°ƒã€‚
func Open(cfg config.DBConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector
	dsn := cfg.URL
	if len(dsn) >= 7 && dsn[:7] == "sqlite:" {
		dialector = sqlite.Open(dsn[7:])
	} else {
		dialector = postgres.Open(dsn)
	}
	level := logger.Silent
	if cfg.Echo {
		level = logger.Info
	}
	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger:                                   logger.Default.LogMode(level),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.PoolSize + cfg.MaxOverflow)
	sqlDB.SetMaxIdleConns(cfg.PoolSize)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	_ = gdb.Callback().Create().Before("gorm:create").Register("hei:audit_create", auditCreate)
	_ = gdb.Callback().Update().Before("gorm:update").Register("hei:audit_update", auditUpdate)
	return gdb, nil
}

func auditCreate(db *gorm.DB) {
	if db.Statement.Context == nil {
		return
	}
	aid := contextx.AccountID(db.Statement.Context)
	if aid == "" {
		return
	}
	db.Statement.SetColumn("created_by", aid)
	db.Statement.SetColumn("updated_by", aid)
}

func auditUpdate(db *gorm.DB) {
	if db.Statement.Context == nil {
		return
	}
	aid := contextx.AccountID(db.Statement.Context)
	if aid == "" {
		return
	}
	db.Statement.SetColumn("updated_by", aid)
}

// WithTx åœ¨äº‹åŠ¡ä¸­æ‰§è¡Œ fnã€‚
func WithTx(ctx context.Context, gdb *gorm.DB, fn func(tx *gorm.DB) error) error {
	return gdb.WithContext(ctx).Transaction(fn)
}
