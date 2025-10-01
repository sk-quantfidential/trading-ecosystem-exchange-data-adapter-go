# Exchange Data Adapter Behavior Tests

This directory contains comprehensive behavior tests for the exchange-data-adapter-go project. The tests are designed to verify the behavior of the data adapter across different scenarios, including normal operations, error conditions, and performance characteristics.

## Test Structure

### Test Suites

1. **AccountBehaviorTestSuite** - Tests account repository operations
   - CRUD operations for accounts
   - Account type handling (SPOT, MARGIN, FUTURES)
   - Account status lifecycle (ACTIVE, SUSPENDED, CLOSED)
   - KYC status management (PENDING, APPROVED, REJECTED)
   - Account queries and filtering
   - Bulk account operations
   - Performance characteristics

2. **OrderBehaviorTestSuite** - Tests order repository operations
   - CRUD operations for orders
   - Order type handling (MARKET, LIMIT, STOP)
   - Order side management (BUY, SELL)
   - Order status lifecycle (PENDING, OPEN, FILLED, PARTIAL, CANCELLED, REJECTED, EXPIRED)
   - Filled quantity tracking and average price calculation
   - Order cancellation workflow
   - Pending orders by account queries
   - Bulk order operations
   - Performance characteristics

3. **TradeBehaviorTestSuite** - Tests trade repository operations
   - CRUD operations for trades
   - Trade execution records by order
   - Trade history queries by symbol and account
   - Fee tracking and calculation
   - Bulk trade operations
   - Performance characteristics

4. **BalanceBehaviorTestSuite** - Tests balance repository operations
   - CRUD operations for balances
   - Available/locked/total balance consistency
   - Atomic balance updates for concurrent operations
   - Balance queries by account and symbol
   - Multi-currency balance management
   - Bulk balance operations
   - Performance characteristics

5. **ServiceDiscoveryBehaviorTestSuite** - Tests service discovery repository operations
   - Service registration and discovery
   - Heartbeat management
   - Service metrics tracking
   - Multi-instance management
   - Load balancing scenarios
   - Stale service cleanup

6. **CacheBehaviorTestSuite** - Tests cache repository operations
   - Basic key-value operations
   - Complex data type caching
   - TTL (Time-To-Live) management
   - Bulk operations
   - Pattern-based operations
   - Performance testing

7. **IntegrationBehaviorTestSuite** - Tests complete system integration
   - Cross-repository data consistency
   - Full workflow scenarios (order → trade → balance)
   - Transaction consistency
   - Concurrent operations
   - Error recovery
   - Large dataset operations

8. **ComprehensiveBehaviorTestSuite** - Runs all tests in a comprehensive manner
   - Complete system validation
   - Performance benchmarking
   - Scalability testing
   - Error condition handling

### Test Framework Features

- **Behavior-Driven Testing**: Uses Given/When/Then pattern for clear test scenarios
- **Automatic Cleanup**: Tracks created resources and cleans them up automatically
- **Performance Assertions**: Built-in performance testing with configurable thresholds
- **Environment Configuration**: Flexible configuration through environment variables
- **CI/CD Ready**: Automatically adapts behavior for CI environments

## Prerequisites

Before running the tests, ensure you have:

1. **PostgreSQL**: Running instance for exchange data storage
   - Default: `postgres://postgres:postgres@localhost:5432/exchange_test?sslmode=disable`
   - Configure with `TEST_POSTGRES_URL` environment variable

2. **Redis**: Running instance for service discovery and caching
   - Default: `redis://localhost:6379/15` (uses database 15 for tests)
   - Configure with `TEST_REDIS_URL` environment variable

## Running Tests

### Quick Start

```bash
# Run all behavior tests
go test -v ./tests

# Run specific test suite
go test -v ./tests -run TestAccountBehaviorSuite
go test -v ./tests -run TestOrderBehaviorSuite
go test -v ./tests -run TestTradeBehaviorSuite
go test -v ./tests -run TestBalanceBehaviorSuite
go test -v ./tests -run TestServiceDiscoveryBehaviorSuite
go test -v ./tests -run TestCacheBehaviorSuite
go test -v ./tests -run TestIntegrationBehaviorSuite
go test -v ./tests -run TestComprehensiveBehaviorSuite
```

### Environment Configuration

