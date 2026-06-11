package result

import (
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"

	"github.com/gin-gonic/gin"
)

func getTraceID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	v, exists := c.Get("trace_id")
	if !exists {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func httpStatus(code int) int {
	switch {
	case code >= 600 || code <= 0:
		return 200
	case code >= 500:
		return 500
	case code >= 429:
		return 429
	case code >= 400:
		return code
	default:
		return 200
	}
}

func writeJSON(c *gin.Context, httpCode int, body gin.H) {
	if c.Writer.Written() {
		return
	}
	c.JSON(httpCode, body)
}

// Success writes a successful response.
func Success(c *gin.Context, data any) {
	writeJSON(c, 200, gin.H{
		"code":     200,
		"message":  "请求成功",
		"data":     data,
		"success":  true,
		"trace_id": getTraceID(c),
	})
}

// Failure writes an error response.
func Failure(c *gin.Context, message string, code int) {
	writeJSON(c, httpStatus(code), gin.H{
		"code":     code,
		"message":  message,
		"data":     nil,
		"success":  false,
		"trace_id": getTraceID(c),
	})
}

// PageData is the structured pagination response body.
type PageData struct {
	Records any   `json:"records"`
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	Size    int   `json:"size"`
	Pages   int   `json:"pages"`
}

// PageDataResult writes a paginated response.
func PageDataResult(c *gin.Context, records any, total int64, page, size int) {
	if c.IsAborted() {
		return
	}
	pages := 0
	if size > 0 {
		pages = int((total + int64(size) - 1) / int64(size))
	}
	writeJSON(c, 200, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			string(enums.PageDataFieldRecords): records,
			string(enums.PageDataFieldTotal):   total,
			string(enums.PageDataFieldPage):    page,
			string(enums.PageDataFieldSize):    size,
			string(enums.PageDataFieldPages):   pages,
		},
		"success":  true,
		"trace_id": getTraceID(c),
	})
}

// WriteError writes a BusinessError or generic server error to the response.
func WriteError(c *gin.Context, err error) {
	if err == nil || c.IsAborted() {
		return
	}
	if be, ok := err.(*exception.BusinessError); ok {
		Failure(c, be.Message, be.Code)
	} else {
		Failure(c, "服务器内部错误", 500)
	}
	c.Abort()
}
