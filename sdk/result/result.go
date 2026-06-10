package result

import (
	"hei-gin/sdk/enums"

	"github.com/gin-gonic/gin"
)

// Responder wraps *gin.Context for unified JSON response formatting.
type Responder struct {
	c *gin.Context
}

// Wrap creates a Responder from a gin context.
func Wrap(c *gin.Context) *Responder {
	return &Responder{c: c}
}

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

// httpStatus maps a business code to an HTTP status code.
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

// OK writes a successful response with HTTP 200.
func (r *Responder) OK(data any) {
	r.c.JSON(200, gin.H{
		"code":     200,
		"message":  "请求成功",
		"data":     data,
		"success":  true,
		"trace_id": getTraceID(r.c),
	})
}

// Fail writes an error response with proper HTTP status code.
func (r *Responder) Fail(msg string, code int) {
	r.c.JSON(httpStatus(code), gin.H{
		"code":     code,
		"message":  msg,
		"data":     nil,
		"success":  false,
		"trace_id": getTraceID(r.c),
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

// Page writes a paginated success response.
func (r *Responder) Page(records any, total int64, page, size int) {
	pages := 0
	if size > 0 {
		pages = int((total + int64(size) - 1) / int64(size))
	}
	r.c.JSON(200, gin.H{
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
		"trace_id": getTraceID(r.c),
	})
}

// --- Convenience functions (delegate to Responder) ---

// Success writes a successful response.
func Success(c *gin.Context, data any) {
	Wrap(c).OK(data)
}

// Failure writes an error response.
func Failure(c *gin.Context, message string, code int) {
	Wrap(c).Fail(message, code)
}

// PageDataResult writes a paginated response.
func PageDataResult(c *gin.Context, records any, total int64, page, size int) {
	Wrap(c).Page(records, total, page, size)
}
