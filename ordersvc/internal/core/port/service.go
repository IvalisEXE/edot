package port

import (
	"ordersvc/internal/core/domain"
	"ordersvc/pkg/common"
)

type Service interface {
	Order(ctx common.ServiceContextManager, request domain.OrderRequest) error
	LockStock(ctx common.ServiceContextManager) error
	ReleaseStock(ctx common.ServiceContextManager) error
}
