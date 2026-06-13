package plugin_sys

import (
	stdlog "log"
	"time"

	"hei-gin/plugins/plugin-sys/provider"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/kernel/plugin"
	"hei-gin/sdk/log"
	"hei-gin/sdk/shared/contracts"
	"hei-gin/sdk/utils"
)

type SysPlugin struct {
	plugin.NoopPlugin
	permProvider *provider.PermissionProvider
	userProvider *provider.UserProvider
}

func (p *SysPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name:        "plugin-sys",
		Version:     "1.0.0",
		Description: "System management plugin (user, role, org, permission, etc.)",
	}
}

func (p *SysPlugin) Name() string { return "plugin-sys" }

func (p *SysPlugin) Init() error {
	p.permProvider = &provider.PermissionProvider{}
	p.userProvider = &provider.UserProvider{}

	auth.RegisterPermissionProvider(p.permProvider)

	var persister contracts.LogPersistenceAPI = &logPersister{}

	log.RegisterPersistence(func(ctx interface{}, category, name, exeStatus, exeMessage, opIP, opAddress, opBrowser, opOS, opUser, traceID, signData, method, url, params string, opTime interface{}) {
		opTimeStr := ""
		if t, ok := opTime.(time.Time); ok {
			opTimeStr = t.Format("2006-01-02 15:04:05")
		} else if s, ok := opTime.(string); ok {
			opTimeStr = s
		}

		entry := contracts.LogEntry{
			ID:         utils.GenerateID(),
			Category:   category,
			Name:       name,
			ExeStatus:  exeStatus,
			ExeMessage: exeMessage,
			OpIP:       opIP,
			OpAddress:  opAddress,
			OpBrowser:  opBrowser,
			OpOS:       opOS,
			OpUser:     opUser,
			TraceID:    traceID,
			SignData:   signData,
			ReqMethod:  method,
			ReqURL:     url,
			ParamJSON:  params,
			OpTime:     opTimeStr,
		}
		if err := persister.SaveLog(entry); err != nil {
			stdlog.Printf("[SYSLOG] Failed to persist log: %v", err)
		}
	})
	stdlog.Println("[plugin-sys] initialized")
	return nil
}

func init() {
	plugin.Register(&SysPlugin{})
}
