package registry

import (
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/auth/middleware"
)

// ===== Route Registry =====

// RouteRegistrar is a function that registers routes on the given gin engine.
type RouteRegistrar func(r *gin.Engine)

var routeRegistrars []RouteRegistrar

// RegisterRoute registers a route registrar function.
func RegisterRoute(reg RouteRegistrar) {
	routeRegistrars = append(routeRegistrars, reg)
}

// ExecuteRoutes runs all registered route registrars.
func ExecuteRoutes(r *gin.Engine) {
	for _, reg := range routeRegistrars {
		reg(r)
	}
}

// ===== Middleware Registry =====

var middlewareRegs []func(r *gin.Engine)

// RegisterMiddleware registers a global middleware function.
func RegisterMiddleware(reg func(r *gin.Engine)) {
	middlewareRegs = append(middlewareRegs, reg)
}

// ApplyMiddlewares runs all registered global middleware registrars.
func ApplyMiddlewares(r *gin.Engine) {
	for _, reg := range middlewareRegs {
		reg(r)
	}
}

// ===== Permission Registry =====

var registeredPerms sync.Map

// Perm registers a permission entry and returns a permission-checking middleware.
// code: permission code, e.g. "sys:banner:page"
// name: human-readable name, e.g. "横幅分页"
func Perm(code, name string) gin.HandlerFunc {
	if _, loaded := registeredPerms.LoadOrStore(code, true); !loaded {
		module := moduleFromCode(code)
		auth.RegisterPermission(auth.PermissionEntry{
			Code:   code,
			Module: module,
			Name:   name,
		})
	}
	return middleware.HeiCheckPermission([]string{code})
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
	return middleware.HeiClientCheckPermission([]string{code})
}

func moduleFromCode(code string) string {
	parts := strings.Split(code, ":")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ":")
	}
	return code
}
