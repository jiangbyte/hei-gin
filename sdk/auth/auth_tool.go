package auth

import (
	"log"

	"hei-gin/sdk/config"
	"hei-gin/sdk/kernel/plugin"
)

var Business = &Realm{ID: BusinessID, tool: newBaseAuthTool(BusinessID)}
var Consumer = &Realm{ID: ConsumerID, tool: newBaseAuthTool(ConsumerID)}

// ---- plugin registration ----

type authPlugin struct{ plugin.NoopPlugin }

func (m *authPlugin) Name() string { return "auth" }

func (m *authPlugin) Init() error {
	Business.tool.Init(config.C.Token.ExpireSeconds, config.C.Token.TokenName)
	Consumer.tool.Init(config.C.Token.ExpireSeconds, config.C.Token.TokenName)
	log.Println("[auth] plugin initialized")
	return nil
}

func (m *authPlugin) Start() error {
	if err := RunPermissionScan(); err != nil {
		return err
	}
	return nil
}

func init() {
	plugin.Register(&authPlugin{})
}
