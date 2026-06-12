package db

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"hei-gin/sdk/config"
	"hei-gin/sdk/utils"
)

var DB *gorm.DB

func InitDB() error {
	cfg := config.C.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.PoolSize + cfg.MaxOverflow)
	sqlDB.SetMaxIdleConns(cfg.PoolSize)

	maxLifetime := time.Duration(cfg.PoolRecycle) * time.Second
	if maxLifetime <= 0 || maxLifetime > 1*time.Hour {
		maxLifetime = 1 * time.Hour
	}
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxLifetime / 2)

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return fmt.Errorf("database ping failed: %w", err)
	}

	log.Printf("[Database] MySQL connection verified, max_conns=%d, max_lifetime=%v",
		cfg.PoolSize+cfg.MaxOverflow, maxLifetime)

	registerModelHooks()
	registerUpdateHook()

	return nil
}

type CtxKeyLoginID struct{}

func ContextWithLoginID(ctx context.Context, uid string) context.Context {
	if uid != "" {
		return context.WithValue(ctx, CtxKeyLoginID{}, uid)
	}
	return ctx
}

func registerModelHooks() {
	DB.Callback().Create().Before("gorm:before_create").Register("model_defaults", func(db *gorm.DB) {
		v := db.Statement.Dest
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Slice {
			for i := 0; i < rv.Len(); i++ {
				elem := rv.Index(i)
				if elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				if elem.Kind() == reflect.Struct {
					fillModelDefaults(elem, db.Statement.Context)
				}
			}
		} else if rv.Kind() == reflect.Struct {
			fillModelDefaults(rv, db.Statement.Context)
		}
	})
}

func registerUpdateHook() {
	DB.Callback().Update().Before("gorm:before_update").Register("model_update_defaults", func(db *gorm.DB) {
		now := time.Now()
		uid, _ := db.Statement.Context.Value(CtxKeyLoginID{}).(string)

		switch v := db.Statement.Dest.(type) {
		case map[string]interface{}:
			v["updated_at"] = now
			if uid != "" {
				v["updated_by"] = uid
			}
		default:
			rv := reflect.ValueOf(db.Statement.Dest)
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			if rv.Kind() == reflect.Struct {
				setField(rv, "UpdatedAt", reflect.ValueOf(&now))
				if uid != "" {
					setField(rv, "UpdatedBy", reflect.ValueOf(&uid))
				}
			}
		}
	})
}

func setField(v reflect.Value, name string, val reflect.Value) {
	f := v.FieldByName(name)
	if !f.IsValid() || !f.CanSet() {
		return
	}
	if f.Kind() == reflect.Ptr && f.IsNil() {
		f.Set(val)
	} else if f.Kind() != reflect.Ptr {
		f.Set(val.Elem())
	}
}

func fillModelDefaults(v reflect.Value, ctx context.Context) {
	now := time.Now()
	uid, _ := ctx.Value(CtxKeyLoginID{}).(string)

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}
		name := v.Type().Field(i).Name

		switch name {
		case "ID":
			if f.Kind() == reflect.String && f.String() == "" {
				f.SetString(utils.GenerateID())
			}
		case "CreatedAt":
			if f.Kind() == reflect.Ptr && f.Type().Elem() == reflect.TypeOf(time.Time{}) && f.IsNil() {
				f.Set(reflect.ValueOf(&now))
			}
		case "UpdatedAt":
			if f.Kind() == reflect.Ptr && f.Type().Elem() == reflect.TypeOf(time.Time{}) && f.IsNil() {
				f.Set(reflect.ValueOf(&now))
			}
		case "CreatedBy":
			if f.Kind() == reflect.Ptr && f.Type().Elem().Kind() == reflect.String && f.IsNil() && uid != "" {
				f.Set(reflect.ValueOf(&uid))
			}
		case "UpdatedBy":
			if f.Kind() == reflect.Ptr && f.Type().Elem().Kind() == reflect.String && f.IsNil() && uid != "" {
				f.Set(reflect.ValueOf(&uid))
			}
		}
	}
}

func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
		DB = nil
	}
}
