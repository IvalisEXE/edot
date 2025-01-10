package middleware

import (
	"strconv"
	"warehousesvc/internal/core/domain"
	"warehousesvc/pkg/common"

	"github.com/labstack/echo/v4"
)

func HeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)
		ctx.SetContext(c.Request().Context())

		userID := c.Request().Header.Get("X-User-Id")
		userIDConv, _ := strconv.Atoi(userID)
		shopID := c.Request().Header.Get("X-Shop-Id")
		shopIDConv, _ := strconv.Atoi(shopID)

		ctx.SetUserHeader(&domain.UserHeader{
			ID:     userIDConv,
			ShopID: shopIDConv,
			Name:   c.Request().Header.Get("X-Username"),
			Phone:  c.Request().Header.Get("X-Phone"),
		})
		return next(c)
	}
}
