package handler

import (
	"net/http"
	"warehousesvc/internal/core/service"
	"warehousesvc/internal/infrastructure/repository"
	"warehousesvc/internal/interfaces/api/handler"
	"warehousesvc/pkg/cache"
	"warehousesvc/pkg/common"
	"warehousesvc/pkg/config"
	"warehousesvc/pkg/database"
	"warehousesvc/pkg/logger"
	"warehousesvc/pkg/middleware"
	"warehousesvc/pkg/response"
	"warehousesvc/pkg/validator"

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

	routes := r.Group("/warehouse", middleware.HeaderMiddleware)
	routes.GET("/stock", handlerImpl.FindAllStock)
	routes.POST("/stock-transfer", handlerImpl.TransferStock)
	routes.POST("/unactive/:id", handlerImpl.UnactiveWarehouse)

}

// Common is the common routes
func Common(r *echo.Echo) {
	r.GET("/", func(c echo.Context) error {
		response := response.Response{
			Message: "Welcome to the Product API",
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
