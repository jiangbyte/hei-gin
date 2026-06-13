package plugin_client

import (
	"sync"

	authcaptchav1 "hei-gin/plugins/plugin-client/auth/captcha/api/v1"
	authsm2v1 "hei-gin/plugins/plugin-client/auth/sm2/api/v1"
	authusernamev1 "hei-gin/plugins/plugin-client/auth/username/api/v1"
	sessionv1 "hei-gin/plugins/plugin-client/session/api/v1"
	userv1 "hei-gin/plugins/plugin-client/user/api/v1"
	"hei-gin/sdk/kernel/plugin"
)

type ClientPlugin struct {
	plugin.NoopPlugin
}

var registerOnce sync.Once

func (p *ClientPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name:        "plugin-client",
		Version:     "1.0.0",
		Description: "Client-side management plugin (user, session, message, auth)",
	}
}

func (p *ClientPlugin) Name() string { return "plugin-client" }

func RegisterPlugin() {
	registerOnce.Do(func() {
		plugin.Register(&ClientPlugin{})
	})
}

func RegisterRoutes() {
	sessionv1.Register()
	userv1.Register()
	authcaptchav1.Register()
	authsm2v1.Register()
	authusernamev1.Register()
}
