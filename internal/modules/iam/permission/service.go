// internal/modules/iam/permission/service.go 业务服务。
//
// Author: Charlie

package permission

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// New 构建 iam.permission 模块。
func New(_ *shared.Deps) module.Module {
	return module.Module{
		Name:   "iam.permission",
		Routes: nil,
		Models: nil,
	}
}
