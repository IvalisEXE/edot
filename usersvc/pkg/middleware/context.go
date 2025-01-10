package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"usersvc/pkg/common"
)

func ContextManager(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the request from the context
			request := c.Request()

			// Add request-specific fields to the logger
			requestLogger := logger.With(
				zap.String("uri", request.RequestURI),
				zap.String("method", request.Method),
			)

			// Create a new ServiceContext with the logger and audi trail
			serviceContext := common.NewServiceContext()
			serviceContext.SetLogger(requestLogger)
			serviceContext.SetContext(c.Request().Context())
			serviceContext.SetAudiTrail(common.NewAudiTrail(request, c.RealIP()))

			// Set the ServiceContext to the request context
			c.Set(common.KEY_API_CONTEXT, serviceContext)

			// Process the request
			err := next(c)

			// Return the error
			return err
		}
	}
}
