package repository

import (
	"context"
	"fmt"

	"github.com/riad/banksystemendtoend/apperrors"
	"github.com/riad/banksystemendtoend/pkg/cache"
	logger "github.com/riad/banksystemendtoend/pkg/log"
	"go.uber.org/zap"
)

// CacheableRepository provides caching functionality for repositories
type CacheableRepository struct {
	cacheService *cache.Service
}

// NewCacheableRepository creates a new cacheable repository
func NewCacheableRepository(cacheService *cache.Service) *CacheableRepository {
	return &CacheableRepository{cacheService: cacheService}
}

// GetCached retrieves a value from the cache by key. If the key is not found, it calls the getter function to get the value and caches it.
func (r *CacheableRepository) GetCached(ctx context.Context, key string, target interface{}, getter func() (interface{}, error)) error {
	//!Cache Side
	err := r.cacheService.Get(ctx, key, target)

	if err != nil && !apperrors.IsRedisError(err) {
		return nil
	}

	//!Db Side
	item, err := getter()
	if err != nil {
		logger.GetLogger().Error("failed to get item from getter", zap.Error(err))
		return err
	}

	go func() {
		if err := r.cacheService.Set(context.Background(), key, item); err != nil {
			logger.GetLogger().Error("failed to set item in cache", zap.Error(err))
		}
	}()

	tempCtx := context.Background()
	tempKey := fmt.Sprintf("temp:%s", key)

	if err := r.cacheService.Set(tempCtx, tempKey, item); err != nil {
		logger.GetLogger().Error("failed to set item in temp cache", zap.Error(err))
		return nil
	}

	err = r.cacheService.Get(tempCtx, tempKey, target)
	r.cacheService.Delete(tempCtx, tempKey)
	return err
}

// InvalidateCache invalidates cache entries by pattern
func (r *CacheableRepository) InvalidateCache(ctx context.Context, pattern string) error {
	return r.cacheService.DeleteByPattern(ctx, pattern)
}
