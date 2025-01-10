package port

import (
	"usersvc/internal/core/domain"
	"usersvc/pkg/common"
)

type Repository interface {
	FindByID(ctx common.ServiceContextManager, id int) (*domain.User, error)
	FindByPhone(ctx common.ServiceContextManager, phone string) (*domain.User, error)
	UpdateLastLogin(ctx common.ServiceContextManager, userID int) error
}
