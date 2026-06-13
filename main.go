package main

// @title           Hei Gin API
// @version         1.0.0
// @description     Hei Gin 后台管理系统 API 文档
// @termsOfService  https://github.com/jiangbyte/hei-gin
// @contact.name    Charlie
// @contact.url     https://github.com/jiangbyte
// @license.name    MIT
// @license.url     https://github.com/jiangbyte/hei-gin/blob/main/LICENSE
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer token authentication

import (
	plugin_client "hei-gin/plugins/plugin-client"
	plugin_im "hei-gin/plugins/plugin-im"
	plugin_sys "hei-gin/plugins/plugin-sys"
	"hei-gin/sdk/kernel/app"
)

func main() {
	plugin_sys.RegisterPlugin()
	plugin_sys.RegisterRoutes()
	plugin_sys.RegisterMigrations()

	plugin_client.RegisterPlugin()
	plugin_client.RegisterRoutes()
	plugin_client.RegisterMigrations()

	plugin_im.RegisterPlugin()
	plugin_im.RegisterRoutes()
	plugin_im.RegisterMigrations()

	app.Run()
}
