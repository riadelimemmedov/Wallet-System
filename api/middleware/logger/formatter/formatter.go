package formatter

import (
	"fmt"

	"github.com/riad/banksystemendtoend/api/middleware/logger/entity"
)

type LogFormatter interface {
	Format(entry entity.LogEntry) string
}

type DefaultFormatter struct {
	TimeFormat string
}

func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{
		TimeFormat: "2006-01-02 15:04:05",
	}
}

func (f *DefaultFormatter) Format(entry entity.LogEntry) string {
	return fmt.Sprintf("[%s] %s %s %s %d %d %v %s",
		entry.Timestamp.Format(f.TimeFormat),
		entry.Method,
		entry.Path,
		entry.ClientIP,
		entry.Status,
		entry.Size,
		entry.Latency,
		entry.Errors,
	)
}
