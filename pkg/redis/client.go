package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
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
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	//? Check if the connection is working
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		zap.L().Fatal("Failed to connect to Redis", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return &Client{client: rdb}, nil
}

// Get retrieves a value from Redis
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set stores a value in Redis with an expiration time
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
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
		zap.L().Error("Failed to get keys", zap.Error(err))
		return fmt.Errorf("failed to get keys when delete by pattern: %w", err)
	}
	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	return nil
}

// Exists checks if a key existsz`	`
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
