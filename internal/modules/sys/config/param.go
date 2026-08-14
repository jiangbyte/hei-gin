// internal/modules/sys/config/param.go 入参定义。
//
// Author: Charlie

package config

import "hei-gin/internal/framework/core/schema"

// AddParam åˆ›å»ºé…ç½®å…¥å‚ã€‚
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

// EditParam æ›´æ–°é…ç½®å…¥å‚ã€‚
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam é…ç½®åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	ConfigKey string `form:"config_key"`
	Category  string `form:"category"`
}

// ListParam é…ç½®åˆ—è¡¨æŸ¥è¯¢ã€‚
//
// Author: Charlie
type ListParam struct {
	Category string `form:"category"`
	Scope    string `form:"scope"`
}

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// BatchItemParam æ‰¹é‡ä¿å­˜å•é¡¹å…¥å‚ï¼ˆremark ä¸Ž description å…¼å®¹ï¼‰ã€‚
//
// Author: Charlie
type BatchItemParam struct {
	ConfigKey   string  `json:"config_key" binding:"required"`
	ConfigValue *string `json:"config_value"`
	Description *string `json:"description"`
	Remark      *string `json:"remark"`
	Category    *string `json:"category"`
}

// BatchSaveParam æ‰¹é‡ä¿å­˜å…¥å‚ã€‚
//
// Author: Charlie
type BatchSaveParam struct {
	Items []BatchItemParam `json:"items" binding:"required"`
}

// TestWebhookParam å®¡è®¡å‘Šè­¦ Webhook æµ‹è¯•å…¥å‚ï¼ˆwebhook_url/webhook_secret ä¸Ž url/secret å…¼å®¹ï¼‰ã€‚
//
// Author: Charlie
type TestWebhookParam struct {
	URL           string `json:"url"`
	Secret        string `json:"secret"`
	WebhookURL    string `json:"webhook_url"`
	WebhookSecret string `json:"webhook_secret"`
}
