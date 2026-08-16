// Package errors 定义 handler / 中间件使用的类型化 API 错误。
//
// Author: Charlie
package errors

import "fmt"

// AppError 是带业务码与文案的 API 错误，可包装底层 error。
//
// Author: Charlie
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error 实现 error 接口。
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 返回被包装的底层错误。
func (e *AppError) Unwrap() error { return e.Err }
