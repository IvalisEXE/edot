package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"usersvc/internal/core/domain"
	"usersvc/internal/core/port"
	"usersvc/pkg/cache"
	"usersvc/pkg/common"
	"usersvc/pkg/response"
)

type Auth struct {
	service port.Service
	cache   cache.CacheManager
}

func NewAuth(service port.Service, cache cache.CacheManager) *Auth {
	return &Auth{
		service: service,
		cache:   cache,
	}
}

func (m *Auth) ValidateAuthDashboard() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(viper.GetString("JWT_SECRET_DASHBOARD")),
		SuccessHandler: func(c echo.Context) {
			// Get the user from the context
			ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)

			// Get the user from the token
			payload := c.Get("user").(*jwt.Token)
			claims := payload.Claims.(jwt.MapClaims)

			// Set the user claim
			claim := domain.UserCustomClaims{
				ID:        int(claims["id"].(float64)),
				ExpiredAt: int(claims["exp"].(float64)),
			}

			ctx.SetUserClaim(&claim)

			// Get user from cache
			user, err := m.cache.Get(ctx.GetContext(), "user:"+strconv.Itoa(claim.ID))

			// Regenerate cache profile user if not exists
			if user == "" || err != nil {
				user = m.Refresh(ctx, claim)
			}

			// Check if user is empty
			if user == "" {
				m.UnAuthorized(c)
				return
			}

			// Set the user session
			var session domain.UserSession
			err = json.Unmarshal([]byte(user), &session)
			if err != nil {
				m.UnAuthorized(c)
				return
			}

			ctx.SetUserSession(&session)
		},
	})
}

func (m *Auth) ValidateAuthCustomer() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(viper.GetString("JWT_SECRET_CUSTOMER")),
		SuccessHandler: func(c echo.Context) {
			// Get the user from the context
			ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)

			// Get the user from the token
			payload := c.Get("user").(*jwt.Token)
			claims := payload.Claims.(jwt.MapClaims)

			// Set the user claim
			claim := domain.UserCustomClaims{
				ID:        int(claims["id"].(float64)),
				ExpiredAt: int(claims["exp"].(float64)),
			}

			ctx.SetUserClaim(&claim)

			// Get user from cache
			user, err := m.cache.Get(ctx.GetContext(), "user:"+strconv.Itoa(claim.ID))

			// Regenerate cache profile user if not exists
			if user == "" || err != nil {
				user = m.Refresh(ctx, claim)
			}

			// Check if user is empty
			if user == "" {
				m.UnAuthorized(c)
				return
			}

			// Set the user session
			var session domain.UserSession
			err = json.Unmarshal([]byte(user), &session)
			if err != nil {
				m.UnAuthorized(c)
				return
			}

			ctx.SetUserSession(&session)
		},
	})
}

// Get user from DB and save it into redis if not exists
func (m *Auth) Refresh(ctx common.ServiceContextManager, claim domain.UserCustomClaims) string {
	// Get user from DB
	user, err := m.service.FindByID(ctx, claim.ID)
	if err != nil {
		return ""
	}

	// Save user to cache
	if err := m.service.SaveSession(ctx, user, nil); err != nil {
		return ""
	}

	// Get user from cache
	str, err := m.cache.Get(ctx.GetContext(), "user:"+strconv.Itoa(claim.ID))
	if err != nil {
		return ""
	}

	return str
}

func (m *Auth) UnAuthorized(c echo.Context) {
	err := c.JSON(http.StatusUnauthorized, response.Response{
		Message: echo.ErrUnauthorized.Error(),
	})

	if err != nil {
		ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)
		ctx.Logger().Error("Failed to send response", zap.Error(err))
	}
}
