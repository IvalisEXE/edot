package port

import "github.com/labstack/echo/v4"

type Handler interface {
	Order(c echo.Context) error
}
