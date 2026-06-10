package v1

import (
	"hei-gin/sdk/auth"
	"hei-gin/sdk/registry"
	middleware "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	clientuser "hei-gin/plugins/plugin-client/user"

	"github.com/gin-gonic/gin"
)

var clientAuth = auth.Consumer

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

func pageHandler(c *gin.Context) {
	var param clientuser.ClientUserPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	clientuser.ClientUserPage(c, &param)
}

func createHandler(c *gin.Context) {
	var vo clientuser.ClientUserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	clientuser.ClientUserCreate(c, &vo, auth.GetLoginIDDefaultNull(c))
	result.Success(c, nil)
}

func modifyHandler(c *gin.Context) {
	var vo clientuser.ClientUserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	clientuser.ClientUserModify(c, &vo, auth.GetLoginIDDefaultNull(c))
	result.Success(c, nil)
}

func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	clientuser.ClientUserRemove(c, param.IDs)
	result.Success(c, nil)
}

func detailHandler(c *gin.Context) {
	id := c.Query("id")
	vo := clientuser.ClientUserDetail(c, id)
	result.Success(c, vo)
}

func currentHandler(c *gin.Context) {
	vo := clientuser.Current(c, clientAuth.GetLoginIDDefaultNull(c))
	result.Success(c, vo)
}

func updateProfileHandler(c *gin.Context) {
	var param clientuser.UpdateProfileParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	clientuser.UpdateProfile(c, clientAuth.GetLoginIDDefaultNull(c), &param)
	result.Success(c, nil)
}

func updateAvatarHandler(c *gin.Context) {
	var param clientuser.UpdateAvatarParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	clientuser.UpdateAvatar(c, clientAuth.GetLoginIDDefaultNull(c), &param)
	result.Success(c, nil)
}

func updatePasswordHandler(c *gin.Context) {
	var param clientuser.UpdatePasswordParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	clientuser.UpdatePassword(c, clientAuth.GetLoginIDDefaultNull(c), &param)
	result.Success(c, nil)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
