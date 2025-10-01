package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/interfaces"
	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type PostgresBalanceRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresBalanceRepository(db *sql.DB, logger *logrus.Logger) interfaces.BalanceRepository {
	return &PostgresBalanceRepository{db: db, logger: logger}
}

func (r *PostgresBalanceRepository) Upsert(ctx context.Context, balance *models.Balance) error {
	query := `
		INSERT INTO exchange.balances (balance_id, account_id, symbol, available_balance, locked_balance, total_balance, last_updated, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (account_id, symbol) DO UPDATE SET
			available_balance = $4, locked_balance = $5, total_balance = $6, last_updated = $7, metadata = $8
	`
	_, err := r.db.ExecContext(ctx, query, balance.BalanceID, balance.AccountID, balance.Symbol,
		balance.AvailableBalance, balance.LockedBalance, balance.TotalBalance, balance.LastUpdated, balance.Metadata)
	if err != nil {
		r.logger.WithError(err).Error("Failed to upsert balance")
		return fmt.Errorf("failed to upsert balance: %w", err)
	}
	return nil
}

func (r *PostgresBalanceRepository) GetByID(ctx context.Context, balanceID string) (*models.Balance, error) {
	query := `SELECT balance_id, account_id, symbol, available_balance, locked_balance, total_balance, last_updated, metadata
		FROM exchange.balances WHERE balance_id = $1`
	balance := &models.Balance{}
	err := r.db.QueryRowContext(ctx, query, balanceID).Scan(&balance.BalanceID, &balance.AccountID,
		&balance.Symbol, &balance.AvailableBalance, &balance.LockedBalance, &balance.TotalBalance,
		&balance.LastUpdated, &balance.Metadata)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("balance not found: %s", balanceID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

func (r *PostgresBalanceRepository) GetByAccountAndSymbol(ctx context.Context, accountID, symbol string) (*models.Balance, error) {
	query := `SELECT balance_id, account_id, symbol, available_balance, locked_balance, total_balance, last_updated, metadata
		FROM exchange.balances WHERE account_id = $1 AND symbol = $2`
	balance := &models.Balance{}
	err := r.db.QueryRowContext(ctx, query, accountID, symbol).Scan(&balance.BalanceID, &balance.AccountID,
		&balance.Symbol, &balance.AvailableBalance, &balance.LockedBalance, &balance.TotalBalance,
		&balance.LastUpdated, &balance.Metadata)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("balance not found for account %s and symbol %s", accountID, symbol)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

func (r *PostgresBalanceRepository) Query(ctx context.Context, query *models.BalanceQuery) ([]*models.Balance, error) {
	sqlQuery := `SELECT balance_id, account_id, symbol, available_balance, locked_balance, total_balance, last_updated, metadata
		FROM exchange.balances WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if query.AccountID != nil {
		sqlQuery += fmt.Sprintf(" AND account_id = $%d", argCount)
		args = append(args, *query.AccountID)
		argCount++
	}
	if query.Symbol != nil {
		sqlQuery += fmt.Sprintf(" AND symbol = $%d", argCount)
		args = append(args, *query.Symbol)
		argCount++
	}

	sqlQuery += " ORDER BY last_updated DESC"
	if query.Limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, query.Limit)
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query balances: %w", err)
	}
	defer rows.Close()

	balances := []*models.Balance{}
	for rows.Next() {
		balance := &models.Balance{}
		if err := rows.Scan(&balance.BalanceID, &balance.AccountID, &balance.Symbol,
			&balance.AvailableBalance, &balance.LockedBalance, &balance.TotalBalance,
			&balance.LastUpdated, &balance.Metadata); err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}
	return balances, nil
}

func (r *PostgresBalanceRepository) UpdateAvailableBalance(ctx context.Context, balanceID string, availableBalance, lockedBalance decimal.Decimal) error {
	totalBalance := availableBalance.Add(lockedBalance)
	query := `UPDATE exchange.balances SET available_balance = $1, locked_balance = $2, total_balance = $3, last_updated = $4 WHERE balance_id = $5`
	_, err := r.db.ExecContext(ctx, query, availableBalance, lockedBalance, totalBalance, time.Now(), balanceID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to update balance")
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}

func (r *PostgresBalanceRepository) GetByAccount(ctx context.Context, accountID string) ([]*models.Balance, error) {
	query := `SELECT balance_id, account_id, symbol, available_balance, locked_balance, total_balance, last_updated, metadata
		FROM exchange.balances WHERE account_id = $1 ORDER BY symbol`
	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balances by account: %w", err)
	}
	defer rows.Close()

	balances := []*models.Balance{}
	for rows.Next() {
		balance := &models.Balance{}
		if err := rows.Scan(&balance.BalanceID, &balance.AccountID, &balance.Symbol,
			&balance.AvailableBalance, &balance.LockedBalance, &balance.TotalBalance,
			&balance.LastUpdated, &balance.Metadata); err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}
	return balances, nil
}

func (r *PostgresBalanceRepository) AtomicUpdate(ctx context.Context, accountID, symbol string, availableDelta, lockedDelta decimal.Decimal) error {
	query := `
		UPDATE exchange.balances
		SET available_balance = available_balance + $1,
			locked_balance = locked_balance + $2,
			total_balance = total_balance + $1 + $2,
			last_updated = $3
		WHERE account_id = $4 AND symbol = $5
	`
	_, err := r.db.ExecContext(ctx, query, availableDelta, lockedDelta, time.Now(), accountID, symbol)
	if err != nil {
		r.logger.WithError(err).Error("Failed to atomically update balance")
		return fmt.Errorf("failed to atomically update balance: %w", err)
	}
	return nil
}
