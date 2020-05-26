package account

import (
	"time"

	"github.com/instapi/client-go/user"
)

// Account represents an account.
type Account struct {
	Limit     int        `json:"limit,omitempty"`
	Count     int        `json:"count,omitempty"`
	Company   *string    `json:"company,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// CreateAccountRequest represents a create account request.
type CreateAccountRequest struct {
	Name    string  `json:"name"`
	Company *string `json:"company,omitempty"`
	Limit   int     `json:"limit"`
	user.User
}

// CreateAccountResponse represents a create account response.
type CreateAccountResponse struct {
	ID    uint64     `json:"id"`
	Name  string     `json:"name"`
	Limit int        `json:"limit"`
	User  *user.User `json:"user"`
}
