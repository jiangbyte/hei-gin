package middleware

import (
	"log"
	"runtime/debug"

	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case *exception.BusinessError:
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

func SafeCall(fn func()) (err error) {
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
	fn()
	return nil
}

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
