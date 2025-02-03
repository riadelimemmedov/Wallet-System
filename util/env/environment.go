package env

import (
	"fmt"

	"github.com/riad/banksystemendtoend/util/config"
)

// ! Returns config with path
func NewAppEnvironmentConfig(environment string) (config.AppConfig, error) {
	appConfigMap := map[string]config.AppConfig{
		"test": {
			ConfigFilePath: "../../../.env.test",
		},
		"dev": {
			ConfigFilePath: "../../../.env.dev",
		},
	}
	appConfig, exists := appConfigMap[environment]
	if !exists {
		return config.AppConfig{}, fmt.Errorf("environment %s not found", environment)
	}
	return appConfig, nil
}
