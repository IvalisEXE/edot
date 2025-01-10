package domain

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrWarehouseNotFound   = echo.NewHTTPError(http.StatusBadRequest, "Warehouse not found")
	ErrWarehouseIsUnactive = echo.NewHTTPError(http.StatusBadRequest, "Warehouse is already unactivated")
	ErrItemNotFound        = echo.NewHTTPError(http.StatusBadRequest, "Item not found")
	ErrInsufficientItem    = echo.NewHTTPError(http.StatusBadRequest, "Insufficient item")
)
