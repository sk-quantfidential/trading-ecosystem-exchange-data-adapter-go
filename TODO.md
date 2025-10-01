# exchange-data-adapter-go - TSE-0001.4 Data Adapters and Orchestrator Integration

## Milestone: TSE-0001.4 - Data Adapters and Orchestrator Integration
**Status**: 📝 **PENDING** - Ready to Start
**Goal**: Create exchange data adapter following audit-data-adapter-go and custodian-data-adapter-go proven pattern
**Components**: Exchange Data Adapter Go
**Dependencies**: TSE-0001.3a (Core Infrastructure Setup) ✅, audit-data-adapter-go pattern ✅, custodian-data-adapter-go pattern ✅
**Estimated Time**: 8-10 hours following established pattern

## 🎯 BDD Acceptance Criteria
> The exchange data adapter can connect to orchestrator PostgreSQL and Redis services, handle exchange-specific operations (accounts, orders, trades, balances), and pass comprehensive behavior tests with proper environment configuration management.

## 📋 Repository Creation and Setup

### Initial Repository Structure
```
exchange-data-adapter-go/
├── .gitignore
├── .env.example
├── go.mod
├── go.sum
├── README.md
├── TODO.md (this file)
├── cmd/
│   └── example/
│       └── main.go                    # Example usage
├── internal/
│   ├── config/
│   │   └── config.go                  # Environment configuration
│   ├── database/
│   │   └── postgres.go                # PostgreSQL connection
│   └── cache/
│       └── redis.go                   # Redis connection
├── pkg/
│   ├── adapters/
│   │   ├── factory.go                 # DataAdapter factory
│   │   ├── postgres_adapter.go        # PostgreSQL implementation
│   │   └── redis_adapter.go           # Redis implementation
│   ├── interfaces/
│   │   ├── account_repository.go      # Account operations
│   │   ├── order_repository.go        # Order operations
│   │   ├── trade_repository.go        # Trade execution operations
│   │   ├── balance_repository.go      # Balance tracking
│   │   ├── service_discovery.go       # Service discovery (shared)
│   │   └── cache.go                   # Cache operations (shared)
│   └── models/
│       ├── account.go                 # Account model
│       ├── order.go                   # Order model
│       ├── trade.go                   # Trade execution model
│       └── balance.go                 # Balance model
└── tests/
    ├── init_test.go                   # Test initialization with godotenv
    ├── behavior_test_suite.go         # BDD test framework
    ├── account_behavior_test.go       # Account tests
    ├── order_behavior_test.go         # Order tests
    ├── trade_behavior_test.go         # Trade tests
    ├── balance_behavior_test.go       # Balance tests
    ├── service_discovery_behavior_test.go
    ├── cache_behavior_test.go
    ├── integration_behavior_test.go
    └── test_utils.go                  # Test utilities
```

## 📋 Task Checklist

### Task 0: Repository Creation and Foundation
**Goal**: Create repository structure and base configuration
**Estimated Time**: 1 hour

#### Steps:
- [ ] Create repository directory structure
- [ ] Initialize go.mod with dependencies:
  ```go
  module github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go

  go 1.24

  require (
      github.com/lib/pq v1.10.9                    // PostgreSQL driver
      github.com/redis/go-redis/v9 v9.15.0        // Redis client
      github.com/sirupsen/logrus v1.9.3           // Logging
      github.com/joho/godotenv v1.5.1             // Environment loading
      github.com/stretchr/testify v1.8.4          // Testing framework
      github.com/shopspring/decimal v1.3.1        // Decimal precision for prices
      google.golang.org/grpc v1.58.3              // gRPC (for models compatibility)
      google.golang.org/protobuf v1.31.0          // Protobuf (for models)
  )
  ```
- [ ] Create .gitignore (copy from audit-data-adapter-go)
- [ ] Create README.md with overview and usage instructions
- [ ] Create .env.example (see configuration below)
- [ ] Create Makefile with test automation

**Evidence to Check**:
- Repository structure created
- go.mod initialized with correct dependencies
- .env.example ready for configuration
- Makefile with test targets

---

