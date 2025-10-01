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

type RedisServiceDiscovery struct {
	client    *redis.Client
	namespace string
	logger    *logrus.Logger
}

func NewRedisServiceDiscovery(client *redis.Client, namespace string, logger *logrus.Logger) interfaces.ServiceDiscoveryRepository {
	return &RedisServiceDiscovery{
		client:    client,
		namespace: namespace,
		logger:    logger,
	}
}

func (r *RedisServiceDiscovery) serviceKey(serviceID string) string {
	return fmt.Sprintf("%s:service:%s", r.namespace, serviceID)
}

func (r *RedisServiceDiscovery) heartbeatKey(serviceID string) string {
	return fmt.Sprintf("%s:heartbeat:%s", r.namespace, serviceID)
}

func (r *RedisServiceDiscovery) Register(ctx context.Context, info *interfaces.ServiceInfo) error {
	key := r.serviceKey(info.ServiceID)
	heartbeatKey := r.heartbeatKey(info.ServiceID)

	data, err := json.Marshal(info)
	if err != nil {
		r.logger.WithError(err).Error("Failed to marshal service info")
		return fmt.Errorf("failed to marshal service info: %w", err)
	}

	// Set service info with 90s TTL
	if err := r.client.Set(ctx, key, data, 90*time.Second).Err(); err != nil {
		r.logger.WithError(err).Error("Failed to register service")
		return fmt.Errorf("failed to register service: %w", err)
	}

	// Set initial heartbeat
	if err := r.client.Set(ctx, heartbeatKey, time.Now().Unix(), 90*time.Second).Err(); err != nil {
		r.logger.WithError(err).Error("Failed to set heartbeat")
		return fmt.Errorf("failed to set heartbeat: %w", err)
	}

	r.logger.WithField("service_id", info.ServiceID).Info("Service registered")
	return nil
}

func (r *RedisServiceDiscovery) Deregister(ctx context.Context, serviceID string) error {
	key := r.serviceKey(serviceID)
	heartbeatKey := r.heartbeatKey(serviceID)

	if err := r.client.Del(ctx, key, heartbeatKey).Err(); err != nil {
		r.logger.WithError(err).Error("Failed to deregister service")
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	r.logger.WithField("service_id", serviceID).Info("Service deregistered")
	return nil
}

func (r *RedisServiceDiscovery) Heartbeat(ctx context.Context, serviceID string) error {
	heartbeatKey := r.heartbeatKey(serviceID)
	serviceKey := r.serviceKey(serviceID)

	// Update heartbeat timestamp
	if err := r.client.Set(ctx, heartbeatKey, time.Now().Unix(), 90*time.Second).Err(); err != nil {
		r.logger.WithError(err).Error("Failed to update heartbeat")
		return fmt.Errorf("failed to update heartbeat: %w", err)
	}

	// Refresh service key TTL
	if err := r.client.Expire(ctx, serviceKey, 90*time.Second).Err(); err != nil {
		r.logger.WithError(err).Error("Failed to refresh service TTL")
		return fmt.Errorf("failed to refresh service TTL: %w", err)
	}

	return nil
}

func (r *RedisServiceDiscovery) Discover(ctx context.Context, serviceName string) ([]*interfaces.ServiceInfo, error) {
	pattern := fmt.Sprintf("%s:service:*", r.namespace)

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		r.logger.WithError(err).Error("Failed to discover services")
		return nil, fmt.Errorf("failed to discover services: %w", err)
	}

	services := []*interfaces.ServiceInfo{}
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			r.logger.WithError(err).WithField("key", key).Warn("Failed to get service data")
			continue
		}

		var info interfaces.ServiceInfo
		if err := json.Unmarshal([]byte(data), &info); err != nil {
			r.logger.WithError(err).WithField("key", key).Warn("Failed to unmarshal service info")
			continue
		}

		if info.ServiceName == serviceName {
			services = append(services, &info)
		}
	}

	return services, nil
}

func (r *RedisServiceDiscovery) GetServiceInfo(ctx context.Context, serviceID string) (*interfaces.ServiceInfo, error) {
	key := r.serviceKey(serviceID)

	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("service not found: %s", serviceID)
	}
	if err != nil {
		r.logger.WithError(err).Error("Failed to get service info")
		return nil, fmt.Errorf("failed to get service info: %w", err)
	}

	var info interfaces.ServiceInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		r.logger.WithError(err).Error("Failed to unmarshal service info")
		return nil, fmt.Errorf("failed to unmarshal service info: %w", err)
	}

	return &info, nil
}

func (r *RedisServiceDiscovery) ListServices(ctx context.Context) ([]*interfaces.ServiceInfo, error) {
	pattern := fmt.Sprintf("%s:service:*", r.namespace)

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		r.logger.WithError(err).Error("Failed to list services")
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	services := []*interfaces.ServiceInfo{}
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			r.logger.WithError(err).WithField("key", key).Warn("Failed to get service data")
			continue
		}

		var info interfaces.ServiceInfo
		if err := json.Unmarshal([]byte(data), &info); err != nil {
			r.logger.WithError(err).WithField("key", key).Warn("Failed to unmarshal service info")
			continue
		}

		services = append(services, &info)
	}

	return services, nil
}

func (r *RedisServiceDiscovery) HealthCheck(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("service discovery health check failed: %w", err)
	}
	return nil
}
