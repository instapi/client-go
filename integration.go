package instapi

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/instapi/client-go/integration"
	"github.com/instapi/client-go/types"
)

// GetIntegrations gets all third party integrations.
func (c *Client) GetIntegrations(ctx context.Context, options ...RequestOption) ([]*integration.Integration, error) {
	var i []*integration.Integration
	_, _, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"integrations",
		http.StatusOK,
		nil,
		&i,
		options...,
	)

	return i, err
}

// Integration gets a third party integration.
func (c *Client) Integration(ctx context.Context, id uint64, options ...RequestOption) (*integration.Integration, error) {
	var i *integration.Integration
	_, _, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"integrations/"+strconv.FormatUint(id, 10),
		http.StatusOK,
		nil,
		&i,
		options...,
	)

	return i, err
}

// GetAccountIntegrations gets integrations attached to the current account.
func (c *Client) GetAccountIntegrations(ctx context.Context, account string, options ...RequestOption) ([]*integration.Account, error) {
	var i []*integration.Account
	_, _, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/integrations",
		http.StatusOK,
		nil,
		&i,
		options...,
	)

	return i, err
}

// AttachIntegration attaches an integration for use with a schema.
func (c *Client) AttachIntegration(ctx context.Context, account, schema string, id uint64, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(schema)+"/integrations",
		http.StatusNoContent,
		struct {
			ID uint64 `json:"id"`
		}{ID: id},
		nil,
		options...,
	)

	return err
}
