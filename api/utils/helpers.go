package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/riad/banksystemendtoend/util/config"
	"golang.org/x/crypto/bcrypt"
)

// LoadEnvFile attempts to load environment variables from a .env file in the project root.
func LoadEnvFile() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword checks if a password matches a hashed password
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsDuplicateError checks if an error is a duplicate error
func IsDuplicateError(err error, data string, instanceType ...string) (bool, error) {
	isDuplicate := strings.Contains(err.Error(), "SQLSTATE 23505")
	if !isDuplicate {
		return false, nil
	}
	if len(instanceType) == 0 || instanceType[0] == "" {
		return true, fmt.Errorf("record already exists %s", data)
	}
	return true, fmt.Errorf("%s already exists - %s", instanceType[0], data)
}

// IsForeignKeyError checks if an error is a foreign key error
func IsForeignKeyError(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23503")
}

// IsUniqueViolationError checks if an error is a unique violation error
func IsUniqueViolationError(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}

// GetValidAccountTypesMessage returns a message with valid account types
func GetValidAccountTypesMessage() string {
	types := []string{}

	for _, v := range config.AccountTypesMap {
		types = append(types, v)
	}
	return fmt.Sprintf("invalid account type: must be one of %s", strings.Join(types, ", "))
}

// GetEnvAsInt returns the value of the environment variable as an integer
func GetEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		intVal, err := strconv.Atoi(value)
		if err == nil {
			return intVal
		}
	}
	return defaultVal
}

// GetEnvAsDuration returns the value of the environment variable as a duration
func GetEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		durationVal, err := time.ParseDuration(value)
		if err == nil {
			return durationVal
		}
	}
	return defaultVal
}
