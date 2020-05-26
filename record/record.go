package record

import "time"

// Record represents a record.
type Record struct {
	ID        string     `json:"instapi:id"`
	CreatedAt time.Time  `json:"instapi:createdAt"`
	UpdatedAt *time.Time `json:"instapi:updatedAt"`
}
