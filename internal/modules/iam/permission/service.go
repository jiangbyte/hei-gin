package permission

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// New æž„å»º iam.permission æ¨¡å—ã€‚
func New(_ *shared.Deps) module.Module {
	return module.Module{
		Name:   "iam.permission",
		Routes: nil,
		Models: nil,
	}
}
