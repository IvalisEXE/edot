package handler

import (
	"net/http"
	"warehousesvc/internal/core/domain"
	"warehousesvc/internal/core/port"
	"warehousesvc/pkg/common"
	"warehousesvc/pkg/response"
	"warehousesvc/pkg/validator"

	"github.com/labstack/echo/v4"
	"go.elastic.co/apm"
)

type handler struct {
	service   port.Service
	validator validator.Validator
}

func New(
	service port.Service,
	validator validator.Validator,
) port.Handler {
	return &handler{
		service:   service,
		validator: validator,
	}
}

func (h *handler) FindAllStock(c echo.Context) error {
	span, context := apm.StartSpan(c.Request().Context(), "handler.warehouse.FindAllStock", "handler")
	defer span.End()

	ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)
	ctx.SetContext(context)

	var (
		request domain.FindAllStockRequest
	)

	if err := c.Bind(&request); err != nil {
		return err
	}

	if err := h.validator.Validate(&request); err != nil {
		return err
	}

	ctx.SetPagination(common.AssertPaginationPayload(request.QueryParam))
	ctx.SetSortParam(common.AssertSortPayload(request.QueryParam))

	stocks, err := h.service.FindAllStock(ctx, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Response{
		Message: "Get list stock is successful",
		Data:    stocks,
	})
}

func (h *handler) TransferStock(c echo.Context) error {
	span, context := apm.StartSpan(c.Request().Context(), "handler.warehouse.TransferStock", "handler")
	defer span.End()

	ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)
	ctx.SetContext(context)

	var (
		request domain.TransferStockRequest
	)

	if err := c.Bind(&request); err != nil {
		return err
	}

	if err := h.validator.Validate(&request); err != nil {
		return err
	}

	err := h.service.TransferStock(ctx, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Response{
		Message: "Transfer stock is successful",
	})
}

func (h *handler) UnactiveWarehouse(c echo.Context) error {
	span, context := apm.StartSpan(c.Request().Context(), "handler.warehouse.UnactiveWarehouse", "handler")
	defer span.End()

	ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)
	ctx.SetContext(context)

	var (
		request domain.ParamID
	)

	if err := c.Bind(&request); err != nil {
		return err
	}

	if err := h.validator.Validate(&request); err != nil {
		return err
	}

	err := h.service.UnactiveWarehouse(ctx, request.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Response{
		Message: "Unactive warehouse is successful",
	})
}
