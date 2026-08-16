// internal/modules/sys/file/register.go 模块自注册。
//
// Author: Charlie

package file

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册 sys.file 模块。
func init() {
	module.Register("sys.file", 50, func(d *module.Deps) module.Module {
		s := NewService(d.DB, d.Storage, d.Runtime)
		d.Provide(ServiceKey, s)
		return module.Module{
			Name:   "sys.file",
			Models: []any{&File{}},
			Routes: []module.RouteRegistrar{s.registerRoutes(d)},
		}
	})
}
