package middleware

import (
	"log"
	"runtime/debug"

	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"

	"github.com/gin-gonic/gin"
)

// Recovery returns a Gin middleware that catches panics and converts them
// into structured JSON responses.
//
// Panic handling:
//   - *exception.BusinessError → 200 with business code and message (no stack trace logged)
//   - Any other panic          → 500 with "服务器内部错误" (full stack trace logged)
//
// NOTE: Must be outermost middleware. Use with gin.New(), NOT gin.Default().
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case *exception.BusinessError:
					// Business errors: no stack trace, just return JSON
					result.Failure(c, e.Message, e.Code)
				case error:
					log.Printf("[PANIC] %v\n%s", e, string(debug.Stack()))
					result.Failure(c, "服务器内部错误", 500)
				default:
					log.Printf("[PANIC] %v\n%s", rec, string(debug.Stack()))
					result.Failure(c, "服务器内部错误", 500)
				}
				c.Abort()
			}
		}()

		c.Next()

		if err := c.Errors.Last(); err != nil {
			result.Failure(c, err.Error(), 400)
			c.Abort()
		}
	}
}

// SafeCall executes fn with panic recovery, returning any error.
// BusinessError panics are converted to errors; other panics are re-panicked
// to be caught by the top-level Recovery middleware.
//
// Usage:
//
//	err := SafeCall(func() {
//	    riskyOperation()
//	})
func SafeCall(fn func()) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			switch e := rec.(type) {
			case *exception.BusinessError:
				err = e
			default:
				// Re-panic non-business errors to be caught by top-level Recovery
				panic(rec)
			}
		}
	}()
	fn()
	return nil
}

// SafeCallCtx is like SafeCall but passes a context-aware function.
func SafeCallCtx(fn func(ctx *gin.Context) error, c *gin.Context) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			switch e := rec.(type) {
			case *exception.BusinessError:
				err = e
			default:
				panic(rec)
			}
		}
	}()
	return fn(c)
}

// GoSafe executes fn with panic recovery in a goroutine.
// Unlike SafeCall, it re-panics non-business errors so the top-level Recovery
// middleware (which only covers gin handler goroutines) is NOT available here.
// Instead, GoSafe logs all panics and only suppresses BusinessError.
// Use GoSafe for background goroutines (event bus subscribers, timers, etc.)
func GoSafe(fn func()) {
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case *exception.BusinessError:
					log.Printf("[GoSafe] BusinessError suppressed: %v", e)
				case error:
					log.Printf("[GoSafe] Panic recovered: %v\n%s", e, string(debug.Stack()))
				default:
					log.Printf("[GoSafe] Panic recovered: %v\n%s", rec, string(debug.Stack()))
				}
			}
		}()
		fn()
	}()
}
