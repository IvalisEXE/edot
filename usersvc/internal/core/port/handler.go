package port

import "github.com/labstack/echo/v4"

type Handler interface {
	Login(c echo.Context) error
	Authenticate(c echo.Context) error
}
