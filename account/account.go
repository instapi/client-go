package account

import (
	"time"
)

// Account represents an account.
type Account struct {
	ID        uint64     `json:"id,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Count     int        `json:"count,omitempty"`
	Name      string     `json:"name,omitempty"`
	Company   string     `json:"company,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// CreateAccountRequest represents a create account request.
type CreateAccountRequest struct {
	Name    string `json:"name"`
	Company string `json:"company,omitempty"`
	Limit   int    `json:"limit"`
}
