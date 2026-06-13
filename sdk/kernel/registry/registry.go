package registry

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/auth/middleware"
)

// RouteRegistrar is a function that registers routes on the given gin engine.
type RouteRegistrar func(r *gin.Engine)

type Snapshot struct {
	Routes      []string `json:"routes"`
	Middlewares []string `json:"middlewares"`
	Frozen      bool     `json:"frozen"`
}

type routeRegistry struct {
	mu              sync.RWMutex
	routeRegistrars []namedRoute
	routeIndex      map[string]struct{}
	middlewareRegs  []namedMiddleware
	middlewareIndex map[string]struct{}
	frozen          bool
}

type namedRoute struct {
	name string
	reg  RouteRegistrar
}

type namedMiddleware struct {
	name string
	reg  func(r *gin.Engine)
}

var global = newRouteRegistry()
var registeredPerms sync.Map

func newRouteRegistry() *routeRegistry {
	return &routeRegistry{
		routeIndex:      make(map[string]struct{}),
		middlewareIndex: make(map[string]struct{}),
	}
}

// RegisterRoute registers a route registrar function.
func RegisterRoute(reg RouteRegistrar) {
	name := runtimeName(reg)
	if err := global.registerRoute(name, reg); err != nil {
		panic(err)
	}
}

// ExecuteRoutes runs all registered route registrars.
func ExecuteRoutes(r *gin.Engine) {
	global.executeRoutes(r)
}

// RegisterMiddleware registers a global middleware function.
func RegisterMiddleware(reg func(r *gin.Engine)) {
	name := runtimeName(reg)
	if err := global.registerMiddleware(name, reg); err != nil {
		panic(err)
	}
}

// ApplyMiddlewares runs all registered global middleware registrars.
func ApplyMiddlewares(r *gin.Engine) {
	global.applyMiddlewares(r)
}

func Freeze() {
	global.freeze()
}

func SnapshotState() Snapshot {
	return global.snapshot()
}

func ResetForTest() {
	global = newRouteRegistry()
	registeredPerms = sync.Map{}
}

// Perm registers a permission entry and returns a permission-checking middleware.
func Perm(code, name string) gin.HandlerFunc {
	if _, loaded := registeredPerms.LoadOrStore(code, true); !loaded {
		module := moduleFromCode(code)
		auth.RegisterPermission(auth.PermissionEntry{
			Code:   code,
			Module: module,
			Name:   name,
		})
	}
	return middleware.CheckPermission(auth.Business, []string{code})
}

// ClientPerm registers a permission entry and returns a permission-checking middleware for CONSUMER.
func ClientPerm(code, name string) gin.HandlerFunc {
	if _, loaded := registeredPerms.LoadOrStore(code, true); !loaded {
		module := moduleFromCode(code)
		auth.RegisterPermission(auth.PermissionEntry{
			Code:   code,
			Module: module,
			Name:   name,
		})
	}
	return middleware.CheckPermission(auth.Consumer, []string{code})
}

func moduleFromCode(code string) string {
	parts := strings.Split(code, ":")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ":")
	}
	return code
}

func runtimeName(fn any) string {
	if fn == nil {
		return "<nil>"
	}
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return fmt.Sprintf("%T", fn)
	}
	if rf := runtime.FuncForPC(v.Pointer()); rf != nil {
		return rf.Name()
	}
	return fmt.Sprintf("%T@%p", fn, fn)
}

func (r *routeRegistry) registerRoute(name string, reg RouteRegistrar) error {
	if reg == nil {
		return fmt.Errorf("route register failed: nil registrar")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.frozen {
		return fmt.Errorf("route register failed: registry frozen, route=%s", name)
	}
	if _, exists := r.routeIndex[name]; exists {
		return fmt.Errorf("route register failed: duplicate registrar=%s", name)
	}

	r.routeIndex[name] = struct{}{}
	r.routeRegistrars = append(r.routeRegistrars, namedRoute{name: name, reg: reg})
	return nil
}

func (r *routeRegistry) executeRoutes(engine *gin.Engine) {
	r.freeze()

	r.mu.RLock()
	routes := append([]namedRoute(nil), r.routeRegistrars...)
	r.mu.RUnlock()

	for _, item := range routes {
		item.reg(engine)
	}
}

func (r *routeRegistry) registerMiddleware(name string, reg func(r *gin.Engine)) error {
	if reg == nil {
		return fmt.Errorf("middleware register failed: nil registrar")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.frozen {
		return fmt.Errorf("middleware register failed: registry frozen, middleware=%s", name)
	}
	if _, exists := r.middlewareIndex[name]; exists {
		return fmt.Errorf("middleware register failed: duplicate registrar=%s", name)
	}

	r.middlewareIndex[name] = struct{}{}
	r.middlewareRegs = append(r.middlewareRegs, namedMiddleware{name: name, reg: reg})
	return nil
}

func (r *routeRegistry) applyMiddlewares(engine *gin.Engine) {
	r.freeze()

	r.mu.RLock()
	items := append([]namedMiddleware(nil), r.middlewareRegs...)
	r.mu.RUnlock()

	for _, item := range items {
		item.reg(engine)
	}
}

func (r *routeRegistry) freeze() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.frozen = true
}

func (r *routeRegistry) snapshot() Snapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routes := make([]string, 0, len(r.routeRegistrars))
	for _, item := range r.routeRegistrars {
		routes = append(routes, item.name)
	}
	mws := make([]string, 0, len(r.middlewareRegs))
	for _, item := range r.middlewareRegs {
		mws = append(mws, item.name)
	}

	return Snapshot{
		Routes:      routes,
		Middlewares: mws,
		Frozen:      r.frozen,
	}
}
