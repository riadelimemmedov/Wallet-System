package apperrors

import (
	"errors"
)

// IsRedisConnectionError checks if the error is one of the known Redis connection-related errors
func IsRedisConnectionError(err error) bool {
	return errors.Is(err, ErrRedisConnectionUnavailable)
}

// IsRedisConnectionAvailableError checks if the error indicates Redis connection is available
func IsRedisConnectionAvailableError(err error) bool {
	return errors.Is(err, ErrRedisConnectionAvailable)
}

// IsRedisDataError checks if the error is related to Redis data handling
func IsRedisDataError(err error) bool {
	return errors.Is(err, ErrRedisNilResponse) ||
		errors.Is(err, ErrRedisUnmarshalFailed) ||
		errors.Is(err, ErrRedisKeyFetchFailed) ||
		errors.Is(err, ErrCacheMiss)
}

// IsRedisError checks if the error is any Redis-related error
func IsRedisError(err error) bool {
	return IsRedisConnectionError(err) || IsRedisDataError(err)
}
