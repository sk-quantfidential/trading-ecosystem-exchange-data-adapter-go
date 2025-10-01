package models

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// Balance represents account balance for a specific symbol
type Balance struct {
	BalanceID        string          `json:"balance_id" db:"balance_id"`
	AccountID        string          `json:"account_id" db:"account_id"`
	Symbol           string          `json:"symbol" db:"symbol"`
	AvailableBalance decimal.Decimal `json:"available_balance" db:"available_balance"`
	LockedBalance    decimal.Decimal `json:"locked_balance" db:"locked_balance"`
	TotalBalance     decimal.Decimal `json:"total_balance" db:"total_balance"`
	LastUpdated      time.Time       `json:"last_updated" db:"last_updated"`
	Metadata         json.RawMessage `json:"metadata,omitempty" db:"metadata"`
}

// BalanceQuery defines query parameters for balance lookups
type BalanceQuery struct {
	AccountID    *string
	Symbol       *string
	MinBalance   *decimal.Decimal
	UpdatedAfter *time.Time
	Limit        int
	Offset       int
	SortBy       string
	SortOrder    string
}
