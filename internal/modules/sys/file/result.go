// internal/modules/sys/file/result.go 出参定义。
//
// Author: Charlie

package file

// ObjectNameParam 按对象名获取 URL 入参。
//
// Author: Charlie
type ObjectNameParam struct {
	ObjectName string `json:"object_name" binding:"required"`
}

// URLResult 文件访问 URL 结果。
//
// Author: Charlie
type URLResult struct {
	ObjectName string `json:"object_name"`
	URL        string `json:"url"`
	ExpiresIn  *int64 `json:"expires_in,omitempty"`
}
