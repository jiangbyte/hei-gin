// Package bind æä¾› Gin è¯·æ±‚ç»‘å®šï¼›JSON èµ°å…¨å±€ stringly ç¼–è§£ç ã€‚
//
// Author: Charlie
package bind

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"hei-gin/internal/framework/core/stringly"
)

// JSON è¯»å–è¯·æ±‚ä½“å¹¶ç”¨ stringly è§£ç ï¼Œå†æŒ‰ binding tag æ ¡éªŒã€‚
func JSON(c *gin.Context, obj any) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if len(body) > 0 {
		if err := stringly.Unmarshal(body, obj); err != nil {
			return err
		}
	}
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
