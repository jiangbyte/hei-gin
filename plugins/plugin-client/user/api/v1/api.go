package v1

import (
	clientuser "hei-gin/plugins/plugin-client/user"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *clientuser.Service
}

var defaultHandler = newHandler(clientuser.DefaultModule)

func newHandler(module *clientuser.Module) *handler {
	return &handler{service: module.Service()}
}

func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/client-user/page",
		registry.Perm("client:user:page", "C端用户分页"),
		defaultHandler.page,
	)

	r.POST("/api/v1/client-user/create",
		registry.Perm("client:user:create", "添加C端用户"),
		defaultHandler.create,
	)

	r.POST("/api/v1/client-user/modify",
		registry.Perm("client:user:modify", "编辑C端用户"),
		defaultHandler.modify,
	)

	r.POST("/api/v1/client-user/remove",
		registry.Perm("client:user:remove", "删除C端用户"),
		defaultHandler.remove,
	)

	r.GET("/api/v1/client-user/detail",
		registry.Perm("client:user:detail", "C端用户详情"),
		defaultHandler.detail,
	)

	r.GET("/api/v1/c/client-user/current",
		middleware.HeiClientCheckLogin(),
		defaultHandler.current,
	)

	r.POST("/api/v1/c/client-user/update-profile",
		middleware.HeiClientCheckLogin(),
		log.SysLog("C端用户更新个人信息"),
		middleware.NoRepeat(3000),
		defaultHandler.updateProfile,
	)

	r.POST("/api/v1/c/client-user/update-avatar",
		middleware.HeiClientCheckLogin(),
		log.SysLog("C端用户更新头像"),
		defaultHandler.updateAvatar,
	)

	r.POST("/api/v1/c/client-user/update-password",
		middleware.HeiClientCheckLogin(),
		log.SysLog("C端用户修改密码"),
		middleware.NoRepeat(3000),
		defaultHandler.updatePassword,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/client-user/page
// @Summary      C端用户分页查询
// @Description  访问 /api/v1/client-user/page，C端用户分页查询
// @Tags         C端用户
// @Accept       json
// @Produce      json
// @Param        query  query  clientuser.ClientUserPageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client-user/page [get]
func (h *handler) page(c *gin.Context) {
	var param clientuser.ClientUserPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Page(c, &param)
}

// createHandler handles POST /api/v1/client-user/create
// @Summary      C端用户创建
// @Description  访问 /api/v1/client-user/create，C端用户创建
// @Tags         C端用户
// @Accept       json
// @Produce      json
// @Param        body  body  clientuser.ClientUserVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client-user/create [post]
func (h *handler) create(c *gin.Context) {
	var vo clientuser.ClientUserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Create(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/client-user/modify
// @Summary      C端用户修改
// @Description  访问 /api/v1/client-user/modify，C端用户修改
// @Tags         C端用户
// @Accept       json
// @Produce      json
// @Param        body  body  clientuser.ClientUserVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client-user/modify [post]
func (h *handler) modify(c *gin.Context) {
	var vo clientuser.ClientUserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Modify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/client-user/remove
// @Summary      C端用户删除
// @Description  访问 /api/v1/client-user/remove，C端用户删除
// @Tags         C端用户
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client-user/remove [post]
func (h *handler) remove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Remove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/client-user/detail
// @Summary      C端用户详情查询
// @Description  访问 /api/v1/client-user/detail，C端用户详情查询
// @Tags         C端用户
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client-user/detail [get]
func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
	result.Success(c, vo)
}

// currentHandler handles GET /api/v1/c/client-user/current
// @Summary      C端用户当前信息
// @Description  访问 /api/v1/c/client-user/current，C端用户当前信息
// @Tags         C端用户
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/client-user/current [get]
func (h *handler) current(c *gin.Context) {
	vo := h.service.Current(c)
	result.Success(c, vo)
}

// updateProfileHandler handles POST /api/v1/c/client-user/update-profile
// @Summary      C端用户更新个人信息
// @Description  访问 /api/v1/c/client-user/update-profile，C端用户更新个人信息
// @Tags         C端用户
// @Accept       json
// @Produce      json
// @Param        body  body  clientuser.UpdateProfileParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/client-user/update-profile [post]
func (h *handler) updateProfile(c *gin.Context) {
	var param clientuser.UpdateProfileParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.UpdateProfile(c, &param)
	result.Success(c, nil)
}

// updateAvatarHandler handles POST /api/v1/c/client-user/update-avatar
// @Summary      C端用户更新头像
// @Description  访问 /api/v1/c/client-user/update-avatar，C端用户更新头像
// @Tags         C端用户
// @Accept       json
// @Produce      json
// @Param        body  body  clientuser.UpdateAvatarParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/client-user/update-avatar [post]
func (h *handler) updateAvatar(c *gin.Context) {
	var param clientuser.UpdateAvatarParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.UpdateAvatar(c, &param)
	result.Success(c, nil)
}

// updatePasswordHandler handles POST /api/v1/c/client-user/update-password
// @Summary      C端用户修改密码
// @Description  访问 /api/v1/c/client-user/update-password，C端用户修改密码
// @Tags         C端用户
// @Accept       json
// @Produce      json
// @Param        body  body  clientuser.UpdatePasswordParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/client-user/update-password [post]
func (h *handler) updatePassword(c *gin.Context) {
	var param clientuser.UpdatePasswordParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.UpdatePassword(c, &param)
	result.Success(c, nil)
}
