package cg_test_order

import (
	"time"

	"hei-gin/internal/framework/core/schema"
)

// AddParam åˆ›å»ºè®¢å•å…¥å‚ã€‚
//
// Author: Charlie
type AddParam struct {
	OrderNo       string         `json:"order_no" binding:"required"`
	Name          string         `json:"name" binding:"required"`
	CustomerName  string         `json:"customer_name" binding:"required"`
	CustomerPhone *string        `json:"customer_phone"`
	Status        string         `json:"status" binding:"required"`
	Type          string         `json:"type" binding:"required"`
	OrderedAt     time.Time      `json:"ordered_at" binding:"required"`
	PaidAt        *time.Time     `json:"paid_at"`
	TotalAmount   float64        `json:"total_amount"`
	ItemCount     int            `json:"item_count"`
	NeedInvoice   bool           `json:"need_invoice"`
	InvoiceConfig map[string]any `json:"invoice_config"`
	Remark        *string        `json:"remark"`
	Extra         map[string]any `json:"extra"`
}

// EditParam æ›´æ–°è®¢å•å…¥å‚ã€‚
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// PageParam è®¢å•åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	OrderNo      string `form:"order_no"`
	Name         string `form:"name"`
	CustomerName string `form:"customer_name"`
	Status       string `form:"status"`
	Type         string `form:"type"`
}

// ItemAddParam åˆ›å»ºè®¢å•æ˜Žç»†å…¥å‚ã€‚
//
// Author: Charlie
type ItemAddParam struct {
	OrderID    string         `json:"order_id" binding:"required"`
	SKUCode    string         `json:"sku_code" binding:"required"`
	Name       string         `json:"name" binding:"required"`
	Category   *string        `json:"category"`
	Status     string         `json:"status" binding:"required"`
	Quantity   int            `json:"quantity"`
	UnitPrice  float64        `json:"unit_price"`
	ShippedAt  *time.Time     `json:"shipped_at"`
	IsGift     bool           `json:"is_gift"`
	ItemConfig map[string]any `json:"item_config"`
	Remark     *string        `json:"remark"`
	Extra      map[string]any `json:"extra"`
}

// ItemEditParam æ›´æ–°è®¢å•æ˜Žç»†å…¥å‚ã€‚
//
// Author: Charlie
type ItemEditParam struct {
	ID string `json:"id" binding:"required"`
	ItemAddParam
}

// ItemPageParam è®¢å•æ˜Žç»†åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type ItemPageParam struct {
	schema.PageQuery
	OrderID string `form:"order_id"`
}
