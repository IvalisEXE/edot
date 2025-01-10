package port

import "github.com/labstack/echo/v4"

type Handler interface {
	FindAllStock(c echo.Context) error
	TransferStock(c echo.Context) error
	UnactiveWarehouse(c echo.Context) error
}
