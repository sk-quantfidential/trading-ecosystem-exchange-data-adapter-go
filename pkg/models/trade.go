package models

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// Trade represents an executed trade
type Trade struct {
	TradeID     string          `json:"trade_id" db:"trade_id"`
	OrderID     string          `json:"order_id" db:"order_id"`
	AccountID   string          `json:"account_id" db:"account_id"`
	Symbol      string          `json:"symbol" db:"symbol"`
	Side        OrderSide       `json:"side" db:"side"`
	Quantity    decimal.Decimal `json:"quantity" db:"quantity"`
	Price       decimal.Decimal `json:"price" db:"price"`
	Fee         decimal.Decimal `json:"fee" db:"fee"`
	FeeCurrency string          `json:"fee_currency" db:"fee_currency"`
	ExecutedAt  time.Time       `json:"executed_at" db:"executed_at"`
	Metadata    json.RawMessage `json:"metadata,omitempty" db:"metadata"`
}

// TradeQuery defines query parameters for trade lookups
type TradeQuery struct {
	OrderID       *string
	AccountID     *string
	Symbol        *string
	Side          *OrderSide
	ExecutedAfter *time.Time
	ExecutedBefore *time.Time
	Limit         int
	Offset        int
	SortBy        string
	SortOrder     string
}
