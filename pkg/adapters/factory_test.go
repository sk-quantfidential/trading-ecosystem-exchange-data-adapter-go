package adapters

import (
	"io"
	"testing"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/internal/config"
	"github.com/sirupsen/logrus"
)

// TestDeriveSchemaName tests PostgreSQL schema name derivation
func TestDeriveSchemaName(t *testing.T) {
	tests := []struct {
		name         string
		serviceName  string
		instanceName string
		expected     string
	}{
		// Singleton service tests
		{
			name:         "singleton service: exchange-simulator",
			serviceName:  "exchange-simulator",
			instanceName: "exchange-simulator",
			expected:     "exchange",
		},
		{
			name:         "singleton service: exchange-data-adapter",
			serviceName:  "exchange-data-adapter",
			instanceName: "exchange-data-adapter",
			expected:     "exchange",
		},

		// Multi-instance service tests
		{
			name:         "multi-instance: exchange-OKX",
			serviceName:  "exchange-simulator",
			instanceName: "exchange-OKX",
			expected:     "exchange_okx",
		},
		{
			name:         "multi-instance: exchange-Binance",
			serviceName:  "exchange-simulator",
			instanceName: "exchange-Binance",
			expected:     "exchange_binance",
		},
		{
			name:         "multi-instance: exchange-Kraken",
			serviceName:  "exchange-simulator",
			instanceName: "exchange-Kraken",
			expected:     "exchange_kraken",
		},

		// Edge cases
		{
			name:         "edge case: single word instance",
			serviceName:  "exchange-simulator",
			instanceName: "exchange",
			expected:     "exchange",
		},
		{
			name:         "edge case: three part instance",
			serviceName:  "exchange-simulator",
			instanceName: "exchange-OKX-Primary",
			expected:     "exchange_okx",
		},
		{
			name:         "edge case: uppercase service",
			serviceName:  "EXCHANGE-SIMULATOR",
			instanceName: "EXCHANGE-SIMULATOR",
			expected:     "EXCHANGE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deriveSchemaName(tt.serviceName, tt.instanceName)
			if result != tt.expected {
				t.Errorf("deriveSchemaName(%s, %s) = %s, expected %s",
					tt.serviceName, tt.instanceName, result, tt.expected)
			}
		})
	}
}

// TestDeriveRedisNamespace tests Redis namespace derivation
func TestDeriveRedisNamespace(t *testing.T) {
	tests := []struct {
		name         string
		serviceName  string
		instanceName string
		expected     string
	}{
		// Singleton service tests
		{
			name:         "singleton service: exchange-simulator",
			serviceName:  "exchange-simulator",
			instanceName: "exchange-simulator",
			expected:     "exchange",
		},
		{
			name:         "singleton service: exchange-data-adapter",
			serviceName:  "exchange-data-adapter",
			instanceName: "exchange-data-adapter",
			expected:     "exchange",
		},

		// Multi-instance service tests
		{
			name:         "multi-instance: exchange-OKX",
			serviceName:  "exchange-simulator",
			instanceName: "exchange-OKX",
			expected:     "exchange:OKX",
		},
		{
			name:         "multi-instance: exchange-Binance",
			serviceName:  "exchange-simulator",
			instanceName: "exchange-Binance",
			expected:     "exchange:Binance",
		},
		{
			name:         "multi-instance: exchange-Kraken",
			serviceName:  "exchange-simulator",
			instanceName: "exchange-Kraken",
			expected:     "exchange:Kraken",
		},

		// Edge cases
		{
			name:         "edge case: single word instance",
			serviceName:  "exchange-simulator",
			instanceName: "exchange",
			expected:     "exchange",
		},
		{
			name:         "edge case: three part instance",
			serviceName:  "exchange-simulator",
			instanceName: "exchange-OKX-Primary",
			expected:     "exchange:OKX",
		},
		{
			name:         "edge case: uppercase service",
			serviceName:  "EXCHANGE-SIMULATOR",
			instanceName: "EXCHANGE-SIMULATOR",
			expected:     "EXCHANGE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deriveRedisNamespace(tt.serviceName, tt.instanceName)
			if result != tt.expected {
				t.Errorf("deriveRedisNamespace(%s, %s) = %s, expected %s",
					tt.serviceName, tt.instanceName, result, tt.expected)
			}
		})
	}
}

// TestNewExchangeDataAdapter tests the factory with derivation
func TestNewExchangeDataAdapter(t *testing.T) {
	tests := []struct {
		name              string
		config            *config.Config
		expectedSchema    string
		expectedNamespace string
	}{
		{
			name: "uses derived schema when not provided",
			config: &config.Config{
				ServiceName:         "exchange-simulator",
				ServiceInstanceName: "exchange-OKX",
				// SchemaName empty - should be derived
				RedisNamespace: "exchange:OKX",
			},
			expectedSchema:    "exchange_okx",
			expectedNamespace: "exchange:OKX",
		},
		{
			name: "uses derived namespace when not provided",
			config: &config.Config{
				ServiceName:         "exchange-simulator",
				ServiceInstanceName: "exchange-Binance",
				SchemaName:          "exchange_binance",
				// RedisNamespace empty - should be derived
			},
			expectedSchema:    "exchange_binance",
			expectedNamespace: "exchange:Binance",
		},
		{
			name: "uses provided values when both specified",
			config: &config.Config{
				ServiceName:         "custom-service",
				ServiceInstanceName: "custom-instance",
				SchemaName:          "explicit_schema",
				RedisNamespace:      "explicit:namespace",
			},
			expectedSchema:    "explicit_schema",
			expectedNamespace: "explicit:namespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logrus.New()
			logger.SetOutput(io.Discard) // Suppress log output in tests

			adapter, err := NewExchangeDataAdapter(tt.config, logger)

			if err != nil {
				t.Fatalf("NewExchangeDataAdapter failed: %v", err)
			}

			// Verify schema name was applied
			if tt.config.SchemaName != tt.expectedSchema {
				t.Errorf("Schema = %s, expected %s",
					tt.config.SchemaName, tt.expectedSchema)
			}

			// Verify Redis namespace was applied
			if tt.config.RedisNamespace != tt.expectedNamespace {
				t.Errorf("Namespace = %s, expected %s",
					tt.config.RedisNamespace, tt.expectedNamespace)
			}

			_ = adapter // Suppress unused variable warning
		})
	}
}
