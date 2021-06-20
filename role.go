package instapi

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/instapi/client-go/types"
)

// AssignRole assigns a account role for the given user.
func (c *Client) AssignRole(ctx context.Context, account, email, role string, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/roles",
		0,
		struct {
			Email string `json:"email"`
			Role  string `json:"role"`
		}{Email: email, Role: role},
		nil,
		options...,
	)

	return err
}

// Subscribe creates a subscription to a schema.
func (c *Client) Subscribe(ctx context.Context, src, dst, schema, role string, expiresAt time.Time, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(src)+"/schemas/"+schema+"/roles",
		0,
		struct {
			Account   string    `json:"account"`
			Role      string    `json:"role"`
			ExpiresAt time.Time `json:"expiresAt"`
		}{
			Account:   dst,
			Role:      role,
			ExpiresAt: expiresAt,
		},
		nil,
		options...,
	)

	return err
}
