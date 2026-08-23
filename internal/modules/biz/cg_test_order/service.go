// internal/modules/biz/cg_test_order/service.go 业务服务。
//
// Author: Charlie

package cg_test_order

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/biz/scope"
)

// Service 订单服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 biz.cg_test_order 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "biz.cg_test_order",
		Order:  92,
		Models: []any{&Order{}, &OrderItem{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建订单。
func (s *Service) Create(ctx context.Context, accountID string, req AddParam, sess *security.SessionPayload) error {
	row := Order{
		ID: idgen.Next(), OrderNo: req.OrderNo, Name: req.Name, CustomerName: req.CustomerName,
		CustomerPhone: req.CustomerPhone, Status: req.Status, Type: req.Type, OrderedAt: req.OrderedAt,
		PaidAt: req.PaidAt, TotalAmount: req.TotalAmount, ItemCount: req.ItemCount, NeedInvoice: req.NeedInvoice,
		InvoiceConfig: mustJSON(req.InvoiceConfig), Remark: req.Remark, Extra: mustJSON(req.Extra),
		CreatedBy: &accountID, UpdatedBy: &accountID, OwnerDeptID: scope.DefaultOwnerDeptID(sess),
	}
	return s.repo.CreateOrder(ctx, &row)
}

// Update 更新订单。
func (s *Service) Update(ctx context.Context, accountID string, req EditParam) error {
	return s.repo.UpdateOrder(ctx, req.ID, map[string]any{
		"order_no": req.OrderNo, "name": req.Name, "customer_name": req.CustomerName, "customer_phone": req.CustomerPhone,
		"status": req.Status, "type": req.Type, "ordered_at": req.OrderedAt, "paid_at": req.PaidAt,
		"total_amount": req.TotalAmount, "item_count": req.ItemCount, "need_invoice": req.NeedInvoice,
		"invoice_config": mustJSON(req.InvoiceConfig), "remark": req.Remark, "extra": mustJSON(req.Extra), "updated_by": accountID,
	})
}

// Delete 批量删除订单。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteOrdersByIDs(ctx, ids)
}

// Detail 订单详情。
func (s *Service) Detail(ctx context.Context, id string) (*Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []Order, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageOrders(ctx, p, sess)
	return rows, total, current, size, err
}

// CreateItem 创建明细。
func (s *Service) CreateItem(ctx context.Context, accountID string, req ItemAddParam) error {
	row := OrderItem{
		ID: idgen.Next(), OrderID: req.OrderID, SKUCode: req.SKUCode, Name: req.Name, Category: req.Category,
		Status: req.Status, Quantity: req.Quantity, UnitPrice: req.UnitPrice, ShippedAt: req.ShippedAt,
		IsGift: req.IsGift, ItemConfig: mustJSON(req.ItemConfig), Remark: req.Remark, Extra: mustJSON(req.Extra),
		CreatedBy: &accountID, UpdatedBy: &accountID,
	}
	return s.repo.CreateItem(ctx, &row)
}

// UpdateItem 更新明细。
func (s *Service) UpdateItem(ctx context.Context, accountID string, req ItemEditParam) error {
	return s.repo.UpdateItem(ctx, req.ID, map[string]any{
		"order_id": req.OrderID, "sku_code": req.SKUCode, "name": req.Name, "category": req.Category,
		"status": req.Status, "quantity": req.Quantity, "unit_price": req.UnitPrice, "shipped_at": req.ShippedAt,
		"is_gift": req.IsGift, "item_config": mustJSON(req.ItemConfig), "remark": req.Remark,
		"extra": mustJSON(req.Extra), "updated_by": accountID,
	})
}

// DeleteItems 批量删除明细。
func (s *Service) DeleteItems(ctx context.Context, ids []string) error {
	return s.repo.DeleteItemsByIDs(ctx, ids)
}

// DetailItem 明细详情。
func (s *Service) DetailItem(ctx context.Context, id string) (*OrderItem, error) {
	return s.repo.GetItemByID(ctx, id)
}

// PageItems 明细分页。
func (s *Service) PageItems(ctx context.Context, p ItemPageParam) (rows []OrderItem, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageItems(ctx, p)
	return rows, total, current, size, err
}

func mustJSON(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte("{}"))
	}
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return b
}
