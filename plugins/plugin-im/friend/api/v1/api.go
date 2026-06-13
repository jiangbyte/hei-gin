package v1

import (
	"hei-gin/sdk/auth"
	"strconv"

	"hei-gin/plugins/plugin-im/friend"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
)

type handler struct {
	service *friend.Service
}

var defaultHandler = newHandler(friend.DefaultModule)

func newHandler(module *friend.Module) *handler {
	return &handler{service: module.Service()}
}

func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/sys/im/friend/send-request",
		authMW.CheckLogin(auth.Business),
		authMW.NoRepeat(3000),
		defaultHandler.sendRequest,
	)
	r.POST("/api/v1/sys/im/friend/accept",
		authMW.CheckLogin(auth.Business),
		defaultHandler.accept,
	)
	r.POST("/api/v1/sys/im/friend/reject",
		authMW.CheckLogin(auth.Business),
		defaultHandler.reject,
	)
	r.GET("/api/v1/sys/im/friend/list",
		authMW.CheckLogin(auth.Business),
		defaultHandler.list,
	)
	r.GET("/api/v1/sys/im/friend/pending-requests",
		authMW.CheckLogin(auth.Business),
		defaultHandler.pendingRequests,
	)
	r.POST("/api/v1/sys/im/friend/remove",
		authMW.CheckLogin(auth.Business),
		defaultHandler.remove,
	)
	r.POST("/api/v1/sys/im/friend/block",
		authMW.CheckLogin(auth.Business),
		defaultHandler.block,
	)
	r.POST("/api/v1/sys/im/friend/unblock",
		authMW.CheckLogin(auth.Business),
		defaultHandler.unblock,
	)
	r.GET("/api/v1/sys/im/friend/block-list",
		authMW.CheckLogin(auth.Business),
		defaultHandler.blockList,
	)
	r.POST("/api/v1/sys/im/friend/remark",
		authMW.CheckLogin(auth.Business),
		defaultHandler.remark,
	)
	r.GET("/api/v1/sys/im/friend/search",
		authMW.CheckLogin(auth.Business),
		defaultHandler.search,
	)
}

func RegisterClientRoutes(r *gin.Engine) {
	r.POST("/api/v1/c/im/friend/send-request",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientSendRequest,
	)
	r.POST("/api/v1/c/im/friend/accept",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientAccept,
	)
	r.POST("/api/v1/c/im/friend/reject",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientReject,
	)
	r.GET("/api/v1/c/im/friend/list",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientList,
	)
	r.GET("/api/v1/c/im/friend/pending-requests",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientPendingRequests,
	)
	r.POST("/api/v1/c/im/friend/remove",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientRemove,
	)
	r.POST("/api/v1/c/im/friend/block",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientBlock,
	)
	r.POST("/api/v1/c/im/friend/unblock",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientUnblock,
	)
	r.GET("/api/v1/c/im/friend/block-list",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientBlockList,
	)
	r.POST("/api/v1/c/im/friend/remark",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientRemark,
	)
	r.GET("/api/v1/c/im/friend/search",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.search,
	)
}

func Register() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}

// @Summary      即时通讯好友发送申请
// @Description  访问 /api/v1/sys/im/friend/send-request，即时通讯好友发送申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.SendRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/send-request [post]
func (h *handler) sendRequest(c *gin.Context) {
	h.sendRequestByType(c, string(auth.BusinessID))
}

// @Summary      即时通讯好友发送申请
// @Description  访问 /api/v1/c/im/friend/send-request，即时通讯好友发送申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.SendRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/send-request [post]
func (h *handler) clientSendRequest(c *gin.Context) {
	h.sendRequestByType(c, string(auth.ConsumerID))
}

