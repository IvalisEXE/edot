package port

import (
	"warehousesvc/internal/core/domain"
	"warehousesvc/pkg/common"
)

type Service interface {
	FindAllStock(ctx common.ServiceContextManager, filter domain.FindAllStockRequest) ([]domain.FindAllStockResponse, error)

	TransferStock(ctx common.ServiceContextManager, request domain.TransferStockRequest) error

	UnactiveWarehouse(ctx common.ServiceContextManager, id int) error
}
