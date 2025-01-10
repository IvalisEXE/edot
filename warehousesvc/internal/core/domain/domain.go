package domain

import (
	"time"
)

type Warehouse struct {
	ID         int        `json:"id"`
	ShopID     int        `json:"shop_id"`
	Name       string     `json:"name"`
	Address    string     `json:"address"`
	IsActive   bool       `json:"is_active"`
	StockTotal int        `json:"stock_total"`
	StockedAt  *time.Time `json:"stocked_at"`
	Created    UserMeta   `json:"created"`
	Updated    UserMeta   `json:"updated"`
}

type WarehouseStockItem struct {
	ID          int        `json:"id"`
	ShopID      int        `json:"shop_id"`
	WarehouseID int        `json:"warehouse_id"`
	Warehouse   Warehouse  `json:"warehouse"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Price       float32    `json:"price"`
	StockTotal  int        `json:"stock_total"`
	StockedAt   *time.Time `json:"stocked_at"`
	Created     UserMeta   `json:"created"`
	Updated     UserMeta   `json:"updated"`
}

type WarehouseStock struct {
	ID                 int                `json:"id"`
	ShopID             int                `json:"shop_id"`
	Warehouse          Warehouse          `json:"wareshouse"`
	WarehouseStockItem WarehouseStockItem `json:"wareshouse_stock_item"`
	StockTotal         int                `json:"stock_total"`
	StockedAt          *time.Time         `json:"stocked_at"`
	Created            UserMeta           `json:"created"`
	Updated            UserMeta           `json:"updated"`
}

type WarehouseStockTransfer struct {
	ID                     int                `json:"id"`
	ShopID                 int                `json:"shop_id"`
	WarehouseOriginID      int                `json:"wareshouse_origin_id"`
	WarehouseOrigin        Warehouse          `json:"wareshouse_origin"`
	WarehouseDestinationID int                `json:"warehouse_destination_id"`
	WarehouseDestination   Warehouse          `json:"warehouse_destination"`
	WarehouseStockItemID   int                `json:"warehouse_stock_item_id"`
	WarehouseStockItem     WarehouseStockItem `json:"warehouse_stock_item"`
	TransferTotal          int                `json:"transfer_total"`
	TransferAt             *time.Time         `json:"transfer_at"`
	Created                UserMeta           `json:"created"`
}
