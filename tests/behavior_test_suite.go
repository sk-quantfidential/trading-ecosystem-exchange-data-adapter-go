package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/adapters"
	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
)

// BehaviorTestSuite provides the base test suite for exchange behavior tests
type BehaviorTestSuite struct {
	suite.Suite
	ctx     context.Context
	adapter *adapters.ExchangeDataAdapter
	config  *adapters.Config

	// Tracking for cleanup
	createdOrders    []string
	createdTrades    []string
	createdAccounts  []string
	createdBalances  []string
	createdServices  []string
}

// SetupSuite runs once before all tests
func (suite *BehaviorTestSuite) SetupSuite() {
	suite.T().Log("Setting up behavior test suite")

	// Get configuration from environment
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgres://exchange_adapter:exchange-adapter-db-pass@localhost:5432/trading_ecosystem?sslmode=disable"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://exchange-adapter:exchange-pass@localhost:6379/0"
	}

	// Create configuration
	suite.config = &adapters.Config{
		PostgresURL:         postgresURL,
		RedisURL:            redisURL,
		ServiceName:         "exchange-test",
		ServiceInstanceName: "exchange-test-" + GenerateTestID("instance"),
		Environment:         "testing",
	}

	// Create adapter
	adapter, err := adapters.NewExchangeDataAdapter(suite.config)
	suite.Require().NoError(err, "Failed to create exchange data adapter")
	suite.adapter = adapter

	suite.T().Logf("Exchange Data Adapter initialized with schema: %s, namespace: %s",
		suite.config.SchemaName, suite.config.RedisNamespace)
}

// SetupTest runs before each test
func (suite *BehaviorTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.createdOrders = []string{}
	suite.createdTrades = []string{}
	suite.createdAccounts = []string{}
	suite.createdBalances = []string{}
	suite.createdServices = []string{}
}

// TearDownTest runs after each test
func (suite *BehaviorTestSuite) TearDownTest() {
	suite.T().Log("Cleaning up test data")

	// Clean up created test data
	for _, orderID := range suite.createdOrders {
		if err := suite.adapter.DeleteOrder(suite.ctx, orderID); err != nil {
			suite.T().Logf("Warning: Failed to delete order %s: %v", orderID, err)
		}
	}

	for _, tradeID := range suite.createdTrades {
		if err := suite.adapter.DeleteTrade(suite.ctx, tradeID); err != nil {
			suite.T().Logf("Warning: Failed to delete trade %s: %v", tradeID, err)
		}
	}

	for _, accountID := range suite.createdAccounts {
		if err := suite.adapter.DeleteAccount(suite.ctx, accountID); err != nil {
			suite.T().Logf("Warning: Failed to delete account %s: %v", accountID, err)
		}
	}

	for _, balanceKey := range suite.createdBalances {
		if err := suite.adapter.DeleteBalance(suite.ctx, balanceKey); err != nil {
			suite.T().Logf("Warning: Failed to delete balance %s: %v", balanceKey, err)
		}
	}

	for _, serviceID := range suite.createdServices {
		if err := suite.adapter.DeregisterService(suite.ctx, serviceID); err != nil {
			suite.T().Logf("Warning: Failed to deregister service %s: %v", serviceID, err)
		}
	}
}

// TearDownSuite runs once after all tests
func (suite *BehaviorTestSuite) TearDownSuite() {
	suite.T().Log("Tearing down behavior test suite")

	if suite.adapter != nil {
		err := suite.adapter.Disconnect(suite.ctx)
		if err != nil {
			suite.T().Logf("Warning: Error disconnecting adapter: %v", err)
		}
	}
}

// BDD-style test helpers

// Given starts a BDD scenario with a given condition
func (suite *BehaviorTestSuite) Given(description string, fn func()) *BDDScenario {
	scenario := &BDDScenario{suite: suite}
	scenario.givens = append(scenario.givens, BDDStep{description, fn})
	return scenario
}

