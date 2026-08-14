// Package schema 提供通用查询 DTO（分页、单 ID）。
//
// Author: Charlie
package schema

// PageQuery 通用分页查询。
//
// Author: Charlie
type PageQuery struct {
	Current int `form:"current" json:"current"`
	Size    int `form:"size" json:"size"`
}

// Normalize 将 current/size 规范到合法范围（默认 1/20，size 上限 100）。
func (q PageQuery) Normalize() (current, size int) {
	current = q.Current
	size = q.Size
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

// IDQuery 单 ID 请求。
//
// Author: Charlie
type IDQuery struct {
	ID string `form:"id" json:"id" binding:"required"`
}
