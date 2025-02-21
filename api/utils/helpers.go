package utils

import (
	"fmt"

	"github.com/joho/godotenv"
)

// LoadEnvFile attempts to load environment variables from a .env file in the project root.
func LoadEnvFile() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}
