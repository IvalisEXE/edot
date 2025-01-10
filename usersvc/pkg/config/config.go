package config

import (
	"log"

	"github.com/spf13/viper"
)

type (
	Container struct {
		// App configuration
		App *App

		// DB configuration
		DB *DB

		// Redis configuration
		Redis *Redis

		// Logger configuration
		Logger *Logger
	}

	App struct {
		Name string
		Port int
	}

	DB struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
		Debug    string
		SSLMode  string
	}

	Redis struct {
		Address    string
		Password   string
		DB         int
		Port       int
		SessionTTL int
	}

	Logger struct {
		Level   string
		Handler string
	}
)

// New creates a new configuration
func New(path string) *Container {
	var cfg Config

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	return &Container{
		App: &App{
			Name: cfg.AppName,
			Port: cfg.AppPort,
		},
		DB: &DB{
			Host:     cfg.DBHost,
			Port:     cfg.DBPort,
			Username: cfg.DBUser,
			Password: cfg.DBPassword,
			Database: cfg.DBName,
			Debug:    cfg.DBDebug,
			SSLMode:  cfg.SSLMode,
		},
		Redis: &Redis{
			Address:    cfg.RedisAddress,
			Password:   cfg.RedisPassword,
			DB:         cfg.RedisDB,
			Port:       cfg.RedisPort,
			SessionTTL: cfg.RedisSessionTTL,
		},
		Logger: &Logger{
			Level:   cfg.LogLevel,
			Handler: cfg.LogHandler,
		},
	}
}

type Config struct {
	// App configuration
	AppName string `mapstructure:"APP_NAME"`
	AppPort int    `mapstructure:"APP_PORT"`

	// DB configuration
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBDebug    string `mapstructure:"DB_DEBUG"`
	SSLMode    string `mapstructure:"SSL_MODE"`

	// Redis configuration
	RedisAddress    string `mapstructure:"REDIS_HOST_URL"`
	RedisPassword   string `mapstructure:"REDIS_PASSWORD"`
	RedisDB         int    `mapstructure:"REDIS_DB"`
	RedisPort       int    `mapstructure:"REDIS_PORT"`
	RedisSessionTTL int    `mapstructure:"REDIS_SESSION_TTL"`

	// Logger configuration
	LogLevel   string `mapstructure:"LOG_LEVEL"`
	LogHandler string `mapstructure:"LOG_HANDLER"`
}
