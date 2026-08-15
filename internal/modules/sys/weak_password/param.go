// internal/modules/sys/weak_password/param.go 入参定义。
//
// Author: Charlie

package weakpassword

import (
	"hei-gin/internal/framework/core/schema"
)

// AddParam 创建弱密码入参。
//
// Author: Charlie
type AddParam struct {
	Password string `json:"password" binding:"required,max=255"`
}

// EditParam 更新弱密码入参。
//
// Author: Charlie
type EditParam struct {
	ID       string `json:"id" binding:"required,max=64"`
	Password string `json:"password" binding:"required,max=255"`
}

// PageParam 弱密码分页查询（兼容 web 发送的 page 与通用 current）。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Page     int    `form:"page" json:"page"`
	Password string `form:"password" json:"password"`
	Keyword  string `form:"keyword" json:"keyword"`
}

// Normalize 分页规范化（page 别名优先）。
func (q PageParam) Normalize() (current, size int) {
	current, size = q.Current, q.Size
	if current < 1 {
		current = q.Page
	}
	if current < 1 {
		current = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return current, size
}

// ListParam 弱密码列表查询。
//
// Author: Charlie
type ListParam struct {
	Password string `form:"password"`
	Keyword  string `form:"keyword"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
