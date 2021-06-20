package user

import "time"

// Credentials represents user credentials.
type Credentials struct {
	Email         string `json:"email"`
	Password      string `json:"password,omitempty"`
	Account       string `json:"account,omitempty"`
	SchemaAccount string `json:"schemaAccount,omitempty"`
	Schema        string `json:"schema,omitempty"`
}

// User represents a user.
type User struct {
	Credentials
	ID        uint64     `json:"id,omitempty"`
	Status    string     `json:"status,omitempty"`
	Role      string     `json:"role"`
	Forename  string     `json:"forename"`
	Surname   string     `json:"surname"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
