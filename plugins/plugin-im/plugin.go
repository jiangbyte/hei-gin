package plugin_im

import (
	message "hei-gin/plugins/plugin-im/message"
	ws "hei-gin/plugins/plugin-im/ws"
	file "hei-gin/plugins/plugin-sys/file"
	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/plugin"
	"hei-gin/sdk/registry"
	"hei-gin/sdk/result"

	"github.com/gin-gonic/gin"

	_ "hei-gin/plugins/plugin-im/broadcast"
	_ "hei-gin/plugins/plugin-im/broadcast/api/v1"
	_ "hei-gin/plugins/plugin-im/friend"
	_ "hei-gin/plugins/plugin-im/friend/api/v1"
	_ "hei-gin/plugins/plugin-im/group"
	_ "hei-gin/plugins/plugin-im/group/api/v1"
	_ "hei-gin/plugins/plugin-im/message/api/v1"
)

type IMPlugin struct{}

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
func (p *IMPlugin) Start() error { return nil }
func (p *IMPlugin) Stop() error {
	if ws.GlobalCrossHub != nil {
		ws.GlobalCrossHub.Close()
	}
	return nil
}

func init() {
	plugin.Register(&IMPlugin{})

	registry.RegisterRoute(func(r *gin.Engine) {
		r.GET("/uploads/:bucket/:file_key", uploadHandler)
		r.GET("/api/v1/sys/im/ws", sysWSHandler)
		r.GET("/api/v1/c/im/ws", clientWSHandler)
	})
}

func uploadHandler(c *gin.Context) {
	bucket := c.Param("bucket")
	fileKey := c.Param("file_key")
	if err := file.FileDownloadByKey(c, bucket, fileKey); err == nil {
		return
	} else if err.Error() == "未授权/未登录" {
		result.Failure(c, err.Error(), 401)
		return
	}
	if err := message.ServeUploadedFile(c, bucket, fileKey); err != nil {
		result.Failure(c, err.Error(), 404)
	}
}

func sysWSHandler(c *gin.Context) {
	result := ws.AuthenticateFromToken(c, enums.LoginTypeBusiness)
	if !result.OK {
		return
	}
	ws.GlobalHub.HandleWebSocket(c.Writer, c.Request, result.UserID, result.UserType)
}

func clientWSHandler(c *gin.Context) {
	result := ws.AuthenticateFromToken(c, enums.LoginTypeConsumer)
	if !result.OK {
		return
	}
	ws.GlobalHub.HandleWebSocket(c.Writer, c.Request, result.UserID, result.UserType)
}
