package interfaces

import (
	"context"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
)

// AccountRepository defines the interface for account data operations
type AccountRepository interface {
	// Create a new account
	Create(ctx context.Context, account *models.Account) error

	// GetByID retrieves an account by its ID
	GetByID(ctx context.Context, accountID string) (*models.Account, error)

	// GetByUserID retrieves accounts for a specific user
	GetByUserID(ctx context.Context, userID string) ([]*models.Account, error)

	// Query retrieves accounts based on query parameters
	Query(ctx context.Context, query *models.AccountQuery) ([]*models.Account, error)

	// Update updates an existing account
	Update(ctx context.Context, account *models.Account) error

	// UpdateStatus updates the status of an account
	UpdateStatus(ctx context.Context, accountID string, status models.AccountStatus) error

	// Delete deletes an account by ID
	Delete(ctx context.Context, accountID string) error
}
