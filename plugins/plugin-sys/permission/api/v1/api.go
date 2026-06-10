package v1

import (
		"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	"hei-gin/plugins/plugin-sys/permission"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all permission routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/permission/modules
	r.GET("/api/v1/sys/permission/modules", registry.Perm("sys:permission:modules", "权限模块列表"), permListModules)
	// GET /api/v1/sys/permission/by-module
	r.GET("/api/v1/sys/permission/by-module", registry.Perm("sys:permission:by-module", "按模块查询权限"), permByModule)
}

// permListModules handles GET /api/v1/sys/permission/modules
func permListModules(c *gin.Context) {
	modules := permission.ListModules(c)
	result.Success(c, modules)
}

// permByModule handles GET /api/v1/sys/permission/by-module
func permByModule(c *gin.Context) {
	module := c.Query("module")
	perms := permission.ListByModule(c, module)
	result.Success(c, perms)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
