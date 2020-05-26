package instapi

import (
	"context"
	"net/http"

	"github.com/instapi/client-go/role"
	"github.com/instapi/client-go/types"
)

// CreateRole creates a role for the user with the given email address.
func (c *Client) CreateRole(ctx context.Context, accountID, email, r string, options ...RequestOption) (*role.CreateRoleResponse, error) {
	var n *role.CreateRoleResponse
	_, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"accounts/"+accountID+"/roles",
		0,
		struct {
			Email string `json:"email"`
			Role  string `json:"role"`
		}{Email: email, Role: r},
		&n,
		options...,
	)

	return n, err
}
