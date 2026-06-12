package safe

import (
	"log"
	"runtime/debug"

	"hei-gin/sdk/web/exception"
)

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
