package dto

import (
	"database/sql"
	"time"
)

// CreateUserAccountResponse represents the combined response after creating both user and account
type CreateUserAccountResponse struct {
	User    UserResponse    `json:"user"`
	Account AccountResponse `json:"account"`
}

// UserResponse represents the user details in the response
type UserResponse struct {
	UserID          int64          `json:"user_id"`
	Username        string         `json:"username"`
	Email           string         `json:"email"`
	FirstName       sql.NullString `json:"first_name,omitempty"`
	LastName        sql.NullString `json:"last_name,omitempty"`
	PhoneNumber     sql.NullString `json:"phone_number,omitempty"`
	ProfileImageUrl sql.NullString `json:"profile_image_url,omitempty"`
	IsActive        bool           `json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// AccountResponse represents the account details in the response
type AccountResponse struct {
	AccountID      int64     `json:"account_id"`
	UserID         int64     `json:"user_id"`
	AccountNumber  string    `json:"account_number"`
	AccountType    string    `json:"account_type"`
	CurrencyCode   string    `json:"currency_code"`
	Balance        float64   `json:"balance"`
	InterestRate   float64   `json:"interest_rate"`
	OverdraftLimit float64   `json:"overdraft_limit"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AccountTypeResponse represents the response for an account type
type AccountTypeResponse struct {
	ID          int64     `json:"id"`
	AccountType string    `json:"account_type"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateAccountTypeParams struct {
	AccountType string `json:"account_type"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}
