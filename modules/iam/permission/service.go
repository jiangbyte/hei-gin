package permission

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// New 构建 iam.permission 模块。
func New(_ *shared.Deps) module.Module {
	return module.Module{
		Name:   "iam.permission",
		Routes: nil,
		Models: nil,
	}
}
