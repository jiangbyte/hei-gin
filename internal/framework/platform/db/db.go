// Package db 提供 GORM 打开、审计字段混入与事务辅助。
//
// Author: Charlie
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

// Model 为共享审计字段混入。
//
// Author: Charlie
type Model struct {
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// OwnerDept 行级归属部门字段混入。
//
// Author: Charlie
type OwnerDept struct {
	OwnerDeptID *string `gorm:"column:owner_dept_id;size:64;index" json:"owner_dept_id"`
}

// Open 按配置打开 Postgres 或 sqlite: 前缀的 SQLite，并注册审计回调。
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
	if schemaHasField(db, "created_by") {
		db.Statement.SetColumn("created_by", aid)
	}
	if schemaHasField(db, "updated_by") {
		db.Statement.SetColumn("updated_by", aid)
	}
}

func auditUpdate(db *gorm.DB) {
	if db.Statement.Context == nil {
		return
	}
	aid := contextx.AccountID(db.Statement.Context)
	if aid == "" {
		return
	}
	if schemaHasField(db, "updated_by") {
		db.Statement.SetColumn("updated_by", aid)
	}
}

// schemaHasField 判断当前语句目标模型/表是否有指定列（避免对无审计列表如 sys_account_password_history 报 invalid field）。
func schemaHasField(db *gorm.DB, column string) bool {
	schema := db.Statement.Schema
	if schema == nil {
		return false
	}
	for _, f := range schema.Fields {
		if f.DBName == column {
			return true
		}
	}
	return false
}

// WithTx 在事务中执行 fn。
func WithTx(ctx context.Context, gdb *gorm.DB, fn func(tx *gorm.DB) error) error {
	return gdb.WithContext(ctx).Transaction(fn)
}
