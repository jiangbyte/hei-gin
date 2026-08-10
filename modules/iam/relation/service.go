package relation

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// New 构建 iam.relation 模块。
func New(_ *shared.Deps) module.Module {
	return module.Module{
		Name:   "iam.relation",
		Models: []any{&Relation{}},
	}
}
