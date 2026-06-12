package plugin_client

import (
	"hei-gin/sdk/kernel/plugin"

	// Blank-import to trigger model registration and route registration
	_ "hei-gin/plugins/plugin-client/session"
	_ "hei-gin/plugins/plugin-client/session/api/v1"
	_ "hei-gin/plugins/plugin-client/user"
	_ "hei-gin/plugins/plugin-client/user/api/v1"

	_ "hei-gin/plugins/plugin-client/auth/captcha/api/v1"
	_ "hei-gin/plugins/plugin-client/auth/sm2/api/v1"
	_ "hei-gin/plugins/plugin-client/auth/username/api/v1"
)

type ClientPlugin struct {
	plugin.NoopPlugin
}

func (p *ClientPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name:        "plugin-client",
		Version:     "1.0.0",
		Description: "Client-side management plugin (user, session, message, auth)",
	}
}

func (p *ClientPlugin) Name() string { return "plugin-client" }

func init() {
	plugin.Register(&ClientPlugin{})
}
