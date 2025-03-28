package setup

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/riad/banksystemendtoend/pkg/redis"
	"github.com/riad/banksystemendtoend/util/config"
	"go.uber.org/zap"
)

var (
	redisClient *redis.Client
	log         = zap.L()
)

// InitializeRedis initializes a new Redis client from start to end
func InitializeRedis(config config.AppConfig, envPrefix string) (*redis.Client, error) {
	if config.ConfigFilePath != "" {
		if err := godotenv.Load(config.ConfigFilePath); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	client, err := SetupRedisClient(envPrefix)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// SetupRedisClient creates a new Redis client
func SetupRedisClient(envPrefix string) (*redis.Client, error) {
	redisConfig := redis.Config{
		Host:            os.Getenv(envPrefix + "_REDIS_HOST"),
		Port:            os.Getenv(envPrefix + "_REDIS_PORT"),
		Password:        os.Getenv(envPrefix + "_REDIS_PASSWORD"),
		DB:              GetEnvAsInt(os.Getenv(envPrefix+"_REDIS_DB"), 0),
		PoolSize:        GetEnvAsInt(os.Getenv(envPrefix+"_REDIS_POOL_SIZE"), 10),
		MinIdleConns:    GetEnvAsInt(os.Getenv(envPrefix+"_REDIS_MIN_IDLE_CONNS"), 5),
		MaxConnLifetime: GetEnvAsDuration(os.Getenv(envPrefix+"_REDIS_MAX_CONN_LIFETIME"), 10*time.Minute),
		IdleTimeout:     GetEnvAsDuration(os.Getenv(envPrefix+"_REDIS_IDLE_TIMEOUT"), 5*time.Minute),
		DialTimeout:     GetEnvAsDuration(os.Getenv(envPrefix+"_REDIS_DIAL_TIMEOUT"), 5*time.Second),
		ReadTimeout:     GetEnvAsDuration(os.Getenv(envPrefix+"_REDIS_READ_TIMEOUT"), 3*time.Second),
		WriteTimeout:    GetEnvAsDuration(os.Getenv(envPrefix+"_REDIS_WRITE_TIMEOUT"), 3*time.Second),
	}
	if redisConfig.Host == "" {
		redisConfig.Host = "localhost"
	}
	if redisConfig.Port == "" {
		redisConfig.Port = "6379"
	}

	maxRetries := GetEnvAsInt(os.Getenv(envPrefix+"_REDIS_MAX_RETRIES"), 3)
	var client *redis.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		client, err = redis.NewClient(redisConfig)
		if err == nil {
			zap.L().Info("Created mew redis client successfully")
			break
		}
		zap.L().Error("Failed to connect to Redis", zap.Error(err))
		if i < maxRetries-1 {
			zap.L().Info("Retrying to connect to Redis", zap.Int("retry", i+1))
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis after %d attempts: %w", maxRetries, err)
	}

	zap.L().Info("Redis connection established",
		zap.String("host", redisConfig.Host),
		zap.String("port", redisConfig.Port))
	return client, nil

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

// GetRedisClient returns the initialized Redis client
func GetRedisClient() (*redis.Client, error) {
	if redisClient == nil {
		log.Error("Redis client is not initialized")
		return nil, fmt.Errorf("redis client is not initialized")
	}
	return redisClient, nil
}
