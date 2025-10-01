package interfaces

import (
	"context"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
)

// TradeRepository defines the interface for trade data operations
type TradeRepository interface {
	// Create creates a new trade record
	Create(ctx context.Context, trade *models.Trade) error

	// GetByID retrieves a trade by its ID
	GetByID(ctx context.Context, tradeID string) (*models.Trade, error)

	// GetByOrderID retrieves all trades for a specific order
	GetByOrderID(ctx context.Context, orderID string) ([]*models.Trade, error)

	// Query retrieves trades based on query parameters
	Query(ctx context.Context, query *models.TradeQuery) ([]*models.Trade, error)

	// GetBySymbol retrieves trades for a specific symbol
	GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.Trade, error)

	// GetByAccount retrieves trades for a specific account
	GetByAccount(ctx context.Context, accountID string) ([]*models.Trade, error)
}
