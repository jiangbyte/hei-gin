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
	"hei-gin/internal/framework/platform/storage"
)

// RouteRegistrar åœ¨ /api ä¸‹æŒ‚è½½è·¯ç”±ï¼ˆå®Œæ•´è·¯å¾„å†™åœ¨ handler ä¸Šï¼Œå¦‚ /v1/admin/...ï¼‰ã€‚
//
// Author: Charlie
type RouteRegistrar func(api *gin.RouterGroup)

// Schedule ä¿ç•™å­—æ®µåå…¼å®¹ï¼›æ–°ä»»åŠ¡è¯·ç”¨ Jobsï¼ˆSnailJob Handlerï¼‰ã€‚
//
// Deprecated: ä½¿ç”¨ Jobsã€‚
//
// Author: Charlie
type Schedule struct {
	Name     string
	Interval string
	Run      func(ctx context.Context) error
}

// Job æ˜¯ SnailJob æ‰§è¡Œå™¨ Handlerï¼ˆName é¡»ä¸ŽæŽ§åˆ¶å° executor_info ä¸€è‡´ï¼‰ã€‚
//
// Author: Charlie
type Job struct {
	Name string
	Run  func(ctx context.Context, param string) error
}

// EventHandler è®¢é˜…å¹³å°ç”Ÿå‘½å‘¨æœŸ / é¢†åŸŸäº‹ä»¶ã€‚
//
// Author: Charlie
type EventHandler struct {
	Event   string
	Handler func(ctx context.Context, payload any) error
}

// Module ä¸ºæ’ä»¶å¥‘çº¦ã€‚å®˜æ–¹æ¨¡å—ç» Register æ³¨å†Œï¼›å‹¿åšæ‰«ç›˜å‘çŽ°ã€‚
//
// Author: Charlie
type Module struct {
	Name          string
	Order         int
	Routes        []RouteRegistrar
	Models        []any
	Jobs          []Job
	Schedules     []Schedule // Deprecated: è¯·ç”¨ Jobs
	OnStart       []func(ctx context.Context) error
	OnStop        []func(ctx context.Context) error
	EventHandlers []EventHandler
}

// Deps æ˜¯ä¼ ç»™å„æ¨¡å—æž„é€ å™¨çš„è¿è¡Œæ—¶ä¾èµ–å›¾ã€‚
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
	services map[string]any
}

// Provide å‘ä¾èµ–è¢‹æ”¾å…¥å‘½åæœåŠ¡ï¼ˆä¾›è·¨æ¨¡å—æŽ¥çº¿ï¼Œå¦‚ account_finderï¼‰ã€‚
func (d *Deps) Provide(name string, v any) {
	if d.services == nil {
		d.services = make(map[string]any)
	}
	d.services[name] = v
}

// Service å–å‡ºå‘½åæœåŠ¡ã€‚
func (d *Deps) Service(name string) (any, bool) {
	if d.services == nil {
		return nil, false
	}
	v, ok := d.services[name]
	return v, ok
}

// Builder æ ¹æ® Deps æž„é€  Moduleã€‚
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

// Register åœ¨ä¸šåŠ¡åŒ… init ä¸­è°ƒç”¨ï¼Œç™»è®°æ¨¡å—æž„é€ å™¨ã€‚
func Register(name string, order int, build Builder) {
	regMu.Lock()
	defer regMu.Unlock()
	regList = append(regList, registration{name: name, order: order, build: build})
}

// BuildAll æŒ‰ order è°ƒç”¨å…¨éƒ¨å·²æ³¨å†Œæž„é€ å™¨ï¼Œå†æŒ‰ disabled/enabled è¿‡æ»¤ã€‚
func BuildAll(d *Deps, disabled, enabledOnly []string) *Registry {
	regMu.Lock()
	list := append([]registration{}, regList...)
	regMu.Unlock()

	// æŒ‰ orderã€name æŽ’åºåŽä¾æ¬¡ Buildï¼ˆä¿è¯ account å…ˆäºŽ auth æä¾›æœåŠ¡ï¼‰
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

// Registry æ˜¯è¿‡æ»¤æŽ’åºåŽçš„æ¨¡å—åˆ—è¡¨ã€‚
//
// Author: Charlie
type Registry struct {
	Modules []Module
}

// NewRegistry æŒ‰åç§°è¿‡æ»¤å¹¶æŽ’åºã€‚
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

// AllModels æ±‡æ€»å…¨éƒ¨æ¨¡å—çš„ GORM æ¨¡åž‹ã€‚
func (r *Registry) AllModels() []any {
	var models []any
	for _, m := range r.Modules {
		models = append(models, m.Models...)
	}
	return models
}

// MountRoutes ä¾æ¬¡æŒ‚è½½å„æ¨¡å—è·¯ç”±ã€‚
func (r *Registry) MountRoutes(api *gin.RouterGroup) {
	for _, m := range r.Modules {
		for _, reg := range m.Routes {
			reg(api)
		}
	}
}

// RunStart æŒ‰æ¨¡å—é¡ºåºæ‰§è¡Œ OnStart é’©å­ã€‚
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

// RunStop é€†åºæ‰§è¡Œ OnStop é’©å­ã€‚
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
