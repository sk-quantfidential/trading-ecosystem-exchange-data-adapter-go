package interfaces

import (
	"context"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
	"github.com/shopspring/decimal"
)

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	// Create creates a new order
	Create(ctx context.Context, order *models.Order) error

	// GetByID retrieves an order by its ID
	GetByID(ctx context.Context, orderID string) (*models.Order, error)

	// Query retrieves orders based on query parameters
	Query(ctx context.Context, query *models.OrderQuery) ([]*models.Order, error)

	// UpdateStatus updates the status of an order
	UpdateStatus(ctx context.Context, orderID string, status models.OrderStatus) error

	// UpdateFilled updates the filled quantity and average price of an order
	UpdateFilled(ctx context.Context, orderID string, filledQuantity, averagePrice decimal.Decimal) error

	// Cancel cancels an order
	Cancel(ctx context.Context, orderID string) error

	// GetPendingByAccount retrieves all pending orders for an account
	GetPendingByAccount(ctx context.Context, accountID string) ([]*models.Order, error)

	// GetByAccountAndSymbol retrieves orders for a specific account and symbol
	GetByAccountAndSymbol(ctx context.Context, accountID, symbol string) ([]*models.Order, error)
}
