package adapters

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/interfaces"
	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
	"github.com/sirupsen/logrus"
)

type PostgresTradeRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresTradeRepository(db *sql.DB, logger *logrus.Logger) interfaces.TradeRepository {
	return &PostgresTradeRepository{db: db, logger: logger}
}

func (r *PostgresTradeRepository) Create(ctx context.Context, trade *models.Trade) error {
	query := `
		INSERT INTO exchange.trades (
			trade_id, order_id, account_id, symbol, side, quantity, price, fee, fee_currency, executed_at, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query, trade.TradeID, trade.OrderID, trade.AccountID,
		trade.Symbol, trade.Side, trade.Quantity, trade.Price, trade.Fee, trade.FeeCurrency,
		trade.ExecutedAt, trade.Metadata)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create trade")
		return fmt.Errorf("failed to create trade: %w", err)
	}
	return nil
}

func (r *PostgresTradeRepository) GetByID(ctx context.Context, tradeID string) (*models.Trade, error) {
	query := `SELECT trade_id, order_id, account_id, symbol, side, quantity, price, fee, fee_currency, executed_at, metadata
		FROM exchange.trades WHERE trade_id = $1`
	trade := &models.Trade{}
	err := r.db.QueryRowContext(ctx, query, tradeID).Scan(&trade.TradeID, &trade.OrderID,
		&trade.AccountID, &trade.Symbol, &trade.Side, &trade.Quantity, &trade.Price, &trade.Fee,
		&trade.FeeCurrency, &trade.ExecutedAt, &trade.Metadata)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("trade not found: %s", tradeID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}
	return trade, nil
}

func (r *PostgresTradeRepository) GetByOrderID(ctx context.Context, orderID string) ([]*models.Trade, error) {
	query := `SELECT trade_id, order_id, account_id, symbol, side, quantity, price, fee, fee_currency, executed_at, metadata
		FROM exchange.trades WHERE order_id = $1 ORDER BY executed_at DESC`
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}
	defer rows.Close()

	trades := []*models.Trade{}
	for rows.Next() {
		trade := &models.Trade{}
		if err := rows.Scan(&trade.TradeID, &trade.OrderID, &trade.AccountID, &trade.Symbol,
			&trade.Side, &trade.Quantity, &trade.Price, &trade.Fee, &trade.FeeCurrency,
			&trade.ExecutedAt, &trade.Metadata); err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}
	return trades, nil
}

func (r *PostgresTradeRepository) Query(ctx context.Context, query *models.TradeQuery) ([]*models.Trade, error) {
	// Simplified query implementation
	sqlQuery := `SELECT trade_id, order_id, account_id, symbol, side, quantity, price, fee, fee_currency, executed_at, metadata
		FROM exchange.trades WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if query.OrderID != nil {
		sqlQuery += fmt.Sprintf(" AND order_id = $%d", argCount)
		args = append(args, *query.OrderID)
		argCount++
	}
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

	sqlQuery += " ORDER BY executed_at DESC"
	if query.Limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, query.Limit)
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()

	trades := []*models.Trade{}
	for rows.Next() {
		trade := &models.Trade{}
		if err := rows.Scan(&trade.TradeID, &trade.OrderID, &trade.AccountID, &trade.Symbol,
			&trade.Side, &trade.Quantity, &trade.Price, &trade.Fee, &trade.FeeCurrency,
			&trade.ExecutedAt, &trade.Metadata); err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}
	return trades, nil
}

func (r *PostgresTradeRepository) GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.Trade, error) {
	query := `SELECT trade_id, order_id, account_id, symbol, side, quantity, price, fee, fee_currency, executed_at, metadata
		FROM exchange.trades WHERE symbol = $1 ORDER BY executed_at DESC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, symbol, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades by symbol: %w", err)
	}
	defer rows.Close()

	trades := []*models.Trade{}
	for rows.Next() {
		trade := &models.Trade{}
		if err := rows.Scan(&trade.TradeID, &trade.OrderID, &trade.AccountID, &trade.Symbol,
			&trade.Side, &trade.Quantity, &trade.Price, &trade.Fee, &trade.FeeCurrency,
			&trade.ExecutedAt, &trade.Metadata); err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}
	return trades, nil
}

func (r *PostgresTradeRepository) GetByAccount(ctx context.Context, accountID string) ([]*models.Trade, error) {
	query := `SELECT trade_id, order_id, account_id, symbol, side, quantity, price, fee, fee_currency, executed_at, metadata
		FROM exchange.trades WHERE account_id = $1 ORDER BY executed_at DESC`
	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades by account: %w", err)
	}
	defer rows.Close()

	trades := []*models.Trade{}
	for rows.Next() {
		trade := &models.Trade{}
		if err := rows.Scan(&trade.TradeID, &trade.OrderID, &trade.AccountID, &trade.Symbol,
			&trade.Side, &trade.Quantity, &trade.Price, &trade.Fee, &trade.FeeCurrency,
			&trade.ExecutedAt, &trade.Metadata); err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}
	return trades, nil
}
