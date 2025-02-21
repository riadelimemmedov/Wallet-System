package middleware

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/riad/banksystemendtoend/api/middleware/logger/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger struct holds the logging configuration and zap logger instance
type Logger struct {
	config    config.LoggerConfig
	ZapLogger *zap.Logger
}

// NewLogger creates and configures a new Logger instance
func NewLogger(cfg config.LoggerConfig) (*Logger, error) {
	dir := filepath.Dir(cfg.Filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	if _, err := os.Stat(cfg.Filename); os.IsNotExist(err) {
		file, err := os.Create(cfg.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file: %v", err)
		}
		file.Close()
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{cfg.Filename}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{
		config:    cfg,
		ZapLogger: logger,
	}, nil
}

// Logger returns a Gin middleware function that logs HTTP request details
func (l *Logger) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		l.ZapLogger.Info("http_request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("size", int64(c.Writer.Size())),
			zap.Duration("latency", time.Since(start)),
			zap.String("errors", c.Errors.String()),
		)
	}
}
