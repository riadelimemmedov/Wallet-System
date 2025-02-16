package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/riad/banksystemendtoend/api/middleware/logger/config"
	"github.com/riad/banksystemendtoend/api/middleware/logger/entity"
	"github.com/riad/banksystemendtoend/api/middleware/logger/formatter"
	"github.com/riad/banksystemendtoend/api/middleware/logger/storage"
)

// Logger handles HTTP request logging with configurable storage and formatting
type Logger struct {
	config    config.LoggerConfig
	storage   *storage.FileStorage
	formatter formatter.LogFormatter
}

// NewLogger creates a logger instance with given config
func NewLogger(cfg config.LoggerConfig) (*Logger, error) {
	storage, err := storage.NewFileStorage(cfg.Filename)
	if err != nil {
		return nil, err
	}
	return &Logger{
		config:    cfg,
		storage:   storage,
		formatter: formatter.NewDefaultFormatter(),
	}, nil
}

// LoggingMiddleware returns a Gin middleware function that logs HTTP requests
func (l *Logger) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		entry := entity.NewLogEntry(
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			c.Writer.Size(),
			time.Since(start),
			c.Errors.String(),
		)
		formattedLog := l.formatter.Format(entry)
		l.storage.Write(formattedLog)
	}
}
