package models

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// OrderType represents the type of order
type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeStop   OrderType = "STOP"
)

// OrderSide represents whether the order is a buy or sell
type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// OrderStatus represents the current status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusOpen      OrderStatus = "OPEN"
	OrderStatusFilled    OrderStatus = "FILLED"
	OrderStatusPartial   OrderStatus = "PARTIALLY_FILLED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusRejected  OrderStatus = "REJECTED"
	OrderStatusExpired   OrderStatus = "EXPIRED"
)

// Order represents a trading order
type Order struct {
	OrderID        string          `json:"order_id" db:"order_id"`
	AccountID      string          `json:"account_id" db:"account_id"`
	Symbol         string          `json:"symbol" db:"symbol"`
	OrderType      OrderType       `json:"order_type" db:"order_type"`
	Side           OrderSide       `json:"side" db:"side"`
	Quantity       decimal.Decimal `json:"quantity" db:"quantity"`
	Price          *decimal.Decimal `json:"price,omitempty" db:"price"` // NULL for market orders
	FilledQuantity decimal.Decimal `json:"filled_quantity" db:"filled_quantity"`
	AveragePrice   *decimal.Decimal `json:"average_price,omitempty" db:"average_price"`
	Status         OrderStatus     `json:"status" db:"status"`
	TimeInForce    string          `json:"time_in_force" db:"time_in_force"` // GTC, IOC, FOK
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
	FilledAt       *time.Time      `json:"filled_at,omitempty" db:"filled_at"`
	CancelledAt    *time.Time      `json:"cancelled_at,omitempty" db:"cancelled_at"`
	Metadata       json.RawMessage `json:"metadata,omitempty" db:"metadata"`
}

// OrderQuery defines query parameters for order lookups
type OrderQuery struct {
	AccountID    *string
	Symbol       *string
	OrderType    *OrderType
	Side         *OrderSide
	Status       *OrderStatus
	CreatedAfter *time.Time
	CreatedBefore *time.Time
	Limit        int
	Offset       int
	SortBy       string
	SortOrder    string
}
