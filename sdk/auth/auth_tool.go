package auth

import (
	"log"
	"sync"

	"hei-gin/sdk/config"
	"hei-gin/sdk/kernel/plugin"
)

var Business = &Realm{ID: BusinessID, tool: newBaseAuthTool(BusinessID)}
var Consumer = &Realm{ID: ConsumerID, tool: newBaseAuthTool(ConsumerID)}

// ---- plugin registration ----

type authPlugin struct{ plugin.NoopPlugin }

var registerOnce sync.Once

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

func RegisterPlugin() {
	registerOnce.Do(func() {
		plugin.Register(&authPlugin{})
	})
}
