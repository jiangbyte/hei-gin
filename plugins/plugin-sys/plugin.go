package plugin_sys

import (
	stdlog "log"
	"sync"
	"time"

	analyzerv1 "hei-gin/plugins/plugin-sys/analyze/api/v1"
	authcaptchav1 "hei-gin/plugins/plugin-sys/auth/captcha/api/v1"
	authsm2v1 "hei-gin/plugins/plugin-sys/auth/sm2/api/v1"
	authusernamev1 "hei-gin/plugins/plugin-sys/auth/username/api/v1"
	bannerv1 "hei-gin/plugins/plugin-sys/banner/api/v1"
	configv1 "hei-gin/plugins/plugin-sys/config/api/v1"
	dictv1 "hei-gin/plugins/plugin-sys/dict/api/v1"
	filev1 "hei-gin/plugins/plugin-sys/file/api/v1"
	groupv1 "hei-gin/plugins/plugin-sys/group/api/v1"
	homev1 "hei-gin/plugins/plugin-sys/home/api/v1"
	logv1 "hei-gin/plugins/plugin-sys/log/api/v1"
	noticev1 "hei-gin/plugins/plugin-sys/notice/api/v1"
	orgv1 "hei-gin/plugins/plugin-sys/org/api/v1"
	permissionv1 "hei-gin/plugins/plugin-sys/permission/api/v1"
	positionv1 "hei-gin/plugins/plugin-sys/position/api/v1"
	"hei-gin/plugins/plugin-sys/provider"
	resourcev1 "hei-gin/plugins/plugin-sys/resource/api/v1"
	rolev1 "hei-gin/plugins/plugin-sys/role/api/v1"
	sessionv1 "hei-gin/plugins/plugin-sys/session/api/v1"
	userv1 "hei-gin/plugins/plugin-sys/user/api/v1"
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

var registerOnce sync.Once

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

func RegisterPlugin() {
	registerOnce.Do(func() {
		plugin.Register(&SysPlugin{})
	})
}

func RegisterRoutes() {
	userv1.Register()
	rolev1.Register()
	orgv1.Register()
	groupv1.Register()
	positionv1.Register()
	dictv1.Register()
	configv1.Register()
	bannerv1.Register()
	homev1.Register()
	logv1.Register()
	noticev1.Register()
	filev1.Register()
	resourcev1.Register()
	sessionv1.Register()
	permissionv1.Register()
	analyzerv1.Register()
	authusernamev1.Register()
	authcaptchav1.Register()
	authsm2v1.Register()
}
