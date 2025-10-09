package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
)

// OrderBehaviorTestSuite tests the behavior of order repository operations
type OrderBehaviorTestSuite struct {
	BehaviorTestSuite
}

// TestOrderBehaviorSuite runs the order behavior test suite
func TestOrderBehaviorSuite(t *testing.T) {
	suite.Run(t, new(OrderBehaviorTestSuite))
}

// TestOrderCRUDOperations tests basic order CRUD operations
func (suite *OrderBehaviorTestSuite) TestOrderCRUDOperations() {
	var orderID = GenerateTestUUID()

	suite.Given("a new order to create", func() {
		// Order defined below
	}).When("creating the order", func() {
		order := suite.CreateTestOrder(orderID, func(o *models.Order) {
			o.Symbol = "ETH-USD"
			o.Side = "buy"
			o.Type = "limit"
			o.Quantity = "10.0"
			o.Price = "3000.00"
			o.Status = "pending"
		})

		err := suite.adapter.CreateOrder(suite.ctx, order)
		suite.Require().NoError(err)
		suite.trackCreatedOrder(orderID)
	}).Then("the order should be retrievable", func() {
		retrieved, err := suite.adapter.GetOrder(suite.ctx, orderID)
		suite.Require().NoError(err)
		suite.Equal(orderID, retrieved.ID)
		suite.Equal("ETH-USD", retrieved.Symbol)
		suite.Equal("10.0", retrieved.Quantity)
	}).And("the order can be updated", func() {
		err := suite.adapter.UpdateOrder(suite.ctx, orderID, func(o *models.Order) {
			o.Status = "filled"
			o.FilledQuantity = "10.0"
		})
		suite.Require().NoError(err)

		updated, err := suite.adapter.GetOrder(suite.ctx, orderID)
		suite.Require().NoError(err)
		suite.Equal("filled", updated.Status)
	})
}

// TestOrderQueryByAccount tests querying orders by account
func (suite *OrderBehaviorTestSuite) TestOrderQueryByAccount() {
	var (
		accountID = "test-account-" + GenerateTestID("query")
		orderID1  = GenerateTestUUID()
		orderID2  = GenerateTestUUID()
	)

	suite.Given("multiple orders for an account", func() {
		// Create first order
		order1 := suite.CreateTestOrder(orderID1, func(o *models.Order) {
			o.AccountID = accountID
			o.Symbol = "BTC-USD"
			o.Side = "buy"
		})
		err := suite.adapter.CreateOrder(suite.ctx, order1)
		suite.Require().NoError(err)
		suite.trackCreatedOrder(orderID1)

		// Create second order
		order2 := suite.CreateTestOrder(orderID2, func(o *models.Order) {
			o.AccountID = accountID
			o.Symbol = "ETH-USD"
			o.Side = "sell"
		})
		err = suite.adapter.CreateOrder(suite.ctx, order2)
		suite.Require().NoError(err)
		suite.trackCreatedOrder(orderID2)
	}).When("querying orders by account", func() {
		orders, err := suite.adapter.GetOrdersByAccount(suite.ctx, accountID)
		suite.Require().NoError(err)

		suite.Then("all account orders should be returned", func() {
			suite.GreaterOrEqual(len(orders), 2)

			symbols := make(map[string]bool)
			for _, order := range orders {
				symbols[order.Symbol] = true
			}
			suite.True(symbols["BTC-USD"])
			suite.True(symbols["ETH-USD"])
		})
	})
}

// TestOrderStatusTransitions tests order status lifecycle
func (suite *OrderBehaviorTestSuite) TestOrderStatusTransitions() {
	var orderID = GenerateTestUUID()

	suite.Given("a pending order", func() {
		order := suite.CreateTestOrder(orderID, func(o *models.Order) {
			o.Status = "pending"
		})
		err := suite.adapter.CreateOrder(suite.ctx, order)
		suite.Require().NoError(err)
		suite.trackCreatedOrder(orderID)
	}).When("the order is partially filled", func() {
		err := suite.adapter.UpdateOrder(suite.ctx, orderID, func(o *models.Order) {
			o.Status = "partially_filled"
			o.FilledQuantity = "0.5"
		})
		suite.Require().NoError(err)
	}).Then("the status should be updated", func() {
		order, err := suite.adapter.GetOrder(suite.ctx, orderID)
		suite.Require().NoError(err)
		suite.Equal("partially_filled", order.Status)
		suite.Equal("0.5", order.FilledQuantity)
	}).And("when fully filled", func() {
		err := suite.adapter.UpdateOrder(suite.ctx, orderID, func(o *models.Order) {
			o.Status = "filled"
			o.FilledQuantity = o.Quantity
		})
		suite.Require().NoError(err)

		order, err := suite.adapter.GetOrder(suite.ctx, orderID)
		suite.Require().NoError(err)
		suite.Equal("filled", order.Status)
	})
}

// TestOrderCancellation tests order cancellation
func (suite *OrderBehaviorTestSuite) TestOrderCancellation() {
	var orderID = GenerateTestUUID()

	suite.Given("a pending order", func() {
		order := suite.CreateTestOrder(orderID, func(o *models.Order) {
			o.Status = "pending"
		})
		err := suite.adapter.CreateOrder(suite.ctx, order)
		suite.Require().NoError(err)
		suite.trackCreatedOrder(orderID)
	}).When("cancelling the order", func() {
		err := suite.adapter.UpdateOrder(suite.ctx, orderID, func(o *models.Order) {
			o.Status = "cancelled"
			o.CancelledAt = &time.Time{}
			*o.CancelledAt = time.Now()
		})
		suite.Require().NoError(err)
	}).Then("the order should be cancelled", func() {
		order, err := suite.adapter.GetOrder(suite.ctx, orderID)
		suite.Require().NoError(err)
		suite.Equal("cancelled", order.Status)
		suite.NotNil(order.CancelledAt)
	})
}
