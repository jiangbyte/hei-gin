// Package cg_test_order 为代码生成演示的订单业务模块。
package cg_test_order

import (
	"time"

	"gorm.io/datatypes"
)

// Order 演示订单实体，对应表 cg_test_order。
//
// Author: Charlie
type Order struct {
	ID            string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	OrderNo       string         `gorm:"column:order_no;size:64" json:"order_no"`
	Name          string         `gorm:"column:name;size:120" json:"name"`
	CustomerName  string         `gorm:"column:customer_name;size:120" json:"customer_name"`
	CustomerPhone *string        `gorm:"column:customer_phone;size:32" json:"customer_phone"`
	Status        string         `gorm:"column:status;size:32" json:"status"`
	Type          string         `gorm:"column:type;size:32" json:"type"`
	OrderedAt     time.Time      `gorm:"column:ordered_at" json:"ordered_at"`
	PaidAt        *time.Time     `gorm:"column:paid_at" json:"paid_at"`
	TotalAmount   float64        `gorm:"column:total_amount" json:"total_amount"`
	ItemCount     int            `gorm:"column:item_count" json:"item_count"`
	NeedInvoice   bool           `gorm:"column:need_invoice" json:"need_invoice"`
	InvoiceConfig datatypes.JSON `gorm:"column:invoice_config;type:jsonb" json:"invoice_config"`
	Remark        *string        `gorm:"column:remark;type:text" json:"remark"`
	Extra         datatypes.JSON `gorm:"column:extra;type:jsonb" json:"extra"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy     *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy     *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
	OwnerDeptID   *string        `gorm:"column:owner_dept_id;size:64" json:"owner_dept_id"`
}

// TableName 返回 Order 对应的数据库表名。
func (Order) TableName() string { return "cg_test_order" }

// OrderItem 演示订单明细实体，对应表 cg_test_order_item。
//
// Author: Charlie
type OrderItem struct {
	ID         string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	OrderID    string         `gorm:"column:order_id;size:64" json:"order_id"`
	SKUCode    string         `gorm:"column:sku_code;size:64" json:"sku_code"`
	Name       string         `gorm:"column:name;size:120" json:"name"`
	Category   *string        `gorm:"column:category;size:32" json:"category"`
	Status     string         `gorm:"column:status;size:32" json:"status"`
	Quantity   int            `gorm:"column:quantity" json:"quantity"`
	UnitPrice  float64        `gorm:"column:unit_price" json:"unit_price"`
	ShippedAt  *time.Time     `gorm:"column:shipped_at" json:"shipped_at"`
	IsGift     bool           `gorm:"column:is_gift" json:"is_gift"`
	ItemConfig datatypes.JSON `gorm:"column:item_config;type:jsonb" json:"item_config"`
	Remark     *string        `gorm:"column:remark;type:text" json:"remark"`
	Extra      datatypes.JSON `gorm:"column:extra;type:jsonb" json:"extra"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy  *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy  *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 OrderItem 对应的数据库表名。
func (OrderItem) TableName() string { return "cg_test_order_item" }
