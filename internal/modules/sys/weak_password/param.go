// internal/modules/sys/weak_password/param.go 入参定义。
//
// Author: Charlie

package weakpassword

import "hei-gin/internal/framework/core/schema"

// AddParam 创建弱密码入参。
//
// Author: Charlie
type AddParam struct {
	Password string `json:"password" binding:"required"`
}

// EditParam 更新弱密码入参。
//
// Author: Charlie
type EditParam struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PageParam 弱密码分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Password string `form:"password"`
}

// ListParam 弱密码列表查询。
//
// Author: Charlie
type ListParam struct {
	Password string `form:"password"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
