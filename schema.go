package instapi

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/instapi/client-go/schema"
	"github.com/instapi/client-go/types"
)

// GetSchema gets the given schema.
func (c *Client) GetSchema(ctx context.Context, name string, options ...RequestOption) (*schema.Schema, error) {
	var s *schema.Schema
	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"schemas/"+name,
		http.StatusOK,
		nil,
		&s,
		options...,
	)

	return s, err
}

// GetSchemas gets the schema collection.
func (c *Client) GetSchemas(ctx context.Context, options ...RequestOption) ([]*schema.Schema, string, error) {
	var s []*schema.Schema
	resp, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"schemas",
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

// DetectSchemaFromFile attempts to detect the schema for the given file.
func (c *Client) DetectSchemaFromFile(ctx context.Context, name, filename string, options ...RequestOption) (*schema.Schema, error) {
	contentType, err := getContentType(filename)

	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer f.Close() // nolint: errcheck

	return c.DetectSchema(ctx, name, contentType, f, options...)
}

// DetectSchema attempts to detect the schema for the given reader.
func (c *Client) DetectSchema(ctx context.Context, name, contentType string, r io.Reader, options ...RequestOption) (*schema.Schema, error) {
	const fileSizeLimit = 16 * 1024 * 1024

	var s *schema.Schema
	_, err := c.doRequest(
		ctx,
		http.MethodPost,
		contentType,
		c.endpoint+"schemas/detect",
		http.StatusOK,
		io.LimitReader(r, fileSizeLimit),
		&s,
		append(options, Name(name))...,
	)

	return s, err
}

// CreateSchema creates a new schema.
func (c *Client) CreateSchema(ctx context.Context, s *schema.Schema, options ...RequestOption) error {
	_, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"schemas",
		http.StatusCreated,
		s,
		nil,
		options...,
	)

	return err
}

// DetectAndCreateSchemaFromFile attempts to detect and create the schema for
// the given file.
func (c *Client) DetectAndCreateSchemaFromFile(ctx context.Context, name, filename string, options ...RequestOption) (*schema.Schema, error) {
	s, err := c.DetectSchemaFromFile(ctx, name, filename, options...)

	if err != nil {
		return nil, err
	}

	err = c.CreateSchema(ctx, s, options...)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// DetectAndCreateSchema attempts to detect and create the schema for a reader.
func (c *Client) DetectAndCreateSchema(ctx context.Context, name, contentType string, r io.Reader, options ...RequestOption) (*schema.Schema, error) {
	s, err := c.DetectSchema(ctx, name, contentType, r, options...)

	if err != nil {
		return nil, err
	}

	err = c.CreateSchema(ctx, s, options...)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// DeleteSchema deletes a schema.
func (c *Client) DeleteSchema(ctx context.Context, name string, options ...RequestOption) error {
	_, err := c.doRequest(
		ctx,
		http.MethodDelete,
		types.JSON,
		c.endpoint+"schemas/"+url.PathEscape(name),
		0,
		nil,
		nil,
		options...,
	)

	return err
}
