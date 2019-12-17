package instapi

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/google/go-querystring/query"
	"github.com/instapi/client-go/schema"
	"github.com/instapi/client-go/types"
)

// GetSchema gets the given schema.
func (c *Client) GetSchema(ctx context.Context, name string) (*schema.Schema, error) {
	var s *schema.Schema

	return s, c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"schemas/"+name,
		http.StatusOK,
		nil,
		nil,
		&s,
	)
}

// DetectSchemaForFile attempts to detect the schema for the given file.
func (c *Client) DetectSchemaFromFile(ctx context.Context, options *Options, filename string) (*schema.Schema, error) {
	contentType, err := getContentType(filename)

	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer f.Close() // nolint: errcheck

	return c.DetectSchema(ctx, options, contentType, f)
}

// DetectSchema attempts to detect the schema for the given reader.
func (c *Client) DetectSchema(ctx context.Context, options *Options, contentType string, r io.Reader) (*schema.Schema, error) {
	query, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	var s *schema.Schema

	return s, c.doRequest(
		ctx,
		http.MethodPost,
		contentType,
		c.endpoint+"schemas/detect",
		http.StatusOK,
		query,
		r,
		&s,
	)
}

// CreateSchema creates a new schema.
func (c *Client) CreateSchema(ctx context.Context, s *schema.Schema) error {
	return c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"schemas",
		http.StatusCreated,
		nil,
		s,
		nil,
	)
}

// DetectAndCreateSchemaForFile attempts to detect and create
// the schema for the given file.
func (c *Client) DetectAndCreateSchemaFromFile(ctx context.Context, options *Options, filename string) (*schema.Schema, error) {
	s, err := c.DetectSchemaFromFile(ctx, options, filename)

	if err != nil {
		return nil, err
	}

	err = c.CreateSchema(ctx, s)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// DetectAndCreateSchema attempts to detect and create
// the schema for the given reader.
func (c *Client) DetectAndCreateSchema(ctx context.Context, options *Options, contentType string, r io.Reader) (*schema.Schema, error) {
	s, err := c.DetectSchema(ctx, options, contentType, r)

	if err != nil {
		return nil, err
	}

	err = c.CreateSchema(ctx, s)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// DeleteSchema deletes a schema.
func (c *Client) DeleteSchema(ctx context.Context, name string) error {
	return c.doRequest(
		ctx,
		http.MethodDelete,
		types.JSON,
		c.endpoint+"schemas/"+url.PathEscape(name),
		http.StatusAccepted,
		nil,
		nil,
		nil,
	)
}
