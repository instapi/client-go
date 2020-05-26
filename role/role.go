package role

// Available roles.
const (
	System = "SYSTEM"
	Admin  = "ADMIN"
	Create = "CREATE"
	Write  = "WRITE"
	Read   = "READ"
)

// CreateRoleResponse represents a create role response.
type CreateRoleResponse struct {
	APIKey string `json:"apiKey"`
	Role   string `json:"role"`
}
