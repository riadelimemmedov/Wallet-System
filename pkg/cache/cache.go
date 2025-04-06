package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	logger "github.com/riad/banksystemendtoend/pkg/log"
	"github.com/riad/banksystemendtoend/pkg/redis"
	"go.uber.org/zap"
)

var (
	// ErrCacheMiss is returned when a key is not found in the cache
	ErrCacheMiss = errors.New("cache miss")
)

// DefaultExpiration is the default expiration time for cache entries
const DefaultExpiration = 60 * time.Minute

// Service provides caching functionality using Redis
type Service struct {
	redisClient    *redis.Client
	prefix         string
	defaultTTL     time.Duration
	checkInterval  time.Duration
	redisAvailable atomic.Bool
	lastCheckTime  atomic.Int64
}

// NewService creates a new cache service
func NewService(redisClient *redis.Client, prefix string, defaultTTL time.Duration) *Service {
	if defaultTTL == 0 {
		defaultTTL = DefaultExpiration
	}

	service := &Service{
		redisClient:    redisClient,
		prefix:         prefix,
		defaultTTL:     defaultTTL,
		checkInterval:  1 * time.Second,
		redisAvailable: atomic.Bool{},
		lastCheckTime:  atomic.Int64{},
	}
	return service
}

// Get retrieves a value from the cache and unmarshals it into the provided target
func (s *Service) Get(ctx context.Context, key string, target interface{}) error {
	if !s.checkConnection(ctx) {
		return ErrCacheMiss
	}
	prefixedKey := s.buildKey(key)
	data, err := s.redisClient.Get(ctx, prefixedKey)
	if err != nil {
		if err.Error() == "redis: nil" {
			return ErrCacheMiss
		}
		logger.GetLogger().Error("Failed to get key from cache",
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to get key %s from cache: %w", key, err)
	}

	if err := json.Unmarshal([]byte(data), target); err != nil {
		logger.GetLogger().Error("Failed to unmarshal data",
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to unmarshal data for key %s: %w", key, err)
	}
	return nil
}

// Set stores a value in the cache with the default TTL
func (s *Service) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	if !s.checkConnection(ctx) {
		return ErrCacheMiss
	}

	expiration := s.defaultTTL
	if len(ttl) > 0 {
		expiration = ttl[0]
	}
	return s.SetWithTTL(ctx, key, value, expiration)
}

// SetWithTTL stores a value in the cache with a custom TTL
func (s *Service) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if s.redisClient == nil {
		return nil
	}

	prefixedKey := s.buildKey(key)
	data, err := json.Marshal(value)
	if err != nil {
		logger.GetLogger().Error("Failed to marshal value",
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}
	if err := s.redisClient.Set(ctx, prefixedKey, data, ttl); err != nil {
		logger.GetLogger().Error("Failed to set key in cache",
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to set key %s in cache: %w", key, err)
	}
	return nil
}

// Delete removes a value from the cache
func (s *Service) Delete(ctx context.Context, key string) error {
	prefixedKey := s.buildKey(key)
	err := s.redisClient.Delete(ctx, prefixedKey)
	if err != nil {
		logger.GetLogger().Error("Failed to delete key from cache",
			zap.String("key", key),
			zap.Error(err))
	}
	return err
}

// DeleteByPattern removes values matching a pattern from the cache
func (s *Service) DeleteByPattern(ctx context.Context, pattern string) error {
	prefixedPattern := s.buildKey(pattern + "*")
	err := s.redisClient.DeleteByPattern(ctx, prefixedPattern)
	if err != nil {
		logger.GetLogger().Error("Failed to delete keys by pattern",
			zap.String("pattern", pattern),
			zap.Error(err))
	}
	return err
}

// Exists checks if a key exists in the cache
func (s *Service) Exists(ctx context.Context, key string) (bool, error) {
	if !s.checkConnection(ctx) {
		return false, nil
	}

	prefixedKey := s.buildKey(key)
	exists, err := s.redisClient.Exists(ctx, prefixedKey)
	if err != nil {
		logger.GetLogger().Error("Failed to check key existence",
			zap.String("key", key),
			zap.Error(err))
	}
	return exists, err
}

// GetRedisClient returns the underlying Redis client
func (s *Service) GetRedisClient() *redis.Client {
	return s.redisClient
}

// GetDefaultTTL returns the default TTL for cache entries
func (s *Service) GetDefaultTTL() time.Duration {
	return s.defaultTTL
}

// CheckRedisConnection checks if the Redis connection is working
func (s *Service) CheckRedisConnection() bool {
	if s.redisClient == nil {
		return false
	}
	err := s.redisClient.CheckRedisConnection()
	return err == nil
}