### Task 1: Environment Configuration System
**Goal**: Create production-ready .env configuration following 12-factor app principles
**Estimated Time**: 30 minutes

#### .env.example Template:
```bash
# Exchange Data Adapter Configuration
# Copy this to .env and update with your orchestrator credentials

# Service Identity
SERVICE_NAME=exchange-data-adapter
SERVICE_VERSION=1.0.0
ENVIRONMENT=development

# PostgreSQL Configuration (orchestrator credentials)
POSTGRES_URL=postgres://exchange_adapter:exchange-adapter-db-pass@localhost:5432/trading_ecosystem?sslmode=disable

# PostgreSQL Connection Pool
MAX_CONNECTIONS=25
MAX_IDLE_CONNECTIONS=10
CONNECTION_MAX_LIFETIME=300s
CONNECTION_MAX_IDLE_TIME=60s

# Redis Configuration (orchestrator credentials)
# Production: Use exchange-adapter user
# Testing: Use admin user for full access
REDIS_URL=redis://exchange-adapter:exchange-pass@localhost:6379/0

# Redis Connection Pool
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=2
REDIS_MAX_RETRIES=3
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s

# Cache Configuration
CACHE_TTL=300s                          # 5 minutes default TTL
CACHE_NAMESPACE=exchange                # Redis key prefix

# Service Discovery
SERVICE_DISCOVERY_NAMESPACE=exchange    # Service registry namespace
HEARTBEAT_INTERVAL=30s                  # Service heartbeat frequency
SERVICE_TTL=90s                         # Service registration TTL

# Test Environment (for integration tests)
TEST_POSTGRES_URL=postgres://exchange_adapter:exchange-adapter-db-pass@localhost:5432/trading_ecosystem?sslmode=disable
TEST_REDIS_URL=redis://admin:admin-secure-pass@localhost:6379/0

# Logging
LOG_LEVEL=info                          # debug, info, warn, error
LOG_FORMAT=json                         # json, text

# Performance Testing
PERF_TEST_SIZE=1000                     # Number of items for performance tests
PERF_THROUGHPUT_MIN=100                 # Minimum ops/second
PERF_LATENCY_MAX=100ms                  # Maximum average latency

# CI/CD
SKIP_INTEGRATION_TESTS=false            # Set to true in CI without infrastructure
```

#### Configuration Implementation (internal/config/config.go):

Follow audit-data-adapter-go pattern with:
- Environment variable loading with defaults
- godotenv integration for .env file loading
- Type-safe configuration struct
- Helper functions: `getEnv()`, `getEnvInt()`, `getEnvDuration()`, `getEnvBool()`

**Acceptance Criteria**:
- [ ] .env.example created with orchestrator credentials
- [ ] Configuration loading working with defaults
- [ ] godotenv integration for test environment
- [ ] All configuration values accessible via Config struct
- [ ] .gitignore includes .env for security

---

### Task 2: Database Schema and Models
**Goal**: Define exchange-specific database schema and Go models
**Estimated Time**: 2 hours

#### Database Schema (PostgreSQL)

**Schema**: `exchange` (to be created in orchestrator)

**Tables**:

