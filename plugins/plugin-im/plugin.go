package plugin_im

import (
	"sync"

	broadcastv1 "hei-gin/plugins/plugin-im/broadcast/api/v1"
	friendv1 "hei-gin/plugins/plugin-im/friend/api/v1"
	groupv1 "hei-gin/plugins/plugin-im/group/api/v1"
	message "hei-gin/plugins/plugin-im/message"
	messagev1 "hei-gin/plugins/plugin-im/message/api/v1"
	ws "hei-gin/plugins/plugin-im/ws"
	file "hei-gin/plugins/plugin-sys/file"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/kernel/plugin"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type IMPlugin struct{}

var registerOnce sync.Once

func (p *IMPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name:        "plugin-im",
		Version:     "1.0.0",
		Description: "Instant messaging plugin (WebSocket + group chat + friend + broadcast)",
	}
}

func (p *IMPlugin) Name() string { return "plugin-im" }
func (p *IMPlugin) Init() error {
	ws.InitCrossHub(ws.GlobalHub, db.Redis)
	return nil
}
func (p *IMPlugin) Start() error {
	ws.GlobalHub.StartOnlineBroadcast()
	return nil
}
func (p *IMPlugin) Stop() error {
	ws.GlobalHub.StopOnlineBroadcast()
	if runtime := ws.Runtime(); runtime != nil {
		runtime.Close()
	}
	return nil
}

func RegisterPlugin() {
	registerOnce.Do(func() {
		plugin.Register(&IMPlugin{})
	})
}

func RegisterRoutes() {
	registry.RegisterRoute(func(r *gin.Engine) {
		r.GET("/uploads/:bucket/:file_key", uploadHandler)
		r.GET("/api/v1/sys/im/ws", sysWSHandler)
		r.GET("/api/v1/c/im/ws", clientWSHandler)
	})
	broadcastv1.Register()
	friendv1.Register()
	groupv1.Register()
	messagev1.Register()
}

// @Summary      即时通讯连接接口调用
// @Description  访问 /uploads/:bucket/:file_key，即时通讯连接接口调用
// @Tags         即时通讯连接
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /uploads/:bucket/:file_key [get]
func uploadHandler(c *gin.Context) {
	bucket := c.Param("bucket")
	fileKey := c.Param("file_key")
	if err := file.DefaultModule.Service().DownloadByKey(c, bucket, fileKey); err == nil {
		return
	} else if err.Error() == "未授权/未登录" {
		result.Failure(c, err.Error(), 401)
		return
	}
	if err := message.DefaultModule.Service().ServeUploadedFile(c, bucket, fileKey); err != nil {
		result.Failure(c, err.Error(), 404)
	}
}

// @Summary      即时通讯连接WebSocket连接
// @Description  访问 /api/v1/sys/im/ws，即时通讯连接WebSocket连接
// @Tags         即时通讯连接
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/ws [get]
func sysWSHandler(c *gin.Context) {
	result := ws.AuthenticateFromToken(c, auth.BusinessID)
	if !result.OK {
		return
	}
	ws.GlobalHub.HandleWebSocket(c.Writer, c.Request, result.UserID, result.UserType)
}

// @Summary      即时通讯连接WebSocket连接
// @Description  访问 /api/v1/c/im/ws，即时通讯连接WebSocket连接
// @Tags         即时通讯连接
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/ws [get]
func clientWSHandler(c *gin.Context) {
	result := ws.AuthenticateFromToken(c, auth.ConsumerID)
	if !result.OK {
		return
	}
	ws.GlobalHub.HandleWebSocket(c.Writer, c.Request, result.UserID, result.UserType)
}
