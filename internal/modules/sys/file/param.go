// internal/modules/sys/file/param.go 入参定义。
//
// Author: Charlie

package file

import "hei-gin/internal/framework/core/schema"

// EditParam 更新文件元数据入参（对齐 hei-boot SysFileEditParam：original_name @NotBlank）。
//
// Author: Charlie
type EditParam struct {
	ID           string `json:"id" binding:"required"`
	OriginalName string `json:"original_name" binding:"required"`
}

// PageParam 文件分页查询（对齐 hei-boot SysFilePageParam）。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	OriginalName    string `form:"original_name"`
	ObjectName      string `form:"object_name"`
	StorageProvider string `form:"storage_provider"`
	ContentType     string `form:"content_type"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
