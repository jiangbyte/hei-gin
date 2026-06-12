package exception

import "fmt"

// BusinessError represents a business-logic error with a message and code.
type BusinessError struct {
	Message string
	Code    int
}

// NewBusinessError creates a BusinessError with the given message and code.
// If message is empty, it defaults to "业务异常". If code is 0, it defaults to 400.
func NewBusinessError(message string, code int) *BusinessError {
	if message == "" {
		message = "业务异常"
	}
	if code == 0 {
		code = 400
	}
	return &BusinessError{Message: message, Code: code}
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("BusinessError{code=%d, message=%s}", e.Code, e.Message)
}
