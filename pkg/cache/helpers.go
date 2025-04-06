package cache

import (
	"context"
	"fmt"
	"time"

	logger "github.com/riad/banksystemendtoend/pkg/log"
	"go.uber.org/zap"
)

// buildKey creates a prefixed key for the cache
func (s *Service) buildKey(key string) string {
	if s.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", s.prefix, key)
}

// checkConnection checks if Redis is available and implements circuit-breaker pattern
func (s *Service) checkConnection(ctx context.Context) bool {
	// If we recently checked and Redis was available, return true quickly
	if s.redisAvailable.Load() {
		// Only recheck every checkInterval
		now := time.Now().Unix()
		lastCheck := s.lastCheckTime.Load()
		if now-lastCheck < int64(s.checkInterval.Seconds()) {
			return true
		}
	}

	// It's time to check the connection
	connectionOK := s.CheckRedisConnection()

	// Update our state
	s.redisAvailable.Store(connectionOK)
	s.lastCheckTime.Store(time.Now().Unix())

	if !connectionOK {
		logger.GetLogger().Warn("Redis connection unavailable",
			zap.String("prefix", s.prefix),
			zap.Time("timestamp", time.Now()))
	} else if !s.redisAvailable.Load() {
		logger.GetLogger().Info("Redis connection recovered",
			zap.String("prefix", s.prefix),
			zap.Time("timestamp", time.Now()))
	}

	return connectionOK
}
