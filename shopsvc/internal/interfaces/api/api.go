package handler

import (
	"net/http"
	"shopsvc/internal/core/service"
	"shopsvc/internal/infrastructure/repository"
	"shopsvc/internal/interfaces/api/handler"
	"shopsvc/pkg/cache"
	"shopsvc/pkg/common"
	"shopsvc/pkg/config"
	"shopsvc/pkg/database"
	"shopsvc/pkg/logger"
	"shopsvc/pkg/response"
	"shopsvc/pkg/validator"

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
	_ = handler.New(serviceImpl, validator)

	// Routes
	r := app.API
	Common(r)

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
