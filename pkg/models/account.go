package models

import (
	"encoding/json"
	"time"
)

// AccountType represents the type of trading account
type AccountType string

const (
	AccountTypeSpot    AccountType = "SPOT"
	AccountTypeMargin  AccountType = "MARGIN"
	AccountTypeFutures AccountType = "FUTURES"
)

// AccountStatus represents the current status of an account
type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "ACTIVE"
	AccountStatusSuspended AccountStatus = "SUSPENDED"
	AccountStatusClosed    AccountStatus = "CLOSED"
)

// KYCStatus represents the KYC verification status
type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "PENDING"
	KYCStatusApproved KYCStatus = "APPROVED"
	KYCStatusRejected KYCStatus = "REJECTED"
)

// Account represents a user trading account
type Account struct {
	AccountID   string        `json:"account_id" db:"account_id"`
	UserID      string        `json:"user_id" db:"user_id"`
	AccountType AccountType   `json:"account_type" db:"account_type"`
	Status      AccountStatus `json:"status" db:"status"`
	KYCStatus   KYCStatus     `json:"kyc_status" db:"kyc_status"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
	Metadata    json.RawMessage `json:"metadata,omitempty" db:"metadata"`
}

// AccountQuery defines query parameters for account lookups
type AccountQuery struct {
	UserID       *string
	AccountType  *AccountType
	Status       *AccountStatus
	KYCStatus    *KYCStatus
	CreatedAfter *time.Time
	Limit        int
	Offset       int
	SortBy       string
	SortOrder    string
}
