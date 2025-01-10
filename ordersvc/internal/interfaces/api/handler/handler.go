package handler

import (
	"net/http"
	"ordersvc/internal/core/domain"
	"ordersvc/internal/core/port"
	"ordersvc/pkg/common"
	"ordersvc/pkg/response"
	"ordersvc/pkg/validator"

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

func (h *handler) Order(c echo.Context) error {
	span, context := apm.StartSpan(c.Request().Context(), "handler.order.Order", "handler")
	defer span.End()

	ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)
	ctx.SetContext(context)

	var (
		request domain.OrderRequest
	)

	if err := c.Bind(&request); err != nil {
		return err
	}

	if err := h.validator.Validate(&request); err != nil {
		return err
	}

	err := h.service.Order(ctx, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Response{
		Message: "Order is successful",
	})
}
