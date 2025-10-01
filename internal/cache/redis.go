package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisClient struct {
	Client *redis.Client
	config *config.Config
	logger *logrus.Logger
}

func NewRedisClient(cfg *config.Config, logger *logrus.Logger) (*RedisClient, error) {
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is required")
	}

	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	opts.PoolSize = cfg.RedisPoolSize
	opts.MinIdleConns = cfg.RedisMinIdleConns
	opts.MaxRetries = cfg.RedisMaxRetries
	opts.DialTimeout = cfg.RedisDialTimeout
	opts.ReadTimeout = cfg.RedisReadTimeout
	opts.WriteTimeout = cfg.RedisWriteTimeout

	client := redis.NewClient(opts)

	return &RedisClient{
		Client: client,
		config: cfg,
		logger: logger,
	}, nil
}

func (r *RedisClient) Connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	r.logger.Info("Redis connection established")
	return nil
}

func (r *RedisClient) Disconnect(ctx context.Context) error {
	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			return fmt.Errorf("failed to close Redis: %w", err)
		}
		r.logger.Info("Redis connection closed")
	}
	return nil
}

func (r *RedisClient) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	return nil
}
