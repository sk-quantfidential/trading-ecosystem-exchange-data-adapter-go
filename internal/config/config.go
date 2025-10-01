package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Service Identity
	ServiceName    string
	ServiceVersion string
	Environment    string

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

	return &Config{
		ServiceName:               getEnv("SERVICE_NAME", "exchange-data-adapter"),
		ServiceVersion:            getEnv("SERVICE_VERSION", "1.0.0"),
		Environment:               getEnv("ENVIRONMENT", "development"),
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
	}, nil
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
