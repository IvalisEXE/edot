package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"go.elastic.co/apm/module/apmechov4/v2"
	"go.uber.org/zap"

	"ordersvc/pkg/cache"
	"ordersvc/pkg/config"
	"ordersvc/pkg/database"

	routes "ordersvc/internal/interfaces/api"
	myLogger "ordersvc/pkg/logger"
	mdw "ordersvc/pkg/middleware"
)

// API
var API = &cobra.Command{
	Use:   "api",
	Short: "Run API Services",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration from the current directory
		config := config.New(".")

		// Create a new database connection
		database.Load(config.DB)
		defer database.Close()

		// Create a new Redis connection as cache manager
		redis := cache.NewRedis(config.Redis)

		// Initialize logger
		logger, err := myLogger.New(config.Logger)
		if err != nil {
			panic(err)
		}

		// Create new Echo instance
		api := echo.New()
		api.HideBanner = true
		api.HidePort = true
		api.HTTPErrorHandler = mdw.ErrorHandler

		// Set up middleware
		api.Use(middleware.Recover())
		api.Use(middleware.Secure())
		api.Use(middleware.CORS())
		api.Use(mdw.ContextManager(logger))
		api.Use(apmechov4.Middleware())

		// Set up API routes
		routes.API(routes.Provider{
			API:   api,
			ENV:   config,
			Cache: redis,
		})

		// Run API in a goroutine
		go func() {
			logger.Info("Starting API server")
			addr := fmt.Sprintf(":%d", config.App.Port)
			if err := api.Start(addr); err != nil && err != http.ErrServerClosed {
				logger.Fatal("API server failure", zap.Error(err))
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		// Graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Shutdown API server
		if err := api.Shutdown(ctx); err != nil {
			logger.Fatal("API server shutdown failed", zap.Error(err))
		}

		logger.Info("Server shutdown gracefully")
	},
}
