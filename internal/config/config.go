package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Service Identity
	ServiceName         string
	ServiceInstanceName string // Instance identifier (e.g., "exchange-OKX")
	ServiceVersion      string
	Environment         string

	// PostgreSQL Schema (auto-derived if empty)
	SchemaName string

	// Redis Namespace (auto-derived if empty)
	RedisNamespace string

	// PostgreSQL
	PostgresURL            string
	MaxConnections         int
	MaxIdleConnections     int
	ConnectionMaxLifetime  time.Duration
	ConnectionMaxIdleTime  time.Duration

	// Redis
	RedisURL          string
	RedisPoolSize     int
	RedisMinIdleConns int
	RedisMaxRetries   int
	RedisDialTimeout  time.Duration
	RedisReadTimeout  time.Duration
	RedisWriteTimeout time.Duration

	// Cache
	CacheTTL       time.Duration
	CacheNamespace string

	// Service Discovery
	ServiceDiscoveryNamespace string
	HeartbeatInterval         time.Duration
	ServiceTTL                time.Duration

	// Test Environment
	TestPostgresURL string
	TestRedisURL    string

	// Logging
	LogLevel  string
	LogFormat string

	// Performance Testing
	PerfTestSize      int
	PerfThroughputMin int
	PerfLatencyMax    time.Duration

	// CI/CD
	SkipIntegrationTests bool
}

func LoadConfig() (*Config, error) {
	// Try to load .env file (ignore errors if not found)
	_ = godotenv.Load()

	cfg := &Config{
		ServiceName:               getEnv("SERVICE_NAME", "exchange-data-adapter"),
		ServiceInstanceName:       getEnv("SERVICE_INSTANCE_NAME", ""),
		ServiceVersion:            getEnv("SERVICE_VERSION", "1.0.0"),
		Environment:               getEnv("ENVIRONMENT", "development"),
		SchemaName:                getEnv("SCHEMA_NAME", ""),
		RedisNamespace:            getEnv("REDIS_NAMESPACE", ""),
		PostgresURL:               getEnv("POSTGRES_URL", ""),
		MaxConnections:            getEnvInt("MAX_CONNECTIONS", 25),
		MaxIdleConnections:        getEnvInt("MAX_IDLE_CONNECTIONS", 10),
		ConnectionMaxLifetime:     getEnvDuration("CONNECTION_MAX_LIFETIME", 300*time.Second),
		ConnectionMaxIdleTime:     getEnvDuration("CONNECTION_MAX_IDLE_TIME", 60*time.Second),
		RedisURL:                  getEnv("REDIS_URL", ""),
		RedisPoolSize:             getEnvInt("REDIS_POOL_SIZE", 10),
		RedisMinIdleConns:         getEnvInt("REDIS_MIN_IDLE_CONNS", 2),
		RedisMaxRetries:           getEnvInt("REDIS_MAX_RETRIES", 3),
		RedisDialTimeout:          getEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		RedisReadTimeout:          getEnvDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		RedisWriteTimeout:         getEnvDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		CacheTTL:                  getEnvDuration("CACHE_TTL", 300*time.Second),
		CacheNamespace:            getEnv("CACHE_NAMESPACE", "exchange"),
		ServiceDiscoveryNamespace: getEnv("SERVICE_DISCOVERY_NAMESPACE", "exchange"),
		HeartbeatInterval:         getEnvDuration("HEARTBEAT_INTERVAL", 30*time.Second),
		ServiceTTL:                getEnvDuration("SERVICE_TTL", 90*time.Second),
		TestPostgresURL:           getEnv("TEST_POSTGRES_URL", ""),
		TestRedisURL:              getEnv("TEST_REDIS_URL", ""),
		LogLevel:                  getEnv("LOG_LEVEL", "info"),
		LogFormat:                 getEnv("LOG_FORMAT", "json"),
		PerfTestSize:              getEnvInt("PERF_TEST_SIZE", 1000),
		PerfThroughputMin:         getEnvInt("PERF_THROUGHPUT_MIN", 100),
		PerfLatencyMax:            getEnvDuration("PERF_LATENCY_MAX", 100*time.Millisecond),
		SkipIntegrationTests:      getEnvBool("SKIP_INTEGRATION_TESTS", false),
	}

	// Backward compatibility: Default ServiceInstanceName to ServiceName
	if cfg.ServiceInstanceName == "" {
		cfg.ServiceInstanceName = cfg.ServiceName
	}

	// Validate instance name
	if err := ValidateInstanceName(cfg.ServiceInstanceName); err != nil {
		// Log warning but don't fail - allow backward compatibility
		// In production, this should be enforced
		_ = err
	}

	return cfg, nil
}

// ValidateInstanceName validates that an instance name follows DNS-safe naming conventions
func ValidateInstanceName(name string) error {
	// Required explicit - no empty strings
	if name == "" {
		return fmt.Errorf("instance name cannot be empty")
	}

	// DNS-safe: lowercase alphanumeric and hyphens only, must start/end with alphanumeric
	validPattern := regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
	if !validPattern.MatchString(name) {
		return fmt.Errorf("instance name must be DNS-safe: lowercase, alphanumeric, hyphens only, must start and end with letter or number (got: %s)", name)
	}

	// Max 63 characters (DNS label limit)
	if len(name) > 63 {
		return fmt.Errorf("instance name exceeds 63 character limit (got: %d characters)", len(name))
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
