package handler

import (
	"net/http"
	"usersvc/internal/core/domain"
	"usersvc/internal/core/port"
	"usersvc/pkg/common"
	"usersvc/pkg/response"
	"usersvc/pkg/validator"

	"github.com/labstack/echo/v4"
	"go.elastic.co/apm"
)

// handler is the user handler
type handler struct {
	service   port.Service
	validator validator.Validator
}

// New creates a new user handler
func New(
	service port.Service,
	validator validator.Validator,
) port.Handler {
	return &handler{
		service:   service,
		validator: validator,
	}
}

func (h *handler) Login(c echo.Context) error {
	span, context := apm.StartSpan(c.Request().Context(), "handler.user.Login", "handler")
	defer span.End()

	var (
		request domain.UserLoginRequest
		err     error
	)

	if err = c.Bind(&request); err != nil {
		return err
	}

	if err = h.validator.Validate(request); err != nil {
		return err
	}

	ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)
	ctx.SetContext(context)

	token, err := h.service.Login(ctx, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Response{
		Message: "Login is successful",
		Data:    token,
	})
}

func (h *handler) Authenticate(c echo.Context) error {
	span, context := apm.StartSpan(c.Request().Context(), "handler.user.Authenticate", "handler")
	defer span.End()

	ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)
	ctx.SetContext(context)

	user, err := h.service.FindByID(ctx, ctx.GetUserSession().ID)
	if err != nil {
		return err
	}

	if user == nil {
		return domain.ErrUserNotFound
	}

	return c.JSON(http.StatusOK, response.Response{
		Message: "Validate token is successful",
		Data: domain.User{
			ID:       user.ID,
			ShopID:   user.ShopID,
			Name:     user.Name,
			Phone:    user.Phone,
			Platform: user.Platform,
		},
	})
}