```sql
-- accounts: User trading accounts
CREATE TABLE exchange.accounts (
    account_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(100) NOT NULL,
    account_type VARCHAR(50) NOT NULL, -- 'SPOT', 'MARGIN', 'FUTURES'
    status VARCHAR(50) NOT NULL, -- 'ACTIVE', 'SUSPENDED', 'CLOSED'
    kyc_status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- 'PENDING', 'APPROVED', 'REJECTED'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB,

    CONSTRAINT unique_user_account_type UNIQUE (user_id, account_type)
);

CREATE INDEX idx_accounts_user ON exchange.accounts(user_id);
CREATE INDEX idx_accounts_status ON exchange.accounts(status);
CREATE INDEX idx_accounts_created ON exchange.accounts(created_at);

-- orders: Trading orders
CREATE TABLE exchange.orders (
    order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(100) UNIQUE,
    account_id UUID NOT NULL REFERENCES exchange.accounts(account_id),
    symbol VARCHAR(50) NOT NULL,
    side VARCHAR(10) NOT NULL, -- 'BUY', 'SELL'
    order_type VARCHAR(50) NOT NULL, -- 'MARKET', 'LIMIT', 'STOP_LOSS', 'STOP_LIMIT'
    quantity DECIMAL(24, 8) NOT NULL,
    filled_quantity DECIMAL(24, 8) NOT NULL DEFAULT 0,
    remaining_quantity DECIMAL(24, 8) NOT NULL,
    price DECIMAL(24, 8), -- NULL for market orders
    stop_price DECIMAL(24, 8), -- For stop orders
    status VARCHAR(50) NOT NULL, -- 'PENDING', 'OPEN', 'PARTIALLY_FILLED', 'FILLED', 'CANCELLED', 'REJECTED'
    time_in_force VARCHAR(50) DEFAULT 'GTC', -- 'GTC', 'IOC', 'FOK', 'DAY'
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    filled_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    metadata JSONB,

    CONSTRAINT positive_quantity CHECK (quantity > 0),
    CONSTRAINT valid_filled_quantity CHECK (filled_quantity >= 0 AND filled_quantity <= quantity),
    CONSTRAINT remaining_equals_unfilled CHECK (remaining_quantity = quantity - filled_quantity)
);

CREATE INDEX idx_orders_account ON exchange.orders(account_id);
CREATE INDEX idx_orders_symbol ON exchange.orders(symbol);
CREATE INDEX idx_orders_status ON exchange.orders(status);
CREATE INDEX idx_orders_submitted ON exchange.orders(submitted_at);
CREATE INDEX idx_orders_external ON exchange.orders(external_id);

-- trades: Executed trades (fills)
CREATE TABLE exchange.trades (
    trade_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES exchange.orders(order_id),
    account_id UUID NOT NULL REFERENCES exchange.accounts(account_id),
    symbol VARCHAR(50) NOT NULL,
    side VARCHAR(10) NOT NULL, -- 'BUY', 'SELL'
    quantity DECIMAL(24, 8) NOT NULL,
    price DECIMAL(24, 8) NOT NULL,
    value DECIMAL(24, 8) NOT NULL, -- quantity * price
    fee DECIMAL(24, 8) NOT NULL DEFAULT 0,
    fee_currency VARCHAR(10),
    executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB,

    CONSTRAINT positive_trade_quantity CHECK (quantity > 0),
    CONSTRAINT positive_price CHECK (price > 0),
    CONSTRAINT positive_value CHECK (value > 0),
    CONSTRAINT non_negative_fee CHECK (fee >= 0)
);

CREATE INDEX idx_trades_order ON exchange.trades(order_id);
CREATE INDEX idx_trades_account ON exchange.trades(account_id);
CREATE INDEX idx_trades_symbol ON exchange.trades(symbol);
CREATE INDEX idx_trades_executed ON exchange.trades(executed_at);

-- balances: Account balances per symbol
CREATE TABLE exchange.balances (
    balance_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES exchange.accounts(account_id),
    symbol VARCHAR(50) NOT NULL,
    available_balance DECIMAL(24, 8) NOT NULL DEFAULT 0,
    locked_balance DECIMAL(24, 8) NOT NULL DEFAULT 0,
    total_balance DECIMAL(24, 8) NOT NULL DEFAULT 0,
    last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB,

    CONSTRAINT positive_available_balance CHECK (available_balance >= 0),
    CONSTRAINT positive_locked_balance CHECK (locked_balance >= 0),
    CONSTRAINT total_equals_sum CHECK (total_balance = available_balance + locked_balance),
    CONSTRAINT unique_account_symbol UNIQUE (account_id, symbol)
);

CREATE INDEX idx_balances_account ON exchange.balances(account_id);
CREATE INDEX idx_balances_symbol ON exchange.balances(symbol);
CREATE INDEX idx_balances_updated ON exchange.balances(last_updated);

-- order_history: Order state changes audit trail
CREATE TABLE exchange.order_history (
    history_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES exchange.orders(order_id),
    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,
    old_filled_quantity DECIMAL(24, 8),
    new_filled_quantity DECIMAL(24, 8),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reason VARCHAR(255),
    metadata JSONB
);

CREATE INDEX idx_order_history_order ON exchange.order_history(order_id);
CREATE INDEX idx_order_history_changed ON exchange.order_history(changed_at);
```

