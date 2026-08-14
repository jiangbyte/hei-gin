// internal/framework/core/response/response.go 响应信封。
//
// Author: Charlie

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/stringly"
)

// ApiResponse 标准信封：code/message/data（线上 code 经 stringly 为字符串）。
//
// Author: Charlie
type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// PageData 为统一分页结构。
//
// Author: Charlie
type PageData struct {
	Size    int64 `json:"size"`
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
	Pages   int64 `json:"pages"`
	Records any   `json:"records"`
}

// writeJSON 以 stringly JSON 写入响应。
func writeJSON(c *gin.Context, httpStatus int, v any) {
	b, err := stringly.Marshal(v)
	if err != nil {
		c.Data(http.StatusInternalServerError, "application/json; charset=utf-8", []byte(`{"code":"500","message":"marshal error","data":null}`))
		return
	}
	c.Data(httpStatus, "application/json; charset=utf-8", b)
}

// OK 写入成功信封（HTTP 200）。
func OK(c *gin.Context, data any) {
	writeJSON(c, http.StatusOK, ApiResponse{Code: 200, Message: "success", Data: data})
}

// Fail 写入失败信封（指定 HTTP 状态与业务码）。
func Fail(c *gin.Context, httpStatus int, code int, message string) {
	writeJSON(c, httpStatus, ApiResponse{Code: code, Message: message, Data: nil})
}

// Page 写入分页成功信封。
func Page(c *gin.Context, current, size, total int64, records any) {
	pages := int64(0)
	if size > 0 {
		pages = (total + size - 1) / size
	}
	OK(c, PageData{
		Size:    size,
		Current: current,
		Total:   total,
		Pages:   pages,
		Records: records,
	})
}
