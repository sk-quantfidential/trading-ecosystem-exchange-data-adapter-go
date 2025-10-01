package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/interfaces"
	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
	"github.com/sirupsen/logrus"
)

type PostgresAccountRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresAccountRepository(db *sql.DB, logger *logrus.Logger) interfaces.AccountRepository {
	return &PostgresAccountRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresAccountRepository) Create(ctx context.Context, account *models.Account) error {
	query := `
		INSERT INTO exchange.accounts (
			account_id, user_id, account_type, status, kyc_status, created_at, updated_at, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		account.AccountID, account.UserID, account.AccountType, account.Status,
		account.KYCStatus, account.CreatedAt, account.UpdatedAt, account.Metadata,
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to create account")
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (r *PostgresAccountRepository) GetByID(ctx context.Context, accountID string) (*models.Account, error) {
	query := `
		SELECT account_id, user_id, account_type, status, kyc_status, created_at, updated_at, metadata
		FROM exchange.accounts
		WHERE account_id = $1
	`

	account := &models.Account{}
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(
		&account.AccountID, &account.UserID, &account.AccountType, &account.Status,
		&account.KYCStatus, &account.CreatedAt, &account.UpdatedAt, &account.Metadata,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found: %s", accountID)
	}
	if err != nil {
		r.logger.WithError(err).Error("Failed to get account by ID")
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

func (r *PostgresAccountRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Account, error) {
	query := `
		SELECT account_id, user_id, account_type, status, kyc_status, created_at, updated_at, metadata
		FROM exchange.accounts
		WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get accounts by user ID")
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}
	defer rows.Close()

	accounts := []*models.Account{}
	for rows.Next() {
		account := &models.Account{}
		if err := rows.Scan(&account.AccountID, &account.UserID, &account.AccountType,
			&account.Status, &account.KYCStatus, &account.CreatedAt, &account.UpdatedAt,
			&account.Metadata); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *PostgresAccountRepository) Query(ctx context.Context, query *models.AccountQuery) ([]*models.Account, error) {
	sqlQuery := `
		SELECT account_id, user_id, account_type, status, kyc_status, created_at, updated_at, metadata
		FROM exchange.accounts
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if query.UserID != nil {
		sqlQuery += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *query.UserID)
		argCount++
	}

	if query.AccountType != nil {
		sqlQuery += fmt.Sprintf(" AND account_type = $%d", argCount)
		args = append(args, *query.AccountType)
		argCount++
	}

	if query.Status != nil {
		sqlQuery += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *query.Status)
		argCount++
	}

	if query.KYCStatus != nil {
		sqlQuery += fmt.Sprintf(" AND kyc_status = $%d", argCount)
		args = append(args, *query.KYCStatus)
		argCount++
	}

	if query.CreatedAfter != nil {
		sqlQuery += fmt.Sprintf(" AND created_at > $%d", argCount)
		args = append(args, *query.CreatedAfter)
		argCount++
	}

	// Add sorting
	if query.SortBy != "" {
		sortOrder := "ASC"
		if strings.ToUpper(query.SortOrder) == "DESC" {
			sortOrder = "DESC"
		}
		sqlQuery += fmt.Sprintf(" ORDER BY %s %s", query.SortBy, sortOrder)
	} else {
		sqlQuery += " ORDER BY created_at DESC"
	}

	// Add pagination
	if query.Limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, query.Limit)
		argCount++
	}

	if query.Offset > 0 {
		sqlQuery += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, query.Offset)
		argCount++
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to query accounts")
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	defer rows.Close()

	accounts := []*models.Account{}
	for rows.Next() {
		account := &models.Account{}
		if err := rows.Scan(&account.AccountID, &account.UserID, &account.AccountType,
			&account.Status, &account.KYCStatus, &account.CreatedAt, &account.UpdatedAt,
			&account.Metadata); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *PostgresAccountRepository) Update(ctx context.Context, account *models.Account) error {
	query := `
		UPDATE exchange.accounts
		SET user_id = $1, account_type = $2, status = $3, kyc_status = $4, updated_at = $5, metadata = $6
		WHERE account_id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		account.UserID, account.AccountType, account.Status, account.KYCStatus,
		account.UpdatedAt, account.Metadata, account.AccountID,
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to update account")
		return fmt.Errorf("failed to update account: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("account not found: %s", account.AccountID)
	}

	return nil
}

func (r *PostgresAccountRepository) UpdateStatus(ctx context.Context, accountID string, status models.AccountStatus) error {
	query := `
		UPDATE exchange.accounts
		SET status = $1, updated_at = $2
		WHERE account_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, status, sql.NullTime{}, accountID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to update account status")
		return fmt.Errorf("failed to update account status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("account not found: %s", accountID)
	}

	return nil
}

func (r *PostgresAccountRepository) Delete(ctx context.Context, accountID string) error {
	query := `DELETE FROM exchange.accounts WHERE account_id = $1`

	result, err := r.db.ExecContext(ctx, query, accountID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete account")
		return fmt.Errorf("failed to delete account: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("account not found: %s", accountID)
	}

	return nil
}
