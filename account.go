package instapi

import (
	"context"
	"net/http"

	"github.com/instapi/client-go/account"
	"github.com/instapi/client-go/types"
	"github.com/instapi/client-go/user"
)

// CreateAccount creates a new account.
func (c *Client) CreateAccount(ctx context.Context, req *account.CreateAccountRequest, options ...RequestOption) (*account.CreateAccountResponse, error) {
	var resp *account.CreateAccountResponse
	_, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"accounts",
		http.StatusCreated,
		req,
		&resp,
		options...,
	)

	return resp, err
}

// GetAccounts gets accounts.
func (c *Client) GetAccounts(ctx context.Context, options ...RequestOption) ([]*account.Account, string, error) {
	var a []*account.Account
	resp, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"accounts",
		http.StatusOK,
		nil,
		&a,
		options...,
	)

	if err != nil {
		return nil, "", err
	}

	next, err := nextLink(resp)

	if err != nil {
		return nil, "", err
	}

	return a, next, err
}

// Account gets the current acting account.
func (c *Client) Account(ctx context.Context, options ...RequestOption) (*account.Account, error) {
	var a *account.Account
	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"accounts/me",
		http.StatusOK,
		nil,
		&a,
		options...,
	)

	return a, err
}

// GetAccount gets the account.
func (c *Client) GetAccount(ctx context.Context, id string, options ...RequestOption) (*account.Account, error) {
	var a *account.Account
	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"accounts/"+id,
		http.StatusOK,
		nil,
		&a,
		options...,
	)

	return a, err
}

// UpdateAccount updates an account.
func (c *Client) UpdateAccount(ctx context.Context, id string, a *account.Account, options ...RequestOption) (*account.Account, error) {
	var n *account.Account
	_, err := c.doRequest(
		ctx,
		http.MethodPut,
		types.JSON,
		c.endpoint+"accounts/"+id,
		0,
		a,
		&n,
		options...,
	)

	return n, err
}

// DeleteAccount deletes an account.
func (c *Client) DeleteAccount(ctx context.Context, id string, options ...RequestOption) error {
	_, err := c.doRequest(
		ctx,
		http.MethodDelete,
		types.JSON,
		c.endpoint+"accounts/"+id,
		0,
		nil,
		nil,
		options...,
	)

	return err
}

// GetAccountUsers gets account users.
func (c *Client) GetAccountUsers(ctx context.Context, id string, options ...RequestOption) ([]*user.User, string, error) {
	var u []*user.User
	resp, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"accounts/"+id+"/users",
		http.StatusOK,
		nil,
		&u,
		options...,
	)

	if err != nil {
		return nil, "", err
	}

	next, err := nextLink(resp)

	if err != nil {
		return nil, "", err
	}

	return u, next, err
}

// CreateUserWithRole creates a user with a role on an account.
func (c *Client) CreateUserWithRole(ctx context.Context, accountID string, u *user.User, role string, options ...RequestOption) (*user.User, string, error) {
	u, err := c.CreateUser(ctx, u)

	if err != nil {
		return nil, "", err
	}

	r, err := c.CreateRole(ctx, accountID, u.Email, role, options...)

	if err != nil {
		return nil, "", err
	}

	return u, r.APIKey, nil
}
