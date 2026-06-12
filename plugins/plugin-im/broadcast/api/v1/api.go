package v1

import (
	"strconv"

	"hei-gin/plugins/plugin-im/broadcast"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
)

type handler struct {
	service *broadcast.Service
}

var defaultHandler = newHandler(broadcast.DefaultModule)

func newHandler(module *broadcast.Module) *handler {
	return &handler{service: module.Service()}
}

func RegisterRoutes(r *gin.Engine) {
	// POST /api/v1/sys/im/broadcast/send
	r.POST("/api/v1/sys/im/broadcast/send",
		authMW.HeiCheckLogin(),
		registry.Perm("sys:im:broadcast:send", "发送通知"),
		log.SysLog("发送通知"),
		authMW.NoRepeat(5000),
		defaultHandler.send,
	)

	// GET /api/v1/sys/im/broadcast/list
	r.GET("/api/v1/sys/im/broadcast/list",
		authMW.HeiCheckLogin(),
		registry.Perm("sys:im:broadcast:list", "通知列表"),
		defaultHandler.list,
	)

	// GET /api/v1/sys/im/broadcast/unread-list
	r.GET("/api/v1/sys/im/broadcast/unread-list",
		authMW.HeiCheckLogin(),
		defaultHandler.unreadList,
	)

	// POST /api/v1/sys/im/broadcast/read
	r.POST("/api/v1/sys/im/broadcast/read",
		authMW.HeiCheckLogin(),
		defaultHandler.read,
	)

	// GET /api/v1/sys/im/broadcast/detail
	r.GET("/api/v1/sys/im/broadcast/detail",
		authMW.HeiCheckLogin(),
		defaultHandler.detail,
	)
}

func RegisterClientRoutes(r *gin.Engine) {
	// GET /api/v1/c/im/broadcast/unread-list
	r.GET("/api/v1/c/im/broadcast/unread-list",
		authMW.HeiClientCheckLogin(),
		defaultHandler.clientUnreadList,
	)

	// POST /api/v1/c/im/broadcast/read
	r.POST("/api/v1/c/im/broadcast/read",
		authMW.HeiClientCheckLogin(),
		defaultHandler.clientRead,
	)

	// GET /api/v1/c/im/broadcast/detail
	r.GET("/api/v1/c/im/broadcast/detail",
		authMW.HeiClientCheckLogin(),
		defaultHandler.detail,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}

// sendHandler handles POST /api/v1/sys/im/broadcast/send
// @Summary      即时通讯广播发送消息
// @Description  访问 /api/v1/sys/im/broadcast/send，即时通讯广播发送消息
// @Tags         即时通讯广播
// @Accept       json
// @Produce      json
// @Param        body  body  broadcast.SendBroadcastParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/broadcast/send [post]
func (h *handler) send(c *gin.Context) {
	var param broadcast.SendBroadcastParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Send(c, &param)
	result.Success(c, nil)
}

// listHandler handles GET /api/v1/sys/im/broadcast/list
// @Summary      即时通讯广播列表查询
// @Description  访问 /api/v1/sys/im/broadcast/list，即时通讯广播列表查询
// @Tags         即时通讯广播
// @Accept       json
// @Produce      json
// @Param        cursor  query  string  false  "cursor"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/broadcast/list [get]
func (h *handler) list(c *gin.Context) {
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	list, hasMore := h.service.Page(c, cursor, size)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

// unreadListHandler handles GET /api/v1/sys/im/broadcast/unread-list
// @Summary      即时通讯广播未读列表
// @Description  访问 /api/v1/sys/im/broadcast/unread-list，即时通讯广播未读列表
// @Tags         即时通讯广播
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/broadcast/unread-list [get]
func (h *handler) unreadList(c *gin.Context) {
	list := h.service.UnreadList(c, string(enums.LoginTypeBusiness))
	result.Success(c, list)
}

// clientUnreadListHandler handles GET /api/v1/c/im/broadcast/unread-list
// @Summary      即时通讯广播未读列表
// @Description  访问 /api/v1/c/im/broadcast/unread-list，即时通讯广播未读列表
// @Tags         即时通讯广播
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/broadcast/unread-list [get]
func (h *handler) clientUnreadList(c *gin.Context) {
	list := h.service.UnreadList(c, string(enums.LoginTypeConsumer))
	result.Success(c, list)
}

// readHandler handles POST /api/v1/sys/im/broadcast/read
// @Summary      即时通讯广播标记已读
// @Description  访问 /api/v1/sys/im/broadcast/read，即时通讯广播标记已读
// @Tags         即时通讯广播
// @Accept       json
// @Produce      json
// @Param        body  body  broadcast.ReadParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/broadcast/read [post]
func (h *handler) read(c *gin.Context) {
	var param broadcast.ReadParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.MarkRead(c, &param)
	result.Success(c, nil)
}

// clientReadHandler handles POST /api/v1/c/im/broadcast/read
// @Summary      即时通讯广播标记已读
// @Description  访问 /api/v1/c/im/broadcast/read，即时通讯广播标记已读
// @Tags         即时通讯广播
// @Accept       json
// @Produce      json
// @Param        body  body  broadcast.ReadParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/broadcast/read [post]
func (h *handler) clientRead(c *gin.Context) {
	var param broadcast.ReadParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.MarkReadConsumer(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/im/broadcast/detail and GET /api/v1/c/im/broadcast/detail
// @Summary      即时通讯广播详情查询
// @Description  访问 /api/v1/sys/im/broadcast/detail，即时通讯广播详情查询
// @Tags         即时通讯广播
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/broadcast/detail [get]
// @Router       /api/v1/c/im/broadcast/detail [get]
func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
	result.Success(c, vo)
}