#### Go Models (pkg/models/)

**pkg/models/account.go**:
```go
package models

import (
    "encoding/json"
    "time"
)

type AccountType string

const (
    AccountTypeSpot    AccountType = "SPOT"
    AccountTypeMargin  AccountType = "MARGIN"
    AccountTypeFutures AccountType = "FUTURES"
)

type AccountStatus string

const (
    AccountStatusActive    AccountStatus = "ACTIVE"
    AccountStatusSuspended AccountStatus = "SUSPENDED"
    AccountStatusClosed    AccountStatus = "CLOSED"
)

type KYCStatus string

const (
    KYCStatusPending  KYCStatus = "PENDING"
    KYCStatusApproved KYCStatus = "APPROVED"
    KYCStatusRejected KYCStatus = "REJECTED"
)

type Account struct {
    AccountID   string        `json:"account_id" db:"account_id"`
    UserID      string        `json:"user_id" db:"user_id"`
    AccountType AccountType   `json:"account_type" db:"account_type"`
    Status      AccountStatus `json:"status" db:"status"`
    KYCStatus   KYCStatus     `json:"kyc_status" db:"kyc_status"`
    CreatedAt   time.Time     `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
    Metadata    json.RawMessage `json:"metadata,omitempty" db:"metadata"`
}

type AccountQuery struct {
    UserID      *string
    AccountType *AccountType
    Status      *AccountStatus
    KYCStatus   *KYCStatus
    CreatedAfter *time.Time
    Limit       int
    Offset      int
    SortBy      string
    SortOrder   string
}
```

**pkg/models/order.go**:
```go
package models

import (
    "encoding/json"
    "time"
    "github.com/shopspring/decimal"
)

type OrderSide string

const (
    OrderSideBuy  OrderSide = "BUY"
    OrderSideSell OrderSide = "SELL"
)

type OrderType string

const (
    OrderTypeMarket    OrderType = "MARKET"
    OrderTypeLimit     OrderType = "LIMIT"
    OrderTypeStopLoss  OrderType = "STOP_LOSS"
    OrderTypeStopLimit OrderType = "STOP_LIMIT"
)

type OrderStatus string

const (
    OrderStatusPending         OrderStatus = "PENDING"
    OrderStatusOpen            OrderStatus = "OPEN"
    OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
    OrderStatusFilled          OrderStatus = "FILLED"
    OrderStatusCancelled       OrderStatus = "CANCELLED"
    OrderStatusRejected        OrderStatus = "REJECTED"
)

type TimeInForce string

const (
    TimeInForceGTC TimeInForce = "GTC" // Good Till Cancelled
    TimeInForceIOC TimeInForce = "IOC" // Immediate Or Cancel
    TimeInForceFOK TimeInForce = "FOK" // Fill Or Kill
    TimeInForceDAY TimeInForce = "DAY" // Day order
)

type Order struct {
    OrderID           string          `json:"order_id" db:"order_id"`
    ExternalID        *string         `json:"external_id,omitempty" db:"external_id"`
    AccountID         string          `json:"account_id" db:"account_id"`
    Symbol            string          `json:"symbol" db:"symbol"`
    Side              OrderSide       `json:"side" db:"side"`
    OrderType         OrderType       `json:"order_type" db:"order_type"`
    Quantity          decimal.Decimal `json:"quantity" db:"quantity"`
    FilledQuantity    decimal.Decimal `json:"filled_quantity" db:"filled_quantity"`
    RemainingQuantity decimal.Decimal `json:"remaining_quantity" db:"remaining_quantity"`
    Price             *decimal.Decimal `json:"price,omitempty" db:"price"`
    StopPrice         *decimal.Decimal `json:"stop_price,omitempty" db:"stop_price"`
    Status            OrderStatus     `json:"status" db:"status"`
    TimeInForce       TimeInForce     `json:"time_in_force" db:"time_in_force"`
    SubmittedAt       time.Time       `json:"submitted_at" db:"submitted_at"`
    UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
    FilledAt          *time.Time      `json:"filled_at,omitempty" db:"filled_at"`
    CancelledAt       *time.Time      `json:"cancelled_at,omitempty" db:"cancelled_at"`
    Metadata          json.RawMessage `json:"metadata,omitempty" db:"metadata"`
}

