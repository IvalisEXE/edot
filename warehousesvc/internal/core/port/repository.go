package port

import (
	"warehousesvc/internal/core/domain"
	"warehousesvc/pkg/common"
)

type Repository interface {
	CreateStockItem(ctx common.ServiceContextManager, request domain.WarehouseStockItem) (*int, error)

	CreateStockItemTransfer(ctx common.ServiceContextManager, request domain.WarehouseStockTransfer) (*int, error)

	UpdateStock(ctx common.ServiceContextManager, id, totalStock int) error

	UpdateStockItem(ctx common.ServiceContextManager, itemID, totalStockItem int) error

	FindAllStock(ctx common.ServiceContextManager, filter domain.FindAllStockRequest) ([]domain.FindAllStockResponse, error)

	FindStockByID(ctx common.ServiceContextManager, id int) (*domain.Warehouse, error)

	FindStockItemByID(ctx common.ServiceContextManager, itemID, warehouseID int) (*domain.WarehouseStockItem, error)

	UnactiveWarehouse(ctx common.ServiceContextManager, id int) error
}
