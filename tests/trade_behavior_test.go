package tests

import (
	"testing"
	"strconv"

	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
)

// TradeBehaviorTestSuite tests the behavior of trade repository operations
type TradeBehaviorTestSuite struct {
	BehaviorTestSuite
}

// TestTradeBehaviorSuite runs the trade behavior test suite
func TestTradeBehaviorSuite(t *testing.T) {
	suite.Run(t, new(TradeBehaviorTestSuite))
}

// TestTradeCRUDOperations tests basic trade CRUD operations
func (suite *TradeBehaviorTestSuite) TestTradeCRUDOperations() {
	var tradeID = GenerateTestUUID()

	suite.Given("a new trade to create", func() {
		// Trade defined below
	}).When("creating the trade", func() {
		trade := suite.CreateTestTrade(tradeID, func(t *models.Trade) {
			t.Symbol = "BTC-USD"
			t.Side = "buy"
			t.Quantity = "0.5"
			t.Price = "50000.00"
			t.Fee = "25.00"
		})

		err := suite.adapter.CreateTrade(suite.ctx, trade)
		suite.Require().NoError(err)
		suite.trackCreatedTrade(tradeID)
	}).Then("the trade should be retrievable", func() {
		retrieved, err := suite.adapter.GetTrade(suite.ctx, tradeID)
		suite.Require().NoError(err)
		suite.Equal(tradeID, retrieved.ID)
		suite.Equal("BTC-USD", retrieved.Symbol)
		suite.Equal("0.5", retrieved.Quantity)
	})
}

// TestTradeQueryByOrder tests querying trades by order
func (suite *TradeBehaviorTestSuite) TestTradeQueryByOrder() {
	var (
		orderID  = "test-order-" + GenerateTestID("trades")
		tradeID1 = GenerateTestUUID()
		tradeID2 = GenerateTestUUID()
	)

	suite.Given("multiple trades for an order", func() {
		// Create first trade
		trade1 := suite.CreateTestTrade(tradeID1, func(t *models.Trade) {
			t.OrderID = orderID
			t.Quantity = "0.3"
		})
		err := suite.adapter.CreateTrade(suite.ctx, trade1)
		suite.Require().NoError(err)
		suite.trackCreatedTrade(tradeID1)

		// Create second trade
		trade2 := suite.CreateTestTrade(tradeID2, func(t *models.Trade) {
			t.OrderID = orderID
			t.Quantity = "0.7"
		})
		err = suite.adapter.CreateTrade(suite.ctx, trade2)
		suite.Require().NoError(err)
		suite.trackCreatedTrade(tradeID2)
	}).When("querying trades by order", func() {
		trades, err := suite.adapter.GetTradesByOrder(suite.ctx, orderID)
		suite.Require().NoError(err)

		suite.Then("all order trades should be returned", func() {
			suite.GreaterOrEqual(len(trades), 2)

			totalQuantity := 0.0
			for _, trade := range trades {
				if qty, err := strconv.ParseFloat(trade.Quantity, 64); err == nil {
					totalQuantity += qty
				}
			}
			suite.Equal(1.0, totalQuantity)
		})
	})
}

// TestTradeQueryByAccount tests querying trades by account
func (suite *TradeBehaviorTestSuite) TestTradeQueryByAccount() {
	var (
		accountID = "test-account-" + GenerateTestID("trades")
		tradeID   = GenerateTestUUID()
	)

	suite.Given("a trade for an account", func() {
		trade := suite.CreateTestTrade(tradeID, func(t *models.Trade) {
			t.AccountID = accountID
		})
		err := suite.adapter.CreateTrade(suite.ctx, trade)
		suite.Require().NoError(err)
		suite.trackCreatedTrade(tradeID)
	}).When("querying trades by account", func() {
		trades, err := suite.adapter.GetTradesByAccount(suite.ctx, accountID)
		suite.Require().NoError(err)

		suite.Then("the account trade should be found", func() {
			suite.GreaterOrEqual(len(trades), 1)
			var found bool
			for _, trade := range trades {
				if trade.ID == tradeID {
					found = true
					break
				}
			}
			suite.True(found)
		})
	})
}
