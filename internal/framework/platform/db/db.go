// Package db 提供 GORM 打开、审计字段混入与事务辅助。
//
// Author: Charlie
package db

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
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

// Open 按 db.driver 或 DSN scheme 打开 MySQL / PostgreSQL，并注册审计回调。
func Open(cfg config.DBConfig) (*gorm.DB, error) {
	driver, dsn, err := resolveDriverDSN(cfg)
	if err != nil {
		return nil, err
	}
	var dialector gorm.Dialector
	switch driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported db driver: %s (use mysql or postgres)", driver)
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

// resolveDriverDSN 解析 driver 与最终传给驱动的 DSN。
// driver 优先用 cfg.Driver；否则从 URL scheme 推断（mysql / postgres|postgresql）。
func resolveDriverDSN(cfg config.DBConfig) (driver, dsn string, err error) {
	raw := strings.TrimSpace(cfg.URL)
	if raw == "" {
		return "", "", fmt.Errorf("db.url is required")
	}
	driver = strings.ToLower(strings.TrimSpace(cfg.Driver))
	if driver == "" {
		driver, err = inferDriver(raw)
		if err != nil {
			return "", "", err
		}
	}
	switch driver {
	case "mysql", "postgres":
	default:
		return "", "", fmt.Errorf("unsupported db.driver %q (use mysql or postgres)", cfg.Driver)
	}
	dsn, err = normalizeDSN(driver, raw)
	if err != nil {
		return "", "", err
	}
	return driver, dsn, nil
}

func inferDriver(raw string) (string, error) {
	lower := strings.ToLower(raw)
	switch {
	case strings.HasPrefix(lower, "mysql://"), strings.HasPrefix(lower, "mysql:"):
		return "mysql", nil
	case strings.HasPrefix(lower, "postgres://"), strings.HasPrefix(lower, "postgresql://"):
		return "postgres", nil
	default:
		return "", fmt.Errorf("cannot infer db driver from url %q; set db.driver to mysql or postgres", raw)
	}
}

// normalizeDSN 将 URL 转为各驱动接受的 DSN。
// MySQL go-sql-driver 使用 user:pass@tcp(host:port)/db?params；也接受已是该格式的 DSN。
func normalizeDSN(driver, raw string) (string, error) {
	if driver == "postgres" {
		return raw, nil
	}
	// mysql
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "mysql://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", fmt.Errorf("parse mysql url: %w", err)
		}
		user := ""
		if u.User != nil {
			user = u.User.Username()
			if p, ok := u.User.Password(); ok {
				user += ":" + p
			}
		}
		host := u.Host
		if host == "" {
			host = "127.0.0.1:3306"
		}
		dbName := strings.TrimPrefix(u.Path, "/")
		q := u.RawQuery
		dsn := fmt.Sprintf("%s@tcp(%s)/%s", user, host, dbName)
		if q != "" {
			dsn += "?" + q
		} else {
			dsn += "?parseTime=true&loc=Local&charset=utf8mb4"
		}
		return dsn, nil
	}
	// 已是 go-sql-driver 风格或其它：原样返回
	return raw, nil
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
