package repository

import (
	"ordersvc/internal/core/domain"
	"ordersvc/pkg/common"
	"time"
)

type Order struct {
	ID                     int       `gorm:"column:id;primaryKey;autoIncrement"`
	ShopID                 int       `gorm:"column:shop_id;index;not null"`
	WarehouseID            int       `gorm:"column:warehouse_id;index;not null"`
	WarehouseStockItemID   int       `gorm:"column:warehouse_stock_item_id;index;not null"`
	WarehouseStockItemCode string    `gorm:"column:warehouse_stock_item_code"`
	ProductName            string    `gorm:"column:product_name"`
	Price                  float32   `gorm:"column:price;type:numeric(12,2);not null;default:0"`
	Quantity               int       `gorm:"column:quantity"`
	PaymentMethod          string    `gorm:"column:payment_method"`
	PaymentExpired         time.Time `gorm:"column:payment_expired"`
	TotalPayment           float32   `gorm:"column:total_payment"`
	Status                 string    `gorm:"column:status"`
	OrderDate              time.Time `gorm:"column:order_date"`
	CreateMeta
}

func (Order) TableName() string {
	return "order"
}

func (o *Order) FromDomain(ctx common.ServiceContextManager, order domain.OrderRequest) {
	now := time.Now()
	o.ShopID = order.ShopID
	o.WarehouseID = order.WarehouseID
	o.WarehouseStockItemID = order.WarehouseStockItemID
	o.WarehouseStockItemCode = order.WarehouseStockItemCode
	o.ProductName = order.ProductName
	o.Price = order.Price
	o.Quantity = order.Quantity
	o.PaymentMethod = order.PaymentMethod
	o.PaymentExpired = now.Add(time.Hour)
	o.TotalPayment = order.TotalPayment
	o.Status = "OPEN"
	o.OrderDate = now
	o.CreateMeta = CreateMeta{
		CreatedAt: now,
		CreatedBy: ctx.GetUserHeader().ID,
	}
}

func (o *Order) ToDomain() domain.Order {
	return domain.Order{
		ID:                     o.ID,
		ShopID:                 o.ShopID,
		WarehouseID:            o.WarehouseID,
		WarehouseStockItemID:   o.WarehouseStockItemID,
		WarehouseStockItemCode: o.WarehouseStockItemCode,
		ProductName:            o.ProductName,
		Price:                  o.Price,
		Quantity:               o.Quantity,
		PaymentMethod:          o.PaymentMethod,
		PaymentExpired:         o.PaymentExpired,
		TotalPayment:           o.TotalPayment,
		Status:                 o.Status,
		OrderDate:              o.OrderDate,
	}
}
