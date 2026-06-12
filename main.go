package main

// @title           Hei Gin API
// @version         1.0.0
// @description     Hei Gin 后台管理系统 API 文档
// @termsOfService  https://github.com/jiangbyte/hei-gin
// @contact.name    Charlie
// @contact.url     https://github.com/jiangbyte
// @license.name    MIT
// @license.url     https://github.com/jiangbyte/hei-gin/blob/main/LICENSE
// @host            localhost:18885
// @BasePath        /api/v1
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer token authentication

import (
	"hei-gin/sdk/app"

	// Plugin route/permission self-registration
	_ "hei-gin/plugins/plugin-client"
	_ "hei-gin/plugins/plugin-im"
	_ "hei-gin/plugins/plugin-sys"
)

func main() {
	app.Run()
}
