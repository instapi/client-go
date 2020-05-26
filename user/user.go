package user

import "time"

// SignIn represents a sign in request.
type SignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents a user.
type User struct {
	ID        uint64     `json:"id,omitempty"`
	Status    string     `json:"status,omitempty"`
	Role      string     `json:"role"`
	Email     string     `json:"email"`
	Forename  string     `json:"forename"`
	Surname   string     `json:"surname"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
