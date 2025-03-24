package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/riad/banksystemendtoend/pkg/redis"
)

var (
	// ErrCacheMiss is returned when a key is not found in the cache
	ErrCacheMiss = errors.New("cache mis")
)

// DefaultExpiration is the default expiration time for cache entries
const DefaultExpiration = 30 * time.Minute

// Service provides caching functionality using Redis
type Service struct {
	redisClient *redis.Client
	prefix      string
	defaultTTL  time.Duration
}

// NewService creates a new cache service
func NewService(redisClient *redis.Client, prefix string, defaultTTL time.Duration) *Service {
	if defaultTTL == 0 {
		defaultTTL = DefaultExpiration
	}
	return &Service{
		redisClient: redisClient,
		prefix:      prefix,
		defaultTTL:  defaultTTL,
	}
}

// buildKey creates a prefixed key for the cache
func (s *Service) buildKey(key string) string {
	return fmt.Sprintf("%s:%s", s.prefix, key)
}

// Get retrieves a value from the cache and unmarshals it into the provided target
func (s *Service) Get(ctx context.Context, key string, target interface{}) error {
	prefixedKey := s.buildKey(key)

	data, err := s.redisClient.Get(ctx, prefixedKey)
	if err != nil {
		if err.Error() == "redis: nil" {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get key %s from cache: %w", key, err)
	}

	if err := json.Unmarshal([]byte(data), target); err != nil {
		return fmt.Errorf("failed to unmarshal data for key %s: %w", key, err)
	}
	return nil
}

// Set stores a value in the cache with the default TTL
func (s *Service) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	expiration := s.defaultTTL
	if len(ttl) > 0 {
		expiration = ttl[0]
	}
	return s.SetWithTTL(ctx, key, value, expiration)
}

// SetWithTTL stores a value in the cache with a custom TTL
func (s *Service) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	prefixedKey := s.buildKey(key)

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}
	if err := s.redisClient.Set(ctx, prefixedKey, data, ttl); err != nil {
		return fmt.Errorf("failed to set key %s in cache: %w", key, err)
	}
	return nil
}

// Delete removes a value from the cache
func (s *Service) Delete(ctx context.Context, key string) error {
	prefixedKey := s.buildKey(key)
	return s.redisClient.Delete(ctx, prefixedKey)
}

// DeleteByPattern removes values matching a pattern from the cache
func (s *Service) DeleteByPattern(ctx context.Context, pattern string) error {
	prefixedPattern := s.buildKey(pattern + "*")
	return s.redisClient.DeleteByPattern(ctx, prefixedPattern)
}

// Exists checks if a key exists in the cache
func (s *Service) Exists(ctx context.Context, key string) (bool, error) {
	prefixedKey := s.buildKey(key)
	return s.redisClient.Exists(ctx, prefixedKey)
}