// CreateTestOrder creates an order for testing
func (suite *BehaviorTestSuite) CreateTestOrder(orderID string, modifiers ...func(*models.Order)) *models.Order {
	order := &models.Order{
		ID:          orderID,
		AccountID:   "test-account-" + GenerateTestID("account"),
		Symbol:      "BTC-USD",
		Side:        "buy",
		Type:        "limit",
		Quantity:    "1.0",
		Price:       "50000.00",
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, modifier := range modifiers {
		modifier(order)
	}

	return order
}

// CreateTestTrade creates a trade for testing
func (suite *BehaviorTestSuite) CreateTestTrade(tradeID string, modifiers ...func(*models.Trade)) *models.Trade {
	trade := &models.Trade{
		ID:         tradeID,
		OrderID:    "test-order-" + GenerateTestID("order"),
		AccountID:  "test-account-" + GenerateTestID("account"),
		Symbol:     "BTC-USD",
		Side:       "buy",
		Quantity:   "1.0",
		Price:      "50000.00",
		Fee:        "50.00",
		FeeCurrency: "USD",
		ExecutedAt: time.Now(),
		CreatedAt:  time.Now(),
	}

	for _, modifier := range modifiers {
		modifier(trade)
	}

	return trade
}

// CreateTestAccount creates an account for testing
func (suite *BehaviorTestSuite) CreateTestAccount(accountID string, modifiers ...func(*models.Account)) *models.Account {
	account := &models.Account{
		ID:          accountID,
		UserID:      "test-user-" + GenerateTestID("user"),
		Type:        "spot",
		Status:      "active",
		Permissions: []string{"trade", "withdraw"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, modifier := range modifiers {
		modifier(account)
	}

	return account
}

// CreateTestBalance creates a balance for testing
func (suite *BehaviorTestSuite) CreateTestBalance(accountID, symbol string, modifiers ...func(*models.Balance)) *models.Balance {
	balance := &models.Balance{
		AccountID:   accountID,
		Symbol:      symbol,
		Available:   "100000.00",
		Reserved:    "0.00",
		Total:       "100000.00",
		LastUpdated: time.Now(),
	}

	for _, modifier := range modifiers {
		modifier(balance)
	}

	return balance
}

// CreateTestServiceRegistration creates a service registration for testing
func (suite *BehaviorTestSuite) CreateTestServiceRegistration(serviceID string, modifiers ...func(*models.ServiceRegistration)) *models.ServiceRegistration {
	service := &models.ServiceRegistration{
		ID:            serviceID,
		Name:          "test-exchange-service",
		Version:       "1.0.0",
		Host:          "localhost",
		Port:          8080,
		Protocol:      "http",
		Status:        "healthy",
		Metadata:      map[string]string{"environment": "testing"},
		RegisteredAt:  time.Now(),
		LastHeartbeat: time.Now(),
	}

	for _, modifier := range modifiers {
		modifier(service)
	}

	return service
}

// Track created entities for cleanup
func (suite *BehaviorTestSuite) trackCreatedOrder(orderID string) {
	suite.createdOrders = append(suite.createdOrders, orderID)
}

func (suite *BehaviorTestSuite) trackCreatedTrade(tradeID string) {
	suite.createdTrades = append(suite.createdTrades, tradeID)
}

func (suite *BehaviorTestSuite) trackCreatedAccount(accountID string) {
	suite.createdAccounts = append(suite.createdAccounts, accountID)
}

func (suite *BehaviorTestSuite) trackCreatedBalance(key string) {
	suite.createdBalances = append(suite.createdBalances, key)
}

func (suite *BehaviorTestSuite) trackCreatedService(serviceID string) {
	suite.createdServices = append(suite.createdServices, serviceID)
}

// BDDScenario represents a BDD-style test scenario
type BDDScenario struct {
	suite  *BehaviorTestSuite
	givens []BDDStep
	whens  []BDDStep
	thens  []BDDStep
}

// BDDStep represents a single step in a BDD scenario
type BDDStep struct {
	description string
	fn          func()
}

// When adds a when condition
func (s *BDDScenario) When(description string, fn func()) *BDDScenario {
	s.whens = append(s.whens, BDDStep{description, fn})
	return s
}

// Then adds a then assertion
func (s *BDDScenario) Then(description string, fn func()) *BDDScenario {
	s.thens = append(s.thens, BDDStep{description, fn})
	s.execute()
	return s
}

// And adds another condition to the current step type
func (s *BDDScenario) And(description string, fn func()) *BDDScenario {
	if len(s.thens) > 0 {
		s.thens = append(s.thens, BDDStep{description, fn})
	} else if len(s.whens) > 0 {
		s.whens = append(s.whens, BDDStep{description, fn})
	} else {
		s.givens = append(s.givens, BDDStep{description, fn})
	}
	return s
}

// execute runs all steps in the scenario
func (s *BDDScenario) execute() {
	for _, given := range s.givens {
		s.suite.T().Logf("  Given %s", given.description)
		given.fn()
	}

	for _, when := range s.whens {
		s.suite.T().Logf("  When %s", when.description)
		when.fn()
	}

	for _, then := range s.thens {
		s.suite.T().Logf("  Then %s", then.description)
		then.fn()
	}
}

// RunScenario runs a predefined scenario
func (suite *BehaviorTestSuite) RunScenario(scenarioName string) {
	log.Printf("Running scenario: %s", scenarioName)
}
