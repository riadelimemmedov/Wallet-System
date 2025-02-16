package config

// LoggerConfig holds configuration settings for the logger
type LoggerConfig struct {
	Filename   string
	TimeFormat string
}

// NewDefaultConfig returns default logger settings with
func NewDefaultConfig() LoggerConfig {
	return LoggerConfig{
		Filename:   "logs/api.log",
		TimeFormat: "2006-01-02 15:04:05",
	}
}