Configure tests using environment variables:

```bash
# Database connections
export TEST_POSTGRES_URL="postgres://user:pass@localhost:5432/exchange_test?sslmode=disable"
export TEST_REDIS_URL="redis://localhost:6379/15"

# Test behavior
export TEST_LOG_LEVEL="info"                # debug, info, warn, error
export TEST_TIMEOUT="10m"                   # Test suite timeout
export SKIP_INTEGRATION_TESTS="false"       # Skip integration tests
export SKIP_PERFORMANCE_TESTS="false"       # Skip performance tests

# Performance testing
export TEST_THROUGHPUT_SIZE="50"            # Number of operations for throughput tests
export TEST_MAX_CONCURRENT_OPS="25"         # Max concurrent operations
export TEST_LARGE_DATASET_SIZE="100"        # Large dataset size for testing

# Run tests
go test -v ./tests
```

### Docker Setup

For easy testing with Docker:

```bash
# Start test databases
docker run -d --name exchange-test-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=exchange_test \
  -p 5432:5432 postgres:17-alpine

docker run -d --name exchange-test-redis \
  -p 6379:6379 redis:8-alpine

# Wait for containers to be ready
sleep 5

# Run tests
go test -v ./tests

# Cleanup
docker rm -f exchange-test-postgres exchange-test-redis
```

### CI/CD Configuration

The tests automatically detect CI environments and adjust behavior:

- Skip performance tests by default in CI
- Use shorter timeouts
- Reduce dataset sizes
- Enable more verbose logging

Example GitHub Actions workflow:

```yaml
name: Behavior Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:17-alpine
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: exchange_test
        options: --health-cmd pg_isready --health-interval 10s --health-timeout 5s --health-retries 5
        ports:
          - 5432:5432

      redis:
        image: redis:8-alpine
        options: --health-cmd "redis-cli ping" --health-interval 10s --health-timeout 5s --health-retries 5
        ports:
          - 6379:6379

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'

    - name: Run Behavior Tests
      env:
        TEST_POSTGRES_URL: postgres://postgres:postgres@localhost:5432/exchange_test?sslmode=disable
        TEST_REDIS_URL: redis://localhost:6379/15
        TEST_LOG_LEVEL: info
      run: go test -v ./tests
```

## Test Scenarios

### Core Scenarios

The tests cover these key scenarios:

1. **Account Lifecycle**
   - Create → Read → Update → Delete
   - Account type and status validation
   - KYC status transitions
   - Account query operations

2. **Order Lifecycle**
   - Create → Read → Update Status → Fill/Cancel
   - Order type and side handling
   - Filled quantity and average price tracking
   - Order status transitions

3. **Trade Lifecycle**
   - Create → Read → Query
   - Trade execution records
   - Fee tracking
   - Trade history queries

4. **Balance Lifecycle**
   - Create → Read → Update → Delete
   - Available/locked balance management
   - Atomic updates for concurrent operations
   - Multi-currency handling

5. **Service Discovery Lifecycle**
   - Register → Discover → Update → Unregister
   - Heartbeat management
   - Health status tracking

6. **Cache Operations**
   - Set → Get → Update → Delete
   - TTL management
   - Pattern operations

7. **Bulk Operations**
   - Batch creation and updates
   - Performance optimization
   - Error handling in batches

8. **Transaction Rollback**
   - Transaction consistency
   - Rollback on errors
   - Data integrity

### Advanced Scenarios

1. **Full Exchange Workflow**
   - Create account
   - Place order (buy/sell)
   - Execute trade
   - Update balances
   - Complete workflow with caching

2. **Concurrent Operations**
   - Multiple simultaneous operations
   - Thread safety validation
   - Consistency under concurrency

3. **Data Consistency**
   - Cross-repository consistency
   - Account-order-trade-balance reconciliation
   - Complex queries

4. **Performance Testing**
   - Throughput measurements
   - Latency validation
   - Scalability testing

## Test Output

### Successful Run Example

