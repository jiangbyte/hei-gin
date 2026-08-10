package app

import (
	"hei-gin/framework/platform/module"
)

// AttachRegisteredModules 根据已 Register 的构造器装配模块（需先 blank import app/internal/modules/all）。
func AttachRegisteredModules(d *Deps) {
	md := &module.Deps{
		Cfg:      d.Cfg,
		DB:       d.DB,
		Redis:    d.Redis,
		Sessions: d.Sessions,
		Perms:    d.Perms,
		Storage:  d.Storage,
	}
	reg := module.BuildAll(md, d.Cfg.Modules.Disabled, d.Cfg.Modules.Enabled)
	for _, m := range reg.Modules {
		for _, eh := range m.EventHandlers {
			d.Events.Subscribe(eh.Event, eh.Handler)
		}
	}
	d.Modules = reg
}