type OrderQuery struct {
    AccountID      *string
    Symbol         *string
    Side           *OrderSide
    OrderType      *OrderType
    Status         *OrderStatus
    SubmittedAfter *time.Time
    Limit          int
    Offset         int
    SortBy         string
    SortOrder      string
}
```

**pkg/models/trade.go**:
```go
package models

import (
    "encoding/json"
    "time"
    "github.com/shopspring/decimal"
)

type Trade struct {
    TradeID     string          `json:"trade_id" db:"trade_id"`
    OrderID     string          `json:"order_id" db:"order_id"`
    AccountID   string          `json:"account_id" db:"account_id"`
    Symbol      string          `json:"symbol" db:"symbol"`
    Side        OrderSide       `json:"side" db:"side"`
    Quantity    decimal.Decimal `json:"quantity" db:"quantity"`
    Price       decimal.Decimal `json:"price" db:"price"`
    Value       decimal.Decimal `json:"value" db:"value"`
    Fee         decimal.Decimal `json:"fee" db:"fee"`
    FeeCurrency *string         `json:"fee_currency,omitempty" db:"fee_currency"`
    ExecutedAt  time.Time       `json:"executed_at" db:"executed_at"`
    Metadata    json.RawMessage `json:"metadata,omitempty" db:"metadata"`
}

type TradeQuery struct {
    OrderID       *string
    AccountID     *string
    Symbol        *string
    Side          *OrderSide
    ExecutedAfter *time.Time
    Limit         int
    Offset        int
    SortBy        string
    SortOrder     string
}
```

**pkg/models/balance.go**:
```go
package models

import (
    "encoding/json"
    "time"
    "github.com/shopspring/decimal"
)

type Balance struct {
    BalanceID        string          `json:"balance_id" db:"balance_id"`
    AccountID        string          `json:"account_id" db:"account_id"`
    Symbol           string          `json:"symbol" db:"symbol"`
    AvailableBalance decimal.Decimal `json:"available_balance" db:"available_balance"`
    LockedBalance    decimal.Decimal `json:"locked_balance" db:"locked_balance"`
    TotalBalance     decimal.Decimal `json:"total_balance" db:"total_balance"`
    LastUpdated      time.Time       `json:"last_updated" db:"last_updated"`
    Metadata         json.RawMessage `json:"metadata,omitempty" db:"metadata"`
}

type BalanceQuery struct {
    AccountID    *string
    Symbol       *string
    MinBalance   *decimal.Decimal
    UpdatedAfter *time.Time
    Limit        int
    Offset       int
    SortBy       string
    SortOrder    string
}
```

**Acceptance Criteria**:
- [ ] Database schema defined for exchange domain (5 tables)
- [ ] Go models created with proper JSON tags
- [ ] Query models for flexible filtering
- [ ] Enums for order types, statuses, sides, time in force
- [ ] Proper use of decimal.Decimal for prices and quantities
- [ ] Proper use of json.RawMessage for metadata

---

### Task 3: Repository Interfaces
**Goal**: Define clean interfaces for all exchange operations
**Estimated Time**: 1 hour

#### Account Repository (pkg/interfaces/account_repository.go):
```go
package interfaces

import (
    "context"
    "github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
)

type AccountRepository interface {
    // Create a new account
    Create(ctx context.Context, account *models.Account) error

    // Get account by ID
    GetByID(ctx context.Context, accountID string) (*models.Account, error)

    // Get accounts by user ID
    GetByUserID(ctx context.Context, userID string) ([]*models.Account, error)

    // Get account by user ID and type
    GetByUserAndType(ctx context.Context, userID string, accountType models.AccountType) (*models.Account, error)

    // Query accounts with filters
    Query(ctx context.Context, query *models.AccountQuery) ([]*models.Account, error)

    // Update account status
    UpdateStatus(ctx context.Context, accountID string, status models.AccountStatus) error

    // Update KYC status
    UpdateKYCStatus(ctx context.Context, accountID string, kycStatus models.KYCStatus) error

    // Update account
    Update(ctx context.Context, account *models.Account) error

    // Delete account
    Delete(ctx context.Context, accountID string) error
}
```

#### Order Repository (pkg/interfaces/order_repository.go):
```go
package interfaces

