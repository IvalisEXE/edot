package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	postgres "go.elastic.co/apm/module/apmgormv2/v2/driver/postgres"

	"ordersvc/pkg/config"
)

// DB is the database connection
var DB *gorm.DB

// Load creates a new database connection
func Load(config *config.DB) {
	var (
		enableDebug logger.LogLevel = logger.Silent
	)

	if config.Debug == "true" {
		enableDebug = logger.Info
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC connect_timeout=10",
		config.Host, config.Username, config.Password, config.Database, config.Port, config.SSLMode,
	)

	// Open the connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  enableDebug, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      false,       // Don't include params in the SQL log
				Colorful:                  true,
			},
		),
	})

	if err != nil {
		log.Fatal(err)
	}

	DB = db

	log.Println("Database connection successfully")
}

// Close closes the connection to the database
func Close() error {
	// Get the underlying sql.DB object
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// Close the sql.DB object
	err = sqlDB.Close()
	if err != nil {
		return err
	}

	// Set the DB object to nil
	DB = nil

	return nil
}
