// internal/framework/core/response/response.go 响应信封。
//
// Author: Charlie

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/stringly"
)

// ApiResponse æ ‡å‡†ä¿¡å°ï¼šcode/message/dataï¼ˆçº¿ä¸Š code ç» stringly ä¸ºå­—ç¬¦ä¸²ï¼‰ã€‚
//
// Author: Charlie
type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// PageData ä¸ºç»Ÿä¸€åˆ†é¡µç»“æž„ã€‚
//
// Author: Charlie
type PageData struct {
	Size    int64 `json:"size"`
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
	Pages   int64 `json:"pages"`
	Records any   `json:"records"`
}

// writeJSON ä»¥ stringly JSON å†™å…¥å“åº”ã€‚
func writeJSON(c *gin.Context, httpStatus int, v any) {
	b, err := stringly.Marshal(v)
	if err != nil {
		c.Data(http.StatusInternalServerError, "application/json; charset=utf-8", []byte(`{"code":"500","message":"marshal error","data":null}`))
		return
	}
	c.Data(httpStatus, "application/json; charset=utf-8", b)
}

// OK å†™å…¥æˆåŠŸä¿¡å°ï¼ˆHTTP 200ï¼‰ã€‚
func OK(c *gin.Context, data any) {
	writeJSON(c, http.StatusOK, ApiResponse{Code: 200, Message: "success", Data: data})
}

// Fail å†™å…¥å¤±è´¥ä¿¡å°ï¼ˆæŒ‡å®š HTTP çŠ¶æ€ä¸Žä¸šåŠ¡ç ï¼‰ã€‚
func Fail(c *gin.Context, httpStatus int, code int, message string) {
	writeJSON(c, httpStatus, ApiResponse{Code: code, Message: message, Data: nil})
}

// Page å†™å…¥åˆ†é¡µæˆåŠŸä¿¡å°ã€‚
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
