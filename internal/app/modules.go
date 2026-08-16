// internal/app/modules.go 模块装配。
//
// Author: Charlie

package app

import (
	"hei-gin/internal/framework/platform/module"
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
		Notify:   d.Notify,
		Audit:    d.Audit,
		Runtime:  d.Runtime,
		Jobs:     d.Jobs,
	}
	reg := module.BuildAll(md, d.Cfg.Modules.Disabled, d.Cfg.Modules.Enabled)
	d.Modules = reg
}
