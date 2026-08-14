// internal/modules/sys/file/param.go 入参定义。
//
// Author: Charlie

package file

import "hei-gin/internal/framework/core/schema"

// EditParam 更新文件元数据入参。
//
// Author: Charlie
type EditParam struct {
	ID           string  `json:"id" binding:"required"`
	OriginalName *string `json:"original_name"`
}

// PageParam 文件分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	OriginalName string `form:"original_name"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
