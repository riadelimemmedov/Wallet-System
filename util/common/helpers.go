package common

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgtype"
)

// !SetNumeric safely convert string to decimal format
func SetNumeric(value string) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	err := numeric.Set(value)
	if err != nil {
		return numeric, fmt.Errorf("failed to set numeric value for exchange rate")
	}
	return numeric, nil
}

// !GetEnvAsInt safely gets on environment variable as integer
func GetEnvAsInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return defaultVal
}

// !GetEnvAsDuration safely gets on environment variable as duration
func GetEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if v, err := time.ParseDuration(value); err == nil {
			return v
		}
	}
	return defaultVal
}
