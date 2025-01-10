package port

import (
	"usersvc/internal/core/domain"
	"usersvc/pkg/common"
)

type Service interface {
	Login(ctx common.ServiceContextManager, request domain.UserLoginRequest) (*domain.UserTokenData, error)
	SaveSession(ctx common.ServiceContextManager, user *domain.User, token *domain.UserTokenData) error
	FindByID(ctx common.ServiceContextManager, id int) (*domain.User, error)
}
