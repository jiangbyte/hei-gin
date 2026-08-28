// internal/modules/biz/cg_test_order/repo.go 持久化仓储。
//
// Author: Charlie

package cg_test_order

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/core/security/datascope"
	"hei-gin/internal/framework/platform/db/dialect"
)

// Repo 订单持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// CreateOrder 创建订单。
func (r *Repo) CreateOrder(ctx context.Context, row *Order) error {
	return r.with(ctx).Create(row).Error
}

// UpdateOrder 更新订单。
func (r *Repo) UpdateOrder(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Order{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteOrdersByIDs 批量删除订单（含明细）。
func (r *Repo) DeleteOrdersByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id IN ?", ids).Delete(&OrderItem{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", ids).Delete(&Order{}).Error
	})
}

// GetOrderByID 按主键查订单。
func (r *Repo) GetOrderByID(ctx context.Context, id string) (*Order, error) {
	var row Order
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageOrders 分页查订单；sess 非空时按 owner_dept_id 数据范围过滤。
func (r *Repo) PageOrders(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []Order, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Order{})
	if sess != nil {
		db = datascope.ApplyKey(db, sess, "biz:cgtestorder:page", "owner_dept_id", "created_by")
	}
	if p.OrderNo != "" {
		db = db.Where(dialect.ILike(db, "order_no"), dialect.Contains(p.OrderNo))
	}
	if p.Name != "" {
		db = db.Where(dialect.ILike(db, "name"), dialect.Contains(p.Name))
	}
	if p.CustomerName != "" {
		db = db.Where(dialect.ILike(db, "customer_name"), dialect.Contains(p.CustomerName))
	}
	if p.Status != "" {
		db = db.Where("status = ?", p.Status)
	}
	if p.Type != "" {
		db = db.Where("type = ?", p.Type)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("ordered_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// CreateItem 创建明细。
func (r *Repo) CreateItem(ctx context.Context, row *OrderItem) error {
	return r.with(ctx).Create(row).Error
}

// UpdateItem 更新明细。
func (r *Repo) UpdateItem(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&OrderItem{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteItemsByIDs 批量删除明细。
func (r *Repo) DeleteItemsByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&OrderItem{}).Error
}

// GetItemByID 按主键查明细。
func (r *Repo) GetItemByID(ctx context.Context, id string) (*OrderItem, error) {
	var row OrderItem
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageItems 分页查明细。
func (r *Repo) PageItems(ctx context.Context, p ItemPageParam) (rows []OrderItem, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&OrderItem{})
	if p.OrderID != "" {
		db = db.Where("order_id = ?", p.OrderID)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("created_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
