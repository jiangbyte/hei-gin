package v1

import (
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	clientuser "hei-gin/plugins/plugin-client/user"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/client-user/page",
		registry.Perm("client:user:page", "C端用户分页"),
		pageHandler,
	)

	r.POST("/api/v1/client-user/create",
		registry.Perm("client:user:create", "添加C端用户"),
		createHandler,
	)

	r.POST("/api/v1/client-user/modify",
		registry.Perm("client:user:modify", "编辑C端用户"),
		modifyHandler,
	)

	r.POST("/api/v1/client-user/remove",
		registry.Perm("client:user:remove", "删除C端用户"),
		removeHandler,
	)

	r.GET("/api/v1/client-user/detail",
		registry.Perm("client:user:detail", "C端用户详情"),
		detailHandler,
	)

	r.GET("/api/v1/c/client-user/current",
		middleware.HeiClientCheckLogin(),
		currentHandler,
	)

	r.POST("/api/v1/c/client-user/update-profile",
		middleware.HeiClientCheckLogin(),
		log.SysLog("C端用户更新个人信息"),
		middleware.NoRepeat(3000),
		updateProfileHandler,
	)

	r.POST("/api/v1/c/client-user/update-avatar",
		middleware.HeiClientCheckLogin(),
		log.SysLog("C端用户更新头像"),
		updateAvatarHandler,
	)

	r.POST("/api/v1/c/client-user/update-password",
		middleware.HeiClientCheckLogin(),
		log.SysLog("C端用户修改密码"),
		middleware.NoRepeat(3000),
		updatePasswordHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/client-user/page
func pageHandler(c *gin.Context) {
	var param clientuser.ClientUserPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientuser.ClientUserPage(c, &param)
}

// createHandler handles POST /api/v1/client-user/create
func createHandler(c *gin.Context) {
	var vo clientuser.ClientUserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientuser.ClientUserCreate(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/client-user/modify
func modifyHandler(c *gin.Context) {
	var vo clientuser.ClientUserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientuser.ClientUserModify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/client-user/remove
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientuser.ClientUserRemove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/client-user/detail
func detailHandler(c *gin.Context) {
	vo := clientuser.ClientUserDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// currentHandler handles GET /api/v1/c/client-user/current
func currentHandler(c *gin.Context) {
	vo := clientuser.ClientUserCurrent(c)
	result.Success(c, vo)
}

// updateProfileHandler handles POST /api/v1/c/client-user/update-profile
func updateProfileHandler(c *gin.Context) {
	var param clientuser.UpdateProfileParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientuser.ClientUserUpdateProfile(c, &param)
	result.Success(c, nil)
}

// updateAvatarHandler handles POST /api/v1/c/client-user/update-avatar
func updateAvatarHandler(c *gin.Context) {
	var param clientuser.UpdateAvatarParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientuser.ClientUserUpdateAvatar(c, &param)
	result.Success(c, nil)
}

// updatePasswordHandler handles POST /api/v1/c/client-user/update-password
func updatePasswordHandler(c *gin.Context) {
	var param clientuser.UpdatePasswordParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientuser.ClientUserUpdatePassword(c, &param)
	result.Success(c, nil)
}