import (
    "context"
    "github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
    "github.com/shopspring/decimal"
)

type OrderRepository interface {
    // Create a new order
    Create(ctx context.Context, order *models.Order) error

    // Get order by ID
    GetByID(ctx context.Context, orderID string) (*models.Order, error)

    // Get order by external ID
    GetByExternalID(ctx context.Context, externalID string) (*models.Order, error)

    // Query orders with filters
    Query(ctx context.Context, query *models.OrderQuery) ([]*models.Order, error)

    // Update order status
    UpdateStatus(ctx context.Context, orderID string, status models.OrderStatus) error

    // Update order fill
    UpdateFill(ctx context.Context, orderID string, filledQty, remainingQty decimal.Decimal) error

    // Complete order (mark as filled)
    Complete(ctx context.Context, orderID string) error

    // Cancel order
    Cancel(ctx context.Context, orderID string, reason string) error

    // Get open orders for account
    GetOpenByAccount(ctx context.Context, accountID string) ([]*models.Order, error)

    // Get orders by account and symbol
    GetByAccountAndSymbol(ctx context.Context, accountID, symbol string) ([]*models.Order, error)
}
```

#### Trade Repository (pkg/interfaces/trade_repository.go):
```go
package interfaces

import (
    "context"
    "github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
)

type TradeRepository interface {
    // Create a new trade
    Create(ctx context.Context, trade *models.Trade) error

    // Get trade by ID
    GetByID(ctx context.Context, tradeID string) (*models.Trade, error)

    // Get trades by order ID
    GetByOrderID(ctx context.Context, orderID string) ([]*models.Trade, error)

    // Query trades with filters
    Query(ctx context.Context, query *models.TradeQuery) ([]*models.Trade, error)

    // Get trades by account
    GetByAccount(ctx context.Context, accountID string) ([]*models.Trade, error)

    // Get trades by account and symbol
    GetByAccountAndSymbol(ctx context.Context, accountID, symbol string) ([]*models.Trade, error)
}
```

#### Balance Repository (pkg/interfaces/balance_repository.go):
```go
package interfaces

import (
    "context"
    "github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
    "github.com/shopspring/decimal"
)

