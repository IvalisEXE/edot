package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"usersvc/pkg/config"
)

// New creates a new logger
func New(cfg *config.Logger) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.DebugLevel
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %v", err)
	}

	cfgHandler, err := getHandlers(cfg.Handler)
	if err != nil {
		return nil, fmt.Errorf("failed to get handlers: %v", err)
	}

	config := zap.Config{
		OutputPaths:      cfgHandler,
		ErrorOutputPaths: []string{"stderr"},
		Encoding:         "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			CallerKey:     "file",
			StacktraceKey: "stacktrace",
			EncodeCaller:  zapcore.ShortCallerEncoder,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime:    zapcore.RFC3339TimeEncoder,
		},
		Level: zap.NewAtomicLevelAt(level),
	}

	return config.Build()
}

// WithData returns a Field for structured logging
func WithData(data interface{}) zapcore.Field {
	return zap.Any("data", data)
}

// WithError returns a Field for error logging
func WithError(err error) zapcore.Field {
	return zap.Error(err)
}

func getHandlers(input string) ([]string, error) {
	if input == "" {
		return []string{"stdout"}, nil
	}

	var handlers []string
	for _, h := range strings.Split(input, ",") {
		switch strings.TrimSpace(h) {
		case "file":
			logFile := filepath.Join("logs", fmt.Sprintf("log-%s.log", time.Now().Format("2006-01-02")))
			handlers = append(handlers, logFile)
		case "stdout":
			handlers = append(handlers, h)
		default:
			return nil, fmt.Errorf("unknown handler: %s", h)
		}
	}

	return handlers, nil
}
