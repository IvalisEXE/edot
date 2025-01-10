package service

import (
	"strconv"
	"time"
	"usersvc/internal/core/domain"
	"usersvc/internal/core/port"
	"usersvc/pkg/cache"
	"usersvc/pkg/common"

	"github.com/spf13/viper"
	"go.elastic.co/apm"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository port.Repository
	cache      cache.CacheManager
}

// New creates a new user service
func New(
	repository port.Repository,
	cache cache.CacheManager,
) port.Service {
	return &service{
		repository: repository,
		cache:      cache,
	}
}

func (s *service) Login(ctx common.ServiceContextManager, request domain.UserLoginRequest) (*domain.UserTokenData, error) {
	span, context := apm.StartSpan(ctx.GetContext(), "service.user.Login", "service")
	defer span.End()

	ctx.SetContext(context)

	user, err := s.repository.FindByPhone(ctx, request.Phone)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domain.ErrPhoneOrEmailNotRegistered
	}

	if user.Password == nil {
		return nil, domain.ErrPasswordNotSet
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(request.Password)); err != nil {
		return nil, domain.ErrPasswordIncorrect
	}

	if err := s.repository.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, err
	}

	token, err := user.GenerateToken()
	if err != nil {
		return nil, err
	}

	if err = s.SaveSession(ctx, user, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *service) SaveSession(ctx common.ServiceContextManager, user *domain.User, token *domain.UserTokenData) error {

	if token != nil {
		return s.saveSessionID(ctx, token)
	}

	if user != nil {
		return s.saveSessionProfile(ctx, user)
	}

	return nil
}

func (s *service) saveSessionID(ctx common.ServiceContextManager, token *domain.UserTokenData) error {
	userSession := domain.UserSessionInfo{
		IpAddress: ctx.GetAudiTrail().IpAddress,
		UserAgent: ctx.GetAudiTrail().UserAgent,
		ExpiredAt: token.ExpiredAt.Format(time.RFC3339),
	}
	return s.cache.Set(ctx.GetContext(), token.SessionID, userSession,
		time.Duration(viper.GetInt("REDIS_SESSION_TTL")))
}

func (s *service) saveSessionProfile(ctx common.ServiceContextManager, user *domain.User) error {
	key := "user:" + strconv.Itoa(user.ID)
	return s.cache.Set(ctx.GetContext(), key, user.ToUserSession(),
		time.Duration(viper.GetInt("REDIS_SESSION_TTL")))
}

func (s *service) FindByID(ctx common.ServiceContextManager, id int) (*domain.User, error) {
	span, context := apm.StartSpan(ctx.GetContext(), "service.user.FindByID", "service")
	defer span.End()

	ctx.SetContext(context)

	return s.repository.FindByID(ctx, id)
}
