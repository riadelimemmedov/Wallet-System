package logger

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	err     error
	log     *zap.Logger
	logPath string
	once    sync.Once
	mu      sync.Mutex
)

// Initialize sets up the global logger
func Initialize(env string) *zap.Logger {
	once.Do(func() {
		var config zap.Config

		if env == "production" {
			config = zap.NewProductionConfig()
			logPath = "logs/production"
			config.OutputPaths = []string{filepath.Join(logPath, "app.log")}
		} else {
			config = zap.NewProductionConfig()
			logPath = "logs/development"
			config.OutputPaths = []string{"stdout", filepath.Join(logPath, "app.log")}
		}

		if err := os.MkdirAll(logPath, 0755); err != nil {
			panic(err)
		}

		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.StacktraceKey = "stacktrace"

		log, err = config.Build(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
		zap.ReplaceGlobals(log)
	})
	return log

}

// Get returns the global logger
func GetLogger() *zap.Logger {
	mu.Lock()
	defer mu.Unlock()

	if log == nil {
		Initialize(os.Getenv("APP_ENV"))
	}
	return log
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() {
	mu.Lock()
	defer mu.Unlock()

	if log != nil {
		_ = log.Sync()
	}
}
