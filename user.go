package instapi

import (
	"context"
	"net/http"
	"strconv"

	"github.com/instapi/client-go/types"
	"github.com/instapi/client-go/user"
)

// CreateUser creates a new user.
func (c *Client) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	var n *user.User
	_, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"users",
		http.StatusCreated,
		u,
		&n,
	)

	return n, err
}

// User gets the current acting user.
func (c *Client) User(ctx context.Context, options ...RequestOption) (*user.User, error) {
	var u *user.User
	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"users/me",
		http.StatusOK,
		nil,
		&u,
		options...,
	)

	return u, err
}

// GetUser gets a user by ID.
func (c *Client) GetUser(ctx context.Context, userID uint64, options ...RequestOption) (*user.User, error) {
	var u *user.User
	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"users/"+strconv.FormatUint(userID, 10),
		http.StatusOK,
		nil,
		&u,
		options...,
	)

	return u, err
}

// UpdateUser updates a user.
func (c *Client) UpdateUser(ctx context.Context, userID uint64, u *user.User) (*user.User, error) {
	var n *user.User
	_, err := c.doRequest(
		ctx,
		http.MethodPut,
		types.JSON,
		c.endpoint+"users/"+strconv.FormatUint(userID, 10),
		0,
		u,
		&n,
	)

	return n, err
}

// DeleteUser deletes a user.
func (c *Client) DeleteUser(ctx context.Context, userID uint64) error {
	_, err := c.doRequest(
		ctx,
		http.MethodDelete,
		types.JSON,
		c.endpoint+"users/"+strconv.FormatUint(userID, 10),
		0,
		nil,
		nil,
	)

	return err
}
