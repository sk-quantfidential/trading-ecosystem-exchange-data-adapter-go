package interfaces

import (
	"context"
	"time"
)

type CacheRepository interface {
	// Set a value with optional TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Get a value
	Get(ctx context.Context, key string) (string, error)

	// Delete a key
	Delete(ctx context.Context, key string) error

	// Check if key exists
	Exists(ctx context.Context, key string) (bool, error)

	// Set expiration on key
	Expire(ctx context.Context, key string, ttl time.Duration) error

	// Get keys matching pattern
	Keys(ctx context.Context, pattern string) ([]string, error)

	// Delete keys matching pattern
	DeletePattern(ctx context.Context, pattern string) error

	// Health check
	HealthCheck(ctx context.Context) error
}
