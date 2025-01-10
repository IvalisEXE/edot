package cron

import (
	"ordersvc/internal/core/service"
	"ordersvc/internal/infrastructure/repository"
	"ordersvc/pkg/cache"
	"ordersvc/pkg/common"
	"ordersvc/pkg/config"
	"ordersvc/pkg/database"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var cfg *config.Container

func init() {
	// Load configuration
	cfg = config.New(".")
}

// Cron
var Cron = &cobra.Command{
	Use:   "cron",
	Short: "Run Cron Services",
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new database connection
		database.Load(cfg.DB)
		defer database.Close()

		// Create a new Redis connection as cache manager
		cache := cache.NewRedis(cfg.Redis)

		// Repository
		repositoryImpl := repository.New(database.DB)

		// Service
		serviceImpl := service.New(repositoryImpl, cache)

		c := cron.New()

		// Every 1 minute
		c.AddFunc("*/1 * * * *", func() {
			serviceImpl.LockStock(&common.ServiceContext{})
			serviceImpl.ReleaseStock(&common.ServiceContext{})
		})

		c.Start()

		select {}
	},
}
