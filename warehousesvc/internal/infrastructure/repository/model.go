package repository

import (
	"time"
	"warehousesvc/internal/core/domain"
)

const (
	ACTIVITY_STOCK_WAREHOUSE          = "warehouse"
	ACTIVITY_STOCK_WAREHOUSE_TRANSFER = "warehouse-transfer"
	ACTIVITY_STOCK_ORDER              = "order"
)

type Warehouse struct {
	ID         int        `gorm:"column:id;primaryKey;autoIncrement"`
	ShopID     int        `gorm:"column:shop_id;index;not null" json:"shop_id"`
	Name       string     `gorm:"column:name;size:255"`
	Address    string     `gorm:"column:address"`
	IsActive   bool       `gorm:"column:is_active" json:"is_active"`
	StockTotal int        `gorm:"column:stock_total" json:"stock_total"`
	StockedAt  *time.Time `gorm:"column:stocked_at" json:"stocked_at"`
	CreateMeta
	UpdateMeta
	DeleteMeta
}

func (Warehouse) TableName() string {
	return "warehouse"
}

func (w Warehouse) ToDomain() *domain.Warehouse {
	return &domain.Warehouse{
		ID:         w.ID,
		ShopID:     w.ShopID,
		Name:       w.Name,
		Address:    w.Address,
		IsActive:   w.IsActive,
		StockTotal: w.StockTotal,
		StockedAt:  w.StockedAt,
		Created: domain.UserMeta{
			At: &w.CreatedAt,
			By: &w.CreatedBy,
		},
		Updated: domain.UserMeta{
			At: w.UpdatedAt,
			By: w.UpdatedBy,
		},
	}
}

type WarehouseStockItem struct {
	ID          int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID      int        `gorm:"column:shop_id;index;not null" json:"shop_id"`
	WarehouseID int        `gorm:"column:warehouse_id;index;not null" json:"warehouse_id"`
	Warehouse   Warehouse  `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
	Code        string     `gorm:"column:code" json:"code"`
	Name        string     `gorm:"column:name" json:"name"`
	Price       float32    `gorm:"column:price;type:numeric(12,2);not null;default:0" json:"price"`
	StockTotal  int        `gorm:"column:stock_total" json:"stock_total"`
	StockedAt   *time.Time `gorm:"column:stocked_at" json:"stocked_at"`
	CreateMeta
	UpdateMeta
}

func (WarehouseStockItem) TableName() string {
	return "warehouse_stock_item"
}

func (w *WarehouseStockItem) ToDomain() *domain.WarehouseStockItem {
	return &domain.WarehouseStockItem{
		ID:     w.ID,
		ShopID: w.ShopID,
		Warehouse: domain.Warehouse{
			ID:       w.Warehouse.ID,
			ShopID:   w.Warehouse.ShopID,
			Name:     w.Warehouse.Name,
			Address:  w.Warehouse.Address,
			IsActive: w.Warehouse.IsActive,
		},
		Name:       w.Name,
		Price:      w.Price,
		StockTotal: w.StockTotal,
		StockedAt:  w.StockedAt,
		Created: domain.UserMeta{
			At: &w.CreatedAt,
			By: &w.CreatedBy,
		},
		Updated: domain.UserMeta{
			At: w.UpdatedAt,
			By: w.UpdatedBy,
		},
	}
}

func (w WarehouseStockItem) ToStockList() domain.FindAllStockResponse {
	return domain.FindAllStockResponse{
		ID:         w.ID,
		Name:       w.Name,
		Price:      w.Price,
		StockTotal: w.StockTotal,
	}
}

type WarehouseStockTransfer struct {
	ID                     int                `gorm:"column:id;primaryKey;autoIncrement"`
	ShopID                 int                `gorm:"column:shop_id;index;not null"`
	WarehouseOriginID      int                `gorm:"column:warehouse_origin_id;index;not null"`
	WarehouseOrigin        Warehouse          `gorm:"foreignKey:WarehouseOriginID;references:ID"`
	WarehouseDestinationID int                `gorm:"column:warehouse_destination_id;index;not null"`
	WarehouseDestination   Warehouse          `gorm:"foreignKey:WarehouseDestinationID;references:ID"`
	WarehouseStockItemID   int                `gorm:"column:warehouse_destination_id;index;not null"`
	WarehouseStockItem     WarehouseStockItem `gorm:"foreignKey:WarehouseStockItemID;references:ID"`
	TransferTotal          int                `gorm:"column:transfer_total"`
	TransferAt             *time.Time         `gorm:"column:transfer_at"`
	CreateMeta
}

func (WarehouseStockTransfer) TableName() string {
	return "warehouse_stock_transfer"
}

type WarehouseAndWarehouseStockItem struct {
	ID                  int        `gorm:"column:id;primaryKey;autoIncrement"`
	ShopID              int        `gorm:"column:shop_id;index;not null"`
	WarehouseID         int        `gorm:"column:warehouse_id;index;not null"`
	Code                string     `gorm:"column:code"`
	Name                string     `gorm:"column:name"`
	Price               float32    `gorm:"column:price;type:numeric(12,2);not null;default:0"`
	StockTotal          int        `gorm:"column:stock_total"`
	WarehouseStockTotal int        `gorm:"column:warehouse_stock_total"`
	StockedAt           *time.Time `gorm:"column:stocked_at" json:"stocked_at"`
	CreateMeta
	UpdateMeta
}

func (w WarehouseAndWarehouseStockItem) ToDomain() *domain.WarehouseStockItem {
	return &domain.WarehouseStockItem{
		ID:          w.ID,
		ShopID:      w.ShopID,
		WarehouseID: w.WarehouseID,
		Name:        w.Name,
		Price:       w.Price,
		StockTotal:  w.StockTotal,
		Warehouse:   domain.Warehouse{StockTotal: w.WarehouseStockTotal},
		StockedAt:   w.StockedAt,
		Created: domain.UserMeta{
			At: &w.CreatedAt,
			By: &w.CreatedBy,
		},
		Updated: domain.UserMeta{
			At: w.UpdatedAt,
			By: w.UpdatedBy,
		},
	}
}
