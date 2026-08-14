// internal/modules/message/notice/register.go 模块自注册。
//
// Author: Charlie

package notice

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册 message.notice 模块。
func init() {
	module.Register("message.notice", 60, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
