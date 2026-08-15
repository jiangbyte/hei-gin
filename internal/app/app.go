// Package app 是应用装配根：基础设施、自注册模块、HTTP 与 SnailJob 执行器。
//
// 默认 blank import internal/modules/all；复杂场景可直接改 framework。
//
// Author: Charlie
package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/logger"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/audit"
	"hei-gin/internal/framework/platform/cache"
	"hei-gin/internal/framework/platform/db"
	"hei-gin/internal/framework/platform/events"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/notify"
	"hei-gin/internal/framework/platform/otel"
	"hei-gin/internal/framework/platform/runtimecfg"
	"hei-gin/internal/framework/platform/snailjob"
	"hei-gin/internal/framework/platform/storage"
)

// Deps 应用进程级依赖。
//
// Author: Charlie
type Deps struct {
	Cfg      *config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	Sessions *security.SessionStore
	Perms    *security.PermissionRegistry
	Storage  *storage.Manager
	Events   *events.Bus
	Audit    *audit.Queue
	Notify   *notify.Facade
	Runtime  *runtimecfg.Settings
	Modules  *module.Registry
}

// API 单体进程：HTTP + 模块钩子 + SnailJob 执行器。
//
// Author: Charlie
type API struct {
	Deps     *Deps
	Engine   *gin.Engine
	Server   *http.Server
	Audit    *audit.Queue
	SnailJob *snailjob.Manager
}

// OpenInfra 连接 DB/Redis/存储，准备空 Deps（随后 AttachRegisteredModules）。
func OpenInfra(cfg *config.Config) (*Deps, error) {
	if err := logger.Setup(cfg.App.Debug); err != nil {
		return nil, err
	}
	if err := idgen.Init(cfg.IDGen.WorkerID, cfg.IDGen.DatacenterID); err != nil {
		return nil, err
	}
	gdb, err := db.Open(cfg.DB)
	if err != nil {
		return nil, err
	}
	rdb, err := cache.Open(cfg.Redis)
	if err != nil {
		return nil, err
	}
	store, err := storage.NewManager(cfg.Storage)
	if err != nil {
		return nil, err
	}
	if err := otel.Init(cfg.OTel); err != nil {
		return nil, err
	}
	nf := notify.NewFacade(cfg.Notify, gdb)
	return &Deps{
		Cfg:      cfg,
		DB:       gdb,
		Redis:    rdb,
		Sessions: security.NewSessionStore(rdb),
		Perms:    security.NewPermissionRegistry(rdb),
		Storage:  store,
		Events:   events.NewBus(),
		Audit:    audit.NewQueue(gdb, rdb, cfg.Audit),
		Notify:   nf,
		Runtime:  runtimecfg.New(gdb),
	}, nil
}

// NewAPI 构建 Gin 引擎与 HTTP Server。
func NewAPI(d *Deps) *API {
	if d.Cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.SecurityHeaders(d.Cfg.Security))
	r.Use(middleware.AccessLog())
	r.Use(middleware.AuthContext(d.Cfg.Auth, d.Sessions))
	r.Use(middleware.AuthWhitelist(d.Cfg.Auth.AuthWhitelist))
	r.Use(middleware.Trace())
	r.Use(middleware.CORS(d.Cfg.CORS))
	r.Use(middleware.ErrorHandler())

	r.GET("/", func(c *gin.Context) {
		response.OK(c, gin.H{"name": d.Cfg.App.Name, "status": "ok"})
	})
	if d.Cfg.Metrics.Enabled {
		path := d.Cfg.Metrics.Path
		if path == "" {
			path = "/metrics"
		}
		r.GET(path, middleware.PrometheusHandler())
	}

	api := r.Group("/api")
	d.Modules.MountRoutes(api)

	srv := &http.Server{
		Addr:              d.Cfg.Addr(),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return &API{
		Deps:     d,
		Engine:   r,
		Server:   srv,
		Audit:    d.Audit,
		SnailJob: snailjob.NewManager(d.Cfg.SnailJob, d.Modules),
	}
}

// Start 启动审计队列、模块钩子、SnailJob 执行器与 HTTP 监听。
func (a *API) Start(ctx context.Context) error {
	a.Audit.Start(ctx)
	if err := a.Deps.Modules.RunStart(ctx); err != nil {
		return err
	}
	if err := a.Deps.Events.Emit(ctx, events.OnDBReady, nil); err != nil {
		return err
	}
	if err := a.Deps.Perms.Sync(ctx); err != nil {
		logger.L.Warn("权限注册表同步失败", zap.Error(err))
	} else {
		_ = a.Deps.Events.Emit(ctx, events.OnPermissionsSynced, nil)
	}
	if a.Deps.Cfg.SnailJob.Enabled {
		if err := a.SnailJob.Start(); err != nil {
			return err
		}
	}
	logger.L.Info("api 正在启动", zap.String("addr", a.Server.Addr))
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.Server.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	case <-time.After(200 * time.Millisecond):
		logger.L.Info("api 已监听", zap.String("addr", a.Server.Addr))
	}
	go func() {
		if err := <-errCh; err != nil && err != http.ErrServerClosed {
			logger.L.Fatal("监听失败", zap.Error(err))
		}
	}()
	return nil
}

// Stop 优雅关闭。
func (a *API) Stop(ctx context.Context) error {
	if a.SnailJob != nil {
		_ = a.SnailJob.Stop(ctx)
	}
	_ = a.Deps.Modules.RunStop(ctx)
	a.Audit.Stop()
	err := a.Server.Shutdown(ctx)
	_ = a.Deps.Redis.Close()
	sqlDB, e := a.Deps.DB.DB()
	if e == nil {
		_ = sqlDB.Close()
	}
	logger.Sync()
	return err
}

// CloseIdle 关闭空闲连接。
func CloseIdle(d *Deps) error {
	if d == nil {
		return nil
	}
	var err error
	if d.Redis != nil {
		err = d.Redis.Close()
	}
	if d.DB != nil {
		sqlDB, e := d.DB.DB()
		if e == nil {
			if e2 := sqlDB.Close(); e2 != nil && err == nil {
				err = e2
			}
		}
	}
	return err
}

// LoadOrDie 加载配置，失败则 panic。
func LoadOrDie(path string) *config.Config {
	cfg, err := config.Load(path)
	if err != nil {
		panic(fmt.Errorf("config: %w", err))
	}
	return cfg
}
