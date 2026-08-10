// Package bind 提供 Gin 请求绑定；JSON 走全局 stringly 编解码。
//
// Author: Charlie
package bind

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"hei-gin/framework/core/stringly"
)

// JSON 读取请求体并用 stringly 解码，再按 binding tag 校验。
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
