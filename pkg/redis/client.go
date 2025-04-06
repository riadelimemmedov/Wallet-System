package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	log = zap.L()
)

// Client wraps redis client with additional functionality
type Client struct {
	client *redis.Client
}

// Config contains Redis connection configuration
type Config struct {
	Host            string
	Port            string
	Password        string
	DB              int
	PoolSize        int
	MinIdleConns    int
	MaxRetries      int
	MaxConnLifetime time.Duration
	IdleTimeout     time.Duration
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

// NewClient creates a new (Redis) client instance
func NewClient(config Config) (*Client, error) {
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})
	return &Client{client: rdb}, nil
}

// Get retrieves a value from Redis
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set stores a value in Redis with an expiration time
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Delete removes a key from Redis
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// DeleteByPattern removes keys matching a pattern
func (c *Client) DeleteByPattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		log.Error("Failed to get keys", zap.Error(err))
		return fmt.Errorf("failed to get keys when delete by pattern: %w", err)
	}
	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	return nil
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// Close closes the Redis client connection
func (c *Client) Close() error {
	return c.client.Close()
}

// GetClient returns the underlying Redis client
func (c *Client) GetClient() *redis.Client {
	return c.client
}

// FlushDB remove all keys from the current database
func (c *Client) FlushDB(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// Info returns information and statistics about the server
func (c *Client) Info(ctx context.Context) (string, error) {
	return c.client.Info(ctx).Result()
}

// GetConnectionStats returns connection pool statistics
func (c *Client) GetConnectionStats(ctx context.Context) *redis.PoolStats {
	return c.client.PoolStats()
}

// CheckRedisConnection checks if the Redis connection is working
func (c *Client) CheckRedisConnection() error {
	if c == nil || c.client == nil {
		return fmt.Errorf("redis client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.client.Ping(ctx).Err(); err != nil {
		log.Error("Failed to connect to Redis", zap.Error(err))
		return fmt.Errorf("failed to connect to Redis %w", err)
	}
	return nil
}

// UpdateClient modifies the Redis client configuration
func (c *Client) UpdateClient(ctx context.Context, redisConfig Config) error {
	if err := c.Close(); err != nil {
		log.Error("Failed to close Redis client", zap.Error(err))
		return fmt.Errorf("error closing existing Redis client: %w", err)
	}
	newClient, err := NewClient(redisConfig)
	if err != nil {
		log.Error("Failed to update Redis client", zap.Error(err))
		return fmt.Errorf("error updatingh new Redis client %w", err)
	}
	*c = *newClient
	return nil
}
