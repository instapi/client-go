package integration

import "time"

// Integration represents a third party integration, e.g. Google Sheets.
type Integration struct {
	ID   uint64 `json:"id"`
	Type string `json:"type,omitempty"`
}

// Account represents an account integration.
type Account struct {
	Integration
	Scopes    []string   `json:"scopes"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
