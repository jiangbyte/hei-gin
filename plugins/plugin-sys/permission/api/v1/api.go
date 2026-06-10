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
	r.GET("/api/v1/sys/permission/modules",
		registry.Perm("sys:permission:modules", "权限模块列表"),
		listModulesHandler,
	)
	// GET /api/v1/sys/permission/by-module
	r.GET("/api/v1/sys/permission/by-module",
		registry.Perm("sys:permission:by-module", "按模块查询权限"),
		byModuleHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// listModulesHandler handles GET /api/v1/sys/permission/modules
func listModulesHandler(c *gin.Context) {
	result.Success(c, permission.PermissionListModules(c))
}

// byModuleHandler handles GET /api/v1/sys/permission/by-module
func byModuleHandler(c *gin.Context) {
	result.Success(c, permission.PermissionListByModule(c, c.Query("module")))
}