type BalanceRepository interface {
    // Create or update balance
    Upsert(ctx context.Context, balance *models.Balance) error

    // Get balance by ID
    GetByID(ctx context.Context, balanceID string) (*models.Balance, error)

    // Get balance by account and symbol
    GetByAccountAndSymbol(ctx context.Context, accountID, symbol string) (*models.Balance, error)

    // Query balances with filters
    Query(ctx context.Context, query *models.BalanceQuery) ([]*models.Balance, error)

    // Update available balance (for locking/unlocking)
    UpdateAvailableBalance(ctx context.Context, balanceID string, availableBalance, lockedBalance decimal.Decimal) error

    // Get all balances for account
    GetByAccount(ctx context.Context, accountID string) ([]*models.Balance, error)

    // Atomic balance update (for concurrent operations)
    AtomicUpdate(ctx context.Context, accountID, symbol string, availableDelta, lockedDelta decimal.Decimal) error

    // Lock balance for order placement
    LockBalance(ctx context.Context, accountID, symbol string, amount decimal.Decimal) error

    // Unlock balance (order cancelled or filled)
    UnlockBalance(ctx context.Context, accountID, symbol string, amount decimal.Decimal) error
}
```

#### Shared Interfaces (copy from audit-data-adapter-go):

**pkg/interfaces/service_discovery.go** - Same as audit-data-adapter-go
**pkg/interfaces/cache.go** - Same as audit-data-adapter-go

**Acceptance Criteria**:
- [ ] All repository interfaces defined
- [ ] Methods follow CRUD + domain-specific operations pattern
- [ ] Context passed to all methods
- [ ] Proper error handling signatures
- [ ] Query methods use query models for flexibility
- [ ] Decimal types for financial precision

---

### Task 4: PostgreSQL Implementation
**Goal**: Implement repository interfaces using PostgreSQL
**Estimated Time**: 3 hours

Follow audit-data-adapter-go pattern for:
- Connection management (internal/database/postgres.go)
- Repository implementations (pkg/adapters/postgres_*.go)
- Transaction support for atomic operations
- Error handling and logging
- Connection pooling

**Files to Create**:
- `internal/database/postgres.go` - Connection management
- `pkg/adapters/postgres_account_repository.go` - Account operations
- `pkg/adapters/postgres_order_repository.go` - Order operations
- `pkg/adapters/postgres_trade_repository.go` - Trade operations
- `pkg/adapters/postgres_balance_repository.go` - Balance operations

**Key Implementation Notes**:
- Use prepared statements for performance
- Handle decimal.Decimal properly in PostgreSQL queries
- Implement order history tracking on status changes
- Atomic balance updates with row-level locking
- Proper constraint validation

**Acceptance Criteria**:
- [ ] PostgreSQL connection with pooling
- [ ] All repository interfaces implemented
- [ ] Proper error handling and logging
- [ ] Transaction support for atomic operations
- [ ] Decimal precision maintained
- [ ] Health check implementation

---

### Task 5: Redis Implementation
**Goal**: Implement caching and service discovery using Redis
**Estimated Time**: 2 hours

Follow audit-data-adapter-go pattern for:
- Redis connection management (internal/cache/redis.go)
- Service discovery implementation (pkg/adapters/redis_service_discovery.go)
- Cache repository implementation (pkg/adapters/redis_cache_repository.go)

**Acceptance Criteria**:
- [ ] Redis connection with pooling
- [ ] Service discovery working with exchange:* namespace
- [ ] Cache operations with TTL management
- [ ] Health check implementation
- [ ] Graceful fallback when Redis unavailable

---

### Task 6: DataAdapter Factory
**Goal**: Create factory pattern for adapter initialization
**Estimated Time**: 1 hour

#### pkg/adapters/factory.go:
```go
package adapters

