package handler

import (
	"net/http"
	"ordersvc/internal/core/service"
	"ordersvc/internal/infrastructure/repository"
	"ordersvc/internal/interfaces/api/handler"
	"ordersvc/pkg/cache"
	"ordersvc/pkg/common"
	"ordersvc/pkg/config"
	"ordersvc/pkg/database"
	"ordersvc/pkg/logger"
	"ordersvc/pkg/middleware"
	"ordersvc/pkg/response"
	"ordersvc/pkg/validator"

	"github.com/labstack/echo/v4"
	"go.elastic.co/apm"
)

type Provider struct {
	API   *echo.Echo
	ENV   *config.Container
	Cache cache.CacheManager
}

// API is the API routes
func API(app Provider) {

	// Load configuration
	cache := app.Cache
	db := database.DB

	// Validator
	validator := validator.New()

	// Repository
	repositoryImpl := repository.New(db)

	// Service
	serviceImpl := service.New(repositoryImpl, cache)

	// Handler
	handlerImpl := handler.New(serviceImpl, validator)

	// Routes
	r := app.API
	Common(r)

	routes := r.Group("/order", middleware.HeaderMiddleware)
	routes.POST("", handlerImpl.Order)

}

// Common is the common routes
func Common(r *echo.Echo) {
	r.GET("/", func(c echo.Context) error {
		response := response.Response{
			Message: "Welcome to the User API",
		}

		return c.JSON(http.StatusOK, response)
	})

	r.GET("/health", func(c echo.Context) error {
		span, context := apm.StartSpan(c.Request().Context(), "handler.Health", "handler")
		defer span.End()

		ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)
		ctx.SetContext(context)

		ctx.Logger().Info("Checking server health", logger.WithData(map[string]interface{}{"status": "OK"}))

		response := response.Response{
			Message: "OK",
		}

		return c.JSON(http.StatusOK, response)
	})
}
