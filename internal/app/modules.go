// internal/app/modules.go 模块装配。
//
// Author: Charlie

package app

import (
	"hei-gin/internal/framework/platform/module"
)

// AttachRegisteredModules æ ¹æ®å·² Register çš„æž„é€ å™¨è£…é…æ¨¡å—ï¼ˆéœ€å…ˆ blank import app/internal/modules/allï¼‰ã€‚
func AttachRegisteredModules(d *Deps) {
	md := &module.Deps{
		Cfg:      d.Cfg,
		DB:       d.DB,
		Redis:    d.Redis,
		Sessions: d.Sessions,
		Perms:    d.Perms,
		Storage:  d.Storage,
		Notify:   d.Notify,
		Audit:    d.Audit,
	}
	reg := module.BuildAll(md, d.Cfg.Modules.Disabled, d.Cfg.Modules.Enabled)
	for _, m := range reg.Modules {
		for _, eh := range m.EventHandlers {
			d.Events.Subscribe(eh.Event, eh.Handler)
		}
	}
	d.Modules = reg
}
