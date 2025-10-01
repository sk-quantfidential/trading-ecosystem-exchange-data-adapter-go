package interfaces

import (
	"context"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
	"github.com/shopspring/decimal"
)

// BalanceRepository defines the interface for balance data operations
type BalanceRepository interface {
	// Upsert creates or updates a balance record
	Upsert(ctx context.Context, balance *models.Balance) error

	// GetByID retrieves a balance by its ID
	GetByID(ctx context.Context, balanceID string) (*models.Balance, error)

	// GetByAccountAndSymbol retrieves balance for a specific account and symbol
	GetByAccountAndSymbol(ctx context.Context, accountID, symbol string) (*models.Balance, error)

	// Query retrieves balances based on query parameters
	Query(ctx context.Context, query *models.BalanceQuery) ([]*models.Balance, error)

	// UpdateAvailableBalance updates the available and locked balances
	UpdateAvailableBalance(ctx context.Context, balanceID string, availableBalance, lockedBalance decimal.Decimal) error

	// GetByAccount retrieves all balances for a specific account
	GetByAccount(ctx context.Context, accountID string) ([]*models.Balance, error)

	// AtomicUpdate performs an atomic update on balance (for concurrent operations)
	AtomicUpdate(ctx context.Context, accountID, symbol string, availableDelta, lockedDelta decimal.Decimal) error
}
