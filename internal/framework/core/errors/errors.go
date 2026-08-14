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

// New 构造指定业务码与文案的 AppError。
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap 构造包装了底层错误的 AppError。
func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// BadRequest 返回 400 错误。
func BadRequest(msg string) *AppError { return New(400, msg) }

// Unauthorized 返回 401 错误。
func Unauthorized(msg string) *AppError { return New(401, msg) }

// Forbidden 返回 403 错误。
func Forbidden(msg string) *AppError { return New(403, msg) }

// NotFound 返回 404 错误。
func NotFound(msg string) *AppError { return New(404, msg) }

// Conflict 返回 409 错误。
func Conflict(msg string) *AppError { return New(409, msg) }

// Internal 返回 500 错误。
func Internal(msg string) *AppError { return New(500, msg) }
