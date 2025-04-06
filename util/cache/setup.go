package setup

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/riad/banksystemendtoend/api/utils"
	logger "github.com/riad/banksystemendtoend/pkg/log"
	"github.com/riad/banksystemendtoend/pkg/redis"
	"github.com/riad/banksystemendtoend/util/config"
	"go.uber.org/zap"
)

var (
	redisClient *redis.Client
)

// InitializeRedis initializes a new Redis client from start to end
func InitializeRedis(config config.AppConfig, _ string) (*redis.Client, error) {
	if config.ConfigFilePath != "" {
		if err := godotenv.Load(config.ConfigFilePath); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	client, err := SetupRedisClient()
	if err != nil {
		return nil, err
	}
	redisClient = client
	return client, nil
}

// SetupRedisClient creates a new Redis client
func SetupRedisClient() (*redis.Client, error) {
	redisConfig := redis.Config{
		Host:            os.Getenv("REDIS_HOST"),
		Port:            os.Getenv("REDIS_PORT"),
		Password:        os.Getenv("REDIS_PASSWORD"),
		DB:              utils.GetEnvAsInt("REDIS_DB", 0),
		PoolSize:        utils.GetEnvAsInt("REDIS_POOL_SIZE", 10),
		MinIdleConns:    utils.GetEnvAsInt("REDIS_MIN_IDLE_CONNS", 3),
		MaxConnLifetime: utils.GetEnvAsDuration("REDIS_MAX_CONN_LIFETIME", 30*time.Minute),
		IdleTimeout:     utils.GetEnvAsDuration("REDIS_IDLE_TIMEOUT", 5*time.Minute),
		DialTimeout:     utils.GetEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:     utils.GetEnvAsDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		WriteTimeout:    utils.GetEnvAsDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		MaxRetries:      utils.GetEnvAsInt("REDIS_MAX_RETRIES", 3),
	}
	if redisConfig.Host == "" {
		redisConfig.Host = "localhost"
	}
	if redisConfig.Port == "" {
		redisConfig.Port = "6379"
	}
	client, err := redis.NewClient(redisConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to setup Redis client: %w", err)
	}

	if err := client.CheckRedisConnection(); err != nil {
		return nil, err
	}

	logger.GetLogger().Info("Redis connection established",
		zap.String("host", redisConfig.Host),
		zap.String("port", redisConfig.Port))
	return client, nil
}

// GetRedisClient returns the initialized Redis client
func GetRedisClient() *redis.Client {
	if redisClient == nil {
		logger.GetLogger().Error("Redis client is not initialized")
		return nil
	}
	return redisClient
}
