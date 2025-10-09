package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
)

// AccountBehaviorTestSuite tests the behavior of account repository operations
type AccountBehaviorTestSuite struct {
	BehaviorTestSuite
}

// TestAccountBehaviorSuite runs the account behavior test suite
func TestAccountBehaviorSuite(t *testing.T) {
	suite.Run(t, new(AccountBehaviorTestSuite))
}

// TestAccountCRUDOperations tests basic account CRUD operations
func (suite *AccountBehaviorTestSuite) TestAccountCRUDOperations() {
	var accountID = GenerateTestUUID()

	suite.Given("a new account to create", func() {
		// Account defined below
	}).When("creating the account", func() {
		account := suite.CreateTestAccount(accountID, func(a *models.Account) {
			a.Type = "spot"
			a.Status = "active"
			a.Permissions = []string{"trade", "withdraw"}
		})

		err := suite.adapter.CreateAccount(suite.ctx, account)
		suite.Require().NoError(err)
		suite.trackCreatedAccount(accountID)
	}).Then("the account should be retrievable", func() {
		retrieved, err := suite.adapter.GetAccount(suite.ctx, accountID)
		suite.Require().NoError(err)
		suite.Equal(accountID, retrieved.ID)
		suite.Equal("spot", retrieved.Type)
		suite.Contains(retrieved.Permissions, "trade")
	}).And("the account can be updated", func() {
		err := suite.adapter.UpdateAccount(suite.ctx, accountID, func(a *models.Account) {
			a.Status = "suspended"
		})
		suite.Require().NoError(err)

		updated, err := suite.adapter.GetAccount(suite.ctx, accountID)
		suite.Require().NoError(err)
		suite.Equal("suspended", updated.Status)
	})
}

// TestAccountQueryByUser tests querying accounts by user
func (suite *AccountBehaviorTestSuite) TestAccountQueryByUser() {
	var (
		userID     = "test-user-" + GenerateTestID("accounts")
		accountID1 = GenerateTestUUID()
		accountID2 = GenerateTestUUID()
	)

	suite.Given("multiple accounts for a user", func() {
		account1 := suite.CreateTestAccount(accountID1, func(a *models.Account) {
			a.UserID = userID
			a.Type = "spot"
		})
		err := suite.adapter.CreateAccount(suite.ctx, account1)
		suite.Require().NoError(err)
		suite.trackCreatedAccount(accountID1)

		account2 := suite.CreateTestAccount(accountID2, func(a *models.Account) {
			a.UserID = userID
			a.Type = "margin"
		})
		err = suite.adapter.CreateAccount(suite.ctx, account2)
		suite.Require().NoError(err)
		suite.trackCreatedAccount(accountID2)
	}).When("querying accounts by user", func() {
		accounts, err := suite.adapter.GetAccountsByUser(suite.ctx, userID)
		suite.Require().NoError(err)

		suite.Then("all user accounts should be returned", func() {
			suite.GreaterOrEqual(len(accounts), 2)

			types := make(map[string]bool)
			for _, account := range accounts {
				types[account.Type] = true
			}
			suite.True(types["spot"])
			suite.True(types["margin"])
		})
	})
}

// TestAccountStatusTransitions tests account status lifecycle
func (suite *AccountBehaviorTestSuite) TestAccountStatusTransitions() {
	var accountID = GenerateTestUUID()

	suite.Given("an active account", func() {
		account := suite.CreateTestAccount(accountID, func(a *models.Account) {
			a.Status = "active"
		})
		err := suite.adapter.CreateAccount(suite.ctx, account)
		suite.Require().NoError(err)
		suite.trackCreatedAccount(accountID)
	}).When("suspending the account", func() {
		err := suite.adapter.UpdateAccount(suite.ctx, accountID, func(a *models.Account) {
			a.Status = "suspended"
		})
		suite.Require().NoError(err)
	}).Then("the status should be updated", func() {
		account, err := suite.adapter.GetAccount(suite.ctx, accountID)
		suite.Require().NoError(err)
		suite.Equal("suspended", account.Status)
	})
}
