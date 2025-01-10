package port

import (
	"ordersvc/internal/core/domain"
	"ordersvc/pkg/common"
	"time"
)

type Repository interface {
	Order(ctx common.ServiceContextManager, request domain.OrderRequest) error
	CheckOrder(ctx common.ServiceContextManager, status string, startTime, endTime time.Time) ([]domain.Order, error)
}
