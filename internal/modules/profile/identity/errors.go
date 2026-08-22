// internal/modules/profile/identity/errors.go 业务错误。
//
// Author: Charlie
package identity

// BizError 带 HTTP 状态的业务错误。
type BizError struct {
	HTTPStatus int
	Code       int
	Message    string
}

func (e *BizError) Error() string { return e.Message }

func bizErr(httpStatus, code int, msg string) error {
	return &BizError{HTTPStatus: httpStatus, Code: code, Message: msg}
}
