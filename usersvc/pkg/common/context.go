package common

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"usersvc/internal/core/domain"
)

type ServiceContext struct {
	context     context.Context
	pagination  *Pagination
	sortParam   *SortParam
	audiTrail   *AudiTrail
	userClaim   *domain.UserCustomClaims
	userSession *domain.UserSession
	logger      *zap.Logger
	database    *gorm.DB
}

type ServiceContextManager interface {
	SetContext(c context.Context)
	GetContext() context.Context
	SetPagination(pagination *Pagination)
	GetPagination() *Pagination
	SetSortParam(param *SortParam)
	GetSortParam() *SortParam
	SetAudiTrail(audiTrail AudiTrail)
	GetAudiTrail() *AudiTrail
	SetUserClaim(userClaim *domain.UserCustomClaims)
	GetUserClaim() *domain.UserCustomClaims
	SetUserSession(userSession *domain.UserSession)
	GetUserSession() *domain.UserSession
	SetLogger(logger *zap.Logger)
	Logger() *zap.Logger
	BeginTransaction(db *gorm.DB) *gorm.DB
	GetDB() *gorm.DB
}

func NewServiceContext() ServiceContextManager {
	serviceContext := new(ServiceContext)
	serviceContext.SetContext(context.Background())

	return serviceContext
}

func (s *ServiceContext) SetContext(c context.Context) {
	s.context = c
}

func (s *ServiceContext) SetPagination(pagination *Pagination) {
	s.pagination = pagination
}

func (s *ServiceContext) SetMetaPagination(pagination *Pagination) {
	s.pagination = pagination
}

func (s *ServiceContext) SetSortParam(param *SortParam) {
	s.sortParam = param
}

func (s *ServiceContext) SetAudiTrail(audiTrail AudiTrail) {
	s.audiTrail = &audiTrail
}

func (s *ServiceContext) SetUserClaim(userClaim *domain.UserCustomClaims) {
	s.userClaim = userClaim
}

func (s *ServiceContext) SetUserSession(userSession *domain.UserSession) {
	s.userSession = userSession
}

func (s *ServiceContext) GetContext() context.Context {
	return s.context
}

func (s *ServiceContext) GetPagination() *Pagination {
	if s.pagination == nil {
		return new(Pagination)
	}
	return s.pagination
}

func (s *ServiceContext) GetSortParam() *SortParam {
	return s.sortParam
}

func (s *ServiceContext) GetAudiTrail() *AudiTrail {
	if s.audiTrail == nil {
		return new(AudiTrail)
	}
	return s.audiTrail
}

func (s *ServiceContext) GetUserClaim() *domain.UserCustomClaims {
	if s.userClaim == nil {
		return new(domain.UserCustomClaims)
	}
	return s.userClaim
}

func (s *ServiceContext) GetUserSession() *domain.UserSession {
	if s.userSession == nil {
		return new(domain.UserSession)
	}
	return s.userSession
}

func (s *ServiceContext) SetLogger(logger *zap.Logger) {
	s.logger = logger
}

func (s *ServiceContext) Logger() *zap.Logger {
	return s.logger
}

func (s *ServiceContext) BeginTransaction(db *gorm.DB) *gorm.DB {
	s.database = db.Begin()
	return s.database
}

func (s *ServiceContext) GetDB() *gorm.DB {
	return s.database
}
