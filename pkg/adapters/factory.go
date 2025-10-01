package adapters

import (
	"context"
	"fmt"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/internal/cache"
	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/internal/config"
	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/internal/database"
	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/interfaces"
	"github.com/sirupsen/logrus"
)

type DataAdapter interface {
	// Repository access
	AccountRepository() interfaces.AccountRepository
	OrderRepository() interfaces.OrderRepository
	TradeRepository() interfaces.TradeRepository
	BalanceRepository() interfaces.BalanceRepository
	ServiceDiscoveryRepository() interfaces.ServiceDiscoveryRepository
	CacheRepository() interfaces.CacheRepository

	// Lifecycle
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	HealthCheck(ctx context.Context) error
}

type ExchangeDataAdapter struct {
	config *config.Config
	logger *logrus.Logger

	// Infrastructure
	postgresDB  *database.PostgresDB
	redisClient *cache.RedisClient

	// Repositories
	accountRepo          interfaces.AccountRepository
	orderRepo            interfaces.OrderRepository
	tradeRepo            interfaces.TradeRepository
	balanceRepo          interfaces.BalanceRepository
	serviceDiscoveryRepo interfaces.ServiceDiscoveryRepository
	cacheRepo            interfaces.CacheRepository
}

func NewExchangeDataAdapter(cfg *config.Config, logger *logrus.Logger) (DataAdapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	adapter := &ExchangeDataAdapter{
		config: cfg,
		logger: logger,
	}

	// Initialize PostgreSQL
	if cfg.PostgresURL != "" {
		postgresDB, err := database.NewPostgresDB(cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create PostgreSQL client: %w", err)
		}
		adapter.postgresDB = postgresDB

		// Initialize PostgreSQL repositories
		adapter.accountRepo = NewPostgresAccountRepository(postgresDB.DB, logger)
		adapter.orderRepo = NewPostgresOrderRepository(postgresDB.DB, logger)
		adapter.tradeRepo = NewPostgresTradeRepository(postgresDB.DB, logger)
		adapter.balanceRepo = NewPostgresBalanceRepository(postgresDB.DB, logger)
	} else {
		logger.Warn("PostgreSQL URL not configured, repositories will not be available")
	}

	// Initialize Redis
	if cfg.RedisURL != "" {
		redisClient, err := cache.NewRedisClient(cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create Redis client: %w", err)
		}
		adapter.redisClient = redisClient

		// Initialize Redis repositories
		adapter.serviceDiscoveryRepo = NewRedisServiceDiscovery(redisClient.Client, cfg.ServiceDiscoveryNamespace, logger)
		adapter.cacheRepo = NewRedisCacheRepository(redisClient.Client, cfg.CacheNamespace, logger)
	} else {
		logger.Warn("Redis URL not configured, cache and service discovery will not be available")
	}

	return adapter, nil
}

func NewExchangeDataAdapterFromEnv(logger *logrus.Logger) (DataAdapter, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return NewExchangeDataAdapter(cfg, logger)
}

func (a *ExchangeDataAdapter) Connect(ctx context.Context) error {
	// Connect to PostgreSQL
	if a.postgresDB != nil {
		if err := a.postgresDB.Connect(ctx); err != nil {
			a.logger.WithError(err).Warn("Failed to connect to PostgreSQL (stub mode)")
		}
	}

	// Connect to Redis
	if a.redisClient != nil {
		if err := a.redisClient.Connect(ctx); err != nil {
			a.logger.WithError(err).Warn("Failed to connect to Redis (stub mode)")
		}
	}

	a.logger.Info("Exchange data adapter connected")
	return nil
}

func (a *ExchangeDataAdapter) Disconnect(ctx context.Context) error {
	var errors []error

	// Disconnect from PostgreSQL
	if a.postgresDB != nil {
		if err := a.postgresDB.Disconnect(ctx); err != nil {
			errors = append(errors, fmt.Errorf("PostgreSQL disconnect error: %w", err))
		}
	}

	// Disconnect from Redis
	if a.redisClient != nil {
		if err := a.redisClient.Disconnect(ctx); err != nil {
			errors = append(errors, fmt.Errorf("Redis disconnect error: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("disconnect errors: %v", errors)
	}

	a.logger.Info("Exchange data adapter disconnected")
	return nil
}

func (a *ExchangeDataAdapter) HealthCheck(ctx context.Context) error {
	// Check PostgreSQL health
	if a.postgresDB != nil {
		if err := a.postgresDB.HealthCheck(ctx); err != nil {
			return fmt.Errorf("PostgreSQL health check failed: %w", err)
		}
	}

	// Check Redis health
	if a.redisClient != nil {
		if err := a.redisClient.HealthCheck(ctx); err != nil {
			return fmt.Errorf("Redis health check failed: %w", err)
		}
	}

	return nil
}

// Repository accessors
func (a *ExchangeDataAdapter) AccountRepository() interfaces.AccountRepository {
	return a.accountRepo
}

func (a *ExchangeDataAdapter) OrderRepository() interfaces.OrderRepository {
	return a.orderRepo
}

func (a *ExchangeDataAdapter) TradeRepository() interfaces.TradeRepository {
	return a.tradeRepo
}

func (a *ExchangeDataAdapter) BalanceRepository() interfaces.BalanceRepository {
	return a.balanceRepo
}

func (a *ExchangeDataAdapter) ServiceDiscoveryRepository() interfaces.ServiceDiscoveryRepository {
	return a.serviceDiscoveryRepo
}

func (a *ExchangeDataAdapter) CacheRepository() interfaces.CacheRepository {
	return a.cacheRepo
}