```
=== RUN   TestComprehensiveBehaviorSuite
=== Starting Comprehensive Behavior Test Suite ===
=== Behavior Test Environment Information ===
INFO[0000] Environment setting                           key=CI value=false
INFO[0000] Environment setting                           key="Skip Performance" value=false
INFO[0000] Environment setting                           key="Test Timeout" value=5m0s
=== End Environment Information ===
=== RUN   TestComprehensiveBehaviorSuite/TestExchangeDataAdapterBehavior
Testing comprehensive exchange data adapter behavior
=== RUN   TestComprehensiveBehaviorSuite/TestExchangeDataAdapterBehavior/BasicAccounts
Running behavior scenario: account_lifecycle
=== RUN   TestComprehensiveBehaviorSuite/TestExchangeDataAdapterBehavior/BasicOrders
Running behavior scenario: order_lifecycle
=== Comprehensive Behavior Test Suite Completed ===
--- PASS: TestComprehensiveBehaviorSuite (2.45s)
PASS
ok      github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/tests 2.451s
```

### Performance Metrics

The tests provide performance metrics:

```
INFO[0001] Performance measurement                       duration=145.123ms operation="create 100 accounts individually"
INFO[0001] Performance measurement                       duration=89.456ms operation="query 100 orders"
INFO[0002] Performance measurement                       duration=234.789ms operation="bulk create 100 trades"
```

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   ```
   Error: Failed to connect to PostgreSQL
   Solution: Ensure PostgreSQL is running and TEST_POSTGRES_URL is correct
   ```

2. **Redis Connection Errors**
   ```
   Error: Failed to connect to Redis
   Solution: Ensure Redis is running and TEST_REDIS_URL is correct
   ```

3. **Test Timeouts**
   ```
   Error: Test timed out
   Solution: Increase TEST_TIMEOUT or run tests individually
   ```

4. **Performance Test Failures**
   ```
   Error: Operation took longer than expected
   Solution: Run with SKIP_PERFORMANCE_TESTS=true or increase thresholds
   ```

### Debug Mode

Enable debug logging for detailed test execution:

```bash
export TEST_LOG_LEVEL=debug
go test -v ./tests -run TestComprehensiveBehaviorSuite
```

### Individual Test Debugging

Run specific test methods:

```bash
# Run only account tests
go test -v ./tests -run TestAccountBehaviorSuite/TestAccountBasicCRUD

# Run only integration tests
go test -v ./tests -run TestIntegrationBehaviorSuite/TestFullExchangeWorkflow
```

## Contributing

When adding new behavior tests:

1. Follow the Given/When/Then pattern
2. Use the test framework helpers
3. Ensure proper cleanup with tracking methods
4. Add performance assertions where appropriate
5. Update this README with new test scenarios

## Test Coverage

The behavior tests provide comprehensive coverage of:

- ✅ All repository interfaces
- ✅ CRUD operations
- ✅ Complex queries
- ✅ Bulk operations
- ✅ Transaction handling
- ✅ Error conditions
- ✅ Performance characteristics
- ✅ Concurrent operations
- ✅ Data consistency
- ✅ Integration scenarios

For code coverage analysis:

```bash
go test -v -coverprofile=coverage.out ./tests
go tool cover -html=coverage.out -o coverage.html
```

## Epic TSE-0001.4.2: Exchange Data Adapter & Orchestrator Integration

This testing suite is part of epic TSE-0001.4.2, which follows the proven pattern from custodian-data-adapter-go. The goal is to create comprehensive BDD tests for the exchange data adapter.

### Implementation Status

- [ ] Phase 1: Test Infrastructure Setup (init_test.go, base test suite)
- [ ] Phase 2: Account Behavior Tests (account_behavior_test.go)
- [ ] Phase 3: Order Behavior Tests (order_behavior_test.go)
- [ ] Phase 4: Trade Behavior Tests (trade_behavior_test.go)
- [ ] Phase 5: Balance Behavior Tests (balance_behavior_test.go)
- [ ] Phase 6: Service Discovery Tests (service_discovery_behavior_test.go)
- [ ] Phase 7: Cache Behavior Tests (cache_behavior_test.go)
- [ ] Phase 8: Integration Tests (integration_behavior_test.go)
- [ ] Phase 9: Comprehensive Tests (comprehensive_behavior_test.go)
- [ ] Phase 10: Makefile and CI/CD Integration

### Success Criteria

- All test suites passing with >90% success rate
- Test coverage >80% for all repository implementations
- Performance tests validating latency < 100ms for individual operations
- Integration tests validating full exchange workflows
- CI/CD ready with automatic environment detection
