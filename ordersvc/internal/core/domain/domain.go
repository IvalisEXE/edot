package domain

import "time"

type Order struct {
	ID                     int       `json:"id"`
	ShopID                 int       `json:"shop_id"`
	WarehouseID            int       `json:"warehouse_id"`
	WarehouseStockItemID   int       `json:"warehouse_stock_item_id"`
	WarehouseStockItemCode string    `json:"warehouse_stock_item_code"`
	ProductName            string    `json:"product_name"`
	Price                  float32   `json:"price"`
	Quantity               int       `json:"quantity"`
	PaymentMethod          string    `json:"payment_method"`
	PaymentExpired         time.Time `json:"payment_expired"`
	TotalPayment           float32   `json:"total_payment"`
	Status                 string    `json:"status"`
	OrderDate              time.Time `json:"order_date"`
}

type OrderRequest struct {
	ShopID                 int     `json:"shop_id" validate:"required"`
	WarehouseID            int     `json:"warehouse_id" validate:"required"`
	WarehouseStockItemID   int     `json:"warehouse_stock_item_id" validate:"required"`
	WarehouseStockItemCode string  `json:"warehouse_stock_item_code" validate:"required"`
	ProductName            string  `json:"product_name" validate:"required"`
	Price                  float32 `json:"price" validate:"required"`
	Quantity               int     `json:"quantity" validate:"required"`
	PaymentMethod          string  `json:"payment_method" validate:"required"`
	TotalPayment           float32 `json:"total_payment" validate:"required"`
}
