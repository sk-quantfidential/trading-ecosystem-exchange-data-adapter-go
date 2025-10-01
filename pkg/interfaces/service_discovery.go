package interfaces

import (
	"context"
	"time"
)

type ServiceInfo struct {
	ServiceName string
	ServiceID   string
	Address     string
	Port        int
	Version     string
	Metadata    map[string]string
	RegisteredAt time.Time
	LastHeartbeat time.Time
}

type ServiceDiscoveryRepository interface {
	// Register a service instance
	Register(ctx context.Context, info *ServiceInfo) error

	// Deregister a service instance
	Deregister(ctx context.Context, serviceID string) error

	// Update heartbeat for a service
	Heartbeat(ctx context.Context, serviceID string) error

	// Discover service instances by name
	Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error)

	// Get service info by ID
	GetServiceInfo(ctx context.Context, serviceID string) (*ServiceInfo, error)

	// List all registered services
	ListServices(ctx context.Context) ([]*ServiceInfo, error)

	// Health check
	HealthCheck(ctx context.Context) error
}
