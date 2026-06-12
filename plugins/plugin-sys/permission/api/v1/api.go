package v1

import (
	"hei-gin/plugins/plugin-sys/permission"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/web/result"

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
// @Summary      权限管理模块列表
// @Description  访问 /api/v1/sys/permission/modules，权限管理模块列表
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/permission/modules [get]
func listModulesHandler(c *gin.Context) {
	result.Success(c, permission.PermissionListModules(c))
}

// byModuleHandler handles GET /api/v1/sys/permission/by-module
// @Summary      权限管理按模块查询
// @Description  访问 /api/v1/sys/permission/by-module，权限管理按模块查询
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Param        module  query  string  false  "module"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/permission/by-module [get]
func byModuleHandler(c *gin.Context) {
	result.Success(c, permission.PermissionListByModule(c, c.Query("module")))
}