func (h *handler) sendRequestByType(c *gin.Context, userType string) {
	var param friend.SendRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.SendRequest(c, userType, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯好友接受申请
// @Description  访问 /api/v1/sys/im/friend/accept，即时通讯好友接受申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.HandleRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/accept [post]
func (h *handler) accept(c *gin.Context) {
	h.acceptByType(c, string(auth.BusinessID))
}

// @Summary      即时通讯好友接受申请
// @Description  访问 /api/v1/c/im/friend/accept，即时通讯好友接受申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.HandleRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/accept [post]
func (h *handler) clientAccept(c *gin.Context) {
	h.acceptByType(c, string(auth.ConsumerID))
}

func (h *handler) acceptByType(c *gin.Context, userType string) {
	var param friend.HandleRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.AcceptRequest(c, userType, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯好友拒绝申请
// @Description  访问 /api/v1/sys/im/friend/reject，即时通讯好友拒绝申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.HandleRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/reject [post]
func (h *handler) reject(c *gin.Context) {
	h.rejectByType(c, string(auth.BusinessID))
}

// @Summary      即时通讯好友拒绝申请
// @Description  访问 /api/v1/c/im/friend/reject，即时通讯好友拒绝申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.HandleRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/reject [post]
func (h *handler) clientReject(c *gin.Context) {
	h.rejectByType(c, string(auth.ConsumerID))
}

func (h *handler) rejectByType(c *gin.Context, userType string) {
	var param friend.HandleRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.RejectRequest(c, userType, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯好友列表查询
// @Description  访问 /api/v1/sys/im/friend/list，即时通讯好友列表查询
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/list [get]
func (h *handler) list(c *gin.Context) {
	result.Success(c, h.service.List(c, string(auth.BusinessID)))
}

// @Summary      即时通讯好友列表查询
// @Description  访问 /api/v1/c/im/friend/list，即时通讯好友列表查询
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/list [get]
func (h *handler) clientList(c *gin.Context) {
	result.Success(c, h.service.List(c, string(auth.ConsumerID)))
}

// @Summary      即时通讯好友待处理申请列表
// @Description  访问 /api/v1/sys/im/friend/pending-requests，即时通讯好友待处理申请列表
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/pending-requests [get]
func (h *handler) pendingRequests(c *gin.Context) {
	h.pendingRequestsByType(c, string(auth.BusinessID))
}

// @Summary      即时通讯好友待处理申请列表
// @Description  访问 /api/v1/c/im/friend/pending-requests，即时通讯好友待处理申请列表
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/pending-requests [get]
func (h *handler) clientPendingRequests(c *gin.Context) {
	h.pendingRequestsByType(c, string(auth.ConsumerID))
}

func (h *handler) pendingRequestsByType(c *gin.Context, userType string) {
	incoming, outgoing := h.service.PendingRequests(c, userType)
	result.Success(c, gin.H{"incoming": incoming, "outgoing": outgoing})
}

// @Summary      即时通讯好友删除
// @Description  访问 /api/v1/sys/im/friend/remove，即时通讯好友删除
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.RemoveFriendParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/remove [post]
func (h *handler) remove(c *gin.Context) {
	h.removeByType(c, string(auth.BusinessID))
}

// @Summary      即时通讯好友删除
// @Description  访问 /api/v1/c/im/friend/remove，即时通讯好友删除
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.RemoveFriendParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/remove [post]
func (h *handler) clientRemove(c *gin.Context) {
	h.removeByType(c, string(auth.ConsumerID))
}

func (h *handler) removeByType(c *gin.Context, userType string) {
	var param friend.RemoveFriendParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Remove(c, userType, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯好友拉黑
// @Description  访问 /api/v1/sys/im/friend/block，即时通讯好友拉黑
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.BlockParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/block [post]
func (h *handler) block(c *gin.Context) {
	h.blockByType(c, string(auth.BusinessID))
}

// @Summary      即时通讯好友拉黑
// @Description  访问 /api/v1/c/im/friend/block，即时通讯好友拉黑
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.BlockParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/block [post]
func (h *handler) clientBlock(c *gin.Context) {
	h.blockByType(c, string(auth.ConsumerID))
}

func (h *handler) blockByType(c *gin.Context, userType string) {
	var param friend.BlockParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Block(c, userType, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯好友取消拉黑
// @Description  访问 /api/v1/sys/im/friend/unblock，即时通讯好友取消拉黑
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.BlockParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/unblock [post]
func (h *handler) unblock(c *gin.Context) {
	h.unblockByType(c, string(auth.BusinessID))
}

// @Summary      即时通讯好友取消拉黑
// @Description  访问 /api/v1/c/im/friend/unblock，即时通讯好友取消拉黑
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.BlockParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/unblock [post]
func (h *handler) clientUnblock(c *gin.Context) {
	h.unblockByType(c, string(auth.ConsumerID))
}

func (h *handler) unblockByType(c *gin.Context, userType string) {
	var param friend.BlockParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Unblock(c, userType, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯好友拉黑列表
// @Description  访问 /api/v1/sys/im/friend/block-list，即时通讯好友拉黑列表
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/block-list [get]
func (h *handler) blockList(c *gin.Context) {
	result.Success(c, h.service.BlockList(c, string(auth.BusinessID)))
}

// @Summary      即时通讯好友拉黑列表
// @Description  访问 /api/v1/c/im/friend/block-list，即时通讯好友拉黑列表
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/block-list [get]
func (h *handler) clientBlockList(c *gin.Context) {
	result.Success(c, h.service.BlockList(c, string(auth.ConsumerID)))
}

// @Summary      即时通讯好友设置备注
// @Description  访问 /api/v1/sys/im/friend/remark，即时通讯好友设置备注
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.RemarkParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/remark [post]
func (h *handler) remark(c *gin.Context) {
	h.remarkByType(c, string(auth.BusinessID))
}

// @Summary      即时通讯好友设置备注
// @Description  访问 /api/v1/c/im/friend/remark，即时通讯好友设置备注
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.RemarkParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/remark [post]
func (h *handler) clientRemark(c *gin.Context) {
	h.remarkByType(c, string(auth.ConsumerID))
}

func (h *handler) remarkByType(c *gin.Context, userType string) {
	var param friend.RemarkParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.UpdateRemark(c, userType, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯好友搜索
// @Description  访问 /api/v1/sys/im/friend/search，即时通讯好友搜索
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        keyword  query  string  false  "keyword"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/search [get]
// @Router       /api/v1/c/im/friend/search [get]
func (h *handler) search(c *gin.Context) {
	keyword := c.Query("keyword")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	result.Success(c, h.service.Search(c, keyword, size))
}