import (
    "context"
    "github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/internal/config"
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

func NewExchangeDataAdapter(cfg *config.Config, logger *logrus.Logger) (DataAdapter, error) {
    // Implementation following audit-data-adapter-go pattern
}

func NewExchangeDataAdapterFromEnv(logger *logrus.Logger) (DataAdapter, error) {
    // Load config from environment and create adapter
}
```

**Acceptance Criteria**:
- [ ] Factory pattern implemented
- [ ] Environment-based initialization
- [ ] Proper lifecycle management
- [ ] Health check aggregation

---

### Task 7: BDD Behavior Testing Framework
**Goal**: Create comprehensive test suite following audit-data-adapter-go pattern
**Estimated Time**: 3 hours

#### Test Files to Create:
- `tests/init_test.go` - godotenv loading and test setup
- `tests/behavior_test_suite.go` - BDD framework with Given/When/Then
- `tests/account_behavior_test.go` - Account CRUD and query tests
- `tests/order_behavior_test.go` - Order lifecycle and fill tests
- `tests/trade_behavior_test.go` - Trade creation and query tests
- `tests/balance_behavior_test.go` - Balance operations and atomic updates
- `tests/service_discovery_behavior_test.go` - Service registration tests
- `tests/cache_behavior_test.go` - Cache operations tests
- `tests/integration_behavior_test.go` - Cross-repository consistency tests
- `tests/test_utils.go` - Test utilities and factories

#### Makefile Test Automation:
```makefile
.PHONY: test test-quick test-account test-order test-trade test-balance test-service test-cache test-integration test-all test-coverage check-env

# Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

check-env:
	@if [ ! -f .env ]; then \
		echo "Warning: .env not found. Copy .env.example to .env"; \
		exit 1; \
	fi

test-quick:
	@if [ -f .env ]; then set -a && . ./.env && set +a; fi && \
	go test -v ./tests -run TestAccountBehavior -timeout=2m

test-account: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestAccountBehaviorSuite -timeout=5m

test-order: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestOrderBehaviorSuite -timeout=5m

test-trade: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestTradeBehaviorSuite -timeout=5m

test-balance: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestBalanceBehaviorSuite -timeout=5m

test-service: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestServiceDiscoveryBehaviorSuite -timeout=5m

test-cache: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestCacheBehaviorSuite -timeout=5m

test-integration: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestIntegrationBehaviorSuite -timeout=10m

test-all: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -timeout=15m

test-coverage: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -coverprofile=coverage.out -timeout=15m
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

build:
	go build -v ./...

clean:
	rm -f coverage.out coverage.html
	go clean -testcache
```

**Test Coverage Goals**:
- Account operations: CRUD, status updates, KYC management, queries
- Order operations: Creation, fills, cancellations, status transitions
- Trade operations: Creation, queries, account lookups
- Balance operations: Upsert, atomic updates, locking/unlocking
- Service discovery: Registration, heartbeat, cleanup
- Cache operations: Set/Get, TTL, pattern operations
- Integration: Order placement → balance locking → trade execution → balance update

**Acceptance Criteria**:
- [ ] BDD test framework established
- [ ] 25+ test scenarios covering all repositories
- [ ] Performance tests with configurable thresholds
- [ ] 80%+ average test pass rate
- [ ] CI/CD adaptation (SKIP_INTEGRATION_TESTS)
- [ ] Automatic .env loading in tests

---

### Task 8: Documentation
**Goal**: Create comprehensive documentation for developers
**Estimated Time**: 1 hour

#### README.md:
- Overview of exchange data adapter
- Architecture and repository pattern
- Installation and setup instructions
- Usage examples for all repositories
- Testing guide
- Environment configuration reference

#### tests/README.md:
- Testing framework overview
- How to run different test suites
- Environment setup for tests
- CI/CD integration
- Performance testing configuration

**Acceptance Criteria**:
- [ ] README.md with complete usage guide
- [ ] tests/README.md with testing instructions
- [ ] Code examples for all repositories
- [ ] Environment configuration documented

---

## 📊 Success Metrics

| Metric | Target | Status |
|--------|--------|--------|
| Repository Interfaces | 6 | ⏳ Pending |
| PostgreSQL Tables | 5 | ⏳ Pending |
| Test Scenarios | 25+ | ⏳ Pending |
| Test Pass Rate | 80%+ | ⏳ Pending |
| Code Coverage | 70%+ | ⏳ Pending |
| Build Status | Pass | ⏳ Pending |
| Documentation | Complete | ⏳ Pending |

---

## 🔧 Validation Commands

### Environment Setup
```bash
# Copy environment template
cp .env.example .env

# Edit with orchestrator credentials
vim .env

# Validate environment
make check-env
```

### Testing
```bash
# Quick smoke test
make test-quick

# Individual test suites
make test-account
make test-order
make test-trade
make test-balance
make test-service
make test-cache

# All tests
make test-all

# With coverage
make test-coverage
```

### Build Validation
```bash
# Build all packages
make build

# Run example
go run cmd/example/main.go
```

---

## 🚀 Integration with exchange-simulator-go

Once complete, exchange-simulator-go will integrate by:
1. Adding dependency: `require github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go v0.1.0`
2. Using `replace` directive for local development
3. Initializing adapter in config layer
4. Using repository interfaces in service layer
5. Following audit-correlator-go integration pattern

---

## ✅ Completion Checklist

- [ ] All 8 tasks completed
- [ ] Build passes without errors
- [ ] 25+ test scenarios passing (80%+ success rate)
- [ ] Documentation complete
- [ ] Example code working
- [ ] Ready for exchange-simulator-go integration

---

**Epic**: TSE-0001 Foundation Services & Infrastructure
**Milestone**: TSE-0001.4 Data Adapters & Orchestrator Integration
**Status**: 📝 READY TO START
**Pattern**: Following audit-data-adapter-go and custodian-data-adapter-go proven approach
**Estimated Completion**: 8-10 hours following established pattern

**Last Updated**: 2025-09-30
