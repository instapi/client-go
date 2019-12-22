package instapi

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
)

// CreateRecords makes a create records request.
func (c *Client) CreateRecords(ctx context.Context, schema, contentType string, r io.Reader, options ...RequestOption) (int, error) {
	var resp struct {
		Count int `json:"count"`
	}
	_, err := c.doRequest(
		ctx,
		http.MethodPost,
		contentType,
		c.endpoint+"schemas/"+url.PathEscape(schema)+"/records",
		http.StatusCreated,
		r,
		&resp,
		options...,
	)

	return resp.Count, err
}

// CreateRecordsFromFile makes a create records request for the given file.
func (c *Client) CreateRecordsFromFile(ctx context.Context, schema, filename string, options ...RequestOption) (int, error) {
	contentType, err := getContentType(filename)

	if err != nil {
		return 0, err
	}

	f, err := os.Open(filename)

	if err != nil {
		return 0, err
	}

	defer f.Close() // nolint: errcheck

	return c.CreateRecords(ctx, schema, contentType, f, options...)
}
