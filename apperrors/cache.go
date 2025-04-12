package apperrors

import "errors"

var (
	ErrCacheMiss                  = errors.New("cache miss")
	ErrRedisConnectionUnavailable = errors.New("redis connection is not available")
	ErrRedisConnectionAvailable   = errors.New("redis connection is available")
	ErrRedisNilResponse           = errors.New("redis: nil")
	ErrRedisUnmarshalFailed       = errors.New("failed to unmarshal data from Redis")
	ErrRedisKeyFetchFailed        = errors.New("failed to get key from Redis cache")
)
