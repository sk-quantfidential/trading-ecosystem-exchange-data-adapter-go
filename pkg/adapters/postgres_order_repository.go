package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/interfaces"
	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type PostgresOrderRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresOrderRepository(db *sql.DB, logger *logrus.Logger) interfaces.OrderRepository {
	return &PostgresOrderRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresOrderRepository) Create(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO exchange.orders (
			order_id, account_id, symbol, order_type, side, quantity, price, filled_quantity,
			average_price, status, time_in_force, created_at, updated_at, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.ExecContext(ctx, query,
		order.OrderID, order.AccountID, order.Symbol, order.OrderType, order.Side,
		order.Quantity, order.Price, order.FilledQuantity, order.AveragePrice,
		order.Status, order.TimeInForce, order.CreatedAt, order.UpdatedAt, order.Metadata,
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to create order")
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (r *PostgresOrderRepository) GetByID(ctx context.Context, orderID string) (*models.Order, error) {
	query := `
		SELECT order_id, account_id, symbol, order_type, side, quantity, price, filled_quantity,
			   average_price, status, time_in_force, created_at, updated_at, filled_at, cancelled_at, metadata
		FROM exchange.orders
		WHERE order_id = $1
	`

	order := &models.Order{}
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.OrderID, &order.AccountID, &order.Symbol, &order.OrderType, &order.Side,
		&order.Quantity, &order.Price, &order.FilledQuantity, &order.AveragePrice,
		&order.Status, &order.TimeInForce, &order.CreatedAt, &order.UpdatedAt,
		&order.FilledAt, &order.CancelledAt, &order.Metadata,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}
	if err != nil {
		r.logger.WithError(err).Error("Failed to get order by ID")
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

func (r *PostgresOrderRepository) Query(ctx context.Context, query *models.OrderQuery) ([]*models.Order, error) {
	sqlQuery := `
		SELECT order_id, account_id, symbol, order_type, side, quantity, price, filled_quantity,
			   average_price, status, time_in_force, created_at, updated_at, filled_at, cancelled_at, metadata
		FROM exchange.orders
		WHERE 1=1
	`

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

	if query.OrderType != nil {
		sqlQuery += fmt.Sprintf(" AND order_type = $%d", argCount)
		args = append(args, *query.OrderType)
		argCount++
	}

	if query.Side != nil {
		sqlQuery += fmt.Sprintf(" AND side = $%d", argCount)
		args = append(args, *query.Side)
		argCount++
	}

	if query.Status != nil {
		sqlQuery += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *query.Status)
		argCount++
	}

	if query.CreatedAfter != nil {
		sqlQuery += fmt.Sprintf(" AND created_at > $%d", argCount)
		args = append(args, *query.CreatedAfter)
		argCount++
	}

	if query.CreatedBefore != nil {
		sqlQuery += fmt.Sprintf(" AND created_at < $%d", argCount)
		args = append(args, *query.CreatedBefore)
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
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to query orders")
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	orders := []*models.Order{}
	for rows.Next() {
		order := &models.Order{}
		if err := rows.Scan(&order.OrderID, &order.AccountID, &order.Symbol, &order.OrderType,
			&order.Side, &order.Quantity, &order.Price, &order.FilledQuantity, &order.AveragePrice,
			&order.Status, &order.TimeInForce, &order.CreatedAt, &order.UpdatedAt,
			&order.FilledAt, &order.CancelledAt, &order.Metadata); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) UpdateStatus(ctx context.Context, orderID string, status models.OrderStatus) error {
	query := `
		UPDATE exchange.orders
		SET status = $1, updated_at = $2
		WHERE order_id = $3
	`

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), orderID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to update order status")
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (r *PostgresOrderRepository) UpdateFilled(ctx context.Context, orderID string, filledQuantity, averagePrice decimal.Decimal) error {
	query := `
		UPDATE exchange.orders
		SET filled_quantity = $1, average_price = $2, updated_at = $3, filled_at = $4
		WHERE order_id = $5
	`

	_, err := r.db.ExecContext(ctx, query, filledQuantity, averagePrice, time.Now(), time.Now(), orderID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to update order filled")
		return fmt.Errorf("failed to update order filled: %w", err)
	}

	return nil
}

func (r *PostgresOrderRepository) Cancel(ctx context.Context, orderID string) error {
	query := `
		UPDATE exchange.orders
		SET status = $1, updated_at = $2, cancelled_at = $3
		WHERE order_id = $4
	`

	_, err := r.db.ExecContext(ctx, query, models.OrderStatusCancelled, time.Now(), time.Now(), orderID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to cancel order")
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

func (r *PostgresOrderRepository) GetPendingByAccount(ctx context.Context, accountID string) ([]*models.Order, error) {
	query := `
		SELECT order_id, account_id, symbol, order_type, side, quantity, price, filled_quantity,
			   average_price, status, time_in_force, created_at, updated_at, filled_at, cancelled_at, metadata
		FROM exchange.orders
		WHERE account_id = $1 AND status IN ($2, $3)
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, models.OrderStatusPending, models.OrderStatusOpen)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get pending orders")
		return nil, fmt.Errorf("failed to get pending orders: %w", err)
	}
	defer rows.Close()

	orders := []*models.Order{}
	for rows.Next() {
		order := &models.Order{}
		if err := rows.Scan(&order.OrderID, &order.AccountID, &order.Symbol, &order.OrderType,
			&order.Side, &order.Quantity, &order.Price, &order.FilledQuantity, &order.AveragePrice,
			&order.Status, &order.TimeInForce, &order.CreatedAt, &order.UpdatedAt,
			&order.FilledAt, &order.CancelledAt, &order.Metadata); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) GetByAccountAndSymbol(ctx context.Context, accountID, symbol string) ([]*models.Order, error) {
	query := `
		SELECT order_id, account_id, symbol, order_type, side, quantity, price, filled_quantity,
			   average_price, status, time_in_force, created_at, updated_at, filled_at, cancelled_at, metadata
		FROM exchange.orders
		WHERE account_id = $1 AND symbol = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, symbol)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get orders by account and symbol")
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	orders := []*models.Order{}
	for rows.Next() {
		order := &models.Order{}
		if err := rows.Scan(&order.OrderID, &order.AccountID, &order.Symbol, &order.OrderType,
			&order.Side, &order.Quantity, &order.Price, &order.FilledQuantity, &order.AveragePrice,
			&order.Status, &order.TimeInForce, &order.CreatedAt, &order.UpdatedAt,
			&order.FilledAt, &order.CancelledAt, &order.Metadata); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}
