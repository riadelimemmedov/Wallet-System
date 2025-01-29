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

// !RandomNumeric generates a random numeric value between 0-10000 with 2 decimal places
func RandomNumeric() (pgtype.Numeric, error) {
	return SetNumeric(fmt.Sprintf("%.2f", RandomFloat(0, 10000)))
}

// ! RandomInterestRate generates a random interest rate between 1.1%-99.99%
func RandomInterestRate() (pgtype.Numeric, error) {
	rate := RandomFloat(1.1, 99.99)
	return SetNumeric(fmt.Sprintf("%.2f", rate))
}
