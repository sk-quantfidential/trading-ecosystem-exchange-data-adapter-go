package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/interfaces"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisCacheRepository struct {
	client    *redis.Client
	namespace string
	logger    *logrus.Logger
}

func NewRedisCacheRepository(client *redis.Client, namespace string, logger *logrus.Logger) interfaces.CacheRepository {
	return &RedisCacheRepository{
		client:    client,
		namespace: namespace,
		logger:    logger,
	}
}

func (r *RedisCacheRepository) keyWithNamespace(key string) string {
	return fmt.Sprintf("%s:%s", r.namespace, key)
}

func (r *RedisCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := r.keyWithNamespace(key)

	var data []byte
	var err error

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			r.logger.WithError(err).Error("Failed to marshal value")
			return fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	if err := r.client.Set(ctx, fullKey, data, ttl).Err(); err != nil {
		r.logger.WithError(err).WithField("key", fullKey).Error("Failed to set cache")
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (r *RedisCacheRepository) Get(ctx context.Context, key string) (string, error) {
	fullKey := r.keyWithNamespace(key)

	result, err := r.client.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		r.logger.WithError(err).WithField("key", fullKey).Error("Failed to get cache")
		return "", fmt.Errorf("failed to get cache: %w", err)
	}

	return result, nil
}

func (r *RedisCacheRepository) Delete(ctx context.Context, key string) error {
	fullKey := r.keyWithNamespace(key)

	if err := r.client.Del(ctx, fullKey).Err(); err != nil {
		r.logger.WithError(err).WithField("key", fullKey).Error("Failed to delete cache")
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	return nil
}

func (r *RedisCacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := r.keyWithNamespace(key)

	count, err := r.client.Exists(ctx, fullKey).Result()
	if err != nil {
		r.logger.WithError(err).WithField("key", fullKey).Error("Failed to check existence")
		return false, fmt.Errorf("failed to check existence: %w", err)
	}

	return count > 0, nil
}

func (r *RedisCacheRepository) Expire(ctx context.Context, key string, ttl time.Duration) error {
	fullKey := r.keyWithNamespace(key)

	if err := r.client.Expire(ctx, fullKey, ttl).Err(); err != nil {
		r.logger.WithError(err).WithField("key", fullKey).Error("Failed to set expiration")
		return fmt.Errorf("failed to set expiration: %w", err)
	}

	return nil
}

func (r *RedisCacheRepository) Keys(ctx context.Context, pattern string) ([]string, error) {
	fullPattern := r.keyWithNamespace(pattern)

	keys, err := r.client.Keys(ctx, fullPattern).Result()
	if err != nil {
		r.logger.WithError(err).WithField("pattern", fullPattern).Error("Failed to get keys")
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	// Remove namespace prefix from keys
	namespacePrefix := r.namespace + ":"
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = key[len(namespacePrefix):]
	}

	return result, nil
}

func (r *RedisCacheRepository) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.Keys(ctx, pattern)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	// Delete all keys matching pattern
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.keyWithNamespace(key)
	}

	if err := r.client.Del(ctx, fullKeys...).Err(); err != nil {
		r.logger.WithError(err).WithField("pattern", pattern).Error("Failed to delete pattern")
		return fmt.Errorf("failed to delete pattern: %w", err)
	}

	return nil
}

func (r *RedisCacheRepository) HealthCheck(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("cache health check failed: %w", err)
	}
	return nil
}
