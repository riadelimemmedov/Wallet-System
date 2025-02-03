package config

// AppConfig holds application configuration settings
type AppConfig struct {
	ConfigFilePath string
}

// Environment defines environment-specific settings
type Environment struct {
	AppEnv string
	Prefix string
}

// Pre-defined environment configurations
var (
	DevEnvironment = Environment{
		AppEnv: "dev",
		Prefix: "DEV",
	}
	TestEnvironment = Environment{
		AppEnv: "test",
		Prefix: "TEST",
	}
)
