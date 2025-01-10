package repository

import (
	"time"

	"go.elastic.co/apm"
	"gorm.io/gorm"

	"ordersvc/internal/core/domain"
	"ordersvc/internal/core/port"
	"ordersvc/pkg/common"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) port.Repository {
	return &repository{db}
}

func (r *repository) Order(ctx common.ServiceContextManager, request domain.OrderRequest) error {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.order.Order", "repository")
	defer span.End()

	ctx.SetContext(context)

	tx := ctx.GetDB()
	if tx == nil {
		tx = r.db
	}

	var order Order

	order.FromDomain(ctx, request)
	return tx.WithContext(ctx.GetContext()).Create(&order).Error
}

func (r *repository) CheckOrder(ctx common.ServiceContextManager, status string, startTime, endTime time.Time) ([]domain.Order, error) {
	var orders []Order

	tx := ctx.GetDB()
	if tx == nil {
		tx = r.db
	}

	if err := tx.Where("status = ? AND payment_expired < ? AND payment_expired > ?",
		status,
		startTime,
		endTime,
	).Find(&orders).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var result []domain.Order
	for _, order := range orders {
		result = append(result, order.ToDomain())
	}

	return result, nil
}
