// internal/framework/platform/module/module.go 模块注册表。
//
// Author: Charlie

package module

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/audit"
	"hei-gin/internal/framework/platform/notify"
	"hei-gin/internal/framework/platform/runtimecfg"
	"hei-gin/internal/framework/platform/storage"
)

// RouteRegistrar 在 /api 下挂载路由（完整路径写在 handler 上，如 /v1/admin/...）。
//
// Author: Charlie
type RouteRegistrar func(api *gin.RouterGroup)

// Schedule 保留字段名兼容；新任务请用 Jobs（SnailJob Handler）。
//
// Deprecated: 使用 Jobs。
//
// Author: Charlie
type Schedule struct {
	Name     string
	Interval string
	Run      func(ctx context.Context) error
}

// Job 是 SnailJob 执行器 Handler（Name 须与控制台 executor_info 一致）。
//
// Author: Charlie
type Job struct {
	Name string
	Run  func(ctx context.Context, param string) error
}

// EventHandler 订阅平台生命周期 / 领域事件。
//
// Author: Charlie
type EventHandler struct {
	Event   string
	Handler func(ctx context.Context, payload any) error
}

// Module 为插件契约。官方模块经 Register 注册；勿做扫盘发现。
//
// Author: Charlie
type Module struct {
	Name          string
	Order         int
	Routes        []RouteRegistrar
	Models        []any
	Jobs          []Job
	Schedules     []Schedule // Deprecated: 请用 Jobs
	OnStart       []func(ctx context.Context) error
	OnStop        []func(ctx context.Context) error
	EventHandlers []EventHandler
}

// Deps 是传给各模块构造器的运行时依赖图。
//
// Author: Charlie
type Deps struct {
	Cfg      *config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	Sessions *security.SessionStore
	Perms    *security.PermissionRegistry
	Storage  *storage.Manager
	Notify   *notify.Facade
	Audit    *audit.Queue
	Runtime  *runtimecfg.Settings
	services map[string]any
}

// Provide 向依赖袋放入命名服务（供跨模块接线，如 account_finder）。
func (d *Deps) Provide(name string, v any) {
	if d.services == nil {
		d.services = make(map[string]any)
	}
	d.services[name] = v
}

// Service 取出命名服务。
func (d *Deps) Service(name string) (any, bool) {
	if d.services == nil {
		return nil, false
	}
	v, ok := d.services[name]
	return v, ok
}

// Builder 根据 Deps 构造 Module。
//
// Author: Charlie
type Builder func(d *Deps) Module

type registration struct {
	name  string
	order int
	build Builder
}

var (
	regMu   sync.Mutex
	regList []registration
)

// Register 在业务包 init 中调用，登记模块构造器。
func Register(name string, order int, build Builder) {
	regMu.Lock()
	defer regMu.Unlock()
	regList = append(regList, registration{name: name, order: order, build: build})
}

// BuildAll 按 order 调用全部已注册构造器，再按 disabled/enabled 过滤。
func BuildAll(d *Deps, disabled, enabledOnly []string) *Registry {
	regMu.Lock()
	list := append([]registration{}, regList...)
	regMu.Unlock()

	// 按 order、name 排序后依次 Build（保证 account 先于 auth 提供服务）
	for i := 1; i < len(list); i++ {
		j := i
		for j > 0 && (list[j].order < list[j-1].order || (list[j].order == list[j-1].order && list[j].name < list[j-1].name)) {
			list[j], list[j-1] = list[j-1], list[j]
			j--
		}
	}

	var mods []Module
	for _, r := range list {
		m := r.build(d)
		if m.Name == "" {
			m.Name = r.name
		}
		if m.Order == 0 {
			m.Order = r.order
		}
		mods = append(mods, m)
	}
	return NewRegistry(mods, disabled, enabledOnly)
}

// Registry 是过滤排序后的模块列表。
//
// Author: Charlie
type Registry struct {
	Modules []Module
}

// NewRegistry 按名称过滤并排序。
func NewRegistry(mods []Module, disabled, enabledOnly []string) *Registry {
	dis := map[string]struct{}{}
	for _, n := range disabled {
		dis[n] = struct{}{}
	}
	en := map[string]struct{}{}
	for _, n := range enabledOnly {
		en[n] = struct{}{}
	}
	var out []Module
	for _, m := range mods {
		if _, ok := dis[m.Name]; ok {
			continue
		}
		if len(en) > 0 {
			if _, ok := en[m.Name]; !ok {
				continue
			}
		}
		out = append(out, m)
	}
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && (out[j].Order < out[j-1].Order || (out[j].Order == out[j-1].Order && out[j].Name < out[j-1].Name)) {
			out[j], out[j-1] = out[j-1], out[j]
			j--
		}
	}
	return &Registry{Modules: out}
}

// AllModels 汇总全部模块的 GORM 模型。
func (r *Registry) AllModels() []any {
	var models []any
	for _, m := range r.Modules {
		models = append(models, m.Models...)
	}
	return models
}

// MountRoutes 依次挂载各模块路由。
func (r *Registry) MountRoutes(api *gin.RouterGroup) {
	for _, m := range r.Modules {
		for _, reg := range m.Routes {
			reg(api)
		}
	}
}

// RunStart 按模块顺序执行 OnStart 钩子。
func (r *Registry) RunStart(ctx context.Context) error {
	for _, m := range r.Modules {
		for _, h := range m.OnStart {
			if err := h(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

// RunStop 逆序执行 OnStop 钩子。
func (r *Registry) RunStop(ctx context.Context) error {
	for i := len(r.Modules) - 1; i >= 0; i-- {
		m := r.Modules[i]
		for j := len(m.OnStop) - 1; j >= 0; j-- {
			if err := m.OnStop[j](ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
