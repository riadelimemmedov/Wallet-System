package repository

import (
	"context"
	"fmt"

	"github.com/riad/banksystemendtoend/pkg/cache"
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
func (r *CacheableRepository) GetCached(ctx context.Context, key string, result interface{}, getter func() (interface{}, error)) error {
	err := r.cacheService.Get(ctx, key, result)
	if err == nil {
		return nil
	}
	if err != cache.ErrCacheMiss {
		zap.L().Error("failed to get item from cache miss", zap.Error(err))
		return err
	}
	item, err := getter()
	if err != nil {
		zap.L().Error("failed to get item from getter", zap.Error(err))
		return err
	}
	go func() {
		if err := r.cacheService.Set(context.Background(), key, item); err != nil {
			zap.L().Error("failed to set item in cache", zap.Error(err))
			fmt.Printf("failed to set item in cache: %v", err)
		}
	}()

	tempCtx := context.Background()
	tempKey := fmt.Sprintf("temp:%s", key)

	if err := r.cacheService.Set(tempCtx, tempKey, item); err != nil {
		zap.L().Error("failed to set item in temp cache", zap.Error(err))
		return err
	}

	err = r.cacheService.Get(tempCtx, tempKey, result)
	r.cacheService.Delete(tempCtx, tempKey)
	return err
}

// InvalidateCache invalidates cache entries by pattern
func (r *CacheableRepository) InvalidateCache(ctx context.Context, pattern string) error {
	return r.cacheService.DeleteByPattern(ctx, pattern)
}
