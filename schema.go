package instapi

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/instapi/client-go/schema"
	"github.com/instapi/client-go/types"
)

const detectSizeLimit = 4 * 1024 * 1024

// GetSchema gets the given schema.
func (c *Client) GetSchema(ctx context.Context, account, name string, options ...RequestOption) (*schema.Schema, error) {
	var s *schema.Schema
	_, _, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(name),
		http.StatusOK,
		nil,
		&s,
		options...,
	)

	return s, err
}

// GetSchemas gets the schema collection.
func (c *Client) GetSchemas(ctx context.Context, account string, options ...RequestOption) ([]*schema.Schema, string, error) {
	var s []*schema.Schema
	resp, _, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas",
		http.StatusOK,
		nil,
		&s,
		options...,
	)

	if err != nil {
		return nil, "", err
	}

	next, err := nextLink(resp)

	if err != nil {
		return nil, "", err
	}

	return s, next, nil
}

// ImportSchemasFromFile attempts to import the schemas and records from the given file.
func (c *Client) ImportSchemasFromFile(ctx context.Context, account, filename string, options ...RequestOption) ([]*schema.Import, error) {
	contentType, err := getContentType(filename)

	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename) // nolint: gosec

	if err != nil {
		return nil, err
	}

	defer f.Close() // nolint: gosec

	var s []*schema.Import
	_, _, err = c.doRequest(
		ctx,
		http.MethodPost,
		contentType,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/import",
		http.StatusOK,
		io.LimitReader(f, detectSizeLimit),
		&s,
		options...,
	)

	return s, err
}

// DetectSchemasFromFile attempts to detect the schema for the given file.
func (c *Client) DetectSchemasFromFile(ctx context.Context, name, filename string, options ...RequestOption) ([]*schema.Schema, error) {
	contentType, err := getContentType(filename)

	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename) // nolint: gosec

	if err != nil {
		return nil, err
	}

	defer f.Close() // nolint: gosec

	return c.DetectSchemas(ctx, name, contentType, f, options...)
}

// DetectSchemas attempts to detect the schema for the given reader.
func (c *Client) DetectSchemas(ctx context.Context, name, contentType string, r io.Reader, options ...RequestOption) ([]*schema.Schema, error) {
	var s []*schema.Schema
	_, _, err := c.doRequest(
		ctx,
		http.MethodPost,
		contentType,
		c.endpoint+"detect",
		http.StatusOK,
		io.LimitReader(r, detectSizeLimit),
		&s,
		append(options, Name(name))...,
	)

	return s, err
}

// CreateSchema creates a new schema.
func (c *Client) CreateSchema(ctx context.Context, account string, s *schema.Schema, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas",
		http.StatusCreated,
		s,
		nil,
		options...,
	)

	return err
}

func (c *Client) createSchemas(ctx context.Context, account string, s []*schema.Schema, options ...RequestOption) error {
	var g errgroup.Group

	for _, v := range s {
		v := v

		g.Go(func() error {
			return c.CreateSchema(ctx, account, v, options...)
		})
	}

	return g.Wait()
}

// DetectAndCreateSchemasFromFile attempts to detect and create the schema for
// the given file.
func (c *Client) DetectAndCreateSchemasFromFile(ctx context.Context, account, name, filename string, options ...RequestOption) ([]*schema.Schema, error) {
	s, err := c.DetectSchemasFromFile(ctx, name, filename, options...)

	if err != nil {
		return nil, err
	}

	err = c.createSchemas(ctx, account, s, options...)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// DetectAndCreateSchemas attempts to detect and create the schema for a reader.
func (c *Client) DetectAndCreateSchemas(ctx context.Context, account, name, contentType string, r io.Reader, options ...RequestOption) ([]*schema.Schema, error) {
	s, err := c.DetectSchemas(ctx, name, contentType, r, options...)

	if err != nil {
		return nil, err
	}

	err = c.createSchemas(ctx, account, s, options...)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// DeleteSchema deletes a schema.
func (c *Client) DeleteSchema(ctx context.Context, account, name string, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodDelete,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(name),
		0,
		nil,
		nil,
		options...,
	)

	return err
}
