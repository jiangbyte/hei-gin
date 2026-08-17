// Package app 是应用装配根：基础设施、自注册模块、HTTP 与内嵌任务调度器。
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
	"hei-gin/internal/framework/core/crypto"
	"hei-gin/internal/framework/core/logger"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/audit"
	"hei-gin/internal/framework/platform/cache"
	"hei-gin/internal/framework/platform/db"
	"hei-gin/internal/framework/platform/gojob"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/notify"
	"hei-gin/internal/framework/platform/otel"
	"hei-gin/internal/framework/platform/runtimecfg"
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
	Audit    *audit.Queue
	Notify   *notify.Facade
	Runtime  *runtimecfg.Settings
	Modules  *module.Registry
	Jobs     *gojob.Manager
}

// API 单体进程：HTTP + 模块钩子 + 内嵌任务调度器。
//
// Author: Charlie
type API struct {
	Deps   *Deps
	Engine *gin.Engine
	Server *http.Server
	Audit  *audit.Queue
	Jobs   *gojob.Manager
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
	store := storage.NewManager()
	if err := otel.Init(cfg.OTel); err != nil {
		return nil, err
	}
	nf := notify.NewFacade(cfg.Notify, gdb)
	rt := runtimecfg.New(gdb)
	if codec, err := crypto.NewFernetFromConfig(cfg.Crypto.FernetKey, cfg.Crypto.VaultAddr); err == nil {
		rt.WithCodec(codec)
	}
	store.SetRuntime(rt)
	return &Deps{
		Cfg:      cfg,
		DB:       gdb,
		Redis:    rdb,
		Sessions: security.NewSessionStore(rdb),
		Perms:    security.NewPermissionRegistry(rdb),
		Storage:  store,
		Audit:  audit.NewQueue(gdb, rdb, cfg.Audit),
		Notify: nf,
		Runtime:  rt,
		// 任务调度器（handlers 在 NewAPI 装配完成后填充）
		Jobs: gojob.NewManager(gdb, rdb, gojob.Config{
			ScanIntervalMS:   cfg.Job.ScanIntervalMS,
			PoolSize:         cfg.Job.PoolSize,
			LogRetentionDays: cfg.Job.LogRetentionDays,
			LogBatchSize:     cfg.Job.LogBatchSize,
		}, nil),
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
	// 操作审计由各路由挂载 middleware.OperationAudit(d.Audit, resourceType, action)
	d.Modules.MountRoutes(api)

	if d.Cfg.App.Debug {
		r.GET("/api/v1/internal/debug/routes", func(c *gin.Context) {
			type routeInfo struct {
				Method string `json:"method"`
				Path   string `json:"path"`
			}
			out := make([]routeInfo, 0, len(r.Routes()))
			for _, rt := range r.Routes() {
				out = append(out, routeInfo{Method: rt.Method, Path: rt.Path})
			}
			response.OK(c, out)
		})
	}

	srv := &http.Server{
		Addr:              d.Cfg.Addr(),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
	// 模块装配完成后填充任务处理器
	if d.Jobs != nil {
		d.Jobs.SetHandlers(collectHandlers(d.Modules))
	}
	return &API{
		Deps:   d,
		Engine: r,
		Server: srv,
		Audit:  d.Audit,
		Jobs:   d.Jobs,
	}
}

// collectHandlers 从模块注册表提取任务处理器（module.Job → gojob.HandlerDef）。
func collectHandlers(regs *module.Registry) []gojob.HandlerDef {
	if regs == nil {
		return nil
	}
	var handlers []gojob.HandlerDef
	for _, mod := range regs.Modules {
		for _, j := range mod.Jobs {
			j := j
			handlers = append(handlers, gojob.HandlerDef{Key: j.Name, Name: j.Name, Run: j.Run})
		}
	}
	return handlers
}

// Start 启动审计队列、模块钩子、任务调度器与 HTTP 监听。
func (a *API) Start(ctx context.Context) error {
	a.Audit.Start(ctx)
	if err := a.Deps.Modules.RunStart(ctx); err != nil {
		return err
	}
	if err := a.Deps.Perms.Sync(ctx); err != nil {
		logger.L.Warn("权限注册表同步失败", zap.Error(err))
	}
	if a.Jobs != nil {
		if err := a.Jobs.Start(ctx); err != nil {
			return err
		}
		logger.L.Info("任务调度器已启动")
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
	if a.Jobs != nil {
		_ = a.Jobs.Stop(ctx)
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
