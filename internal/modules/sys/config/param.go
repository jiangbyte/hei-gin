// internal/modules/sys/config/param.go 入参定义。
//
// Author: Charlie

package config

import "hei-gin/internal/framework/core/schema"

// AddParam 创建配置入参。
//
// Author: Charlie
type AddParam struct {
	ConfigKey   string  `json:"config_key" binding:"required"`
	ConfigValue *string `json:"config_value"`
	Category    *string `json:"category"`
	Remark      *string `json:"remark"`
	SortCode    int     `json:"sort_code"`
	ValueType   string  `json:"value_type"`
	Label       *string `json:"label"`
	Scope       *string `json:"scope"`
	Scene       *string `json:"scene"`
}

// EditParam 更新配置入参。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam 配置分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	ConfigKey string `form:"config_key"`
	Category  string `form:"category"`
}

// ListParam 配置列表查询。
//
// Author: Charlie
type ListParam struct {
	Category string `form:"category"`
	Scope    string `form:"scope"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// BatchItemParam 批量保存单项入参（remark 与 description 兼容）。
//
// Author: Charlie
// BatchItemParam 批量保存单项入参（对齐 hei-boot SysConfigBatchItemParam；remark 与 description 兼容）。
//
// Author: Charlie
type BatchItemParam struct {
	ConfigKey   string  `json:"config_key" binding:"required"`
	ConfigValue *string `json:"config_value"`
	Description *string `json:"description"`
	Remark      *string `json:"remark"`
	Category    *string `json:"category"`
	ValueType   *string `json:"value_type"`
	Label       *string `json:"label"`
	Scope       *string `json:"scope"`
	Scene       *string `json:"scene"`
	IsBuiltin   *bool   `json:"is_builtin"`
	SortCode    *int    `json:"sort_code"`
}

// BatchSaveParam 批量保存入参。
//
// Author: Charlie
type BatchSaveParam struct {
	Items []BatchItemParam `json:"items" binding:"required"`
}

// TestWebhookParam 审计告警 Webhook 测试入参（webhook_url/webhook_secret 与 url/secret 兼容）。
//
// Author: Charlie
type TestWebhookParam struct {
	URL           string `json:"url"`
	Secret        string `json:"secret"`
	WebhookURL    string `json:"webhook_url"`
	WebhookSecret string `json:"webhook_secret"`
}
